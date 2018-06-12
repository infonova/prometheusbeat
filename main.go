package main

import (
	"os"

	"github.com/infonova/prometheusbeat/cmd"

	_ "github.com/infonova/prometheusbeat/include"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
