package service_professional_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestDeleteServiceProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing service-professional association
		WHEN:  DELETE /service-professional/{serviceId}/{professionalId}/ is called with valid IDs
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
	err = db.Create(professional).Error
	assert.NoError(t, err)

	// Create service-professional association
	serviceProfessional := &model.ServiceProfessional{
		ServiceId:      service.Id,
		ProfessionalId: professional.Id,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(serviceProfessional).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/"+service.Id.String()+"/"+professional.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the association was deleted
	var count int64
	db.Model(&model.ServiceProfessional{}).Where("service_id = ? AND professional_id = ?", service.Id, professional.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteServiceProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: Non-existent service and professional IDs
		WHEN:  DELETE /service-professional/{serviceId}/{professionalId}/ is called with invalid IDs
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	serviceId := uuid.New()
	professionalId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/"+serviceId.String()+"/"+professionalId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteServiceProfessionalInvalidIds(t *testing.T) {
	/*
		GIVEN: Invalid UUID formats for service and professional IDs
		WHEN:  DELETE /service-professional/{serviceId}/{professionalId}/ is called with invalid UUIDs
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidServiceId := "invalid-service-id"
	invalidProfessionalId := "invalid-professional-id"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/service-professional/"+invalidServiceId+"/"+invalidProfessionalId+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
