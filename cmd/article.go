package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var articleCmd = &cobra.Command{
	Use:     "article",
	Short:   "获取文章详情",
	Long:    `使用 dedao-dl article 获取文章详情`,
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		if classID > 0 && articleID == 0 {
			app.ArticleList(classID)
		}

		if articleID > 0 {
			err := app.ArticleDetail(classID, articleID)
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(articleCmd)
	articleCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程id")
	articleCmd.PersistentFlags().IntVarP(&articleID, "aid", "a", 0, "文章id")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
