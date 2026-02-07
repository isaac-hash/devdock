package cmd

import (
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

type Config struct {
	ProjectName string `yaml:"project_name"`
	Type string `yaml:"type,omitempty"`
}

func LoadConfig() (Config, error) {
	var cfg Config
	data, err := os.ReadFile("devdock.yaml")
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(data, &cfg)
	return cfg, err
}

func WriteConfig(cfg Config) error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile("devdock.yaml", data, 0644)
}

func GetProjectName() string{
	wd, _ := os.Getwd()
	return filepath.Base(wd)
}