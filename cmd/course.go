package cmd

import (
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var cType string

var courseCmd = &cobra.Command{
	Use:       "course",
	Short:     "`dedao-dl course`获取购买过的课程",
	Long:      `使用 dedao-dl course 获取购买过的课程，课程分类，章节等`,
	Args:      cobra.OnlyValidArgs,
	ValidArgs: []string{"all", "bauhinia", "ebook", "compass"},
	Example:   "dedao-dl course -t ebook",
	Run: func(cmd *cobra.Command, args []string) {
		app.CourseList(cType)
	},
}

var courseTypeCmd = &cobra.Command{
	Use:   "cat",
	Short: "`dedao-dl cat` 获取购买过的课程分类",
	Long:  `使用 dedao-dl cat 获取购买过的课程分类`,
	Run: func(cmd *cobra.Command, args []string) {
		app.CourseType()
	},
}

func init() {
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(courseTypeCmd)

	rootCmd.PersistentFlags().StringVarP(&cType, "type", "t", "bauhinia", "课程类型(all, bauhinia, ebook, compass)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
}

func argFuncs(funcs ...cobra.PositionalArgs) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		for _, f := range funcs {
			err := f(cmd, args)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
