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
	Long:    `获取当前登录用户信息。默认展示：昵称(nickname)、用户标识(uid_hazy)、今日学习时长(today_study_time)、连续学习天数(study_serial_days)、是否会员(is_v)、是否讲师(is_teacher)；使用 --json 可查看完整字段（如头像 avatar、会员信息 vip_user）。`,
	Example: "dedao-dl user -h\n" + "dedao-dl user\n" + "dedao-dl --json user",
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
		table.Append([]string{"结果", "接口返回为空，请使用 --json 排查"})
		table.Render()
		return nil
	}
	table.Append([]string{"昵称(nickname)", info.Nickname})
	table.Append([]string{"用户标识(uid_hazy)", info.UIDHazy})
	table.Append([]string{"今日学习时长(today_study_time)", strconv.Itoa(info.TodayStudyTime)})
	table.Append([]string{"连续学习天数(study_serial_days)", strconv.Itoa(info.StudySerialDays)})
	table.Append([]string{"是否会员(is_v)", strconv.Itoa(info.IsV)})
	table.Append([]string{"是否讲师(is_teacher)", strconv.Itoa(info.IsTeacher)})
	table.Render()
	return nil
}
