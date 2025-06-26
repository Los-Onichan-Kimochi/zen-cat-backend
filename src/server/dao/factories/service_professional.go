package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ServiceProfessionalModelF struct {
	Id             *uuid.UUID
	ServiceId      *uuid.UUID
	ProfessionalId *uuid.UUID
}

// Create a new service professional on DB
func NewServiceProfessionalModel(db *gorm.DB, option ...ServiceProfessionalModelF) *model.ServiceProfessional {
	// Create default service if not provided
	service := NewServiceModel(db)

	// Create default professional if not provided
	professional := NewProfessionalModel(db)

	serviceProfessional := &model.ServiceProfessional{
		Id:             uuid.New(),
		ServiceId:      service.Id,
		ProfessionalId: professional.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			serviceProfessional.Id = *parameters.Id
		}
		if parameters.ServiceId != nil {
			serviceProfessional.ServiceId = *parameters.ServiceId
		}
		if parameters.ProfessionalId != nil {
			serviceProfessional.ProfessionalId = *parameters.ProfessionalId
		}
	}

	result := db.Create(serviceProfessional)
	if result.Error != nil {
		log.Fatalf("Error when trying to create service professional: %v", result.Error)
	}

	return serviceProfessional
}

// Create size number of new service professionals on DB
func NewServiceProfessionalModelBatch(
	db *gorm.DB,
	size int,
	option ...ServiceProfessionalModelF,
) []*model.ServiceProfessional {
	serviceProfessionals := []*model.ServiceProfessional{}
	for i := 0; i < size; i++ {
		var serviceProfessional *model.ServiceProfessional
		if len(option) > 0 {
			serviceProfessional = NewServiceProfessionalModel(db, option[0])
		} else {
			serviceProfessional = NewServiceProfessionalModel(db)
		}
		serviceProfessionals = append(serviceProfessionals, serviceProfessional)
	}
	return serviceProfessionals
}
