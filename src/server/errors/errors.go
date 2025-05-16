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
		LocalNotFound        Error
		UserNotFound         Error
		ServiceNotFound      Error
		PlanNotFound         Error
		MembershipNotFound   Error
		OnboardingNotFound   Error
	}{
		CommunityNotFound: Error{
			Code:    "COMMUNITY_ERROR_001",
			Message: "Community not found",
		},
		ProfessionalNotFound: Error{
			Code:    "PROFESSIONAL_ERROR_001",
			Message: "Professional not found",
		},
		LocalNotFound: Error{
			Code:    "LOCAL_ERROR_001",
			Message: "Local not found",
		},
		ServiceNotFound: Error{
			Code:    "SERVICE_ERROR_001",
			Message: "Service not found",
		},
		PlanNotFound: Error{
			Code:    "PLAN_ERROR_001",
			Message: "Plan not found",
		},
		UserNotFound: Error{
			Code:    "USER_ERROR_001",
			Message: "User not found",
		},
		MembershipNotFound: Error{
			Code:    "MEMBERSHIP_ERROR_001",
			Message: "Membership not found",
		},
		OnboardingNotFound: Error{
			Code:    "ONBOARDING_ERROR_001",
			Message: "Onboarding not found",
		},
	}

	// For 422 Unprocessable Entity errors
	UnprocessableEntityError = struct {
		InvalidCommunityId    Error
		InvalidRequestBody    Error
		InvalidProfessionalId Error
		InvalidLocalId        Error
		InvalidServiceId      Error
		InvalidPlanId         Error
		InvalidMembershipId   Error
		InvalidOnboardingId   Error
		InvalidUserEmail      Error
		InvalidUserId         Error
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
		InvalidLocalId: Error{
			Code:    "LOCAL_ERROR_004",
			Message: "Invalid local id",
		},
		InvalidServiceId: Error{
			Code:    "SERVICE_ERROR_004",
			Message: "Invalid service id",
		},
		InvalidPlanId: Error{
			Code:    "PLAN_ERROR_004",
			Message: "Invalid plan id",
		},
		InvalidMembershipId: Error{
			Code:    "MEMBERSHIP_ERROR_001",
			Message: "Invalid membership id",
		},
		InvalidOnboardingId: Error{
			Code:    "ONBOARDING_ERROR_001",
			Message: "Invalid onboarding id",
		},
		InvalidUserEmail: Error{
			Code:    "USER_ERROR_001",
			Message: "Invalid user email",
		},
		InvalidUserId: Error{
			Code:    "USER_ERROR_004",
			Message: "Invalid user id",
		},
	}

	// For 400 Bad Request errors
	BadRequestError = struct {
		InvalidUpdatedByValue  Error
		CommunityNotCreated    Error
		CommunityNotUpdated    Error
		LocalNotCreated        Error
		LocalNotUpdated        Error
		LocalNotSoftDeleted    Error
		ProfessionalNotCreated Error
		ProfessionalNotUpdated Error
		ServiceNotCreated      Error
		ServiceNotUpdated      Error
		PlanNotCreated         Error
		PlanNotUpdated         Error
		PlanNotSoftDeleted     Error
		InvalidPlanType        Error
		MembershipNotCreated   Error
		MembershipNotUpdated   Error
		OnboardingNotCreated   Error
		OnboardingNotUpdated   Error
		UserNotCreated         Error
		UserNotUpdated         Error
		UserNotSoftDeleted     Error
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
		ProfessionalNotCreated: Error{
			Code:    "PROFESSIONAL_ERROR_002",
			Message: "Professional not created",
		},
		ProfessionalNotUpdated: Error{
			Code:    "PROFESSIONAL_ERROR_003",
			Message: "Professional not updated",
		},
		LocalNotCreated: Error{
			Code:    "LOCAL_ERROR_002",
			Message: "Local not created",
		},
		LocalNotUpdated: Error{
			Code:    "LOCAL_ERROR_003",
			Message: "Local not updated",
		},
		LocalNotSoftDeleted: Error{
			Code:    "LOCAL_ERROR_005",
			Message: "Local not soft deleted",
		},
		MembershipNotCreated: Error{
			Code:    "MEMBERSHIP_ERROR_002",
			Message: "Membership not created",
		},
		MembershipNotUpdated: Error{
			Code:    "MEMBERSHIP_ERROR_003",
			Message: "Membership not updated",
		},
		OnboardingNotCreated: Error{
			Code:    "ONBOARDING_ERROR_002",
			Message: "Onboarding not created",
		},
		OnboardingNotUpdated: Error{
			Code:    "ONBOARDING_ERROR_003",
			Message: "Onboarding not updated",
		},
		UserNotCreated: Error{
			Code:    "USER_ERROR_002",
			Message: "User not created",
		},
		UserNotUpdated: Error{
			Code:    "USER_ERROR_003",
			Message: "User not updated",
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
		InvalidPlanType: Error{
			Code:    "PLAN_ERROR_005",
			Message: "Invalid plan type",
		},
		PlanNotSoftDeleted: Error{
			Code:    "PLAN_ERROR_006",
			Message: "Plan not soft deleted",
		},
		UserNotSoftDeleted: Error{
			Code:    "USER_ERROR_005",
			Message: "User not soft deleted",
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
