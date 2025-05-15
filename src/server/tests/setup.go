package tests

import (
	"context"
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
	logger logging.Logger,
	astroCatPsqlDB *gorm.DB,
	envSetting *schemas.EnvSettings,
	t *testing.T,
) {
	if envSetting.AstroCatPostgresHost != "localhost" {
		msg := "Not allow clear Levels Postgres DB into instance different to localhost"
		if t == nil {
			panic(msg)
		} else {
			t.Fatalf("%s. This function should only be used for tests in local environment", msg)
		}

		return
	}

	if astroCatPsqlDB != nil {
		astroCatPsqlDB.Where("true").Delete(&model.Community{})

		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = &CustomLogger{}
		}

		if t == nil {
			// Add sample data for testing
		}
	}
}
