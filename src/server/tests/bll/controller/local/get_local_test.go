package local_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: A local record exists in the database
		WHEN:  GetLocal is called with valid ID
		THEN:  The local record should be returned
	*/
	// GIVEN
	localController, _, db := controllerTest.NewLocalControllerTestWrapper(t)

	// Create a local record
	local := &model.Local{
		LocalName:      "Test Gym",
		StreetName:     "Main Street",
		BuildingNumber: "123",
		District:       "Downtown",
		Province:       "Lima",
		Region:         "Lima",
		Reference:      "Near the park",
		Capacity:       50,
		ImageUrl:       "https://example.com/gym.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(local).Error
	assert.NoError(t, err)

	// WHEN
	result, errResult := localController.GetLocal(local.Id)

	// THEN
	assert.Nil(t, errResult)
	assert.NotNil(t, result)
	assert.Equal(t, local.Id, result.Id)
	assert.Equal(t, local.LocalName, result.LocalName)
	assert.Equal(t, local.StreetName, result.StreetName)
	assert.Equal(t, local.BuildingNumber, result.BuildingNumber)
	assert.Equal(t, local.District, result.District)
}

func TestGetLocalNotFound(t *testing.T) {
	/*
		GIVEN: No local record exists with the given ID
		WHEN:  GetLocal is called with non-existent ID
		THEN:  An error should be returned
	*/
	// GIVEN
	localController, _, _ := controllerTest.NewLocalControllerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	result, errResult := localController.GetLocal(nonExistentId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
	assert.Contains(t, errResult.Message, "not found")
}

func TestGetLocalWithNilId(t *testing.T) {
	/*
		GIVEN: A nil UUID
		WHEN:  GetLocal is called with nil UUID
		THEN:  An error should be returned
	*/
	// GIVEN
	localController, _, _ := controllerTest.NewLocalControllerTestWrapper(t)
	nilId := uuid.Nil

	// WHEN
	result, errResult := localController.GetLocal(nilId)

	// THEN
	assert.NotNil(t, errResult)
	assert.Nil(t, result)
}
