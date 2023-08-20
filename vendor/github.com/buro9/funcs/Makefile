.PHONY: all
all: install

.PHONY: install
install:
	$Q go install

.PHONY: test
test:
	$Q go test -race ./...
	$Q go vet ./...
	$Q golint ./...

.PHONY: cover
cover: clean
	@echo "NOTE: make cover does not exit 1 on failure, don't use it to check for tests success"
	$(if $V,@echo "-- go test -coverpkg=./... -coverprofile=cover/... ./...")
	@for MOD in $(allpackages); do \
		go test -coverpkg=`echo $(allpackages)|tr " " ","` \
			-coverprofile=cover/`echo $$MOD|tr "/" "_"`.out \
			$$MOD 2>&1 | grep -v "no packages being tested depend on"; \
	done
	$Q gocovmerge cover/*.out > cover/all.merged
	$Q go tool cover -html cover/all.merged
	@echo ""
	@echo "=====> Total test coverage: <====="
	@echo ""
	$Q go tool cover -func cover/all.merged

.PHONY: list
list:
	@echo $(allpackages)

# memoize allpackages, so that it's executed only once and only if used
_allpackages = $(shell ( go list ./... ))
allpackages = $(if $(__allpackages),,$(eval __allpackages := $$(_allpackages)))$(__allpackages)

.PHONY: clean
clean:
	rm -f ./cover/*.out
	rm -f ./cover/all.merged
