package config

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Addr   string `toml:"bind_addr" env:"MTG_BIND_ADDR"`
	DbType string `toml:"db_type" env:"MTG_DB_TYPE"`
	DbPath string `toml:"db_path" env:"MTG_DB_PATH"`
}

func New() *Config {
	c := Config{Addr: ":3000", DbType: "sqlite3", DbPath: "db/development.db"}
	return &c
}

func (c *Config) LoadFile(path string) error {
	_, err := toml.DecodeFile(path, c)
	return err
}

func (c *Config) LoadEnv() error {
	value := reflect.Indirect(reflect.ValueOf(c))
	typ := value.Type()
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Retrieve environment variable.
		v := strings.TrimSpace(os.Getenv(field.Tag.Get("env")))
		if v == "" {
			continue
		}

		// Set the appropriate type.
		switch field.Type.Kind() {
		case reflect.Bool:
			value.Field(i).SetBool(v != "0" && v != "false")
		case reflect.Int:
			newValue, err := strconv.ParseInt(v, 10, 0)
			if err != nil {
				return fmt.Errorf("Parse error: %s: %s", field.Tag.Get("env"), err)
			}
			value.Field(i).SetInt(newValue)
		case reflect.String:
			value.Field(i).SetString(v)
		}
	}
	return nil
}
