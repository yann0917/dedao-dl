package main

import (
	"fmt"

	"github.com/yann0917/dedao-dl/cmd"
	"github.com/yann0917/dedao-dl/config"
)

func init() {

	err := config.Instance.Init()
	if err != nil {
		fmt.Println(err)
	}
}
func main() {
	cmd.Execute()
}
