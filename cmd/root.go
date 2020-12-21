package cmd

import (
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/config"
)

var rootCmd = &cobra.Command{
	Use:   "dedao-dl",
	Short: "dedao-dl is a very fast dedao app course article download tools",
	Long: `A Fast dedao app course article download tools built with
		love by spf13 and friends in Go.`,
}

// AuthFunc check login
var AuthFunc = func(cmd *cobra.Command, args []string) error {
	if config.Instance.AcitveUID == "" {
		if len(config.Instance.Users) > 0 {
			return errors.New("存在登录的用户，可以进行切换登录用户")
		}
		return errors.New("请先前往 https://www.dedao.cn 登录得到账户")
	}
	return nil
}

// Execute exec cmd
func Execute() error {
	return rootCmd.Execute()
}
