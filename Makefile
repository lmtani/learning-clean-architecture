.SHELLFLAGS = -e

wire:
	cd cmd/server && wire

generate:
	go run github.com/99designs/gqlgen generate

grpc:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

build: wire generate grpc
	go mod tidy
	go build ./...

.PHONY: wire generate build
