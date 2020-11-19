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

// GitOptionT ...
type GitOptionT struct {
	Name         string
	RequireValue bool
	IsPrefix     bool
}

// GitOption ...
type GitOption = *GitOptionT

// CommandLineT ...
type CommandLineT struct {
	SubCommand            string
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
	c.Printf("fgit %d.%d.%d - 快50倍的git clone github.com。\n", VersionMajor, VersionMinor, VersionFix)
}

func (me CommandLine) String() string {
	return JSONPretty(me)
}

func filterExtendedArguments(r CommandLine) {
	valueTrue := true

	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]

		if arg == "--private" {
			r.IsPrivate = &valueTrue
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

		r.Args = append(r.Args, arg)
	}
}

func resolveSubCommand(r CommandLine) {
	valueTrue := true

	arg0 := r.Args[0]
	if isOptionArg(arg0) {
		return
	}

	if arg0 == "clone" {
		r.SubCommand = "clone"

		r.IsGitClone = true
		r.PerhapsNeedInstrument = true
		return
	}

	if arg0 == "pull" || arg0 == "push" || arg0 == "fetch" {
		r.SubCommand = arg0
		r.PerhapsNeedInstrument = true
		if arg0 == "push" {
			r.IsPrivate = &valueTrue
		}
	}
}

// ParseCommandLine ...
func ParseCommandLine() CommandLine {
	valueFalse := false

	r := &CommandLineT{
		SubCommand:            "",
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
	resolveSubCommand(r)

	argSize := len(r.Args)
	if argSize == 0 {
		r.PerhapsNeedInstrument = false
	}

	if r.PerhapsNeedInstrument == false {
		return r
	}

	if r.SubCommand == "clone" {
		parseGitCloneCommandLine(r)
	} else if r.SubCommand == "fetch" || r.SubCommand == "pull" {
		parsePullOrFetchCommand(r)
	} else if r.SubCommand == "push" {
		parsePushCommand(r)
	} else {
		r.PerhapsNeedInstrument = false
	}

	if r.PerhapsNeedInstrument == false {
		return r
	}

	if !r.IsGitClone {
		if len(r.GitRemoteName) == 0 {
			r.GitRemoteName = ResolveGitRemoteName("")
		}
	}

	if len(r.GitURLText) == 0 {
		r.GitURLText = ResolveGitURLText(r.GitURLText, r.GitRemoteName, r.IsGitClone)
	}

	return r
}

func parsePullOrFetchCommand(r CommandLine) {
	argSize := len(r.Args)

	argsWithoutOptions := []string{}
	for i := 1; i < argSize; i++ {
		arg := r.Args[i]

		if !isOptionArg(arg) {
			argsWithoutOptions = append(argsWithoutOptions, arg)
		}
	}

	if len(argsWithoutOptions) > 0 {
		r.GitRemoteName = argsWithoutOptions[0]
	}

	r.PerhapsNeedInstrument = true
}

func parsePushCommand(r CommandLine) {
	valueTrue := true
	argSize := len(r.Args)

	argsWithoutOptions := []string{}
	for i := 1; i < argSize; i++ {
		arg := r.Args[i]

		if !isOptionArg(arg) {
			argsWithoutOptions = append(argsWithoutOptions, arg)
		} else {
			lenOfRepoOptionPrefix := len("--repo=")
			if len(arg) > lenOfRepoOptionPrefix && arg[0:lenOfRepoOptionPrefix] == "--repo=" {
				r.GitRemoteName = arg[lenOfRepoOptionPrefix:]
			}
		}
	}

	if len(r.GitRemoteName) == 0 {
		if len(argsWithoutOptions) > 0 {
			r.GitRemoteName = argsWithoutOptions[0]
		}
	}

	r.IsPrivate = &valueTrue
	r.PerhapsNeedInstrument = true
}

var cloneOptions = []GitOptionT{
	{"-v", false, false},
	{"--verbose", false, false},
	{"-q", false, false}, {"--quiet", false, false},
	{"--progress, false", false, false},
	{"-n", false, false}, {"--no-checkout", false, false},
	{"--bare", false, false},
	{"--mirror", false, false},
	{"-l", false, false}, {"--local", false, false},
	{"--no-hardlinks", false, false},
	{"-s", false, false}, {"--shared", false, false},
	{"--recursive", false, true},
	{"-j", true, false}, {"--jobs", true, false},
	{"--template", true, false},
	{"--reference", true, false},
	{"--reference-if-able", true, false},
	{"--dissociate", false, false},
	{"-o", true, false}, {"--origin", true, false},
	{"-b", true, false}, {"--branch", true, false},
	{"-u", true, false}, {"--upload-pack", true, false},
	{"--depth", true, false},
	{"--shallow-since", true, false},
	{"--shallow-exclude", true, false},
	{"--single-branch", false, false},
	{"--no-tags", false, false},
	{"--shallow-submodules", false, false},
	{"--separate-git-dir", true, false},
	{"-c", true, false}, {"--config", true, false},
	{"--server-option", true, false},
	{"-4", false, false}, {"--ipv4", false, false},
	{"-6", false, false}, {"--ipv6", false, false},
	{"--filter", true, false},
	{"--", true, false},
}

func parseGitCloneCommandLine(r CommandLine) {
	argSize := len(r.Args)

	for i := 1; i < argSize; i++ {
		arg := r.Args[i]
		argValue := ""

		var opt GitOption

		for _, t := range cloneOptions {
			if t.Name == arg || (t.IsPrefix && strings.Index(arg, t.Name) == 0) {
				opt = &t
				if t.RequireValue {
					if i < argSize-1 {
						argValue = r.Args[i+1]
					}
					i++
				}
				break
			}
		}

		if opt != nil {
			if arg == "-o" || arg == "--origin" {
				r.GitRemoteName = argValue
			} else if arg == "--" {
				if len(argValue) > 0 {
					r.ArgIndexOfGitURLText = i
					r.GitURLText = argValue
				}
			}
		} else if isOptionArg(arg) == false {
			if len(r.GitURLText) == 0 {
				r.ArgIndexOfGitURLText = i
				r.GitURLText = arg
			} else {
				r.GitCloneDir = arg
			}
		}

	}

	if len(r.GitURLText) == 0 {
		r.PerhapsNeedInstrument = false
		return
	}

	if len(r.GitCloneDir) == 0 {
		r.GitCloneDir = r.GitURLText[strings.LastIndex(r.GitURLText, "/")+1:]
		if strings.HasSuffix(strings.ToLower(r.GitCloneDir), ".git") {
			r.GitCloneDir = r.GitCloneDir[:len(r.GitCloneDir)-4]
		}
	}

}

func isOptionArg(arg string) bool {
	return len(arg) > 0 && arg[0:1] == "-"
}
