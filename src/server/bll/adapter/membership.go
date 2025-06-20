package adapter

import (
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPsql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Membership struct {
	logger        logging.Logger
	DaoPostgresql *daoPsql.AstroCatPsqlCollection
}

// Create Membership adapter
func NewMembershipAdapter(
	logger logging.Logger,
	daoPostgresql *daoPsql.AstroCatPsqlCollection,
) *Membership {
	return &Membership{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a specific reservation and adapts it.
func (m *Membership) GetPostgresqlMembership(
	membershipId uuid.UUID,
) (*schemas.Membership, *errors.Error) {
	membershipModel, err := m.DaoPostgresql.Membership.GetMembership(membershipId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	return &schemas.Membership{
		Id:          membershipModel.Id,
		Description: membershipModel.Description,
		StartDate:   membershipModel.StartDate,
		EndDate:     membershipModel.EndDate,
		Status:      schemas.MembershipStatus(membershipModel.Status),
		Community:   schemas.Community{Id: membershipModel.CommunityId},
		User:        schemas.User{Id: membershipModel.UserId},
		Plan:        schemas.Plan{Id: membershipModel.PlanId},
	}, nil
}