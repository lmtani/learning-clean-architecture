wire:
	cd cmd/server/ && wire && cd -

generate:
	go run github.com/99designs/gqlgen generate

grpc:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

build: wire generate grpc
	go mod tidy && \
	cd cmd/server && \
	go build ./... && \
	cd -

run: build
	cd cmd/server && \
	./server || \
	cd -
	

.PHONY: wire generate build run grpc
