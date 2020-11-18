package main

import (
	"os"
	"strings"

	"github.com/gookit/color"
)

const (
	// VersionMajor ...
	VersionMajor = 1

	// VersionMinor ...
	VersionMinor = 0

	// VersionFix ...
	VersionFix = 0
)

// CommandLineT ...
type CommandLineT struct {
	GitCommand            string
	GitRemoteName         string
	IsGitClone            bool
	GitCloneDir           string
	GitURLText            string
	ArgIndexOfGitURLText  int
	PerhapsNeedInstrument bool
	IsPrivate             *bool
	Args                  []string
}

// CommandLine ...
type CommandLine = *CommandLineT

// PrintHelp ...
func PrintHelp(errorMode bool) {
	var c color.Color
	if errorMode {
		c = color.Red
	} else {
		c = color.Blue
	}
	c.Printf("fgit %d.%d.%d - 让中国开发者git clone https://github.com时提速100倍。\n", VersionMajor, VersionMinor, VersionFix)
}

func (me CommandLine) String() string {
	return JSONPretty(me)
}

func filterExtendedArguments(cmdline CommandLine) {
	valueTrue := true

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == "--private" {
			cmdline.IsPrivate = &valueTrue
			continue
		}
		if arg == "--debug" {
			Debug = true
			continue
		}
		if arg == "--mock" {
			Mock = true
			continue
		}

		cmdline.Args = append(cmdline.Args, arg)
	}
}

func resolveGitCommand(cmdline CommandLine) {
	valueTrue := true

	arg0 := cmdline.Args[0]
	if isOptionArg(arg0) {
		return
	}

	if arg0 == "clone" {
		cmdline.GitCommand = "clone"

		cmdline.IsGitClone = true
		cmdline.PerhapsNeedInstrument = true
		return
	}

	if arg0 == "pull" || arg0 == "push" || arg0 == "fetch" {
		cmdline.GitCommand = arg0
		cmdline.PerhapsNeedInstrument = true
		if arg0 == "push" {
			cmdline.IsPrivate = &valueTrue
		}
	}
}

// ParseCommandLine ...
func ParseCommandLine() CommandLine {
	valueFalse := false

	r := &CommandLineT{
		GitCommand:            "",
		GitRemoteName:         "origin",
		IsGitClone:            false,
		GitCloneDir:           "",
		GitURLText:            "",
		ArgIndexOfGitURLText:  -1,
		PerhapsNeedInstrument: false,
		IsPrivate:             &valueFalse,
		Args:                  []string{},
	}

	filterExtendedArguments(r)
	resolveGitCommand(r)

	argSize := len(r.Args)
	if argSize == 0 {
		r.PerhapsNeedInstrument = false
	}

	if r.PerhapsNeedInstrument == false {
		return r
	}

	if r.GitCommand == "clone" {
		parseGitCloneCommandLine(r)
	} else if r.GitCommand == "fetch" || r.GitCommand == "pull" {
		parsePullOrFetchCommand(r)
	} else if r.GitCommand == "push" {
		parsePushCommand(r)
	} else {
		r.PerhapsNeedInstrument = false
	}

	if r.PerhapsNeedInstrument == false {
		return r
	}

	if !r.IsGitClone {
		if len(Cmdline.GitRemoteName) == 0 {
			Cmdline.GitRemoteName = ResolveGitRemoteName("")
		}
	}

	if len(Cmdline.GitURLText) == 0 {
		Cmdline.GitURLText = ResolveGitURLText(Cmdline.GitURLText, Cmdline.GitRemoteName, Cmdline.IsGitClone)
	}

	return r
}

func parsePullOrFetchCommand(cmdline CommandLine) {
	argSize := len(cmdline.Args)

	argsWithoutOptions := []string{}
	for i := 1; i < argSize; i++ {
		arg := cmdline.Args[i]

		if !isOptionArg(arg) {
			argsWithoutOptions = append(argsWithoutOptions, arg)
		}
	}

	if len(argsWithoutOptions) > 0 {
		cmdline.GitRemoteName = argsWithoutOptions[0]
	}

	cmdline.PerhapsNeedInstrument = true
}

func parsePushCommand(cmdline CommandLine) {
	valueTrue := true
	argSize := len(cmdline.Args)

	argsWithoutOptions := []string{}
	for i := 1; i < argSize; i++ {
		arg := cmdline.Args[i]

		if !isOptionArg(arg) {
			argsWithoutOptions = append(argsWithoutOptions, arg)
		} else {
			lenOfRepoOptionPrefix := len("--repo=")
			if len(arg) > lenOfRepoOptionPrefix && arg[0:lenOfRepoOptionPrefix] == "--repo=" {
				cmdline.GitRemoteName = arg[lenOfRepoOptionPrefix:]
			}
		}
	}

	if len(cmdline.GitRemoteName) == 0 {
		if len(argsWithoutOptions) > 0 {
			cmdline.GitRemoteName = argsWithoutOptions[0]
		}
	}

	cmdline.IsPrivate = &valueTrue
	cmdline.PerhapsNeedInstrument = true
}

func parseGitCloneCommandLine(cmdline CommandLine) {
	argSize := len(cmdline.Args)

	for i := 1; i < argSize; i++ {
		arg := cmdline.Args[i]

		argNext := ""
		if i < argSize-1 {
			argNext = cmdline.Args[i+1]
		}

		if isOptionArg(arg) == false {
			continue
		}

		if arg == "-o" || arg == "--origin" {
			cmdline.GitRemoteName = argNext
		} else if arg == "--" {
			if len(argNext) > 0 {
				cmdline.ArgIndexOfGitURLText = i
				cmdline.GitURLText = argNext
			}
		}
	}

	argLast := cmdline.Args[argSize-1]

	argLastPrev := ""
	if argSize > 1 {
		argLastPrev = cmdline.Args[argSize-2]
	}

	if len(cmdline.GitURLText) == 0 {
		if isOptionArg(argLast) == false {
			if isOptionArg(argLastPrev) {
				cmdline.ArgIndexOfGitURLText = argSize - 1
				cmdline.GitURLText = argLast
			} else {
				if argSize >= 2 {
					cmdline.ArgIndexOfGitURLText = argSize - 2
				}
				cmdline.GitURLText = argLastPrev
				cmdline.GitCloneDir = argLast
			}
		}
	}

	if len(cmdline.GitURLText) == 0 {
		cmdline.PerhapsNeedInstrument = false
		return
	}

	if len(cmdline.GitCloneDir) == 0 {
		cmdline.GitCloneDir = cmdline.GitURLText[strings.LastIndex(cmdline.GitURLText, "/")+1:]
		if strings.HasSuffix(strings.ToLower(cmdline.GitCloneDir), ".git") {
			cmdline.GitCloneDir = cmdline.GitCloneDir[:len(cmdline.GitCloneDir)-4]
		}
	}

}

func isOptionArg(arg string) bool {
	return len(arg) > 0 && arg[0:1] != "-"
}
