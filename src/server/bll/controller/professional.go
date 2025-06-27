package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Professional struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Professional controller
func NewProfessionalController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Professional {
	return &Professional{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a professional.
func (p *Professional) GetProfessional(
	professionalId uuid.UUID,
) (*schemas.Professional, *errors.Error) {
	return p.Adapter.Professional.GetPostgresqlProfessional(professionalId)
}

// Fetch all professionals.
func (p *Professional) FetchProfessionals() (*schemas.Professionals, *errors.Error) {
	professionals, err := p.Adapter.Professional.FetchPostgresqlProfessionals()
	if err != nil {
		return nil, err
	}

	return &schemas.Professionals{Professionals: professionals}, nil
}

// Creates a professional.
func (p *Professional) CreateProfessional(
	createProfessionalData schemas.CreateProfessionalRequest,
	updatedBy string,
) (*schemas.Professional, *errors.Error) {
	// Validate required fields
	if createProfessionalData.Name == "" {
		return nil, &errors.BadRequestError.ProfessionalNotCreated
	}

	var secondLastName *string
	if createProfessionalData.SecondLastName != "" {
		secondLastName = &createProfessionalData.SecondLastName
	}
	return p.Adapter.Professional.CreatePostgresqlProfessional(
		createProfessionalData.Name,
		createProfessionalData.FirstLastName,
		secondLastName,
		createProfessionalData.Specialty,
		createProfessionalData.Email,
		createProfessionalData.PhoneNumber,
		createProfessionalData.Type,
		createProfessionalData.ImageUrl,
		updatedBy,
	)
}

// Updates a professional.
func (p *Professional) UpdateProfessional(
	professionalId uuid.UUID,
	updateProfessionalData schemas.UpdateProfessionalRequest,
	updatedBy string,
) (*schemas.Professional, *errors.Error) {
	return p.Adapter.Professional.UpdatePostgresqlProfessional(
		professionalId,
		updateProfessionalData.Name,
		updateProfessionalData.FirstLastName,
		updateProfessionalData.SecondLastName,
		updateProfessionalData.Specialty,
		updateProfessionalData.Email,
		updateProfessionalData.PhoneNumber,
		updateProfessionalData.Type,
		updateProfessionalData.ImageUrl,
		updatedBy,
	)
}

func (p *Professional) DeleteProfessional(professionalId uuid.UUID) *errors.Error {
	return p.Adapter.Professional.DeletePostgresqlProfessional(professionalId)
}

// Bulk creates professionals.
func (p *Professional) BulkCreateProfessionals(
	createProfessionalsData []*schemas.CreateProfessionalRequest,
	updatedBy string,
) (*schemas.Professionals, *errors.Error) {
	professionals, err := p.Adapter.Professional.BulkCreatePostgresqlProfessionals(
		createProfessionalsData,
		updatedBy,
	)
	if err != nil {
		return nil, err
	}

	return &schemas.Professionals{Professionals: professionals}, nil
}

// Bulk deletes professionals.
func (p *Professional) BulkDeleteProfessionals(
	bulkDeleteProfessionalData schemas.BulkDeleteProfessionalRequest,
) *errors.Error {
	return p.Adapter.Professional.BulkDeletePostgresqlProfessionals(
		bulkDeleteProfessionalData.Professionals,
	)
}
