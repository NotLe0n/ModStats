package config

import (
	"encoding/json"
	"log"

	"github.com/spf13/viper"
)

var C *viper.Viper = setup_config()

func setup_config() *viper.Viper {
	viper.SetDefault("port", "8080")
	viper.SetDefault("useHTTPS", false)
	viper.SetDefault("certPath", "")
	viper.SetDefault("keyPath", "")
	viper.SetDefault("API-URL", "https://tmlapis.le0n.dev")
	viper.SetDefault("timeout", "60")

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SafeWriteConfig()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s\n", err)
	}

	json_string, err := json.MarshalIndent(viper.AllSettings(), "", "\t")
	if err == nil {
		log.Printf("using config: \n%s\n", json_string)
	}

	return viper.GetViper()
}
