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
	err := config.Instance.Init()
	if err != nil {
		fmt.Println(err)
	}
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
