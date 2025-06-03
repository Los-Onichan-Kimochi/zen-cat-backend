package controller

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Local struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Local controller
func NewLocalController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Local {
	return &Local{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a local.
func (l *Local) GetLocal(localId uuid.UUID) (*schemas.Local, *errors.Error) {
	return l.Adapter.Local.GetPostgresqlLocal(localId)
}

// Fetch all locals.
func (l *Local) FetchLocals() (*schemas.Locals, *errors.Error) {
	locals, err := l.Adapter.Local.FetchPostgresqlLocals()
	if err != nil {
		return nil, err
	}

	return &schemas.Locals{Locals: locals}, nil
}

// Creates a local.
func (l *Local) CreateLocal(
	createLocalData schemas.CreateLocalRequest,
	updatedBy string,
) (*schemas.Local, *errors.Error) {
	return l.Adapter.Local.CreatePostgresqlLocal(
		createLocalData.LocalName,
		createLocalData.StreetName,
		createLocalData.BuildingNumber,
		createLocalData.District,
		createLocalData.Province,
		createLocalData.Region,
		createLocalData.Reference,
		createLocalData.Capacity,
		createLocalData.ImageUrl,
		updatedBy,
	)
}

// Update a local.
func (l *Local) UpdateLocal(
	localId uuid.UUID,
	updateLocalData schemas.UpdateLocalRequest,
	updatedBy string,
) (*schemas.Local, *errors.Error) {
	return l.Adapter.Local.UpdatePostgresqlLocal(
		localId,
		updateLocalData.LocalName,
		updateLocalData.StreetName,
		updateLocalData.BuildingNumber,
		updateLocalData.District,
		updateLocalData.Province,
		updateLocalData.Region,
		updateLocalData.Reference,
		updateLocalData.Capacity,
		updateLocalData.ImageUrl,
		updatedBy,
	)
}

// Deletes a local.
func (l *Local) DeleteLocal(localId uuid.UUID) *errors.Error {
	return l.Adapter.Local.DeletePostgresqlLocal(localId)
}

// Bulk creates locals.
func (l *Local) BulkCreateLocals(
	createLocalsData []*schemas.CreateLocalRequest,
	updatedBy string,
)([]*schemas.Local, *errors.Error) {
	return l.Adapter.Local.BulkCreatePostgresqlLocals(
		createLocalsData,
		updatedBy,
	)
}
// Bulk deletes locals.
func (l *Local) BulkDeleteLocals(
	bulkDeleteLocalData schemas.BulkDeleteLocalRequest,
) *errors.Error {
	return l.Adapter.Local.BulkDeletePostgresqlLocals(
		bulkDeleteLocalData.Locals,
	)
}
