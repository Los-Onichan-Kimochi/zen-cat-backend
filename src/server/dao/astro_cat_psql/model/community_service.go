package model

import "github.com/google/uuid"

type CommunityService struct {
	CommunityId uuid.UUID `gorm:"primaryKey"`
	Community   Community `gorm:"foreignKey:CommunityId;references:Id"`
	ServiceId   uuid.UUID `gorm:"primaryKey"`
	Service     Service   `gorm:"foreignKey:ServiceId;references:Id"`
	AuditFields
}

func (CommunityService) TableName() string {
	return "astro_cat_community_service"
}
