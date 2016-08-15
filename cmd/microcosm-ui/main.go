package main

import (
	"flag"
	"log"

	"github.com/microcosm-cc/microcosm-ui"
)

func main() {
	ui.RegisterFlags()
	flag.Parse()

	ui.ParseTemplates()

	log.Fatal(ui.ListenAndServe())
}
