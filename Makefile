.PHONY: all generate fmt

all: bin/kuma

bin/kuma: generate
	go build -o bin/kuma .

generate:
	go generate -x ./...

fmt:
	gofmt -w -l -s $$(git ls-files '*.go' | grep -Ev '^vendor/')

loc:
	cloc $$(git ls-files | grep -Ev '^vendor/|.pb.go$$')
