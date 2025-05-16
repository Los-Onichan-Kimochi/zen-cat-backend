package controller

import (
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
