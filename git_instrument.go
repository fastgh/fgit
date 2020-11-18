package main

import (
	"log"
	"net/url"

	"github.com/pkg/errors"
)

// GithubInstrument ...GithubInstrument
func GithubInstrument(isPrivate bool, config Config) {
	var global bool
	var workDir string
	if cmdline.IsGitClone {
		global = true
		workDir = cmdline.GitCloneDir
	} else {
		global = false
		workDir = ""
	}

	if Debug {
		log.Printf("global: %v, workDir: %v", global, workDir)
	}

	oldHTTPProxy, oldHTTPSProxy := ConfigGitHTTPProxy(workDir, global, config.Proxy, config.Proxy)
	oldGit(true, false)
	ResetGitHTTPProxy(workDir, global, oldHTTPProxy, oldHTTPSProxy)
}

// InstrumentURLwithMirror ...
func InstrumentURLwithMirror(gitURLText string, mirrorURLText string) string {
	var err error

	var mirrorURL *url.URL
	if mirrorURL, err = url.Parse(mirrorURLText); err != nil {
		panic(errors.Wrapf(err, "无法解析URL: %s", mirrorURLText))
	}

	var gitURL *url.URL
	if gitURL, err = url.Parse(gitURLText); err != nil {
		panic(errors.Wrapf(err, "无法解析URL: %s", gitURLText))
	}

	gitURL.Scheme = mirrorURL.Scheme
	gitURL.Host = mirrorURL.Host
	gitURL.User = mirrorURL.User

	return gitURL.String()
}
