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
	"onichankimochi.com/astro_cat_backend/src/server/utils"
)

type CustomLogger struct{}

// Helper function to create string pointers
func strPtr(s string) *string {
	return &s
}

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
			{"ServiceProfessional", &model.ServiceProfessional{}},
			{"ServiceLocal", &model.ServiceLocal{}},
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
	// contrasenha test123
	// lo hasheo
	hashedPassword, err := utils.HashPassword("test123")
	if err != nil {
		appLogger.Errorf("Error hashing password: %v", err)
		return
	}

	// Create dummy users
	users := []*model.User{
		{
			Id:             uuid.New(),
			Name:           "Test-1",
			FirstLastName:  "User",
			SecondLastName: nil,
			Password:       hashedPassword,
			Email:          "test-1@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "TestAdmin",
			FirstLastName:  "User",
			SecondLastName: nil,
			Password:       hashedPassword,
			Email:          "testAdmin@zen-cat.com",
			Rol:            model.UserRolAdmin,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "María",
			FirstLastName:  "González",
			SecondLastName: strPtr("López"),
			Password:       "maria123",
			Email:          "maria.gonzalez@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Carlos",
			FirstLastName:  "Mendoza",
			SecondLastName: strPtr("Ruiz"),
			Password:       "carlos123",
			Email:          "carlos.mendoza@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Ana",
			FirstLastName:  "Martínez",
			SecondLastName: nil,
			Password:       "ana123",
			Email:          "ana.martinez@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Luis",
			FirstLastName:  "Rodríguez",
			SecondLastName: strPtr("Flores"),
			Password:       "luis123",
			Email:          "luis.rodriguez@zen-cat.com",
			Rol:            model.UserRolClient,
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
		{
			Id:             uuid.New(),
			Name:           "Pedro",
			FirstLastName:  "Sánchez",
			SecondLastName: strPtr("García"),
			Specialty:      "Entrenamiento Personal",
			Email:          "pedro.sanchez@gym.com",
			PhoneNumber:    "555-0123",
			Type:           model.ProfessionalTypeGymTrainer,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Laura",
			FirstLastName:  "Fernández",
			SecondLastName: nil,
			Specialty:      "Yoga Avanzado",
			Email:          "laura.fernandez@yoga.com",
			PhoneNumber:    "555-0456",
			Type:           model.ProfessionalTypeYogaTrainer,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			Name:           "Roberto",
			FirstLastName:  "Díaz",
			SecondLastName: strPtr("Morales"),
			Specialty:      "Medicina General",
			Email:          "roberto.diaz@medic.com",
			PhoneNumber:    "555-0789",
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
		{
			Id:             uuid.New(),
			LocalName:      "Studio Zen",
			StreetName:     "Wellness Blvd",
			BuildingNumber: "789",
			District:       "Wellness",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Top floor wellness center",
			Capacity:       12,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Fitness Plus",
			StreetName:     "Sports Ave",
			BuildingNumber: "321",
			District:       "Sports",
			Province:       "Central",
			Region:         "Metropolitan",
			Reference:      "Next to Sports Complex",
			Capacity:       25,
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
	// Obtener la zona horaria de Lima (UTC-5)
	limaLocation, _ := time.LoadLocation("America/Lima")
	now := time.Now().In(limaLocation)
	baseDate := now.Truncate(24 * time.Hour) // Start of today in Lima timezone

	sessions := []*model.Session{
		// Sesión que ya terminó (ayer)
		{
			Id:              uuid.New(),
			Title:           "Morning Yoga",
			Date:            baseDate.Add(-24 * time.Hour),             // Ayer
			StartTime:       baseDate.Add(-24*time.Hour + 8*time.Hour), // Ayer 8:00 AM
			EndTime:         baseDate.Add(-24*time.Hour + 9*time.Hour), // Ayer 9:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 5, // Test-1, María, Carlos(anulled), Ana, Luis
			Capacity:        20,
			SessionLink:     nil,
			ProfessionalId:  professionals[0].Id, // John - Yoga Trainer
			LocalId:         &locals[1].Id,       // Local Yoga
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Sesión de hoy por la tarde (futura si es mañana, pasada si es noche)
		{
			Id:              uuid.New(),
			Title:           "Evening Gym",
			Date:            baseDate,
			StartTime:       baseDate.Add(18 * time.Hour), // Hoy 6:00 PM
			EndTime:         baseDate.Add(19 * time.Hour), // Hoy 7:00 PM
			State:           model.SessionStateScheduled,
			RegisteredCount: 3, // Admin, María, Carlos(cancelled)
			Capacity:        15,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[0].Id,       // Local Gym
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Sesión que está en curso (hora actual +/- 30 min)
		{
			Id:              uuid.New(),
			Title:           "Advanced Yoga Workshop",
			Date:            baseDate,
			StartTime:       now.Add(-30 * time.Minute), // Empezó hace 30 minutos
			EndTime:         now.Add(30 * time.Minute),  // Termina en 30 minutos
			State:           model.SessionStateScheduled,
			RegisteredCount: 4, // Test-1, María, Ana, Luis(done)
			Capacity:        12,
			SessionLink:     nil,
			ProfessionalId:  professionals[3].Id, // Laura - Advanced Yoga Trainer
			LocalId:         &locals[2].Id,       // Studio Zen
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Sesión futura para mañana
		{
			Id:              uuid.New(),
			Title:           "Personal Training Session",
			Date:            baseDate.Add(24 * time.Hour),              // Tomorrow
			StartTime:       baseDate.Add(24*time.Hour + 16*time.Hour), // 4:00 PM tomorrow
			EndTime:         baseDate.Add(24*time.Hour + 17*time.Hour), // 5:00 PM tomorrow
			State:           model.SessionStateScheduled,
			RegisteredCount: 1, // Carlos
			Capacity:        4,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[3].Id,       // Fitness Plus
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Medical Consultation",
			Date:            baseDate.Add(48 * time.Hour),              // Day after tomorrow
			StartTime:       baseDate.Add(48*time.Hour + 14*time.Hour), // 2:00 PM
			EndTime:         baseDate.Add(48*time.Hour + 15*time.Hour), // 3:00 PM
			State:           model.SessionStateScheduled,
			RegisteredCount: 2, // Ana, Luis
			Capacity:        5,
			SessionLink:     strPtr("https://meet.example.com/medical-session"),
			ProfessionalId:  professionals[1].Id, // Jane - Medic
			LocalId:         nil,                 // Virtual session
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Weekend Yoga Flow",
			Date:            baseDate.Add(72 * time.Hour),              // 3 days from now
			StartTime:       baseDate.Add(72*time.Hour + 9*time.Hour),  // 9:00 AM
			EndTime:         baseDate.Add(72*time.Hour + 10*time.Hour), // 10:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 5, // Test-1, Admin, María, Carlos, Ana
			Capacity:        15,
			SessionLink:     nil,
			ProfessionalId:  professionals[0].Id, // John - Yoga Trainer
			LocalId:         &locals[1].Id,       // Local Yoga
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Strength Training Bootcamp",
			Date:            baseDate.Add(96 * time.Hour),             // 4 days from now
			StartTime:       baseDate.Add(96*time.Hour + 7*time.Hour), // 7:00 AM
			EndTime:         baseDate.Add(96*time.Hour + 8*time.Hour), // 8:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 5, // Test-1, María, Carlos, Ana, Luis(done)
			Capacity:        25,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[3].Id,       // Fitness Plus
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "General Health Checkup",
			Date:            baseDate.Add(120 * time.Hour),              // 5 days from now
			StartTime:       baseDate.Add(120*time.Hour + 11*time.Hour), // 11:00 AM
			EndTime:         baseDate.Add(120*time.Hour + 12*time.Hour), // 12:00 PM
			State:           model.SessionStateScheduled,
			RegisteredCount: 1, // Admin
			Capacity:        3,
			SessionLink:     strPtr("https://meet.example.com/health-checkup"),
			ProfessionalId:  professionals[4].Id, // Roberto - General Medicine
			LocalId:         nil,                 // Virtual session
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
		// Reservations for Morning Yoga (sessions[0])
		{
			Id:               uuid.New(),
			Name:             "Morning Yoga - Test User",
			ReservationTime:  baseDate.Add(8 * time.Hour), // Same as session start time
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[0].Id, // Test-1
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Morning Yoga - María",
			ReservationTime:  baseDate.Add(8 * time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[2].Id, // María
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Morning Yoga - Carlos",
			ReservationTime:  baseDate.Add(8 * time.Hour),
			State:            model.ReservationStateAnulled,
			LastModification: time.Now(),
			UserId:           users[3].Id, // Carlos
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Morning Yoga - Ana",
			ReservationTime:  baseDate.Add(8 * time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[4].Id, // Ana
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Morning Yoga - Luis",
			ReservationTime:  baseDate.Add(8 * time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[5].Id, // Luis
			SessionId:        sessions[0].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Evening Gym (sessions[1])
		{
			Id:               uuid.New(),
			Name:             "Evening Gym - Admin User",
			ReservationTime:  baseDate.Add(18 * time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[1].Id, // Test-2 (Admin)
			SessionId:        sessions[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Evening Gym - María",
			ReservationTime:  baseDate.Add(18 * time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[2].Id, // María
			SessionId:        sessions[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Evening Gym - Carlos",
			ReservationTime:  baseDate.Add(18 * time.Hour),
			State:            model.ReservationStateCancelled,
			LastModification: time.Now(),
			UserId:           users[3].Id, // Carlos
			SessionId:        sessions[1].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Advanced Yoga Workshop (sessions[2])
		{
			Id:               uuid.New(),
			Name:             "Advanced Yoga Workshop - Test User",
			ReservationTime:  baseDate.Add(24*time.Hour + 10*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[0].Id, // Test-1
			SessionId:        sessions[2].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Advanced Yoga Workshop - María",
			ReservationTime:  baseDate.Add(24*time.Hour + 10*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[2].Id, // María
			SessionId:        sessions[2].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Advanced Yoga Workshop - Ana",
			ReservationTime:  baseDate.Add(24*time.Hour + 10*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[4].Id, // Ana
			SessionId:        sessions[2].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Advanced Yoga Workshop - Luis",
			ReservationTime:  baseDate.Add(24*time.Hour + 10*time.Hour),
			State:            model.ReservationStateDone,
			LastModification: time.Now(),
			UserId:           users[5].Id, // Luis
			SessionId:        sessions[2].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Personal Training Session (sessions[3])
		{
			Id:               uuid.New(),
			Name:             "Personal Training - Carlos",
			ReservationTime:  baseDate.Add(24*time.Hour + 16*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[3].Id, // Carlos
			SessionId:        sessions[3].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Medical Consultation (sessions[4])
		{
			Id:               uuid.New(),
			Name:             "Medical Consultation - Ana",
			ReservationTime:  baseDate.Add(48*time.Hour + 14*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[4].Id, // Ana
			SessionId:        sessions[4].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Medical Consultation - Luis",
			ReservationTime:  baseDate.Add(48*time.Hour + 14*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[5].Id, // Luis
			SessionId:        sessions[4].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Weekend Yoga Flow (sessions[5])
		{
			Id:               uuid.New(),
			Name:             "Weekend Yoga Flow - Test User",
			ReservationTime:  baseDate.Add(72*time.Hour + 9*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[0].Id, // Test-1
			SessionId:        sessions[5].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Weekend Yoga Flow - Admin",
			ReservationTime:  baseDate.Add(72*time.Hour + 9*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[1].Id, // Test-2 (Admin)
			SessionId:        sessions[5].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Weekend Yoga Flow - María",
			ReservationTime:  baseDate.Add(72*time.Hour + 9*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[2].Id, // María
			SessionId:        sessions[5].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Weekend Yoga Flow - Carlos",
			ReservationTime:  baseDate.Add(72*time.Hour + 9*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[3].Id, // Carlos
			SessionId:        sessions[5].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Weekend Yoga Flow - Ana",
			ReservationTime:  baseDate.Add(72*time.Hour + 9*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[4].Id, // Ana
			SessionId:        sessions[5].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for Strength Training Bootcamp (sessions[6])
		{
			Id:               uuid.New(),
			Name:             "Strength Training Bootcamp - Test User",
			ReservationTime:  baseDate.Add(96*time.Hour + 7*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[0].Id, // Test-1
			SessionId:        sessions[6].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Strength Training Bootcamp - María",
			ReservationTime:  baseDate.Add(96*time.Hour + 7*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[2].Id, // María
			SessionId:        sessions[6].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Strength Training Bootcamp - Carlos",
			ReservationTime:  baseDate.Add(96*time.Hour + 7*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[3].Id, // Carlos
			SessionId:        sessions[6].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Strength Training Bootcamp - Ana",
			ReservationTime:  baseDate.Add(96*time.Hour + 7*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[4].Id, // Ana
			SessionId:        sessions[6].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Name:             "Strength Training Bootcamp - Luis",
			ReservationTime:  baseDate.Add(96*time.Hour + 7*time.Hour),
			State:            model.ReservationStateDone,
			LastModification: time.Now(),
			UserId:           users[5].Id, // Luis
			SessionId:        sessions[6].Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},

		// Reservations for General Health Checkup (sessions[7])
		{
			Id:               uuid.New(),
			Name:             "General Health Checkup - Admin",
			ReservationTime:  baseDate.Add(120*time.Hour + 11*time.Hour),
			State:            model.ReservationStateConfirmed,
			LastModification: time.Now(),
			UserId:           users[1].Id, // Test-2 (Admin)
			SessionId:        sessions[7].Id,
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
			District:       "Downtown",
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
			District:       "Business",
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
