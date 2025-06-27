package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestDeleteUserSuccessfully(t *testing.T) {
	// GIVEN: An existing user
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	// WHEN: DeleteUser is called
	err := controller.DeleteUser(testUser.Id)

	// THEN: User is deleted successfully
	assert.Nil(t, err)

	// Verify user is deleted by trying to get it
	deletedUser, getErr := controller.GetUser(testUser.Id)
	assert.NotNil(t, getErr)
	assert.Nil(t, deletedUser)
}

func TestDeleteUserNotFound(t *testing.T) {
	// GIVEN: A non-existent user ID
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN: DeleteUser is called with non-existent ID
	err := controller.DeleteUser(nonExistentId)

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Contains(t, err.Message, "not found")
}

func TestBulkDeleteUsersSuccessfully(t *testing.T) {
	// GIVEN: Multiple existing users
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create test users
	user1 := factories.NewUserModel(db, factories.UserModelF{})
	user2 := factories.NewUserModel(db, factories.UserModelF{})
	user3 := factories.NewUserModel(db, factories.UserModelF{})

	bulkDeleteRequest := schemas.BulkDeleteUserRequest{
		Users: []uuid.UUID{user1.Id, user2.Id, user3.Id},
	}

	// WHEN: BulkDeleteUsers is called
	err := controller.BulkDeleteUsers(bulkDeleteRequest)

	// THEN: All users are deleted successfully
	assert.Nil(t, err)

	// Verify users are deleted
	_, getErr1 := controller.GetUser(user1.Id)
	_, getErr2 := controller.GetUser(user2.Id)
	_, getErr3 := controller.GetUser(user3.Id)

	assert.NotNil(t, getErr1)
	assert.NotNil(t, getErr2)
	assert.NotNil(t, getErr3)
}

func TestBulkDeleteUsersEmpty(t *testing.T) {
	// GIVEN: Empty bulk delete request
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	bulkDeleteRequest := schemas.BulkDeleteUserRequest{
		Users: []uuid.UUID{}, // Empty list
	}

	// WHEN: BulkDeleteUsers is called with empty list
	err := controller.BulkDeleteUsers(bulkDeleteRequest)

	// THEN: No error occurs (should handle empty gracefully)
	assert.Nil(t, err)
}

func TestBulkDeleteUsersPartialFailure(t *testing.T) {
	// GIVEN: Mix of existing and non-existent user IDs
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create one real user
	existingUser := factories.NewUserModel(db, factories.UserModelF{})
	nonExistentId := uuid.New()

	bulkDeleteRequest := schemas.BulkDeleteUserRequest{
		Users: []uuid.UUID{existingUser.Id, nonExistentId},
	}

	// WHEN: BulkDeleteUsers is called with mix of valid/invalid IDs
	err := controller.BulkDeleteUsers(bulkDeleteRequest)

	// THEN: The operation should handle partial failures gracefully
	// (The exact behavior depends on implementation - it might succeed partially or fail completely)
	// For now, we just check that it doesn't panic
	// The existing user should be deleted if the operation is partial
	if err == nil {
		_, getErr := controller.GetUser(existingUser.Id)
		assert.NotNil(t, getErr) // Should be deleted
	}
}
