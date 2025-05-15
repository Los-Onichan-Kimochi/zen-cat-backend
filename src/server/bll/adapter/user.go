package adapter

import (
	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type User struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

func NewUserAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *User {
	return &User{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

func (u *User) GetPostgresqlUser(
	userId uuid.UUID,
) (*schemas.User, *errors.Error) {
	userModel, err := u.DaoPostgresql.User.GetUser(userId)
	if err != nil {
		return nil, &errors.ObjectNotFoundError.UserNotFound
	}

	// Mapear memberships
	var memberships []*schemas.Membership
	for _, m := range userModel.Memberships {
		memberships = append(memberships, &schemas.Membership{
			Id:          m.Id,
			Description: m.Description,
			StartDate:   m.StartDate,
			EndDate:     m.EndDate,
			Status:      schemas.MembershipStatus(m.Status),
			Community: schemas.Community{
				Id:                  m.Community.Id,
				Name:                m.Community.Name,
				Purpose:             m.Community.Purpose,
				ImageUrl:            m.Community.ImageUrl,
				NumberSubscriptions: m.Community.NumberSubscriptions,
			},
			Plan: schemas.Plan{
				Id:               m.Plan.Id,
				Fee:              m.Plan.Fee,
				Type:             schemas.PlanType(m.Plan.Type),
				ReservationLimit: m.Plan.ReservationLimit,
			},
		})
	}

	return &schemas.User{
		Id:             userModel.Id,
		Name:           userModel.Name,
		FirstLastName:  userModel.FirstLastName,
		SecondLastName: userModel.SecondLastName,
		Password:       userModel.Password,
		Email:          userModel.Email,
		Rol:            schemas.UserRol(userModel.Rol),
		ImageUrl:       userModel.ImageUrl,
		Memberships:    memberships,
	}, nil
}

func (u *User) FetchPostgresqlUsers() ([]*schemas.User, *errors.Error) {
	usersModel, err := u.DaoPostgresql.User.FetchUsers()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.UserNotFound
	}

	users := make([]*schemas.User, len(usersModel))
	for i, userModel := range usersModel {
		var memberships []*schemas.Membership
		for _, m := range userModel.Memberships {
			memberships = append(memberships, &schemas.Membership{
				Id:          m.Id,
				Description: m.Description,
				StartDate:   m.StartDate,
				EndDate:     m.EndDate,
				Status:      schemas.MembershipStatus(m.Status),
				Community: schemas.Community{
					Id:                  m.Community.Id,
					Name:                m.Community.Name,
					Purpose:             m.Community.Purpose,
					ImageUrl:            m.Community.ImageUrl,
					NumberSubscriptions: m.Community.NumberSubscriptions,
				},
				Plan: schemas.Plan{
					Id:               m.Plan.Id,
					Fee:              m.Plan.Fee,
					Type:             schemas.PlanType(m.Plan.Type),
					ReservationLimit: m.Plan.ReservationLimit,
				},
			})
		}
		users[i] = &schemas.User{
			Id:             userModel.Id,
			Name:           userModel.Name,
			FirstLastName:  userModel.FirstLastName,
			SecondLastName: userModel.SecondLastName,
			Password:       userModel.Password,
			Email:          userModel.Email,
			Rol:            schemas.UserRol(userModel.Rol),
			ImageUrl:       userModel.ImageUrl,
			Memberships:    memberships,
		}
	}

	return users, nil
}

func (u *User) CreatePostgresqlUser(
	name string,
	firstLastName string,
	secondLastName *string,
	password string,
	email string,
	rol string,
	imageUrl string,
	updatedBy string,
) (*schemas.User, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	userModel := &model.User{
		Id:             uuid.New(),
		Name:           name,
		FirstLastName:  firstLastName,
		SecondLastName: secondLastName,
		Password:       password,
		Email:          email,
		Rol:            model.UserRol(rol),
		ImageUrl:       imageUrl,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := u.DaoPostgresql.User.CreateUser(userModel); err != nil {
		return nil, &errors.BadRequestError.UserNotCreated
	}

	return &schemas.User{
		Id:             userModel.Id,
		Name:           userModel.Name,
		FirstLastName:  userModel.FirstLastName,
		SecondLastName: userModel.SecondLastName,
		Password:       userModel.Password,
		Email:          userModel.Email,
		Rol:            schemas.UserRol(userModel.Rol),
		ImageUrl:       userModel.ImageUrl,
		// Memberships:    userModel.Memberships,
		// Onboarding:     userModel.Onboarding,
	}, nil
}

func (u *User) UpdatePostgresqlUser(
	userId uuid.UUID,
	name *string,
	firstLastName *string,
	secondLastName *string,
	password *string,
	email *string,
	rol *string,
	imageUrl *string,
	updatedBy string,
) (*schemas.User, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	userModel, err := u.DaoPostgresql.User.UpdateUser(
		userId,
		name,
		firstLastName,
		secondLastName,
		password,
		email,
		rol,
		imageUrl,
		updatedBy,
	)
	if err != nil {
		return nil, &errors.BadRequestError.UserNotUpdated
	}

	return &schemas.User{
		Id:             userModel.Id,
		Name:           userModel.Name,
		FirstLastName:  userModel.FirstLastName,
		SecondLastName: userModel.SecondLastName,
		Password:       userModel.Password,
		Email:          userModel.Email,
		Rol:            schemas.UserRol(userModel.Rol),
		ImageUrl:       userModel.ImageUrl,
		// Memberships:    userModel.Memberships,
	}, nil
}
