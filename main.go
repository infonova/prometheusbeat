package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/infonova/prometheusbeat/beater"
)

func main() {
	err := beat.Run("prometheusbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
