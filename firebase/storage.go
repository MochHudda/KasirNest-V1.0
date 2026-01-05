package firebase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
)

// StorageService handles Firebase Storage operations
type StorageService struct {
	client *storage.Client
	bucket string
	ctx    context.Context
}

// NewStorageService creates a new storage service
func NewStorageService(client *Client, bucketName string) *StorageService {
	return &StorageService{
		client: client.Storage,
		bucket: bucketName,
		ctx:    client.GetContext(),
	}
}

// UploadFile uploads a file to Firebase Storage
func (s *StorageService) UploadFile(fileName string, data io.Reader, contentType string) (string, error) {
	if s.client == nil {
		return "", errors.New("storage client not initialized")
	}

	// Create a bucket handle
	bucket := s.client.Bucket(s.bucket)

	// Create an object handle
	obj := bucket.Object(fileName)

	// Create a writer
	writer := obj.NewWriter(s.ctx)
	writer.ContentType = contentType

	// Copy the file data to storage
	if _, err := io.Copy(writer, data); err != nil {
		writer.Close()
		return "", err
	}

	// Close the writer
	if err := writer.Close(); err != nil {
		return "", err
	}

	// Generate download URL
	downloadURL, err := s.GetDownloadURL(fileName)
	if err != nil {
		return "", err
	}

	return downloadURL, nil
}

// GetDownloadURL generates a signed download URL for a file
func (s *StorageService) GetDownloadURL(fileName string) (string, error) {
	if s.client == nil {
		return "", errors.New("storage client not initialized")
	}

	// Generate signed URL (valid for 7 days)
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  "GET",
		Expires: time.Now().Add(7 * 24 * time.Hour),
	}

	return s.client.Bucket(s.bucket).SignedURL(fileName, opts)
}

// DeleteFile deletes a file from Firebase Storage
func (s *StorageService) DeleteFile(fileName string) error {
	if s.client == nil {
		return errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(fileName)

	return obj.Delete(s.ctx)
}

// ListFiles lists files in the storage bucket with optional prefix
func (s *StorageService) ListFiles(prefix string) ([]string, error) {
	if s.client == nil {
		return nil, errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)

	query := &storage.Query{Prefix: prefix}
	iter := bucket.Objects(s.ctx, query)

	var files []string
	for {
		attrs, err := iter.Next()
		if err == storage.ErrObjectNotExist {
			break
		}
		if err != nil {
			return nil, err
		}
		files = append(files, attrs.Name)
	}

	return files, nil
}

// GetFileInfo gets metadata information about a file
func (s *StorageService) GetFileInfo(fileName string) (*FileInfo, error) {
	if s.client == nil {
		return nil, errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(fileName)

	attrs, err := obj.Attrs(s.ctx)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		Name:        attrs.Name,
		Size:        attrs.Size,
		ContentType: attrs.ContentType,
		Created:     attrs.Created,
		Updated:     attrs.Updated,
		MD5:         fmt.Sprintf("%x", attrs.MD5),
	}, nil
}

// DownloadFile downloads a file from Firebase Storage
func (s *StorageService) DownloadFile(fileName string) (io.ReadCloser, error) {
	if s.client == nil {
		return nil, errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(fileName)

	return obj.NewReader(s.ctx)
}

// CopyFile copies a file within the storage bucket
func (s *StorageService) CopyFile(srcFileName, destFileName string) error {
	if s.client == nil {
		return errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)
	src := bucket.Object(srcFileName)
	dst := bucket.Object(destFileName)

	_, err := dst.CopierFrom(src).Run(s.ctx)
	return err
}

// MoveFile moves a file within the storage bucket
func (s *StorageService) MoveFile(srcFileName, destFileName string) error {
	// Copy file to new location
	if err := s.CopyFile(srcFileName, destFileName); err != nil {
		return err
	}

	// Delete original file
	return s.DeleteFile(srcFileName)
}

// FileExists checks if a file exists in storage
func (s *StorageService) FileExists(fileName string) (bool, error) {
	if s.client == nil {
		return false, errors.New("storage client not initialized")
	}

	bucket := s.client.Bucket(s.bucket)
	obj := bucket.Object(fileName)

	_, err := obj.Attrs(s.ctx)
	if err == storage.ErrObjectNotExist {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// FileInfo represents file metadata
type FileInfo struct {
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	MD5         string    `json:"md5"`
}

// UploadProductImage uploads a product image with proper naming
func (s *StorageService) UploadProductImage(productID string, data io.Reader, contentType string) (string, error) {
	fileName := fmt.Sprintf("products/%s/%d.jpg", productID, time.Now().Unix())
	return s.UploadFile(fileName, data, contentType)
}

// DeleteProductImage deletes a product image
func (s *StorageService) DeleteProductImage(productID string) error {
	// List files with product prefix
	files, err := s.ListFiles(fmt.Sprintf("products/%s/", productID))
	if err != nil {
		return err
	}

	// Delete all files for this product
	for _, file := range files {
		if err := s.DeleteFile(file); err != nil {
			return err
		}
	}

	return nil
}
