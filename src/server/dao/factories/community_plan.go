package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type CommunityPlanModelF struct {
	Id          *uuid.UUID
	CommunityId *uuid.UUID
	PlanId      *uuid.UUID
}

// Create a new community plan on DB
func NewCommunityPlanModel(db *gorm.DB, option ...CommunityPlanModelF) *model.CommunityPlan {
	// Create default community if not provided
	community := NewCommunityModel(db)

	// Create default plan if not provided
	plan := NewPlanModel(db)

	communityPlan := &model.CommunityPlan{
		Id:          uuid.New(),
		CommunityId: community.Id,
		PlanId:      plan.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			communityPlan.Id = *parameters.Id
		}
		if parameters.CommunityId != nil {
			communityPlan.CommunityId = *parameters.CommunityId
		}
		if parameters.PlanId != nil {
			communityPlan.PlanId = *parameters.PlanId
		}
	}

	result := db.Create(communityPlan)
	if result.Error != nil {
		log.Fatalf("Error when trying to create community plan: %v", result.Error)
	}

	return communityPlan
}

// Create size number of new community plans on DB
func NewCommunityPlanModelBatch(
	db *gorm.DB,
	size int,
	option ...CommunityPlanModelF,
) []*model.CommunityPlan {
	communityPlans := []*model.CommunityPlan{}
	for i := 0; i < size; i++ {
		var communityPlan *model.CommunityPlan
		if len(option) > 0 {
			communityPlan = NewCommunityPlanModel(db, option[0])
		} else {
			communityPlan = NewCommunityPlanModel(db)
		}
		communityPlans = append(communityPlans, communityPlan)
	}
	return communityPlans
}
