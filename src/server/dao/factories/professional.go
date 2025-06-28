package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ProfessionalModelF struct {
	Id             *uuid.UUID
	Name           *string
	FirstLastName  *string
	SecondLastName *string
	Specialty      *string
	Email          *string
	PhoneNumber    *string
	Type           *model.ProfessionalType
	ImageUrl       *string
}

// Create a new professional on DB
func NewProfessionalModel(db *gorm.DB, option ...ProfessionalModelF) *model.Professional {
	secondLastName := "Smith"
	professional := &model.Professional{
		Id:             uuid.New(),
		Name:           "Test Professional",
		FirstLastName:  "Test",
		SecondLastName: &secondLastName,
		Specialty:      "General Medicine",
		Email:          "professional@example.com",
		PhoneNumber:    "+1234567890",
		Type:           model.ProfessionalTypeMedic,
		ImageUrl:       "https://example.com/professional.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			professional.Id = *parameters.Id
		}
		if parameters.Name != nil {
			professional.Name = *parameters.Name
		}
		if parameters.FirstLastName != nil {
			professional.FirstLastName = *parameters.FirstLastName
		}
		if parameters.SecondLastName != nil {
			professional.SecondLastName = parameters.SecondLastName
		}
		if parameters.Specialty != nil {
			professional.Specialty = *parameters.Specialty
		}
		if parameters.Email != nil {
			professional.Email = *parameters.Email
		}
		if parameters.PhoneNumber != nil {
			professional.PhoneNumber = *parameters.PhoneNumber
		}
		if parameters.Type != nil {
			professional.Type = *parameters.Type
		}
		if parameters.ImageUrl != nil {
			professional.ImageUrl = *parameters.ImageUrl
		}
	}

	result := db.Create(professional)
	if result.Error != nil {
		log.Fatalf("Error when trying to create professional: %v", result.Error)
	}

	return professional
}

// Create size number of new professionals on DB
func NewProfessionalModelBatch(
	db *gorm.DB,
	size int,
	option ...ProfessionalModelF,
) []*model.Professional {
	professionals := []*model.Professional{}
	for i := 0; i < size; i++ {
		var professional *model.Professional
		if len(option) > 0 {
			professional = NewProfessionalModel(db, option[0])
		} else {
			professional = NewProfessionalModel(db)
		}
		professionals = append(professionals, professional)
	}
	return professionals
}
