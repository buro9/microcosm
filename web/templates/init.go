package templates

import "github.com/oxtoacart/bpool"

func init() {
	// Used within templates.go
	bufpool = bpool.NewBufferPool(64)

	// produce the definitions of which templates include which other templates
	loadDefinitions()
}
