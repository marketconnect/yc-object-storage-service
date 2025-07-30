package api

import (
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/marketconnect/yc-object-storage-service/s3"

	"github.com/gin-gonic/gin"
)

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
