package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestUpdateUserSuccessfully(t *testing.T) {
	// GIVEN: An existing user and valid update request
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	name := "John"
	firstName := "Doe"
	email := "john.original@example.com"
	rol := model.UserRolClient

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Rol:           &rol,
	})

	// Prepare update request
	newName := "Johnny"
	newEmail := "johnny.updated@example.com"
	newImageUrl := "https://example.com/new-avatar.jpg"

	updateRequest := schemas.UpdateUserRequest{
		Name:     &newName,
		Email:    &newEmail,
		ImageUrl: &newImageUrl,
	}

	// WHEN: UpdateUser is called
	result, err := controller.UpdateUser(testUser.Id, updateRequest, "test_admin")

	// THEN: User is updated successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Id, result.Id)
	assert.Equal(t, newName, result.Name)
	assert.Equal(t, newEmail, result.Email)
	assert.Equal(t, newImageUrl, result.ImageUrl)
	// Unchanged fields should remain the same
	assert.Equal(t, testUser.FirstLastName, result.FirstLastName)
}

func TestUpdateUserPartialFields(t *testing.T) {
	// GIVEN: An existing user and partial update request
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	name := "Jane"
	firstName := "Smith"
	email := "jane.smith@example.com"
	rol := model.UserRolClient

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:          &name,
		FirstLastName: &firstName,
		Email:         &email,
		Rol:           &rol,
	})

	// Prepare partial update request (only name)
	newName := "Janet"
	updateRequest := schemas.UpdateUserRequest{
		Name: &newName,
		// Other fields are nil, should not be updated
	}

	// WHEN: UpdateUser is called
	result, err := controller.UpdateUser(testUser.Id, updateRequest, "test_admin")

	// THEN: Only the specified field is updated
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Id, result.Id)
	assert.Equal(t, newName, result.Name)
	// Unchanged fields should remain the same
	assert.Equal(t, testUser.Email, result.Email)
	assert.Equal(t, testUser.FirstLastName, result.FirstLastName)
}

func TestUpdateUserEmptyUpdatedBy(t *testing.T) {
	// GIVEN: An existing user but empty updatedBy
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	testUser := factories.NewUserModel(db, factories.UserModelF{})

	updateRequest := schemas.UpdateUserRequest{
		Name: strPtr("New Name"),
	}

	// WHEN: UpdateUser is called with empty updatedBy
	result, err := controller.UpdateUser(testUser.Id, updateRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestUpdateUserNotFound(t *testing.T) {
	// GIVEN: A non-existent user ID
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	// Use a random UUID that doesn't exist
	nonExistentId := uuid.New()

	updateRequest := schemas.UpdateUserRequest{
		Name: strPtr("New Name"),
	}

	// WHEN: UpdateUser is called with non-existent ID
	result, err := controller.UpdateUser(nonExistentId, updateRequest, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Message, "not found")
}

func TestUpdateUserRole(t *testing.T) {
	// GIVEN: An existing user and role update request
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user with client role
	name := "John"
	email := "john.role@example.com"
	rol := model.UserRolClient

	testUser := factories.NewUserModel(db, factories.UserModelF{
		Name:  &name,
		Email: &email,
		Rol:   &rol,
	})

	// Prepare role update request
	newRole := string(schemas.UserRolAdmin)
	updateRequest := schemas.UpdateUserRequest{
		Rol: &newRole,
	}

	// WHEN: UpdateUser is called
	result, err := controller.UpdateUser(testUser.Id, updateRequest, "system_admin")

	// THEN: User role is updated successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Id, result.Id)
	assert.Equal(t, schemas.UserRolAdmin, result.Rol)
	// Other fields should remain unchanged
	assert.Equal(t, testUser.Name, result.Name)
	assert.Equal(t, testUser.Email, result.Email)
}

func TestUpdateUserPassword(t *testing.T) {
	// GIVEN: An existing user and password update request
	controller, _, db := controllerTest.NewUserControllerTestWrapper(t)

	// Create a test user
	testUser := factories.NewUserModel(db, factories.UserModelF{})
	originalPassword := testUser.Password

	// Prepare password update request
	newPassword := "newSecurePassword123"
	updateRequest := schemas.UpdateUserRequest{
		Password: &newPassword,
	}

	// WHEN: UpdateUser is called
	result, err := controller.UpdateUser(testUser.Id, updateRequest, "test_admin")

	// THEN: User password is updated successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, testUser.Id, result.Id)
	// Password should be different from original
	assert.NotEqual(t, originalPassword, result.Password)
	// In current implementation, password might not be hashed - just ensure it's updated
}
