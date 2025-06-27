package service_local_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteServiceLocalSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing service-local association
		WHEN:  DELETE /service-local/{serviceId}/{localId}/ is called with valid IDs
		THEN:  A HTTP_204_NO_CONTENT status should be returned and association should be deleted
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
	req := httptest.NewRequest(http.MethodDelete, "/service-local/"+service.Id.String()+"/"+local.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the association was deleted
	var count int64
	db.Model(&model.ServiceLocal{}).Where("service_id = ? AND local_id = ?", service.Id, local.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteServiceLocalNotFound(t *testing.T) {
	/*
		GIVEN: Non-existent service and local IDs
		WHEN:  DELETE /service-local/{serviceId}/{localId}/ is called with invalid IDs
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	serviceId := uuid.New()
	localId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/"+serviceId.String()+"/"+localId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteServiceLocalInvalidIds(t *testing.T) {
	/*
		GIVEN: Invalid UUID formats for service and local IDs
		WHEN:  DELETE /service-local/{serviceId}/{localId}/ is called with invalid UUIDs
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-service-id"
	invalidLocalId := "invalid-local-id"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/"+invalidServiceId+"/"+invalidLocalId+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
