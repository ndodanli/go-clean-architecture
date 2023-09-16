package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func exportConfig() error {
	viper.SetConfigType("yaml")
	var configPath string
	allEnvironments := os.Environ()
	fmt.Println(allEnvironments)
	switch os.Getenv("APP_ENV") {
	case "test":
		configPath = "../environments/"
		viper.SetConfigName("test")
	case "dev":
		configPath = "./environments/"
		viper.SetConfigName("dev")
	case "prod":
		configPath = "./environments/"
		viper.SetConfigName("prod")
	default:
		configPath = "./environments/"
		viper.SetConfigName("dev")
	}

	viper.AddConfigPath(configPath)

	if err := viper.ReadInConfig(); err != nil {
		return err
	}
	return nil
}

// ParseConfig Parse config file
func ParseConfig() (*Config, error) {
	if err := exportConfig(); err != nil {
		return nil, err
	}

	var c Config
	err := viper.Unmarshal(&c)
	if err != nil {
		log.Printf("unable to decode into struct, %v", err)
		return nil, err
	}

	return &c, nil
}
