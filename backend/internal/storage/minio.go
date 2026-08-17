package storage

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioClient MinIO 对象存储客户端：用于导出审核报告快照。
type MinioClient struct {
	client *minio.Client
	bucket string
	log    *slog.Logger
}

// NewMinioClient 创建 MinIO 客户端并确保存储桶存在。
func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool, log *slog.Logger) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("init minio client: %w", err)
	}
	m := &MinioClient{client: client, bucket: bucket, log: log}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("check minio bucket: %w", err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("make minio bucket: %w", err)
		}
	}
	return m, nil
}

// UploadJSON 上传 JSON 快照并返回预签名下载地址。
func (m *MinioClient) UploadJSON(ctx context.Context, objectName string, data []byte) (string, error) {
	_, err := m.client.PutObject(ctx, m.bucket, objectName, bytes.NewReader(data), int64(len(data)),
		minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return "", fmt.Errorf("put object: %w", err)
	}
	u, err := m.client.PresignedGetObject(ctx, m.bucket, objectName, 24*time.Hour, nil)
	if err != nil {
		return "", fmt.Errorf("presign object: %w", err)
	}
	return u.String(), nil
}
