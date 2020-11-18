package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	resp, err := http.Get("http://control.fastgithub.com:7000/api/v1/proxies")
	if err != nil {
		panic(errors.Wrap(err, "failed to fetch available proxies"))
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var r []ProxyServerInfo
	err = json.Unmarshal(body, &r)
	if err != nil {
		panic(errors.Wrapf(err, "failed to parse proxy list. got: %s", string(body)))
	}
	return r
}

// SelectProxy ...
func SelectProxy() ProxyServerInfo {
	proxies := FetchAvailableProxies()
	if proxies == nil || len(proxies) == 0 {
		return nil
	}

	r := proxies[rand.Intn(len(proxies))]
	if r.VersionMajor != VersionMajor {
		panic(fmt.Errorf("当前的fgit客户端版本是%d.%d.%d。该版本和服务器要求的版本%d.%d.%d不兼容。请在https://github.com/fastgh/fgit/releases重新下载",
			VersionMajor, VersionMinor, VersionFix,
			r.VersionMajor, r.VersionMinor, r.VersionFix,
		))
	}

	return r
}
