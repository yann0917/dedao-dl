package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buyCmd = &cobra.Command{
	Use:   "buy",
	Short: "use `dedao-dl buy` to login https://www.dedao.cn",
	Long:  `use dedao-dl buy to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("buy cmd")
	},
}

func init() {
	rootCmd.AddCommand(buyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
