package professional_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional exists in the database
		WHEN:  GET /professional/{professionalId}/ is called with a valid professional ID
		THEN:  A HTTP_200_OK status should be returned with the professional data
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test professional using factory
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/professional/"+professional.Id.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Professional
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct professional data
	assert.Equal(t, professional.Id, response.Id)
	assert.Equal(t, professional.Name, response.Name)
	assert.Equal(t, professional.FirstLastName, response.FirstLastName)
	assert.Equal(t, professional.SecondLastName, response.SecondLastName)
	assert.Equal(t, professional.Email, response.Email)
	assert.Equal(t, professional.PhoneNumber, response.PhoneNumber)
	assert.Equal(t, professional.Specialty, response.Specialty)
	assert.Equal(t, string(professional.Type), response.Type)
	assert.Equal(t, professional.ImageUrl, response.ImageUrl)
}

func TestGetProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No professional exists with the provided ID
		WHEN:  GET /professional/{professionalId}/ is called with a non-existent professional ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentProfessionalId := uuid.New()

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/professional/"+nonExistentProfessionalId.String()+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestGetProfessionalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  GET /professional/{professionalId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidProfessionalId := "invalid-uuid"

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/professional/"+invalidProfessionalId+"/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
