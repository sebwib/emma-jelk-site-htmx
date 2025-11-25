package services

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
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
	"github.com/nfnt/resize"
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

func resizeImage(data []byte, size uint) ([]byte, error) {
	log.Println("trying to decode")
	originalImage, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("decode error %w", err)
	}

	newImage := resize.Resize(size, 0, originalImage, resize.Lanczos3)

	var outData bytes.Buffer
	err = jpeg.Encode(&outData, newImage, nil)
	log.Println("encoding")

	if err != nil {
		return nil, err
	}

	return outData.Bytes(), nil
}

func (u *ImageUploader) UploadImage(file multipart.File, header *multipart.FileHeader) (string, string, error) {
	// Generate unique filename
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%s-%d%s", uuid.New().String(), time.Now().Unix(), ext)

	log.Println("read bytes")
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", err
	}

	log.Println("resizing")
	thumbFileBytes, err := resizeImage(fileBytes, 600)
	if err != nil {
		return "", "", err
	}

	thumbFileName := fmt.Sprintf("thumb-%s", filename)
	// replace thumb file extension with .jpg
	thumbFileName = strings.TrimSuffix(thumbFileName, filepath.Ext(thumbFileName)) + ".jpg"

	if u.useS3 {
		thumbUploadedFileName, err := u.uploadToS3(thumbFileBytes, thumbFileName, header)
		if err != nil {
			return "", "", err
		}

		uploadedFileName, err := u.uploadToS3(fileBytes, filename, header)
		if err != nil {
			return "", "", err
		}
		return uploadedFileName, thumbUploadedFileName, nil
	}

	thumbUploadedFileName, err := u.uploadToLocal(thumbFileBytes, thumbFileName)
	if err != nil {
		return "", "", err
	}
	uploadedFileName, err := u.uploadToLocal(fileBytes, filename)
	if err != nil {
		return "", "", err
	}

	return uploadedFileName, thumbUploadedFileName, nil
}

func (u *ImageUploader) uploadToS3(fileBytes []byte, filename string, header *multipart.FileHeader) (string, error) {
	// Determine content type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	log.Println("uploading: ", filename)

	_, err := u.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(u.bucketName),
		Key:         aws.String(filename),
		Body:        bytes.NewReader(fileBytes),
		ContentType: aws.String(contentType),
		//ACL:         "public-read",
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Return the S3 URL
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.bucketName, filename), nil
}

func (u *ImageUploader) uploadToLocal(fileBytes []byte, filename string) (string, error) {
	log.Println("uploadToLocal: ", filename)

	// Create destination file
	dstPath := filepath.Join(u.localPath, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", fmt.Errorf("failed to create local file: %w", err)
	}
	defer dst.Close()

	// Copy file contents
	if _, err := dst.Write(fileBytes); err != nil {
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
