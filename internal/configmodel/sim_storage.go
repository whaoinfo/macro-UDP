package configmodel

type SimStorageModel struct {
	Endpoint string `json:"endpoint"`
}

func GetSimStorageConfig() *SimStorageModel {
	return confInst.simStorageModel
}
