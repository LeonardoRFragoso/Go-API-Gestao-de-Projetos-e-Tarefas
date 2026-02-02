package router

import (
	"github.com/gin-gonic/gin"
	"github.com/taskmanager/backend/internal/handler"
	"github.com/taskmanager/backend/internal/middleware"
	"github.com/taskmanager/backend/internal/service"
)

type Router struct {
	engine         *gin.Engine
	authHandler    *handler.AuthHandler
	userHandler    *handler.UserHandler
	projectHandler *handler.ProjectHandler
	boardHandler   *handler.BoardHandler
	listHandler    *handler.ListHandler
	taskHandler    *handler.TaskHandler
	commentHandler *handler.CommentHandler
	labelHandler   *handler.LabelHandler
	statsHandler   *handler.StatsHandler
	authService    service.AuthService
	corsOrigins    string
}

func New(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	projectHandler *handler.ProjectHandler,
	boardHandler *handler.BoardHandler,
	listHandler *handler.ListHandler,
	taskHandler *handler.TaskHandler,
	commentHandler *handler.CommentHandler,
	labelHandler *handler.LabelHandler,
	statsHandler *handler.StatsHandler,
	authService service.AuthService,
	corsOrigins string,
) *Router {
	return &Router{
		engine:         gin.Default(),
		authHandler:    authHandler,
		userHandler:    userHandler,
		projectHandler: projectHandler,
		boardHandler:   boardHandler,
		listHandler:    listHandler,
		taskHandler:    taskHandler,
		commentHandler: commentHandler,
		labelHandler:   labelHandler,
		statsHandler:   statsHandler,
		authService:    authService,
		corsOrigins:    corsOrigins,
	}
}

func (r *Router) Setup() *gin.Engine {
	r.engine.Use(middleware.CORSMiddleware(r.corsOrigins))

	api := r.engine.Group("/api/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", r.authHandler.Register)
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(r.authService))
		{
			protected.POST("/auth/logout", r.authHandler.Logout)
			protected.GET("/auth/me", r.authHandler.Me)

			users := protected.Group("/users")
			{
				users.GET("", r.userHandler.List)
				users.GET("/:id", r.userHandler.GetByID)
				users.PUT("/me", r.userHandler.Update)
				users.PUT("/me/password", r.userHandler.UpdatePassword)
			}

			projects := protected.Group("/projects")
			{
				projects.POST("", r.projectHandler.Create)
				projects.GET("", r.projectHandler.List)
				projects.GET("/:id", r.projectHandler.GetByID)
				projects.PUT("/:id", r.projectHandler.Update)
				projects.DELETE("/:id", r.projectHandler.Delete)
				projects.GET("/:id/members", r.projectHandler.GetMembers)
				projects.POST("/:id/members", r.projectHandler.AddMember)
				projects.DELETE("/:id/members/:memberId", r.projectHandler.RemoveMember)
				projects.GET("/:id/boards", r.boardHandler.ListByProject)
				projects.GET("/:id/tasks/search", r.taskHandler.Search)
			}

			boards := protected.Group("/boards")
			{
				boards.POST("", r.boardHandler.Create)
				boards.GET("/:id", r.boardHandler.GetByID)
				boards.PUT("/:id", r.boardHandler.Update)
				boards.DELETE("/:id", r.boardHandler.Delete)
				boards.PUT("/:id/lists/reorder", r.listHandler.Reorder)
				boards.GET("/:id/labels", r.labelHandler.ListByBoard)
			}

			lists := protected.Group("/lists")
			{
				lists.POST("", r.listHandler.Create)
				lists.GET("/:id", r.listHandler.GetByID)
				lists.PUT("/:id", r.listHandler.Update)
				lists.DELETE("/:id", r.listHandler.Delete)
				lists.PUT("/:id/tasks/reorder", r.taskHandler.Reorder)
			}

			tasks := protected.Group("/tasks")
			{
				tasks.POST("", r.taskHandler.Create)
				tasks.GET("/:id", r.taskHandler.GetByID)
				tasks.PUT("/:id", r.taskHandler.Update)
				tasks.DELETE("/:id", r.taskHandler.Delete)
				tasks.PUT("/:id/move", r.taskHandler.Move)
				tasks.POST("/:id/assignees", r.taskHandler.AddAssignee)
				tasks.DELETE("/:id/assignees/:userId", r.taskHandler.RemoveAssignee)
				tasks.POST("/:id/labels", r.taskHandler.AddLabel)
				tasks.DELETE("/:id/labels/:labelId", r.taskHandler.RemoveLabel)
				tasks.GET("/:id/comments", r.commentHandler.ListByTask)
			}

			comments := protected.Group("/comments")
			{
				comments.POST("", r.commentHandler.Create)
				comments.GET("/:id", r.commentHandler.GetByID)
				comments.PUT("/:id", r.commentHandler.Update)
				comments.DELETE("/:id", r.commentHandler.Delete)
			}

			labels := protected.Group("/labels")
			{
				labels.POST("", r.labelHandler.Create)
				labels.GET("/:id", r.labelHandler.GetByID)
				labels.PUT("/:id", r.labelHandler.Update)
				labels.DELETE("/:id", r.labelHandler.Delete)
			}

			stats := protected.Group("/stats")
			{
				stats.GET("/dashboard", r.statsHandler.GetDashboardStats)
				stats.GET("/tasks", r.statsHandler.GetTaskStats)
				stats.GET("/weekly", r.statsHandler.GetWeeklyProgress)
			}
		}
	}

	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	return r.engine
}
