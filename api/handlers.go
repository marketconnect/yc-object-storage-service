package api

import (
	"net/http"
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
	if prefix == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prefix query parameter is required"})
		return
	}

	keys, err := h.S3Client.ListObjects(prefix)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list objects", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"objects": keys})
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
