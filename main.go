package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/yann0917/dedao-dl/cmd"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/utils"
)

func init() {
	// clean 只删 output/.cache 目录，不依赖配置和缓存库。
	// 库损坏时 Init 会因打开失败直接退出（log.Fatalf / badger 后台协程 panic），
	// 若不跳过，用来救急的 clean cache 自己也起不来
	if isCleanCommand(os.Args[1:]) {
		return
	}

	err := config.Instance.Init()
	if err != nil {
		fmt.Println(err)
	}
}

func isCleanCommand(args []string) bool {
	return len(args) > 0 && args[0] == "clean"
}

func main() {
	if !isWebCommand(os.Args[1:]) {
		// 非 web 模式继续沿用信号清理，避免中断时遗留数据库句柄。
		setupCleanupOnExit()
	}

	defer closeBadgerDB()

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func setupCleanupOnExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("正在关闭程序...")

		closeBadgerDB()

		os.Exit(0)
	}()
}

func closeBadgerDB() {
	if err := utils.CloseBadgerDB(); err != nil {
		fmt.Printf("关闭数据库时出错: %v\n", err)
		return
	}
}

func isWebCommand(args []string) bool {
	return len(args) > 0 && args[0] == "web"
}
