package api

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// @Summary 			Get Professional.
// @Description 		Gets a professional given its id.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Success 			200 {object} schemas.Professional "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/ [get]
func (a *Api) GetProfessional(c echo.Context) error {
	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	response, err := a.BllController.Professional.GetProfessional(professionalId)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Fetch Professionals.
// @Description 		Fetch all professionals, filtered by params.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.Professionals "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/ [get]
func (a *Api) FetchProfessionals(c echo.Context) error {
	response, err := a.BllController.Professional.FetchProfessionals()
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// @Summary 			Create Professional.
// @Description 		Creates a new professional.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               request	body   schemas.CreateProfessionalRequest true  "Create Professional Request"
// @Success 			201 {object} schemas.Professional "Created"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/ [post]
func (a *Api) CreateProfessional(c echo.Context) error {
	updateBy := "ADMIN"

	var request schemas.CreateProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	response, err := a.BllController.Professional.CreateProfessional(request, updateBy)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(http.StatusCreated, response)
}

// @Summary 			Update Professional.
// @Description 		Updates a professional given its id.
// @Tags 				Professional
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Param               professionalId    path   string  true  "Professional ID"
// @Param               request	body   schemas.UpdateProfessionalRequest true  "Update Professional Request"
// @Success 			200 {object} schemas.Professional "OK"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Missing or malformed JWT"
// @Failure 			404 {object} errors.Error "Not Found"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/professional/{professionalId}/ [patch]
func (a *Api) UpdateProfessional(c echo.Context) error {
	updateBy := "ADMIN"

	professionalId, parseErr := uuid.Parse(c.Param("professionalId"))
	if parseErr != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidProfessionalId, c)
	}

	var request schemas.UpdateProfessionalRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}
	response, err := a.BllController.Professional.UpdateProfessional(
		professionalId,
		request,
		updateBy,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}
	return c.JSON(http.StatusOK, response)
}
