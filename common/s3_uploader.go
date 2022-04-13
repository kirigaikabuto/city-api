package common

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	se "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"strings"
)

type FileInfo struct {
	FileUrl string `json:"file_url,omitempty"`
}

type Avatar struct {
	Original  string `json:"original"`
	Normal    string `json:"normal"`
	Thumbnail string `json:"thumbnail"`
}

type S3Uploader interface {
	UploadFile(file []byte, key, fileType, contentType string) (*FileInfo, error)
	DeleteFile(key, fileType string) error
	FileExist(key, fileType string) (bool, error)
	ListFile() []string
	GetFile(key, fileType string) (string, error)
}

type defaultS3Uploader struct {
	Bucket         string
	UploadFilesURI string
	Uploader       *s3.S3
}

func NewS3Uploader(endpoint, accessKey, secretKey, bucket, uploadedFilesURI, region string) (S3Uploader, error) {
	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Endpoint:         aws.String(endpoint),
		Region:           aws.String(region),
		DisableSSL:       aws.Bool(true),
		S3ForcePathStyle: aws.Bool(true),
	}
	newSession, err := se.NewSession(s3Config)
	if err != nil {
		return nil, err
	}
	s3Client := s3.New(newSession)

	return &defaultS3Uploader{
		Bucket:         bucket,
		UploadFilesURI: uploadedFilesURI,
		Uploader:       s3Client,
	}, nil
}

func (updr *defaultS3Uploader) UploadFile(file []byte, key, fileType, contentType string) (*FileInfo, error) {
	//var contentType string
	//if fileType == "png" {
	//	contentType = "image/png"
	//} else {
	//	contentType = "applications/octet-stream"
	//}
	_, err := updr.Uploader.PutObject(&s3.PutObjectInput{
		Body:        bytes.NewReader(file),
		Bucket:      aws.String(updr.Bucket),
		Key:         aws.String(key + "_original." + fileType),
		ACL:         aws.String("public-read"),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return nil, err
	}

	fileURI := updr.UploadFilesURI + "/"
	return &FileInfo{
		fileURI + key + "_original." + fileType,
	}, nil
}

func (updr *defaultS3Uploader) DeleteFile(key, fileType string) error {
	filePath := key + "_original." + fileType
	_, err := updr.Uploader.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(updr.Bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		return err
	}
	return nil
}

func (updr *defaultS3Uploader) FileExist(key, fileType string) (bool, error) {
	filePath := key + "_original." + fileType
	_, err := updr.Uploader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(updr.Bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey:") {
			return false, nil
		} else {
			return false, err
		}
	}
	return true, nil
}

func (updr *defaultS3Uploader) GetFile(key, fileType string) (string, error) {
	filePath := key + "_original." + fileType
	_, err := updr.Uploader.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(updr.Bucket),
		Key:    aws.String(filePath),
	})
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchKey:") {
			return "", nil
		} else {
			return "", err
		}
	}
	return updr.UploadFilesURI + "/" + filePath, nil
}

func (updr *defaultS3Uploader) ListFile() []string {
	ctx := context.Background()
	objects := []string{}
	err := updr.Uploader.ListObjectsPagesWithContext(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(updr.Bucket),
	}, func(p *s3.ListObjectsOutput, lastPage bool) bool {
		for _, o := range p.Contents {
			objects = append(objects, aws.StringValue(o.Key))
		}
		return true // continue paging
	})
	if err != nil {
		panic(fmt.Sprintf("failed to list objects for bucket, %s, %v", updr.Bucket, err))
	}

	return objects
}

var magicTable = map[string]string{
	"\xff\xd8\xff":      "image/jpeg",
	"\x89PNG\r\n\x1a\n": "image/png",
	"GIF87a":            "image/gif",
	"GIF89a":            "image/gif",
}

var imageTable = map[string]string{
	"image/jpeg": ".jpeg",
	"image/png":  ".png",
	"image/gif":  ".gif",
}

func isImageType(in string) bool {
	for _, v := range magicTable {
		if v == in {
			return true
		}
	}
	return false
}
