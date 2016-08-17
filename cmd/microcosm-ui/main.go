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

	// Listen and wait for errors (none should ever be received, so we run
	// forever)
	errs := ui.ListenAndServe()
	select {
	case err := <-errs:
		log.Fatal(err)
	}
}
