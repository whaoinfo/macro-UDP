package storageagent

import (
	"github.com/whaoinfo/macro-UDP/pkg/storageagent/s3"
	"io"
)

const (
	AWSS3ClientType = "s3"
	SimClientType   = "sim"
)

type ClientType string

type NewClientFunc func() IClient

var (
	registerClientMap = map[ClientType]NewClientFunc{
		AWSS3ClientType: func() IClient {
			return s3.NewAmazonS3Client()
		},
	}
)

type IClient interface {
	Initialize(args ...interface{}) error
	Authenticate() error
	Upload(bucket, key string, reader io.Reader) error
	UploadLarge(bucket, key string, reader io.Reader, PartSize int64) error
}

type ClientInfo struct {
	ClientType ClientType
	Args       []interface{}
}
