package adapter

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Membership struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

func NewMembershipAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Membership {
	return &Membership{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (m *Membership) GetPostgresqlMembership(
	membershipId uuid.UUID,
) (*schemas.Membership, *errors.Error) {
	membershipModel, err := m.DaoPostgresql.Membership.GetMembership(membershipId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	// Convertir model a schema
	return m.convertModelToSchema(membershipModel), nil
}

func (m *Membership) GetPostgresqlMembershipsByUserId(
	userId uuid.UUID,
) ([]*schemas.Membership, *errors.Error) {
	membershipsModel, err := m.DaoPostgresql.Membership.GetMembershipsByUserId(userId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	memberships := make([]*schemas.Membership, len(membershipsModel))
	for i, membershipModel := range membershipsModel {
		memberships[i] = m.convertModelToSchema(membershipModel)
	}

	return memberships, nil
}

func (m *Membership) GetPostgresqlMembershipByUserAndCommunity(
	userId uuid.UUID,
	communityId uuid.UUID,
) (*schemas.Membership, *errors.Error) {
	membershipModel, err := m.DaoPostgresql.Membership.GetMembershipByUserAndCommunity(userId, communityId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	return m.convertModelToSchema(membershipModel), nil
}

func (m *Membership) GetPostgresqlMembershipsByCommunityId(
	communityId uuid.UUID,
) ([]*schemas.Membership, *errors.Error) {
	membershipsModel, err := m.DaoPostgresql.Membership.GetMembershipsByCommunityId(communityId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	memberships := make([]*schemas.Membership, len(membershipsModel))
	for i, membershipModel := range membershipsModel {
		memberships[i] = m.convertModelToSchema(membershipModel)
	}

	return memberships, nil
}

func (m *Membership) FetchPostgresqlMemberships() ([]*schemas.Membership, *errors.Error) {
	membershipsModel, err := m.DaoPostgresql.Membership.FetchMemberships()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	memberships := make([]*schemas.Membership, len(membershipsModel))
	for i, membershipModel := range membershipsModel {
		memberships[i] = m.convertModelToSchema(membershipModel)
	}

	return memberships, nil
}

func (m *Membership) CreatePostgresqlMembership(
	description string,
	startDate time.Time,
	endDate time.Time,
	status schemas.MembershipStatus,
	communityId uuid.UUID,
	userId uuid.UUID,
	planId uuid.UUID,
	updatedBy string,
) (*schemas.Membership, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	membershipModel := &model.Membership{
		Id:          uuid.New(),
		Description: description,
		StartDate:   startDate,
		EndDate:     endDate,
		Status:      model.MembershipStatus(status),
		CommunityId: communityId,
		UserId:      userId,
		PlanId:      planId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := m.DaoPostgresql.Membership.CreateMembership(membershipModel); err != nil {
		return nil, &errors.BadRequestError.MembershipNotCreated
	}

	// Obtener la membership creada con las relaciones
	createdMembership, err := m.DaoPostgresql.Membership.GetMembership(membershipModel.Id)
	if err != nil {
		return nil, &errors.BadRequestError.MembershipNotCreated
	}

	return m.convertModelToSchema(createdMembership), nil
}

func (m *Membership) UpdatePostgresqlMembership(
	membershipId uuid.UUID,
	description *string,
	startDate *time.Time,
	endDate *time.Time,
	status *schemas.MembershipStatus,
	communityId *uuid.UUID,
	userId *uuid.UUID,
	planId *uuid.UUID,
	updatedBy string,
) (*schemas.Membership, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Obtener la membership existente
	existingMembership, err := m.DaoPostgresql.Membership.GetMembership(membershipId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.MembershipNotFound
	}

	// Aplicar updates solo a campos no nil
	if description != nil {
		existingMembership.Description = *description
	}
	if startDate != nil {
		existingMembership.StartDate = *startDate
	}
	if endDate != nil {
		existingMembership.EndDate = *endDate
	}
	if status != nil {
		existingMembership.Status = model.MembershipStatus(*status)
	}
	if communityId != nil {
		existingMembership.CommunityId = *communityId
	}
	if userId != nil {
		existingMembership.UserId = *userId
	}
	if planId != nil {
		existingMembership.PlanId = *planId
	}

	existingMembership.AuditFields.UpdatedBy = updatedBy

	if err := m.DaoPostgresql.Membership.UpdateMembership(existingMembership); err != nil {
		return nil, &errors.BadRequestError.MembershipNotUpdated
	}

	// Obtener la membership actualizada con las relaciones
	updatedMembership, err := m.DaoPostgresql.Membership.GetMembership(membershipId)
	if err != nil {
		return nil, &errors.BadRequestError.MembershipNotUpdated
	}

	return m.convertModelToSchema(updatedMembership), nil
}

func (m *Membership) DeletePostgresqlMembership(membershipId uuid.UUID) *errors.Error {
	if err := m.DaoPostgresql.Membership.DeleteMembership(membershipId); err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.MembershipNotFound
		}
		return &errors.BadRequestError.MembershipNotDeleted
	}
	return nil
}

// Función helper para convertir model a schema
func (m *Membership) convertModelToSchema(membershipModel *model.Membership) *schemas.Membership {
	return &schemas.Membership{
		Id:          membershipModel.Id,
		Description: membershipModel.Description,
		StartDate:   membershipModel.StartDate,
		EndDate:     membershipModel.EndDate,
		Status:      schemas.MembershipStatus(membershipModel.Status),
		CommunityId: membershipModel.CommunityId,
		Community: schemas.Community{
			Id:                  membershipModel.Community.Id,
			Name:                membershipModel.Community.Name,
			Purpose:             membershipModel.Community.Purpose,
			ImageUrl:            membershipModel.Community.ImageUrl,
			NumberSubscriptions: membershipModel.Community.NumberSubscriptions,
		},
		UserId: membershipModel.UserId,
		User: schemas.User{
			Id:             membershipModel.User.Id,
			Name:           membershipModel.User.Name,
			FirstLastName:  membershipModel.User.FirstLastName,
			SecondLastName: membershipModel.User.SecondLastName,
			Email:          membershipModel.User.Email,
			Rol:            schemas.UserRol(membershipModel.User.Rol),
			ImageUrl:       membershipModel.User.ImageUrl,
		},
		PlanId: membershipModel.PlanId,
		Plan: schemas.Plan{
			Id:               membershipModel.Plan.Id,
			Fee:              membershipModel.Plan.Fee,
			Type:             membershipModel.Plan.Type,
			ReservationLimit: membershipModel.Plan.ReservationLimit,
		},
	}
}
