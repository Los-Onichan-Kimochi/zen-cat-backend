package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ServiceLocalModelF struct {
	Id        *uuid.UUID
	ServiceId *uuid.UUID
	LocalId   *uuid.UUID
}

// Create a new service local on DB
func NewServiceLocalModel(db *gorm.DB, option ...ServiceLocalModelF) *model.ServiceLocal {
	// Create default service if not provided
	service := NewServiceModel(db)

	// Create default local if not provided
	local := NewLocalModel(db)

	serviceLocal := &model.ServiceLocal{
		Id:        uuid.New(),
		ServiceId: service.Id,
		LocalId:   local.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			serviceLocal.Id = *parameters.Id
		}
		if parameters.ServiceId != nil {
			serviceLocal.ServiceId = *parameters.ServiceId
		}
		if parameters.LocalId != nil {
			serviceLocal.LocalId = *parameters.LocalId
		}
	}

	result := db.Create(serviceLocal)
	if result.Error != nil {
		log.Fatalf("Error when trying to create service local: %v", result.Error)
	}

	return serviceLocal
}

// Create size number of new service locals on DB
func NewServiceLocalModelBatch(
	db *gorm.DB,
	size int,
	option ...ServiceLocalModelF,
) []*model.ServiceLocal {
	serviceLocals := []*model.ServiceLocal{}
	for i := 0; i < size; i++ {
		var serviceLocal *model.ServiceLocal
		if len(option) > 0 {
			serviceLocal = NewServiceLocalModel(db, option[0])
		} else {
			serviceLocal = NewServiceLocalModel(db)
		}
		serviceLocals = append(serviceLocals, serviceLocal)
	}
	return serviceLocals
}
