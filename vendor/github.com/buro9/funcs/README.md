# Funcs

Funcs is a package that provides a `template.FuncMap` for use in your templates.

## Usage

Either just download a copy using `go get`:

```bash
go get -u github.com/buro9/funcs
```

Or vendor it within your solution (using [`gvt`](https://github.com/FiloSottile/gvt) here as an example):

```bash
gvt fetch github.com/buro9/funcs
```

Once you have pulled the code, you are able to import it.

The following example is taken from the [text/template documentation](https://golang.org/pkg/text/template/#example_Template_func) but has been modified to use the `template.FuncMap` provided by this package:

```go
package main

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/buro9/funcs"
)

func main() {
	// A simple template definition to test our function.
	// We print the input text several ways:
	// - the original
	// - title-cased
	// - title-cased and then printed with %q
	// - printed with %q and then title-cased.
	const templateText = `
Input: {{printf "%q" .}}
Output 0: {{title .}}
Output 1: {{title . | printf "%q"}}
Output 2: {{printf "%q" . | title}}
`

	// Create a template, add the function map, and parse the text.
	tmpl, err := template.New("titleTest").Funcs(funcs.Map()).Parse(templateText)
	if err != nil {
		log.Fatalf("parsing: %s", err)
	}

	// Run the template to verify the output.
	err = tmpl.Execute(os.Stdout, "the go programming language")
	if err != nil {
		log.Fatalf("execution: %s", err)
	}
}
```

## Contributing

By contributing code to this package, you are agreeing to release it under the [MIT License] (https://github.com/buro9/funcs/LICENSE.md).

Package layout:

* `funcs/` only provides the `template.FuncMap`
* `funcs/{funcName|funcGroupName}/` provides one or more funcs that is then included in the `template.FuncMap` provided by `funcs.Map()`

When contributing funcs, the important thing is that child directories contain either an individual func or a logical group of funcs are fully self-contained and tested. Only include more than one func per child directory where the funcs share common code that couples the implementation, otherwise provide more child directories.

The goal is to allow other developers to treat this package as a library of funcs, wherein the developer can cherry pick specific funcs to use by referencing the leaf directories, or can just use the entire library of funcs by calling `funcs.Map()` at the package root.

Your funcs should not have dependencies beyond core Go packages if at all possible. Developers who wish to include funcs in their work will naturally desire as few dependencies as possible.

How to contribute:

1. Fork this repo
2. Add your changes
3. Ensure you have tests and that they pass [`go test -race`](https://golang.org/doc/articles/race_detector.html)
4. Ensure your documentation passes [`golint`](https://github.com/golang/lint) and [`go vet`](https://golang.org/cmd/vet/) with zero warnings or errors
5. Create a pull request describing changes

Failure to include tests and documentation will see your PR closed. Including funcs that can only work with knowledge of your web app (i.e. link URL rewriting in HTML templates that rewrite to some path only your app knows about) will see your PR closed. Just remember, the point of this package is that everyone can use it.

Simple funcs, easy to test and understand are better than complex funcs. If a simple func does not do what you wish, feel free to contribute an improved version of it.

Note: 

Checklist:

* Your code is in a child directory
* Your func(s) all have tests
* Your func(s) are documented
* Your func(s) are generalised and will work with any project that implements templates (they do not require any knowledge of your application)
* You have reduced dependencies as far as possible
* You accept that any contributions to be released under the MIT License

We do *NOT* vendor dependencies, these should be managed by your vendoring tools of choice when you include this package in your code.