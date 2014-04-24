package config

import "github.com/BurntSushi/toml"

type Config struct {
	Addr   string `toml:"bind_addr"`
	DbType string `toml:"db_type"`
	DbPath string `toml:"db_path"`
}

func New() *Config {
	c := Config{Addr: ":3000", DbType: "sqlite3", DbPath: "db/development.db"}
	return &c
}

func (c *Config) LoadFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	return err
}
