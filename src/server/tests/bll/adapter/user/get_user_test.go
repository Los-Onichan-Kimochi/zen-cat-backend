package user_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"onichankimochi.com/astro_cat_backend/src/server/dao/factories"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	adapterTest "onichankimochi.com/astro_cat_backend/src/server/tests/bll/adapter"
)

func TestGetUserSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists in the database
		WHEN:  GetPostgresqlUser is called with the user ID
		THEN:  The user is returned with all related data
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewUserAdapterTestWrapper(t)

	existingUser := factories.NewUserModel(db, factories.UserModelF{})

	// WHEN
	user, err := adapter.GetPostgresqlUser(existingUser.Id)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.Id, user.Id)
	assert.Equal(t, existingUser.Name, user.Name)
	assert.Equal(t, existingUser.FirstLastName, user.FirstLastName)
	assert.Equal(t, existingUser.SecondLastName, user.SecondLastName)
	assert.Equal(t, existingUser.Email, user.Email)
	assert.Equal(t, string(existingUser.Rol), string(user.Rol))
	assert.Equal(t, existingUser.ImageUrl, user.ImageUrl)
}

func TestGetUserNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the given ID
		WHEN:  GetPostgresqlUser is called with non-existent ID
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewUserAdapterTestWrapper(t)

	nonExistentId := uuid.New()

	// WHEN
	user, err := adapter.GetPostgresqlUser(nonExistentId)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, errors.ObjectNotFoundError.UserNotFound, *err)
}

func TestGetUserByEmailSuccessfully(t *testing.T) {
	/*
		GIVEN: A user exists in the database
		WHEN:  GetPostgresqlUserByEmail is called with the user email
		THEN:  The user is returned with all related data
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewUserAdapterTestWrapper(t)

	userEmail := "test@example.com"
	existingUser := factories.NewUserModel(db, factories.UserModelF{
		Email: &userEmail,
	})

	// WHEN
	user, err := adapter.GetPostgresqlUserByEmail(existingUser.Email)

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, existingUser.Id, user.Id)
	assert.Equal(t, existingUser.Name, user.Name)
	assert.Equal(t, existingUser.Email, user.Email)
}

func TestGetUserByEmailNotFound(t *testing.T) {
	/*
		GIVEN: No user exists with the given email
		WHEN:  GetPostgresqlUserByEmail is called with non-existent email
		THEN:  A not found error is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewUserAdapterTestWrapper(t)

	nonExistentEmail := "nonexistent@example.com"

	// WHEN
	user, err := adapter.GetPostgresqlUserByEmail(nonExistentEmail)

	// THEN
	assert.NotNil(t, err)
	assert.Nil(t, user)
	assert.Equal(t, errors.ObjectNotFoundError.UserNotFound, *err)
}

func TestFetchUsersSuccessfully(t *testing.T) {
	/*
		GIVEN: Multiple users exist in the database
		WHEN:  FetchPostgresqlUsers is called
		THEN:  All users are returned
	*/
	// GIVEN
	adapter, _, db := adapterTest.NewUserAdapterTestWrapper(t)

	user1Email := "user1@example.com"
	user2Email := "user2@example.com"
	user1 := factories.NewUserModel(db, factories.UserModelF{
		Email: &user1Email,
	})
	user2 := factories.NewUserModel(db, factories.UserModelF{
		Email: &user2Email,
	})

	// WHEN
	users, err := adapter.FetchPostgresqlUsers()

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.GreaterOrEqual(t, len(users), 2)

	// Find our created users
	foundUser1 := false
	foundUser2 := false
	for _, user := range users {
		if user.Id == user1.Id {
			foundUser1 = true
			assert.Equal(t, user1.Email, user.Email)
		}
		if user.Id == user2.Id {
			foundUser2 = true
			assert.Equal(t, user2.Email, user.Email)
		}
	}
	assert.True(t, foundUser1)
	assert.True(t, foundUser2)
}

func TestFetchUsersEmpty(t *testing.T) {
	/*
		GIVEN: No users exist in the database
		WHEN:  FetchPostgresqlUsers is called
		THEN:  An empty list is returned
	*/
	// GIVEN
	adapter, _, _ := adapterTest.NewUserAdapterTestWrapper(t)

	// WHEN
	users, err := adapter.FetchPostgresqlUsers()

	// THEN
	assert.Nil(t, err)
	assert.NotNil(t, users)
	assert.Equal(t, 0, len(users))
}
