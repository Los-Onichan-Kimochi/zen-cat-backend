package user_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
	controllerTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/controller"
)

func TestBulkCreateUsersSuccessfully(t *testing.T) {
	// GIVEN: Valid bulk user creation request
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequests := []*schemas.CreateUserRequest{
		{
			Name:          "John",
			FirstLastName: "Doe",
			Password:      "password123",
			Email:         "john.bulk@example.com",
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/john.jpg",
		},
		{
			Name:          "Jane",
			FirstLastName: "Smith",
			Password:      "password456",
			Email:         "jane.bulk@example.com",
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/jane.jpg",
		},
		{
			Name:          "Admin",
			FirstLastName: "User",
			Password:      "adminpass789",
			Email:         "admin.bulk@example.com",
			Rol:           string(schemas.UserRolAdmin),
			ImageUrl:      "https://example.com/admin.jpg",
		},
	}

	// WHEN: BulkCreateUsers is called
	results, err := controller.BulkCreateUsers(createRequests, "test_admin")

	// THEN: All users are created successfully
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 3)

	// Verify each user
	assert.Equal(t, "John", results[0].Name)
	assert.Equal(t, "john.bulk@example.com", results[0].Email)
	assert.Equal(t, schemas.UserRolClient, results[0].Rol)

	assert.Equal(t, "Jane", results[1].Name)
	assert.Equal(t, "jane.bulk@example.com", results[1].Email)
	assert.Equal(t, schemas.UserRolClient, results[1].Rol)

	assert.Equal(t, "Admin", results[2].Name)
	assert.Equal(t, "admin.bulk@example.com", results[2].Email)
	assert.Equal(t, schemas.UserRolAdmin, results[2].Rol)
}

func TestBulkCreateUsersEmptyList(t *testing.T) {
	// GIVEN: Empty bulk creation request
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequests := []*schemas.CreateUserRequest{}

	// WHEN: BulkCreateUsers is called with empty list
	results, err := controller.BulkCreateUsers(createRequests, "test_admin")

	// THEN: An error is returned for empty list
	assert.NotNil(t, err)
	assert.Nil(t, results)
}

func TestBulkCreateUsersEmptyUpdatedBy(t *testing.T) {
	// GIVEN: Valid requests but empty updatedBy
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequests := []*schemas.CreateUserRequest{
		{
			Name:          "John",
			FirstLastName: "Doe",
			Password:      "password123",
			Email:         "john.empty@example.com",
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/john.jpg",
		},
	}

	// WHEN: BulkCreateUsers is called with empty updatedBy
	results, err := controller.BulkCreateUsers(createRequests, "")

	// THEN: An error is returned
	assert.NotNil(t, err)
	assert.Nil(t, results)
	assert.Equal(t, "Invalid updated by value", err.Message)
}

func TestBulkCreateUsersWithDuplicateEmails(t *testing.T) {
	// GIVEN: Requests with duplicate emails
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequests := []*schemas.CreateUserRequest{
		{
			Name:          "John",
			FirstLastName: "Doe",
			Password:      "password123",
			Email:         "duplicate.bulk@example.com",
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/john.jpg",
		},
		{
			Name:          "Jane",
			FirstLastName: "Smith",
			Password:      "password456",
			Email:         "duplicate.bulk@example.com", // Same email
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/jane.jpg",
		},
	}

	// WHEN: BulkCreateUsers is called with duplicate emails
	results, err := controller.BulkCreateUsers(createRequests, "test_admin")

	// THEN: Users are created successfully but with unique emails generated
	assert.Nil(t, err)
	assert.NotNil(t, results)
}

func TestBulkCreateUsersMixedRoles(t *testing.T) {
	// GIVEN: Users with different roles
	controller, _, _ := controllerTest.NewUserControllerTestWrapper(t)

	createRequests := []*schemas.CreateUserRequest{
		{
			Name:          "Client",
			FirstLastName: "User",
			Password:      "clientpass",
			Email:         "client.mixed@example.com",
			Rol:           string(schemas.UserRolClient),
			ImageUrl:      "https://example.com/client.jpg",
		},
		{
			Name:          "Admin",
			FirstLastName: "User",
			Password:      "adminpass",
			Email:         "admin.mixed@example.com",
			Rol:           string(schemas.UserRolAdmin),
			ImageUrl:      "https://example.com/admin.jpg",
		},
		{
			Name:          "Guest",
			FirstLastName: "User",
			Password:      "guestpass",
			Email:         "guest.mixed@example.com",
			Rol:           string(schemas.UserRolGuest),
			ImageUrl:      "https://example.com/guest.jpg",
		},
	}

	// WHEN: BulkCreateUsers is called
	results, err := controller.BulkCreateUsers(createRequests, "system_admin")

	// THEN: All users with different roles are created successfully
	assert.Nil(t, err)
	assert.NotNil(t, results)
	assert.Len(t, results, 3)

	// Verify roles
	assert.Equal(t, schemas.UserRolClient, results[0].Rol)
	assert.Equal(t, schemas.UserRolAdmin, results[1].Rol)
	assert.Equal(t, schemas.UserRolGuest, results[2].Rol)
}
