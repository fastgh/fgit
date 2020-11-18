package main

import (
	"log"
	"net/url"

	"github.com/pkg/errors"
)

// InstrumentContext ...
type InstrumentContext struct {
	oldHTTPProxy  string
	oldHTTPSProxy string
	useMirror     bool
	workDir       string
	global        bool
	shouldReset   bool
	mirroredURL   string
	originalURL   string
	remoteName    string
}

var instrumentContext = InstrumentContext{}

// GithubInstrument ...GithubInstrument
func GithubInstrument(isPrivate bool, config Config) {
	if isPrivate {
		instrumentContext.useMirror = false
	} else {
		instrumentContext.useMirror = true
	}

	if Cmdline.IsGitClone {
		instrumentContext.global = true
		instrumentContext.workDir = Cmdline.GitCloneDir
	} else {
		instrumentContext.global = false
		instrumentContext.workDir = ""
	}

	if Debug {
		log.Printf("[fgit] 修改前: %s\n", JSONPretty(instrumentContext))
	}

	defer ResetGithubRemote()

	if instrumentContext.useMirror {
		if Debug {
			log.Println("[fgit] 设置镜像")
		}

		instrumentContext.originalURL = Cmdline.GitURLText
		instrumentContext.mirroredURL = GeneratedMirroredURL(instrumentContext.originalURL, config.Mirror)

		instrumentContext.remoteName = Cmdline.GitRemoteName

		SetGitRemoteURL(instrumentContext.workDir, instrumentContext.remoteName, instrumentContext.mirroredURL)
	} else {
		if Debug {
			log.Println("[fgit] 设置代理")
		}
		instrumentContext.oldHTTPProxy, instrumentContext.oldHTTPSProxy = ConfigGitHTTPProxy(instrumentContext.workDir, instrumentContext.global, config.Proxy, config.Proxy)
	}
	instrumentContext.shouldReset = true

	if Debug {
		log.Printf("[fgit] 修改后: %s\n", JSONPretty(instrumentContext))
	}

	oldGit(true, false)
}

// ResetGithubRemote ...
func ResetGithubRemote() {
	if Debug {
		log.Printf("[fgit] 重置修改: %s\n", JSONPretty(instrumentContext))
	}

	if instrumentContext.shouldReset {
		if instrumentContext.useMirror {
			if Debug {
				log.Println("[fgit] 重置镜像")
			}
			SetGitRemoteURL(instrumentContext.workDir, instrumentContext.remoteName, instrumentContext.originalURL)
		} else {
			if Debug {
				log.Println("[fgit] 重置代理")
			}
			SetGitHTTPProxy(instrumentContext.workDir, instrumentContext.global, instrumentContext.oldHTTPProxy, instrumentContext.oldHTTPSProxy)
		}
	}
}

// GeneratedMirroredURL ...
func GeneratedMirroredURL(gitURLText string, mirrorURLText string) string {
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
