package s3

import (
	"log"
	"strings"
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

type ListObjectsOutput struct {
	Folders []string
	Files   []string
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

func (c *Client) ListObjects(prefix string, delimiter string) (*ListObjectsOutput, error) {
	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(c.BucketName),
		Prefix: aws.String(prefix),
	}
	if delimiter != "" {
		input.Delimiter = aws.String(delimiter)
	}

	result, err := c.S3Svc.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, item := range result.Contents {
		// When listing a "folder", S3 might return the folder key itself.
		// We want to ignore it if it's the same as the prefix we are querying.
		if *item.Key != prefix {
			files = append(files, *item.Key)
		}
	}

	var folders []string
	for _, p := range result.CommonPrefixes {
		if p.Prefix != nil {
			folders = append(folders, *p.Prefix)
		}
	}

	return &ListObjectsOutput{Folders: folders, Files: files}, nil
}

func (c *Client) ListAllFolders() ([]string, error) {
	folderSet := make(map[string]struct{})
	var lastKey *string

	for {
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(c.BucketName),
			ContinuationToken: lastKey,
		}

		result, err := c.S3Svc.ListObjectsV2(input)
		if err != nil {
			return nil, err
		}

		for _, item := range result.Contents {
			if strings.Contains(*item.Key, "/") {
				pathParts := strings.Split(*item.Key, "/")
				// Iterate through path parts to add all parent directories
				for i := 1; i < len(pathParts); i++ {
					folderSet[strings.Join(pathParts[:i], "/")+"/"] = struct{}{}
				}
			}
		}

		if !*result.IsTruncated {
			break
		}
		lastKey = result.NextContinuationToken
	}

	folders := make([]string, 0, len(folderSet))
	for folder := range folderSet {
		folders = append(folders, folder)
	}

	return folders, nil
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
