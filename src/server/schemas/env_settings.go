package schemas

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"onichankimochi.com/zen_cat_backend/src/logging"
	"onichankimochi.com/zen_cat_backend/src/server/utils/env"
)

type EnvSettings struct {
	// Logs
	EnableSqlLogs bool

	// Swagger
	EnableSwagger bool

	// Ports
	MainPort string

	// ZenCat DB
	ZenCatPostgresHost     string
	ZenCatPostgresPort     string
	ZenCatPostgresUser     string
	ZenCatPostgresPassword string
	ZenCatPostgresName     string
}

// Create a new env settings defined on .env file
func NewEnvSettings(logger logging.Logger) *EnvSettings {
	// STAGE is an env var to be use in arquitecture
	if stage := os.Getenv("STAGE"); stage == "local" || stage == "" {
		if envPath, err := env.FindEnvPath(); err != nil {
			logger.Panicln(".env", err)
		} else if err := godotenv.Load(envPath); err != nil {
			logger.Panicln("Failed to load .env file", err)
		}
	}

	enableSqlLogs, err := strconv.ParseBool(os.Getenv("ENABLE_SQL_LOGS"))
	if err != nil {
		logger.Panicln("Invalid value for ENABLE_SQL_LOGS, must be boolean", err)
	}

	enableSwagger, err := strconv.ParseBool(os.Getenv("ENABLE_SWAGGER"))
	if err != nil {
		logger.Panicln("Invalid value for ENABLE_SWAGGER, must be boolean", err)
	}

	mainPort := os.Getenv("MAIN_PORT")

	zenCatPostgresHost := os.Getenv("ZEN_CAT_POSTGRES_HOST")
	zenCatPostgresPort := os.Getenv("ZEN_CAT_POSTGRES_PORT")
	zenCatPostgresUser := os.Getenv("ZEN_CAT_POSTGRES_USER")
	zenCatPostgresPassword := os.Getenv("ZEN_CAT_POSTGRES_PASSWORD")
	zenCatPostgresName := os.Getenv("ZEN_CAT_POSTGRES_NAME")

	return &EnvSettings{
		EnableSqlLogs: enableSqlLogs,

		EnableSwagger: enableSwagger,

		MainPort: mainPort,

		ZenCatPostgresHost:     zenCatPostgresHost,
		ZenCatPostgresPort:     zenCatPostgresPort,
		ZenCatPostgresUser:     zenCatPostgresUser,
		ZenCatPostgresPassword: zenCatPostgresPassword,
		ZenCatPostgresName:     zenCatPostgresName,
	}
}
