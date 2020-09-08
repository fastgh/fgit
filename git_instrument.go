package main

// SetGitRemoteURL ...
func SetGitRemoteURL(workDir string, remoteName string, remoteURL string) {
	ExecGit(workDir, []string{"remote", "set-url", "origin", remoteURL})
}

// GithubInstrument ...GithubInstrument
func GithubInstrument(isPrivate bool, config Config) {
	var global bool
	var workDir string
	if cmdline.IsGitClone {
		global = true
		workDir = cmdline.GitCloneDir
	} else {
		global = false
		workDir = ""
	}

	if isPrivate {
		//TODO: recover
		oldHTTPProxy, oldHTTPSProxy := ConfigGitHTTPProxy(workDir, global, config.Proxy, config.Proxy)
		oldGit(true, false)
		ResetGitHTTPProxy(workDir, global, oldHTTPProxy, oldHTTPSProxy)
	} else {
		//TODO:  drecover
		newCloneURL := InstrumentURLwithMirror(cmdline.GitURLText, config.Mirror)

		if !cmdline.IsGitClone {
			SetGitRemoteURL(workDir, cmdline.GitRemoteName, newCloneURL)
		} else {
			cmdline.Args[cmdline.ArgIndexOfGitURLText-1] = newCloneURL
		}

		oldGit(true, false)
		SetGitRemoteURL(workDir, cmdline.GitRemoteName, cmdline.GitURLText)
	}
}
