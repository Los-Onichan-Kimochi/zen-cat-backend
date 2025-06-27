package model

import "github.com/google/uuid"

type UserRol string

const (
	UserRolAdmin  UserRol = "ADMINISTRATOR"
	UserRolClient UserRol = "CLIENT"
	UserRolGuest  UserRol = "GUEST"
)

type User struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name           string
	FirstLastName  string
	SecondLastName *string
	Password       string
	Email          string
	Rol            UserRol
	ImageUrl       string
	AuditFields

	Onboarding  *Onboarding   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Memberships []*Membership `gorm:"foreignKey:UserId"`
}

func (User) TableName() string {
	return "astro_cat_user"
}
