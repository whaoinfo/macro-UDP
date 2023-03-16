package configmodel

type ConfigAWSS3Model struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type ConfigStorageAgentModel struct {
	//ImportClientTypes []string         `json:"importClientTypes"`
	AmazonS3 ConfigAWSS3Model `json:"s3Storage"`
}

type ConfigStorageModel struct {
	StorageAgent ConfigStorageAgentModel `json:"storageAgent"`
}
