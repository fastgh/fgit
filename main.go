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

	rand.Seed(time.Now().UnixNano())

	//TODO: recover

	cmdline = ParseCommandLine()
	if Debug {
		log.Printf("Mock: %v\n", Mock)
		log.Printf("命令行: \n%v\n", cmdline)
	}

	if cmdline.PerhapsNeedInstrument == false {
		if Debug {
			log.Println("无需设置代理")
		}
		oldGit(false, false)
		return
	}

	defer func() {
		if p := recover(); p != nil {
			oldGit(false, true)
			color.Red.Printf("出错: %v\n", p)
			return
		}
	}()

	cmdline.GitURLText = ResolveGitURLText(cmdline.GitURLText, cmdline.GitRemoteName, cmdline.IsGitClone)
	gitURL := ResolveGitURL(cmdline.GitURLText)

	if Debug {
		log.Printf("GitURLText: %s, gitURL=%v\n", cmdline.GitURLText, gitURL)
	}

	if strings.ToLower(gitURL.Host) != "github.com" {
		if Debug {
			log.Println("not github.com, so skipped")
		}
		oldGit(false, false)
		return
	}

	if strings.ToLower(gitURL.Scheme) != "https" {
		fmt.Printf("不支持%s (仅支持https)\n", gitURL.Scheme)
		oldGit(false, false)
		return
	}

	var isPrivate bool
	if cmdline.IsPrivate != nil {
		isPrivate = *cmdline.IsPrivate
	} else if len(gitURL.User.Username()) > 0 {
		isPrivate = true

		if Debug {
			log.Println("发现URL中嵌入有用户名，因此设置为私有库模式")
		}
	}

	cfg := LoadConfig()
	if Debug {
		log.Printf("配置：%v\n", cfg)
	}

	HookInterruptSignal()

	GithubInstrument(isPrivate, cfg)
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
				fmt.Printf("程序崩溃, 错误原因: %s\n堆栈:\n%s", e, string(debug.Stack()))
			}
		}()
		for range signalChan {
			if Debug {
				log.Println("收到中断信号，退出前恢复原先的GITHUB设置...")
			}
			oldGit(false, true)
			if Debug {
				log.Println("完成恢复原先的GITHUB设置")
			}
			os.Exit(0)
		}
	}()
}
