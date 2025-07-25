package tracex

// TraceName represents the tracing name.
const TraceName = "HelpStudent"

// A Config is an opentelemetry config.
type Config struct {
	Name     string  `yaml:"Name"`
	Endpoint string  `yaml:"Endpoint"`
	Auth     string  `yaml:"Auth"`
	Sampler  float64 `yaml:"Sampler"`
	Proto    string  `yaml:"Proto"`
	UrlPath  string  `yaml:"UrlPath"`
}
