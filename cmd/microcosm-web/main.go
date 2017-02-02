package main

import (
	"flag"
	"log"

	"github.com/buro9/microcosm/web/opts"
	"github.com/buro9/microcosm/web/server"
	"github.com/buro9/microcosm/web/templates"
)

func main() {
	opts.RegisterFlags()
	flag.Parse()

	templates.Load()

	// Listen and wait for errors (none should ever be received, so we run
	// forever)
	errs := server.ListenAndServe()
	select {
	case err := <-errs:
		log.Fatal(err)
	}
}
