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

// Create Membership postgresql controller
func NewMembershipController(logger logging.Logger, postgresqlDB *gorm.DB) *Membership {
	return &Membership{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a specific reservation by ID.
func (m *Membership) GetMembership(membershipId uuid.UUID) (*model.Membership, error) {
	var membership model.Membership
	result := m.PostgresqlDB.Preload("User").
		Preload("Community").
		Preload("Plan").
		Where("id = ?", membershipId).
		First(&membership)
	if result.Error != nil {
		return nil, result.Error
	}

	return &membership, nil
}
