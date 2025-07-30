package main

import (
	"fmt"
	"log"

	"github.com/marketconnect/yc-object-storage-service/api"
	"github.com/marketconnect/yc-object-storage-service/config"
	"github.com/marketconnect/yc-object-storage-service/s3"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	s3Client := s3.NewClient(cfg)

	handler := &api.Handler{
		S3Client: s3Client,
	}

	router := gin.Default()

	// Internal API routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/list", handler.ListObjectsHandler)
		v1.GET("/generate-url", handler.GeneratePresignedURLHandler)
		v1.GET("/list-all-folders", handler.ListAllFoldersHandler)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	addr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Starting storage service on %s", addr)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}
