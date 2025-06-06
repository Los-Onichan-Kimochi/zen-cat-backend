package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

// Gets a community model given params.
func (c *Community) GetCommunity(communityId uuid.UUID) (*model.Community, error) {
	community := &model.Community{}

	result := c.PostgresqlDB.First(&community, "id = ?", communityId)
	if result.Error != nil {
		return nil, result.Error
	}

	return community, nil
}

// TODO: Add filters and sorting.
// Fetch all communities.
func (c *Community) FetchCommunities() ([]*model.Community, error) {
	communities := []*model.Community{}

	result := c.PostgresqlDB.Find(&communities)
	if result.Error != nil {
		return nil, result.Error
	}

	return communities, nil
}

// Creates a community given its model.
func (c *Community) CreateCommunity(community *model.Community) error {
	return c.PostgresqlDB.Create(community).Error
}

// Updates community given fields to update.
func (c *Community) UpdateCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) (*model.Community, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
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

	// Check if there are any fields to update
	var community model.Community
	if len(updateFields) == 1 {
		if err := c.PostgresqlDB.First(&community, "id = ?", id).Error; err != nil {
			return nil, err
		}

		return &community, nil
	}

	// Perform the update and return the model
	result := c.PostgresqlDB.Model(&community).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &community, nil
}

// Soft deletes a community given its ID.
func (c *Community) DeleteCommunity(communityId uuid.UUID) error {
	result := c.PostgresqlDB.Delete(&model.Community{}, "id = ?", communityId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Creates communities given their models.
func (c *Community) BulkCreateCommunities(communities []*model.Community) error {
	return c.PostgresqlDB.Create(&communities).Error
}

// Batch deletes multiple communities given their IDs.
func (c *Community) BulkDeleteCommunities(communityIds []uuid.UUID) error {
	if len(communityIds) == 0 {
		return nil
	}

	result := c.PostgresqlDB.Where("id IN ?", communityIds).Delete(&model.Community{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
