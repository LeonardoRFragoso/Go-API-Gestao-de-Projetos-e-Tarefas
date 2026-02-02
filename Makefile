.PHONY: help dev dev-backend dev-frontend build up down logs clean test migrate

help:
	@echo "Comandos disponíveis:"
	@echo "  make dev           - Inicia desenvolvimento (backend + frontend)"
	@echo "  make dev-backend   - Inicia apenas o backend"
	@echo "  make dev-frontend  - Inicia apenas o frontend"
	@echo "  make build         - Build das imagens Docker"
	@echo "  make up            - Inicia todos os containers"
	@echo "  make down          - Para todos os containers"
	@echo "  make logs          - Exibe logs dos containers"
	@echo "  make clean         - Remove containers e volumes"
	@echo "  make test          - Executa testes"
	@echo "  make migrate       - Executa migrations"

dev:
	@echo "Iniciando desenvolvimento..."
	docker-compose up postgres -d
	@echo "Aguardando PostgreSQL..."
	@sleep 3
	$(MAKE) dev-backend &
	$(MAKE) dev-frontend

dev-backend:
	cd backend && go run ./cmd/api

dev-frontend:
	cd frontend && npm run dev

build:
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	docker system prune -f

test:
	cd backend && go test ./...

migrate:
	cd backend && go run ./cmd/api migrate

# Backend commands
backend-deps:
	cd backend && go mod tidy

backend-build:
	cd backend && go build -o bin/api ./cmd/api

# Frontend commands
frontend-deps:
	cd frontend && npm install

frontend-build:
	cd frontend && npm run build

# Full setup
setup: frontend-deps backend-deps
	@echo "Dependências instaladas!"
	cp backend/.env.example backend/.env
	cp frontend/.env.example frontend/.env
	@echo "Arquivos .env criados. Configure-os antes de iniciar."
