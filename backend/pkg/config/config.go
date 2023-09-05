package config

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
)

func (cfg *Config) ParseData(data []byte) error {
	return yaml.Unmarshal([]byte(data), cfg)
}

func (cfg *Config) ParseFile() error {
	data, err := ioutil.ReadFile(cfg.File)
	if err != nil {
		return err
	}

	return cfg.ParseData(data)
}

func (cfg *Config) WriteFile() error {
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(cfg.File, data, 0640)
	if err != nil {
		return err
	}

	return nil
}
