package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
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

// Gets a community from postgresql DB and adapts it to a Community schema.
func (c *Community) GetPostgresqlCommunity(
	communityId uuid.UUID,
) (*schemas.Community, *errors.Error) {
	communityModel, err := c.DaoPostgresql.Community.GetCommunity(communityId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityNotFound
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
func (c *Community) FetchPostgresqlCommunities() ([]*schemas.Community, *errors.Error) {
	communitiesModel, err := c.DaoPostgresql.Community.FetchCommunities()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityNotFound
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

// Creates a community into postgresql DB and returns it.
func (c *Community) CreatePostgresqlCommunity(
	name string,
	purpose string,
	imageUrl string,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityModel := &model.Community{
		Id:                  uuid.New(),
		Name:                name,
		Purpose:             purpose,
		ImageUrl:            imageUrl,
		NumberSubscriptions: 0, // Default number of initial subscriptions
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := c.DaoPostgresql.Community.CreateCommunity(communityModel); err != nil {
		return nil, &errors.BadRequestError.CommunityNotCreated
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}

// Updates a community given fields in postgresql DB and returns it.
func (c *Community) UpdatePostgresqlCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityModel, err := c.DaoPostgresql.Community.UpdateCommunity(
		id,
		name,
		purpose,
		imageUrl,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.CommunityNotUpdated
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}
