package service_professional_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  GET /service-professional/ is called
		THEN:  A HTTP_200_OK status should be returned with all associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create service-professional associations using factories
	_ = factories.NewServiceProfessionalModel(db)
	_ = factories.NewServiceProfessionalModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.ServiceProfessionals), 2)
}

func TestFetchServiceProfessionalsByServiceId(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  GET /service-professional/?serviceId={serviceId} is called
		THEN:  A HTTP_200_OK status should be returned with filtered associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create services and professionals using factories
	service1 := factories.NewServiceModel(db)
	service2 := factories.NewServiceModel(db)
	professional1 := factories.NewProfessionalModel(db)
	professional2 := factories.NewProfessionalModel(db)

	// Create service-professional associations
	_ = factories.NewServiceProfessionalModel(db, factories.ServiceProfessionalModelF{
		ServiceId:      &service1.Id,
		ProfessionalId: &professional1.Id,
	})
	_ = factories.NewServiceProfessionalModel(db, factories.ServiceProfessionalModelF{
		ServiceId:      &service2.Id,
		ProfessionalId: &professional2.Id,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/?serviceId="+service1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.ServiceProfessionals))
	if len(response.ServiceProfessionals) > 0 {
		assert.Equal(t, service1.Id, response.ServiceProfessionals[0].ServiceId)
	}
}

func TestFetchServiceProfessionalsEmpty(t *testing.T) {
	/*
		GIVEN: No service-professional associations exist
		WHEN:  GET /service-professional/ is called
		THEN:  A HTTP_200_OK status should be returned with empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.ServiceProfessionals))
}
