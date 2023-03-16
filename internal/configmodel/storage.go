package configmodel

type ConfigStorageAgentModel struct {
	//ImportClientTypes []string         `json:"importClientTypes"`
	AmazonS3 ConfigAWSS3Model `json:"s3Storage"`
}

type ConfigStorageModel struct {
	StorageAgent ConfigStorageAgentModel `json:"storageAgent"`
}
