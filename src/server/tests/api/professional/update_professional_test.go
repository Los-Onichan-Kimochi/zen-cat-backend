package professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestUpdateProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: A professional exists in the database
		WHEN:  PATCH /professional/{professionalId}/ is called with valid update data
		THEN:  A HTTP_200_OK status should be returned with the updated professional
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create a test professional using factory
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Prepare update data
	newName := "Dr. Updated"
	newSpecialty := "Updated Specialty"
	updateRequest := schemas.UpdateProfessionalRequest{
		Name:      &newName,
		Specialty: &newSpecialty,
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/professional/"+professional.Id.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Professional
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the updated data
	assert.Equal(t, professional.Id, response.Id)
	assert.Equal(t, newName, response.Name)
	assert.Equal(t, newSpecialty, response.Specialty)
	// Other fields should remain unchanged
	assert.Equal(t, professional.Email, response.Email)
}

func TestUpdateProfessionalNotFound(t *testing.T) {
	/*
		GIVEN: No professional exists with the provided ID
		WHEN:  PATCH /professional/{professionalId}/ is called with a non-existent professional ID
		THEN:  A HTTP_404_NOT_FOUND status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	nonExistentProfessionalId := uuid.New()

	updateRequest := schemas.UpdateProfessionalRequest{
		Name: func() *string { s := "Dr. Updated"; return &s }(),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/professional/"+nonExistentProfessionalId.String()+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusNotFound, rec.Code)
}

func TestUpdateProfessionalInvalidUUID(t *testing.T) {
	/*
		GIVEN: An invalid UUID format is provided
		WHEN:  PATCH /professional/{professionalId}/ is called with an invalid UUID
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)
	invalidProfessionalId := "invalid-uuid"

	updateRequest := schemas.UpdateProfessionalRequest{
		Name: func() *string { s := "Dr. Updated"; return &s }(),
	}

	requestBody, _ := json.Marshal(updateRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/professional/"+invalidProfessionalId+"/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestUpdateProfessionalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  PATCH /professional/{professionalId}/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	professional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodPatch, "/professional/"+professional.Id.String()+"/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}
