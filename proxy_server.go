package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"

	"github.com/pkg/errors"
)

// ProxyServerInfoT ...
type ProxyServerInfoT struct {
	Protocol string `json:"protocol"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Weight   int    `json:"Weight"`

	VersionMajor int `json:"versionMajor"`
	VersionMinor int `json:"versionMinor"`
	VersionFix   int `json:"versionFix"`
}

// ProxyServerInfo ...
type ProxyServerInfo = *ProxyServerInfoT

// FetchAvailableProxies ...
func FetchAvailableProxies() []ProxyServerInfo {
	controlServerURL := "http://control.fastgithub.com:7000/api/v1/proxies?type=github.com"

	if Debug {
		log.Printf("[fgit] 正在查询可用的代理服务器: %s\n", controlServerURL)
	}

	resp, err := http.Get(controlServerURL)
	if err != nil {
		panic(errors.Wrap(err, "查询可用的代理服务器时失败"))
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var r []ProxyServerInfo
	JSONUnmarshal(string(body), &r)

	if Debug {
		log.Printf("[fgit] 查询到可用的代理服务器: \n%s\n", JSONPretty(r))
	}

	return r
}

// SelectProxy ...
func SelectProxy() ProxyServerInfo {
	proxies := FetchAvailableProxies()
	if proxies == nil || len(proxies) == 0 {
		if Debug {
			log.Println("[fgit] 没有可用的代理服务器")
		}
		return nil
	}

	r := proxies[rand.Intn(len(proxies))]
	if r.VersionMajor != VersionMajor {
		panic(fmt.Errorf("当前的fgit客户端版本是%d.%d.%d。该版本和服务器要求的版本%d.%d.%d不兼容。请在https://github.com/fastgh/fgit/releases重新下载",
			VersionMajor, VersionMinor, VersionFix,
			r.VersionMajor, r.VersionMinor, r.VersionFix,
		))
	}

	if Debug {
		log.Printf("[fgit] 使用代理服务器: \n%s\n", JSONPretty(r))
	}

	return r
}
