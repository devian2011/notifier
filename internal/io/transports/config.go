package transports

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Smtp map[string]SmtpConfig `yaml:"smtp" json:"smtp"`
	File map[string]FileConfig `yaml:"file" json:"file"`
}

func loadConfigFromFile(filePath string) (*Config, error) {
	file, openFileErr := os.OpenFile(filePath, os.O_RDONLY, 0755)
	if openFileErr != nil {
		return nil, openFileErr
	}
	cfg := &Config{}

	decodeErr := yaml.NewDecoder(file).Decode(cfg)

	return cfg, decodeErr
}
