package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var classID string

var articleCmd = &cobra.Command{
	Use:   "article",
	Short: "`dedao-dl article` 获取文章详情",
	Long:  `使用 dedao-dl article 获取文章详情`,
	Run: func(cmd *cobra.Command, args []string) {
		app.ArticleList(classID)
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.PersistentFlags().StringVarP(&classID, "id", "i", "", "课程id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
