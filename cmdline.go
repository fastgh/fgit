package main

import (
	"encoding/json"
	"os"

	"net/url"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

const (
	// AppVersion ...
	AppVersion = "0.9.0"
)

// CommandLineT ...
type CommandLineT struct {
	GitCommand               string
	GitRemoteName            string
	IsGitClone               bool
	GitURLText               string
	ArgIndexOfGitURLText     int
	PerhapsNeedProxyOrMirror bool
	IsPrivate                *bool
	Args                     []string
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
	c.Printf("fgit: a fastgithub git client, version %s\n", AppVersion)
}

func (me CommandLine) String() string {
	json, err := json.MarshalIndent(me, "", "  ")
	if err != nil {
		panic(errors.Wrapf(err, "failed to json mashal"))
	}
	return string(json)
}

// ParseCommandLine ...
func ParseCommandLine() CommandLine {
	valueTrue := true
	valueFalse := false

	r := &CommandLineT{
		GitCommand:               "",
		GitRemoteName:            "",
		IsGitClone:               false,
		GitURLText:               "",
		ArgIndexOfGitURLText:     -1,
		PerhapsNeedProxyOrMirror: false,
		IsPrivate:                &valueFalse,
		Args:                     []string{},
	}

	hasCmd := false

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

		if arg[0:1] == "-" {
			continue
		}

		if hasCmd {
			continue
		}
		hasCmd = true
		r.GitCommand = arg

		if r.GitCommand == "clone" {
			r.IsGitClone = true
			r.PerhapsNeedProxyOrMirror = true

			if i < len(os.Args)-1 {
				r.ArgIndexOfGitURLText = i + 1
				r.GitURLText = os.Args[r.ArgIndexOfGitURLText]
			}
		} else if r.GitCommand == "pull" || r.GitCommand == "push" || r.GitCommand == "fetch" {
			r.PerhapsNeedProxyOrMirror = true
			if i < len(os.Args)-1 {
				r.GitRemoteName = os.Args[i+1]
			}
			if r.GitCommand == "push" {
				r.IsPrivate = &valueTrue
			}
		}
	}

	return r
}

// InstrumentURLwithMirror ...
func InstrumentURLwithMirror(gitURLText string, mirrorURLText string) string {
	var err error

	var mirrorURL *url.URL
	if mirrorURL, err = url.Parse(mirrorURLText); err != nil {
		panic(errors.Wrapf(err, "failed to parse url: %s", mirrorURLText))
	}

	var gitURL *url.URL
	if gitURL, err = url.Parse(gitURLText); err != nil {
		panic(errors.Wrapf(err, "failed to parse url: %s", gitURLText))
	}

	gitURL.Scheme = mirrorURL.Scheme
	gitURL.Host = mirrorURL.Host
	gitURL.User = mirrorURL.User

	return gitURL.String()
}
