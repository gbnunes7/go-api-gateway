## Objetivo desse projeto de estudo.

Construir um API Gateway em Go que agregue dados de 3 serviços com concorrência, resiliência e observabilidade.

## Escopo

Unico endpoint disponivel vai ser /dashboard, de onde vao vir dos 3 clientes externos.
Clientes -> Users / Orders / Billings

## Contratos

Users: {
  id: uuid (PK)
  name: varchar(255)
  email: varchar(255)

  orderId: uuid (FK)
  billingId: uuid (FK)
}

Orders: {
  id: uuid (PK)
  totalPrice: int
  createdAt: string

  userId: uuid (FK)
}

Billings: {
  id: uuid (PK)
  paymentType: varchar(255)
  paidValue: int

  orderId: uuid (FK)
}


## Arquitetura

Client  →  GET /dashboard  →  Gateway (:8080)
                                  │
                                  ├→ Users   (:8081)
                                  ├→ Orders  (:8082)
                                  └→ Billing (:8083)
                                  │
                                  ↓
Client  ←  JSON agregado   ←  Gateway
