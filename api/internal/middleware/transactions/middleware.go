package transactions

import (
	"net/http"

	"github.com/MAJIAXIT/api_base/api/pkg/logger"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Middleware interface {
	TransactionMiddleware() gin.HandlerFunc
	GetTx(c *gin.Context) *gorm.DB
}

type middleware struct {
	db *gorm.DB
}

func New(db *gorm.DB) Middleware {
	return &middleware{
		db: db,
	}
}

func (m *middleware) TransactionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tx := m.db.Begin()
		if tx.Error != nil {
			logger.Error("Failed to start transaction: %v", tx.Error)
			c.Status(http.StatusInternalServerError)
			c.Abort()
			return
		}

		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
				logger.Error("Panic recovered in transaction: %v", r)
				c.Status(http.StatusInternalServerError)
			}
		}()

		// Store transaction in context
		c.Set("tx", tx)

		// Process request
		c.Next()

		if len(c.Errors) > 0 || c.Writer.Status() >= 400 {
			tx.Rollback()
		} else {
			if err := tx.Commit().Error; err != nil {
				logger.Error("Failed to commit transaction: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			}
		}
	}
}

func (m *middleware) GetTx(c *gin.Context) *gorm.DB {
	if tx, exists := c.Get("tx"); exists {
		return tx.(*gorm.DB)
	}
	// Fallback to regular db
	return m.db
}
