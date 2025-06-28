package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type CommunityServiceModelF struct {
	Id          *uuid.UUID
	CommunityId *uuid.UUID
	ServiceId   *uuid.UUID
}

// Create a new community service on DB
func NewCommunityServiceModel(db *gorm.DB, option ...CommunityServiceModelF) *model.CommunityService {
	var communityService *model.CommunityService

	// Use a transaction to ensure atomicity
	err := db.Transaction(func(tx *gorm.DB) error {
		// Create default community if not provided
		community := NewCommunityModel(tx)

		// Create default service if not provided
		service := NewServiceModel(tx)

		communityService = &model.CommunityService{
			Id:          uuid.New(),
			CommunityId: community.Id,
			ServiceId:   service.Id,
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		}

		if len(option) > 0 {
			parameters := option[0]
			if parameters.Id != nil {
				communityService.Id = *parameters.Id
			}
			if parameters.CommunityId != nil {
				communityService.CommunityId = *parameters.CommunityId
			}
			if parameters.ServiceId != nil {
				communityService.ServiceId = *parameters.ServiceId
			}
		}

		result := tx.Create(communityService)
		if result.Error != nil {
			return result.Error
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Error when trying to create community service: %v", err)
	}

	return communityService
}

// Create size number of new community services on DB
func NewCommunityServiceModelBatch(
	db *gorm.DB,
	size int,
	option ...CommunityServiceModelF,
) []*model.CommunityService {
	communityServices := []*model.CommunityService{}
	for i := 0; i < size; i++ {
		var communityService *model.CommunityService
		if len(option) > 0 {
			communityService = NewCommunityServiceModel(db, option[0])
		} else {
			communityService = NewCommunityServiceModel(db)
		}
		communityServices = append(communityServices, communityService)
	}
	return communityServices
}
