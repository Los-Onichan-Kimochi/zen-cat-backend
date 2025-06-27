package service_local_test

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

func TestFetchServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-local associations exist
		WHEN:  GET /service-local/ is called
		THEN:  A HTTP_200_OK status should be returned with all associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create service-local associations using factories
	_ = factories.NewServiceLocalModel(db)
	_ = factories.NewServiceLocalModel(db)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.ServiceLocals), 2)
}

func TestFetchServiceLocalsByServiceId(t *testing.T) {
	/*
		GIVEN: Multiple service-local associations exist
		WHEN:  GET /service-local/?serviceId={serviceId} is called
		THEN:  A HTTP_200_OK status should be returned with filtered associations
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create services using factories
	service1 := factories.NewServiceModel(db)
	service2 := factories.NewServiceModel(db)

	// Create locals using factories
	local1 := factories.NewLocalModel(db)
	local2 := factories.NewLocalModel(db)

	// Create service-local associations
	_ = factories.NewServiceLocalModel(db, factories.ServiceLocalModelF{
		ServiceId: &service1.Id,
		LocalId:   &local1.Id,
	})
	_ = factories.NewServiceLocalModel(db, factories.ServiceLocalModelF{
		ServiceId: &service2.Id,
		LocalId:   &local2.Id,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/?serviceId="+service1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.ServiceLocals))
	assert.Equal(t, service1.Id, response.ServiceLocals[0].ServiceId)
}

func TestFetchServiceLocalsEmpty(t *testing.T) {
	/*
		GIVEN: No service-local associations exist
		WHEN:  GET /service-local/ is called
		THEN:  A HTTP_200_OK status should be returned with empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.ServiceLocals))
}
