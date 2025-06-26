package tests

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

func ClearPostgresqlDatabaseTesting(
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
		// fmt.Println("...Clearing AstroCatPsql database (hard delete)...")

		originalLogger := astroCatPsqlDB.Logger
		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = originalLogger.LogMode(logger.Silent)
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
			{"AuditLog", &model.AuditLog{}}, // Clear audit logs first to avoid FK constraints
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

	} else {
		appLogger.Warn("astroCatPsqlDB is nil, skipping database clearing.")
	}
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
		// fmt.Println("...Clearing AstroCatPsql database (hard delete)...")

		originalLogger := astroCatPsqlDB.Logger
		if !envSetting.EnableSqlLogs {
			astroCatPsqlDB.Logger = originalLogger.LogMode(logger.Silent)
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
			{"AuditLog", &model.AuditLog{}}, // Clear audit logs first to avoid FK constraints
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
	reservationLimitBasic := 5
	reservationLimitPremium := 15
	plans := []*model.Plan{
		{
			Id:               uuid.New(),
			Fee:              70.0,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimit,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Fee:              1000.0,
			Type:             model.PlanTypeAnual,
			ReservationLimit: nil,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Planes para runners
		{
			Id:               uuid.New(),
			Fee:              49.90,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimitBasic,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Fee:              89.90,
			Type:             model.PlanTypeMonthly,
			ReservationLimit: &reservationLimitPremium,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Fee:              499.00,
			Type:             model.PlanTypeAnual,
			ReservationLimit: nil,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:               uuid.New(),
			Fee:              899.00,
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
	// Define a fixed UUID for the main test community that will be used in frontend
	mainCommunityId, _ := uuid.Parse("ade8c5e1-ab82-47e0-b48b-3f8f2324c450")

	communities := []*model.Community{
		{
			Id:                  mainCommunityId, // Fixed UUID for frontend integration
			Name:                "Runners",
			Purpose:             "Comunidad principal de bienestar que ofrece servicios de yoga, atención médica y fitness para mejorar tu calidad de vida",
			ImageUrl:            "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			NumberSubscriptions: 150,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:                  uuid.New(),
			Name:                "Maternal Care Community",
			Purpose:             "Comunidad especializada en cuidado maternal y servicios de lactario",
			ImageUrl:            "https://images.unsplash.com/photo-1555252333-9f8e92e65df9?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			NumberSubscriptions: 75,
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
	// Define fixed UUIDs for system and main test users
	systemUserId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000") // System user for anonymous events
	mainUserId, _ := uuid.Parse("11111111-1111-1111-1111-111111111111")

	users := []*model.User{
		// SYSTEM user for anonymous/error events
		{
			Id:             systemUserId,
			Name:           "SYSTEM",
			FirstLastName:  "ANONYMOUS",
			SecondLastName: nil,
			Password:       hashedPassword, // Same hash but this user can't actually login
			Email:          "system@zen-cat.internal",
			Rol:            model.UserRolAdmin, // Admin role for system operations
			ImageUrl:       "system-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             mainUserId, // Fixed UUID for frontend integration
			Name:           "Usuario",
			FirstLastName:  "Demo",
			SecondLastName: nil,
			Password:       hashedPassword,
			Email:          "demo@zen-cat.com",
			Rol:            model.UserRolClient,
			ImageUrl:       "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
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
			Password:       hashedPassword,
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
			Password:       hashedPassword,
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
			Password:       hashedPassword,
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
			Password:       hashedPassword,
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
			Description: "Clases para practicar posturas, respiración y meditación. Mejora tu flexibilidad, reduce el estrés y encuentra equilibrio interior.",
			ImageUrl:    "https://images.unsplash.com/photo-1544367567-0f2fcb009e0b?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			IsVirtual:   false,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "Cita médica",
			Description: "Atención personalizada con profesionales de la salud para consultas, diagnósticos y tratamientos. Agenda tu cita fácilmente.",
			ImageUrl:    "https://images.unsplash.com/photo-1559757148-5c350d0d3c56?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			IsVirtual:   true,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "Lactario",
			Description: "Área privada y cómoda para que las mamás puedan amamantar o extraer leche materna en un entorno seguro y tranquilo.",
			ImageUrl:    "https://images.unsplash.com/photo-1555252333-9f8e92e65df9?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			IsVirtual:   false,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "Gimnasio",
			Description: "Espacio equipado con máquinas y pesas para entrenamiento físico. Mejora tu condición física con rutinas personalizadas.",
			ImageUrl:    "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			IsVirtual:   false,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			Name:        "Entrenamiento Funcional",
			Description: "Ejercicios que imitan movimientos de la vida diaria para mejorar fuerza, coordinación y resistencia de manera integral.",
			ImageUrl:    "https://images.unsplash.com/photo-1571019613454-1cb2f99b2d8b?ixlib=rb-4.0.3&ixid=M3wxMjA3fDB8MHxwaG90by1wYWdlfHx8fGVufDB8fHx8fA%3D%3D&auto=format&fit=crop&w=1000&q=80",
			IsVirtual:   false,
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
			LocalName:      "Pabellón A",
			StreetName:     "Av. Almirante Cornejo",
			BuildingNumber: "1504",
			District:       "San Miguel",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Cerca al parque central",
			Capacity:       20,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Pabellón A",
			StreetName:     "Av. Constructores",
			BuildingNumber: "12345",
			District:       "San Juan de Lurigancho",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Frente al centro comercial",
			Capacity:       15,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Pabellón A",
			StreetName:     "Prolong. Santa María del Carmen",
			BuildingNumber: "1234",
			District:       "San Miguel",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Al lado del hospital",
			Capacity:       12,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Pabellón B",
			StreetName:     "Av. Universitaria",
			BuildingNumber: "100",
			District:       "San Miguel",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Campus universitario",
			Capacity:       25,
			ImageUrl:       "test-image",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			LocalName:      "Studio Zen",
			StreetName:     "Av. La Marina",
			BuildingNumber: "2000",
			District:       "San Miguel",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Centro comercial Plaza San Miguel",
			Capacity:       18,
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
			Description: "Yearly Maternal Membership",
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
		{
			Id:          uuid.New(),
			Description: "Yearly Wellness Membership",
			StartDate:   time.Now(),
			EndDate:     time.Now().AddDate(1, 0, 0),
			Status:      model.MembershipStatusActive,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
			CommunityId: communities[0].Id,
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
		// Main ZenCat Wellness Community services
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // ZenCat Wellness Community
			ServiceId:   services[0].Id,    // Yoga
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // ZenCat Wellness Community
			ServiceId:   services[1].Id,    // Cita médica
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // ZenCat Wellness Community
			ServiceId:   services[3].Id,    // Gimnasio
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Maternal Care Community services
		{
			Id:          uuid.New(),
			CommunityId: communities[1].Id, // Maternal Care Community
			ServiceId:   services[2].Id,    // Lactario
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[1].Id, // Maternal Care Community
			ServiceId:   services[4].Id,    // Entrenamiento Funcional
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
		// Runners Community plans (4 total - 2 monthly tiers + 2 annual tiers)
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[2].Id,       // Monthly Basic Plan ($49.90)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[3].Id,       // Monthly Premium Plan ($89.90)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[4].Id,       // Annual Basic Plan ($499.00)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[5].Id,       // Annual Premium Plan ($899.00)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Keep original plans for backward compatibility
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[0].Id,       // Original Monthly Plan ($70.0)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:          uuid.New(),
			CommunityId: communities[0].Id, // Runners Community
			PlanId:      plans[1].Id,       // Original Annual Plan ($1000.0)
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Maternal Care Community plan
		{
			Id:          uuid.New(),
			CommunityId: communities[1].Id, // Maternal Care Community
			PlanId:      plans[0].Id,       // Monthly Plan
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

	// Create dummy service-professional relationships
	serviceProfessionals := []*model.ServiceProfessional{
		// Yoga service with yoga trainers
		{
			Id:             uuid.New(),
			ServiceId:      services[0].Id,      // Yoga
			ProfessionalId: professionals[0].Id, // John - Yoga Trainer
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			ServiceId:      services[0].Id,      // Yoga
			ProfessionalId: professionals[3].Id, // Laura - Advanced Yoga Trainer
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Medical service with doctors
		{
			Id:             uuid.New(),
			ServiceId:      services[1].Id,      // Cita médica
			ProfessionalId: professionals[1].Id, // Jane - Medic
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:             uuid.New(),
			ServiceId:      services[1].Id,      // Cita médica
			ProfessionalId: professionals[4].Id, // Roberto - General Medicine
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Gym service with gym trainers
		{
			Id:             uuid.New(),
			ServiceId:      services[3].Id,      // Gimnasio
			ProfessionalId: professionals[2].Id, // Pedro - Gym Trainer
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Functional training with gym trainer
		{
			Id:             uuid.New(),
			ServiceId:      services[4].Id,      // Entrenamiento Funcional
			ProfessionalId: professionals[2].Id, // Pedro - Gym Trainer
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, sp := range serviceProfessionals {
		if err := astroCatPsqlDB.Create(sp).Error; err != nil {
			appLogger.Errorf("Error creating dummy service professional: %v", err)
			return
		}
	}

	// Create dummy service-local relationships
	serviceLocals := []*model.ServiceLocal{
		// Yoga service available in multiple locations
		{
			Id:        uuid.New(),
			ServiceId: services[0].Id, // Yoga
			LocalId:   locals[0].Id,   // Pabellón A - San Miguel
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[0].Id, // Yoga
			LocalId:   locals[1].Id,   // Pabellón A - San Juan de Lurigancho
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[0].Id, // Yoga
			LocalId:   locals[2].Id,   // Pabellón A - Santa María del Carmen
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[0].Id, // Yoga
			LocalId:   locals[4].Id,   // Studio Zen
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Lactario service in specific locations
		{
			Id:        uuid.New(),
			ServiceId: services[2].Id, // Lactario
			LocalId:   locals[0].Id,   // Pabellón A - San Miguel
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[2].Id, // Lactario
			LocalId:   locals[1].Id,   // Pabellón A - San Juan de Lurigancho
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Gym service in specific locations
		{
			Id:        uuid.New(),
			ServiceId: services[3].Id, // Gimnasio
			LocalId:   locals[3].Id,   // Pabellón B
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[3].Id, // Gimnasio
			LocalId:   locals[0].Id,   // Pabellón A - San Miguel
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Functional training in multiple locations
		{
			Id:        uuid.New(),
			ServiceId: services[4].Id, // Entrenamiento Funcional
			LocalId:   locals[1].Id,   // Pabellón A - San Juan de Lurigancho
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:        uuid.New(),
			ServiceId: services[4].Id, // Entrenamiento Funcional
			LocalId:   locals[3].Id,   // Pabellón B
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	for _, sl := range serviceLocals {
		if err := astroCatPsqlDB.Create(sl).Error; err != nil {
			appLogger.Errorf("Error creating dummy service local: %v", err)
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
		// Más sesiones para hoy con diferentes horarios
		{
			Id:              uuid.New(),
			Title:           "Morning Yoga Session",
			Date:            baseDate,
			StartTime:       baseDate.Add(6 * time.Hour), // 6:00 AM
			EndTime:         baseDate.Add(7 * time.Hour), // 7:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 2,
			Capacity:        15,
			SessionLink:     nil,
			ProfessionalId:  professionals[0].Id, // John - Yoga Trainer
			LocalId:         &locals[0].Id,       // Pabellón A - San Miguel
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Functional Training",
			Date:            baseDate,
			StartTime:       baseDate.Add(7 * time.Hour), // 7:00 AM
			EndTime:         baseDate.Add(8 * time.Hour), // 8:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 3,
			Capacity:        20,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[1].Id,       // Pabellón A - San Juan de Lurigancho
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Gym Session",
			Date:            baseDate,
			StartTime:       baseDate.Add(8 * time.Hour), // 8:00 AM
			EndTime:         baseDate.Add(9 * time.Hour), // 9:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 5,
			Capacity:        25,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[3].Id,       // Pabellón B
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Zen Yoga",
			Date:            baseDate,
			StartTime:       baseDate.Add(9 * time.Hour),  // 9:00 AM
			EndTime:         baseDate.Add(10 * time.Hour), // 10:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 1,
			Capacity:        18,
			SessionLink:     nil,
			ProfessionalId:  professionals[3].Id, // Laura - Advanced Yoga Trainer
			LocalId:         &locals[4].Id,       // Studio Zen
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		// Sesiones para mañana con más horarios
		{
			Id:              uuid.New(),
			Title:           "Early Morning Yoga",
			Date:            baseDate.Add(24 * time.Hour),             // Tomorrow
			StartTime:       baseDate.Add(24*time.Hour + 5*time.Hour), // 5:00 AM
			EndTime:         baseDate.Add(24*time.Hour + 6*time.Hour), // 6:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 0,
			Capacity:        15,
			SessionLink:     nil,
			ProfessionalId:  professionals[0].Id, // John - Yoga Trainer
			LocalId:         &locals[0].Id,       // Pabellón A - San Miguel
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Morning Gym",
			Date:            baseDate.Add(24 * time.Hour),             // Tomorrow
			StartTime:       baseDate.Add(24*time.Hour + 6*time.Hour), // 6:00 AM
			EndTime:         baseDate.Add(24*time.Hour + 7*time.Hour), // 7:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 1,
			Capacity:        20,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[1].Id,       // Pabellón A - San Juan de Lurigancho
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Yoga Flow",
			Date:            baseDate.Add(24 * time.Hour),             // Tomorrow
			StartTime:       baseDate.Add(24*time.Hour + 7*time.Hour), // 7:00 AM
			EndTime:         baseDate.Add(24*time.Hour + 8*time.Hour), // 8:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 2,
			Capacity:        12,
			SessionLink:     nil,
			ProfessionalId:  professionals[3].Id, // Laura - Advanced Yoga Trainer
			LocalId:         &locals[2].Id,       // Pabellón A - Santa María del Carmen
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			Id:              uuid.New(),
			Title:           "Power Training",
			Date:            baseDate.Add(24 * time.Hour),             // Tomorrow
			StartTime:       baseDate.Add(24*time.Hour + 8*time.Hour), // 8:00 AM
			EndTime:         baseDate.Add(24*time.Hour + 9*time.Hour), // 9:00 AM
			State:           model.SessionStateScheduled,
			RegisteredCount: 3,
			Capacity:        25,
			SessionLink:     nil,
			ProfessionalId:  professionals[2].Id, // Pedro - Gym Trainer
			LocalId:         &locals[3].Id,       // Pabellón B
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
	district := "Lince"
	province := "Lima"
	region := "Lima"
	onboardings := []*model.Onboarding{
		{
			Id:             uuid.New(),
			PhoneNumber:    "123456789",
			DocumentType:   model.DocumentTypeDni,
			DocumentNumber: "12345678",
			PostalCode:     "15001",
			District:       &district,
			Province:       &province,
			Region:         &region,
			Address:        "Main St 123, Near Central Park",
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
			PostalCode:     "15002",
			District:       &district,
			Province:       &province,
			Region:         &region,
			Address:        "Downtown Ave 456, Near Business Center",
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
