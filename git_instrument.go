package main

import (
	"log"
	"net/url"

	"github.com/pkg/errors"
)

// InstrumentContext ...
type InstrumentContext struct {
	OldHTTPProxy  string
	OldHTTPSProxy string
	UseMirror     bool
	WorkDir       string
	Global        bool
	ShouldReset   bool
	MirroredURL   string
	OriginalURL   string
	RemoteName    string
}

var instrumentContext = InstrumentContext{}

// GithubInstrument ...GithubInstrument
func GithubInstrument(isPrivate bool, config Config) {
	if isPrivate {
		instrumentContext.UseMirror = false
	} else {
		instrumentContext.UseMirror = true
	}

	if Cmdline.IsGitClone {
		instrumentContext.Global = true
		instrumentContext.WorkDir = Cmdline.GitCloneDir
	} else {
		instrumentContext.Global = false
		instrumentContext.WorkDir = ""
	}

	if Debug {
		log.Printf("[fgit] 修改前: %s\n", JSONPretty(instrumentContext))
	}

	defer ResetGithubRemote()

	if instrumentContext.useMirror {
		if Debug {
			log.Println("[fgit] 设置镜像")
		}

		instrumentContext.OriginalURL = Cmdline.GitURLText
		instrumentContext.MirroredURL = GeneratedMirroredURL(instrumentContext.OriginalURL, config.Mirror)

		instrumentContext.RemoteName = Cmdline.GitRemoteName

		SetGitRemoteURL(instrumentContext.WorkDir, instrumentContext.RemoteName, instrumentContext.MirroredURL)
	} else {
		if Debug {
			log.Println("[fgit] 设置代理")
		}
		instrumentContext.OldHTTPProxy, instrumentContext.OldHTTPSProxy = ConfigGitHTTPProxy(instrumentContext.WorkDir, instrumentContext.Global, config.Proxy, config.Proxy)
	}
	instrumentContext.ShouldReset = true

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

	if instrumentContext.ShouldReset {
		if instrumentContext.UseMirror {
			if Debug {
				log.Println("[fgit] 重置镜像")
			}
			SetGitRemoteURL(instrumentContext.WorkDir, instrumentContext.RemoteName, instrumentContext.OriginalURL)
		} else {
			if Debug {
				log.Println("[fgit] 重置代理")
			}
			SetGitHTTPProxy(instrumentContext.WorkDir, instrumentContext.Global, instrumentContext.OldHTTPProxy, instrumentContext.OldHTTPSProxy)
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
