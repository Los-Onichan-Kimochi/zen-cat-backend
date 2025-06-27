package factories

import (
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type LocalModelF struct {
	Id             *uuid.UUID
	LocalName      *string
	StreetName     *string
	BuildingNumber *string
	District       *string
	Province       *string
	Region         *string
	Reference      *string
	Capacity       *int
	ImageUrl       *string
}

// Create a new local on DB
func NewLocalModel(db *gorm.DB, option ...LocalModelF) *model.Local {
	capacity := 20
	local := &model.Local{
		Id:             uuid.New(),
		LocalName:      "Test Local",
		StreetName:     "Test Street",
		BuildingNumber: "123",
		District:       "Test District",
		Province:       "Test Province",
		Region:         "Test Region",
		Reference:      "Near Test Landmark",
		Capacity:       capacity,
		ImageUrl:       "https://example.com/local.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}

	if len(option) > 0 {
		parameters := option[0]
		if parameters.Id != nil {
			local.Id = *parameters.Id
		}
		if parameters.LocalName != nil {
			local.LocalName = *parameters.LocalName
		}
		if parameters.StreetName != nil {
			local.StreetName = *parameters.StreetName
		}
		if parameters.BuildingNumber != nil {
			local.BuildingNumber = *parameters.BuildingNumber
		}
		if parameters.District != nil {
			local.District = *parameters.District
		}
		if parameters.Province != nil {
			local.Province = *parameters.Province
		}
		if parameters.Region != nil {
			local.Region = *parameters.Region
		}
		if parameters.Reference != nil {
			local.Reference = *parameters.Reference
		}
		if parameters.Capacity != nil {
			local.Capacity = *parameters.Capacity
		}
		if parameters.ImageUrl != nil {
			local.ImageUrl = *parameters.ImageUrl
		}
	}

	result := db.Create(local)
	if result.Error != nil {
		log.Fatalf("Error when trying to create local: %v", result.Error)
	}

	return local
}

// Create size number of new locals on DB
func NewLocalModelBatch(
	db *gorm.DB,
	size int,
	option ...LocalModelF,
) []*model.Local {
	locals := []*model.Local{}
	for i := 0; i < size; i++ {
		var local *model.Local
		if len(option) > 0 {
			local = NewLocalModel(db, option[0])
		} else {
			local = NewLocalModel(db)
		}
		locals = append(locals, local)
	}
	return locals
}
