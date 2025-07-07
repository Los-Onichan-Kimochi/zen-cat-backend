package adapter

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type MembershipSuspension struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

func NewMembershipSuspensionAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *MembershipSuspension {
	return &MembershipSuspension{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (m *MembershipSuspension) CreatePostgresqlMembershipSuspension(membershipId uuid.UUID) (*schemas.MembershipSuspension, *errors.Error) {
	suspensionModel, err := m.DaoPostgresql.MembershipSuspension.CreateMembershipSuspension(membershipId)
	if err != nil {
		return nil, &errors.BadRequestError.MembershipSuspensionNotCreated
	}

	return m.convertModelToSchema(suspensionModel), nil
}

func (m *MembershipSuspension) GetLatestOpenPostgresqlMembershipSuspension(membershipId uuid.UUID) (*model.MembershipSuspension, *errors.Error) {
	suspensionModel, err := m.DaoPostgresql.MembershipSuspension.GetLatestOpenMembershipSuspension(membershipId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.MembershipSuspensionNotFound
		}
		return nil, &errors.InternalServerError.DatabaseError
	}

	return suspensionModel, nil
}

func (m *MembershipSuspension) UpdatePostgresqlMembershipSuspension(suspension *model.MembershipSuspension) (*schemas.MembershipSuspension, *errors.Error) {
	err := m.DaoPostgresql.MembershipSuspension.UpdateMembershipSuspension(suspension)
	if err != nil {
		return nil, &errors.BadRequestError.MembershipSuspensionNotUpdated
	}
	return m.convertModelToSchema(suspension), nil
}

func (m *MembershipSuspension) convertModelToSchema(suspensionModel *model.MembershipSuspension) *schemas.MembershipSuspension {
	return &schemas.MembershipSuspension{
		Id:           suspensionModel.Id,
		MembershipId: suspensionModel.MembershipId,
		SuspendedAt:  suspensionModel.SuspendedAt,
		ResumedAt:    suspensionModel.ResumedAt,
	}
}
