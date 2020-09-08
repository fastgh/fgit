package main

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// ResolveGitURLText ...
func ResolveGitURLText(gitURLText string, remoteName string, isGitClone bool) string {
	if len(gitURLText) == 0 {
		if !isGitClone {
			gitURLText = GetGitRemoteURL(remoteName)
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

// ConfigGitHTTPProxy ...
func ConfigGitHTTPProxy(global bool, newHTTPProxy, newHTTPSProxy string) (oldHTTPProxy, oldHTTPSProxy string) {
	var scope string
	if global {
		scope = "--global"
	} else {
		scope = "--local"
	}

	ExecGit([]string{"config", scope, "--set", "http.proxy", newHTTPProxy})
	ExecGit([]string{"config", scope, "--set", "https.proxy", newHTTPSProxy})

	oldHTTPProxy = ExecGit([]string{"config", scope, "--get", "http.proxy"})
	oldHTTPSProxy = ExecGit([]string{"config", scope, "--get", "https.proxy"})

	return
}

// ResetGitHTTPProxy ...
func ResetGitHTTPProxy(global bool, oldHTTPProxy, oldHTTPSProxy string) {
	var scope string
	if global {
		scope = "--global"
	} else {
		scope = "--local"
	}

	if len(oldHTTPProxy) > 0 {
		ExecGit([]string{"config", scope, "--set", "http.proxy", oldHTTPProxy})
	} else {
		ExecGit([]string{"config", scope, "--unset-all", "http.proxy"})
	}

	if len(oldHTTPSProxy) > 0 {
		ExecGit([]string{"config", scope, "--set", "https.proxy", oldHTTPSProxy})
	} else {
		ExecGit([]string{"config", scope, "--unset-all", "https.proxy"})
	}
}

// GetGitRemoteURL ...
func GetGitRemoteURL(remoteName string) string {
	return ExecGit([]string{"remote", "get-url", "origin"})
}

// SetGitRemoteURL ...
func SetGitRemoteURL(remoteName string, remoteURL string) {
	ExecGit([]string{"remote", "set-url", "origin", remoteURL})
}

// ExecGit ...
func ExecGit(args []string) string {
	if Debug {
		fmt.Println("git " + strings.Join(args, " "))
	}

	if Mock {
		fmt.Println("mocking run: git " + strings.Join(args, " "))
		return ""
	}

	var command = exec.Command("git", args...)
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
