package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Community struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Community adapter
func NewCommunityAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Community {
	return &Community{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Creates a community into postgresql DB.
func (c *Community) CreatePostgresqlCommunity(
	community *schemas.Community,
	updatedBy string,
) *schemas.Error {
	if updatedBy == "" {
		return &schemas.BadRequestError.InvalidUpdatedByValue
	}

	communityModel := &model.Community{
		Id:                  community.Id,
		Name:                community.Name,
		Purpose:             community.Purpose,
		ImageUrl:            community.ImageUrl,
		NumberSubscriptions: 0, // Default number of subscriptions
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := c.DaoPostgresql.Community.CreateCommunity(communityModel); err != nil {
		return &schemas.BadRequestError.CommunityNotCreated
	}

	return nil
}

// Gets a community from postgresql DB and adapts it to a Community schema.
func (c *Community) GetPostgresqlCommunity(id uuid.UUID) (*schemas.Community, *schemas.Error) {
	communityModel, err := c.DaoPostgresql.Community.GetCommunity(id)
	if err != nil {
		return nil, &schemas.ObjectNotFoundError.CommunityNotFound
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}

// Fetch communities from postgresql DB and adapts them to a Community schema.
func (c *Community) FetchPostgresqlCommunities() ([]*schemas.Community, *schemas.Error) {
	communitiesModel, err := c.DaoPostgresql.Community.FetchCommunities()
	if err != nil {
		return nil, &schemas.ObjectNotFoundError.CommunityNotFound
	}

	communities := make([]*schemas.Community, len(communitiesModel))
	for i, communityModel := range communitiesModel {
		communities[i] = &schemas.Community{
			Id:                  communityModel.Id,
			Name:                communityModel.Name,
			Purpose:             communityModel.Purpose,
			ImageUrl:            communityModel.ImageUrl,
			NumberSubscriptions: communityModel.NumberSubscriptions,
		}
	}

	return communities, nil
}

// Updates a community given fields in postgresql DB.
func (c *Community) UpdatePostgresqlCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) *schemas.Error {
	if updatedBy == "" {
		return &schemas.BadRequestError.InvalidUpdatedByValue
	}

	if err := c.DaoPostgresql.Community.UpdateCommunity(id, name, purpose, imageUrl, updatedBy); err != nil {
		return &schemas.BadRequestError.CommunityNotUpdated
	}

	return nil
}
