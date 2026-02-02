package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/taskmanager/backend/internal/config"
	"github.com/taskmanager/backend/internal/database"
	"github.com/taskmanager/backend/internal/handler"
	"github.com/taskmanager/backend/internal/repository"
	"github.com/taskmanager/backend/internal/router"
	"github.com/taskmanager/backend/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	gin.SetMode(cfg.Server.GinMode)

	db, err := database.Connect(&cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := database.Migrate(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	boardRepo := repository.NewBoardRepository(db)
	listRepo := repository.NewListRepository(db)
	taskRepo := repository.NewTaskRepository(db)
	commentRepo := repository.NewCommentRepository(db)
	labelRepo := repository.NewLabelRepository(db)

	authService := service.NewAuthService(userRepo, tokenRepo, &cfg.JWT)
	userService := service.NewUserService(userRepo)
	projectService := service.NewProjectService(projectRepo)
	boardService := service.NewBoardService(boardRepo, listRepo)
	listService := service.NewListService(listRepo)
	taskService := service.NewTaskService(taskRepo)
	commentService := service.NewCommentService(commentRepo)
	labelService := service.NewLabelService(labelRepo)

	authHandler := handler.NewAuthHandler(authService, userService)
	userHandler := handler.NewUserHandler(userService)
	projectHandler := handler.NewProjectHandler(projectService)
	boardHandler := handler.NewBoardHandler(boardService, projectService)
	listHandler := handler.NewListHandler(listService, boardService, projectService)
	taskHandler := handler.NewTaskHandler(taskService, listService, boardService, projectService)
	commentHandler := handler.NewCommentHandler(commentService)
	labelHandler := handler.NewLabelHandler(labelService)
	statsHandler := handler.NewStatsHandler(taskRepo, projectService)

	r := router.New(
		authHandler,
		userHandler,
		projectHandler,
		boardHandler,
		listHandler,
		taskHandler,
		commentHandler,
		labelHandler,
		statsHandler,
		authService,
		cfg.CORS.Origins,
	)

	engine := r.Setup()

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := engine.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
