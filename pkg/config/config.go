package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Host    string    `yaml:"host"`
	Port    int       `yaml:"port"`
	Plugins yaml.Node `yaml:"plugins,omitempty"`
}

func ParseFile(f string) (*Config, error) {
	b, err := os.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

func Parse(b []byte) (*Config, error) {
	config := &Config{}
	err := yaml.Unmarshal(b, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}
