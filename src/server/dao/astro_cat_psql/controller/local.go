package controller

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Local struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

// Create Local postgresql controller
func NewLocalController(logger logging.Logger, postgresqlDB *gorm.DB) *Local {
	return &Local{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

// Gets a local model given params.
func (l *Local) GetLocal(localId uuid.UUID) (*model.Local, error) {
	local := &model.Local{}

	result := l.PostgresqlDB.First(&local, "id = ?", localId)
	if result.Error != nil {
		return nil, result.Error
	}

	return local, nil
}

// Fetch all locals.
func (l *Local) FetchLocals() ([]*model.Local, error) {
	locals := []*model.Local{}

	result := l.PostgresqlDB.Find(&locals)
	if result.Error != nil {
		return nil, result.Error
	}

	return locals, nil
}

// Creates a professional given its model.
func (l *Local) CreateLocal(local *model.Local) error {
	return l.PostgresqlDB.Create(local).Error
}

// Updates a local given its model.
func (l *Local) UpdateLocal(
	id uuid.UUID,
	localName *string,
	streetName *string,
	buildingNumber *string,
	district *string,
	province *string,
	region *string,
	reference *string,
	capacity *int,
	imageUrl *string,
	updatedBy string,
) (*model.Local, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}
	if localName != nil {
		updateFields["local_name"] = *localName
	}
	if streetName != nil {
		updateFields["street_name"] = *streetName
	}
	if buildingNumber != nil {
		updateFields["building_number"] = *buildingNumber
	}
	if district != nil {
		updateFields["district"] = *district
	}
	if province != nil {
		updateFields["province"] = *province
	}
	if region != nil {
		updateFields["region"] = *region
	}
	if reference != nil {
		updateFields["reference"] = *reference
	}
	if capacity != nil {
		updateFields["capacity"] = *capacity
	}
	if imageUrl != nil {
		updateFields["image_url"] = *imageUrl
	}

	// Check if there are any fields to update
	var local model.Local
	if len(updateFields) == 1 {
		if err := l.PostgresqlDB.First(&local, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &local, nil
	}

	// Perform the update
	result := l.PostgresqlDB.Model(&local).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	return &local, nil
}

// Soft deletes a local given its ID.
func (l *Local) DeleteLocal(id uuid.UUID) error {
	result := l.PostgresqlDB.Delete(&model.Local{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

//Creates multiple locals in a batch.
func (l *Local) BulkCreateLocals(locals []*model.Local) error {
	return l.PostgresqlDB.Create(&locals).Error
}

// Batch deletes multiple professionals given their IDs.
func (l *Local) BulkDeleteLocals(localIds []uuid.UUID) error {
	if len(localIds) == 0 {
		return nil
	}

	result := l.PostgresqlDB.Where("id IN ?", localIds).Delete(&model.Local{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}
