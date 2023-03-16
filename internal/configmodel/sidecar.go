package configmodel

type SidecarConfigModel struct {
	S3Storage ConfigAWSS3Model `json:"s3Storage"`
}

type ConfigAWSS3Model struct {
	Endpoint string `json:"endpoint"`
	//Region          string `json:"region"`
	AccessKeyID     string `json:"accessKeyId"`
	SecretAccessKey string `json:"secretAccessKey"`
}

func GetSidecarConfig() *SidecarConfigModel {
	return confInst.sidecarConfigModel
}
