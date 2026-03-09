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

Fase 1 
1	Definir escopo (rotas, serviços, contratos mínimos)
2	Desenhar arquitetura (README ou Mermaid)
3	Subir 3 serviços mockados (3 portas, JSON fixo)
4	Gateway básico (servidor + rota /health ou similar)
5	Rota /dashboard retornando resposta fixa

Sim. Fase 1 e Fase 2 estão fechadas; o próximo passo é a Fase 3 — Concorrência e resiliência.
Resumo do que entra nela:
Ordem	Item	O que fazer
1	Propagar contexto	Garantir que o context.Context da request vá do handler → use case → clients (você já usa r.Context() e Execute(ctx); conferir se todos os clients recebem e usam o ctx nas chamadas HTTP).
2	Adicionar timeout	Criar um context com timeout (ex.: 5s) no handler ou no use case, derivado do r.Context(), e passá-lo para o use case; assim as chamadas não ficam penduradas para sempre.
3	Paralelizar chamadas	No use case, em vez de chamar GetUsers → GetOrders → GetBillings em sequência, chamar os três em paralelo com goroutines + errgroup ou sync.WaitGroup, e só depois agregar o resultado.
4	Tratar falhas	Decidir o que fazer quando um dos 3 serviços falha: retornar erro 502, ou degradar (ex.: retornar só os que responderam e indicar falha parcial).
5	Degradação graciosa	Se um serviço falhar, retornar o que der (ex.: users e orders ok, billings em falha) com status 200 ou 207 e um campo indicando quais partes falharam.
6	Circuit breaker	Adicionar um circuit breaker por client (ex.: após N falhas consecutivas, parar de chamar por X segundos e retornar erro ou fallback direto).
7	Simular falhas	Ter como derrubar ou atrasar um mock (ex.: query param ou rota nos mocks) para testar timeout, circuit breaker e degradação.
Ordem sugerida: primeiro contexto + timeout, depois paralelizar, em seguida tratamento de falhas e degradação, e por último circuit breaker e simulação de falhas. Assim você vai evoluindo a resiliência em camadas.