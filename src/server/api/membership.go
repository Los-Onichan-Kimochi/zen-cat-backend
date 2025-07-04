package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Membership by ID.
// @Description 		Gets a membership given its id.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               membershipId    path   string  true  "Membership ID"
// @Success 			200 {object} schemas.Membership "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/{membershipId}/ [get]
func (a *Api) GetMembership(c echo.Context) error {
	membershipId, parseErr := uuid.Parse(c.Param("membershipId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidMembershipId, c)
	}

	response, err := a.BllController.Membership.GetMembership(membershipId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Memberships by User ID.
// @Description 		Gets all memberships for a given user id.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Success 			200 {object} schemas.Memberships "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/user/{userId}/ [get]
func (a *Api) GetMembershipsByUserId(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	response, err := a.BllController.Membership.GetMembershipsByUserId(userId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Memberships by Community ID.
// @Description 		Gets all memberships for a given community id.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Success 			200 {object} schemas.Memberships "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/community/{communityId}/ [get]
func (a *Api) GetMembershipsByCommunityId(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	response, err := a.BllController.Membership.GetMembershipsByCommunityId(communityId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Memberships.
// @Description 		Fetch all memberships.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Memberships "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/ [get]
func (a *Api) FetchMemberships(c echo.Context) error {
	response, err := a.BllController.Membership.FetchMemberships()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Membership.
// @Description 		Creates a new membership.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               createMembershipRequest    body   schemas.CreateMembershipRequest  true  "Create Membership Request"
// @Success 			201 {object} schemas.Membership "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/ [post]
func (a *Api) CreateMembership(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	var request schemas.CreateMembershipRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Membership.CreateMembership(request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Create Membership for User.
// @Description 		Creates a new membership for a specific user.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               createMembershipForUserRequest    body   schemas.CreateMembershipForUserRequest  true  "Create Membership For User Request"
// @Success 			201 {object} schemas.Membership "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/user/{userId}/ [post]
func (a *Api) CreateMembershipForUser(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	var request schemas.CreateMembershipForUserRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Membership.CreateMembershipForUser(userId, request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Membership by ID.
// @Description 		Updates a membership given its id.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               membershipId    path   string  true  "Membership ID"
// @Param               updateMembershipRequest    body   schemas.UpdateMembershipRequest  true  "Update Membership Request"
// @Success 			200 {object} schemas.Membership "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/{membershipId}/ [patch]
func (a *Api) UpdateMembership(c echo.Context) error {
	updatedBy := "USER" // Podría ser obtenido del JWT token en producción

	membershipId, parseErr := uuid.Parse(c.Param("membershipId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidMembershipId, c)
	}

	var request schemas.UpdateMembershipRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Membership.UpdateMembership(membershipId, request, updatedBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Delete Membership by ID.
// @Description 		Deletes a membership given its id.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               membershipId    path   string  true  "Membership ID"
// @Success 			204 {string} string "No Content"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/{membershipId}/ [delete]
func (a *Api) DeleteMembership(c echo.Context) error {
	membershipId, parseErr := uuid.Parse(c.Param("membershipId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidMembershipId, c)
	}

	err := a.BllController.Membership.DeleteMembership(membershipId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.NoContent(http.StatusNoContent)
}

// @Summary 			Get Users by Community ID.
// @Description 		Gets all users who have active memberships in a specific community.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               communityId    path   string  true  "Community ID"
// @Success 			200 {object} schemas.Users "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/community/{communityId}/users [get]
func (a *Api) GetUsersByCommunityId(c echo.Context) error {
	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	response, err := a.BllController.Membership.GetUsersByCommunityId(communityId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Get Membership by User ID and Community ID.
// @Description 		Gets a membership for a given user ID and community ID.
// @Tags 				Membership
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               userId    path   string  true  "User ID"
// @Param               communityId    path   string  true  "Community ID"
// @Success 			200 {object} schemas.Membership "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/membership/user/{userId}/community/{communityId} [get]
func (a *Api) GetMembershipByUserAndCommunity(c echo.Context) error {
	userId, parseErr := uuid.Parse(c.Param("userId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidUserId, c)
	}

	communityId, parseErr := uuid.Parse(c.Param("communityId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityId, c)
	}

	response, err := a.BllController.Membership.GetMembershipByUserAndCommunity(userId, communityId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
