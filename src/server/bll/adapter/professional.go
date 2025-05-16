package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Professional struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Professional adapter
func NewProfessionalAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Professional {
	return &Professional{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a professional from postgresql DB.
func (p *Professional) GetPostgresqlProfessional(
	professionalId uuid.UUID,
) (*schemas.Professional, *errors.Error) {
	professionalModel, err := p.DaoPostgresql.Professional.GetProfessional(professionalId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ProfessionalNotFound
	}

	return &schemas.Professional{
		Id:             professionalModel.Id,
		Name:           professionalModel.Name,
		FirstLastName:  professionalModel.FirstLastName,
		SecondLastName: professionalModel.SecondLastName,
		Specialty:      professionalModel.Specialty,
		Email:          professionalModel.Email,
		PhoneNumber:    professionalModel.PhoneNumber,
		Type:           string(professionalModel.Type),
		ImageUrl:       professionalModel.ImageUrl,
	}, nil
}

// Fetch all professionals from postgresql DB.
func (p *Professional) FetchPostgresqlProfessionals() ([]*schemas.Professional, *errors.Error) {
	professionalsModel, err := p.DaoPostgresql.Professional.FetchProfessionals()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ProfessionalNotFound
	}

	professionals := make([]*schemas.Professional, len(professionalsModel))
	for i, professionalModel := range professionalsModel {
		professionals[i] = &schemas.Professional{
			Id:             professionalModel.Id,
			Name:           professionalModel.Name,
			FirstLastName:  professionalModel.FirstLastName,
			SecondLastName: professionalModel.SecondLastName,
			Specialty:      professionalModel.Specialty,
			Email:          professionalModel.Email,
			PhoneNumber:    professionalModel.PhoneNumber,
			Type:           string(professionalModel.Type),
			ImageUrl:       professionalModel.ImageUrl,
		}
	}
	return professionals, nil
}

// Creates a professional into postgresql DB and returns it.
func (p *Professional) CreatePostgresqlProfessional(
	name string,
	firstLastName string,
	secondLastName *string,
	specialty string,
	email string,
	phoneNumber string,
	professionalType string,
	imageUrl string,
	updatedBy string,
) (*schemas.Professional, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	professionalModel := &model.Professional{
		Id:             uuid.New(),
		Name:           name,
		FirstLastName:  firstLastName,
		SecondLastName: secondLastName,
		Specialty:      specialty,
		Email:          email,
		PhoneNumber:    phoneNumber,
		Type:           model.ProfessionalType(professionalType),
		ImageUrl:       imageUrl,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := p.DaoPostgresql.Professional.CreateProfessional(professionalModel); err != nil {
		return nil, &errors.BadRequestError.ProfessionalNotCreated
	}

	return &schemas.Professional{
		Id:             professionalModel.Id,
		Name:           professionalModel.Name,
		FirstLastName:  professionalModel.FirstLastName,
		SecondLastName: professionalModel.SecondLastName,
		Specialty:      professionalModel.Specialty,
		Email:          professionalModel.Email,
		PhoneNumber:    professionalModel.PhoneNumber,
		Type:           string(professionalModel.Type),
		ImageUrl:       professionalModel.ImageUrl,
	}, nil
}

// Updates a professional given fields in postgresql DB and returns it.
func (p *Professional) UpdatePostgresqlProfessional(
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
) (*schemas.Professional, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	professionalModel, err := p.DaoPostgresql.Professional.UpdateProfessional(
		id,
		name,
		firstLastName,
		secondLastName,
		specialty,
		email,
		phoneNumber,
		professionalType,
		imageUrl,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.ProfessionalNotFound
	}

	return &schemas.Professional{
		Id:             professionalModel.Id,
		Name:           professionalModel.Name,
		FirstLastName:  professionalModel.FirstLastName,
		SecondLastName: professionalModel.SecondLastName,
		Specialty:      professionalModel.Specialty,
		Email:          professionalModel.Email,
		PhoneNumber:    professionalModel.PhoneNumber,
		Type:           string(professionalModel.Type),
		ImageUrl:       professionalModel.ImageUrl,
	}, nil
}

// Deletes a professional from postgresql DB.
func (p *Professional) DeletePostgresqlProfessional(id uuid.UUID) *errors.Error {
	if err := p.DaoPostgresql.Professional.DeleteProfessional(id); err != nil {
		return &errors.ObjectNotFoundError.ProfessionalNotFound
	}
	return nil
}
