package fileServer

import "github.com/pkg/errors"

var clients = make(map[string]FileClient)
var creatorMap = map[string]func(Config) (FileClient, error){
	"mock": NewMock,
	"oss":  NewAliOSS,
}

func InitFileServers(cfgs []Config) error {
	for _, cfg := range cfgs {
		if _, ok := creatorMap[cfg.StorageType]; !ok {
			return errors.New("unknown storage type: " + cfg.StorageType)
		}
		if _, ok := clients[cfg.Key]; ok {
			return errors.New("duplicate storage Key: " + cfg.Key)
		}
		client, err := creatorMap[cfg.StorageType](cfg)
		if err != nil {
			return errors.Wrap(err, "create file client failed")
		}
		clients[cfg.Key] = client
		print("init file server: ", cfg.Key, " ", cfg.StorageType)
	}
	if len(clients) == 0 {
		clients["mock"], _ = NewMock(Config{})
	}
	return nil
}

func Client(key string) FileClient {
	if client, ok := clients[key]; ok {
		return client
	}
	return clients["mock"]
}
