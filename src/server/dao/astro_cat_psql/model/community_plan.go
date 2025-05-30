package model

import "github.com/google/uuid"

type CommunityPlan struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CommunityId uuid.UUID `gorm:"type:uuid"`
	Community   Community `gorm:"foreignKey:CommunityId;references:Id"`
	PlanId      uuid.UUID `gorm:"type:uuid"`
	Plan        Plan      `gorm:"foreignKey:PlanId;references:Id"`
	AuditFields
}

func (CommunityPlan) TableName() string {
	return "astro_cat_community_plan"
}
