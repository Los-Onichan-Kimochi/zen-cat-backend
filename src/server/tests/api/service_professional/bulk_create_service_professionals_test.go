package service_professional_test

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
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestBulkCreateServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid services and professionals
		WHEN:  POST /service-professional/bulk/ is called with valid associations
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

	// Create professionals
	professional1 := &model.Professional{
		Name:          "Dr. Smith",
		FirstLastName: "Johnson",
		Specialty:     "Cardiology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "123456789",
		Type:          model.ProfessionalTypeMedic,
		ImageUrl:      "https://example.com/doctor1.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	professional2 := &model.Professional{
		Name:          "Dr. Jane",
		FirstLastName: "Doe",
		Specialty:     "Neurology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "987654321",
		Type:          model.ProfessionalTypeGymTrainer,
		ImageUrl:      "https://example.com/doctor2.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.Professional{professional1, professional2}).Error
	assert.NoError(t, err)

	// Create request
	request := schemas.BatchCreateServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.CreateServiceProfessionalRequest{
			{
				ServiceId:      service1.Id,
				ProfessionalId: professional1.Id,
			},
			{
				ServiceId:      service2.Id,
				ProfessionalId: professional2.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.ServiceProfessionals
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(response.ServiceProfessionals))
}

func TestBulkCreateServiceProfessionalsInvalidRequest(t *testing.T) {
	/*
		GIVEN: Invalid request body
		WHEN:  POST /service-professional/bulk/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateServiceProfessionalsNonExistentService(t *testing.T) {
	/*
		GIVEN: Non-existent service ID
		WHEN:  POST /service-professional/bulk/ is called with invalid service ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a professional
	professional := &model.Professional{
		Name:          "Dr. Smith",
		FirstLastName: "Johnson",
		Specialty:     "Cardiology",
		Email:         utilsTest.GenerateRandomEmail(),
		PhoneNumber:   "123456789",
		Type:          model.ProfessionalTypeMedic,
		ImageUrl:      "https://example.com/doctor.jpg",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(professional).Error
	assert.NoError(t, err)

	// Create request with non-existent service
	nonExistentServiceId, _ := uuid.Parse("00000000-0000-0000-0000-000000000000")
	request := schemas.BatchCreateServiceProfessionalRequest{
		ServiceProfessionals: []*schemas.CreateServiceProfessionalRequest{
			{
				ServiceId:      nonExistentServiceId, // Non-existent service
				ProfessionalId: professional.Id,
			},
		},
	}

	requestBody, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/service-professional/bulk/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
