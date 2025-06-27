package professional_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
	utilsTest "onichankimochi.com/astro_cat_backend/src/server/tests/utils"
)

func TestBulkCreateProfessionalsSuccessfully(t *testing.T) {
	/*
		GIVEN: Valid professionals data
		WHEN:  POST /professional/bulk-create/ is called with valid data
		THEN:  A HTTP_201_CREATED status should be returned with the created professionals
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	request := schemas.BulkCreateProfessionalRequest{
		Professionals: []*schemas.CreateProfessionalRequest{
			{
				Name:           "Dr. John",
				FirstLastName:  "Doe",
				SecondLastName: "Smith",
				Specialty:      "Cardiology",
				Email:          utilsTest.GenerateRandomEmail(),
				PhoneNumber:    "+1234567890",
				Type:           "MEDIC",
				ImageUrl:       "professional1.jpg",
			},
			{
				Name:           "Dr. Jane",
				FirstLastName:  "Smith",
				SecondLastName: "Johnson",
				Specialty:      "Neurology",
				Email:          utilsTest.GenerateRandomEmail(),
				PhoneNumber:    "+0987654321",
				Type:           "MEDIC",
				ImageUrl:       "professional2.jpg",
			},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Professionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Len(t, response.Professionals, 2)

	// Verify the created professionals
	for i, professional := range response.Professionals {
		assert.NotEmpty(t, professional.Id)
		assert.Equal(t, request.Professionals[i].Name, professional.Name)
		assert.Equal(t, request.Professionals[i].Email, professional.Email)
		assert.Equal(t, request.Professionals[i].Specialty, professional.Specialty)
	}
}

func TestBulkCreateProfessionalsEmptyList(t *testing.T) {
	/*
		GIVEN: An empty professionals list
		WHEN:  POST /professional/bulk-create/ is called with empty list
		THEN:  A HTTP_201_CREATED status should be returned with empty response
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	request := schemas.BulkCreateProfessionalRequest{
		Professionals: []*schemas.CreateProfessionalRequest{},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)

	var response schemas.Professionals
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Professionals)
}

func TestBulkCreateProfessionalsInvalidRequestBody(t *testing.T) {
	/*
		GIVEN: An invalid request body
		WHEN:  POST /professional/bulk-create/ is called with invalid JSON
		THEN:  A HTTP_422_UNPROCESSABLE_ENTITY status should be returned
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/bulk-create/", strings.NewReader(`{"invalid": json`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)
}

func TestBulkCreateProfessionalsPartialFailure(t *testing.T) {
	/*
		GIVEN: A mix of valid and invalid professional data
		WHEN:  POST /professional/bulk-create/ is called with some invalid data
		THEN:  A HTTP_201_CREATED status should be returned (API processes what it can)
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	request := schemas.BulkCreateProfessionalRequest{
		Professionals: []*schemas.CreateProfessionalRequest{
			{
				Name:           "Dr. John",
				FirstLastName:  "Doe",
				SecondLastName: "Smith",
				Specialty:      "Cardiology",
				Email:          utilsTest.GenerateRandomEmail(),
				PhoneNumber:    "+1234567890",
				Type:           "MEDIC",
				ImageUrl:       "professional1.jpg",
			},
			{
				// Missing some fields but API will still process it
				Name:      "Dr. Jane",
				Specialty: "Neurology",
			},
		},
	}

	body, _ := json.Marshal(request)

	// WHEN
	req := httptest.NewRequest(http.MethodPost, "/professional/bulk-create/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusCreated, rec.Code)
}
