package controller

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type CommunityService struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create CommunityService postgresql controller
func NewCommunityServiceController(logger logging.Logger, postgresqlDB *gorm.DB) *CommunityService {
	return &CommunityService{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a community-service association.
func (cs *CommunityService) CreateCommunityService(communityService *model.CommunityService) error {
	return cs.PostgresqlDB.Create(communityService).Error
}

// Gets a specific community-service association.
// Gorm does not return the Id of a soft deleted record.
func (cs *CommunityService) GetCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) (*model.CommunityService, error) {
	var communityService model.CommunityService
	result := cs.PostgresqlDB.Where("community_id = ? AND service_id = ?", communityId, serviceId).
		First(&communityService)
	if result.Error != nil {
		return nil, result.Error // Returns gorm.ErrRecordNotFound if not found
	}

	return &communityService, nil
}

// Deletes a specific community-service association.
func (cs *CommunityService) DeleteCommunityService(
	communityId uuid.UUID,
	serviceId uuid.UUID,
) error {
	result := cs.PostgresqlDB.Where("community_id = ? AND service_id = ?", communityId, serviceId).
		Delete(&model.CommunityService{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was deleted
	}

	return nil
}

// Creates multiple community-service associations.
func (cs *CommunityService) BulkCreateCommunityServices(
	communityServices []*model.CommunityService,
) error {
	if len(communityServices) == 0 {
		return nil
	}

	var conditions []string
	var args []any
	var count int64

	for _, communityService := range communityServices {
		conditions = append(conditions, "(community_id = ? AND service_id = ?)")
		args = append(args, communityService.CommunityId, communityService.ServiceId)
	}
	result := cs.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Find(&model.CommunityService{}).
		Count(&count)
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("one or more community-service associations already exist")
	}

	err := cs.PostgresqlDB.Create(communityServices).Error
	if err != nil {
		return err
	}

	return nil
}

// Bulk deletes multiple community-service associations.
func (cs *CommunityService) BulkDeleteCommunityServices(
	communityServices []*model.CommunityService,
) error {
	if len(communityServices) == 0 {
		return nil
	}

	// Build the WHERE clause for the bulk delete
	var conditions []string
	var args []interface{}
	for _, communityService := range communityServices {
		conditions = append(conditions, "(community_id = ? AND service_id = ?)")
		args = append(args, communityService.CommunityId, communityService.ServiceId)
	}

	// Execute the bulk delete
	result := cs.PostgresqlDB.Where(strings.Join(conditions, " OR "), args...).
		Delete(&model.CommunityService{})
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// Fetch all community-service associations, filtered by
//
//   - `communityId` if provided.
//   - `serviceId` if provided.
func (cs *CommunityService) FetchCommunityServices(
	communityId *uuid.UUID,
	serviceId *uuid.UUID,
) ([]*model.CommunityService, error) {
	var communityServices []*model.CommunityService

	query := cs.PostgresqlDB.Model(&model.CommunityService{})

	if communityId != nil {
		query = query.Where("community_id = ?", communityId)
	}
	if serviceId != nil {
		query = query.Where("service_id = ?", serviceId)
	}

	if err := query.Find(&communityServices).Error; err != nil {
		return nil, err
	}

	return communityServices, nil
}
