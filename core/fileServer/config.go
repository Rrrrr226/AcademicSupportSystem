package fileServer

type Config struct {
	Key             string   `yaml:"Key"`
	StorageType     string   `yaml:"StorageType"`
	AccessKeyId     string   `yaml:"AccessKeyId"`
	AccessKeySecret string   `yaml:"AccessKeySecret"`
	EndPoint        string   `yaml:"EndPoint"`
	BucketName      string   `yaml:"BucketName"`
	Schema          string   `yaml:"Schema"`
	Host            string   `yaml:"Host"`
	Prefix          string   `yaml:"Prefix"`
	CallbackUrls    []string `yaml:"CallbackUrls"`
}
