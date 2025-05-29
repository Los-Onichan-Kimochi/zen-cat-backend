package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Service struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Service postgresql controller
func NewServiceController(logger logging.Logger, postgresqlDB *gorm.DB) *Service {
	return &Service{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a service model given params.
func (c *Service) GetService(serviceId uuid.UUID) (*model.Service, error) {
	service := &model.Service{}

	result := c.PostgresqlDB.First(&service, "id = ?", serviceId)
	if result.Error != nil {
		return nil, result.Error
	}

	return service, nil
}

// Fetch all services, filtered by `ids` if provided.
func (c *Service) FetchServices(ids []uuid.UUID) ([]*model.Service, error) {
	services := []*model.Service{}

	query := c.PostgresqlDB.Model(&model.Service{})

	if len(ids) > 0 {
		query = query.Where("id IN (?)", ids)
	}

	if err := query.Find(&services).Error; err != nil {
		return nil, err
	}

	return services, nil
}

// Creates a service given its model.
func (c *Service) CreateService(service *model.Service) error {
	return c.PostgresqlDB.Create(service).Error
}

// Updates service given fields to update.
func (c *Service) UpdateService(
	id uuid.UUID,
	name *string,
	description *string,
	imageUrl *string,
	isVirtual *bool,
	updatedBy string,
) (*model.Service, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}

	if name != nil {
		updateFields["name"] = *name
	}
	if description != nil {
		updateFields["description"] = *description
	}
	if imageUrl != nil {
		updateFields["image_url"] = *imageUrl
	}
	if isVirtual != nil {
		updateFields["is_virtual"] = *isVirtual
	}

	// Check if there are any fields to update
	var service model.Service
	if len(updateFields) == 1 {
		if err := c.PostgresqlDB.First(&service, "id = ?", id).Error; err != nil {
			return nil, err
		}

		return &service, nil
	}

	// Perform the update and return the model
	result := c.PostgresqlDB.Model(&service).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &service, nil
}

// Deletes a service given its id.
func (l *Service) DeleteService(id uuid.UUID) error {
	result := l.PostgresqlDB.Delete(&model.Service{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
