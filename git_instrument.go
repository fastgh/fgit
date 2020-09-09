package main

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

	//TODO: recover
	oldHTTPProxy, oldHTTPSProxy := ConfigGitHTTPProxy(workDir, global, config.Proxy, config.Proxy)
	oldGit(true, false)
	ResetGitHTTPProxy(workDir, global, oldHTTPProxy, oldHTTPSProxy)
}
