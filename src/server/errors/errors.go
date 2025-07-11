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
		CommunityNotFound            Error
		ReservationNotFound          Error
		ProfessionalNotFound         Error
		LocalNotFound                Error
		UserNotFound                 Error
		ServiceNotFound              Error
		PlanNotFound                 Error
		MembershipNotFound           Error
		OnboardingNotFound           Error
		CommunityPlanNotFound        Error
		CommunityServiceNotFound     Error
		ServiceLocalNotFound         Error
		ServiceProfessionalNotFound  Error
		SessionNotFound              Error
		AuditLogNotFound             Error
		MembershipSuspensionNotFound Error
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
		CommunityPlanNotFound: Error{
			Code:    "COMMUNITY_PLAN_ERROR_001",
			Message: "Community-Plan association not found",
		},
		CommunityServiceNotFound: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_001",
			Message: "Community-Service association not found",
		},
		ServiceLocalNotFound: Error{
			Code:    "SERVICE_LOCAL_ERROR_001",
			Message: "Service-Local association not found",
		},
		ServiceProfessionalNotFound: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_001",
			Message: "Service-Professional association not found",
		},
		SessionNotFound: Error{
			Code:    "SESSION_ERROR_001",
			Message: "Session not found",
		},
		AuditLogNotFound: Error{
			Code:    "AUDIT_LOG_ERROR_001",
			Message: "Audit log not found",
		},
		MembershipSuspensionNotFound: Error{
			Code:    "MEMBERSHIP_SUSPENSION_ERROR_001",
			Message: "Membership suspension not found",
		},
	}

	// For 422 Unprocessable Entity errors
	UnprocessableEntityError = struct {
		InvalidCommunityId            Error
		InvalidRequestBody            Error
		InvalidProfessionalId         Error
		InvalidLocalId                Error
		InvalidServiceId              Error
		InvalidPlanId                 Error
		InvalidMembershipId           Error
		InvalidOnboardingId           Error
		InvalidUserEmail              Error
		InvalidUserId                 Error
		InvalidCommunityPlanId        Error
		InvalidCommunityServiceId     Error
		InvalidParsingInteger         Error
		InvalidServiceLocalId         Error
		InvalidServiceProfessionalId  Error
		InvalidSessionId              Error
		InvalidReservationId          Error
		InvalidMembershipSuspensionId Error
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
		InvalidCommunityPlanId: Error{
			Code:    "COMMUNITY_PLAN_ERROR_004",
			Message: "Invalid community_id or plan_id for association",
		},
		InvalidCommunityServiceId: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_004",
			Message: "Invalid community_id or service_id for association",
		},
		InvalidParsingInteger: Error{
			Code:    "REQUEST_ERROR_004",
			Message: "Invalid parsing integer",
		},
		InvalidServiceLocalId: Error{
			Code:    "SERVICE_LOCAL_ERROR_004",
			Message: "Invalid service_id or local_id for association",
		},
		InvalidServiceProfessionalId: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_004",
			Message: "Invalid service_id or professional_id for association",
		},
		InvalidSessionId: Error{
			Code:    "SESSION_ERROR_004",
			Message: "Invalid session id",
		},
		InvalidReservationId: Error{
			Code:    "RESERVATION_ERROR_004",
			Message: "Invalid reservation id",
		},
		InvalidMembershipSuspensionId: Error{
			Code:    "MEMBERSHIP_SUSPENSION_ERROR_004",
			Message: "Invalid membership suspension id",
		},
	}

	// For 400 Bad Request errors
	BadRequestError = struct {
		InvalidUpdatedByValue          Error
		InvalidCommunityName           Error
		InvalidServiceName             Error
		DuplicateCommunityName         Error
		DuplicateUserEmail             Error
		CommunityNotCreated            Error
		CommunityNotUpdated            Error
		CommunityNotSoftDeleted        Error
		LocalNotCreated                Error
		LocalNotUpdated                Error
		LocalNotSoftDeleted            Error
		ProfessionalNotCreated         Error
		ProfessionalNotUpdated         Error
		ProfessionalNotSoftDeleted     Error
		ServiceNotCreated              Error
		ServiceNotUpdated              Error
		ServiceNotSoftDeleted          Error
		PlanNotCreated                 Error
		PlanNotUpdated                 Error
		PlanNotSoftDeleted             Error
		InvalidPlanType                Error
		MembershipNotCreated           Error
		MembershipNotUpdated           Error
		MembershipNotDeleted           Error
		OnboardingNotCreated           Error
		OnboardingNotUpdated           Error
		UserNotCreated                 Error
		UserNotUpdated                 Error
		UserNotSoftDeleted             Error
		UserPasswordNotUpdated         Error
		CommunityPlanNotCreated        Error
		CommunityPlanNotDeleted        Error
		CommunityServiceNotCreated     Error
		CommunityServiceNotDeleted     Error
		ServiceLocalNotCreated         Error
		ServiceLocalNotDeleted         Error
		ServiceProfessionalNotCreated  Error
		ServiceProfessionalNotDeleted  Error
		SessionNotCreated              Error
		SessionNotUpdated              Error
		SessionNotSoftDeleted          Error
		MembershipSuspensionNotCreated Error
		MembershipSuspensionNotUpdated Error
	}{
		InvalidUpdatedByValue: Error{
			Code:    "BAD_REQUEST_ERROR_001",
			Message: "Invalid updated by value",
		},
		InvalidCommunityName: Error{
			Code:    "BAD_REQUEST_ERROR_004",
			Message: "Community name cannot be empty",
		},
		InvalidServiceName: Error{
			Code:    "BAD_REQUEST_ERROR_005",
			Message: "Service name cannot be empty",
		},
		DuplicateCommunityName: Error{
			Code:    "BAD_REQUEST_ERROR_002",
			Message: "Community name already exists",
		},
		DuplicateUserEmail: Error{
			Code:    "BAD_REQUEST_ERROR_003",
			Message: "User email already exists",
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
		ProfessionalNotSoftDeleted: Error{
			Code:    "PROFESSIONAL_ERROR_005",
			Message: "Professional not soft deleted",
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
		MembershipNotDeleted: Error{
			Code:    "MEMBERSHIP_ERROR_005",
			Message: "Membership not deleted",
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
		ServiceNotSoftDeleted: Error{
			Code:    "SERVICE_ERROR_005",
			Message: "Service not soft deleted",
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
		UserNotSoftDeleted: Error{
			Code:    "USER_ERROR_005",
			Message: "User not soft deleted",
		},
		UserPasswordNotUpdated: Error{
			Code:    "USER_ERROR_007",
			Message: "Failed to update user password",
		},
		CommunityPlanNotCreated: Error{
			Code:    "COMMUNITY_PLAN_ERROR_002",
			Message: "Community-Plan association not created",
		},
		CommunityPlanNotDeleted: Error{
			Code:    "COMMUNITY_PLAN_ERROR_005",
			Message: "Community-Plan association not deleted",
		},
		CommunityServiceNotCreated: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_002",
			Message: "Community-Service association not created",
		},
		CommunityServiceNotDeleted: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_005",
			Message: "Community-Service association not deleted",
		},
		ServiceLocalNotCreated: Error{
			Code:    "SERVICE_LOCAL_ERROR_002",
			Message: "Service-Local association not created",
		},
		ServiceLocalNotDeleted: Error{
			Code:    "SERVICE_LOCAL_ERROR_005",
			Message: "Service-Local association not deleted",
		},
		ServiceProfessionalNotCreated: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_002",
			Message: "Service-Professional association not created",
		},
		ServiceProfessionalNotDeleted: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_005",
			Message: "Service-Professional association not deleted",
		},
		SessionNotCreated: Error{
			Code:    "SESSION_ERROR_002",
			Message: "Session not created",
		},
		SessionNotUpdated: Error{
			Code:    "SESSION_ERROR_003",
			Message: "Session not updated",
		},
		SessionNotSoftDeleted: Error{
			Code:    "SESSION_ERROR_005",
			Message: "Session not soft deleted",
		},
		MembershipSuspensionNotCreated: Error{
			Code:    "MEMBERSHIP_SUSPENSION_ERROR_002",
			Message: "Membership suspension not created",
		},
		MembershipSuspensionNotUpdated: Error{
			Code:    "MEMBERSHIP_SUSPension_ERROR_003",
			Message: "Membership suspension not updated",
		},
	}

	ContactError = struct {
		MissingFields      Error
		InvalidEmailFormat Error
		FailedToSendEmail  Error
	}{
		MissingFields: Error{
			Code:    "CONTACT_ERROR_001",
			Message: "Name, email and message are required.",
		},
		InvalidEmailFormat: Error{
			Code:    "CONTACT_ERROR_002",
			Message: "Email format is invalid.",
		},
		FailedToSendEmail: Error{
			Code:    "CONTACT_ERROR_003",
			Message: "Failed to send contact email.",
		},
	}

	// For 401 Unauthorized errors
	AuthenticationError = struct {
		UnauthorizedUser    Error
		InvalidRefreshToken Error
		InvalidAccessToken  Error
	}{
		UnauthorizedUser: Error{
			Code:    "AUTHENTICATION_ERROR_001",
			Message: "Unauthorized",
		},
		InvalidRefreshToken: Error{
			Code:    "AUTHENTICATION_ERROR_002",
			Message: "Invalid refresh token",
		},
		InvalidAccessToken: Error{
			Code:    "AUTHENTICATION_ERROR_003",
			Message: "Invalid access token",
		},
	}

	// For 403 Forbidden errors
	ForbiddenError = struct {
		InsufficientPrivileges Error
	}{
		InsufficientPrivileges: Error{
			Code:    "FORBIDDEN_ERROR_001",
			Message: "Insufficient privileges to access this resource",
		},
	}

	// For 409 Conflict errors
	ConflictError = struct {
		CommunityPlanAlreadyExists       Error
		CommunityServiceAlreadyExists    Error
		ServiceProfessionalAlreadyExists Error
		ServiceLocalAlreadyExists        Error
		UserAlreadyExists                Error
		SessionTimeConflict              Error
		UserReservationTimeConflict      Error
	}{
		CommunityPlanAlreadyExists: Error{
			Code:    "COMMUNITY_PLAN_ERROR_006",
			Message: "Community-Plan association already exists",
		},
		CommunityServiceAlreadyExists: Error{
			Code:    "COMMUNITY_SERVICE_ERROR_006",
			Message: "Community-Service association already exists",
		},
		UserAlreadyExists: Error{
			Code:    "USER_ERROR_006",
			Message: "User already exists with this email",
		},
		ServiceLocalAlreadyExists: Error{
			Code:    "SERVICE_LOCAL_ERROR_003",
			Message: "Service-Local association already exists",
		},
		ServiceProfessionalAlreadyExists: Error{
			Code:    "SERVICE_PROFESSIONAL_ERROR_003",
			Message: "Service-Professional association already exists",
		},
		SessionTimeConflict: Error{
			Code:    "CONFLICT_ERROR_001",
			Message: "Session time conflict detected",
		},
		UserReservationTimeConflict: Error{
			Code:    "CONFLICT_ERROR_002",
			Message: "User has another reservation at the same time",
		},
	}

	// For 500 Internal Server errors
	InternalServerError = struct {
		Default               Error
		FailedToUploadImage   Error
		FailedToDownloadImage Error
		DatabaseError         Error
	}{
		Default: Error{
			Code:    "INTERNAL_SERVER_ERROR_001",
			Message: "An unexpected error occurred.",
		},
		FailedToUploadImage: Error{
			Code:    "INTERNAL_SERVER_ERROR_002",
			Message: "Failed to upload image to S3",
		},
		FailedToDownloadImage: Error{
			Code:    "INTERNAL_SERVER_ERROR_003",
			Message: "Failed to download image from external storage",
		},
		DatabaseError: Error{
			Code:    "INTERNAL_SERVER_ERROR_004",
			Message: "Database error",
		},
	}

	// For forgot password or recovery flows
	ForgotPasswordError = struct {
		InvalidEmail        Error
		FailedToSendEmail   Error
		InvalidOrExpiredPin Error
		FailedToSendSMS     Error
	}{
		InvalidEmail: Error{
			Code:    "FORGOT_PASSWORD_ERROR_001",
			Message: "Email not associated to any account",
		},
		FailedToSendEmail: Error{
			Code:    "FORGOT_PASSWORD_ERROR_002",
			Message: "Failed to send reset email",
		},
		InvalidOrExpiredPin: Error{
			Code:    "FORGOT_PASSWORD_ERROR_003",
			Message: "Invalid or expired PIN code",
		},
		FailedToSendSMS: Error{
			Code:    "FORGOT_PASSWORD_ERROR_004",
			Message: "Failed to send SMS with recovery PIN",
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

	case isInErrorGroup(err, ConflictError):
		statusCode = http.StatusConflict

	case isInErrorGroup(err, InternalServerError):
		statusCode = http.StatusInternalServerError

	case isInErrorGroup(err, AuthenticationError):
		statusCode = http.StatusUnauthorized

	case isInErrorGroup(err, ForbiddenError):
		statusCode = http.StatusForbidden

	default:
		statusCode = http.StatusInternalServerError // Default case for other errors
	}

	// Send JSON response with the error code and message
	return c.JSON(statusCode, err)
}
