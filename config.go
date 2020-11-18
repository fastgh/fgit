package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/google/uuid"
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
		accountID, err := uuid.NewUUID()
		if err != nil {
			panic(errors.Wrap(err, "生成账号ID失败"))
		}
		SaveConfigJSONFile(path, &ConfigT{
			AccountID: accountID.String(),
		})
	}
	r := ConfigWithJSONFile(path)

	if len(r.Proxy) == 0 {
		if len(r.Mirror) == 0 {
			proxy := SelectProxy()
			r.Proxy = fmt.Sprintf("%s://%s:na@%s:%d", r.AccountID, proxy.Protocol, proxy.Host, proxy.Port)
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
	jsonText := JSONPretty(config)
	WriteFile(path, []byte(jsonText))
}

// ConfigWithJSON ...
func ConfigWithJSON(jsonText string) Config {
	r := &ConfigT{}
	JSONUnmarshal(jsonText, &r)
	return r
}

// GetConfigJSONFilePath return (file path)
func GetConfigJSONFilePath() string {
	dir, err := homedir.Dir()
	if err != nil {
		panic(errors.Wrapf(err, "无法获取用户主目录路径"))
	}

	r := filepath.Join(dir, ".fgit.json")
	if Debug {
		log.Printf("配置文件：%s", r)
	}

	return r
}
