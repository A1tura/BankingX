package storage

import (
	"bytes"
	"context"
	"fmt"

	"github.com/minio/minio-go/v7"
)

func UploadDocument(minioClient *minio.Client, documentType string, document []byte, userId int) error {
	found, err := minioClient.BucketExists(context.Background(), "kyc-documents")
	if !found {
		err := minioClient.MakeBucket(context.Background(), "kyc-documents", minio.MakeBucketOptions{
			Region: "us",
		})
		if err != nil {
			return err
		}
	}
	if err != nil {
		return err
	}

	image := bytes.NewReader(document)

	filename := fmt.Sprint(userId, "/", documentType, ".png")
	_, err = minioClient.PutObject(context.Background(), "kyc-documents", filename, image, int64(image.Len()), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}
