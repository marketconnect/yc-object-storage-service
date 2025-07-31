package api

import (
	"archive/zip"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/marketconnect/yc-object-storage-service/s3"

	"github.com/gin-gonic/gin"
)

type ArchiveRequest struct {
	Keys    []string `json:"keys"`
	Folders []string `json:"folders"`
}

type Handler struct {
	S3Client *s3.Client
}

func (h *Handler) ListObjectsHandler(c *gin.Context) {
	prefix := c.Query("prefix")
	delimiter := c.Query("delimiter")

	output, err := h.S3Client.ListObjects(prefix, delimiter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list objects", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"folders": output.Folders, "files": output.Files})
}

func (h *Handler) ListAllFoldersHandler(c *gin.Context) {
	folders, err := h.S3Client.ListAllFolders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list all folders", "details": err.Error()})
		return
	}
	// Sort for consistent ordering in the UI
	sort.Strings(folders)

	c.JSON(http.StatusOK, gin.H{"folders": folders})
}

func (h *Handler) GeneratePresignedURLHandler(c *gin.Context) {
	objectKey := c.Query("objectKey")
	if objectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "objectKey query parameter is required"})
		return
	}

	expiresStr := c.DefaultQuery("expires", "3600") // Default to 1 hour
	expires, err := strconv.ParseInt(expiresStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid expires parameter"})
		return
	}

	lifetime := time.Duration(expires) * time.Second

	url, err := h.S3Client.GeneratePresignedURL(objectKey, lifetime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate presigned URL", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"url": url})
}

func (h *Handler) CreateArchiveHandler(c *gin.Context) {
	var req ArchiveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body", "details": err.Error()})
		return
	}

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename=archive.zip")

	zipWriter := zip.NewWriter(c.Writer)
	defer zipWriter.Close()

	allKeys := make(map[string]struct{})

	// Add individual files
	for _, key := range req.Keys {
		allKeys[key] = struct{}{}
	}

	// Add files from folders
	for _, folderPrefix := range req.Folders {
		files, err := h.S3Client.ListAllObjects(folderPrefix)
		if err != nil {
			// Log the error but try to continue
			log.Printf("Error listing objects for prefix %s: %v", folderPrefix, err)
			continue
		}
		for _, file := range files {
			allKeys[file] = struct{}{}
		}
	}

	// Process all unique keys
	for key := range allKeys {
		obj, err := h.S3Client.GetObject(key)
		if err != nil {
			log.Printf("Error getting object %s: %v", key, err)
			continue
		}

		f, err := zipWriter.Create(key)
		if err != nil {
			obj.Body.Close()
			log.Printf("Error creating zip entry for %s: %v", key, err)
			continue
		}

		if _, err := io.Copy(f, obj.Body); err != nil {
			obj.Body.Close()
			log.Printf("Error copying object body for %s: %v", key, err)
			continue
		}
		obj.Body.Close()
	}
}
