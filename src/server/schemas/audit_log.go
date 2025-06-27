package schemas

import (
	"time"

	"github.com/google/uuid"
)

type AuditActionType string

const (
	// Admin actions
	AuditActionCreate     AuditActionType = "CREATE"
	AuditActionUpdate     AuditActionType = "UPDATE"
	AuditActionDelete     AuditActionType = "DELETE"
	AuditActionBulkCreate AuditActionType = "BULK_CREATE"
	AuditActionBulkDelete AuditActionType = "BULK_DELETE"

	// User actions
	AuditActionLogin             AuditActionType = "LOGIN"
	AuditActionLogout            AuditActionType = "LOGOUT"
	AuditActionRegister          AuditActionType = "REGISTER"
	AuditActionSubscribe         AuditActionType = "SUBSCRIBE"
	AuditActionUnsubscribe       AuditActionType = "UNSUBSCRIBE"
	AuditActionCreateReservation AuditActionType = "CREATE_RESERVATION"
	AuditActionCancelReservation AuditActionType = "CANCEL_RESERVATION"
	AuditActionUpdateProfile     AuditActionType = "UPDATE_PROFILE"
)

type AuditEntityType string

const (
	AuditEntityUser                AuditEntityType = "USER"
	AuditEntityCommunity           AuditEntityType = "COMMUNITY"
	AuditEntityProfessional        AuditEntityType = "PROFESSIONAL"
	AuditEntityLocal               AuditEntityType = "LOCAL"
	AuditEntityPlan                AuditEntityType = "PLAN"
	AuditEntityService             AuditEntityType = "SERVICE"
	AuditEntitySession             AuditEntityType = "SESSION"
	AuditEntityReservation         AuditEntityType = "RESERVATION"
	AuditEntityMembership          AuditEntityType = "MEMBERSHIP"
	AuditEntityOnboarding          AuditEntityType = "ONBOARDING"
	AuditEntityCommunityPlan       AuditEntityType = "COMMUNITY_PLAN"
	AuditEntityCommunityService    AuditEntityType = "COMMUNITY_SERVICE"
	AuditEntityServiceLocal        AuditEntityType = "SERVICE_LOCAL"
	AuditEntityServiceProfessional AuditEntityType = "SERVICE_PROFESSIONAL"
)

type AuditLog struct {
	Id             uuid.UUID       `json:"id"`
	UserId         uuid.UUID       `json:"user_id"`
	UserEmail      string          `json:"user_email"`
	UserRole       UserRol         `json:"user_role"`
	Action         AuditActionType `json:"action"`
	EntityType     AuditEntityType `json:"entity_type"`
	EntityId       *uuid.UUID      `json:"entity_id"`
	EntityName     *string         `json:"entity_name"`
	OldValues      *string         `json:"old_values"`
	NewValues      *string         `json:"new_values"`
	IPAddress      string          `json:"ip_address"`
	UserAgent      *string         `json:"user_agent"`
	AdditionalInfo *string         `json:"additional_info"`
	Success        bool            `json:"success"`
	ErrorMessage   *string         `json:"error_message"`
	CreatedAt      time.Time       `json:"created_at"`
	User           *UserProfile    `json:"user,omitempty"`
}

type AuditLogs struct {
	AuditLogs  []*AuditLog `json:"audit_logs"`
	TotalCount int64       `json:"total_count"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
}

type AuditLogFilters struct {
	UserIds     []string   `json:"user_ids"`
	Actions     []string   `json:"actions"`
	EntityTypes []string   `json:"entity_types"`
	UserRoles   []string   `json:"user_roles"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
	Success     *bool      `json:"success"`
	Page        int        `json:"page"`
	PageSize    int        `json:"page_size"`
}

type AuditStats struct {
	TotalEvents  int64            `json:"total_events"`
	Actions      []ActionStat     `json:"actions"`
	UserRoles    []UserRoleStat   `json:"user_roles"`
	EntityTypes  []EntityTypeStat `json:"entity_types"`
	SuccessCount int64            `json:"success_count"`
	FailureCount int64            `json:"failure_count"`
	ActiveUsers  int64            `json:"active_users"`
}

type ActionStat struct {
	Action string `json:"action"`
	Count  int64  `json:"count"`
}

type UserRoleStat struct {
	UserRole string `json:"user_role"`
	Count    int64  `json:"count"`
}

type EntityTypeStat struct {
	EntityType string `json:"entity_type"`
	Count      int64  `json:"count"`
}

// Helper structs for creating audit logs
type AuditContext struct {
	UserId    uuid.UUID
	UserEmail string
	UserRole  UserRol
	IPAddress string
	UserAgent *string
}

type AuditEvent struct {
	Action         AuditActionType
	EntityType     AuditEntityType
	EntityId       *uuid.UUID
	EntityName     *string
	OldValues      interface{}
	NewValues      interface{}
	AdditionalInfo *string
	Success        bool
	ErrorMessage   *string
}
