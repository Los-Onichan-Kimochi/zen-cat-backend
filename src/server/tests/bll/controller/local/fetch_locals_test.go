package local_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestFetchLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple local records exist in the database
		WHEN:  FetchLocals is called
		THEN:  All local records should be returned
	*/
	// GIVEN
	localController, _, db := controllerTest.NewLocalControllerTestWrapper(t)

	// Create local records
	locals := []*model.Local{
		{
			LocalName:      "Test Gym 1",
			StreetName:     "Main Street",
			BuildingNumber: "123",
			District:       "Downtown",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Near the park",
			Capacity:       50,
			ImageUrl:       "https://example.com/gym1.jpg",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
		{
			LocalName:      "Test Gym 2",
			StreetName:     "Second Street",
			BuildingNumber: "456",
			District:       "Uptown",
			Province:       "Lima",
			Region:         "Lima",
			Reference:      "Near the mall",
			Capacity:       30,
			ImageUrl:       "https://example.com/gym2.jpg",
			AuditFields: model.AuditFields{
				UpdatedBy: "ADMIN",
			},
		},
	}
	err := db.Create(locals).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := localController.FetchLocals()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.GreaterOrEqual(t, len(result.Locals), 2)
}

func TestFetchLocalsEmpty(t *testing.T) {
	/*
		GIVEN: No local records exist in the database
		WHEN:  FetchLocals is called
		THEN:  An empty list should be returned
	*/
	// GIVEN
	localController, _, _ := controllerTest.NewLocalControllerTestWrapper(t)

	// WHEN
	result, errResult := localController.FetchLocals()

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, 0, len(result.Locals))
}
