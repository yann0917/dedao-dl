package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/config"
	"github.com/yann0917/dedao-dl/utils"
)

var freeCmd = &cobra.Command{
	Use:   "free",
	Short: "获取免费专区课程",
	Long:  `使用 dedao-dl free 获取免费专区课程列表`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return freeInfo(args[0])
		}
		return freeList()
	},
}

func init() {
	rootCmd.AddCommand(freeCmd)
}

func freeInfo(enid string) (err error) {
	service := config.Instance.ActiveUserService()
	info, err := service.CourseInfo(enid)
	if err != nil {
		return
	}

	out := os.Stdout
	table := tablewriter.NewWriter(out)

	fmt.Fprint(out, "专栏名称："+info.ClassInfo.Name+"\n")
	fmt.Fprint(out, "专栏作者："+info.ClassInfo.LecturerNameAndTitle+"\n")
	if info.ClassInfo.PhaseNum == 0 {
		fmt.Fprint(out, "共"+strconv.Itoa(info.ClassInfo.CurrentArticleCount)+"讲\n")
	} else {
		fmt.Fprint(out, "更新进度："+strconv.Itoa(info.ClassInfo.CurrentArticleCount)+
			"/"+strconv.Itoa(info.ClassInfo.PhaseNum)+"\n")
	}
	fmt.Fprint(out, "课程亮点："+info.ClassInfo.Highlight+"\n")
	fmt.Fprintln(out)

	table.Header([]string{"#", "ID", "章节", "讲数", "更新时间", "是否更新完成"})

	if len(info.ChapterList) > 0 {
		for i, p := range info.ChapterList {
			isFinished := "❌"
			if p.IsFinished == 1 {
				isFinished = "✔"
			}
			table.Append([]string{strconv.Itoa(i),
				p.IDStr, p.Name, strconv.Itoa(p.PhaseNum),
				utils.Unix2String(int64(p.UpdateTime)),
				isFinished,
			})
		}
	} else if len(info.FlatArticleList) > 0 {
		isFinished := "❌"
		if info.ClassInfo.IsFinished == 1 {
			isFinished = "✔"
		}
		for i, p := range info.FlatArticleList {
			table.Append([]string{strconv.Itoa(i),
				p.IDStr, "-", p.Title,
				utils.Unix2String(int64(p.UpdateTime)),
				isFinished,
			})
		}
		if info.HasMoreFlatArticleList {
			fmt.Fprint(out, "⚠️  更多文章请使用 article -i 查看文章列表...\n")
		}
	}
	table.Render()
	return
}

func freeList() (err error) {
	service := config.Instance.ActiveUserService()
	list, err := service.SunflowerResourceList()
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "ID(Enid)", "Name", "Score", "Intro"})

	for i, p := range list.List {
		table.Append([]string{
			strconv.Itoa(i),
			p.Enid,
			p.Name,
			fmt.Sprintf("%.2f", p.Score),
			p.Intro,
		})
	}
	table.Render()
	return
}
