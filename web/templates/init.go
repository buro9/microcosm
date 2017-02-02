package templates

import "github.com/oxtoacart/bpool"

func init() {
	// Used within templates.go
	bufpool = bpool.NewBufferPool(64)
}
