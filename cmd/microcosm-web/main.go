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

	// MustCompile all templates as a compile error is more preferable than a
	// runtime error
	templates.Load()

	// Listen and wait for errors (none should ever be received, so we should
	// run forever)
	errs := server.ListenAndServe()
	select {
	case err := <-errs:
		log.Fatal(err)
	}
}
