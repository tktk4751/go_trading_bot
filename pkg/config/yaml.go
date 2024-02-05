package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var data = "config.yml"

type Config struct {
	AssetName string `yaml:"assetName"`
	Dration   string `yaml:"duration"`
	Start     string `yaml:"start"`
	End       string `yaml:"end"`
}

func Yaml() (Config, error) {
	t := Config{}

	// Read the YAML file
	file, errer := os.ReadFile(data)
	if errer != nil {
		return t, errer
	}

	err := yaml.Unmarshal([]byte(file), &t)
	if err != nil {
		return t, err
	}
	return t, nil
}
