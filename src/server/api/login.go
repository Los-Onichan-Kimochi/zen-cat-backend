package api

import (
	"bytes"
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// Login 				godoc
// @Summary 			User Login
// @Description 		Authenticate user with email and password, returns user info and tokens
// @Tags 				Login
// @Accept 				json
// @Produce 			json
// @Param               loginRequest    body   schemas.LoginRequest  true  "Login Request"
// @Success 			200 {object} schemas.LoginResponse "Login successful"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Unauthorized - Invalid credentials"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/login/ [post]
func (a *Api) Login(c echo.Context) error {
	var request schemas.LoginRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Login.Login(request.Email, request.Password)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}

// Register 			godoc
// @Summary 			User Registration
// @Description 		Register a new user with email and password, returns user info and tokens
// @Tags 				Login
// @Accept 				json
// @Produce 			json
// @Param               registerRequest    body   schemas.RegisterRequest  true  "Register Request"
// @Success 			201 {object} schemas.LoginResponse "Registration successful"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			409 {object} errors.Error "Conflict - User already exists"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/register/ [post]
func (a *Api) Register(c echo.Context) error {
	var request schemas.RegisterRequest
	if err := c.Bind(&request); err != nil {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Manual validation for required fields
	if request.Name == "" || request.FirstLastName == "" || request.Email == "" || request.Password == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	// Additional validation for password length (min=6 as per schema)
	if len(request.Password) < 6 {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Login.Register(
		request.Name,
		request.FirstLastName,
		request.SecondLastName,
		request.Email,
		request.Password,
		request.ImageUrl,
	)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusCreated, response)
}

// GetCurrentUser 		godoc
// @Summary 			Get Current User Info
// @Description 		Get current authenticated user information from token
// @Tags 				Login
// @Accept 				json
// @Produce 			json
// @Security			JWT
// @Success 			200 {object} schemas.UserProfile "User info"
// @Failure 			401 {object} errors.Error "Unauthorized"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/me/ [get]
func (a *Api) GetCurrentUser(c echo.Context) error {
	_, credentials, authError := a.BllController.Auth.AccessTokenValidation(c)
	if authError != nil {
		return errors.HandleError(*authError, c)
	}

	// Get the first role as primary rol
	var userRol schemas.UserRol
	if len(credentials.UserRoles) > 0 {
		userRol = schemas.UserRol(credentials.UserRoles[0])
	} else {
		userRol = schemas.UserRolClient // Default to client
	}

	userProfile := schemas.UserProfile{
		Id:             credentials.UserId,
		Name:           credentials.UserName,
		FirstLastName:  credentials.UserFirstName,
		SecondLastName: credentials.UserLastName,
		Email:          credentials.UserEmail,
		Rol:            userRol,
		ImageUrl:       credentials.UserImageUrl,
	}

	return c.JSON(http.StatusOK, userProfile)
}

// GoogleLogin 		godoc
// @Summary 			Google Login
// @Description 		Authenticate user using Google ID token, returns user info and tokens
// @Tags 				Login
// @Accept 				json
// @Produce 			json
// @Param               googleLoginRequest  body   schemas.GoogleLoginRequest  true  "Google Login Request"
// @Success 			200 {object} schemas.GoogleLoginResponse "Login successful"
// @Failure 			400 {object} errors.Error "Bad Request"
// @Failure 			401 {object} errors.Error "Unauthorized - Invalid token"
// @Failure 			422 {object} errors.Error "Unprocessable Entity"
// @Failure 			500 {object} errors.Error "Internal Server Error"
// @Router 				/login/google/ [post]
func (a *Api) GoogleLogin(c echo.Context) error {
	// 🔍 Extra para depurar si sigue fallando
	body, _ := io.ReadAll(c.Request().Body)
	c.Request().Body = io.NopCloser(bytes.NewBuffer(body))

	var request schemas.GoogleLoginRequest
	if err := c.Bind(&request); err != nil || request.Token == "" {
		return errors.HandleError(errors.UnprocessableEntityError.InvalidRequestBody, c)
	}

	response, err := a.BllController.Login.GoogleLogin(c.Request().Context(), request.Token)
	if err != nil {
		return errors.HandleError(*err, c)
	}

	return c.JSON(http.StatusOK, response)
}
