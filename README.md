# GoExpert - Clean Architecture

Para rodar o projeto:

```bash
# Iniciar o banco de dados e o serviço de mensageria
docker-compose up

# Copiar o arquivo de configuração
cp configs/local.template.env cmd/server/.env

# Migrar o banco de dados
make migrate

# Buildar e executar o projeto
make run

# Starting web server on port :8000
# Starting gRPC server on port 50051
# Starting GraphQL server on port 8080
```

## HTTP

O serviço HTTP rodará em http://localhost:8000.

- Para criar uma ordem:

  ```bash
  curl -X POST http://localhost:8000/order -d '{
    "id":"b",
    "price": 100.5,
    "tax": 0.5
  }' -H "Content-Type: application/json"
  ```

- Para listar as ordens criadas:

  ```bash
  curl http://localhost:8000/order
  ```

## GraphQL

Para teste com GraphQL, acessar o playground em http://localhost:8080 e rodar:

- Para criar uma ordem:

  ```graphql
  mutation CreateOrder {
    createOrder(input: {id:"b", Price: 100.5, Tax: 0.5}) {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

- Para listar as ordems criadas:

  ```graphql
  query ListOrders {
    orders {
      id
      Price
      Tax
      FinalPrice
    }
  }
  ```

## gRPC

Rodar o client:

```bash
evans -r repl
package pb

# Criar ordem
service OrderService
call CreateOrder

# id (TYPE_STRING) => fgfff
# price (TYPE_FLOAT) => 501
# tax (TYPE_FLOAT) => 0.2
# {
#   "finalPrice": 501.2,
#   "id": "fgfff",
#   "price": 501,
#   "tax": 0.2
# }

# Listar ordens
call ListOrders
# {
#   "orders": [
#     {
#       "finalPrice": 101,
#       "id": "abc",
#       "price": 100.5,
#       "tax": 0.5
#     },
#     ...
#   ]
# }
```


## Requerimentos

- protobuf-compiler (https://grpc.io/docs/protoc-installation/)
- protoc-gen-go (https://grpc.io/docs/languages/go/quickstart/)
  - Nota: precisei executar: `go get -u google.golang.org/grpc` para resolver o erro em `SupportPackageIsVersion9` do .pb.go.
- protoc-gen-go-grpc (https://grpc.io/docs/languages/go/quickstart/)
- evans (https://github.com/ktr0731/evans)