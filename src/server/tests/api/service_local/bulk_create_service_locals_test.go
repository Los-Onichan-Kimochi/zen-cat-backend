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

func TestBulkCreateServiceLocalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid services and locals
		WHEN:  POST /service-local/bulk/ is called with valid associations
		THEN:  A HTTP_201_CREATED status should be returned with created associations
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

	// Create request
	request := schemas.BatchCreateServiceLocalRequest{
		ServiceLocals: []*schemas.CreateServiceLocalRequest{
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
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceLocals
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.ServiceLocals))
}

func TestBulkCreateServiceLocalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  POST /service-local/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateServiceLocalsNonExistentService(t *testing.T) {
	/*
		GIVEN: Non-existent service ID
		WHEN:  POST /service-local/bulk/ is called with invalid service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a local
	local := &model.Local{
		LocalName:      "Test Local",
		StreetName:     "Test Street",
		BuildingNumber: "123",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(local).Error
	assert.NoError(t, err)

	// Create request with non-existent service
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	request := schemas.BatchCreateServiceLocalRequest{
		ServiceLocals: []*schemas.CreateServiceLocalRequest{
			{
				ServiceId: nonExistentServiceId, // Non-existent service
				LocalId:   local.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-local/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
