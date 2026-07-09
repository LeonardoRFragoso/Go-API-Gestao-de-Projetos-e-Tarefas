# � TaskFlow — Plataforma SaaS de Gestão de Projetos, Squads e Tarefas

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Vue.js-3-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue.js" />
  <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
  <img src="https://img.shields.io/badge/JWT-Auth-FF6B6B?style=for-the-badge&logo=json-web-tokens&logoColor=white" alt="JWT" />
  <img src="https://img.shields.io/badge/REST-API-4ECDC4?style=for-the-badge&logo=api&logoColor=white" alt="REST API" />
</p>

**TaskFlow** é uma plataforma SaaS moderna para gestão de projetos, squads e tarefas. Construída com **Go**, **Vue 3**, **PostgreSQL** e **Docker**, oferece um Kanban intuitivo com drag-and-drop, dashboard com métricas em tempo real, sistema de permissões granulares, gerenciamento de equipes e muito mais.

Perfeita para equipes que buscam uma solução robusta, escalável e fácil de usar para organizar projetos e colaboração.

🔗 **Repositório:** [https://github.com/LeonardoRFragoso/Go-API-Gestao-de-Projetos-e-Tarefas](https://github.com/LeonardoRFragoso/Go-API-Gestao-de-Projetos-e-Tarefas)

## ✨ Features

- 🌙 **Dark Mode** - Tema escuro/claro com transições suaves
- 🌐 **Internacionalização** - Suporte para Português (BR) e Inglês
- 📊 **Dashboard com Gráficos** - Visualização de estatísticas com Chart.js
- 📋 **Kanban Board** - Arrastar e soltar tarefas entre colunas
- � **Equipes/Squads** - Crie equipes, adicione membros e vincule projetos
- 🔔 **Sistema de Notificações** - Notificações em tempo real para tarefas, convites e menções
- 🛡️ **Permissões Granulares** - Controle de acesso por projeto e equipe (lead, admin, member)
- � **Busca e Filtros** - Pesquise tarefas por título, descrição e prioridade
- 📱 **Responsivo** - Menu mobile e bottom navigation
- ⌨️ **Atalhos de Teclado** - Navegação rápida
- 🎯 **Onboarding** - Tutorial para novos usuários
- 🔐 **JWT Authentication** - Access + Refresh tokens

## � Tecnologias

### Backend
- **Go 1.23** - Linguagem principal
- **Gin** - Framework web de alta performance
- **GORM** - ORM para PostgreSQL
- **JWT** - Autenticação com access e refresh tokens
- **PostgreSQL 16** - Banco de dados relacional

### Frontend
- **Vue.js 3** - Framework frontend com Composition API
- **Pinia** - Gerenciamento de estado
- **Vue Router** - Roteamento SPA
- **TailwindCSS** - Estilização utilitária
- **Vue I18n** - Internacionalização
- **Chart.js + Vue-ChartJS** - Gráficos interativos
- **Lucide Icons** - Ícones modernos
- **Vue Toastification** - Notificações toast
- **VueDraggable** - Drag and drop para Kanban
- **VueUse** - Composables utilitários

## 🏗️ Arquitetura

### Backend - Arquitetura em Camadas

```
backend/
├── cmd/api/                    # Entrypoint da aplicação
├── internal/
│   ├── config/                 # Configuração centralizada (env, JWT, CORS)
│   ├── database/               # Conexão PostgreSQL e migrations
│   ├── handler/                # Controllers HTTP (request/response)
│   ├── middleware/             # Auth, CORS, logging
│   ├── models/                 # Entidades GORM (User, Project, Task, etc)
│   ├── repository/             # Data access layer (queries GORM)
│   ├── router/                 # Definição de rotas Gin
│   └── service/                # Lógica de negócio (regras, validações)
├── Dockerfile
├── go.mod / go.sum
└── .env.example
```

**Fluxo de Requisição:**
1. `Router` recebe requisição HTTP
2. `Middleware` valida autenticação JWT
3. `Handler` processa request e chama `Service`
4. `Service` executa lógica de negócio
5. `Repository` acessa dados via GORM
6. Resposta retorna ao cliente

### Frontend - Vue 3 com Pinia

```
frontend/src/
├── api/                        # Cliente Axios (requisições HTTP)
├── assets/                     # CSS global, tailwind
├── components/                 # Componentes reutilizáveis
│   ├── ui/                     # Componentes base (Button, Input, etc)
│   └── charts/                 # Gráficos (Chart.js)
├── i18n/                       # Internacionalização (PT-BR, EN)
├── layouts/                    # Layouts (MainLayout, AuthLayout)
├── router/                     # Vue Router (rotas SPA)
├── stores/                     # Pinia stores (estado global)
│   ├── auth.js                 # Autenticação
│   ├── projects.js             # Projetos
│   ├── boards.js               # Boards
│   ├── tasks.js                # Tarefas
│   └── stats.js                # Estatísticas
├── views/                      # Páginas (Dashboard, Projects, Board, etc)
├── App.vue                     # Componente raiz
└── main.js                     # Entrypoint
```

**Fluxo de Dados:**
1. Componente Vue dispara ação em `Store` (Pinia)
2. Store chama `API` (Axios)
3. API faz requisição ao Backend
4. Store atualiza estado
5. Componente reage à mudança de estado

### Banco de Dados - PostgreSQL

**Entidades principais:**
- `users` - Usuários do sistema
- `projects` - Projetos
- `project_members` - Membros de projeto com roles
- `boards` - Quadros Kanban
- `lists` - Listas dentro de boards
- `tasks` - Tarefas
- `task_assignees` - Responsáveis de tarefas
- `comments` - Comentários em tarefas
- `labels` - Labels/tags
- `teams` - Equipes/Squads
- `notifications` - Notificações

## 📁 Estrutura do Projeto

```
├── backend/                    # API Go/Gin
├── frontend/                   # SPA Vue 3
├── docker-compose.yml          # Orquestração de containers
├── Makefile                    # Comandos úteis
└── README.md                   # Este arquivo
```

## 🛠️ Instalação e Execução

### Pré-requisitos
- Go 1.21+
- Node.js 18+
- PostgreSQL 15+ (ou Docker)
- Docker e Docker Compose (opcional)

### Opção 1: Com Docker (Recomendado)

```bash
# Clonar o repositório
git clone https://github.com/LeonardoRFragoso/Go-API-Gestao-de-Projetos-e-Tarefas.git
cd Go-API-Gestao-de-Projetos-e-Tarefas

# Subir todos os serviços
docker-compose up -d

# Acessar:
# - Frontend: http://localhost:3000
# - Backend API: http://localhost:8080
```

### Opção 2: Desenvolvimento Local

1. **Configurar variáveis de ambiente:**

```bash
# Backend
cp backend/.env.example backend/.env
# Edite o arquivo .env com suas configurações

# Frontend
cp frontend/.env.example frontend/.env
```

2. **Iniciar PostgreSQL (via Docker):**

```bash
docker-compose up postgres -d
```

3. **Iniciar Backend:**

```bash
cd backend
go mod tidy
go run ./cmd/api
```

4. **Iniciar Frontend:**

```bash
cd frontend
npm install
npm run dev
```

5. **Acessar:**
   - Frontend: http://localhost:5173
   - Backend API: http://localhost:8080

## 📚 API Endpoints

### Autenticação
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/auth/register` | Cadastrar usuário |
| POST | `/api/v1/auth/login` | Login |
| POST | `/api/v1/auth/refresh` | Renovar token |
| POST | `/api/v1/auth/logout` | Logout |
| GET | `/api/v1/auth/me` | Usuário atual |

### Usuários
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/users` | Listar usuários |
| GET | `/api/v1/users/:id` | Buscar usuário |
| PUT | `/api/v1/users/me` | Atualizar perfil |
| PUT | `/api/v1/users/me/password` | Alterar senha |

### Projetos
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/projects` | Listar projetos |
| POST | `/api/v1/projects` | Criar projeto |
| GET | `/api/v1/projects/:id` | Buscar projeto |
| PUT | `/api/v1/projects/:id` | Atualizar projeto |
| DELETE | `/api/v1/projects/:id` | Excluir projeto |
| GET | `/api/v1/projects/:id/members` | Listar membros |
| POST | `/api/v1/projects/:id/members` | Adicionar membro |
| DELETE | `/api/v1/projects/:id/members/:memberId` | Remover membro |

### Quadros (Boards)
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/projects/:projectId/boards` | Listar quadros |
| POST | `/api/v1/boards` | Criar quadro |
| GET | `/api/v1/boards/:id` | Buscar quadro com listas |
| PUT | `/api/v1/boards/:id` | Atualizar quadro |
| DELETE | `/api/v1/boards/:id` | Excluir quadro |

### Listas
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/lists` | Criar lista |
| GET | `/api/v1/lists/:id` | Buscar lista |
| PUT | `/api/v1/lists/:id` | Atualizar lista |
| DELETE | `/api/v1/lists/:id` | Excluir lista |
| PUT | `/api/v1/boards/:boardId/lists/reorder` | Reordenar listas |

### Tarefas
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/tasks` | Criar tarefa |
| GET | `/api/v1/tasks/:id` | Buscar tarefa |
| PUT | `/api/v1/tasks/:id` | Atualizar tarefa |
| DELETE | `/api/v1/tasks/:id` | Excluir tarefa |
| PUT | `/api/v1/tasks/:id/move` | Mover tarefa |
| POST | `/api/v1/tasks/:id/assignees` | Adicionar responsável |
| DELETE | `/api/v1/tasks/:id/assignees/:userId` | Remover responsável |
| POST | `/api/v1/tasks/:id/labels` | Adicionar label |
| DELETE | `/api/v1/tasks/:id/labels/:labelId` | Remover label |
| GET | `/api/v1/projects/:projectId/tasks/search` | Buscar tarefas |

### Comentários
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/tasks/:taskId/comments` | Listar comentários |
| POST | `/api/v1/comments` | Criar comentário |
| PUT | `/api/v1/comments/:id` | Atualizar comentário |
| DELETE | `/api/v1/comments/:id` | Excluir comentário |

### Labels
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/boards/:boardId/labels` | Listar labels |
| POST | `/api/v1/labels` | Criar label |
| PUT | `/api/v1/labels/:id` | Atualizar label |
| DELETE | `/api/v1/labels/:id` | Excluir label |

### Equipes (Teams/Squads)
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/teams` | Listar equipes |
| POST | `/api/v1/teams` | Criar equipe |
| GET | `/api/v1/teams/:id` | Buscar equipe |
| PUT | `/api/v1/teams/:id` | Atualizar equipe |
| DELETE | `/api/v1/teams/:id` | Excluir equipe |
| GET | `/api/v1/teams/:id/members` | Listar membros |
| POST | `/api/v1/teams/:id/members` | Adicionar membro |
| PUT | `/api/v1/teams/:id/members/:memberId` | Atualizar role do membro |
| DELETE | `/api/v1/teams/:id/members/:memberId` | Remover membro |
| GET | `/api/v1/teams/:id/projects` | Listar projetos da equipe |
| POST | `/api/v1/teams/:id/projects` | Vincular projeto |
| DELETE | `/api/v1/teams/:id/projects/:projectId` | Desvincular projeto |

### Notificações
| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/api/v1/notifications` | Listar notificações |
| GET | `/api/v1/notifications/count` | Contar não lidas |
| PUT | `/api/v1/notifications/:id/read` | Marcar como lida |
| PUT | `/api/v1/notifications/read-all` | Marcar todas como lidas |
| DELETE | `/api/v1/notifications/:id` | Excluir notificação |
| DELETE | `/api/v1/notifications` | Excluir todas |

## 🔐 Autenticação

A API usa JWT com dois tokens:
- **Access Token**: Válido por 15 minutos
- **Refresh Token**: Válido por 7 dias

Inclua o token no header:
```
Authorization: Bearer <access_token>
```

## 🎨 Funcionalidades do Frontend (Detalhadas)

### Autenticação & Perfil
- ✅ **Login/Registro** - Autenticação JWT com validação
- ✅ **Perfil de Usuário** - Editar nome, email, senha
- ✅ **Persistência de Sessão** - Refresh token automático

### Dashboard
- ✅ **Cards de Estatísticas** - Total de projetos, tarefas, progresso
- ✅ **Gráficos Interativos** - Tarefas por status, prioridade, evolução semanal
- ✅ **Projetos Recentes** - Acesso rápido aos últimos projetos
- ✅ **Métricas em Tempo Real** - Dados atualizados automaticamente

### Projetos
- ✅ **CRUD Completo** - Criar, editar, visualizar, excluir projetos
- ✅ **Visão Geral** - Informações consolidadas do projeto
- ✅ **Gerenciamento de Membros** - Adicionar/remover com controle de roles
- ✅ **Configurações** - Editar propriedades e gerenciar acesso

### Kanban Board
- ✅ **Drag-and-Drop** - Mover tarefas entre colunas
- ✅ **Colunas Customizáveis** - To Do, In Progress, Done
- ✅ **Filtros** - Por prioridade, responsável, label
- ✅ **Busca** - Pesquisar tarefas em tempo real

### Tarefas
- ✅ **Criação Rápida** - Modal ou inline
- ✅ **Edição Completa** - Título, descrição, prioridade, data de vencimento
- ✅ **Responsáveis** - Atribuir múltiplos usuários
- ✅ **Labels** - Categorizar com tags customizáveis
- ✅ **Comentários** - Discussão e colaboração
- ✅ **Histórico** - Rastreamento de mudanças

### Equipes
- ✅ **Gerenciamento** - Criar, editar, excluir equipes
- ✅ **Membros** - Adicionar com roles (lead, admin, member)
- ✅ **Projetos de Equipe** - Vincular projetos

### Notificações
- ✅ **Dropdown** - Acesso rápido às notificações
- ✅ **Badge de Contagem** - Número de não lidas
- ✅ **Marcar como Lida** - Individual ou em lote
- ✅ **Excluir** - Remover notificações

### Tema & Idioma
- ✅ **Dark Mode** - Toggle com persistência
- ✅ **Light Mode** - Tema claro padrão
- ✅ **Português (BR)** - Tradução completa
- ✅ **Inglês** - Suporte bilíngue

### UX/UI
- ✅ **Responsivo** - Mobile, tablet, desktop
- ✅ **Skeleton Loading** - Estados de carregamento elegantes
- ✅ **Toast Notifications** - Feedback visual
- ✅ **Animações** - Transições suaves
- ✅ **Atalhos de Teclado** - Navegação rápida
- ✅ **Acessibilidade** - ARIA labels e navegação por teclado

## ⌨️ Atalhos de Teclado

| Atalho | Ação |
|--------|------|
| `?` | Mostrar/ocultar atalhos |
| `N` | Nova tarefa |
| `/` ou `Ctrl+K` | Buscar tarefas |
| `Esc` | Fechar modal/dropdown |
| `Ctrl+D` | Toggle dark mode |
| `Ctrl+L` | Logout |

## 🖼️ Screenshots

### Dashboard
![Dashboard - Em breve](https://via.placeholder.com/800x600?text=Dashboard+Placeholder)

Visualização geral com cards de estatísticas, gráficos de tarefas por status e progresso semanal.

### Kanban Board
![Kanban Board - Em breve](https://via.placeholder.com/800x600?text=Kanban+Board+Placeholder)

Board com colunas (To Do, In Progress, Done), drag-and-drop, filtros e busca.

### Projetos
![Projetos - Em breve](https://via.placeholder.com/800x600?text=Projetos+Placeholder)

Lista de projetos com membros, boards e ações rápidas.

## 📝 Comandos do Makefile

```bash
make help           # Lista todos os comandos
make dev            # Inicia desenvolvimento completo
make dev-backend    # Inicia apenas o backend
make dev-frontend   # Inicia apenas o frontend
make build          # Build das imagens Docker
make up             # Inicia containers
make down           # Para containers
make logs           # Ver logs
make clean          # Limpa containers e volumes
make test           # Executa testes
make setup          # Configura dependências
```

## 🔧 Variáveis de Ambiente

### Backend (.env)
```env
PORT=8080
GIN_MODE=debug
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=taskmanager
DB_SSLMODE=disable
JWT_SECRET=your-super-secret-key
JWT_ACCESS_EXPIRY=15m
JWT_REFRESH_EXPIRY=168h
CORS_ORIGINS=http://localhost:5173
```

### Frontend (.env)
```env
VITE_API_URL=http://localhost:8080/api/v1
```

## 📦 Produção

Para deploy em produção:

1. Altere `JWT_SECRET` para uma chave segura
2. Configure `GIN_MODE=release`
3. Use um banco de dados PostgreSQL gerenciado
4. Configure HTTPS/SSL
5. Use o docker-compose.yml ou deploy em Kubernetes

## 🎯 Diferenciais Técnicos

### Backend
- **Arquitetura em Camadas** - Separação clara de responsabilidades (handler → service → repository)
- **GORM com Migrations** - ORM robusta com versionamento de schema
- **JWT com Refresh Tokens** - Autenticação segura com tokens de curta duração
- **Permissões Granulares** - Sistema de roles e permissões por recurso
- **Tratamento de Erros Consistente** - Response patterns padronizados
- **CORS Configurável** - Suporte a múltiplas origens

### Frontend
- **Vue 3 Composition API** - Código moderno e reativo
- **Pinia com Persistência** - Estado global com suporte a localStorage
- **Responsivo Mobile-First** - Funciona perfeitamente em todos os tamanhos
- **Dark Mode Nativo** - Tema escuro com persistência
- **Internacionalização** - Suporte a múltiplos idiomas (PT-BR, EN)
- **Gráficos Interativos** - Chart.js integrado para análise visual
- **Drag-and-Drop** - VueDraggable para Kanban intuitivo
- **Atalhos de Teclado** - Navegação rápida com atalhos customizáveis

### DevOps
- **Docker Compose** - Ambiente completo em um comando
- **Multi-stage Builds** - Imagens otimizadas e leves
- **Health Checks** - Verificação de saúde dos serviços
- **Volumes Persistentes** - Dados PostgreSQL preservados

## � Próximos Passos (Roadmap)

### Fase 1 (Atual) ✅
- [x] Arquitetura base em camadas
- [x] CRUD completo de projetos, tarefas, boards
- [x] Autenticação JWT
- [x] Dashboard com gráficos
- [x] Kanban com drag-and-drop
- [x] Equipes/Squads
- [x] Notificações
- [x] Permissões granulares
- [ ] Melhorias de UX/UI
- [ ] Tela de configurações de projeto
- [ ] Visão geral aprimorada de projeto

### Fase 2 (Planejado)
- [ ] Filtros avançados e busca full-text
- [ ] Relatórios e exportação de dados
- [ ] Integração com calendário
- [ ] Webhooks para integrações
- [ ] API GraphQL (alternativa)
- [ ] Testes automatizados (unit + integration)

### Fase 3 (Futuro)
- [ ] WebSockets para colaboração em tempo real
- [ ] Automações e workflows
- [ ] Sprints e planejamento
- [ ] Analytics avançado
- [ ] Mobile app (React Native)
- [ ] Autenticação OAuth2 (Google, GitHub)

## �📄 Licença

MIT License

## 👨‍💻 Autor

**Leonardo R. Fragoso**
- GitHub: [@LeonardoRFragoso](https://github.com/LeonardoRFragoso)
- LinkedIn: [linkedin.com/in/leonardo-fragoso](https://linkedin.com/in/leonardo-fragoso)

---

**Desenvolvido com ❤️ para equipes que buscam organização e produtividade.**
