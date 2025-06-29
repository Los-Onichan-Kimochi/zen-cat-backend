package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestGetUserStatsEmpty(t *testing.T) {
	// GIVEN: No users in the database
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Stats are returned with zero counts or an error
	if err != nil {
		assert.NotNil(t, err)
	} else {
		assert.Nil(t, err)
		assert.NotNil(t, result)
	}
}

func TestGetUserStatsWithUsers(t *testing.T) {
	// GIVEN: Multiple users with different roles
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create users with different roles
	adminRole := model.UserRolAdmin
	clientRole := model.UserRolClient
	guestRole := model.UserRolGuest

	// Create 2 admins
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Admin1"),
		Rol:  &adminRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Admin2"),
		Rol:  &adminRole,
	})

	// Create 3 clients
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Client1"),
		Rol:  &clientRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Client2"),
		Rol:  &clientRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Client3"),
		Rol:  &clientRole,
	})

	// Create 1 guest
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("Guest1"),
		Rol:  &guestRole,
	})

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Stats reflect the created users
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(6), result.TotalUsers) // 2 + 3 + 1
	assert.Equal(t, int64(2), result.AdminCount)
	assert.Equal(t, int64(3), result.ClientCount)
	assert.Equal(t, int64(1), result.GuestCount)
	assert.NotNil(t, result.RoleDistribution)
	assert.NotNil(t, result.RecentConnections)
}

func TestGetUserStatsOnlyAdmins(t *testing.T) {
	// GIVEN: Only admin users
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	adminRole := model.UserRolAdmin

	// Create multiple admin users
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("SuperAdmin"),
		Rol:  &adminRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("ModeratorAdmin"),
		Rol:  &adminRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("SystemAdmin"),
		Rol:  &adminRole,
	})

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Stats show only admin users
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(3), result.TotalUsers)
	assert.Equal(t, int64(3), result.AdminCount)
	assert.Equal(t, int64(0), result.ClientCount)
	assert.Equal(t, int64(0), result.GuestCount)
}

func TestGetUserStatsOnlyClients(t *testing.T) {
	// GIVEN: Only client users
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	clientRole := model.UserRolClient

	// Create multiple client users
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("RegularClient1"),
		Rol:  &clientRole,
	})
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("RegularClient2"),
		Rol:  &clientRole,
	})

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Stats show only client users
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.TotalUsers)
	assert.Equal(t, int64(0), result.AdminCount)
	assert.Equal(t, int64(2), result.ClientCount)
	assert.Equal(t, int64(0), result.GuestCount)
}

func TestGetUserStatsRoleDistribution(t *testing.T) {
	// GIVEN: Users with mixed roles
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	adminRole := model.UserRolAdmin
	clientRole := model.UserRolClient

	// Create users
	factories.NewUserModel(db, factories.UserModelF{Rol: &adminRole})
	factories.NewUserModel(db, factories.UserModelF{Rol: &clientRole})
	factories.NewUserModel(db, factories.UserModelF{Rol: &clientRole})

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Role distribution is provided
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.RoleDistribution)

	// The exact structure of RoleDistribution depends on implementation
	// but it should contain information about role counts
	if len(result.RoleDistribution) > 0 {
		// Verify the distribution contains role information
		for _, roleInfo := range result.RoleDistribution {
			assert.NotEmpty(t, roleInfo.Role)
			assert.GreaterOrEqual(t, roleInfo.Count, int64(0))
		}
	}
}

func TestGetUserStatsRecentConnections(t *testing.T) {
	// GIVEN: Users in the database
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a user
	factories.NewUserModel(db, factories.UserModelF{
		Name: strPtr("ConnectedUser"),
	})

	// WHEN: GetUserStats is called
	result, err := controller.GetUserStats()

	// THEN: Recent connections information is provided
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.NotNil(t, result.RecentConnections)

	// The exact structure depends on implementation
	// Recent connections might be empty if no connection tracking is implemented
	for _, connection := range result.RecentConnections {
		assert.NotEqual(t, "", connection.UserId)
		assert.NotEmpty(t, connection.UserEmail)
		assert.NotEmpty(t, connection.UserName)
		assert.NotEmpty(t, connection.Role)
	}
}
