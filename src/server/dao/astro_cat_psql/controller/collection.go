package controller

import (
	"fmt"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	"onichankimochi.com/astro_cat_backend/src/server/utils/psql"
)

type AstroCatPsqlCollection struct {
	Logger               logging.Logger
	Community            *Community
	Professional         *Professional
	Local                *Local
	User                 *User
	Onboarding           *Onboarding
	Membership           *Membership
	Service              *Service
	Plan                 *Plan
	CommunityPlan        *CommunityPlan
	CommunityService     *CommunityService
	ServiceLocal         *ServiceLocal
	ServiceProfessional  *ServiceProfessional
	Session              *Session
	Reservation          *Reservation
	AuditLog             *AuditLog
	MembershipSuspension *MembershipSuspension
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
		Logger:               logger,
		Community:            NewCommunityController(logger, postgresqlDB),
		Professional:         NewProfessionalController(logger, postgresqlDB),
		Local:                NewLocalController(logger, postgresqlDB),
		User:                 NewUserController(logger, postgresqlDB),
		Onboarding:           NewOnboardingController(logger, postgresqlDB),
		Membership:           NewMembershipController(logger, postgresqlDB),
		Service:              NewServiceController(logger, postgresqlDB),
		Plan:                 NewPlanController(logger, postgresqlDB),
		CommunityPlan:        NewCommunityPlanController(logger, postgresqlDB),
		CommunityService:     NewCommunityServiceController(logger, postgresqlDB),
		ServiceLocal:         NewServiceLocalController(logger, postgresqlDB),
		ServiceProfessional:  NewServiceProfessionalController(logger, postgresqlDB),
		Session:              NewSessionController(logger, postgresqlDB),
		Reservation:          NewReservationController(logger, postgresqlDB),
		AuditLog:             NewAuditLogController(logger, postgresqlDB),
		MembershipSuspension: NewMembershipSuspensionController(logger, postgresqlDB),
	}, postgresqlDB
}

// Helper function to create AstroCat tables
func createTables(astroCatPsqlDB *gorm.DB) {
	fmt.Println("Starting table creation...")

	fmt.Println("Creating Plan table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Plan{}); err != nil {
		fmt.Printf("Error creating Plan table: %v\n", err)
		panic(err)
	}
	fmt.Println("Plan table created successfully")

	fmt.Println("Creating Template table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Template{}); err != nil {
		fmt.Printf("Error creating Template table: %v\n", err)
		panic(err)
	}
	fmt.Println("Template table created successfully")

	fmt.Println("Creating Local table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Local{}); err != nil {
		fmt.Printf("Error creating Local table: %v\n", err)
		panic(err)
	}
	fmt.Println("Local table created successfully")

	fmt.Println("Creating Professional table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Professional{}); err != nil {
		fmt.Printf("Error creating Professional table: %v\n", err)
		panic(err)
	}
	fmt.Println("Professional table created successfully")

	fmt.Println("Creating Onboarding table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Onboarding{}); err != nil {
		fmt.Printf("Error creating Onboarding table: %v\n", err)
		panic(err)
	}
	fmt.Println("Onboarding table created successfully")

	fmt.Println("Creating User table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.User{}); err != nil {
		fmt.Printf("Error creating User table: %v\n", err)
		panic(err)
	}
	fmt.Println("User table created successfully")

	fmt.Println("Creating Community table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Community{}); err != nil {
		fmt.Printf("Error creating Community table: %v\n", err)
		panic(err)
	}
	fmt.Println("Community table created successfully")

	fmt.Println("Creating Membership table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Membership{}); err != nil {
		fmt.Printf("Error creating Membership table: %v\n", err)
		panic(err)
	}
	fmt.Println("Membership table created successfully")

	fmt.Println("Creating MembershipSuspension table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.MembershipSuspension{}); err != nil {
		fmt.Printf("Error creating MembershipSuspension table: %v\n", err)
		panic(err)
	}
	fmt.Println("MembershipSuspension table created successfully")

	fmt.Println("Creating Service table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Service{}); err != nil {
		fmt.Printf("Error creating Service table: %v\n", err)
		panic(err)
	}
	fmt.Println("Service table created successfully")

	fmt.Println("Creating Session table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Session{}); err != nil {
		fmt.Printf("Error creating Session table: %v\n", err)
		panic(err)
	}
	fmt.Println("Session table created successfully")

	fmt.Println("Creating Reservation table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.Reservation{}); err != nil {
		fmt.Printf("Error creating Reservation table: %v\n", err)
		panic(err)
	}
	fmt.Println("Reservation table created successfully")

	fmt.Println("Creating CommunityService table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityService{}); err != nil {
		fmt.Printf("Error creating CommunityService table: %v\n", err)
		panic(err)
	}
	fmt.Println("CommunityService table created successfully")

	fmt.Println("Creating CommunityPlan table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.CommunityPlan{}); err != nil {
		fmt.Printf("Error creating CommunityPlan table: %v\n", err)
		panic(err)
	}
	fmt.Println("CommunityPlan table created successfully")

	fmt.Println("Creating ServiceLocal table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.ServiceLocal{}); err != nil {
		fmt.Printf("Error creating ServiceLocal table: %v\n", err)
		panic(err)
	}
	fmt.Println("ServiceLocal table created successfully")

	fmt.Println("Creating ServiceProfessional table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.ServiceProfessional{}); err != nil {
		fmt.Printf("Error creating ServiceProfessional table: %v\n", err)
		panic(err)
	}
	fmt.Println("ServiceProfessional table created successfully")

	fmt.Println("Creating AuditLog table...")
	if err := astroCatPsqlDB.AutoMigrate(&model.AuditLog{}); err != nil {
		fmt.Printf("Error creating AuditLog table: %v\n", err)
		panic(err)
	}
	fmt.Println("AuditLog table created successfully")

	fmt.Println("All tables created successfully!")
}

// Helper function to drop all AstroCat tables without considering constraints
func dropAllTables(astroCatPsqlDB *gorm.DB) {
	// Disable foreign key constraints temporarily
	astroCatPsqlDB.Exec("SET CONSTRAINTS ALL DEFERRED")

	// Drop all tables in reverse order of dependencies
	tablesToDrop := []string{
		"audit_logs",
		"service_professionals",
		"service_locals",
		"community_plans",
		"community_services",
		"reservations",
		"sessions",
		"templates",
		"onboardings",
		"membership_suspensions",
		"memberships",
		"users",
		"communities",
		"services",
		"plans",
		"professionals",
		"locals",
	}

	for _, tableName := range tablesToDrop {
		if err := astroCatPsqlDB.Exec("DROP TABLE IF EXISTS " + tableName + " CASCADE").Error; err != nil {
			fmt.Printf("Warning: Error dropping table %s: %v\n", tableName, err)
		} else {
			fmt.Printf("Dropped table: %s\n", tableName)
		}
	}

	// Reactivate foreign key constraints
	astroCatPsqlDB.Exec("SET CONSTRAINTS ALL IMMEDIATE")
}
