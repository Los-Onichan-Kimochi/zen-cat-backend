package controller

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type MembershipSuspension struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewMembershipSuspensionController(
	logger logging.Logger,
	postgresqlDB *gorm.DB,
) *MembershipSuspension {
	return &MembershipSuspension{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (m *MembershipSuspension) CreateMembershipSuspension(membershipId uuid.UUID) (*model.MembershipSuspension, error) {
	suspension := &model.MembershipSuspension{
		Id:           uuid.New(),
		MembershipId: membershipId,
		SuspendedAt:  time.Now(),
		ResumedAt:    nil, // Explicitly nil
	}

	result := m.PostgresqlDB.Create(suspension)
	if result.Error != nil {
		m.logger.Errorf("failed to create membership suspension: %v", result.Error)
		return nil, result.Error
	}

	return suspension, nil
}

func (m *MembershipSuspension) GetLatestOpenMembershipSuspension(membershipId uuid.UUID) (*model.MembershipSuspension, error) {
	var suspension model.MembershipSuspension
	result := m.PostgresqlDB.Where("membership_id = ? AND resumed_at IS NULL", membershipId).Order("suspended_at desc").First(&suspension)

	if result.Error != nil {
		return nil, result.Error
	}

	return &suspension, nil
}

func (m *MembershipSuspension) UpdateMembershipSuspension(suspension *model.MembershipSuspension) error {
	result := m.PostgresqlDB.Save(suspension)
	if result.Error != nil {
		m.logger.Errorf("failed to update membership suspension: %v", result.Error)
		return result.Error
	}

	return nil
}
