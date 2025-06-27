package service_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteServiceSuccessfully(t *testing.T) {
	/*
		GIVEN: A service exists in the database
		WHEN:  DELETE /service/{serviceId}/ is called with a valid service ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and the service should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test service using factory
	service := factories.NewServiceModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service/"+service.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the service was deleted
	var count int64
	db.Model(&service).Where("id = ?", service.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteServiceNotFound(t *testing.T) {
	/*
		GIVEN: No service exists with the provided ID
		WHEN:  DELETE /service/{serviceId}/ is called with a non-existent service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentServiceId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service/"+nonExistentServiceId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteServiceInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID is provided
		WHEN:  DELETE /service/{serviceId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service/"+invalidServiceId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
