package internal

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// OpenDB is a drop-in replacement for Open()
func OpenDB(path string) (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(path), newConfig())
}
func newLog() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second,  // Slow SQL threshold
			LogLevel:      logger.Error, // Log level
			Colorful:      true,         // Disable color
		},
	)
}

func newConfig() *gorm.Config {
	return &gorm.Config{
		SkipDefaultTransaction:                   false,
		NamingStrategy:                           nil,
		Logger:                                   newLog(),
		NowFunc:                                  func() time.Time { return time.Now().Local() },
		DryRun:                                   false,
		PrepareStmt:                              true,
		DisableAutomaticPing:                     false,
		DisableForeignKeyConstraintWhenMigrating: false,
		ClauseBuilders:                           nil,
		ConnPool:                                 nil,
		Dialector:                                nil,
		Plugins:                                  nil,
	}
}
