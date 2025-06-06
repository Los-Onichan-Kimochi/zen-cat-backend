package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get User.
// @Description 		Gets a user given its id.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/{userId}/ [get]
func (a *Api) GetUser(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	response, err := a.BllController.User.GetUser(userId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Users.
// @Description 		Fetch all users, filtered by params.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Users "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/ [get]
func (a *Api) FetchUsers(c echo.Context) error {
	response, err := a.BllController.User.FetchUsers()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create User.
// @Description 		Creates a new user.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               createUserRequest    body   schemas.CreateUserRequest  true  "Create User Request"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/ [post]
func (a *Api) CreateUser(c echo.Context) error {
	updateBy := "ADMIN"

	var request schemas.CreateUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	response, err := a.BllController.User.CreateUser(request, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update User.
// @Description 		Updates a user given its id.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               updateUserRequest    body   schemas.UpdateUserRequest  true  "Update User Request"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/{userId}/ [patch]
func (a *Api) UpdateUser(c echo.Context) error {
	updateBy := "ADMIN"

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	var request schemas.UpdateUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	response, err := a.BllController.User.UpdateUser(userId, request, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete User.
// @Description 		Deletes a user given its id.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/{userId}/ [delete]
func (a *Api) DeleteUser(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	if err := a.BllController.User.DeleteUser(userId); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Delete Users.
// @Description 		Bulk delete users given their ids.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteUserRequest true  "Bulk Delete User Request"
// @Success 			204 "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/bulk-delete/ [delete]
func (a *Api) BulkDeleteUsers(c echo.Context) error {
	var request schemas.BulkDeleteUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if err := a.BllController.User.BulkDeleteUsers(request); err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Bulk Create Users.
// @Description 		Bulk creates users.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkCreateUserRequest true  "Bulk Create User Request"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/bulk-create/ [post]
func (a *Api) BulkCreateUsers(c echo.Context) error {
	updateBy := "ADMIN"

	var request schemas.BulkCreateUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	response, err := a.BllController.User.BulkCreateUsers(request.Users, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Bulk Delete Users.
// @Description 		Bulk deletes users given their ids.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.BulkDeleteUserRequest true  "Bulk Delete User Request"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/bulk-delete/ [delete]
func (a *Api) BulkDeleteUsers(c echo.Context) error {
	var request schemas.BulkDeleteUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	if err := a.BllController.User.BulkDeleteUsers(request); err != nil {
		return errors.HandleError(*err, c)
	}
	return c.NoContent(http.StatusNoContent)
}
