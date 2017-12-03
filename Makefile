.PHONY: all

all: bin/kuma

bin/kuma: generate
	go build -o bin/kuma .

#
# Tasks
.PHONY: generate fmt dep

generate:
	go generate -x ./...

fmt:
	gofmt -w -l -s $$(git ls-files '*.go' | grep -Ev '^vendor/')

dep:
	dep ensure -v
	dep prune -v

#
# Utilities
.PHONY: loc tree

loc:
	cloc $$(git ls-files | grep -Ev '^vendor/|.pb.go$$')

tree:
	tree -I vendor -N
