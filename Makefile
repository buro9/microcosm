.PHONY: all
all: microcosm-ui

.PHONY: microcosm-ui
microcosm-ui:
	go install github.com/microcosm-cc/microcosm-ui/cmd/microcosm-ui

.PHONY: vendor
vendor:
	# Core dependencies
	-gvt fetch github.com/bep/inflect
	-gvt fetch github.com/dustin/go-humanize
	-gvt fetch github.com/eknkc/amber
	-gvt fetch github.com/gregjones/httpcache
	-gvt fetch github.com/oxtoacart/bpool
	-gvt fetch github.com/pressly/chi
	-gvt fetch github.com/spf13/afero
	-gvt fetch github.com/spf13/cast
	-gvt fetch github.com/spf13/hugo/bufferpool
	-gvt fetch github.com/spf13/hugo/helpers
	-gvt fetch github.com/yosssi/ace
	-gvt fetch gopkg.in/oleiade/reflections.v1