package services

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

type ImageUploader struct {
	s3Client   *s3.Client
	bucketName string
	useS3      bool
	localPath  string
}

func NewImageUploader() (*ImageUploader, error) {
	bucketName := os.Getenv("S3_BUCKET_NAME")
	awsRegion := os.Getenv("AWS_REGION")
	awsAccessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	uploader := &ImageUploader{
		bucketName: bucketName,
		localPath:  "./static/upload",
	}

	// If S3 credentials are provided, use S3, otherwise use local storage
	if bucketName != "" && awsRegion != "" && awsAccessKey != "" && awsSecretKey != "" {
		cfg, err := config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(awsRegion),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
				awsAccessKey,
				awsSecretKey,
				"",
			)),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to load AWS config: %w", err)
		}

		uploader.s3Client = s3.NewFromConfig(cfg)
		uploader.useS3 = true
	}

	// Ensure local directory exists
	if err := os.MkdirAll(uploader.localPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create upload directory: %w", err)
	}

	return uploader, nil
}

func (u *ImageUploader) UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	if u.useS3 {
		return u.uploadToS3(file, filename, header)
	}
	return u.uploadToLocal(file, filename)
}

func (u *ImageUploader) uploadToS3(file multipart.File, filename string, header *multipart.FileHeader) (string, error) {
	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err := u.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(filename),
		Body:        file,
		ContentType: aws.String(contentType),
		ACL:         "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Return the S3 URL
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.bucketName, filename), nil
}

func (u *ImageUploader) uploadToLocal(file multipart.File, filename string) (string, error) {
	// Create destination file
	dstPath := filepath.Join(u.localPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := io.Copy(dst, file); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Return relative URL
	return "/static/upload/" + filename, nil
}

func (u *ImageUploader) DeleteImage(url string) error {
	if u.useS3 && strings.Contains(url, "s3.amazonaws.com") {
		// Extract key from S3 URL
		parts := strings.Split(url, "/")
		if len(parts) == 0 {
			return fmt.Errorf("invalid S3 URL")
		}
		key := parts[len(parts)-1]

		_, err := u.s3Client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
			Bucket: aws.String(u.bucketName),
			Key:    aws.String(key),
		})
		return err
	}

	// Delete from local storage
	if strings.HasPrefix(url, "/static/upload/") {
		filename := strings.TrimPrefix(url, "/static/upload/")
		path := filepath.Join(u.localPath, filename)
		return os.Remove(path)
	}

	return nil
}
