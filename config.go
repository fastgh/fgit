package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
)

// AccountT ...
type AccountT struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	CreatedAt int64  `json:"createdAt"`
}

// Account ...
type Account = *AccountT

// ConfigT ...
type ConfigT struct {
	Account Account `json:"account"`
	Mirror  string  `json:"mirror"`
	Proxy   string  `json:"proxy"`
}

// Config ...
type Config = *ConfigT

func InputEmail() string {
	var email string
	fmt.Scanf("请输入注册邮件：%s", &email)
	return email
}

func InputPassword() string {
	var password string
	fmt.Scanf("请输入密码：%s", &password)
	return password
}

func InputPasswordAgain() string {
	var password string
	fmt.Scanf("请再次输入密码：%s", &password)
	return password
}

//
 LoadConfig ...
func LoadConfig() Config {
	path := GetConfigJSONFilePath()

	if !FileExists(path) {
		var email string
		fmt.Scanf("请输入注册邮件：%s", &email)

		var password string
		fmt.Scanf("请输入密码：%s", &password)

		account := CreateAccount(email, password)
		SaveConfigJSONFile(path, &ConfigT{
			Account: account,
		})
	}
	r := ConfigWithJSONFile(path)

	proxy := SelectProxy()
	r.Proxy = fmt.Sprintf("%s://%s:%s@%s:%d", proxy.Protocol, r.Account.ID, "NA", proxy.Host, proxy.Port)

	if len(r.Mirror) == 0 {
		r.Mirror = "https://github.com.cnpmjs.org"
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
