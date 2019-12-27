package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"

	"github.com/cozy/cozy-stack/client"
	"github.com/cozy/cozy-stack/client/auth"
)

const (
	// UserAgent is used for the HTTP Header User-Agent to identify this
	// software when making requests to the stack in logs.
	UserAgent = "cozy-fuse-client"
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
	OAuthClient auth.Client `json:"client"`
	// Token is the bearer tokens for this OAuth client.
	Token auth.AccessToken `json:"token"`

	client *client.Client
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

// Client returns a client for making requests to cozy-stack.
// TODO manage properly the case where the config is not complete or is invalid
func (c *Config) Client() *client.Client {
	if c.client == nil {
		u, err := url.Parse(c.Instance)
		if err != nil {
			panic(err)
		}
		c.client = &client.Client{
			Addr:       u.Host,
			Domain:     u.Host,
			Scheme:     u.Scheme,
			AuthClient: &c.OAuthClient,
			Authorizer: &c.Token,
			UserAgent:  UserAgent,
		}
	}
	return c.client
}
