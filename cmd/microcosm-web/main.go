package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/buro9/microcosm/web/opts"
	"github.com/buro9/microcosm/web/server"
	"github.com/buro9/microcosm/web/templates"
)

func main() {
	opts.RegisterFlags()
	flag.Parse()
	if err := opts.ValidateFlags(); err != nil {
		fmt.Println(err.Error())
		flag.PrintDefaults()
		os.Exit(1)
	}

	// Compile all templates, these are .MustCompile and so will prevent later
	// runtime errors relating to badly formatted templates
	templates.Compile()

	// Listen and wait for errors (none should ever be received, so we should
	// run forever)
	errs := server.ListenAndServe()
	select {
	case err := <-errs:
		log.Fatal(err)
	}
}
