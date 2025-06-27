package onboarding_test

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

func TestFetchOnboardingsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple onboarding records exist
		WHEN:  GET /onboarding/ is called
		THEN:  A HTTP_200_OK status should be returned with all onboardings
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create users first
	user1 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "John",
		FirstLastName: "Doe",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	user2 := &model.User{
		Email:         utilsTest.GenerateRandomEmail(),
		Name:          "Jane",
		FirstLastName: "Smith",
		Rol:           "MEMBER",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err := db.Create([]*model.User{user1, user2}).Error
	assert.NoError(t, err)

	// Create onboarding records
	onboarding1 := &model.Onboarding{
		UserId:         user1.Id,
		DocumentType:   "DNI",
		DocumentNumber: "12345678",
		PhoneNumber:    "987654321",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	onboarding2 := &model.Onboarding{
		UserId:         user2.Id,
		DocumentType:   "PASSPORT",
		DocumentNumber: "87654321",
		PhoneNumber:    "123456789",
		AuditFields: model.AuditFields{
			UpdatedBy: "ADMIN",
		},
	}
	err = db.Create([]*model.Onboarding{onboarding1, onboarding2}).Error
	assert.NoError(t, err)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/onboarding/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Onboardings
	err = json.NewDecoder(rec.Body).Decode(&response)
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
