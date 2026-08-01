package database

import (
	"fmt"
	"sync"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	GetDB() *gorm.DB
	WithTransaction(fn func(Database) error) error
}

type PostgresDatabase struct {
	DB *gorm.DB
}

var (
	dbOnce     sync.Once
	dbInstance *PostgresDatabase
)

func NewPostgresDatabase(dbConfig *bootstrap.DataBase, dbConst *bootstrap.DBConst) *PostgresDatabase {
	dbOnce.Do(func() {
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
			dbConfig.Host,
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Port,
			dbConfig.SSLMode,
		)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			panic(fmt.Errorf("failed to connect database: %w", err))
		}

		sqlDB, err := db.DB()
		if err != nil {
			panic(fmt.Errorf("failed to get sql.DB from gorm: %w", err))
		}

		sqlDB.SetMaxOpenConns(int(dbConst.MaxOpenDbConn))
		sqlDB.SetConnMaxIdleTime(dbConst.MaxIdleDbConn)
		sqlDB.SetConnMaxLifetime(dbConst.MaxDbLifeTime)

		dbInstance = &PostgresDatabase{DB: db}
	})

	return dbInstance
}

func (pgx *PostgresDatabase) GetDB() *gorm.DB {
	return dbInstance.DB
}

func (pgx *PostgresDatabase) WithTransaction(fn func(Database) error) error {
	return pgx.DB.Transaction(func(tx *gorm.DB) error {
		txWrapper := &PostgresDatabase{DB: tx}
		return fn(txWrapper)
	})
}
