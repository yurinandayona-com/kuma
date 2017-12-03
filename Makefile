.PHONY: all generate fmt

all: bin/kuma

bin/kuma: generate
	go build -o bin/kuma .

generate:
	go generate -x ./...

fmt:
	gofmt -w -l -s $$(git ls-files '*.go' | grep -Ev '^vendor/')
