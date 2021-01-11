package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
)

var topicCmd = &cobra.Command{
	Use:     "topic",
	Short:   "获取推荐话题列表",
	Long:    `使用 dedao-dl topic 获取推荐话题列表，也可获取某个话题前 40 条精选话题内容`,
	Args:    cobra.NoArgs,
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {

		if len(topicID) > 0 {
			err := notesList(topicID)
			return err
		}
		err := topicAll()
		return err
	},
}

func init() {
	rootCmd.AddCommand(topicCmd)
	topicCmd.PersistentFlags().StringVarP(&topicID, "id", "i", "", "话题id")
}

func topicAll() (err error) {
	list, err := app.TopicAll()
	if err != nil {
		return
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"#", "ID", "话题名称", "话题参与数"})
	table.SetAutoWrapText(false)

	for i, p := range list.List {
		table.Append([]string{strconv.Itoa(i),
			p.TopicIDHazy, p.Name,
			strconv.Itoa(p.NotesCount),
		})
	}
	table.Render()
	return
}

func notesList(id string) (err error) {
	detail, err := app.TopicDetail(id)
	if err != nil {
		return
	}
	list, err := app.TopicNotesList(id)
	if err != nil {
		return
	}
	out := os.Stdout
	fmt.Fprint(out, "话题名称："+detail.Name+"\n")
	fmt.Fprint(out, "话题名称："+detail.Intro+"\n")
	fmt.Fprintln(out)

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"#", "昵称", "点赞数", "发布时间", "笔记"})
	table.SetAutoWrapText(true)

	for i, p := range list.NoteDetailList {
		table.Append([]string{strconv.Itoa(i),
			p.FPart.NickName,
			strconv.Itoa(p.NoteCount.LikeCount),
			p.FPart.TimeDesc,
			p.FPart.Note,
		})
	}
	table.Render()
	return
}
