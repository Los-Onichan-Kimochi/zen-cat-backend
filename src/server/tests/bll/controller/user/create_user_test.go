package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestCreateUserSuccessfully(t *testing.T) {
	// GIVEN: Valid user creation request
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequest := schemas.CreateUserRequest{
		Name:           "John",
		FirstLastName:  "Doe",
		SecondLastName: "Smith",
		Password:       "securePassword123",
		Email:          "john.doe@example.com",
		Rol:            string(schemas.UserRolClient),
		ImageUrl:       "https://example.com/avatar.jpg",
	}

	// WHEN: CreateUser is called
	result, err := controller.CreateUser(createRequest, "test_admin")

	// THEN: User is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.FirstLastName, result.FirstLastName)
	assert.Equal(t, createRequest.SecondLastName, *result.SecondLastName)
	assert.Equal(t, createRequest.Email, result.Email)
	assert.Equal(t, createRequest.Rol, result.Rol)
	assert.Equal(t, createRequest.ImageUrl, result.ImageUrl)
	assert.NotEqual(t, "", result.Id)
}

func TestCreateUserWithoutSecondLastName(t *testing.T) {
	// GIVEN: Valid user creation request without second last name
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequest := schemas.CreateUserRequest{
		Name:           "Jane",
		FirstLastName:  "Doe",
		SecondLastName: "", // Empty second last name
		Password:       "securePassword123",
		Email:          "jane.doe@example.com",
		Rol:            string(schemas.UserRolClient),
		ImageUrl:       "https://example.com/jane.jpg",
	}

	// WHEN: CreateUser is called
	result, err := controller.CreateUser(createRequest, "test_admin")

	// THEN: User is created successfully without second last name
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.FirstLastName, result.FirstLastName)
	assert.Nil(t, result.SecondLastName)
	assert.Equal(t, createRequest.Email, result.Email)
	assert.Equal(t, createRequest.Rol, result.Rol)
}

func TestCreateUserEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid user creation request but empty updatedBy
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequest := schemas.CreateUserRequest{
		Name:          "John",
		FirstLastName: "Doe",
		Password:      "securePassword123",
		Email:         "john.test@example.com",
		Rol:           string(schemas.UserRolClient),
		ImageUrl:      "https://example.com/avatar.jpg",
	}

	// WHEN: CreateUser is called with empty updatedBy
	result, err := controller.CreateUser(createRequest, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	// GIVEN: A user already exists with the same email
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequest1 := schemas.CreateUserRequest{
		Name:          "John",
		FirstLastName: "Doe",
		Password:      "securePassword123",
		Email:         "duplicate@example.com",
		Rol:           string(schemas.UserRolClient),
		ImageUrl:      "https://example.com/avatar1.jpg",
	}

	createRequest2 := schemas.CreateUserRequest{
		Name:          "Jane",
		FirstLastName: "Smith",
		Password:      "securePassword456",
		Email:         "duplicate@example.com", // Same email
		Rol:           string(schemas.UserRolClient),
		ImageUrl:      "https://example.com/avatar2.jpg",
	}

	// Create first user
	_, err1 := controller.CreateUser(createRequest1, "test_admin")
	assert.Nil(t, err1)

	// WHEN: CreateUser is called with duplicate email
	result, err := controller.CreateUser(createRequest2, "test_admin")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestCreateUserWithAdminRole(t *testing.T) {
	// GIVEN: Valid user creation request with admin role
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequest := schemas.CreateUserRequest{
		Name:          "Admin",
		FirstLastName: "User",
		Password:      "adminPassword123",
		Email:         "admin@example.com",
		Rol:           string(schemas.UserRolAdmin),
		ImageUrl:      "https://example.com/admin.jpg",
	}

	// WHEN: CreateUser is called
	result, err := controller.CreateUser(createRequest, "system_admin")

	// THEN: Admin user is created successfully
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, createRequest.Name, result.Name)
	assert.Equal(t, createRequest.Email, result.Email)
	assert.Equal(t, schemas.UserRolAdmin, result.Rol)
	assert.NotEqual(t, "", result.Id)
}
