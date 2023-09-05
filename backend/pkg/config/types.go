package config

import "time"

type Config struct {
	File string

	ErrorLevel  string   `yaml:"error-level"`
	Host        string   `yaml:"host"`
	AllowOrigin []string `yaml:"allow-origin"`
	StaticDir   string   `yaml:"static-dir"`
	WebRoot     string   `yaml:"web-root"`

	Mysql struct {
		Master          string        `yaml:"master"`
		Sources         []string      `yaml:"sources"`
		Replicas        []string      `yaml:"replicas"`
		MaxIdleConns    int           `yam:"set-max-idle-conns"`
		MaxOpenConns    int           `yam:"set-max-open-conns"`
		ConnMaxLifetime time.Duration `yaml:"set-conn-max-lifetime"`
		AllowTable      []string      `yaml:"allow-table"`
	} `yaml:"mysql"`
}

func New(file string) (*Config, error) {
	cfg := &Config{File: file}
	err := cfg.ParseFile()
	return cfg, err
}
