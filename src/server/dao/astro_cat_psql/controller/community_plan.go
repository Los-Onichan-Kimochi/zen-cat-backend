package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type CommunityPlan struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create CommunityPlan postgresql controller
func NewCommunityPlanController(logger logging.Logger, postgresqlDB *gorm.DB) *CommunityPlan {
	return &CommunityPlan{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a community-plan association.
func (cp *CommunityPlan) CreateCommunityPlan(communityPlan *model.CommunityPlan) error {
	return cp.PostgresqlDB.Create(communityPlan).Error
}

// Gets a specific community-plan association.
func (cp *CommunityPlan) GetCommunityPlan(
	communityId uuid.UUID,
	planId uuid.UUID,
) (*model.CommunityPlan, error) {
	var communityPlan model.CommunityPlan
	result := cp.PostgresqlDB.Where("community_id = ? AND plan_id = ?", communityId, planId).
		First(&communityPlan)
	if result.Error != nil {
		return nil, result.Error // Returns gorm.ErrRecordNotFound if not found
	}

	return &communityPlan, nil
}
