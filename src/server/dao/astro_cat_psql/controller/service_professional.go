package controller

import (
	"fmt"
	"strings"

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
// Gorm does not return the Id of a soft deleted record.
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

// Fetch all service-professional associations, filtered by
//
//   - `serviceId` if provided.
func (cp *ServiceProfessional) GetProfessionalsByServiceId(
	serviceId uuid.UUID,
) ([]*model.Professional, error) {
	var serviceProfessionals []*model.ServiceProfessional

	// Realizamos la consulta en la base de datos para obtener las asociaciones de servicio y profesional
	query := cp.PostgresqlDB.Model(&model.ServiceProfessional{})

	query = query.Where("service_id = ?", serviceId)

	// Ejecutamos la consulta y almacenamos las asociaciones
	if err := query.Find(&serviceProfessionals).Error; err != nil {
		return nil, err
	}

	// Obtenemos los IDs de los profesionales asociados
	var professionalIds []uuid.UUID
	for _, serviceProfessional := range serviceProfessionals {
		professionalIds = append(professionalIds, serviceProfessional.ProfessionalId)
	}

	// Ahora realizamos una segunda consulta para obtener los servicios
	var professionals []*model.Professional
	if err := cp.PostgresqlDB.Where("id IN (?)", professionalIds).Find(&professionals).Error; err != nil {
		return nil, err
	}

	return professionals, nil
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

// Creates multiple service-professional associations.
func (cp *ServiceProfessional) BulkCreateServiceProfessionals(
	serviceProfessionals []*model.ServiceProfessional,
) error {
	if len(serviceProfessionals) == 0 {
		return nil
	}

	var conditions []string
	var args []interface{}
	var count int64

	for _, serviceProfessional := range serviceProfessionals {
		conditions = append(conditions, "(service_id = ? AND professional_id = ?)")
		args = append(args, serviceProfessional.ServiceId, serviceProfessional.ProfessionalId)
	}
	result := cp.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Find(&model.ServiceProfessional{}).
		Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("one or more service-professional associations already exist")
	}

	err := cp.PostgresqlDB.Create(serviceProfessionals).Error
	if err != nil {
		return err
	}

	return nil
}

// Fetch all service-professional associations, filtered by
//
//   - `serviceId` if provided.
//   - `professionalId` if provided.
func (cp *ServiceProfessional) FetchServiceProfessionals(
	serviceId *uuid.UUID,
	professionalId *uuid.UUID,
) ([]*model.ServiceProfessional, error) {
	var serviceProfessionals []*model.ServiceProfessional

	query := cp.PostgresqlDB.Model(&model.ServiceProfessional{})

	if serviceId != nil {
		query = query.Where("service_id = ?", serviceId)
	}
	if professionalId != nil {
		query = query.Where("professional_id = ?", professionalId)
	}

	if err := query.Find(&serviceProfessionals).Error; err != nil {
		return nil, err
	}

	return serviceProfessionals, nil
}

// Bulk deletes multiple service-professional associations.
func (cp *ServiceProfessional) BulkDeleteServiceProfessionals(
	serviceProfessionals []*model.ServiceProfessional,
) error {
	// Build the WHERE clause for the bulk delete
	var conditions []string
	var args []interface{}
	for _, serviceProfessional := range serviceProfessionals {
		conditions = append(conditions, "(service_id = ? AND professional_id = ?)")
		args = append(args, serviceProfessional.ServiceId, serviceProfessional.ProfessionalId)
	}

	// Execute the bulk delete
	result := cp.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Delete(&model.ServiceProfessional{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}
