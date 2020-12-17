package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "dedao-dl",
	Short: "dedao-dl is a very fast dedao app course article download tools",
	Long: `A Fast dedao app course article download tools built with
		love by spf13 and friends in Go.
		Complete documentation is available at http://hugo.spf13.com`,
}

// Execute exec cmd
func Execute() error {
	return rootCmd.Execute()
}
