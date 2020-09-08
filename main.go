package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/gookit/color"
)

var cmdline CommandLine

func oldGit(fgitHelpFirst bool, errorMode bool) {
	if cmdline == nil {
		fgitHelpFirst = false
	}

	if fgitHelpFirst {
		PrintHelp(errorMode)
		fmt.Println()

		if cmdline == nil {
			ExecGit("", os.Args[1:])
		} else {
			ExecGit("", cmdline.Args)
		}
	} else {
		if cmdline == nil {
			ExecGit("", os.Args[1:])
		} else {
			ExecGit("", cmdline.Args)
		}

		fmt.Println()
		PrintHelp(errorMode)
	}
}

func main() {
	//TODO: recover

	cmdline = ParseCommandLine()
	if Debug {
		fmt.Printf("Mock: %v\n", Mock)
		fmt.Printf("Command line: \n%v\n", cmdline)
	}

	if cmdline.PerhapsNeedInstrument == false {
		oldGit(false, false)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			oldGit(false, true)
			color.Red.Printf("error: %v\n", p)
			return
		}
	}()

	cmdline.GitURLText = ResolveGitURLText(cmdline.GitURLText, cmdline.GitRemoteName, cmdline.IsGitClone)
	gitURL := ResolveGitURL(cmdline.GitURLText)

	if Debug {
		fmt.Printf("GitURLText: %s, gitURL=%v\n", cmdline.GitURLText, gitURL)
	}

	if strings.ToLower(gitURL.Host) != "github.com" {
		if Debug {
			fmt.Println("not github.com")
		}
		oldGit(false, false)
		return
	}

	var isPrivate bool
	if cmdline.IsPrivate != nil {
		isPrivate = *cmdline.IsPrivate
	} else if len(gitURL.User.Username()) > 0 {
		isPrivate = true

		if Debug {
			fmt.Println("set private=true because detected user name in url")
		}
	}

	cfg := LoadConfig()
	GithubInstrument(isPrivate, cfg)
}
