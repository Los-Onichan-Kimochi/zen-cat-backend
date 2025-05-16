package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Create CommunityPlan.
// @Description 		Associates a community with a plan.
// @Tags 				CommunityPlan
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request body schemas.CreateCommunityPlanRequest true "Community-Plan Association Request"
// @Success 			201 {object} schemas.CommunityPlan "Created"
// @Failure 			400 {object} errors.Error "Bad Request (e.g., invalid updatedBy)"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found (Community or Plan does not exist)"
// @Failure 			409 {object} errors.Error "Conflict (Association already exists)"
// @Failure 			422 {object} errors.Error "Unprocessable Entity (Invalid UUIDs or request body)"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/community-plan/ [post]
func (a *Api) CreateCommunityPlan(c echo.Context) error {
	updatedBy := "ADMIN"

	var request schemas.CreateCommunityPlanRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	if request.CommunityId == uuid.Nil || request.PlanId == uuid.Nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidCommunityPlanId, c)
	}

	response, err := a.BllController.CommunityPlan.CreateCommunityPlan(request, updatedBy)
	if err != nil {
		// TODO: Handle error more properly
		if err.Code == errors.BadRequestError.CommunityPlanAlreadyExists.Code {
			return c.JSON(http.StatusConflict, err)
		}
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}
