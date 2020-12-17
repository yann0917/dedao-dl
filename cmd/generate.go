package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var fileName = "dedao.json"

var userInfo = `{
	"user":{
		"phone":"",
		"password":"password",
		"cookie":""
	}
}
`
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "use `dedao-dl generate` to generate config file dedao.json",
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
		fmt.Println("## dedao.json generated ##")
		file, err := os.Create(fileName)
		defer file.Close()
		if err != nil {
			fmt.Println(err.Error())
		} else {
			file.Write([]byte(userInfo))
		}
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
