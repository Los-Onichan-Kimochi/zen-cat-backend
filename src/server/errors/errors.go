package errors

import (
	"net/http"
	"reflect"

	"github.com/labstack/echo/v4"
)

type Error struct {
	Code    string
	Message string
}

var (
	// For 404 Not Found errors
	ObjectNotFoundError = struct {
		CommunityNotFound Error
	}{
		CommunityNotFound: Error{
			Code:    "COMMUNITY_ERROR_001",
			Message: "Community not found",
		},
	}

	// For 422 Unprocessable Entity errors
	UnprocessableEntityError = struct {
		InvalidCommunityId Error
		InvalidRequestBody Error
	}{
		InvalidRequestBody: Error{
			Code:    "REQUEST_ERROR_001",
			Message: "Invalid body request",
		},
		InvalidCommunityId: Error{
			Code:    "COMMUNITY_ERROR_004",
			Message: "Invalid community id",
		},
	}

	// For 400 Bad Request errors
	BadRequestError = struct {
		InvalidUpdatedByValue Error
		CommunityNotCreated   Error
		CommunityNotUpdated   Error
	}{
		InvalidUpdatedByValue: Error{
			Code:    "REQUEST_ERROR_002",
			Message: "Invalid updated by value error",
		},
		CommunityNotCreated: Error{
			Code:    "COMMUNITY_ERROR_002",
			Message: "Community not created",
		},
		CommunityNotUpdated: Error{
			Code:    "COMMUNITY_ERROR_003",
			Message: "Community not updated",
		},
	}

	// For 500 Internal Server errors
	InternalServerError = struct{}{}
)

// Helper function to check if an error is in a specific error group.
func isInErrorGroup(err Error, group interface{}) bool {
	val := reflect.ValueOf(group)
	for i := 0; i < val.NumField(); i++ {
		if val.Field(i).Interface() == err {
			return true
		}
	}
	return false
}

// General error handler function for endpoints.
func HandleError(err Error, c echo.Context) error {
	var statusCode int
	switch {
	case isInErrorGroup(err, ObjectNotFoundError):
		statusCode = http.StatusNotFound

	case isInErrorGroup(err, UnprocessableEntityError):
		statusCode = http.StatusUnprocessableEntity

	case isInErrorGroup(err, BadRequestError):
		statusCode = http.StatusBadRequest

	case isInErrorGroup(err, InternalServerError):
		statusCode = http.StatusInternalServerError

	default:
		statusCode = http.StatusInternalServerError // Default case for other errors
	}

	// Send JSON response with the error code and message
	return c.JSON(statusCode, err)
}
