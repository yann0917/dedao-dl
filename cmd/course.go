package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var courseCmd = &cobra.Command{
	Use:   "course",
	Short: "use `dedao-dl course` to login https://www.dedao.cn",
	Long:  `use dedao-dl course to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("course cmd")
	},
}

func init() {
	rootCmd.AddCommand(courseCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
