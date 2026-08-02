# clean-architecture

Desafio Go-Expert sobre Clean Architecture: listagem de Orders exposta simultaneamente via **REST**, **gRPC** e **GraphQL**, todas consumindo o mesmo `ListOrdersUseCase`.

## Como rodar

Pré-requisito: Docker e Docker Compose instalados.

```bash
docker compose up
```

Esse único comando:

1. Sobe o banco de dados MySQL (com healthcheck).
2. Aguarda o banco ficar saudável (`depends_on: condition: service_healthy`) e, como camada extra de segurança, a própria aplicação faz retry de conexão ao subir — evitando qualquer race condition de inicialização.
3. Executa as migrations automaticamente (criação da tabela `orders`) antes de iniciar os servidores.
4. Sobe os três servidores (REST, gRPC e GraphQL) em paralelo.

Não é necessário rodar nenhum comando manual além do `docker compose up`.

## Portas

| Serviço  | Protocolo | Porta | Endpoint                                    |
|----------|-----------|-------|----------------------------------------------|
| Web      | REST      | 8000  | `POST /order`, `GET /order`                   |
| gRPC     | gRPC      | 50051 | serviço `OrderService` (`CreateOrder`, `ListOrders`) |
| GraphQL  | GraphQL   | 8080  | `POST /query` (playground em `GET /`)         |
| MySQL    | MySQL     | 3306  | banco `orders`                                |

## Testando

O arquivo [api.http](api.http) na raiz já contém requisições REST e GraphQL prontas para:
- Criar uma Order (popular o banco).
- Listar as Orders.

### REST
```
POST http://localhost:8000/order
GET  http://localhost:8000/order
```

### GraphQL
Abra `http://localhost:8080` para o playground, ou envie requisições para `http://localhost:8080/query`:
```graphql
mutation {
  createOrder(input: { id: "1", Price: 100, Tax: 10 }) {
    id
    Price
    Tax
    FinalPrice
  }
}

query {
  ListOrders {
    id
    Price
    Tax
    FinalPrice
  }
}
```

### gRPC
Use [grpcurl](https://github.com/fullstorydev/grpcurl) (reflection está habilitada no servidor):
```bash
grpcurl -plaintext localhost:50051 list
grpcurl -plaintext -d '{"id": "1", "price": 100, "tax": 10}' localhost:50051 pb.OrderService/CreateOrder
grpcurl -plaintext -d '{}' localhost:50051 pb.OrderService/ListOrders
```

## Arquitetura

- `internal/entity`: entidade `Order` e regras de validação.
- `internal/usecase`: `CreateOrderUseCase` e `ListOrdersUseCase`, dependendo apenas da interface `OrderRepositoryInterface`.
- `internal/infra/database`: implementação do repositório em MySQL.
- `internal/infra/web`: handlers REST.
- `internal/infra/grpc`: proto, código gerado e serviço gRPC.
- `internal/infra/graph`: schema, código gerado (gqlgen) e resolvers GraphQL.
- `cmd/ordersystem/main.go`: composição das dependências e inicialização dos três servidores.
- `migrations/`: migrations SQL (aplicadas via `golang-migrate` na inicialização da aplicação).

## Testes

```bash
go test ./...
```

Cobrem a entidade, os use cases (com mock de repositório) e o repositório (com SQLite em memória).
