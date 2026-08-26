package database

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/mehdi1443/cinema-booking-system/bootstrap"
	domainlogger "github.com/mehdi1443/cinema-booking-system/internal/domain/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
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

func NewPostgresDatabase(dbConfig *bootstrap.DataBase, dbConst *bootstrap.DBConst, appLogger domainlogger.Logger) *PostgresDatabase {
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

		customLogger := NewGormLoggerAdapter(appLogger, 200*time.Millisecond)

		gormConfig := &gorm.Config{
			Logger: customLogger.LogMode(gormlogger.Warn),
		}

		db, err := gorm.Open(postgres.Open(dsn), gormConfig)
		if err != nil {
			log.Fatalf("Critical error: failed to connect to Postgres DB '%s': %v", dbConfig.Name, err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Critical error: failed to get sql.DB from gorm: %v", err)
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
