package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ServiceModelF struct {
	Id          *uuid.UUID
	Name        *string
	Description *string
	ImageUrl    *string
	IsVirtual   *bool
}

// Create a new service on DB
func NewServiceModel(db *gorm.DB, option ...ServiceModelF) *model.Service {
	service := &model.Service{
		Id:          uuid.New(),
		Name:        "Test Service",
		Description: "Test service description",
		ImageUrl:    "https://example.com/service.jpg",
		IsVirtual:   false,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			service.Id = *parameters.Id
		}
		if parameters.Name != nil {
			service.Name = *parameters.Name
		}
		if parameters.Description != nil {
			service.Description = *parameters.Description
		}
		if parameters.ImageUrl != nil {
			service.ImageUrl = *parameters.ImageUrl
		}
		if parameters.IsVirtual != nil {
			service.IsVirtual = *parameters.IsVirtual
		}
	}

	result := db.Create(service)
	if result.Error != nil {
		log.Fatalf("Error when trying to create service: %v", result.Error)
	}

	return service
}

// Create size number of new services on DB
func NewServiceModelBatch(
	db *gorm.DB,
	size int,
	option ...ServiceModelF,
) []*model.Service {
	services := []*model.Service{}
	for i := 0; i < size; i++ {
		var service *model.Service
		if len(option) > 0 {
			service = NewServiceModel(db, option[0])
		} else {
			service = NewServiceModel(db)
		}
		services = append(services, service)
	}
	return services
}
