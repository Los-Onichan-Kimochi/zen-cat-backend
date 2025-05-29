package model

import "github.com/google/uuid"

type CommunityService struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	CommunityId uuid.UUID `gorm:"type:uuid"`
	Community   Community `gorm:"foreignKey:CommunityId;references:Id"`
	ServiceId   uuid.UUID `gorm:"type:uuid"`
	Service     Service   `gorm:"foreignKey:ServiceId;references:Id"`
	AuditFields
}

func (CommunityService) TableName() string {
	return "astro_cat_community_service"
}
