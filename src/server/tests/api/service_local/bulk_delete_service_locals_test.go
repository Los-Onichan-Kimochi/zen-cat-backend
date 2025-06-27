package service_local_test

import (
	"bytes"
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

func TestBulkDeleteServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Existing service-local associations
		WHEN:  DELETE /service-local/bulk/ is called with valid associations
		THEN:  A HTTP_204_NO_CONTENT status should be returned and associations should be deleted
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create services
	service1 := &model.Service{
		Name:        "Test Service 1",
		Description: "Test Description 1",
		ImageUrl:    "https://example.com/image1.jpg",
		IsVirtual:   false,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	service2 := &model.Service{
		Name:        "Test Service 2",
		Description: "Test Description 2",
		ImageUrl:    "https://example.com/image2.jpg",
		IsVirtual:   true,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create([]*model.Service{service1, service2}).Error
	assert.NoError(t, err)

	// Create locals
	local1 := &model.Local{
		LocalName:      "Test Local 1",
		StreetName:     "Test Street 1",
		BuildingNumber: "123",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	local2 := &model.Local{
		LocalName:      "Test Local 2",
		StreetName:     "Test Street 2",
		BuildingNumber: "456",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.Local{local1, local2}).Error
	assert.NoError(t, err)

	// Create service-local associations
	serviceLocal1 := &model.ServiceLocal{
		ServiceId: service1.Id,
		LocalId:   local1.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	serviceLocal2 := &model.ServiceLocal{
		ServiceId: service2.Id,
		LocalId:   local2.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.ServiceLocal{serviceLocal1, serviceLocal2}).Error
	assert.NoError(t, err)

	// Create delete request
	request := schemas.BulkDeleteServiceLocalRequest{
		ServiceLocals: []*schemas.DeleteServiceLocalRequest{
			{
				ServiceId: service1.Id,
				LocalId:   local1.Id,
			},
			{
				ServiceId: service2.Id,
				LocalId:   local2.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the associations were deleted
	var count int64
	db.Model(&model.ServiceLocal{}).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestBulkDeleteServiceLocalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  DELETE /service-local/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkDeleteServiceLocalsNonExistent(t *testing.T) {
	/*
		GIVEN: Non-existent service-local associations
		WHEN:  DELETE /service-local/bulk/ is called with non-existent associations
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Create delete request with non-existent associations
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	nonExistentLocalId, _ := uuid.Parse("00000000-0000-0000-0000-000000000001")
	request := schemas.BulkDeleteServiceLocalRequest{
		ServiceLocals: []*schemas.DeleteServiceLocalRequest{
			{
				ServiceId: nonExistentServiceId,
				LocalId:   nonExistentLocalId,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
