.PHONY: all

all: build

#
# Tasks
.PHONY: build generate vet fmt dep test

build: generate
	go build -o bin/kuma ./cmd/...

generate:
	go generate -x ./...

vet:
	go vet ./...

fmt:
	gofmt -w -l -s $$(git ls-files '*.go' | grep -Ev '^vendor/')
	goimports -w $$(git ls-files '*.go' | grep -Ev '^vendor/')

dep:
	dep ensure -v
	dep prune -v

test:
	go test -v ./...

#
# Utilities
.PHONY: loc tree

loc:
	cloc $$(git ls-files | grep -Ev '^vendor/|.pb.go$$')

tree:
	tree -I vendor -N
