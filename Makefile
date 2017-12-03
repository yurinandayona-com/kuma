all: bin/kuma

bin/kuma: generate
	go build -o bin/kuma .

generate:
	go generate -x ./...

fmt:
	gofmt -w -l -s $$(go list ./... | tail -n+2 | sed "s#^$$(go list .)/##")
