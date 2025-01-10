FROM golang:latest AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/server
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-w -s" -o server .

FROM scratch
COPY --from=builder /app/cmd/server/server .
CMD ["./server"]
