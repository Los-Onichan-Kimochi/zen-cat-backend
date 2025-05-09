package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Community struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Community controller
func NewCommunityController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Community {
	return &Community{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a community.
func (c *Community) GetCommunity(communityId uuid.UUID) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.GetPostgresqlCommunity(communityId)
}

// Fetch all communities.
func (c *Community) FetchCommunities() (*schemas.Communities, *errors.Error) {
	communities, err := c.Adapter.Community.FetchPostgresqlCommunities()
	if err != nil {
		return nil, err
	}

	return &schemas.Communities{Communities: communities}, nil
}

// Creates a community.
func (c *Community) CreateCommunity(
	createCommunityData schemas.CreateCommunityRequest,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.CreatePostgresqlCommunity(
		createCommunityData.Name,
		createCommunityData.Purpose,
		createCommunityData.ImageUrl,
		updatedBy,
	)
}

// Updates a community.
func (c *Community) UpdateCommunity(
	communityId uuid.UUID,
	updateCommunityData schemas.UpdateCommunityRequest,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.UpdatePostgresqlCommunity(
		communityId,
		updateCommunityData.Name,
		updateCommunityData.Purpose,
		updateCommunityData.ImageUrl,
		updatedBy,
	)
}

// TODO: Add BulkCreateCommunities (Batch)
