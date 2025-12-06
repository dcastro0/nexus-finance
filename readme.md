# Nexus Finance API

API REST desenvolvida em Go para gerenciamento de transações financeiras. O projeto simula operações bancárias essenciais, como criação de contas, depósitos e transferências entre contas, focando em escalabilidade e boas práticas de engenharia de software.

## 🚀 Tecnologias

* **Linguagem:** Go (Golang)
* **Roteamento:** [Chi](https://github.com/go-chi/chi)
* **ORM:** [Gorm](https://gorm.io/)
* **Banco de Dados:** PostgreSQL
* **Configuração:** [Viper](https://github.com/spf13/viper)
* **Infraestrutura:** Docker & Docker Compose

## 📂 Arquitetura

O projeto segue os princípios da **Clean Architecture**, dividindo responsabilidades para garantir desacoplamento e facilidade de manutenção:

```text
nexus-finance/
├── cmd/api/         # Entrypoint da aplicação (main.go)
├── configs/         # Gerenciamento de variáveis de ambiente
├── internal/
│   ├── domain/      # Entidades do domínio e interfaces (Core)
│   ├── infra/       # Implementações (Banco de dados, HTTP Handlers, Repositórios)
│   └── usecase/     # Regras de negócio (Casos de uso)
├── pkg/             # Pacotes compartilhados
└── docker-compose.yml