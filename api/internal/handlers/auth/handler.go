package auth

import (
	auth_middleware "github.com/MAJIAXIT/projname/api/internal/middleware/auth"
	transactions_mw "github.com/MAJIAXIT/projname/api/internal/middleware/transactions"
	auth_service "github.com/MAJIAXIT/projname/api/internal/service/auth"
	users_service "github.com/MAJIAXIT/projname/api/internal/service/users"
	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterRoutes(router *gin.RouterGroup)
	Login(c *gin.Context)
	Signup(c *gin.Context)
	Refresh(c *gin.Context)
	Logout(c *gin.Context)
	LogoutAll(c *gin.Context)
	Me(c *gin.Context)
}

type handler struct {
	authService            auth_service.Service
	usersService           users_service.Service
	transactionsMiddleware transactions_mw.Middleware
	authMiddleware         auth_middleware.Middleware
}

func New(
	authService auth_service.Service,
	usersService users_service.Service,
	transactionsMiddleware transactions_mw.Middleware,
	authMiddleware auth_middleware.Middleware) Handler {
	return &handler{
		authService:            authService,
		usersService:           usersService,
		transactionsMiddleware: transactionsMiddleware,
		authMiddleware:         authMiddleware,
	}
}

func (h *handler) RegisterRoutes(router *gin.RouterGroup) {
	auth := router.Group("/auth")
	{
		afterAuth := auth.Group("")
		afterAuth.Use(h.authMiddleware.AuthMiddleware())
		{
			afterAuth.POST("/logout", h.Logout)
			afterAuth.GET("/me", h.Me)
		}

		authTx := auth.Group("")
		authTx.Use(h.transactionsMiddleware.TransactionMiddleware())
		{
			authTx.POST("/login", h.Login)
			authTx.POST("/signup", h.Signup)
			authTx.POST("/refresh", h.Refresh)
		}
	}
}
