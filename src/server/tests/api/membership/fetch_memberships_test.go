package membership_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestFetchMembershipsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple memberships exist in the database
		WHEN:  GET /membership/ is called
		THEN:  A list of memberships is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)
	factories.NewMembershipModel(db, factories.MembershipModelF{})
	factories.NewMembershipModel(db, factories.MembershipModelF{})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Memberships
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.GreaterOrEqual(t, len(response.Memberships), 2)
}

func TestFetchMembershipsEmpty(t *testing.T) {
	/*
		GIVEN: No memberships exist in the database
		WHEN:  GET /membership/ is called
		THEN:  An empty list is returned with a HTTP_200_OK status
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/membership/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.Memberships
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)
	assert.Empty(t, response.Memberships)
}
