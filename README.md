# 📋 TaskManager - Sistema de Gestão de Projetos e Tarefas

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.23-00ADD8?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/Vue.js-3-4FC08D?style=for-the-badge&logo=vue.js&logoColor=white" alt="Vue.js" />
  <img src="https://img.shields.io/badge/PostgreSQL-16-336791?style=for-the-badge&logo=postgresql&logoColor=white" alt="PostgreSQL" />
  <img src="https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white" alt="Docker" />
</p>

Sistema completo estilo Trello para gestão de projetos e tarefas, com backend em Go e frontend em Vue.js 3. Inclui **Dark Mode**, **Internacionalização (PT-BR/EN)**, **Gráficos**, **Kanban Board** com drag-and-drop, **Equipes/Squads**, **Notificações** e muito mais!

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

## 📁 Estrutura do Projeto

```
├── backend/
│   ├── cmd/api/           # Entrypoint da aplicação
│   ├── internal/
│   │   ├── config/        # Configurações
│   │   ├── database/      # Conexão e migrations
│   │   ├── handler/       # Controllers/Handlers
│   │   ├── middleware/    # Auth, CORS
│   │   ├── models/        # Entidades do banco
│   │   ├── repository/    # Acesso a dados
│   │   ├── router/        # Definição de rotas
│   │   └── service/       # Lógica de negócio
│   ├── Dockerfile
│   ├── go.mod
│   └── .env.example
├── frontend/
│   ├── src/
│   │   ├── api/           # Cliente Axios
│   │   ├── assets/        # CSS global
│   │   ├── layouts/       # Layouts da aplicação
│   │   ├── router/        # Vue Router
│   │   ├── stores/        # Pinia stores
│   │   └── views/         # Páginas
│   ├── Dockerfile
│   ├── package.json
│   └── .env.example
├── docker-compose.yml
├── Makefile
└── README.md
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

## 🎨 Funcionalidades do Frontend

- ✅ **Autenticação** - Login/registro com JWT
- ✅ **Dashboard** - Métricas, gráficos de pizza e barras
- ✅ **Dark Mode** - Toggle de tema com persistência
- ✅ **Internacionalização** - PT-BR e Inglês
- ✅ **CRUD de Projetos** - Criar, editar, excluir projetos
- ✅ **CRUD de Quadros** - Gerenciar boards do Kanban
- ✅ **Kanban Board** - Drag-and-drop com colunas coloridas
- ✅ **Busca e Filtros** - Pesquisar tarefas, filtrar por prioridade
- ✅ **Prioridades** - Baixa, média, alta, urgente com badges coloridos
- ✅ **Toast Notifications** - Feedback visual para todas as ações
- ✅ **Skeleton Loading** - Estados de carregamento elegantes
- ✅ **Responsivo** - Menu mobile, bottom navigation
- ✅ **Atalhos de Teclado** - `?` para ver atalhos, `N` para nova tarefa
- ✅ **Onboarding** - Tutorial interativo para novos usuários
- ✅ **Animações** - Transições suaves e micro-interações
- ✅ **Equipes/Squads** - Gerenciamento de equipes com membros e projetos
- ✅ **Notificações** - Dropdown com notificações em tempo real
- ✅ **Permissões** - Sistema de roles (lead, admin, member)

## ⌨️ Atalhos de Teclado

| Atalho | Ação |
|--------|------|
| `?` | Mostrar/ocultar atalhos |
| `N` | Nova tarefa |
| `/` ou `Ctrl+K` | Buscar |
| `Esc` | Fechar modal |

## 🖼️ Screenshots

### Dashboard (Light Mode)
O dashboard exibe estatísticas de tarefas com gráficos interativos e lista de projetos recentes.

### Kanban Board (Dark Mode)
Board Kanban com colunas coloridas (To Do, In Progress, Done), drag-and-drop e filtros.

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

## 📄 Licença

MIT License

## 👨‍💻 Autor

**Leonardo R. Fragoso**
- GitHub: [@LeonardoRFragoso](https://github.com/LeonardoRFragoso)
