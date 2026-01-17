package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/config"
)

var freeCmd = &cobra.Command{
	Use:   "free",
	Short: "获取首页免费专区课程",
	Long:  `使用 dedao-dl free 获取首页免费专区课程列表`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return freeList()
	},
}

func init() {
	rootCmd.AddCommand(freeCmd)
}

func freeList() (err error) {
	service := config.Instance.ActiveUserService()
	list, err := service.SunflowerResourceList()
	if err != nil {
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"#", "ID(Enid)", "Type", "Name", "Score", "Intro"})

	for i, p := range list.List {
		table.Append([]string{
			strconv.Itoa(i),
			p.Enid,
			strconv.Itoa(p.ProductType),
			p.Name,
			fmt.Sprintf("%.2f", p.Score),
			p.Intro,
		})
	}
	table.Render()
	return
}
