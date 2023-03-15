package configmodel

type ConfigAWSS3Model struct {
	Fqdn            string `json:"fqdn"`
	Port            int    `json:"port"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

type ConfigStorageAgentModel struct {
	ImportClientTypes []string         `json:"importClientTypes"`
	AmazonS3          ConfigAWSS3Model `json:"s3Storage"`
}

type ConfigStorageModel struct {
	StorageAgent ConfigStorageAgentModel `json:"storageAgent"`
}
