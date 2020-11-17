package main

import (
	"encoding/json"
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
}

// ProxyServerInfo ...
type ProxyServerInfo = *ProxyServerInfoT

// FetchAvailableProxies ...
func FetchAvailableProxies() []ProxyServerInfo {
	resp, err := http.Get("http://fastgithub.com:7000")
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

	return proxies[rand.Intn(len(proxies))]
}
