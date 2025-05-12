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

	Onboarding *Onboarding `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

func (User) TableName() string {
	return "astro_cat_user"
}
