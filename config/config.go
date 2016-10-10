// Package config deals with loading and parsing configuration from disk
package config

import (
	"fmt"
	"github.com/go-ini/ini"
	"github.com/mrtazz/github-unreleased/logger"
)

// Config is a struct to configure
type Config struct {
	iniFile *ini.File
}

// NewConfigFromFile returns a config struct that is filled with the values
// from the passed in ini file
func NewConfigFromFile(path string) (*Config, error) {
	iniFile, err := ini.Load(path)
	ret := &Config{
		iniFile: iniFile,
	}

	return ret, err
}

// GetConfigValue returns a configuration value from the config file from the
// default section.
func (c *Config) GetConfigValue(keyName string) (value string, err error) {
	if c.iniFile == nil {
		logger.Debug("iniFile not initialized")
		return "", fmt.Errorf("No config file open.")
	}

	section, err := c.iniFile.GetSection("default")

	if err != nil {
		return "", err
	}
	key, keyError := section.GetKey(keyName)
	if keyError != nil {
		return "", keyError
	}
	return key.String(), nil
}
