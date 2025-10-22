package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// TODO: Add method on Config struct to set the user on the config

// Reads the config data from the config file
func Read() (Config, error) {

	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return Config{}, fmt.Errorf("Error in finding users home directory: %v", err)
	}

	configPath := fmt.Sprintf("%v/.gatorconfig.json", userHomeDir)

	rawConfigContent, err := os.ReadFile(configPath)

	if err != nil {
		return Config{}, fmt.Errorf("Error in reading config file: %v", err)
	}

	var config Config

	err = json.Unmarshal(rawConfigContent, &config)

	if err != nil {
		return Config{}, fmt.Errorf("Error unmarshalling data from config file: %v", err)
	}

	return config, nil
}
