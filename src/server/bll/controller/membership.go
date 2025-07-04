package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Membership struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

func NewMembershipController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Membership {
	return &Membership{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

func (m *Membership) GetMembership(membershipId uuid.UUID) (*schemas.Membership, *errors.Error) {
	return m.Adapter.Membership.GetPostgresqlMembership(membershipId)
}

func (m *Membership) GetMembershipsByUserId(userId uuid.UUID) (*schemas.Memberships, *errors.Error) {
	memberships, err := m.Adapter.Membership.GetPostgresqlMembershipsByUserId(userId)
	if err != nil {
		return nil, err
	}

	return &schemas.Memberships{Memberships: memberships}, nil
}

func (m *Membership) GetMembershipsByCommunityId(communityId uuid.UUID) (*schemas.Memberships, *errors.Error) {
	memberships, err := m.Adapter.Membership.GetPostgresqlMembershipsByCommunityId(communityId)
	if err != nil {
		return nil, err
	}

	return &schemas.Memberships{Memberships: memberships}, nil
}

func (m *Membership) GetMembershipByUserAndCommunity(userId uuid.UUID, communityId uuid.UUID) (*schemas.Membership, *errors.Error) {
	return m.Adapter.Membership.GetPostgresqlMembershipByUserAndCommunity(userId, communityId)
}

func (m *Membership) FetchMemberships() (*schemas.Memberships, *errors.Error) {
	memberships, err := m.Adapter.Membership.FetchPostgresqlMemberships()
	if err != nil {
		return nil, err
	}

	return &schemas.Memberships{Memberships: memberships}, nil
}

func (m *Membership) CreateMembership(
	createMembershipRequest schemas.CreateMembershipRequest,
	updatedBy string,
) (*schemas.Membership, *errors.Error) {
	// Validar que el usuario existe
	_, userErr := m.Adapter.User.GetPostgresqlUser(createMembershipRequest.UserId)
	if userErr != nil {
		return nil, userErr
	}

	// Validar que la comunidad existe
	_, communityErr := m.Adapter.Community.GetPostgresqlCommunity(createMembershipRequest.CommunityId)
	if communityErr != nil {
		return nil, communityErr
	}

	// Validar que el plan existe
	_, planErr := m.Adapter.Plan.GetPostgresqlPlan(createMembershipRequest.PlanId)
	if planErr != nil {
		return nil, planErr
	}

	// Validar que el plan pertenece a la comunidad (verificar en CommunityPlan)
	_, communityPlanErr := m.Adapter.CommunityPlan.GetPostgresqlCommunityPlan(createMembershipRequest.CommunityId, createMembershipRequest.PlanId)
	if communityPlanErr != nil {
		return nil, communityPlanErr
	}

	return m.Adapter.Membership.CreatePostgresqlMembership(
		createMembershipRequest.Description,
		createMembershipRequest.StartDate,
		createMembershipRequest.EndDate,
		createMembershipRequest.Status,
		createMembershipRequest.CommunityId,
		createMembershipRequest.UserId,
		createMembershipRequest.PlanId,
		updatedBy,
	)
}

// CreateMembershipForUser crea una membership ligada a un usuario específico
func (m *Membership) CreateMembershipForUser(
	userId uuid.UUID,
	createMembershipForUserRequest schemas.CreateMembershipForUserRequest,
	updatedBy string,
) (*schemas.Membership, *errors.Error) {
	// Validar que el usuario existe antes de crear la membership
	_, userErr := m.Adapter.User.GetPostgresqlUser(userId)
	if userErr != nil {
		return nil, userErr
	}

	// Validar que la comunidad existe
	_, communityErr := m.Adapter.Community.GetPostgresqlCommunity(createMembershipForUserRequest.CommunityId)
	if communityErr != nil {
		return nil, communityErr
	}

	// Validar que el plan existe
	_, planErr := m.Adapter.Plan.GetPostgresqlPlan(createMembershipForUserRequest.PlanId)
	if planErr != nil {
		return nil, planErr
	}

	// Validar que el plan pertenece a la comunidad (verificar en CommunityPlan)
	_, communityPlanErr := m.Adapter.CommunityPlan.GetPostgresqlCommunityPlan(createMembershipForUserRequest.CommunityId, createMembershipForUserRequest.PlanId)
	if communityPlanErr != nil {
		return nil, communityPlanErr
	}

	return m.Adapter.Membership.CreatePostgresqlMembership(
		createMembershipForUserRequest.Description,
		createMembershipForUserRequest.StartDate,
		createMembershipForUserRequest.EndDate,
		createMembershipForUserRequest.Status,
		createMembershipForUserRequest.CommunityId,
		userId, // El userId viene del parámetro de la URL, no del body
		createMembershipForUserRequest.PlanId,
		updatedBy,
	)
}

func (m *Membership) UpdateMembership(
	membershipId uuid.UUID,
	updateMembershipRequest schemas.UpdateMembershipRequest,
	updatedBy string,
) (*schemas.Membership, *errors.Error) {
	// Validar que la membership existe antes de actualizar
	_, membershipErr := m.Adapter.Membership.GetPostgresqlMembership(membershipId)
	if membershipErr != nil {
		return nil, membershipErr
	}

	// Si se está actualizando el usuario, validar que existe
	if updateMembershipRequest.UserId != nil {
		_, userErr := m.Adapter.User.GetPostgresqlUser(*updateMembershipRequest.UserId)
		if userErr != nil {
			return nil, userErr
		}
	}

	// Si se está actualizando la comunidad, validar que existe
	if updateMembershipRequest.CommunityId != nil {
		_, communityErr := m.Adapter.Community.GetPostgresqlCommunity(*updateMembershipRequest.CommunityId)
		if communityErr != nil {
			return nil, communityErr
		}
	}

	// Si se está actualizando el plan, validar que existe
	if updateMembershipRequest.PlanId != nil {
		_, planErr := m.Adapter.Plan.GetPostgresqlPlan(*updateMembershipRequest.PlanId)
		if planErr != nil {
			return nil, planErr
		}
	}

	return m.Adapter.Membership.UpdatePostgresqlMembership(
		membershipId,
		updateMembershipRequest.Description,
		updateMembershipRequest.StartDate,
		updateMembershipRequest.EndDate,
		updateMembershipRequest.Status,
		updateMembershipRequest.CommunityId,
		updateMembershipRequest.UserId,
		updateMembershipRequest.PlanId,
		updatedBy,
	)
}

// GetUsersByCommunityId retrieves all users who have active memberships in the specified community
func (m *Membership) GetUsersByCommunityId(communityId uuid.UUID) (*schemas.Users, *errors.Error) {
	// First, check if the community exists
	_, communityErr := m.Adapter.Community.GetPostgresqlCommunity(communityId)
	if communityErr != nil {
		return nil, communityErr
	}

	// Get all memberships for the community
	memberships, err := m.Adapter.Membership.GetPostgresqlMembershipsByCommunityId(communityId)
	if err != nil {
		return nil, err
	}

	// Early return if no memberships are found
	if len(memberships) == 0 {
		return &schemas.Users{Users: make([]*schemas.User, 0)}, nil
	}

	// Extract user IDs from active memberships
	var userIds []uuid.UUID
	for _, membership := range memberships {
		// Only include active memberships
		if membership.Status == "ACTIVE" {
			userIds = append(userIds, membership.UserId)
		}
	}

	// Get all users with those IDs
	users, err := m.Adapter.User.GetPostgresqlUsersByIds(userIds)
	if err != nil {
		return nil, err
	}

	return &schemas.Users{Users: users}, nil
}

func (m *Membership) DeleteMembership(membershipId uuid.UUID) *errors.Error {
	return m.Adapter.Membership.DeletePostgresqlMembership(membershipId)
}
