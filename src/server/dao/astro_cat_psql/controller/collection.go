package controller

import (
	// para agregar data dummy
	"time"

	"github.com/google/uuid"
	//-----------------------
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils/psql"
)

type AstroCatPsqlCollection struct {
	Logger           logging.Logger
	Community        *Community
	Professional     *Professional
	Local            *Local
	User             *User
	Service          *Service
	Plan             *Plan
	CommunityPlan    *CommunityPlan
	CommunityService *CommunityService
}

// Create dao controller collection
func NewAstroCatPsqlCollection(
	logger logging.Logger,
	envSettings *schemas.EnvSettings,
) (*AstroCatPsqlCollection, *gorm.DB) {
	postgresqlDB, err := psql.CreateConnection(
		envSettings.AstroCatPostgresHost,
		envSettings.AstroCatPostgresUser,
		envSettings.AstroCatPostgresPassword,
		envSettings.AstroCatPostgresName,
		envSettings.AstroCatPostgresPort,
		envSettings.EnableSqlLogs,
	)
	if err != nil {
		logger.Panicln("Failed to connect to AstroCat Postgresql database")
	}

	if err := postgresqlDB.Use(otelgorm.NewPlugin()); err != nil {
		logger.Panicln("Failed to instrument AstroCat Postgresql database")
	}

	createTables(postgresqlDB)

	return &AstroCatPsqlCollection{
		Logger:           logger,
		Community:        NewCommunityController(logger, postgresqlDB),
		Professional:     NewProfessionalController(logger, postgresqlDB),
		Local:            NewLocalController(logger, postgresqlDB),
		User:             NewUserController(logger, postgresqlDB),
		Service:          NewServiceController(logger, postgresqlDB),
		Plan:             NewPlanController(logger, postgresqlDB),
		CommunityPlan:    NewCommunityPlanController(logger, postgresqlDB),
		CommunityService: NewCommunityServiceController(logger, postgresqlDB),
	}, postgresqlDB
}

// Helper function to create AstroCat tables
func createTables(astroCatPsqlDB *gorm.DB) {
	/*	// Drop existing tables - lo he usado para dropear
		// provisionalmente y meter data dummy - a borrar mas adelante
		astroCatPsqlDB.Migrator().DropTable(
			&model.Plan{},
			&model.Template{},
			&model.Local{},
			&model.Professional{},
			&model.Onboarding{},
			&model.User{},
			&model.Community{},
			&model.Membership{},
			&model.Service{},
			&model.Session{},
			&model.Reservation{},
		)
	*/
	if err := astroCatPsqlDB.AutoMigrate(&model.Plan{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Template{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Local{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Professional{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Onboarding{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.User{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Community{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Membership{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Service{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Session{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.Reservation{}); err != nil {
		panic(err)
	}

	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityService{}); err != nil {
		panic(err)
	}
	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityPlan{}); err != nil {
		panic(err)
	}

	// Add dummy data only if no users exist - a borrar mas adelante
	var count int64
	astroCatPsqlDB.Model(&model.User{}).Count(&count)
	if count == 0 {
		// Create dummy plan
		plan := &model.Plan{
			Id:               uuid.New(),
			Fee:              0.0,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: nil,
			AuditFields: model.AuditFields{
				UpdatedBy: "system",
			},
		}
		if err := astroCatPsqlDB.Create(plan).Error; err != nil {
			panic(err)
		}

		// Create dummy community
		community := &model.Community{
			Id:                  uuid.New(),
			Name:                "Dummy Community",
			Purpose:             "Community for testing",
			ImageUrl:            "",
			NumberSubscriptions: 0,
			AuditFields: model.AuditFields{
				UpdatedBy: "system",
			},
		}
		if err := astroCatPsqlDB.Create(community).Error; err != nil {
			panic(err)
		}

		// Create dummy user
		user := &model.User{
			Id:             uuid.New(),
			Name:           "Test",
			FirstLastName:  "User",
			SecondLastName: nil,
			Password:       "test123",
			Email:          "test@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "",
			AuditFields: model.AuditFields{
				UpdatedBy: "system",
			},
		}
		if err := astroCatPsqlDB.Create(user).Error; err != nil {
			panic(err)
		}

		// Create dummy membership
		membership := &model.Membership{
			Id:          uuid.New(),
			Description: "Test membership",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(1, 0, 0), // 1 year from now
			Status:      model.MembershipStatusActive,
			AuditFields: model.AuditFields{
				UpdatedBy: "system",
			},
			CommunityId: community.Id,
			UserId:      user.Id,
			PlanId:      plan.Id,
		}
		if err := astroCatPsqlDB.Create(membership).Error; err != nil {
			panic(err)
		}
	}
}
