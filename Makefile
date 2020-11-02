GOARCH=amd64
BINARY=dessego

build:
	go build -ldflags="-s -w" -o bin/${BINARY}-linux-${GOARCH} ./cmd/server/main.go

lint:
	golangci-lint run ./cmd/... ./internal/...

deps:
	go mod verify && \
	go mod tidy && \
	go mod vendor && \
	modvendor -copy="**/*.c **/*.h **/*.m"

.PHONY: pkg build lint deps
