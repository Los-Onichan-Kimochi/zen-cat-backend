package professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestCreateProfessionalSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid professional data
		WHEN:  POST /professional/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the created professional
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	professionalRequest := schemas.CreateProfessionalRequest{
		Name:           "Dr. John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Specialty:      "Cardiology",
		Email:          utilsTest.GenerateRandomEmail(),
		PhoneNumber:    "+1234567890",
		Type:           "MEDIC",
		ImageUrl:       "professional-image.jpg",
	}

	requestBody, _ := json.Marshal(professionalRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Professional
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains the correct data
	assert.NotEmpty(t, response.Id)
	assert.Equal(t, professionalRequest.Name, response.Name)
	assert.Equal(t, professionalRequest.FirstLastName, response.FirstLastName)
	assert.Equal(t, &professionalRequest.SecondLastName, response.SecondLastName)
	assert.Equal(t, professionalRequest.Specialty, response.Specialty)
	assert.Equal(t, professionalRequest.Email, response.Email)
	assert.Equal(t, professionalRequest.PhoneNumber, response.PhoneNumber)
	assert.Equal(t, professionalRequest.Type, response.Type)
	assert.True(t, strings.HasPrefix(response.ImageUrl, professionalRequest.ImageUrl))
}

func TestCreateProfessionalInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /professional/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(
		http.MethodPost,
		"/professional/",
		strings.NewReader(`{"invalid": json`),
	)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestCreateProfessionalMissingRequiredFields(t *testing.T) {
	/*
		GIVEN: A request body missing required fields
		WHEN:  POST /professional/ is called with incomplete data
		THEN:  A HTTP_400_BAD_REQUEST status should be returned due to validation
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// Missing required name field - should trigger validation error
	professionalRequest := schemas.CreateProfessionalRequest{
		Specialty: "Cardiology",
	}

	requestBody, _ := json.Marshal(professionalRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func TestCreateProfessionalDuplicateEmail(t *testing.T) {
	/*
		GIVEN: A professional with a duplicate email
		WHEN:  POST /professional/ is called with an existing email
		THEN:  A HTTP_201_CREATED status should be returned (API doesn't prevent duplicates)
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create an existing professional
	existingProfessional := factories.NewProfessionalModel(db, factories.ProfessionalModelF{})

	// Try to create another professional with the same email
	duplicateProfessionalRequest := schemas.CreateProfessionalRequest{
		Name:           "Dr. Jane",
		FirstLastName:  "Smith",
		SecondLastName: "Johnson",
		Specialty:      "Neurology",
		Email:          existingProfessional.Email, // Same email
		PhoneNumber:    "+0987654321",
		Type:           "MEDIC",
		ImageUrl:       "duplicate-professional.jpg",
	}

	requestBody, _ := json.Marshal(duplicateProfessionalRequest)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)
}
