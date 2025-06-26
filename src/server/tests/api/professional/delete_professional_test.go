package professional_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestDeleteProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional exists in the database
		WHEN:  DELETE /professional/{professionalId}/ is called with a valid professional ID
		THEN:  A HTTP_204_NO_CONTENT status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test professional using factory
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/"+professional.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNoContent, rec.Code)
}

func TestDeleteProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: A non-existent professional ID
		WHEN:  DELETE /professional/{id}/ is called with non-existent ID
		THEN:  A HTTP_400_BAD_REQUEST status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/"+nonExistentId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestDeleteProfessionalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  DELETE /professional/{professionalId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidProfessionalId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodDelete, "/professional/"+invalidProfessionalId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
