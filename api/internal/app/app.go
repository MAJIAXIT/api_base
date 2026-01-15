package app

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MAJIAXIT/api_base/api/config"
	auth_mw "github.com/MAJIAXIT/api_base/api/internal/middleware/auth"
	transactions_mw "github.com/MAJIAXIT/api_base/api/internal/middleware/transactions"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	auth_hnd "github.com/MAJIAXIT/api_base/api/internal/handlers/auth"
	users_hnd "github.com/MAJIAXIT/api_base/api/internal/handlers/users"

	auth_svc "github.com/MAJIAXIT/api_base/api/internal/service/auth"
	users_svc "github.com/MAJIAXIT/api_base/api/internal/service/users"

	"github.com/MAJIAXIT/api_base/api/pkg/database"
	"github.com/MAJIAXIT/api_base/api/pkg/logger"
)

type App struct {
	config      *config.Config
	db          *gorm.DB
	router      *gin.Engine
	httpServer  *http.Server
	httpsServer *http.Server

	// Background job context
	backgroundCtx    context.Context
	backgroundCancel context.CancelFunc

	// Middleware
	transactionsMiddleware transactions_mw.Middleware
	authMiddleware         auth_mw.Middleware

	// Services
	authService  auth_svc.Service
	usersService users_svc.Service

	// Handlers
	authHandler  auth_hnd.Handler
	usersHandler users_hnd.Handler
}

func New() *App {
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgres(&cfg.Database)
	if err != nil {
		logger.Fatal("Failed to connect to database: %v", err)
	}

	// Create background context
	backgroundCtx, backgroundCancel := context.WithCancel(context.Background())

	app := &App{
		config:           cfg,
		db:               db,
		backgroundCtx:    backgroundCtx,
		backgroundCancel: backgroundCancel,
	}

	app.setupServices()
	app.setupMiddlewares()
	app.setupHandlers()
	app.setupServers()
	app.setupRoutes()

	return app
}

func (a *App) setupServices() {
	a.usersService = users_svc.New()
	a.authService = auth_svc.New(
		&a.config.JWT,
		a.usersService)
}

func (a *App) setupMiddlewares() {
	a.transactionsMiddleware = transactions_mw.New(a.db)
	a.authMiddleware = auth_mw.New(a.authService)
}

func (a *App) setupHandlers() {
	a.authHandler = auth_hnd.New(
		a.authService,
		a.usersService,
		a.transactionsMiddleware,
		a.authMiddleware,
	)
	a.usersHandler = users_hnd.New(
		a.usersService,
		a.transactionsMiddleware,
		a.authMiddleware,
	)

}

func (a *App) setupServers() {
	a.router = gin.New()

	a.router.Use(gin.Recovery())
	a.router.Use(gin.Logger())

	a.httpServer = &http.Server{
		Addr:         ":" + a.config.Server.HTTPPort,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}

	a.httpsServer = &http.Server{
		Addr:         ":" + a.config.Server.HTTPSPort,
		Handler:      a.router,
		ReadTimeout:  a.config.Server.ReadTimeout,
		WriteTimeout: a.config.Server.WriteTimeout,
	}
}

func (a *App) setupRoutes() {
	v1 := a.router.Group("/api/v1")

	a.authHandler.RegisterRoutes(v1)
	a.usersHandler.RegisterRoutes(v1)
}

// func (a *App) startBackgroundJobs() {
// 	// Start notification checker goroutine
// 	go func() {
// 		logger.Info("Starting subscription notification checker (every 1 hour)")
// 		ticker := time.NewTicker(1 * time.Hour)
// 		defer ticker.Stop()

// 		// Run initial check after 1 minute (to allow app to fully start)
// 		initialTimer := time.NewTimer(1 * time.Minute)
// 		defer initialTimer.Stop()

// 		for {
// 			select {
// 			case <-a.backgroundCtx.Done():
// 				logger.Info("Stopping subscription notification checker")
// 				return
// 			case <-initialTimer.C:
// 				// Run initial check
// 				if err := a.notificationService.CheckAndSendNotifications(a.db); err != nil {
// 					logger.Error("Error in initial notification check: %v", err)
// 				}
// 				// Disable initial timer after first run
// 				initialTimer.Stop()
// 			case <-ticker.C:
// 				// Run periodic check
// 				if err := a.notificationService.CheckAndSendNotifications(a.db); err != nil {
// 					logger.Error("Error in periodic notification check: %v", err)
// 				}
// 			}
// 		}
// 	}()
// }

func (a *App) Start() {
	logger.Info("Starting servers...")

	// Start background jobs
	// a.startBackgroundJobs()

	// Start HTTP server
	go func() {
		logger.Info("Starting HTTP server on :%s", a.config.Server.HTTPPort)
		if err := a.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("HTTP server error: %v", err)
		}
	}()

	// Start HTTPS server
	// go func() {
	// 	logger.Info("Starting HTTPS server on :%s", a.config.Server.HTTPSPort)
	// 	if err := a.httpsServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
	// 		logger.Fatal("HTTPS server error: %v", err)
	// 	}
	// }()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall. SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down servers...")

	// Cancel background jobs
	a.backgroundCancel()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		logger.Info("HTTP server Shutdown: %v", err)
	}
	// if err := a.httpsServer.Shutdown(ctx); err != nil {
	// 	logger.Info("HTTPS server Shutdown: %v", err)
	// }

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	logger.Info("Server exited gracefully")
	logger.Shutdown()
}
