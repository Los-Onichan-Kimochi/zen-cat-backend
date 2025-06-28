package service_local_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing service-local association
		WHEN:  GET /service-local/{serviceId}/{localId}/ is called with valid IDs
		THEN:  A HTTP_200_OK status should be returned with the association data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a service
	service := &model.Service{
		Name:        "Test Service",
		Description: "Test Description",
		ImageUrl:    "https://example.com/image.jpg",
		IsVirtual:   false,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(service).Error
	assert.NoError(t, err)

	// Create a local
	local := &model.Local{
		LocalName:      "Test Local",
		StreetName:     "Test Street",
		BuildingNumber: "123",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(local).Error
	assert.NoError(t, err)

	// Create service-local association
	serviceLocal := &model.ServiceLocal{
		ServiceId: service.Id,
		LocalId:   local.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(serviceLocal).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/"+service.Id.String()+"/"+local.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocal
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, serviceLocal.ServiceId, response.ServiceId)
	assert.Equal(t, serviceLocal.LocalId, response.LocalId)
}

func TestGetServiceLocalNotFound(t *testing.T) {
	/*
		GIVEN: Non-existent service and local IDs
		WHEN:  GET /service-local/{serviceId}/{localId}/ is called with invalid IDs
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	serviceId := uuid.New()
	localId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/"+serviceId.String()+"/"+localId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetServiceLocalInvalidIds(t *testing.T) {
	/*
		GIVEN: Invalid UUID formats for service and local IDs
		WHEN:  GET /service-local/{serviceId}/{localId}/ is called with invalid UUIDs
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-service-id"
	invalidLocalId := "invalid-local-id"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/"+invalidServiceId+"/"+invalidLocalId+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
