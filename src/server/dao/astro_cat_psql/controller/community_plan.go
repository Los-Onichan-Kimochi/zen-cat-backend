package controller

import (
	"errors"
	"fmt"

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

// Deletes a specific community-plan association.
func (cp *CommunityPlan) DeleteCommunityPlan(
	communityId uuid.UUID,
	planId uuid.UUID,
) error {
	result := cp.PostgresqlDB.Where("community_id = ? AND plan_id = ?", communityId, planId).
		Delete(&model.CommunityPlan{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // Indicate that no record was deleted
	}

	return nil
}

// Creates multiple community-plan associations.
func (cp *CommunityPlan) BulkCreateCommunityPlans(communityPlans []*model.CommunityPlan) error {
	err := cp.PostgresqlDB.Create(communityPlans).Error
	if err != nil {
		// Check if it's a duplicate key error
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return fmt.Errorf("one or more community-plan associations already exist: %w", err)
		}
		return err
	}
	return nil
}

// Fetch all community-plan associations, filtered by
//
//   - `communityId` if provided.
//   - `planId` if provided.
func (cp *CommunityPlan) FetchCommunityPlans(
	communityId *uuid.UUID,
	planId *uuid.UUID,
) ([]*model.CommunityPlan, error) {
	var communityPlans []*model.CommunityPlan

	query := cp.PostgresqlDB.Model(&model.CommunityPlan{})

	if communityId != nil {
		query = query.Where("community_id = ?", communityId)
	}
	if planId != nil {
		query = query.Where("plan_id = ?", planId)
	}

	if err := query.Find(&communityPlans).Error; err != nil {
		return nil, err
	}

	return communityPlans, nil
}
