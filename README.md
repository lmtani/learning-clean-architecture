# GoExpert - Clean Architecture

Para rodar o projeto:

```bash
# Iniciar o banco de dados
docker-compose up -d db

# Copiar o arquivo de configuração
cp configs/local.template.env cmd/server/.env

# Buildar o projeto
make build

# Rodar o projeto
./server

# Para testar o endpoint de criação de ordem:

curl -X POST http://localhost:8000/order -d '{
  "id":"b",
  "price": 100.5,
  "tax": 0.5
}' -H "Content-Type: application/json"
```

Para teste com GraphQL, acessar o playground em http://localhost:8080 e rodar a query:

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


## Requerimentos

- protobuf-compiler (https://grpc.io/docs/protoc-installation/)
- protoc-gen-go (https://grpc.io/docs/languages/go/quickstart/)
  - Nota: precisei executar: `go get -u google.golang.org/grpc` para resolver o erro em `SupportPackageIsVersion9` do .pb.go.
- protoc-gen-go-grpc (https://grpc.io/docs/languages/go/quickstart/)