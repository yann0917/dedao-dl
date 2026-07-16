package cmd

import (
	"github.com/spf13/cobra"
	webserver "github.com/yann0917/dedao-dl/web/server"
)

var (
	webHost        string
	webPort        int
	webOpenBrowser bool
)

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web UI 与 API 服务",
	Long:  "启动基于 gin 的 Web API，并打开内置的 Web 页面。",
	RunE: func(cmd *cobra.Command, args []string) error {
		return webserver.Start(webserver.Options{
			Host:        webHost,
			Port:        webPort,
			OpenBrowser: webOpenBrowser,
		})
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
	webCmd.Flags().StringVar(&webHost, "host", "127.0.0.1", "Web 服务监听地址")
	webCmd.Flags().IntVar(&webPort, "port", 17878, "Web 服务监听端口")
	webCmd.Flags().BoolVar(&webOpenBrowser, "open", true, "启动后自动打开浏览器")
}
