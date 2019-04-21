package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

var (
	// ErrNotExist is the error returned while trying to load a config file
	// that does not exist.
	ErrNotExist = errors.New("Config file does not exist")
)

// Config is the configuration for a mount point handled by cozy-fuse.
type Config struct {
	// Mount is the path to the directory where the user can see their files.
	Mount string `json:"mount"`
	// Data is a the path to a directory where cozy-fuse can persist the files
	// between two executions.
	Data string `json:"data"`
	// Instance is the URL of a cozy instance.
	Instance string `json:"instance"`
	// Client is the parameters of the OAuth client for this instance.
	Client struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
	} `json:"client"`
	// Token is the bearer tokens for this OAuth client.
	Token struct {
		Access  string `json:"access_token"`
		Refresh string `json:"refresh_token"`
	} `json:"token"`
}

// Load returns the config at the given path.
func Load(configPath string) (*Config, error) {
	f, err := ioutil.ReadFile(configPath)
	if err != nil {
		if err == os.ErrNotExist {
			err = ErrNotExist
		}
		return nil, err
	}
	cfg := &Config{}
	if err := json.Unmarshal(f, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
