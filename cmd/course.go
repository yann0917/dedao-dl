package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var (
	cType     string
	classID   int
	articleID int
	bookID    int
	compassID int
)

var courseTypeCmd = &cobra.Command{
	Use:     "cat",
	Short:   "获取课程分类",
	Long:    `使用 dedao-dl cat 获取课程分类`,
	Example: "dedao-dl cat",
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		app.CourseType()
	},
}

var courseCmd = &cobra.Command{
	Use:     "course",
	Short:   "获取我购买过课程",
	Long:    `使用 dedao-dl course 获取我购买过的课程`,
	Example: "dedao-dl course",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("classID", classID)
		if classID > 0 {
			app.CourseInfo(classID)
			return
		}
		app.CourseList("bauhinia")
	},
}

var ebookCmd = &cobra.Command{
	Use:     "ebook",
	Short:   "获取我的电子书架",
	Long:    `使用 dedao-dl ebook 获取我的电子书架`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ebook",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("bookID", bookID)
		if bookID > 0 {
			app.EbookDetail(bookID)
			return
		}
		app.CourseList("ebook")
	},
}

var compassCmd = &cobra.Command{
	Use:     "ace",
	Short:   "获取我的锦囊",
	Long:    `使用 dedao-dl ace 获取我的锦囊`,
	Args:    cobra.OnlyValidArgs,
	Example: "dedao-dl ace",
	PreRunE: AuthFunc,
	Run: func(cmd *cobra.Command, args []string) {
		if compassID > 0 {
			return
		}
		app.CourseList("compass")
	},
}

func init() {
	rootCmd.AddCommand(courseTypeCmd)
	rootCmd.AddCommand(courseCmd)
	rootCmd.AddCommand(ebookCmd)
	rootCmd.AddCommand(compassCmd)
	courseCmd.PersistentFlags().IntVarP(&classID, "id", "i", 0, "课程 ID，获取课程信息")
	ebookCmd.PersistentFlags().IntVarP(&bookID, "id", "i", 0, "电子书ID")
	compassCmd.PersistentFlags().IntVarP(&compassID, "id", "i", 0, "锦囊 ID")
	// rootCmd.PersistentFlags().StringVarP(&cType, "type", "t", "bauhinia", "课程类型(all, bauhinia, ebook, compass")

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
