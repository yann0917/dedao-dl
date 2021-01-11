package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var ebookCmd = &cobra.Command{
	Use:     "ebook",
	Short:   "获取我的电子书架",
	Long:    `使用 dedao-dl ebook 获取我的电子书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ebook",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("bookID", bookID)
		if bookID > 0 {
			return app.EbookDetail(bookID)

		}
		return courseList(app.CateEbook)
	},
}

func init() {
	rootCmd.AddCommand(ebookCmd)

	ebookCmd.PersistentFlags().IntVarP(&bookID, "id", "i", 0, "电子书ID")
	// rootCmd.PersistentFlags().StringVarP(&cType, "type", "t", "bauhinia", "课程类型(all, bauhinia, ebook, compass")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
