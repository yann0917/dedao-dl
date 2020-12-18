package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/config"
)

var rootCmd = &cobra.Command{
	Use:   "dedao-dl",
	Short: "dedao-dl is a very fast dedao app course article download tools",
	Long: `A Fast dedao app course article download tools built with
		love by spf13 and friends in Go.
		Complete documentation is available at http://hugo.spf13.com`,
}

var AuthFunc = func(cmd *cobra.Command, args []string) {
	if config.Instance.AcitveUID == "" {
		fmt.Println("authFunc")
		fmt.Println(len(config.Instance.Users))
		if len(config.Instance.Users) > 0 {
			fmt.Println("存在登录的用户，可以进行切换登录用户")

		}
		fmt.Println("请先登录极客时间账户")
		os.Exit(1)
	}
}

// Execute exec cmd
func Execute() error {
	return rootCmd.Execute()
}
