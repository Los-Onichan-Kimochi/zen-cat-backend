package factories

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type OnboardingModelF struct {
	Id             *uuid.UUID
	DocumentType   *model.DocumentType
	DocumentNumber *string
	PhoneNumber    *string
	BirthDate      *time.Time
	Gender         *model.Gender
	PostalCode     *string
	Address        *string
	District       *string
	Province       *string
	Region         *string
	UserId         *uuid.UUID
}

// Create a new onboarding on DB
func NewOnboardingModel(db *gorm.DB, option ...OnboardingModelF) *model.Onboarding {
	// Create default user if not provided
	user := NewUserModel(db)

	birthDate := time.Now().AddDate(-30, 0, 0) // 30 years ago
	gender := model.GenderMale
	district := "Test District"
	province := "Test Province"
	region := "Test Region"

	onboarding := &model.Onboarding{
		Id:             uuid.New(),
		DocumentType:   model.DocumentTypeDni,
		DocumentNumber: "12345678",
		PhoneNumber:    "+1234567890",
		BirthDate:      &birthDate,
		Gender:         &gender,
		PostalCode:     "12345",
		Address:        "123 Test St",
		District:       &district,
		Province:       &province,
		Region:         &region,
		UserId:         user.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			onboarding.Id = *parameters.Id
		}
		if parameters.DocumentType != nil {
			onboarding.DocumentType = *parameters.DocumentType
		}
		if parameters.DocumentNumber != nil {
			onboarding.DocumentNumber = *parameters.DocumentNumber
		}
		if parameters.PhoneNumber != nil {
			onboarding.PhoneNumber = *parameters.PhoneNumber
		}
		if parameters.BirthDate != nil {
			onboarding.BirthDate = parameters.BirthDate
		}
		if parameters.Gender != nil {
			onboarding.Gender = parameters.Gender
		}
		if parameters.PostalCode != nil {
			onboarding.PostalCode = *parameters.PostalCode
		}
		if parameters.Address != nil {
			onboarding.Address = *parameters.Address
		}
		if parameters.District != nil {
			onboarding.District = parameters.District
		}
		if parameters.Province != nil {
			onboarding.Province = parameters.Province
		}
		if parameters.Region != nil {
			onboarding.Region = parameters.Region
		}
		if parameters.UserId != nil {
			onboarding.UserId = *parameters.UserId
		}
	}

	result := db.Create(onboarding)
	if result.Error != nil {
		log.Fatalf("Error when trying to create onboarding: %v", result.Error)
	}

	return onboarding
}

// Create size number of new onboardings on DB
func NewOnboardingModelBatch(
	db *gorm.DB,
	size int,
	option ...OnboardingModelF,
) []*model.Onboarding {
	onboardings := []*model.Onboarding{}
	for i := 0; i < size; i++ {
		var onboarding *model.Onboarding
		if len(option) > 0 {
			onboarding = NewOnboardingModel(db, option[0])
		} else {
			onboarding = NewOnboardingModel(db)
		}
		onboardings = append(onboardings, onboarding)
	}
	return onboardings
}
