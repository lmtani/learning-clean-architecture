wire:
	@echo "Running wire"
	cd cmd/server/ && wire && cd -

generate:
	@echo "Running gqlgen"
	go run github.com/99designs/gqlgen generate

grpc:
	@echo "Running protoc"
	protoc --go_out=. --go-grpc_out=. internal/infra/grpc/protofiles/order.proto

sqlc:
	@echo "Running sqlc"
	sqlc generate

create-migration:
	@echo "Creating migration"
	migrate create -ext sql -dir internal/infra/database/psql/migrations -seq init

migrate:
	@echo "Running migration"
	migrate -path internal/infra/database/psql/migrations -database "postgresql://root:root@localhost:5432/orders?sslmode=disable" -verbose up

migrate-docker:
	docker run -v $(shell pwd)/internal/infra/database/psql/migrations:/migrations \
	    --network host \
		migrate/migrate -path=/migrations/ -database "postgresql://root:root@localhost:5432/orders?sslmode=disable" up

build: sqlc grpc generate wire 
	@echo "Building server"
	go mod tidy && \
	cd cmd/server && \
	go build ./... && \
	cd -

run: build
	@echo "Running server"
	cd cmd/server && \
	./server || \
	cd -
	
test:
	@echo "Running tests with coverage"
	go test -cover ./...

.PHONY: wire generate build run grpc sqlc migrate migrate-docker test
