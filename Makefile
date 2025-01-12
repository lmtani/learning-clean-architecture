wire:
	cd cmd/server && wire

generate:
	go run github.com/99designs/gqlgen generate

build: wire generate
	go mod tidy
	go build ./...

.PHONY: wire generate build
