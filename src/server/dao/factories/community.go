package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type CommunityModelF struct {
	Id                  *uuid.UUID
	Name                *string
	Purpose             *string
	ImageUrl            *string
	NumberSubscriptions *int
}

// Create a new community on DB
func NewCommunityModel(db *gorm.DB, option ...CommunityModelF) *model.Community {
	community := &model.Community{
		Id:                  uuid.New(),
		Name:                "Test Community",
		Purpose:             "Test Purpose",
		ImageUrl:            "https://example.com/community.jpg",
		NumberSubscriptions: 0,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			community.Id = *parameters.Id
		}
		if parameters.Name != nil {
			community.Name = *parameters.Name
		}
		if parameters.Purpose != nil {
			community.Purpose = *parameters.Purpose
		}
		if parameters.ImageUrl != nil {
			community.ImageUrl = *parameters.ImageUrl
		}
		if parameters.NumberSubscriptions != nil {
			community.NumberSubscriptions = *parameters.NumberSubscriptions
		}
	}

	result := db.Create(community)
	if result.Error != nil {
		log.Fatalf("Error when trying to create community: %v", result.Error)
	}

	return community
}

// Create size number of new communities on DB
func NewCommunityModelBatch(
	db *gorm.DB,
	size int,
	option ...CommunityModelF,
) []*model.Community {
	communities := []*model.Community{}
	for i := 0; i < size; i++ {
		var community *model.Community
		if len(option) > 0 {
			community = NewCommunityModel(db, option[0])
		} else {
			community = NewCommunityModel(db)
		}
		communities = append(communities, community)
	}
	return communities
}
