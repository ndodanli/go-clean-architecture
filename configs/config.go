package configs

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
)

func exportConfig() error {
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./environments/")
	allEnvironments := os.Environ()
	fmt.Println(allEnvironments)
	switch os.Getenv("APP_ENV") {
	case "dev":
		viper.SetConfigName("dev")
	case "prod":
		viper.SetConfigName("prod")
	default:
		viper.SetConfigName("dev")
	}

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
