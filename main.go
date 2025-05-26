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
	// 设置信号处理，确保程序退出前关闭数据库
	setupCleanupOnExit()

	cmd.Execute()
}

func setupCleanupOnExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		fmt.Println("正在关闭程序...")

		// 获取 BadgerDB 实例并关闭
		db, err := utils.GetBadgerDB(utils.GetDefaultBadgerDBPath())
		if err == nil && db != nil {
			if err := db.Close(); err != nil {
				fmt.Printf("关闭数据库时出错: %v\n", err)
			} else {
				fmt.Println("数据库已安全关闭")
			}
		}

		os.Exit(0)
	}()
}
