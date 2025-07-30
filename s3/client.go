package s3

import (
	"log"
	"time"

	"github.com/marketconnect/yc-object-storage-service/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type Client struct {
	S3Svc      *s3.S3
	BucketName string
}

func NewClient(cfg *config.Config) *Client {
	awsCfg := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(cfg.S3AccessKeyID, cfg.S3SecretAccessKey, ""),
		Endpoint:         aws.String(cfg.S3Endpoint),
		Region:           aws.String(cfg.S3Region),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		log.Fatalf("Failed to create S3 session: %v", err)
	}

	return &Client{
		S3Svc:      s3.New(sess),
		BucketName: cfg.S3BucketName,
	}
}

func (c *Client) ListObjects(prefix string) ([]string, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.BucketName),
		Prefix: aws.String(prefix),
	}

	result, err := c.S3Svc.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}

	var objectKeys []string
	for _, item := range result.Contents {
		objectKeys = append(objectKeys, *item.Key)
	}

	return objectKeys, nil
}

func (c *Client) GeneratePresignedURL(objectKey string, lifetime time.Duration) (string, error) {
	req, _ := c.S3Svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(c.BucketName),
		Key:    aws.String(objectKey),
	})

	urlStr, err := req.Presign(lifetime)
	if err != nil {
		return "", err
	}

	return urlStr, nil
}
