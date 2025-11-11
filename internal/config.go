package Config

import (
	"encoding/json"
	"fmt"
	"os"
)

const configFilePath = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (c *Config) SetUser(userName string) error {
	c.CurrentUserName = userName

	err := write(*c)

	if err != nil {
		return err
	}

	return nil
}

// Reads the config data from the config file
func Read() (Config, error) {

	configPath, err := getConfigPath()

	if err != nil {
		return Config{}, err
	}

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

func getConfigPath() (string, error) {

	userHomeDir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("Error in finding users home directory: %v", err)
	}

	configPath := fmt.Sprintf("%v/%v", userHomeDir, configFilePath)

	return configPath, nil
}

// Writes to the config file
func write(c Config) error {

	configPath, err := getConfigPath()

	if err != nil {
		return err
	}

	json, err := json.Marshal(c)

	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, json, 0644)

	if err != nil {
		return err
	}

	return nil
}
