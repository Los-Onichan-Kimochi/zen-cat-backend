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
		CommunityNotFound    Error
		ReservationNotFound  Error
		ProfessionalNotFound Error
		ServiceNotFound      Error
		PlanNotFound         Error
	}{
		CommunityNotFound: Error{
			Code:    "COMMUNITY_ERROR_001",
			Message: "Community not found",
		},
		ProfessionalNotFound: Error{
			Code:    "PROFESSIONAL_ERROR_001",
			Message: "Professional not found",
		},
		ServiceNotFound: Error{
			Code:    "SERVICE_ERROR_001",
			Message: "Service not found",
		},
		PlanNotFound: Error{
			Code:    "PLAN_ERROR_001",
			Message: "Plan not found",
		},
	}

	// For 422 Unprocessable Entity errors
	UnprocessableEntityError = struct {
		InvalidCommunityId    Error
		InvalidRequestBody    Error
		InvalidProfessionalId Error
		InvalidServiceId      Error
		InvalidPlanId         Error
	}{
		InvalidRequestBody: Error{
			Code:    "REQUEST_ERROR_001",
			Message: "Invalid body request",
		},
		InvalidCommunityId: Error{
			Code:    "COMMUNITY_ERROR_004",
			Message: "Invalid community id",
		},
		InvalidProfessionalId: Error{
			Code:    "PROFESSIONAL_ERROR_004",
			Message: "Invalid professional id",
		},
		InvalidServiceId: Error{
			Code:    "SERVICE_ERROR_004",
			Message: "Invalid service id",
		},
		InvalidPlanId: Error{
			Code:    "PLAN_ERROR_004",
			Message: "Invalid plan id",
		},
	}

	// For 400 Bad Request errors
	BadRequestError = struct {
		InvalidUpdatedByValue   Error
		CommunityNotCreated     Error
		CommunityNotUpdated     Error
		CommunityNotSoftDeleted Error
		ProfessionalNotCreated  Error
		ProfessionalNotUpdated  Error
		ServiceNotCreated       Error
		ServiceNotUpdated       Error
		PlanNotCreated          Error
		PlanNotUpdated          Error
		PlanNotSoftDeleted      Error
		InvalidPlanType         Error
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
		CommunityNotSoftDeleted: Error{
			Code:    "COMMUNITY_ERROR_005",
			Message: "Community not soft deleted",
		},
		ProfessionalNotCreated: Error{
			Code:    "PROFESSIONAL_ERROR_002",
			Message: "Professional not created",
		},
		ProfessionalNotUpdated: Error{
			Code:    "PROFESSIONAL_ERROR_003",
			Message: "Professional not updated",
		},
		ServiceNotCreated: Error{
			Code:    "SERVICE_ERROR_002",
			Message: "Service not created",
		},
		ServiceNotUpdated: Error{
			Code:    "SERVICE_ERROR_003",
			Message: "Service not updated",
		},
		PlanNotCreated: Error{
			Code:    "PLAN_ERROR_002",
			Message: "Plan not created",
		},
		PlanNotUpdated: Error{
			Code:    "PLAN_ERROR_003",
			Message: "Plan not updated",
		},
		PlanNotSoftDeleted: Error{
			Code:    "PLAN_ERROR_006",
			Message: "Plan not soft deleted",
		},
		InvalidPlanType: Error{
			Code:    "PLAN_ERROR_005",
			Message: "Invalid plan type",
		},
	}

	// For 500 Internal Server errors
	InternalServerError = struct {
		Default Error
	}{
		Default: Error{
			Code:    "INTERNAL_SERVER_ERROR_001",
			Message: "An unexpected error occurred.",
		},
	}
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
