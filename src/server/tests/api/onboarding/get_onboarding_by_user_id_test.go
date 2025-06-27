package onboarding_test

import (
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

func TestGetOnboardingByUserIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing onboarding record for a user
		WHEN:  GET /onboarding/user/{userId}/ is called with valid user ID
		THEN:  A HTTP_200_OK status should be returned with the onboarding data
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

	// Create an onboarding record with all required fields
	district := "Test District"
	province := "Test Province"
	region := "Test Region"
	birthDate := time.Now().AddDate(-30, 0, 0) // 30 years ago
	gender := model.GenderMale

	onboarding := &model.Onboarding{
		UserId:         user.Id,
		DocumentType:   model.DocumentTypeDni,
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		BirthDate:      &birthDate,
		Gender:         &gender,
		PostalCode:     "12345",
		Address:        "123 Test St",
		District:       &district,
		Province:       &province,
		Region:         &region,
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create(onboarding).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/user/"+user.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Onboarding
	err = json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, onboarding.Id, response.Id)
	assert.Equal(t, onboarding.UserId, response.UserId)
	assert.Equal(t, schemas.DocumentType(onboarding.DocumentType), response.DocumentType)
	assert.Equal(t, onboarding.DocumentNumber, response.DocumentNumber)
	assert.Equal(t, onboarding.PhoneNumber, response.PhoneNumber)
}

func TestGetOnboardingByUserIdNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent user ID
		WHEN:  GET /onboarding/user/{userId}/ is called with invalid user ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/user/"+nonExistentUserId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetOnboardingByUserIdInvalidId(t *testing.T) {
	/*
		GIVEN: An invalid UUID format for user ID
		WHEN:  GET /onboarding/user/{userId}/ is called with invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/user/"+invalidId+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
