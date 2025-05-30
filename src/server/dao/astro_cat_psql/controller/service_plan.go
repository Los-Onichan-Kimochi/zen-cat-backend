package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type ServiceProfessional struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create ServiceProfessional postgresql controller
func NewServiceProfessionalController(logger logging.Logger, postgresqlDB *gorm.DB) *ServiceProfessional {
	return &ServiceProfessional{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a service-professional association.
func (cp *ServiceProfessional) CreateServiceProfessional(serviceProfessional *model.ServiceProfessional) error {
	return cp.PostgresqlDB.Create(serviceProfessional).Error
}

// Gets a specific service-professional association.
func (cp *ServiceProfessional) GetServiceProfessional(
	serviceId uuid.UUID,
	professionalId uuid.UUID,
) (*model.ServiceProfessional, error) {
	var serviceProfessional model.ServiceProfessional
	result := cp.PostgresqlDB.Where("service_id = ? AND professional_id = ?", serviceId, professionalId).
		First(&serviceProfessional)
	if result.Error != nil {
		return nil, result.Error // Returns gorm.ErrRecordNotFound if not found
	}

	return &serviceProfessional, nil
}

// Deletes a specific service-professional association.
func (cp *ServiceProfessional) DeleteServiceProfessional(
	serviceId uuid.UUID,
	professionalId uuid.UUID,
) error {
	result := cp.PostgresqlDB.Where("service_id = ? AND professional_id = ?", serviceId, professionalId).
		Delete(&model.ServiceProfessional{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was deleted
	}

	return nil
}