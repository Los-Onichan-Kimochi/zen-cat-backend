package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
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

// Creates a community.
func (c *Community) CreateCommunity(community *schemas.Community, updatedBy string) *schemas.Error {
	return c.Adapter.Community.CreatePostgresqlCommunity(community, updatedBy)
}

// Gets a community.
func (c *Community) GetCommunity(id uuid.UUID) (*schemas.Community, *schemas.Error) {
	return c.Adapter.Community.GetPostgresqlCommunity(id)
}

// Fetch all communities.
func (c *Community) FetchCommunities() ([]*schemas.Community, *schemas.Error) {
	return c.Adapter.Community.FetchPostgresqlCommunities()
}

// Updates a community.
func (c *Community) UpdateCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) *schemas.Error {
	return c.Adapter.Community.UpdatePostgresqlCommunity(id, name, purpose, imageUrl, updatedBy)
}
