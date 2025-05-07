package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Community struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Community postgresql controller
func NewCommunityController(logger logging.Logger, postgresqlDB *gorm.DB) *Community {
	return &Community{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Creates a community given its model.
func (p *Community) CreateCommunity(community *model.Community) error {
	return p.PostgresqlDB.Create(community).Error
}

// Gets a community model given params.
func (p *Community) GetCommunity(id uuid.UUID) (*model.Community, error) {
	community := &model.Community{}

	result := p.PostgresqlDB.Where(id).Find(&community)
	if result.Error != nil {
		return nil, result.Error
	}

	return community, nil
}

// TODO: Add sorting.
// Fetch all communities.
func (p *Community) FetchCommunities() ([]*model.Community, error) {
	communities := []*model.Community{}

	result := p.PostgresqlDB.Find(&communities)
	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
}

// Updates community given fields to update.
func (r *Community) UpdateCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) error {
	updateFields := map[string]any{}

	// Check if there are any fields to update
	if len(updateFields) == 0 {
		return nil
	}

	if name != nil {
		updateFields["name"] = *name
	}
	if purpose != nil {
		updateFields["purpose"] = *purpose
	}
	if imageUrl != nil {
		updateFields["image_url"] = *imageUrl
	}

	updateFields["updated_by"] = updatedBy

	return r.PostgresqlDB.Model(&model.Community{}).Where(id).Updates(updateFields).Error
}
