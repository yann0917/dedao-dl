package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "use `dedao-dl download` to login https://www.dedao.cn",
	Long:  `use dedao-dl download to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download cmd")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
