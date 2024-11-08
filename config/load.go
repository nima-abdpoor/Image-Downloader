package config

import (
	"fmt"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"strings"
)

func Load(configPath string) *Config {
	var k = koanf.New(".")

	err := k.Load(file.Provider(configPath), yaml.Parser())
	if err != nil {
		fmt.Printf("Cant load Config from yaml file %v", err)
	}

	err = k.Load(env.Provider("GID_", ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, "GID_")), "_", ".", -1)
	}), nil)
	if err != nil {
		fmt.Printf("Cant load Config from env %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
