package cmd

import (
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/cmd/app"
	"github.com/yann0917/dedao-dl/services"
)

func init() {
	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:     "user",
	Short:   "获取当前登录用户信息",
	Long:    `使用 dedao-dl user 获取 /api/pc/user/info 返回的用户信息。`,
	Example: "dedao-dl user\n" + "dedao-dl --json user",
	PreRunE: AuthFunc,
	RunE: func(cmd *cobra.Command, args []string) error {
		info, err := app.User()
		if err != nil {
			return err
		}
		if outputJSON {
			return printJSON(info)
		}
		return renderUserTable(info)
	},
}

func renderUserTable(info *services.User) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.Header([]string{"字段", "值"})
	if info == nil {
		table.Append([]string{"result", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}
	table.Append([]string{"nickname", info.Nickname})
	table.Append([]string{"uid_hazy", info.UIDHazy})
	table.Append([]string{"today_study_time", strconv.Itoa(info.TodayStudyTime)})
	table.Append([]string{"study_serial_days", strconv.Itoa(info.StudySerialDays)})
	table.Append([]string{"is_v", strconv.Itoa(info.IsV)})
	table.Append([]string{"is_teacher", strconv.Itoa(info.IsTeacher)})
	table.Render()
	return nil
}
