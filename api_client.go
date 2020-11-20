package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// HTTPProxyServerInfoT ...
type HTTPProxyServerInfoT struct {
	Protocol     string `json:"protocol"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	Weight       int    `json:"weight"`
	VersionMajor int    `json:"versionMajor"`
	VersionMinor int    `json:"versionMinor"`
	VersionFix   int    `json:"versionFix"`
}

// HTTPProxyServerInfo ...
type HTTPProxyServerInfo = *HTTPProxyServerInfoT

// ListAllHTTPProxyServers ...
func ListAllHTTPProxyServers(controlServerURL string) []HTTPProxyServerInfo {
	apiURL := fmt.Sprintf("%s/servers/proxy?for=github.com", controlServerURL)

	if Debug {
		log.Printf("[fgit] 正在查询可用的代理服务器: %s\n", apiURL)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(errors.Wrap(err, "查询可用的代理服务器时失败"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "查询可用的代理服务器时失败"))
	}

	var r []HTTPProxyServerInfo
	JSONUnmarshal(string(body), &r)

	if Debug {
		log.Printf("[fgit] 查询到可用的代理服务器: %s\n", JSONPretty(r))
	}

	return r
}

// SelectProxy ...
func SelectProxy(controlServerURL string, releaseDownloadURL string) HTTPProxyServerInfo {
	proxies := ListAllHTTPProxyServers(controlServerURL)
	if proxies == nil || len(proxies) == 0 {
		if Debug {
			log.Println("[fgit] 没有可用的代理服务器")
		}
		return nil
	}

	r := proxies[rand.Intn(len(proxies))]
	if r.VersionMajor != VersionMajor {
		panic(fmt.Errorf("当前的fgit客户端版本是%d.%d.%d。该版本和服务器要求的版本%d.%d.%d不兼容。请在%s重新下载",
			VersionMajor, VersionMinor, VersionFix,
			r.VersionMajor, r.VersionMinor, r.VersionFix,
			releaseDownloadURL,
		))
	}

	if Debug {
		log.Printf("[fgit] 选择代理服务器: %s\n", JSONPretty(r))
	}

	return r
}

// ListAllMirrors ...
func ListAllMirrors(controlServerURL string) []string {
	apiURL := fmt.Sprintf("%s/servers/mirror?for=github.com", controlServerURL)

	if Debug {
		log.Printf("[fgit] 正在查询可用的镜像服务器: %s\n", apiURL)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(errors.Wrap(err, "查询可用的镜像服务器时失败"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "查询可用的镜像服务器时失败"))
	}

	var r []string
	JSONUnmarshal(string(body), &r)

	if Debug {
		log.Printf("[fgit] 查询到可用的镜像服务器: %s\n", JSONPretty(r))
	}

	return r
}

// SelectMirror ...
func SelectMirror(controlServerURL string) string {
	mirrors := ListAllMirrors(controlServerURL)
	if mirrors == nil || len(mirrors) == 0 {
		if Debug {
			log.Println("[fgit] 没有可用的镜像服务器")
		}
		return ""
	}

	r := mirrors[rand.Intn(len(mirrors))]
	if Debug {
		log.Printf("[fgit] 选择代理服务器: %s\n", JSONPretty(r))
	}

	return r
}

// CreateAccount ...
func CreateAccount(controlServerURL string, password string) string {
	apiURL := fmt.Sprintf("%s/account?password=%s", controlServerURL, password)

	if Debug {
		log.Printf("[fgit] 正在注册账号: %s\n", apiURL)
	}

	resp, err := http.Post(apiURL, "application/json", strings.NewReader(""))
	if err != nil {
		panic(errors.Wrap(err, "账号注册失败"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "账号注册失败"))
	}

	var accountID string
	JSONUnmarshal(string(body), &accountID)

	if Debug {
		log.Printf("[fgit] 账号注册成功: %s\n", accountID)
	}

	return accountID
}

// LoginByID ...
func LoginByID(controlServerURL string, accountID string, password string) string {
	apiURL := fmt.Sprintf("%s/account/_/%s/LoginByID?password=%s", controlServerURL, accountID, password)

	if Debug {
		log.Printf("[fgit] 正在登录: %s\n", apiURL)
	}

	resp, err := http.Get(apiURL)
	if err != nil {
		panic(errors.Wrap(err, "登录失败"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(errors.Wrap(err, "登录失败"))
	}

	var r string
	JSONUnmarshal(string(body), &r)

	if Debug {
		log.Printf("[fgit] 登录成功: %s\n", r)
	}

	return r
}
