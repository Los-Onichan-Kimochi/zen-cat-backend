package psql

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Create postgresql connection
func CreateConnection(
	host string,
	user string,
	password string,
	dbanme string,
	port string,
	enableLogs bool,
) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbanme, port,
	)
	gormConfig := &gorm.Config{}
	if enableLogs {
		gormConfig.Logger = logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: false,
			Colorful:                  true,
		})
	} else {
		// Explicitly disable all logging when enableLogs is false
		gormConfig.Logger = logger.Default.LogMode(logger.Silent)
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// Configure connection pool to prevent exhaustion during tests
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(10)
	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(5)
	// Set maximum lifetime of connections
	sqlDB.SetConnMaxLifetime(time.Hour)

	return db, nil
}
