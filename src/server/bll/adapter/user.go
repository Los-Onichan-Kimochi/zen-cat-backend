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
				Type:             model.PlanType(m.Plan.Type),
				ReservationLimit: m.Plan.ReservationLimit,
			},
		})
	}
	// Mapear onboarding (si existe)
	var onboarding *schemas.Onboarding
	if userModel.Onboarding != nil {
		onboarding = &schemas.Onboarding{
			Id:             userModel.Onboarding.Id,
			DocumentType:   schemas.DocumentType(userModel.Onboarding.DocumentType),
			DocumentNumber: userModel.Onboarding.DocumentNumber,
			PhoneNumber:    userModel.Onboarding.PhoneNumber,
			BirthDate:      userModel.Onboarding.BirthDate,
			Gender:         (*schemas.Gender)(userModel.Onboarding.Gender),
			City:           userModel.Onboarding.City,
			PostalCode:     userModel.Onboarding.PostalCode,
			District:       userModel.Onboarding.District,
			Address:        userModel.Onboarding.Address,
		}
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
		Onboarding:     onboarding,
	}, nil
}

func (u *User) GetPostgresqlUserByEmail(
	email string,
) (*schemas.User, *errors.Error) {
	userModel, err := u.DaoPostgresql.User.GetUserByEmail(email)
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
				Type:             model.PlanType(m.Plan.Type),
				ReservationLimit: m.Plan.ReservationLimit,
			},
		})
	}
	// Mapear onboarding (si existe)
	var onboarding *schemas.Onboarding
	if userModel.Onboarding != nil {
		onboarding = &schemas.Onboarding{
			Id:             userModel.Onboarding.Id,
			DocumentType:   schemas.DocumentType(userModel.Onboarding.DocumentType),
			DocumentNumber: userModel.Onboarding.DocumentNumber,
			PhoneNumber:    userModel.Onboarding.PhoneNumber,
			BirthDate:      userModel.Onboarding.BirthDate,
			Gender:         (*schemas.Gender)(userModel.Onboarding.Gender),
			City:           userModel.Onboarding.City,
			PostalCode:     userModel.Onboarding.PostalCode,
			District:       userModel.Onboarding.District,
			Address:        userModel.Onboarding.Address,
		}
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
		Onboarding:     onboarding,
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
					Type:             model.PlanType(m.Plan.Type),
					ReservationLimit: m.Plan.ReservationLimit,
				},
			})
		}
		// Mapear onboarding (si existe)
		var onboarding *schemas.Onboarding
		if userModel.Onboarding != nil {
			onboarding = &schemas.Onboarding{
				Id:             userModel.Onboarding.Id,
				DocumentType:   schemas.DocumentType(userModel.Onboarding.DocumentType),
				DocumentNumber: userModel.Onboarding.DocumentNumber,
				PhoneNumber:    userModel.Onboarding.PhoneNumber,
				BirthDate:      userModel.Onboarding.BirthDate,
				Gender:         (*schemas.Gender)(userModel.Onboarding.Gender),
				City:           userModel.Onboarding.City,
				PostalCode:     userModel.Onboarding.PostalCode,
				District:       userModel.Onboarding.District,
				Address:        userModel.Onboarding.Address,
			}
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
			Onboarding:     onboarding,
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
	memberships []*schemas.Membership,
	onboarding *schemas.Onboarding,
) (*schemas.User, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	var membershipsModel []*model.Membership
	for _, m := range memberships {
		membershipsModel = append(membershipsModel, &model.Membership{
			Id:          m.Id,
			Description: m.Description,
			StartDate:   m.StartDate,
			EndDate:     m.EndDate,
			Status:      model.MembershipStatus(m.Status),
			Community: model.Community{
				Id:                  m.Community.Id,
				Name:                m.Community.Name,
				Purpose:             m.Community.Purpose,
				ImageUrl:            m.Community.ImageUrl,
				NumberSubscriptions: m.Community.NumberSubscriptions,
			},
			Plan: model.Plan{
				Id:               m.Plan.Id,
				Fee:              m.Plan.Fee,
				Type:             model.PlanType(m.Plan.Type),
				ReservationLimit: m.Plan.ReservationLimit,
			},
		})
	}

	var onboardingModel *model.Onboarding
	if onboarding != nil {
		onboardingModel = &model.Onboarding{
			Id:             uuid.New(),
			DocumentType:   model.DocumentType(onboarding.DocumentType),
			DocumentNumber: onboarding.DocumentNumber,
			PhoneNumber:    onboarding.PhoneNumber,
			BirthDate:      onboarding.BirthDate,
			Gender:         (*model.Gender)(onboarding.Gender),
			City:           onboarding.City,
			PostalCode:     onboarding.PostalCode,
			District:       onboarding.District,
			Address:        onboarding.Address,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
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
		Memberships:    membershipsModel,
		Onboarding:     onboardingModel,
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	// Establecer la relación UserId en el onboarding
	if onboardingModel != nil {
		onboardingModel.UserId = userModel.Id
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
		Memberships:    memberships,
		Onboarding:     onboarding,
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
	memberships []*schemas.Membership,
	onboarding *schemas.Onboarding,
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
		Memberships:    memberships,
		Onboarding:     onboarding,
	}, nil
}

func (u *User) DeletePostgresqlUser(userId uuid.UUID) *errors.Error {
	// Primero eliminar el onboarding asociado si existe
	onboardingAdapter := NewOnboardingAdapter(u.logger, u.DaoPostgresql)
	onboardingAdapter.DeletePostgresqlOnboardingByUserId(userId) // Ignoramos el error si no existe

	// Luego eliminar el usuario
	if err := u.DaoPostgresql.User.DeleteUser(userId); err != nil {
		return &errors.BadRequestError.UserNotSoftDeleted
	}
	return nil
}

func (u *User) BulkCreatePostgresqlUser(
	usersData []*schemas.CreateUserRequest,
	updatedBy string,
) ([]*schemas.User, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	usersModel := make([]*model.User, len(usersData))
	for i, userData := range usersData {
		var secondLastNamePtr *string
		if userData.SecondLastName != "" {
			secondLastNamePtr = &userData.SecondLastName
		} else {
			secondLastNamePtr = nil
		}
		usersModel[i] = &model.User{
			Id:             uuid.New(),
			Name:           userData.Name,
			FirstLastName:  userData.FirstLastName,
			SecondLastName: secondLastNamePtr,
			Password:       userData.Password,
			Email:          userData.Email,
			Rol:            model.UserRol(userData.Rol),
			ImageUrl:       userData.ImageUrl,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}
	if err := u.DaoPostgresql.User.BulkCreateUsers(usersModel); err != nil {
		return nil, &errors.BadRequestError.UserNotCreated
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
					Type:             model.PlanType(m.Plan.Type),
					ReservationLimit: m.Plan.ReservationLimit,
				},
			})
		}
		// Mapear onboarding (si existe)
		var onboarding *schemas.Onboarding
		if userModel.Onboarding != nil {
			onboarding = &schemas.Onboarding{
				Id:             userModel.Onboarding.Id,
				DocumentType:   schemas.DocumentType(userModel.Onboarding.DocumentType),
				DocumentNumber: userModel.Onboarding.DocumentNumber,
				PhoneNumber:    userModel.Onboarding.PhoneNumber,
				BirthDate:      userModel.Onboarding.BirthDate,
				Gender:         (*schemas.Gender)(userModel.Onboarding.Gender),
				City:           userModel.Onboarding.City,
				PostalCode:     userModel.Onboarding.PostalCode,
				District:       userModel.Onboarding.District,
				Address:        userModel.Onboarding.Address,
			}
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
			Onboarding:     onboarding,
		}
	}

	return users, nil
}

func (u *User) BulkDeletePostgresqlUser(
	userIds []uuid.UUID,
) *errors.Error {
	// Primero eliminar los onboardings asociados si existen
	onboardingAdapter := NewOnboardingAdapter(u.logger, u.DaoPostgresql)
	for _, userId := range userIds {
		onboardingAdapter.DeletePostgresqlOnboardingByUserId(userId) // Ignoramos errores si no existen
	}

	// Luego eliminar los usuarios
	if err := u.DaoPostgresql.User.BulkDeleteUsers(userIds); err != nil {
		return &errors.BadRequestError.UserNotSoftDeleted
	}
	return nil
}
