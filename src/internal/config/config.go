package config

import (
	"encoding/json"
	"io"
	"os"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func (c *Config) Read() {
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
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		panic("error unmarhsalling JSON")
	}
}

func (c *Config) SetUser(user string) {
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
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		panic("error unmarhsalling JSON")
	}

	configFile.Close()
	configFile, err = os.Create(homeDir + configFileName)
	if err != nil {
		panic("could not create config file: " + homeDir + configFileName)
	}
	defer configFile.Close()
	c.CurrentUserName = user
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		panic("error encoding json file")
	}
}
