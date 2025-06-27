package adapter

import (
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Onboarding struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

func NewOnboardingAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Onboarding {
	return &Onboarding{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (o *Onboarding) GetPostgresqlOnboarding(
	onboardingId uuid.UUID,
) (*schemas.Onboarding, *errors.Error) {
	onboardingModel, err := o.DaoPostgresql.Onboarding.GetOnboarding(onboardingId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.OnboardingNotFound
	}

	return &schemas.Onboarding{
		Id:             onboardingModel.Id,
		DocumentType:   schemas.DocumentType(onboardingModel.DocumentType),
		DocumentNumber: onboardingModel.DocumentNumber,
		PhoneNumber:    onboardingModel.PhoneNumber,
		BirthDate:      onboardingModel.BirthDate,
		Gender:         (*schemas.Gender)(onboardingModel.Gender),
		PostalCode:     onboardingModel.PostalCode,
		District:       onboardingModel.District,
		Province:       onboardingModel.Province,
		Region:         onboardingModel.Region,
		Address:        onboardingModel.Address,
		UserId:         onboardingModel.UserId,
	}, nil
}

func (o *Onboarding) GetPostgresqlOnboardingByUserId(
	userId uuid.UUID,
) (*schemas.Onboarding, *errors.Error) {
	onboardingModel, err := o.DaoPostgresql.Onboarding.GetOnboardingByUserId(userId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.OnboardingNotFound
	}

	return &schemas.Onboarding{
		Id:             onboardingModel.Id,
		DocumentType:   schemas.DocumentType(onboardingModel.DocumentType),
		DocumentNumber: onboardingModel.DocumentNumber,
		PhoneNumber:    onboardingModel.PhoneNumber,
		BirthDate:      onboardingModel.BirthDate,
		Gender:         (*schemas.Gender)(onboardingModel.Gender),
		PostalCode:     onboardingModel.PostalCode,
		District:       onboardingModel.District,
		Province:       onboardingModel.Province,
		Region:         onboardingModel.Region,
		Address:        onboardingModel.Address,
		UserId:         onboardingModel.UserId,
	}, nil
}

func (o *Onboarding) FetchPostgresqlOnboardings() ([]*schemas.Onboarding, *errors.Error) {
	onboardingsModel, err := o.DaoPostgresql.Onboarding.FetchOnboardings()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.OnboardingNotFound
	}

	onboardings := make([]*schemas.Onboarding, len(onboardingsModel))
	for i, onboardingModel := range onboardingsModel {
		onboardings[i] = &schemas.Onboarding{
			Id:             onboardingModel.Id,
			DocumentType:   schemas.DocumentType(onboardingModel.DocumentType),
			DocumentNumber: onboardingModel.DocumentNumber,
			PhoneNumber:    onboardingModel.PhoneNumber,
			BirthDate:      onboardingModel.BirthDate,
			Gender:         (*schemas.Gender)(onboardingModel.Gender),
			PostalCode:     onboardingModel.PostalCode,
			District:       onboardingModel.District,
			Province:       onboardingModel.Province,
			Region:         onboardingModel.Region,
			Address:        onboardingModel.Address,
			UserId:         onboardingModel.UserId,
		}
	}

	return onboardings, nil
}

func (o *Onboarding) CreatePostgresqlOnboarding(
	userId uuid.UUID,
	documentType schemas.DocumentType,
	documentNumber string,
	phoneNumber string,
	birthDate *string,
	gender *schemas.Gender,
	postalCode string,
	district *string,
	province *string,
	region *string,
	address string,
	updatedBy string,
) (*schemas.Onboarding, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Validar que el usuario existe
	_, err := o.DaoPostgresql.User.GetUser(userId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.UserNotFound
	}

	// Verificar si ya existe un onboarding para este usuario
	existingOnboarding, _ := o.DaoPostgresql.Onboarding.GetOnboardingByUserId(userId)
	if existingOnboarding != nil {
		return nil, &errors.BadRequestError.OnboardingNotCreated // Ya existe
	}

	// Parsear fecha de nacimiento si se proporciona
	var parsedBirthDate *time.Time
	if birthDate != nil && *birthDate != "" {
		parsed, parseErr := time.Parse("2006-01-02", *birthDate)
		if parseErr != nil {
			return nil, &errors.BadRequestError.OnboardingNotCreated
		}
		parsedBirthDate = &parsed
	}

	onboardingModel := &model.Onboarding{
		Id:             uuid.New(),
		DocumentType:   model.DocumentType(documentType),
		DocumentNumber: documentNumber,
		PhoneNumber:    phoneNumber,
		BirthDate:      parsedBirthDate,
		Gender:         (*model.Gender)(gender),
		PostalCode:     postalCode,
		District:       district,
		Province:       province,
		Region:         region,
		Address:        address,
		UserId:         userId,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := o.DaoPostgresql.Onboarding.CreateOnboarding(onboardingModel); err != nil {
		return nil, &errors.BadRequestError.OnboardingNotCreated
	}

	return &schemas.Onboarding{
		Id:             onboardingModel.Id,
		DocumentType:   schemas.DocumentType(onboardingModel.DocumentType),
		DocumentNumber: onboardingModel.DocumentNumber,
		PhoneNumber:    onboardingModel.PhoneNumber,
		BirthDate:      onboardingModel.BirthDate,
		Gender:         (*schemas.Gender)(onboardingModel.Gender),
		PostalCode:     onboardingModel.PostalCode,
		District:       onboardingModel.District,
		Province:       onboardingModel.Province,
		Region:         onboardingModel.Region,
		Address:        onboardingModel.Address,
		UserId:         onboardingModel.UserId,
	}, nil
}

func (o *Onboarding) UpdatePostgresqlOnboarding(
	onboardingId uuid.UUID,
	documentType *schemas.DocumentType,
	documentNumber *string,
	phoneNumber *string,
	birthDate *string,
	gender *schemas.Gender,
	postalCode *string,
	district *string,
	province *string,
	region *string,
	address *string,
	updatedBy string,
) (*schemas.Onboarding, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Parsear fecha de nacimiento si se proporciona
	var parsedBirthDate *time.Time
	if birthDate != nil && *birthDate != "" {
		parsed, parseErr := time.Parse("2006-01-02", *birthDate)
		if parseErr != nil {
			return nil, &errors.BadRequestError.OnboardingNotUpdated
		}
		parsedBirthDate = &parsed
	}

	onboardingModel, err := o.DaoPostgresql.Onboarding.UpdateOnboarding(
		onboardingId,
		(*model.DocumentType)(documentType),
		documentNumber,
		phoneNumber,
		parsedBirthDate,
		(*model.Gender)(gender),
		postalCode,
		district,
		province,
		region,
		address,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.OnboardingNotUpdated
	}

	return &schemas.Onboarding{
		Id:             onboardingModel.Id,
		DocumentType:   schemas.DocumentType(onboardingModel.DocumentType),
		DocumentNumber: onboardingModel.DocumentNumber,
		PhoneNumber:    onboardingModel.PhoneNumber,
		BirthDate:      onboardingModel.BirthDate,
		Gender:         (*schemas.Gender)(onboardingModel.Gender),
		PostalCode:     onboardingModel.PostalCode,
		District:       onboardingModel.District,
		Province:       onboardingModel.Province,
		Region:         onboardingModel.Region,
		Address:        onboardingModel.Address,
		UserId:         onboardingModel.UserId,
	}, nil
}

func (o *Onboarding) DeletePostgresqlOnboarding(onboardingId uuid.UUID) *errors.Error {
	if err := o.DaoPostgresql.Onboarding.DeleteOnboarding(onboardingId); err != nil {
		return &errors.ObjectNotFoundError.OnboardingNotFound
	}
	return nil
}

func (o *Onboarding) DeletePostgresqlOnboardingByUserId(userId uuid.UUID) *errors.Error {
	if err := o.DaoPostgresql.Onboarding.DeleteOnboardingByUserId(userId); err != nil {
		return &errors.ObjectNotFoundError.OnboardingNotFound
	}
	return nil
}
