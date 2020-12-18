package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var fileName = "config.json"

var userInfo = `config.Dedao{
	User: "",
}`
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "use `dedao-dl generate` to generate config file config.json",
	Long: `use dedao-dl generate to generate dedao.json:
	{
		"user":{
			"phone":"phone",
			"password":"password",
			"cookie":""
		}
	}
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("## config.json generated ##")
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
