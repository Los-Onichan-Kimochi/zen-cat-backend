package controller

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"onichankimochi.com/astro_cat_backend/src/logging"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
)

type Onboarding struct {
	logger       logging.Logger
	PostgresqlDB *gorm.DB
}

func NewOnboardingController(logger logging.Logger, postgresqlDB *gorm.DB) *Onboarding {
	return &Onboarding{
		logger:       logger,
		PostgresqlDB: postgresqlDB,
	}
}

func (o *Onboarding) GetOnboarding(onboardingId uuid.UUID) (*model.Onboarding, error) {
	onboarding := &model.Onboarding{}
	result := o.PostgresqlDB.
		Preload("User").
		First(&onboarding, "id = ?", onboardingId)
	if result.Error != nil {
		return nil, result.Error
	}

	return onboarding, nil
}

func (o *Onboarding) GetOnboardingByUserId(userId uuid.UUID) (*model.Onboarding, error) {
	onboarding := &model.Onboarding{}
	result := o.PostgresqlDB.
		Preload("User").
		First(&onboarding, "user_id = ?", userId)
	if result.Error != nil {
		return nil, result.Error
	}

	return onboarding, nil
}

func (o *Onboarding) FetchOnboardings() ([]*model.Onboarding, error) {
	onboardings := []*model.Onboarding{}
	result := o.PostgresqlDB.
		Preload("User").
		Find(&onboardings)
	if result.Error != nil {
		return nil, result.Error
	}

	return onboardings, nil
}

func (o *Onboarding) CreateOnboarding(onboarding *model.Onboarding) error {
	return o.PostgresqlDB.Create(onboarding).Error
}

func (o *Onboarding) UpdateOnboarding(
	id uuid.UUID,
	documentType *model.DocumentType,
	documentNumber *string,
	phoneNumber *string,
	birthDate *time.Time,
	gender *model.Gender,
	city *string,
	postalCode *string,
	district *string,
	address *string,
	updatedBy string,
) (*model.Onboarding, error) {
	updateFields := map[string]any{
		"updated_by": updatedBy,
	}

	if documentType != nil {
		updateFields["document_type"] = *documentType
	}
	if documentNumber != nil {
		updateFields["document_number"] = *documentNumber
	}
	if phoneNumber != nil {
		updateFields["phone_number"] = *phoneNumber
	}
	if birthDate != nil {
		updateFields["birth_date"] = *birthDate
	}
	if gender != nil {
		updateFields["gender"] = *gender
	}
	if city != nil {
		updateFields["city"] = *city
	}
	if postalCode != nil {
		updateFields["postal_code"] = *postalCode
	}
	if district != nil {
		updateFields["district"] = *district
	}
	if address != nil {
		updateFields["address"] = *address
	}

	var onboarding model.Onboarding
	if len(updateFields) == 1 {
		if err := o.PostgresqlDB.First(&onboarding, "id = ?", id).Error; err != nil {
			return nil, err
		}
		return &onboarding, nil
	}

	result := o.PostgresqlDB.Model(&onboarding).
		Clauses(clause.Returning{}).
		Where("id = ?", id).
		Updates(updateFields)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &onboarding, nil
}

func (o *Onboarding) DeleteOnboarding(onboardingId uuid.UUID) error {
	result := o.PostgresqlDB.Delete(&model.Onboarding{}, "id = ?", onboardingId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (o *Onboarding) DeleteOnboardingByUserId(userId uuid.UUID) error {
	result := o.PostgresqlDB.Delete(&model.Onboarding{}, "user_id = ?", userId)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
