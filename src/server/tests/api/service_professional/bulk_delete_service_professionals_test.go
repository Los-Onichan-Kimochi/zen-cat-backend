package service_professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestBulkDeleteServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid service-professional associations exist
		WHEN:  DELETE /service-professional/bulk/ is called with valid associations
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create service-professional associations using factories
	serviceProfessional1 := factories.NewServiceProfessionalModel(db)
	serviceProfessional2 := factories.NewServiceProfessionalModel(db)

	bulkRequest := schemas.BulkDeleteServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.DeleteServiceProfessionalRequest{
			{
				ServiceId:      serviceProfessional1.ServiceId,
				ProfessionalId: serviceProfessional1.ProfessionalId,
			},
			{
				ServiceId:      serviceProfessional2.ServiceId,
				ProfessionalId: serviceProfessional2.ProfessionalId,
			},
		},
	}

	requestBody, _ := json.Marshal(bulkRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestBulkDeleteServiceProfessionalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  DELETE /service-professional/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteServiceProfessionalsNonExistent(t *testing.T) {
	/*
		GIVEN: Non-existent service-professional associations
		WHEN:  DELETE /service-professional/bulk/ is called with non-existent associations
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create delete request with non-existent associations
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	nonExistentProfessionalId, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	request := schemas.BulkDeleteServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.DeleteServiceProfessionalRequest{
			{
				ServiceId:      nonExistentServiceId,
				ProfessionalId: nonExistentProfessionalId,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
