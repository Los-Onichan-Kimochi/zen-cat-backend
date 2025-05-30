package controller

import (
	"fmt"
	"strings"

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
// Gorm does not return the Id of a soft deleted record.
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

// Creates multiple service-local associations.
func (cp *ServiceLocal) BulkCreateServiceLocals(
	serviceLocals []*model.ServiceLocal,
) error {
	if len(serviceLocals) == 0 {
		return nil
	}

	var conditions []string
	var args []interface{}
	var count int64

	for _, serviceLocal := range serviceLocals {
		conditions = append(conditions, "(service_id = ? AND local_id = ?)")
		args = append(args, serviceLocal.ServiceId, serviceLocal.LocalId)
	}
	result := cp.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Find(&model.ServiceLocal{}).
		Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("one or more service-local associations already exist")
	}

	err := cp.PostgresqlDB.Create(serviceLocals).Error
	if err != nil {
		return err
	}

	return nil
}

// Fetch all service-local associations, filtered by
//
//   - `serviceId` if provided.
//   - `localId` if provided.
func (cp *ServiceLocal) FetchServiceLocals(
	serviceId *uuid.UUID,
	localId *uuid.UUID,
) ([]*model.ServiceLocal, error) {
	var serviceLocals []*model.ServiceLocal

	query := cp.PostgresqlDB.Model(&model.ServiceLocal{})

	if serviceId != nil {
		query = query.Where("service_id = ?", serviceId)
	}
	if localId != nil {
		query = query.Where("local_id = ?", localId)
	}

	if err := query.Find(&serviceLocals).Error; err != nil {
		return nil, err
	}

	return serviceLocals, nil
}

// Bulk deletes multiple service-local associations.
func (cp *ServiceLocal) BulkDeleteServiceLocals(
	serviceLocals []*model.ServiceLocal,
) error {
	// Build the WHERE clause for the bulk delete
	var conditions []string
	var args []interface{}
	for _, serviceLocal := range serviceLocals {
		conditions = append(conditions, "(service_id = ? AND local_id = ?)")
		args = append(args, serviceLocal.ServiceId, serviceLocal.LocalId)
	}

	// Execute the bulk delete
	result := cp.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Delete(&model.ServiceLocal{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}