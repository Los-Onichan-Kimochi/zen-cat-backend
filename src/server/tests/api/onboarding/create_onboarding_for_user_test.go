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
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestCreateOnboardingForUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A valid request to create an onboarding for a user
		WHEN:  POST /onboarding/user/{userId}/ is called
		THEN:  A new onboarding is created and a HTTP_201_CREATED status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})

	gender := schemas.GenderMale
	birthDate := time.Now().UTC().Truncate(24 * time.Hour)
	district := "district"
	province := "province"
	region := "region"

	request := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "123456789",
		PostalCode:     "12345",
		Address:        "123 Main St",
		Gender:         &gender,
		BirthDate:      &birthDate,
		District:       &district,
		Province:       &province,
		Region:         &region,
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/onboarding/user/"+user.Id.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Onboarding
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Id)
	assert.Equal(t, user.Id, response.UserId)
	assert.Equal(t, request.DocumentType, response.DocumentType)
	assert.Equal(t, request.DocumentNumber, response.DocumentNumber)
	assert.Equal(t, request.PhoneNumber, response.PhoneNumber)
	assert.Equal(t, request.PostalCode, response.PostalCode)
	assert.Equal(t, request.Address, response.Address)
	assert.Equal(t, *request.Gender, *response.Gender)
	assert.Equal(t, request.BirthDate.Format("2006-01-02"), response.BirthDate.Format("2006-01-02"))
	assert.Equal(t, *request.District, *response.District)
	assert.Equal(t, *request.Province, *response.Province)
	assert.Equal(t, *request.Region, *response.Region)
}

func TestCreateOnboardingForUserInvalidBody(t *testing.T) {
	/*
		GIVEN: An invalid request to create an onboarding for a user
		WHEN:  POST /onboarding/user/{userId}/ is called
		THEN:  A HTTP_400_BAD_REQUEST status is returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	user := factories.NewUserModel(db, factories.UserModelF{})

	request := schemas.CreateOnboardingRequest{}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/onboarding/user/"+user.Id.String()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateOnboardingForUserInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid user ID is provided
		WHEN:  POST /onboarding/user/{userId}/ is called
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	request := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "123456789",
		PostalCode:     "12345",
		Address:        "123 Main St",
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/onboarding/user/invalid-uuid/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateOnboardingForUserNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent user ID is provided
		WHEN:  POST /onboarding/user/{userId}/ is called
		THEN:  A HTTP_404_NOT_FOUND status is returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	request := schemas.CreateOnboardingRequest{
		DocumentType:   schemas.DocumentTypeDNI,
		DocumentNumber: "12345678",
		PhoneNumber:    "123456789",
		PostalCode:     "12345",
		Address:        "123 Main St",
	}
	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/onboarding/user/"+uuid.NewString()+"/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}
