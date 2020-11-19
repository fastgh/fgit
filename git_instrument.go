package main

import (
	"log"
	"net/url"

	"github.com/pkg/errors"
)

// InstrumentContextT ...
type InstrumentContextT struct {
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

// InstrumentContext ...
type InstrumentContext = *InstrumentContextT

var instrumentContext = &InstrumentContextT{}

// GithubInstrument ...GithubInstrument
func GithubInstrument(cmdline CommandLine, isPrivate bool, config Config) {
	ictx := instrumentContext

	if isPrivate {
		ictx.UseMirror = false
	} else {
		ictx.UseMirror = true
	}

	if cmdline.IsGitClone {
		ictx.Global = true
		ictx.WorkDir = cmdline.GitCloneDir
	} else {
		ictx.Global = false
		ictx.WorkDir = cmdline.GitDir
	}

	if Debug {
		log.Printf("[fgit] 修改前: %s\n", JSONPretty(ictx))
	}

	defer ResetGithubRemote()

	if ictx.UseMirror {
		if Debug {
			log.Println("[fgit] 设置镜像")
		}

		ictx.OriginalURL = cmdline.GitURLText
		ictx.MirroredURL = GeneratedMirroredURL(ictx.OriginalURL, config.Mirror)

		ictx.RemoteName = cmdline.GitRemoteName

		if cmdline.IsGitClone {
			cmdline.Args[cmdline.ArgIndexOfGitURLText] = ictx.MirroredURL
		} else {
			SetGitRemoteURL(ictx.WorkDir, ictx.RemoteName, ictx.MirroredURL)
		}
	} else {
		if Debug {
			log.Println("[fgit] 设置代理")
		}
		ictx.OldHTTPProxy, ictx.OldHTTPSProxy = ConfigGitHTTPProxy(ictx.WorkDir, ictx.Global, config.Proxy, config.Proxy)
	}
	ictx.ShouldReset = true

	if Debug {
		log.Printf("[fgit] 修改后: %s\n", JSONPretty(ictx))
	}

	oldGit(cmdline, true, false)
}

// ResetGithubRemote ...
func ResetGithubRemote() {
	ictx := instrumentContext

	if Debug {
		log.Printf("[fgit] 重置修改: %s\n", JSONPretty(ictx))
	}

	if ictx.ShouldReset {
		if ictx.UseMirror {
			if Debug {
				log.Println("[fgit] 重置镜像")
			}
			SetGitRemoteURL(ictx.WorkDir, ictx.RemoteName, ictx.OriginalURL)
		} else {
			if Debug {
				log.Println("[fgit] 重置代理")
			}
			SetGitHTTPProxy(ictx.WorkDir, ictx.Global, ictx.OldHTTPProxy, ictx.OldHTTPSProxy)
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
