package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Local struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Local adapter
func NewLocalAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Local {
	return &Local{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a Local from postgresql DB.
func (l *Local) GetPostgresqlLocal(
	localId uuid.UUID,
) (*schemas.Local, *errors.Error) {
	localModel, err := l.DaoPostgresql.Local.GetLocal(localId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.LocalNotFound
	}
	return &schemas.Local{
		Id:             localModel.Id,
		LocalName:      localModel.LocalName,
		StreetName:     localModel.StreetName,
		BuildingNumber: localModel.BuildingNumber,
		District:       localModel.District,
		Province:       localModel.Province,
		Region:         localModel.Region,
		Reference:      localModel.Reference,
		Capacity:       localModel.Capacity,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Fetch all locals from postgresql DB.
func (l *Local) FetchPostgresqlLocals() ([]*schemas.Local, *errors.Error) {
	localsModel, err := l.DaoPostgresql.Local.FetchLocals()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.LocalNotFound
	}

	locals := make([]*schemas.Local, len(localsModel))
	for i, localModel := range localsModel {
		locals[i] = &schemas.Local{
			Id:             localModel.Id,
			LocalName:      localModel.LocalName,
			StreetName:     localModel.StreetName,
			BuildingNumber: localModel.BuildingNumber,
			District:       localModel.District,
			Province:       localModel.Province,
			Region:         localModel.Region,
			Reference:      localModel.Reference,
			Capacity:       localModel.Capacity,
			ImageUrl:       localModel.ImageUrl,
		}
	}
	return locals, nil
}

// Creates a local into postgresql DB and returns it.
func (l *Local) CreatePostgresqlLocal(
	localName string,
	streetName string,
	buildingNumber string,
	district string,
	province string,
	region string,
	reference string,
	capacity int,
	imageUrl string,
	updatedBy string,
) (*schemas.Local, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localModel := &model.Local{
		Id:             uuid.New(),
		LocalName:      localName,
		StreetName:     streetName,
		BuildingNumber: buildingNumber,
		District:       district,
		Province:       province,
		Region:         region,
		Reference:      reference,
		Capacity:       capacity,
		ImageUrl:       imageUrl,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := l.DaoPostgresql.Local.CreateLocal(localModel); err != nil {
		return nil, &errors.BadRequestError.LocalNotCreated
	}

	return &schemas.Local{
		Id:             localModel.Id,
		LocalName:      localModel.LocalName,
		StreetName:     localModel.StreetName,
		BuildingNumber: localModel.BuildingNumber,
		District:       localModel.District,
		Province:       localModel.Province,
		Region:         localModel.Region,
		Reference:      localModel.Reference,
		Capacity:       localModel.Capacity,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Creates multiple locals into postgresql DB and returns them.
func (l *Local) BulkCreatePostgresqlLocals(
	localsData []*schemas.CreateLocalRequest,
	updatedBy string,
) ([]*schemas.Local, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localsModel := make([]*model.Local, len(localsData))
	for i, localData := range localsData {
		localsModel[i] = &model.Local{
			Id:             uuid.New(),
			LocalName:      localData.LocalName,
			StreetName:     localData.StreetName,
			BuildingNumber: localData.BuildingNumber,
			District:       localData.District,
			Province:       localData.Province,
			Region:         localData.Region,
			Reference:      localData.Reference,
			Capacity:       localData.Capacity,
			ImageUrl:       localData.ImageUrl,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}

	if err := l.DaoPostgresql.Local.BulkCreateLocals(localsModel); err != nil {
		return nil, &errors.BadRequestError.LocalNotCreated
	}

	locals := make([]*schemas.Local, len(localsModel))
	for i, localModel := range localsModel {
		locals[i] = &schemas.Local{
			Id:             localModel.Id,
			LocalName:      localModel.LocalName,
			StreetName:     localModel.StreetName,
			BuildingNumber: localModel.BuildingNumber,
			District:       localModel.District,
			Province:       localModel.Province,
			Region:         localModel.Region,
			Reference:      localModel.Reference,
			Capacity:       localModel.Capacity,
			ImageUrl:       localModel.ImageUrl,
		}
	}

	return locals, nil
}

// Updates a local given fields in postgresql DB and returns it.
func (l *Local) UpdatePostgresqlLocal(
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
) (*schemas.Local, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localModel, err := l.DaoPostgresql.Local.UpdateLocal(
		id,
		localName,
		streetName,
		buildingNumber,
		district,
		province,
		region,
		reference,
		capacity,
		imageUrl,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.LocalNotUpdated
	}

	return &schemas.Local{
		Id:             localModel.Id,
		LocalName:      localModel.LocalName,
		StreetName:     localModel.StreetName,
		BuildingNumber: localModel.BuildingNumber,
		District:       localModel.District,
		Province:       localModel.Province,
		Region:         localModel.Region,
		Reference:      localModel.Reference,
		Capacity:       localModel.Capacity,
		ImageUrl:       localModel.ImageUrl,
	}, nil
}

// Delets a plan from postgresql BD
func (l *Local) DeletePostgresqlLocal(localId uuid.UUID) *errors.Error {
	if err := l.DaoPostgresql.Local.DeleteLocal(localId); err != nil {
		return &errors.BadRequestError.LocalNotSoftDeleted
	}

	return nil
}

// Bulk deletes locals from postgresql DB
func (l *Local) BulkDeletePostgresqlLocals(localIds []string) *errors.Error {
	// Convert string IDs to UUIDs
	uuidIds := make([]uuid.UUID, len(localIds))
	for i, id := range localIds {
		parsedId, err := uuid.Parse(id)
		if err != nil {
			return &errors.UnprocessableEntityError.InvalidLocalId
		}
		uuidIds[i] = parsedId
	}

	if err := l.DaoPostgresql.Local.BulkDeleteLocals(uuidIds); err != nil {
		return &errors.BadRequestError.LocalNotSoftDeleted
	}

	return nil
}

// Bulk creates professionals into postgresql DB.
func (l *Local) BulkCreatePostgresqlLocals(
	localsData []*schemas.CreateLocalRequest,
	updatedBy string,
) ([]*schemas.Local, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	localsModel := make([]*model.Local, len(localsData))
	for i, localData := range localsData {
		localsModel[i] = &model.Local{
			Id:             uuid.New(),
			LocalName:      localData.LocalName,
			StreetName:     localData.StreetName,
			BuildingNumber: localData.BuildingNumber,
			District:       localData.District,
			Province:       localData.Province,
			Region:         localData.Region,
			Reference:      localData.Reference,
			Capacity:       localData.Capacity,
			ImageUrl:       localData.ImageUrl,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	if err := l.DaoPostgresql.Local.BulkCreateLocals(localsModel); err != nil {
		return nil, &errors.BadRequestError.LocalNotCreated
	}
	locals := make([]*schemas.Local, len(localsModel))
	for i, localModel := range localsModel {
		locals[i] = &schemas.Local{
			Id:             localModel.Id,
			LocalName:      localModel.LocalName,
			StreetName:     localModel.StreetName,
			BuildingNumber: localModel.BuildingNumber,
			District:       localModel.District,
			Province:       localModel.Province,
			Region:         localModel.Region,
			Reference:      localModel.Reference,
			Capacity:       localModel.Capacity,
			ImageUrl:       localModel.ImageUrl,
		}
	}
	return locals, nil
}
