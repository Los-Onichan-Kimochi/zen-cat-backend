package model

import "github.com/google/uuid"

type UserRol string

const (
	UserRolAdmin  UserRol = "ADMIN"
	UserRolClient UserRol = "CLIENT"
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

	OnboardingId *uuid.UUID  `gorm:"type:uuid"`
	Onboarding   *Onboarding `gorm:"foreignKey:OnboardingId"`
}

func (User) TableName() string {
	return "astro_cat_user"
}
