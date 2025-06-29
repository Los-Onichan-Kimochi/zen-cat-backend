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

	// Auth
	DisableAuthForTests bool

	// Ports
	MainPort string

	// AstroCat DB
	AstroCatPostgresHost     string
	AstroCatPostgresPort     string
	AstroCatPostgresUser     string
	AstroCatPostgresPassword string
	AstroCatPostgresName     string
	AstroCatPsqlSslMode      string

	// JWT
	TokenSignatureKey []byte

	// Email
	EmailHost     string
	EmailPort     int
	EmailUser     string
	EmailPassword string
	EmailFrom     string

	// AWS S3
	AwsAccessKeyId     string
	AwsSecretAccessKey string
	AwsSessionToken    string
	AwsRegion          string
	S3BucketName       string

	// Twilio
	TwilioAccountSid  string
	TwilioAuthToken   string
	TwilioPhoneNumber string
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
		enableSqlLogs = false
	}

	enableSwagger, err := strconv.ParseBool(os.Getenv("ENABLE_SWAGGER"))
	if err != nil {
		enableSwagger = false
	}

	mainPort := os.Getenv("MAIN_PORT")
	// Railway uses PORT environment variable
	if mainPort == "" {
		mainPort = os.Getenv("PORT")
	}
	// Default port if none is specified
	if mainPort == "" {
		mainPort = "8080"
	}

	astroCatPostgresHost := os.Getenv("ASTRO_CAT_POSTGRES_HOST")
	astroCatPostgresPort := os.Getenv("ASTRO_CAT_POSTGRES_PORT")
	astroCatPostgresUser := os.Getenv("ASTRO_CAT_POSTGRES_USER")
	astroCatPostgresPassword := os.Getenv("ASTRO_CAT_POSTGRES_PASSWORD")
	astroCatPostgresName := os.Getenv("ASTRO_CAT_POSTGRES_NAME")
	astroCatPsqlSslMode := os.Getenv("ASTRO_CAT_PSQL_SSL_MODE")

	tokenSignatureKey := []byte(os.Getenv("TOKEN_SIGNATURE_KEY"))

	// Email
	emailHost := os.Getenv("EMAIL_HOST")
	emailPortStr := os.Getenv("EMAIL_PORT")
	emailUser := os.Getenv("EMAIL_USER")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailFrom := os.Getenv("EMAIL_FROM")

	emailPort, err := strconv.Atoi(emailPortStr)
	if err != nil {
		emailPort = 0
	}

	// AWS S3
	awsAccessKeyId := os.Getenv("AWS_ACCESS_KEY_ID")
	awsSecretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	awsSessionToken := os.Getenv("AWS_SESSION_TOKEN")
	awsRegion := os.Getenv("AWS_REGION")
	s3BucketName := os.Getenv("S3_BUCKET_NAME")

	// Twilio
	twilioAccountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	twilioAuthToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")

	return &EnvSettings{
		EnableSqlLogs: enableSqlLogs,

		EnableSwagger: enableSwagger,

		MainPort: mainPort,

		AstroCatPostgresHost:     astroCatPostgresHost,
		AstroCatPostgresPort:     astroCatPostgresPort,
		AstroCatPostgresUser:     astroCatPostgresUser,
		AstroCatPostgresPassword: astroCatPostgresPassword,
		AstroCatPostgresName:     astroCatPostgresName,
		AstroCatPsqlSslMode:      astroCatPsqlSslMode,

		TokenSignatureKey: tokenSignatureKey,

		EmailHost:     emailHost,
		EmailPort:     emailPort,
		EmailUser:     emailUser,
		EmailPassword: emailPassword,
		EmailFrom:     emailFrom,

		AwsAccessKeyId:     awsAccessKeyId,
		AwsSecretAccessKey: awsSecretAccessKey,
		AwsSessionToken:    awsSessionToken,
		AwsRegion:          awsRegion,
		S3BucketName:       s3BucketName,

		TwilioAccountSid:  twilioAccountSid,
		TwilioAuthToken:   twilioAuthToken,
		TwilioPhoneNumber: twilioPhoneNumber,
	}
}
