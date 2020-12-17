package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/yann0917/dedao-dl/services"
)

// Login login
type Login struct {
	phone         string
	password      string
	GAT           string `json:"gat"`
	ISID          string `json:"isid"`
	GuardDeviceID string `json:"guard_device_id"`
	SID           string `json:"sid"`
	AcwTc         string `json:"acw_tc"`
	Iget          string `json:"iget"`
	Token         string `json:"token"`
	CookieStr     string `json:"cookieStr"`
}

// IsByCookie cookie login
func (l *Login) IsByCookie() bool {
	return l.GAT != "" && l.ISID != "" && l.SID != ""
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "use `dedao-dl login` to login https://www.dedao.cn",
	Long:  `use dedao-dl login to login https://www.dedao.cn`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("login cmd")
	},
}

func init() {
	rootCmd.AddCommand(loginCmd)
	var cookie services.CookieOptions
	loginCmd.Flags().StringVarP(&cookie.CookieStr, "cookie", "c", "", "cookie from www.dedao.cn")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
