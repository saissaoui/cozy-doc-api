package utils

import (
	"github.com/spf13/viper"
)

type AppConfig struct {
	Port            int
	Environment     string
	CouchDbHost     string
	CouchDbPort     int
	CouchDbUser     string
	CouchDbPassword string
}

func LoadConfig() AppConfig {
	viper.AutomaticEnv()
	cfg := AppConfig{
		Port:            getIntWithDefault("PORT", 8088),
		Environment:     getStringWithDefault("ENVIRONMENT", "development"),
		CouchDbHost:     getStringWithDefault("DB_HOST", "http://localhost"),
		CouchDbPort:     getIntWithDefault("DB_PORT", 5984),
		CouchDbUser:     getStringWithDefault("DB_USER", "user"),
		CouchDbPassword: getStringWithDefault("DB_PASSWORD", "password"),
	}

	return cfg
}

func getStringWithDefault(key, defaultValue string) string {
	viper.SetDefault(key, defaultValue)
	return viper.GetString(key)
}

func getIntWithDefault(key string, defaultValue int) int {
	viper.SetDefault(key, defaultValue)
	return viper.GetInt(key)
}
