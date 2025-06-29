package factories

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type MembershipModelF struct {
	Id          *uuid.UUID
	Description *string
	StartDate   *time.Time
	EndDate     *time.Time
	Status      *model.MembershipStatus
	CommunityId *uuid.UUID
	UserId      *uuid.UUID
	PlanId      *uuid.UUID
}

// Create a new membership on DB
func NewMembershipModel(db *gorm.DB, option ...MembershipModelF) *model.Membership {
	var membership *model.Membership

	// Use a transaction to ensure atomicity
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create default related entities if not provided
		community := NewCommunityModel(tx)
		user := NewUserModel(tx)
		plan := NewPlanModel(tx)

		now := time.Now()
		endDate := now.AddDate(0, 1, 0) // 1 month from now

		membership = &model.Membership{
			Id:          uuid.New(),
			Description: "Test Membership",
			StartDate:   now,
			EndDate:     endDate,
			Status:      model.MembershipStatusActive,
			CommunityId: community.Id,
			UserId:      user.Id,
			PlanId:      plan.Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		}

		if len(option) > 0 {
			parameters := option[0]
			if parameters.Id != nil {
				membership.Id = *parameters.Id
			}
			if parameters.Description != nil {
				membership.Description = *parameters.Description
			}
			if parameters.StartDate != nil {
				membership.StartDate = *parameters.StartDate
			}
			if parameters.EndDate != nil {
				membership.EndDate = *parameters.EndDate
			}
			if parameters.Status != nil {
				membership.Status = *parameters.Status
			}
			if parameters.CommunityId != nil {
				membership.CommunityId = *parameters.CommunityId
			}
			if parameters.UserId != nil {
				membership.UserId = *parameters.UserId
			}
			if parameters.PlanId != nil {
				membership.PlanId = *parameters.PlanId
			}
		}

		result := tx.Create(membership)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error when trying to create membership: %v", err)
	}

	return membership
}

// Create size number of new memberships on DB
func NewMembershipModelBatch(
	db *gorm.DB,
	size int,
	option ...MembershipModelF,
) []*model.Membership {
	memberships := []*model.Membership{}
	for i := 0; i < size; i++ {
		var membership *model.Membership
		if len(option) > 0 {
			membership = NewMembershipModel(db, option[0])
		} else {
			membership = NewMembershipModel(db)
		}
		memberships = append(memberships, membership)
	}
	return memberships
}
