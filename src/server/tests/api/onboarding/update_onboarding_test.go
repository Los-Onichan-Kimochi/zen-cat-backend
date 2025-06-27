package onboarding_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func documentTypePtr(dt schemas.DocumentType) *schemas.DocumentType {
	return &dt
}

func genderPtr(g schemas.Gender) *schemas.Gender {
	return &g
}

func timePtr(s string) *time.Time {
	t, _ := time.Parse("2006-01-02", s)
	return &t
}

func TestUpdateOnboardingSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing onboarding record
		WHEN:  PATCH /onboarding/{onboardingId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with updated onboarding
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a user first
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	// Create an onboarding record
	onboarding := &model.Onboarding{
		UserId:         user.Id,
		DocumentType:   "DNI",
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(onboarding).Error
	assert.NoError(t, err)

	// Create update request
	updateRequest := schemas.UpdateOnboardingRequest{
		DocumentType:   documentTypePtr(schemas.DocumentTypePassport),
		DocumentNumber: stringPtr("87654321"),
		PhoneNumber:    stringPtr("123456789"),
		BirthDate:      timePtr("1990-01-01"),
		Gender:         genderPtr(schemas.GenderMale),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/onboarding/"+onboarding.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Onboarding
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, onboarding.Id, response.Id)
	assert.Equal(t, *updateRequest.DocumentType, response.DocumentType)
	assert.Equal(t, *updateRequest.DocumentNumber, response.DocumentNumber)
	assert.Equal(t, *updateRequest.PhoneNumber, response.PhoneNumber)
}

func TestUpdateOnboardingNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent onboarding ID
		WHEN:  PATCH /onboarding/{onboardingId}/ is called with invalid ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	updateRequest := schemas.UpdateOnboardingRequest{
		DocumentType:   documentTypePtr(schemas.DocumentTypePassport),
		DocumentNumber: stringPtr("87654321"),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/onboarding/"+nonExistentId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateOnboardingInvalidId(t *testing.T) {
	/*
		GIVEN: An invalid UUID format for onboarding ID
		WHEN:  PATCH /onboarding/{onboardingId}/ is called with invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidId := "invalid-uuid"

	updateRequest := schemas.UpdateOnboardingRequest{
		DocumentType: documentTypePtr(schemas.DocumentTypePassport),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/onboarding/"+invalidId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateOnboardingInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An existing onboarding and invalid request body
		WHEN:  PATCH /onboarding/{onboardingId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a user first
	user := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create(user).Error
	assert.NoError(t, err)

	// Create an onboarding record
	onboarding := &model.Onboarding{
		UserId:         user.Id,
		DocumentType:   "DNI",
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(onboarding).Error
	assert.NoError(t, err)

	invalidJSON := `{"invalid": json}`

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/onboarding/"+onboarding.Id.String()+"/", bytes.NewBufferString(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
