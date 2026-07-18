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

func (c *Config) Read() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not obtain user home directory name - %v", err)
	}
	homeDir += "/"

	configFile, err := os.Open(homeDir + configFileName)
	if err != nil {
		return fmt.Errorf("could not open config file: %s - %v", homeDir+configFileName, err)
	}
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file: %s - %v", homeDir+configFileName, err)
	}
	defer configFile.Close()
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return fmt.Errorf("error unmarhsalling JSON - %v", err)
	}
	return nil
}

func (c *Config) SetUser(user string) error {
	homeDir, err := os.UserHomeDir()
	homeDir += "/"
	if err != nil {
		return fmt.Errorf("could not obtain user home directory name - %v", err)
	}
	configFile, err := os.Open(homeDir + configFileName)
	if err != nil {
		return fmt.Errorf("could not open config file: %s - %v", homeDir+configFileName, err)
	}
	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file: %s - %v", homeDir+configFileName, err)
	}
	err = json.Unmarshal(configBytes, c)
	if err != nil {
		return fmt.Errorf("error unmarhsalling JSON - %v", err)
	}

	configFile.Close()
	configFile, err = os.Create(homeDir + configFileName)
	if err != nil {
		return fmt.Errorf("could not create config file: %s - %v", homeDir+configFileName, err)
	}
	defer configFile.Close()
	c.CurrentUserName = user
	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(c)
	if err != nil {
		return fmt.Errorf("error encoding json file - %v", err)
	}
	return nil
}
