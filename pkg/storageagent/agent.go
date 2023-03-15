package storageagent

import (
	"fmt"
	"io"
)

type Agent struct {
	clientMap map[ClientType]IClient
}

func (t *Agent) Initialize(clientInfoList []ClientInfo) error {
	t.clientMap = make(map[ClientType]IClient)

	for _, info := range clientInfoList {
		newFunc, exist := registerClientMap[info.ClientType]
		if !exist {
			return fmt.Errorf("type %s dose not exist", info.ClientType)
		}

		client := newFunc()
		if err := client.Initialize(info.Args...); err != nil {
			return fmt.Errorf("failed to initialize %s client, %v",
				info.ClientType, err)
		}

		t.clientMap[info.ClientType] = client
	}

	return nil
}

func (t *Agent) Authenticate() error {
	for clientType, client := range t.clientMap {
		if err := client.Authenticate(); err != nil {
			return fmt.Errorf("the %s client has failed to authenticate, %v", clientType, err)
		}
	}
	return nil
}

func (t *Agent) Upload(clientType ClientType, bucket, key string, reader io.Reader) error {
	client := t.clientMap[clientType]
	if client == nil {
		return fmt.Errorf("type %s dose not exist", clientType)
	}

	return client.Upload(bucket, key, reader)
}
