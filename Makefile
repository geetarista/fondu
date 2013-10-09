ERROR_COLOR=\033[31;01m
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
WARN_COLOR=\033[33;01m

DEPS = $(go list -f '{{range .TestImports}}{{.}} {{end}}' ./...)
VERSION=$(shell go run $(shell ls *.go | grep -v "_test\.go") -v)

all:
	@echo "$(OK_COLOR)==> Building$(NO_COLOR)"
	@bash --norc -i ./script/build.sh $(VERSION)

clean:
	rm -rf bin/ data/

deps:
	@echo "$(OK_COLOR)==> Installing dependencies$(NO_COLOR)"
	@go get -d -v ./...
	@echo $(DEPS) | xargs -n1 go get -d

format:
	go fmt ./...

updatedeps:
	@echo "$(OK_COLOR)==> Updating all dependencies$(NO_COLOR)"
	@go get -d -v -u ./...
	@echo $(DEPS) | xargs -n1 go get -d -u

release: all
	@echo "$(OK_COLOR)==> Releasing version $(VERSION)$(NO_COLOR)"
	@bash --norc -i ./release.sh $(VERSION)

test: deps format
	go test -cover
	@rm -rf data
	@if [ -f src.test ]; then rm src.test; fi

.PHONY: all clean deps format release test updatedeps
