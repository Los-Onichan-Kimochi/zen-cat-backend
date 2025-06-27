package service_local_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
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

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocals
	err = json.NewDecoder(rec.Body).Decode(&response)
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

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-local/?serviceId="+service1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceLocals
	err = json.NewDecoder(rec.Body).Decode(&response)
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
