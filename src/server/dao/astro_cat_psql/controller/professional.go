package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Professional struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Professional postgresql controller
func NewProfessionalController(logger logging.Logger, postgresqlDB *gorm.DB) *Professional {
	return &Professional{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a professional model given params.
func (p *Professional) GetProfessional(professionalId uuid.UUID) (*model.Professional, error) {
	professional := &model.Professional{}

	result := p.PostgresqlDB.First(&professional, "id = ?", professionalId)
	if result.Error != nil {
		return nil, result.Error
	}

	return professional, nil
}

// Fetch all professionals.
func (p *Professional) FetchProfessionals() ([]*model.Professional, error) {
	professionals := []*model.Professional{}

	result := p.PostgresqlDB.Find(&professionals)
	if result.Error != nil {
		return nil, result.Error
	}

	return professionals, nil
}

// Creates a professional given its model.
func (p *Professional) CreateProfessional(professional *model.Professional) error {
	return p.PostgresqlDB.Create(professional).Error
}

// Updates professional given fields to update.
func (p *Professional) UpdateProfessional(
	id uuid.UUID,
	name *string,
	firstLastName *string,
	secondLastName *string,
	specialty *string,
	email *string,
	phoneNumber *string,
	professionalType *string,
	imageUrl *string,
	updatedBy string,
) (*model.Professional, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}
	if name != nil {
		updateFields["name"] = *name
	}
	if firstLastName != nil {
		updateFields["first_last_name"] = *firstLastName
	}
	if secondLastName != nil {
		updateFields["second_last_name"] = *secondLastName
	}
	if specialty != nil {
		updateFields["specialty"] = *specialty
	}
	if email != nil {
		updateFields["email"] = *email
	}
	if phoneNumber != nil {
		updateFields["phone_number"] = *phoneNumber
	}
	if professionalType != nil {
		updateFields["type"] = *professionalType
	}
	if imageUrl != nil {
		updateFields["image_url"] = *imageUrl
	}
	// Check if there are any fields to update
	var professional model.Professional
	if len(updateFields) == 1 {
		if err := p.PostgresqlDB.First(&professional, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &professional, nil
	}

	// Perform the update
	result := p.PostgresqlDB.Model(&professional).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &professional, nil
}

// Soft deletes a professional given its ID.
func (p *Professional) DeleteProfessional(professionalId uuid.UUID) error {
	result := p.PostgresqlDB.Delete(&model.Professional{}, "id = ?", professionalId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Creates multiple professionals in a batch.
func (p *Professional) BulkCreateProfessionals(professionals []*model.Professional) error {
	if len(professionals) == 0 {
		return nil
	}
	return p.PostgresqlDB.Create(&professionals).Error
}

// Batch deletes multiple professionals given their IDs.
func (p *Professional) BulkDeleteProfessionals(professionalIds []uuid.UUID) error {
	if len(professionalIds) == 0 {
		return nil
	}

	result := p.PostgresqlDB.Where("id IN ?", professionalIds).Delete(&model.Professional{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
