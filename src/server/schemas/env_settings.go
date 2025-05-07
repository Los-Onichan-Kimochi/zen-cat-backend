package schemas

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/utils/env"
)

type EnvSettings struct {
	// Logs
	EnableSqlLogs bool

	// Swagger
	EnableSwagger bool

	// Ports
	MainPort string

	// AstroCat DB
	AstroCatPostgresHost     string
	AstroCatPostgresPort     string
	AstroCatPostgresUser     string
	AstroCatPostgresPassword string
	AstroCatPostgresName     string
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

	astroCatPostgresHost := os.Getenv("ASTRO_CAT_POSTGRES_HOST")
	astroCatPostgresPort := os.Getenv("ASTRO_CAT_POSTGRES_PORT")
	astroCatPostgresUser := os.Getenv("ASTRO_CAT_POSTGRES_USER")
	astroCatPostgresPassword := os.Getenv("ASTRO_CAT_POSTGRES_PASSWORD")
	astroCatPostgresName := os.Getenv("ASTRO_CAT_POSTGRES_NAME")

	return &EnvSettings{
		EnableSqlLogs: enableSqlLogs,

		EnableSwagger: enableSwagger,

		MainPort: mainPort,

		AstroCatPostgresHost:     astroCatPostgresHost,
		AstroCatPostgresPort:     astroCatPostgresPort,
		AstroCatPostgresUser:     astroCatPostgresUser,
		AstroCatPostgresPassword: astroCatPostgresPassword,
		AstroCatPostgresName:     astroCatPostgresName,
	}
}
