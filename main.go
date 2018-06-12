package main

import (
	"os"

	"github.com/infonova/prometheusbeat/cmd"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
