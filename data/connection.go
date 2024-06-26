package data

import (
	"fmt"
	"os"

	"github.com/dogefuzz/dogefuzz/config"
	"github.com/dogefuzz/dogefuzz/entities"
	"go.uber.org/zap"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type connection struct {
	db           *gorm.DB
	databaseName string
	logger       *zap.Logger
}

func NewConnection(cfg *config.Config, logger *zap.Logger) (*connection, error) {
	logger.Info(fmt.Sprintf("Initializing database in \"%s.db\" file", cfg.DatabaseName))
	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("%s.db", cfg.DatabaseName)), &gorm.Config{
		Logger: DataLogger{
			Logger: logger,
		},
	})
	if err != nil {
		return nil, err
	}

	db.Exec("PRAGMA journal_mode=WAL;")

	db.Exec("pragma page_size = 32768;")
	db.Exec("pragma mmap_size = 30000000000;")
	db.Exec("pragma temp_store = memory;")
	db.Exec("pragma synchronous = normal;")

	return &connection{
		db:           db,
		databaseName: cfg.DatabaseName,
		logger:       logger,
	}, nil
}

func (m *connection) Clean() error {
	m.logger.Info(fmt.Sprintf("Cleaning \"%s.db\" database file", m.databaseName))
	return os.Remove(fmt.Sprintf("%s.db", m.databaseName))
}

func (m *connection) GetDB() *gorm.DB {
	return m.db
}

func (m *connection) Migrate() error {
	m.logger.Info("migrating database...")
	return m.db.AutoMigrate(&entities.Contract{}, &entities.Function{}, &entities.Task{}, &entities.Transaction{})
}
