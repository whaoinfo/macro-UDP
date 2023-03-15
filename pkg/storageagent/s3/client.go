package s3

import (
	"context"
	"crypto/tls"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"net/http"
)

func NewAmazonS3Client() *AmazonS3Client {
	return &AmazonS3Client{}
}

type AmazonS3Client struct {
	accessKeyID     string
	secretAccessKey string
	region          string // us-east-2
	endpointURL     string // http://172.0.3.45:9000

	cfg *aws.Config
}

func (t *AmazonS3Client) Initialize(args ...interface{}) error {
	t.accessKeyID = args[0].(string)
	t.secretAccessKey = args[1].(string)
	t.region = args[2].(string)
	t.endpointURL = args[3].(string)

	cfg, loadCfgErr := loadConfig(t.region, t.endpointURL)
	if loadCfgErr != nil {
		return loadCfgErr
	}

	t.cfg = cfg
	return nil
}

func (t *AmazonS3Client) Authenticate() error {

	return nil
}

func loadConfig(region, endpointURL string) (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		// config.WithClientLogMode(aws.LogRequestWithBody|aws.LogResponseWithBody),
		config.WithRegion(region),
		config.WithHTTPClient(&http.Client{Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}}),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					//PartitionID: "aws",
					URL: endpointURL,
					//SigningRegion: "us-east-2",
					//HostnameImmutable: true,
				}, nil
			}),
		))
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
