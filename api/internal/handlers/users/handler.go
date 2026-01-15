package users

import (
	auth_mw "github.com/MAJIAXIT/api_base/api/internal/middleware/auth"
	transactions_mw "github.com/MAJIAXIT/api_base/api/internal/middleware/transactions"
	"github.com/MAJIAXIT/api_base/api/internal/service/users"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(router *gin.RouterGroup)
	UpdateUser(c *gin.Context)
}

type handler struct {
	usersService           users.Service
	transactionsMiddleware transactions_mw.Middleware
	authMiddleware         auth_mw.Middleware
}

func New(
	usersService users.Service,
	transactionsMiddleware transactions_mw.Middleware,
	authMiddleware auth_mw.Middleware) Handler {
	return &handler{
		usersService:           usersService,
		transactionsMiddleware: transactionsMiddleware,
		authMiddleware:         authMiddleware,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	// Authenticated user routes
	users := router.Group("/users")
	users.Use(h.authMiddleware.AuthMiddleware())
	{
		current := users.Group("/current")
		{
			current.GET("", h.GetCurrentUser)
			current.POST("/update", h.UpdateUser)
		}
	}
}
