package user_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	apiTest "onichankimochi.com/astro_cat_backend/src/server/tests/api"
)

func TestGetUserStatsSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple users exist in the database with different roles
		WHEN:  GET /user/stats/ is called
		THEN:  A HTTP_200_OK status should be returned with user statistics
	*/
	// GIVEN
	server, db := apiTest.NewApiServerTestWrapper(t)

	// Create test users with different roles
	adminRole := model.UserRolAdmin
	clientRole := model.UserRolClient
	guestRole := model.UserRolGuest

	// Create 2 admin users
	factories.NewUserModelBatch(db, 2, factories.UserModelF{
		Rol: &adminRole,
	})

	// Create 3 client users
	factories.NewUserModelBatch(db, 3, factories.UserModelF{
		Rol: &clientRole,
	})

	// Create 1 guest user
	factories.NewUserModelBatch(db, 1, factories.UserModelF{
		Rol: &guestRole,
	})

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/stats/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.UserStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains correct statistics
	assert.Equal(t, int64(6), response.TotalUsers) // 2 admin + 3 client + 1 guest
	assert.Equal(t, int64(2), response.AdminCount)
	assert.Equal(t, int64(3), response.ClientCount)
	assert.Equal(t, int64(1), response.GuestCount)

	// Verify role distribution
	assert.Len(t, response.RoleDistribution, 3) // 3 different roles

	// Check that recent connections array exists (might be empty)
	assert.NotNil(t, response.RecentConnections)
}

func TestGetUserStatsNoUsers(t *testing.T) {
	/*
		GIVEN: No users exist in the database
		WHEN:  GET /user/stats/ is called
		THEN:  A HTTP_200_OK status should be returned with empty statistics
	*/
	// GIVEN
	server, _ := apiTest.NewApiServerTestWrapper(t)

	// WHEN
	req := httptest.NewRequest(http.MethodGet, "/user/stats/", nil)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	rec := httptest.NewRecorder()
	server.Echo.ServeHTTP(rec, req)

	// THEN
	assert.Equal(t, http.StatusOK, rec.Code)

	var response schemas.UserStats
	err := json.NewDecoder(rec.Body).Decode(&response)
	assert.NoError(t, err)

	// Verify the response contains empty statistics
	assert.Equal(t, int64(0), response.TotalUsers)
	assert.Equal(t, int64(0), response.AdminCount)
	assert.Equal(t, int64(0), response.ClientCount)
	assert.Equal(t, int64(0), response.GuestCount)

	// Verify empty role distribution and recent connections
	assert.Empty(t, response.RoleDistribution)
	assert.Empty(t, response.RecentConnections)
}
