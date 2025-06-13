package model

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
	Id             uuid.UUID       `gorm:"type:uuid;primaryKey"`
	UserId         uuid.UUID       `gorm:"type:uuid;not null"`     // User who performed the action
	UserEmail      string          `gorm:"size:255;not null"`      // User email for easy filtering
	UserRole       UserRol         `gorm:"size:50;not null"`       // User role (ADMIN/CLIENT)
	Action         AuditActionType `gorm:"size:50;not null;index"` // Action performed
	EntityType     AuditEntityType `gorm:"size:50;not null;index"` // Type of entity affected
	EntityId       *uuid.UUID      `gorm:"type:uuid;index"`        // ID of the entity affected (nullable for bulk operations)
	EntityName     *string         `gorm:"size:255"`               // Name/description of the entity for display
	OldValues      *string         `gorm:"type:text"`              // JSON of old values (for updates)
	NewValues      *string         `gorm:"type:text"`              // JSON of new values (for creates/updates)
	IPAddress      string          `gorm:"size:45"`                // User's IP address
	UserAgent      *string         `gorm:"size:500"`               // User's browser/client info
	AdditionalInfo *string         `gorm:"type:text"`              // Additional context (e.g., bulk operation details)
	Success        bool            `gorm:"not null;default:true"`  // Whether the action was successful
	ErrorMessage   *string         `gorm:"size:1000"`              // Error message if action failed
	CreatedAt      time.Time       `gorm:"autoCreateTime;index"`   // When the action occurred

	// Relations
	User *User `gorm:"foreignKey:UserId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

func (AuditLog) TableName() string {
	return "astro_cat_audit_log"
}
