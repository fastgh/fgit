package main

import (
	"encoding/json"
	"os"
	"strings"

	"net/url"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

const (
	// AppVersion ...
	AppVersion = "1.0.0"
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
	c.Printf("fgit: version %s. Git client with built-in GITHUB proxy, 50x ~ 100x faster for Chinese developers.\n", AppVersion)
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
		GitCommand:            "",
		GitRemoteName:         "",
		IsGitClone:            false,
		GitCloneDir:           "",
		GitURLText:            "",
		ArgIndexOfGitURLText:  -1,
		PerhapsNeedInstrument: false,
		IsPrivate:             &valueFalse,
		Args:                  []string{},
	}

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

		if arg == "clone" {
			r.GitCommand = "clone"

			r.IsGitClone = true
			r.PerhapsNeedInstrument = true
		} else if r.IsGitClone {
			if r.ArgIndexOfGitURLText == -1 {
				r.ArgIndexOfGitURLText = i
				r.GitURLText = os.Args[r.ArgIndexOfGitURLText]
			} else {
				r.GitCloneDir = arg
			}
		} else if arg == "pull" || arg == "push" || arg == "fetch" {
			r.GitCommand = arg
			r.PerhapsNeedInstrument = true
			if i < len(os.Args)-1 {
				r.GitRemoteName = os.Args[i+1]
			}
			if r.GitCommand == "push" {
				r.IsPrivate = &valueTrue
			}
		}
	}

	if r.IsGitClone && len(r.GitCloneDir) == 0 {
		r.GitCloneDir = r.GitURLText[strings.LastIndex(r.GitURLText, "/")+1:]
		if strings.HasSuffix(strings.ToLower(r.GitCloneDir), ".git") {
			r.GitCloneDir = r.GitCloneDir[:len(r.GitCloneDir)-4]
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
