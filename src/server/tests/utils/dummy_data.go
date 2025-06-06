package utils

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
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

		// Start a transaction
		tx := astroCatPsqlDB.Begin()

		// Disable foreign key constraints temporarily
		tx.Exec("SET CONSTRAINTS ALL DEFERRED")

		// First delete tables that have references to other tables
		tablesToClear := []struct {
			name  string
			model any
		}{
			// First delete tables with foreign key dependencies
			{"Membership", &model.Membership{}},
			{"CommunityPlan", &model.CommunityPlan{}},
			{"CommunityService", &model.CommunityService{}},
			{"Reservation", &model.Reservation{}},
			{"Session", &model.Session{}},
			{"Onboarding", &model.Onboarding{}},
			{"Template", &model.Template{}},

			// Then delete independent tables
			{"Professional", &model.Professional{}},
			{"Local", &model.Local{}},
			{"User", &model.User{}},
			{"Plan", &model.Plan{}},
			{"Service", &model.Service{}},
			{"Community", &model.Community{}},
		}

		for _, table := range tablesToClear {
			appLogger.Infof("Attempting to hard delete all records from %s table...", table.name)
			if err := tx.Unscoped().Where("true").Delete(table.model).Error; err != nil {
				tx.Rollback()
				appLogger.Errorf("Error clearing %s table: %v", table.name, err)
				return
			}
		}

		// Reactivate foreign key constraints
		tx.Exec("SET CONSTRAINTS ALL IMMEDIATE")

		// Confirm the transaction
		if err := tx.Commit().Error; err != nil {
			appLogger.Errorf("Error committing transaction: %v", err)
			return
		}

		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = originalLogger
		}

		if t == nil {
			createDummyData(appLogger, astroCatPsqlDB)
		}
	} else {
		appLogger.Warn("astroCatPsqlDB is nil, skipping database clearing.")
	}
}

// Creates sample data for testing purposes
func createDummyData(appLogger logging.Logger, astroCatPsqlDB *gorm.DB) {
	fmt.Println("Creating dummy data...")

	// Create dummy plans
	reservationLimit := 8
	plans := []*model.Plan{
		{
			Id:               uuid.Must(uuid.Parse("d1694efe-9a13-42d7-a9e8-4d629f9f2f35")),
			Fee:              70.0,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.Must(uuid.Parse("6d222f80-8887-4cc2-b6a1-48d08cd2d742")),
			Fee:              1000.0,
			Type:             model.PlanTypeAnual,
			ReservationLimit: nil,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.Must(uuid.Parse("eb71f5e0-589d-4f1b-86e7-696c30e92bfe")),
			Fee:              69.0,
			Type:             model.PlanTypeAnual,
			ReservationLimit: nil,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, plan := range plans {
		if err := astroCatPsqlDB.Create(plan).Error; err != nil {
			appLogger.Errorf("Error creating dummy plan: %v", err)
			return
		}
	}

	// Create dummy communities
	communities := []*model.Community{
		{
			Id:                  uuid.Must(uuid.Parse("e804b95a-a388-4751-b246-96fe97232d35")),
			Name:                "Yoga Community",
			Purpose:             "Community for yoga enthusiasts",
			ImageUrl:            "test-image",
			NumberSubscriptions: 0,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:                  uuid.Must(uuid.Parse("a1570014-f96c-4ba1-9ac6-e2aec2127910")),
			Name:                "Gym Group",
			Purpose:             "Community for meditation practitioners",
			ImageUrl:            "test-image",
			NumberSubscriptions: 0,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:                  uuid.Must(uuid.Parse("76035ca7-1d3b-4d7d-9091-fc55f7410e59")),
			Name:                "Gamers Group",
			Purpose:             "Community for Dota2",
			ImageUrl:            "test-image",
			NumberSubscriptions: 0,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, community := range communities {
		if err := astroCatPsqlDB.Create(community).Error; err != nil {
			appLogger.Errorf("Error creating dummy community: %v", err)
			return
		}
	}

	// Create dummy users
	users := []*model.User{
		{
			Id:             uuid.New(),
			Name:           "Test-1",
			FirstLastName:  "User",
			SecondLastName: nil,
			Password:       "test123",
			Email:          "test-1@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Test-2",
			FirstLastName:  "User",
			SecondLastName: nil,
			Password:       "test123",
			Email:          "test-2@zen-cat.com",
			Rol:            model.UserRolAdmin,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, user := range users {
		if err := astroCatPsqlDB.Create(user).Error; err != nil {
			appLogger.Errorf("Error creating dummy user: %v", err)
			return
		}
	}

	// Create dummy services
	services := []*model.Service{
		{
			Id:          uuid.New(),
			Name:        "Yoga",
			Description: "Servicio de yoga",
			ImageUrl:    "test-image",
			IsVirtual:   false,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "GYM",
			Description: "Servicio de gimnasio",
			ImageUrl:    "test-image",
			IsVirtual:   false,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "Citas Médicas",
			Description: "Servicio online de citas médicas",
			ImageUrl:    "test-image",
			IsVirtual:   true,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, service := range services {
		if err := astroCatPsqlDB.Create(service).Error; err != nil {
			appLogger.Errorf("Error creating dummy service: %v", err)
			return
		}
	}

	// Create dummy professionals
	professionals := []*model.Professional{
		{
			Id:             uuid.New(),
			Name:           "John",
			FirstLastName:  "Doe",
			SecondLastName: nil,
			Specialty:      "Yoga",
			Email:          "john@gmail.com",
			PhoneNumber:    "123456789",
			Type:           model.ProfessionalTypeYogaTrainer,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Jane",
			FirstLastName:  "Smith",
			SecondLastName: nil,
			Specialty:      "Cardiología",
			Email:          "jane@gmail.com",
			PhoneNumber:    "987654321",
			Type:           model.ProfessionalTypeMedic,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, professional := range professionals {
		if err := astroCatPsqlDB.Create(professional).Error; err != nil {
			appLogger.Errorf("Error creating dummy professional: %v", err)
			return
		}
	}

	// Create dummy locals
	locals := []*model.Local{
		{
			Id:             uuid.New(),
			LocalName:      "Local Gym",
			StreetName:     "Main St",
			BuildingNumber: "123",
			District:       "Downtown",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Near Central Park",
			Capacity:       20,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Local Yoga",
			StreetName:     "Downtown Ave",
			BuildingNumber: "456",
			District:       "Business",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Near Business Center",
			Capacity:       15,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, local := range locals {
		if err := astroCatPsqlDB.Create(local).Error; err != nil {
			appLogger.Errorf("Error creating dummy local: %v", err)
			return
		}
	}

	// Create dummy memberships
	memberships := []*model.Membership{
		{
			Id:          uuid.New(),
			Description: "Monthly Yoga Membership",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(0, 1, 0),
			Status:      model.MembershipStatusActive,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
			CommunityId: communities[0].Id,
			UserId:      users[0].Id,
			PlanId:      plans[0].Id,
		},
		{
			Id:          uuid.New(),
			Description: "Yearly Gym Membership",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(1, 0, 0),
			Status:      model.MembershipStatusActive,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
			CommunityId: communities[1].Id,
			UserId:      users[1].Id,
			PlanId:      plans[1].Id,
		},
	}
	for _, membership := range memberships {
		if err := astroCatPsqlDB.Create(membership).Error; err != nil {
			appLogger.Errorf("Error creating dummy membership: %v", err)
			return
		}
	}

	// Create dummy community services
	communityServices := []*model.CommunityService{
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id,
			ServiceId:   services[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[1].Id,
			ServiceId:   services[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, cs := range communityServices {
		if err := astroCatPsqlDB.Create(cs).Error; err != nil {
			appLogger.Errorf("Error creating dummy community service: %v", err)
			return
		}
	}

	// Create dummy community plans
	communityPlans := []*model.CommunityPlan{
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id,
			PlanId:      plans[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[1].Id,
			PlanId:      plans[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[2].Id,
			PlanId:      plans[2].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, cp := range communityPlans {
		if err := astroCatPsqlDB.Create(cp).Error; err != nil {
			appLogger.Errorf("Error creating dummy community plan: %v", err)
			return
		}
	}

	// Create dummy sessions
	sessions := []*model.Session{
		{
			Id:              uuid.New(),
			Title:           "Morning Yoga",
			Date:            time.Now(),
			StartTime:       time.Now(),
			EndTime:         time.Now().Add(time.Hour),
			State:           model.SessionStateOnGoing,
			RegisteredCount: 5,
			Capacity:        20,
			SessionLink:     nil,
			ProfessionalId:  professionals[0].Id,
			LocalId:         &locals[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Evening Gym",
			Date:            time.Now().Add(time.Hour * 2),
			StartTime:       time.Now().Add(time.Hour * 2),
			EndTime:         time.Now().Add(time.Hour * 3),
			State:           model.SessionStateScheduled,
			RegisteredCount: 3,
			Capacity:        15,
			SessionLink:     nil,
			ProfessionalId:  professionals[1].Id,
			LocalId:         &locals[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, session := range sessions {
		if err := astroCatPsqlDB.Create(session).Error; err != nil {
			appLogger.Errorf("Error creating dummy session: %v", err)
			return
		}
	}

	// Create dummy reservations
	reservations := []*model.Reservation{
		{
			Id:               uuid.New(),
			Name:             "Yoga Class Reservation",
			ReservationTime:  time.Now().Add(time.Hour * 24),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[0].Id,
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Gym Session Reservation",
			ReservationTime:  time.Now().Add(time.Hour * 48),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[1].Id,
			SessionId:        sessions[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, reservation := range reservations {
		if err := astroCatPsqlDB.Create(reservation).Error; err != nil {
			appLogger.Errorf("Error creating dummy reservation: %v", err)
			return
		}
	}

	// Create dummy templates
	templates := []*model.Template{
		{
			Id:             uuid.New(),
			Link:           "https://example.com/medic-template",
			ProfessionalId: professionals[1].Id, // Medic
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, template := range templates {
		if err := astroCatPsqlDB.Create(template).Error; err != nil {
			appLogger.Errorf("Error creating dummy template: %v", err)
			return
		}
	}

	// Create dummy onboarding
	onboardings := []*model.Onboarding{
		{
			Id:             uuid.New(),
			PhoneNumber:    "123456789",
			DocumentType:   model.DocumentTypeDni,
			DocumentNumber: "12345678",
			StreetName:     "Main St",
			BuildingNumber: "123",
			District:       "Downtown",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Near Central Park",
			UserId:         users[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			PhoneNumber:    "987654321",
			DocumentType:   model.DocumentTypeForeignerCard,
			DocumentNumber: "87654321",
			StreetName:     "Downtown Ave",
			BuildingNumber: "456",
			District:       "Business",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Near Business Center",
			UserId:         users[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, onboarding := range onboardings {
		if err := astroCatPsqlDB.Create(onboarding).Error; err != nil {
			appLogger.Errorf("Error creating dummy onboarding: %v", err)
			return
		}
	}
}
