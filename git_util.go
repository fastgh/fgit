package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ConfigGitHTTPProxy ...
func ConfigGitHTTPProxy(workDir string, global bool, newHTTPProxy, newHTTPSProxy string) (oldHTTPProxy, oldHTTPSProxy string) {
	var scope string
	if global {
		scope = "--global"
		workDir = ""
	} else {
		scope = "--local"
	}

	oldHTTPProxy = ExecGit(workDir, []string{"config", scope, "--get", "http.https://github.com.proxy"})
	if strings.Index(oldHTTPProxy, "exit") >= 0 {
		oldHTTPProxy = ""
	}
	oldHTTPSProxy = ExecGit(workDir, []string{"config", scope, "--get", "https.https://github.com.proxy"})
	if strings.Index(oldHTTPSProxy, "exit") >= 0 {
		oldHTTPSProxy = ""
	}

	ExecGit(workDir, []string{"config", scope, "http.https://github.com.proxy", newHTTPProxy})
	ExecGit(workDir, []string{"config", scope, "https.https://github.com.proxy", newHTTPSProxy})

	return
}

// ResetGitHTTPProxy ...
func ResetGitHTTPProxy(workDir string, global bool, oldHTTPProxy, oldHTTPSProxy string) {
	var scope string
	if global {
		scope = "--global"
	} else {
		scope = "--local"
	}
	if len(oldHTTPProxy) > 0 {
		ExecGit(workDir, []string{"config", scope, "http.https://github.com.proxy", oldHTTPProxy})
	} else {
		ExecGit(workDir, []string{"config", scope, "--unset-all", "http.https://github.com.proxyy"})
	}

	if len(oldHTTPSProxy) > 0 {
		ExecGit(workDir, []string{"config", scope, "https.https://github.com.proxy", oldHTTPSProxy})
	} else {
		ExecGit(workDir, []string{"config", scope, "--unset-all", "https.https://github.com.proxy"})
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
		panic(errors.New("cannot resolve git url"))
	}

	return gitURLText
}

// ResolveGitURL ...
func ResolveGitURL(gitURLText string) *url.URL {
	var err error
	var r *url.URL

	if r, err = url.Parse(gitURLText); err != nil {
		panic(errors.Wrapf(err, "failed to parse url: %s", gitURLText))
	}

	if strings.ToLower(r.Scheme) != "https" {
		panic(fmt.Errorf("only https is supported but got: %s", r.Scheme))
	}

	return r

}

// GetGitRemoteURL ...
func GetGitRemoteURL(workDir string, remoteName string) string {
	return ExecGit(workDir, []string{"remote", "get-url", "origin"})
}

// ExecGit ...
func ExecGit(workDir string, args []string) string {
	if Debug {
		fmt.Println("git " + strings.Join(args, " "))
	}

	if Mock {
		fmt.Println("mocking run: git " + strings.Join(args, " "))
		return ""
	}

	var command = exec.Command("git", args...)
	if len(workDir) > 0 {
		command.Dir = workDir
	}
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
