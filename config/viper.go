package config

import (
	Format "fmt"

	Viper "github.com/spf13/viper"
)

func Initialize() *Config {

	Viper.SetConfigName("dyndns.yml")
	Viper.SetConfigType("yaml")
	Viper.AddConfigPath("/var/dyndns/") // path to look for the config file in
	Viper.AddConfigPath(".")            // optionally look for config in the working directory
	err := Viper.ReadInConfig()
	if err != nil {
		panic(Format.Errorf("fatal error reading config file dyndns.yml: %w", err))
	}
	var config Config
	if err := Viper.Unmarshal(&config); err != nil {
		Format.Println(err)
		panic(Format.Errorf("fatal error unmarshalling config file dyndns.yml: %w", err))
	}
	return &config
}

type HetznerConfig struct {
	Dns map[string]string
}

type UserAndPassword struct {
	User     string
	Password string
}

type Config struct {
	Hetzner HetznerConfig               `mapstructure:"hetzner"`
	Users   map[string]*UserAndPassword `mapstructure:"users"`
}
