package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"
)

// ConfigGitHTTPProxy ...
func ConfigGitHTTPProxy(workDir string, global bool, newHTTPProxy string) (oldHTTPProxy, oldHTTPProxyAuthMethod string) {
	var scope string
	if global {
		scope = "--global"
		workDir = ""
	} else {
		scope = "--local"
	}

	oldHTTPProxy = ExecGit(false, workDir, []string{"config", scope, "--get", "http.https://github.com.proxy"})
	if strings.Index(oldHTTPProxy, "exit") >= 0 {
		oldHTTPProxy = ""
	}
	oldHTTPProxyAuthMethod = ExecGit(false, workDir, []string{"config", scope, "--get", "http.https://github.com.proxyAuthMethod"})
	if strings.Index(oldHTTPProxyAuthMethod, "exit") >= 0 {
		oldHTTPProxyAuthMethod = ""
	}

	ExecGit(false, workDir, []string{"config", scope, "http.https://github.com.proxy", newHTTPProxy})
	ExecGit(false, workDir, []string{"config", scope, "http.https://github.com.proxyAuthMethod", "basic"})

	return
}

// SetGitHTTPProxy ...
func SetGitHTTPProxy(workDir string, global bool, oldHTTPProxy, oldHTTPProxyAuthMethod string) {
	var scope string
	if global {
		scope = "--global"
	} else {
		scope = "--local"
	}
	if len(oldHTTPProxy) > 0 {
		ExecGit(false, workDir, []string{"config", scope, "http.https://github.com.proxy", oldHTTPProxy})
	} else {
		ExecGit(false, workDir, []string{"config", scope, "--unset-all", "http.https://github.com.proxy"})
	}

	if len(oldHTTPProxyAuthMethod) > 0 {
		ExecGit(false, workDir, []string{"config", scope, "https.https://github.com.proxyAuthMethod", oldHTTPProxyAuthMethod})
	} else {
		ExecGit(false, workDir, []string{"config", scope, "--unset-all", "https.https://github.com.proxyAuthMethod"})
	}
}

// ResolveGitURLText ...
func ResolveGitURLText(gitURLText string, remoteName string, isGitClone bool) string {
	if len(gitURLText) == 0 {
		if !isGitClone {
			gitURLText = GetGitRemoteURL("", remoteName)
		}
	}

	if len(gitURLText) == 0 {
		panic(fmt.Sprintf("获取GIT URL失败: %s", gitURLText))
	}

	return gitURLText
}

// ResolveGitRemoteName ...
func ResolveGitRemoteName(workDir string) string {
	r := ExecGit(false, workDir, []string{"remote"})
	r = strings.Trim(r, "\n\r\t ")

	posOfReturn := strings.Index(r, "\n")
	if posOfReturn > 0 {
		r = r[0:posOfReturn]
	}

	posOfReturn = strings.Index(r, "\r")
	if posOfReturn > 0 {
		r = r[0:posOfReturn]
	}

	return r
}

// ResolveGitURL ...
func ResolveGitURL(gitURLText string) *url.URL {
	var err error
	var r *url.URL

	if r, err = url.Parse(gitURLText); err != nil {
		panic(errors.Wrapf(err, "解析URL失败: %s", gitURLText))
	}

	return r

}

// GetGitRemoteURL ...
func GetGitRemoteURL(workDir string, remoteName string) string {
	return ExecGit(false, workDir, []string{"remote", "get-url", remoteName})
}

// SetGitRemoteURL ...
func SetGitRemoteURL(workDir string, remoteName string, remoteURL string) {
	ExecGit(false, workDir, []string{"remote", "set-url", remoteName, remoteURL})
}

// ExecGit ...
func ExecGit(redirect bool, workDir string, args []string) string {
	if len(workDir) > 0 {
		if DirExists(workDir) == false {
			workDir = path.Join(GetWorkDir(), workDir)
			if DirExists(workDir) == false {
				workDir = GetWorkDir()
			}
		}
	} else {
		workDir = GetWorkDir()
	}

	if Debug {
		log.Printf("[fgit] %s: git %s\n", workDir, strings.Join(args, " "))
	}

	var command = exec.Command("git", args...)

	if len(workDir) > 0 {
		command.Dir = workDir
	}

	if redirect {
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		var err = command.Start()
		if err != nil {
			return err.Error()
		}
		err = command.Wait()
		if err != nil {
			return err.Error()
		}
		return ""
	}

	t, _ := command.Output()
	r := string(t)
	r = strings.Trim(r, "\n\r\t ")

	if Debug {
		log.Printf("[fgit] return: %s\n", r)
	}
	return r
}
