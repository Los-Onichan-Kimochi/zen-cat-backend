package onboarding_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestFetchOnboardingsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple onboarding records exist in the system
		WHEN:  GET /onboarding/ is called
		THEN:  A HTTP_200_OK status should be returned with list of onboarding records
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create onboarding records using factories (which automatically create users)
	email1 := utilsTest.GenerateRandomEmail()
	email2 := utilsTest.GenerateRandomEmail()

	// Create first onboarding with unique user
	docNum1 := "12345678"
	phone1 := "987654321"
	onboarding1 := factories.NewOnboardingModel(db, factories.OnboardingModelF{
		DocumentNumber: &docNum1,
		PhoneNumber:    &phone1,
	})

	// Update the user email to be unique
	db.Model(&onboarding1.User).Update("email", email1)

	// Create second onboarding with unique user
	docNum2 := "87654321"
	phone2 := "123456789"
	onboarding2 := factories.NewOnboardingModel(db, factories.OnboardingModelF{
		DocumentNumber: &docNum2,
		PhoneNumber:    &phone2,
	})

	// Update the user email to be unique
	db.Model(&onboarding2.User).Update("email", email2)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Onboardings
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Onboardings), 2)
}

func TestFetchOnboardingsEmpty(t *testing.T) {
	/*
		GIVEN: No onboarding records exist
		WHEN:  GET /onboarding/ is called
		THEN:  A HTTP_200_OK status should be returned with empty array
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Onboardings
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(response.Onboardings))
}
