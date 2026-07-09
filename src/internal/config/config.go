package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func Read() *Config {
	homeDir, err := os.UserHomeDir()
	homeDir += "/"
	if err != nil {
		panic("could not obtain user home directory name")
	}
	configFile, err := os.Open(homeDir + configFileName)
	if err != nil {
		panic("could not open config file: " + homeDir + configFileName)
	}
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		panic("could not read config file: " + homeDir + configFileName)
	}
	defer configFile.Close()
	config := Config{}
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		panic("error unmarhsalling JSON")
	}
	fmt.Printf("{\ndb_url: \"%s\"\n}\n", config.DBUrl)
	return &config
}
