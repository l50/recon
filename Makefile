all: build

build:
	go build -v
	go test ./...
	go vet
	find . -not \( \( -wholename './.git' -o -wholename '*/vendor/*' \) -prune \) -name '*.go' | xargs gofmt -s -d
	go install

init:
	go mod download