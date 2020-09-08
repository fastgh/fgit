package main

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

// SetGitRemoteURL ...
func SetGitRemoteURL(remoteName string, remoteURL string) {
	ExecGit([]string{"remote", "set-url", "origin", remoteURL})
}

// GithubInstrument ...GithubInstrument
func GithubInstrument(isPrivate bool, config Config) {
	if isPrivate {
		var global bool
		if cmdline.IsGitClone {
			global = true
		} else {
			global = false
		}

		//TODO: recover
		oldHTTPProxy, oldHTTPSProxy := ConfigGitHTTPProxy(global, config.Proxy, config.Proxy)
		oldGit(true, true)
		ResetGitHTTPProxy(global, oldHTTPProxy, oldHTTPSProxy)
	} else {
		//TODO: recover
		newCloneURL := InstrumentURLwithMirror(cmdline.GitURLText, config.Mirror)
		if !cmdline.IsGitClone {
			SetGitRemoteURL(cmdline.GitRemoteName, newCloneURL)
		}

		oldGit(true, true)
		SetGitRemoteURL(cmdline.GitRemoteName, cmdline.GitURLText)
	}
}
