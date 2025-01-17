wire:
	cd cmd/server/ && wire && cd -

generate:
	go run github.com/99designs/gqlgen generate

grpc:
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

sqlc:
	sqlc generate

create-migration:
	migrate create -ext sql -dir internal/infra/database/psql/migrations -seq init

migrate:
	migrate -path internal/infra/database/psql/migrations -database "postgresql://root:root@localhost:5432/orders?sslmode=disable" -verbose up

build: sqlc grpc generate wire
	go mod tidy && \
	cd cmd/server && \
	go build ./... && \
	cd -

run: build
	cd cmd/server && \
	./server || \
	cd -
	

.PHONY: wire generate build run grpc sqlc migrate
