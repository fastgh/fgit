package main

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// ConfigT ...
type ConfigT struct {
	AccountID string `json:"account-id"`
	Mirror    string `json:"mirror"`
	Proxy     string `json:"proxy"`
}

// Config ...
type Config = *ConfigT

// LoadConfig ...
func LoadConfig() Config {
	path := GetConfigJSONFilePath()
	if !FileExists(path) {
		SaveConfigJSONFile(path, &ConfigT{})
	}
	r := ConfigWithJSONFile(path)

	if len(r.Proxy) == 0 {
		if len(r.Mirror) == 0 {
			proxy := SelectProxy()
			r.Proxy = fmt.Sprintf("%s://%s:%d", proxy.Protocol, proxy.Host, proxy.Port)
			//panic("proxy not configured")
		}
	}

	return r
}

// ConfigWithJSONFile ...
func ConfigWithJSONFile(path string) Config {
	jsonText := string(ReadFile(path))
	return ConfigWithJSON(jsonText)
}

// SaveConfigJSONFile ...
func SaveConfigJSONFile(path string, config Config) {
	jsonText, err := json.Marshal(config)
	if err != nil {
		panic(errors.Wrapf(err, "failed to marshal json: %v\n", jsonText))
	}
	WriteFile(path, []byte(jsonText))
}

// ConfigWithJSON ...
func ConfigWithJSON(jsonText string) Config {
	r := &ConfigT{}
	if err := json.Unmarshal([]byte(jsonText), &r); err != nil {
		panic(errors.Wrapf(err, "failed to unmarshal json: %s\n"+jsonText))
	}
	return r
}

// GetConfigJSONFilePath return (file path)
func GetConfigJSONFilePath() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(errors.Wrapf(err, "failed to get home dir"))
	}

	return filepath.Join(dir, ".fgit.json")
}
