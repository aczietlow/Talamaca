package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Token    string
	Username string   `json:"user"`
	Owner    string   `json:"owner"`
	Repo     string   `json:"repo"`
	Users    []string `json:"users""`
}

var osReadFile = os.ReadFile

func LoadConfig(filepath string) (*Config, error) {
	configFile, err := osReadFile(filepath)

	if err != nil {
		fmt.Printf("%v\r\n", "Error reading config file.")
		return nil, err
	}
	var config Config
	if err := json.Unmarshal(configFile, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
