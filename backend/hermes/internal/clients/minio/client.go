package minio

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	s3lib "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type Config struct {
	URL       string `json:"S3_URL"`
	AccessKey string `json:"S3_ACCESS_KEY"`
	SecretKey string `json:"S3_SECRET_KEY"`
	Bucket    string `json:"S3_BUCKET"`
	ChunkSize int    `json:"S3_CHUNK_SIZE" cfgDefault:"5242880"`
}

func NewClient(cfg Config) (*Client, error) {
	awsCfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("eu-north-1"),
		config.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(
				cfg.AccessKey, cfg.SecretKey, "",
			),
		),
		config.WithEndpointResolver(
			aws.EndpointResolverFunc(
				func(service, region string) (aws.Endpoint, error) {
					return aws.Endpoint{
						URL:               cfg.URL,
						HostnameImmutable: true,
					}, nil
				},
			),
		),
	)
	if err != nil {
		return nil, err
	}

	svc := s3lib.NewFromConfig(awsCfg)

	_, err = svc.HeadBucket(
		context.TODO(), &s3lib.HeadBucketInput{
			Bucket: &cfg.Bucket,
		},
	)
	if err != nil {
		_, err = svc.CreateBucket(
			context.TODO(), &s3lib.CreateBucketInput{
				Bucket: &cfg.Bucket,
				CreateBucketConfiguration: &types.CreateBucketConfiguration{
					LocationConstraint: types.BucketLocationConstraintEuNorth1,
				},
			},
		)
		if err != nil {
			return nil, err
		}
	}

	return &Client{S3: svc, cfg: cfg}, nil
}

type Client struct {
	S3  *s3lib.Client
	cfg Config
}

func (c *Client) Release() error {
	return nil
}

func (c *Client) UploadFile(ctx context.Context, fileName string, data []byte) (string, error) {
	initResp, err := c.S3.CreateMultipartUpload(ctx, &s3lib.CreateMultipartUploadInput{
		Bucket: &c.cfg.Bucket,
		Key:    &fileName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to initiate multipart upload: %w", err)
	}

	uploadID := initResp.UploadId
	var parts []types.CompletedPart
	partNumber := int32(1)
	buffer := make([]byte, c.cfg.ChunkSize)

	file := bytes.NewReader(data)

	for {

		n, err := file.Read(buffer)
		if err != nil && err.Error() != "EOF" {
			return "", fmt.Errorf("failed to read file: %w", err)
		}
		if n == 0 {
			break
		}

		uploadResp, err := c.S3.UploadPart(ctx, &s3lib.UploadPartInput{
			Bucket:     &c.cfg.Bucket,
			Key:        &fileName,
			UploadId:   uploadID,
			PartNumber: &partNumber,
			Body:       bytes.NewReader(buffer[:n]),
		})
		if err != nil {
			c.S3.AbortMultipartUpload(ctx, &s3lib.AbortMultipartUploadInput{
				Bucket:   &c.cfg.Bucket,
				Key:      &fileName,
				UploadId: uploadID,
			})
			return "", fmt.Errorf("failed to upload part %d: %w", partNumber, err)
		}

		parts = append(parts, types.CompletedPart{
			ETag:       uploadResp.ETag,
			PartNumber: aws.Int32(partNumber),
		})

		partNumber++
	}

	_, err = c.S3.CompleteMultipartUpload(ctx, &s3lib.CompleteMultipartUploadInput{
		Bucket:          &c.cfg.Bucket,
		Key:             &fileName,
		UploadId:        uploadID,
		MultipartUpload: &types.CompletedMultipartUpload{Parts: parts},
	})
	if err != nil {
		return "", fmt.Errorf("failed to complete multipart upload: %w", err)
	}

	url := fmt.Sprintf("%s/%s/%s", c.cfg.URL, c.cfg.Bucket, fileName)
	return url, nil
}
