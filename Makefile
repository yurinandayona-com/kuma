all: bin/kuma

bin/kuma: generate
	go build -o bin/kuma .

generate:
	go generate -x ./...

fmt:
	go fmt ./...
