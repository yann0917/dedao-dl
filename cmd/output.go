package cmd

import (
	"os"

	jsoniter "github.com/json-iterator/go"
)

func printJSON(v interface{}) error {
	enc := jsoniter.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}
