package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
)

func (t *AmazonS3Client) Upload(bucket, key string, reader io.Reader) error {
	client := s3.NewFromConfig(*t.cfg)
	output, putErr := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})

	if putErr != nil {
		return putErr
	}

	output = output
	return nil
}

func (t *AmazonS3Client) UploadLarge(bucket, key string, reader io.Reader, partSize int64) error {
	client := s3.NewFromConfig(*t.cfg)
	uploader := manager.NewUploader(client, func(u *manager.Uploader) {
		//u.PartSize = partKBs * 1024
		u.PartSize = partSize
	})

	_, putErr := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   reader,
	})

	if putErr != nil {
		return putErr
	}

	//output = output
	return nil
}
