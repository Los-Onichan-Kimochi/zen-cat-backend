package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ServiceLocal struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create ServiceLocal postgresql controller
func NewServiceLocalController(logger logging.Logger, postgresqlDB *gorm.DB) *ServiceLocal {
	return &ServiceLocal{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a service-local association.
func (cp *ServiceLocal) CreateServiceLocal(serviceLocal *model.ServiceLocal) error {
	return cp.PostgresqlDB.Create(serviceLocal).Error
}

// Gets a specific service-local association.
func (cp *ServiceLocal) GetServiceLocal(
	serviceId uuid.UUID,
	localId uuid.UUID,
) (*model.ServiceLocal, error) {
	var serviceLocal model.ServiceLocal
	result := cp.PostgresqlDB.Where("service_id = ? AND local_id = ?", serviceId, localId).
		First(&serviceLocal)
	if result.Error != nil {
		return nil, result.Error // Returns gorm.ErrRecordNotFound if not found
	}

	return &serviceLocal, nil
}

// Deletes a specific service-local association.
func (cp *ServiceLocal) DeleteServiceLocal(
	serviceId uuid.UUID,
	localId uuid.UUID,
) error {
	result := cp.PostgresqlDB.Where("service_id = ? AND local_id = ?", serviceId, localId).
		Delete(&model.ServiceLocal{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was deleted
	}

	return nil
}
