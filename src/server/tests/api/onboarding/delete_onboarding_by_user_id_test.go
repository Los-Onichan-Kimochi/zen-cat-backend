package onboarding_test

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

func TestDeleteOnboardingByUserIdSuccessfully(t *testing.T) {
	/*
		GIVEN: An existing onboarding record for a user
		WHEN:  DELETE /onboarding/user/{userId}/ is called with valid user ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned and onboarding should be deleted
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

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/onboarding/user/"+user.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)

	// Verify the onboarding was deleted
	var count int64
	db.Model(&model.Onboarding{}).Where("user_id = ?", user.Id).Count(&count)
	assert.Equal(t, int64(0), count)
}

func TestDeleteOnboardingByUserIdNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent user ID
		WHEN:  DELETE /onboarding/user/{userId}/ is called with invalid user ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentUserId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/onboarding/user/"+nonExistentUserId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestDeleteOnboardingByUserIdInvalidUserId(t *testing.T) {
	/*
		GIVEN: An invalid UUID format for user ID
		WHEN:  DELETE /onboarding/user/{userId}/ is called with invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidUserId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/onboarding/user/"+invalidUserId+"/", nil)
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
