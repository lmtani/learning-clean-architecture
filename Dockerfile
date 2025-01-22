FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN apt-get update && apt-get install -y protobuf-compiler \
    && go install google.golang.org/protobuf/cmd/protoc-gen-go@latest \
    && go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest \
    && go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/google/wire/cmd/wire@latest \
    && go mod download 

WORKDIR /app/cmd/server
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM scratch
COPY --from=builder /app/cmd/server/server .
CMD ["./server"]
