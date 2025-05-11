package config

import "os"

type Config struct {
	Port string
}

func NewConfig() *Config {
	return &Config{
		Port: getPort(),
	}
}

func getPort() string {
	return getOrDefault("SUBSCRIPTIONS_PORT", DefaultPort)
}

func getOrDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}
