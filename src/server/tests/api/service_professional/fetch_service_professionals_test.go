package service_professional_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestFetchServiceProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple service-professional associations exist
		WHEN:  GET /service-professional/ is called
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

	// Create service-professional associations
	serviceProfessional1 := &model.ServiceProfessional{
		ServiceId:      service1.Id,
		ProfessionalId: professional1.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	serviceProfessional2 := &model.ServiceProfessional{
		ServiceId:      service2.Id,
		ProfessionalId: professional2.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.ServiceProfessional{serviceProfessional1, serviceProfessional2}).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessionals
	err = json.NewDecoder(rec.Body).Decode(&response)
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

	// Create service-professional associations
	serviceProfessional1 := &model.ServiceProfessional{
		ServiceId:      service1.Id,
		ProfessionalId: professional1.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	serviceProfessional2 := &model.ServiceProfessional{
		ServiceId:      service2.Id,
		ProfessionalId: professional2.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.ServiceProfessional{serviceProfessional1, serviceProfessional2}).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/service-professional/?serviceId="+service1.Id.String(), nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.ServiceProfessionals
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(response.ServiceProfessionals))
	assert.Equal(t, service1.Id, response.ServiceProfessionals[0].ServiceId)
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
