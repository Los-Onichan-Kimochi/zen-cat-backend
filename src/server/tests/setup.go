package tests

import (
	"context"
	"fmt"
	"testing"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type CustomLogger struct{}

func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	// Ignore info logs
}

func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	// Ignore warning logs
}

func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	// Ignore error logs
}

func (l *CustomLogger) Trace(
	ctx context.Context,
	begin time.Time,
	fc func() (string, int64),
	err error,
) {
	// Ignore trace logs
}

// Remove all data from AstroCatPsql db.
//   - Note: Only use for tests
func ClearPostgresqlDatabase(
	appLogger logging.Logger,
	astroCatPsqlDB *gorm.DB,
	envSetting *schemas.EnvSettings,
	t *testing.T,
) {
	if envSetting.AstroCatPostgresHost != "localhost" {
		msg := "Not allow clear Levels Postgres DB into instance different to localhost"
		if t == nil {
			appLogger.Panicf(
				"%s. This function should only be used for tests in local environment",
				msg,
			)
		} else {
			t.Fatalf("%s. This function should only be used for tests in local environment", msg)
		}
		return
	}

	if astroCatPsqlDB != nil {
		fmt.Println("...Clearing AstroCatPsql database (hard delete)...")

		originalLogger := astroCatPsqlDB.Logger
		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = originalLogger.LogMode(logger.Info)
		}

		// Iniciar una transacción
		tx := astroCatPsqlDB.Begin()

		// Desactivar temporalmente las restricciones de clave foránea
		tx.Exec("SET CONSTRAINTS ALL DEFERRED")

		// Primero eliminar las membresías que tienen referencias a otras tablas
		if err := tx.Unscoped().Where("true").Delete(&model.Membership{}).Error; err != nil {
			tx.Rollback()
			appLogger.Errorf("Error clearing Membership table: %v", err)
			return
		}

		// Luego eliminar las otras tablas
		tablesToClear := map[string]interface{}{
			"User":         &model.User{},
			"Community":    &model.Community{},
			"Professional": &model.Professional{},
			"Service":      &model.Service{},
			"Plan":         &model.Plan{},
		}

		for tableName, modelInstance := range tablesToClear {
			appLogger.Infof("Attempting to hard delete all records from %s table...", tableName)
			if err := tx.Unscoped().Where("true").Delete(modelInstance).Error; err != nil {
				tx.Rollback()
				appLogger.Errorf("Error clearing %s table: %v", tableName, err)
				return
			}
		}

		// Reactivar las restricciones de clave foránea
		tx.Exec("SET CONSTRAINTS ALL IMMEDIATE")

		// Confirmar la transacción
		if err := tx.Commit().Error; err != nil {
			appLogger.Errorf("Error committing transaction: %v", err)
			return
		}

		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = originalLogger
		}

		if t == nil {
			// Add sample data for testing
		}
	} else {
		appLogger.Warn("astroCatPsqlDB is nil, skipping database clearing.")
	}
}
