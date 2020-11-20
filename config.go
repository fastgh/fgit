package main

import (
	"fmt"
	"log"
	"math/rand"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

const (
	defaultControlServerURL   = "https://control.fastgithub.com:10443/api/v1"
	defaultReleaseDownloadURL = "https://github.com/fastgh/fgit/releases"
)

func init() {
	rand.Seed(time.Now().UnixNano()) //将时间戳设置成种子数
}

// AccountT ...
type AccountT struct {
	ID        string `json:"id"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"createdAt"`
}

// Account ...
type Account = *AccountT

// ConfigT ...
type ConfigT struct {
	Account            Account `json:"account"`
	ControlServerURL   string  `json:"controlServerURL"`
	ReleaseDownloadURL string  `json:"releaseDownloadURL"`
	Mirror             string  `json:"mirror"`
	Proxy              string  `json:"proxy"`
}

// Config ...
type Config = *ConfigT

// LoadConfig ...
func LoadConfig() Config {
	path := GetConfigJSONFilePath()

	if !FileExists(path) {
		password := NewUUID()
		accountID := CreateAccount(defaultControlServerURL, password)

		account := &AccountT{
			ID:       accountID,
			Password: password,
		}

		SaveConfigJSONFile(path, &ConfigT{
			Account:            account,
			ControlServerURL:   defaultControlServerURL,
			ReleaseDownloadURL: defaultReleaseDownloadURL,
		})
	}

	r := ConfigWithJSONFile(path)
	if len(r.ControlServerURL) == 0 {
		r.ControlServerURL = defaultControlServerURL
	}
	if len(r.ReleaseDownloadURL) == 0 {
		r.ReleaseDownloadURL = defaultReleaseDownloadURL
	}

	token := LoginByID(defaultControlServerURL, r.Account.ID, r.Account.Password)

	if len(r.Proxy) == 0 {
		proxy := SelectProxy(r.ControlServerURL, r.ReleaseDownloadURL)
		r.Proxy = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, r.Account.ID, token, proxy.Host, proxy.Port)
	}

	if len(r.Mirror) == 0 {
		r.Mirror = SelectMirror(r.ControlServerURL)
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
		log.Printf("[fgit] 配置文件：%s", r)
	}

	return r
}

// NewUUID ...
func NewUUID() string {
	newUUID, err := uuid.NewUUID()
	if err != nil {
		panic(errors.Wrap(err, "UUID生成失败"))
	}
	return newUUID.String()
}
