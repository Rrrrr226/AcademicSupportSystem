package logx

// A LogConf is a logging config.
type LogConf struct {
	EnableSls   bool          `yaml:"EnableSls"`
	SlsSinkConf []SlsSinkConf `yaml:"SlsSinkConf"`
}

type SlsSinkConf struct {
	Url string `yaml:"Url"`
}
