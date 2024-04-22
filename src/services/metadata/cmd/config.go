package main

type Config struct {
	Api ApiConfig `yaml:"api"`
}

type ApiConfig struct {
	Port int `yaml:"port"`
}
