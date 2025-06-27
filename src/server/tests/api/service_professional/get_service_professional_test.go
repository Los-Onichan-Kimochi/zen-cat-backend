package service_professional_test

import (
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

func TestGetServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A service-professional association exists
		WHEN:  GET /service-professional/{serviceId}/{professionalId}/ is called
		THEN:  A HTTP_200_OK status should be returned with the association
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create dependencies first
	service := factories.NewServiceModel(db, factories.ServiceModelF{})
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Create the association
	serviceProfessional := factories.NewServiceProfessionalModel(db, factories.ServiceProfessionalModelF{
		ServiceId:      &service.Id,
		ProfessionalId: &professional.Id,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/"+service.Id.String()+"/"+professional.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessional
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	assert.Equal(t, serviceProfessional.Id, response.Id)
	assert.Equal(t, service.Id, response.ServiceId)
	assert.Equal(t, professional.Id, response.ProfessionalId)
}

func TestGetServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No association exists between service and professional
		WHEN:  GET /service-professional/{serviceId}/{professionalId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	serviceId := uuid.New()
	professionalId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/"+serviceId.String()+"/"+professionalId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetServiceProfessionalEmptyParams(t *testing.T) {
	/*
		GIVEN: Empty service or professional ID parameters
		WHEN:  GET /service-professional/// is called with empty parameters
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN - Empty service ID
	req := httptest.NewRequest(http.MethodGet, "/service-professional//"+uuid.New().String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
