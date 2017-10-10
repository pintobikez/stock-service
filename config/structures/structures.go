package structures

type DatabaseConfig struct {
	Driver struct {
		Host   string `yaml:"host,omitempty"`
		User   string `yaml:"user,omitempty"`
		Pw     string `yaml:"pw,omitempty"`
		Port   int    `yaml:"port,omitempty"`
		Schema string `yaml:"schema,omitempty"`
	}
}

type PublisherConfig struct {
	Host     string `yaml:"host,omitempty"`
	User     string `yaml:"user,omitempty"`
	Pw       string `yaml:"pw,omitempty"`
	Port     int    `yaml:"port,omitempty"`
	Exchange string `yaml:"exchange,omitempty"`
}
