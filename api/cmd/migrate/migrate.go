package migrate

import (
	"os/user"
	"runtime/debug"

	"github.com/MAJIAXIT/projname/api/config"
	"github.com/MAJIAXIT/projname/api/internal/models/session"
	"github.com/MAJIAXIT/projname/api/pkg/database"
	"github.com/MAJIAXIT/projname/api/pkg/logger"
	"gorm.io/gorm"
)

func main() {
	// var (
	// 	example      = flag.String("ex", "default", "description")
	// )
	// flag.Parse()

	dbConfig := config.LoadDBConfig()
	db, err := database.NewPostgres(dbConfig)
	if err != nil {
		logger.Fatal("Failed to start migration: %v", err)
	}

	tx := db.Begin()
	if tx.Error != nil {
		logger.Fatal("Failed to start transaction: %v", tx.Error)
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			stack := debug.Stack()
			logger.Fatal("Panic recovered: %v\nStack trace:\n%s", r, string(stack))
		}
	}()

	if err := dropTables(tx); err != nil {
		tx.Rollback()
		logger.Fatal("Failed to drop tables: %v", err)
	}
	if err := createTables(tx); err != nil {
		tx.Rollback()
		logger.Fatal("Failed to create tables: %v", err)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		logger.Fatal("Failed to commit transaction: %v, transaction rollbacked", err)
	}

	logger.Info("Migrated")
}

func dropTables(tx *gorm.DB) error {
	if err := tx.Migrator().DropTable(
	// &user.User{},
	// &session.Session{},
	); err != nil {
		return logger.WrapError(err)
	}

	return nil
}

func createTables(tx *gorm.DB) error {
	if err := tx.Migrator().AutoMigrate(
		&user.User{},
		&session.Session{},
	); err != nil {
		return logger.WrapError(err)
	}

	return nil
}
