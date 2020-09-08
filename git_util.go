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

// GetGitRemoteURL ...
func GetGitRemoteURL(remoteName string) string {
	return ExecGit([]string{"remote", "get-url", "origin"})
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
