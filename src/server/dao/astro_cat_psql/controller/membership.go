package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Membership struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewMembershipController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *Membership {
	return &Membership{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (m *Membership) GetMembership(membershipId uuid.UUID) (*model.Membership, error) {
	var membership model.Membership
	result := m.PostgresqlDB.Preload("Community").Preload("User").Preload("Plan").
		Where("id = ?", membershipId).First(&membership)

	if result.Error != nil {
		return nil, result.Error
	}

	return &membership, nil
}

func (m *Membership) GetMembershipsByUserId(userId uuid.UUID) ([]*model.Membership, error) {
	var memberships []*model.Membership
	result := m.PostgresqlDB.Preload("Community").Preload("User").Preload("Plan").
		Where("user_id = ?", userId).Find(&memberships)

	if result.Error != nil {
		return nil, result.Error
	}

	return memberships, nil
}

func (m *Membership) GetMembershipsByCommunityId(communityId uuid.UUID) ([]*model.Membership, error) {
	var memberships []*model.Membership
	result := m.PostgresqlDB.Preload("Community").Preload("User").Preload("Plan").
		Where("community_id = ?", communityId).Find(&memberships)

	if result.Error != nil {
		return nil, result.Error
	}

	return memberships, nil
}

func (m *Membership) FetchMemberships() ([]*model.Membership, error) {
	var memberships []*model.Membership
	result := m.PostgresqlDB.Preload("Community").Preload("User").Preload("Plan").
		Find(&memberships)

	if result.Error != nil {
		return nil, result.Error
	}

	return memberships, nil
}

func (m *Membership) CreateMembership(membership *model.Membership) error {
	result := m.PostgresqlDB.Create(membership)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *Membership) UpdateMembership(membership *model.Membership) error {
	result := m.PostgresqlDB.Save(membership)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (m *Membership) DeleteMembership(membershipId uuid.UUID) error {
	result := m.PostgresqlDB.Where("id = ?", membershipId).Delete(&model.Membership{})
	if result.Error != nil {
		return result.Error
	}

	// Check if any rows were affected
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
