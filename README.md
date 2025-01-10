# GoExpert - Clean Architecture

Para rodar o projeto:

```bash
# Iniciar o banco de dados
docker-compose up -d db

# Copiar o arquivo de configuração
cp configs/local.template.env cmd/server/.env

# Buildar o projeto
cd cmd/server
go build

# Rodar o projeto
./server

# Para testar o endpoint de criação de ordem:

curl -X POST http://localhost:8000/order -d '{
  "id":"b",
  "price": 100.5,
  "tax": 0.5
}' -H "Content-Type: application/json"
```