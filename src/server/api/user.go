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

// @Summary 			Check User Exists by Email.
// @Description 		Checks if a user exists with the given email address.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               email    query   string  true  "Email address"
// @Success 			200 {object} schemas.CheckUserExistsResponse "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/check-email/ [get]
func (a *Api) CheckUserExists(c echo.Context) error {
	email := c.QueryParam("email")
	if email == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserEmail, c)
	}

	response, err := a.BllController.User.CheckUserExistsByEmail(email)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary      Change password with email
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        changePassword body schemas.ChangePasswordInput true "Email and new password"
// @Success      200 {object} echo.Map
// @Failure      400,401,422,500 {object} errors.Error
// @Router       /user/change-password/ [post]
func (a *Api) ChangePassword(c echo.Context) error {
	// Parse body
	var request schemas.ChangePasswordInput
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Cambiar la contraseña
	if err := a.BllController.User.ChangePassword(request.Email, request); err != nil {
		return errors.HandleError(*err, c)
	}

	// Respuesta exitosa
	return c.JSON(http.StatusOK, echo.Map{
		"message": "Password changed successfully",
	})
}

// @Summary 			Change User Role.
// @Description 		Changes the role of a user given its id.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               changeRoleRequest    body   schemas.ChangeUserRoleRequest  true  "Change Role Request"
// @Success 			200 {object} schemas.User "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			403 {object} errors.Error "Forbidden - Admin role required"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/{userId}/role/ [patch]
func (a *Api) ChangeUserRole(c echo.Context) error {
	updateBy := "ADMIN"

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	var request schemas.ChangeUserRoleRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Convert to UpdateUserRequest
	rolStr := string(request.Rol)
	updateRequest := schemas.UpdateUserRequest{
		Rol: &rolStr,
	}

	response, err := a.BllController.User.UpdateUser(userId, updateRequest, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get User Statistics.
// @Description 		Get user statistics including role distribution and recent connections.
// @Tags 				User
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.UserStats "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			403 {object} errors.Error "Forbidden - Admin role required"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/user/stats/ [get]
func (a *Api) GetUserStats(c echo.Context) error {
	response, err := a.BllController.User.GetUserStats()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
