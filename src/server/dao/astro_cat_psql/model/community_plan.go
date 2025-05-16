package model

import "github.com/google/uuid"

type CommunityPlan struct {
	CommunityId uuid.UUID `gorm:"primaryKey"`
	Community   Community `gorm:"foreignKey:CommunityId;references:Id"`
	PlanId      uuid.UUID `gorm:"primaryKey"`
	Plan        Plan      `gorm:"foreignKey:PlanId;references:Id"`
	AuditFields
}

func (CommunityPlan) TableName() string {
	return "astro_cat_community_plan"
}
