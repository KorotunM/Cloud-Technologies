package storage

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"pragma/internal/config"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client *minio.Client
	cfg    config.StorageConfig
}

func NewClient(cfg config.StorageConfig) (*Client, error) {
	if cfg.Bucket == "" || cfg.Endpoint == "" {
		// Storage is not configured; fall back to placeholders.
		return &Client{cfg: cfg}, nil
	}

	var creds *credentials.Credentials

	if cfg.AccessKeyID != "" && cfg.SecretAccessKey != "" {
		creds = credentials.NewStaticV4(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		)
	} else {
		// anonymous access (public bucket)
		creds = nil
	}

	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  creds,
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("init object storage client: %w", err)
	}

	return &Client{client: client, cfg: cfg}, nil
}

// ImageURL returns a public or presigned URL for the object key or a placeholder.
func (c *Client) ImageURL(ctx context.Context, objectKey string) (string, error) {
	key := strings.TrimPrefix(objectKey, "/")
	if key == "" {
		return c.placeholder(), nil
	}

	if c.cfg.CDNBaseURL != "" {
		return strings.TrimSuffix(c.cfg.CDNBaseURL, "/") + "/" + key, nil
	}

	if c.client == nil || c.cfg.Bucket == "" {
		return c.placeholder(), nil
	}

	presigned, err := c.client.PresignedGetObject(ctx, c.cfg.Bucket, key, c.cfg.SignedURLTTL, url.Values{})
	if err != nil {
		if fallback := c.placeholder(); fallback != "" {
			return fallback, nil
		}
		return "", fmt.Errorf("presign object %s: %w", key, err)
	}

	return presigned.String(), nil
}

func (c *Client) placeholder() string {
	if c.cfg.DefaultImageURL != "" {
		return c.cfg.DefaultImageURL
	}
	// Local static fallback to keep UI functional even without storage.
	return "/assets/logo.jpg"
}
