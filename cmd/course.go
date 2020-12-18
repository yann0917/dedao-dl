package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var cList, cType string
var courseCmd = &cobra.Command{
	Use:   "course",
	Short: "`dedao-dl course`获取购买过的课程",
	Long:  `使用 dedao-dl course 获取购买过的课程，课程分类，章节等`,
	Run: func(cmd *cobra.Command, args []string) {
		app.CourseList(cList, 1, 300)
	},
}

var courseTypeCmd = &cobra.Command{
	Use:   "cat",
	Short: "`dedao-dl cat` 获取购买过的课程分类",
	Long:  `使用 dedao-dl type 获取购买过的课程分类`,
	Run: func(cmd *cobra.Command, args []string) {
		app.CourseType()
	},
}

func init() {
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(courseTypeCmd)

	courseCmd.PersistentFlags().StringVarP(&cList, "list", "l", "all", "已购课程列表")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}
