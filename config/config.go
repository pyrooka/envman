package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

// Store the path of the config file.
var configFilePath string

// Config defines the structure of the config file.
type Config struct {
	DefaultBackend string           `json:"defaultBackend"`
	Local          LocalConfig      `json:"local"`
	GitHubGist     GitHubGistConfig `json:"githubgist"`
}

// Helper functions.

// Determines the path of the config file on the current machine.
func getConfigPath() (err error) {
	// Get the home directory.
	currentUser, err := user.Current()
	if err != nil {
		return
	}
	homeDir := currentUser.HomeDir

	configFilePath = filepath.Join(homeDir, ".envman")

	return
}

// Load reads the config file from the disk.
func Load() (c *Config, err error) {
	// Get and set the config path.
	err = getConfigPath()
	if err != nil {
		return
	}

	// Read the config.
	data, err := ioutil.ReadFile(configFilePath)
	if os.IsNotExist(err) {
		err = nil
		// Create a new config.
		c = &Config{
			DefaultBackend: "local",
		}
		return
	} else if err != nil {
		return
	}

	// Parse the file.
	c = &Config{}
	err = json.Unmarshal(data, c)

	return
}

// Save writes the config to the file.
func (c *Config) Save() (err error) {
	// Create JSON from the struct.
	data, err := json.Marshal(c)
	if err != nil {
		return
	}

	// Write to file.
	file, err := os.Create(configFilePath)
	if err != nil {
		return
	}
	defer file.Close()

	_, err = file.Write(data)

	return
}
