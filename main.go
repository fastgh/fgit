package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/gookit/color"
)

func oldGit(cmdline CommandLine, fgitHelpFirst bool, errorMode bool) {
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

	rand.Seed(time.Now().UnixNano())

	//TODO: recover

	cmdline := ParseCommandLine()
	if Debug {
		log.Printf("[fgit] 命令行: \n%s\n", JSONPretty(cmdline))
	}

	if cmdline.PerhapsNeedInstrument == false {
		if Debug {
			log.Println("[fgit] 无需设置代理")
		}
		oldGit(cmdline, false, false)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			color.Red.Printf("[fgit] 出错: %v\n", p)
			ResetGithubRemote()
			return
		}
	}()

	gitURL := ResolveGitURL(cmdline.GitURLText)

	if Debug {
		log.Printf("[fgit] GitURLText: %s, gitURL=%s\n", cmdline.GitURLText, JSONMarshal(gitURL))
	}

	if strings.ToLower(gitURL.Host) != "github.com" {
		if Debug {
			log.Println("[fgit] 忽略非github.com库")
		}
		oldGit(cmdline, false, false)
		return
	}

	if strings.ToLower(gitURL.Scheme) != "https" {
		color.Yellow.Printf("[fgit] 不支持%s (仅支持https)\n", gitURL.Scheme)
		oldGit(cmdline, false, false)
		return
	}

	var isPrivate bool
	if cmdline.IsPrivate != nil {
		isPrivate = *cmdline.IsPrivate
	} else if len(gitURL.User.Username()) > 0 {
		isPrivate = true

		if Debug {
			log.Println("[fgit] 发现URL中嵌入有用户名，因此设置为私有库模式")
		}
	}

	cfg := LoadConfig()
	if Debug {
		log.Printf("[fgit] 配置：\n%s\n", JSONPretty(cfg))
	}

	HookInterruptSignal()

	GithubInstrument(cmdline, isPrivate, cfg)
}

// HookInterruptSignal ...
func HookInterruptSignal() {
	signalChan := make(chan os.Signal, 1)

	signal.Notify(signalChan,
		os.Interrupt,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	go func() {
		defer func() {
			if e := recover(); e != nil {
				color.Red.Printf("[fgit] 程序崩溃, 错误原因: %s\n堆栈:\n%s", e, string(debug.Stack()))
			}
		}()
		for range signalChan {
			if Debug {
				log.Println("[fgit] 收到中断信号，退出前恢复原先的GITHUB设置...")
			}
			ResetGithubRemote()
			if Debug {
				log.Println("[fgit] 完成恢复原先的GITHUB设置")
			}
			os.Exit(0)
		}
	}()
}
