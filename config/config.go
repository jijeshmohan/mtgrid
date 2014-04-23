package config

import (
	"fmt"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   string
	Database string
	Dbpath   string
}

func NewConfig(env string) (*Config, error) {
	var config Config
	if env == "" {
		env = "development"
	}
	filename := fmt.Sprintf("./%s.conf", env)
	if _, err := toml.DecodeFile(filename, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
