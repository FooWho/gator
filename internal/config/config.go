package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	DBUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("could not obtain user home directory: %w", err)
	}
	return filepath.Join(homeDir, configFileName), nil
}

func (c *Config) Read() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	configFile, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("could not open config file %s: %w", path, err)
	}
	defer configFile.Close()

	configBytes, err := io.ReadAll(configFile)
	if err != nil {
		return fmt.Errorf("could not read config file %s: %w", path, err)
	}

	if len(configBytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(configBytes, c); err != nil {
		return fmt.Errorf("error unmarshalling JSON from %s: %w", path, err)
	}
	return nil
}

func (c *Config) SetUser(user string) error {
	if err := c.Read(); err != nil {
		return fmt.Errorf("could not read existing config: %w", err)
	}

	c.CurrentUserName = user

	path, err := getConfigPath()
	if err != nil {
		return err
	}

	configFile, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("could not create config file %s: %w", path, err)
	}
	defer configFile.Close()

	encoder := json.NewEncoder(configFile)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(c); err != nil {
		return fmt.Errorf("error encoding json to file %s: %w", path, err)
	}
	return nil
}
