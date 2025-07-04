package adapter

import (
	"github.com/google/uuid"
	daoPostgresql "onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/controller"
	"onichankimochi.com/astro_cat_backend/src/server/dao/astro_cat_psql/model"
	"onichankimochi.com/astro_cat_backend/src/server/errors"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"

	"fmt"
	"time"

	"gorm.io/gorm"
	"onichankimochi.com/astro_cat_backend/src/logging"
)

type Community struct {
	logger        logging.Logger
	DaoPostgresql *daoPostgresql.AstroCatPsqlCollection
}

// Creates Community adapter
func NewCommunityAdapter(
	logger logging.Logger,
	daoPostgresql *daoPostgresql.AstroCatPsqlCollection,
) *Community {
	return &Community{
		logger:        logger,
		DaoPostgresql: daoPostgresql,
	}
}

// Gets a community from postgresql DB and adapts it to a Community schema.
func (c *Community) GetPostgresqlCommunity(
	communityId uuid.UUID,
) (*schemas.Community, *errors.Error) {
	communityModel, err := c.DaoPostgresql.Community.GetCommunity(communityId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.CommunityNotFound
		}
		return nil, &errors.BadRequestError.CommunityNotCreated
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}

// Fetch communities from postgresql DB and adapts them to a Community schema.
func (c *Community) FetchPostgresqlCommunities() ([]*schemas.Community, *errors.Error) {
	communitiesModel, err := c.DaoPostgresql.Community.FetchCommunities()
	if err != nil {
		return nil, &errors.ObjectNotFoundError.CommunityNotFound
	}

	communities := make([]*schemas.Community, len(communitiesModel))
	for i, communityModel := range communitiesModel {
		communities[i] = &schemas.Community{
			Id:                  communityModel.Id,
			Name:                communityModel.Name,
			Purpose:             communityModel.Purpose,
			ImageUrl:            communityModel.ImageUrl,
			NumberSubscriptions: communityModel.NumberSubscriptions,
		}
	}

	return communities, nil
}

// Creates a community.
func (c *Community) CreatePostgresqlCommunity(
	name string,
	purpose string,
	imageUrl string,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	// Validate name is not empty
	if name == "" {
		return nil, &errors.BadRequestError.InvalidCommunityName
	}

	// Check for duplicate community name
	existingCommunity, _ := c.DaoPostgresql.Community.GetCommunityByName(name)
	if existingCommunity != nil {
		return nil, &errors.BadRequestError.DuplicateCommunityName
	}

	communityModel := &model.Community{
		Id:                  uuid.New(),
		Name:                name,
		Purpose:             purpose,
		ImageUrl:            imageUrl,
		NumberSubscriptions: 0, // Default number of initial subscriptions
		AuditFields: model.AuditFields{
			UpdatedBy: updatedBy,
		},
	}

	if err := c.DaoPostgresql.Community.CreateCommunity(communityModel); err != nil {
		return nil, &errors.BadRequestError.CommunityNotCreated
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}

// Creates multiple communities into postgresql DB and returns them.
func (c *Community) BulkCreatePostgresqlCommunities(
	communitiesData []*schemas.CreateCommunityRequest,
	updatedBy string,
) ([]*schemas.Community, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communitiesModel := make([]*model.Community, len(communitiesData))
	for i, communityData := range communitiesData {
		communitiesModel[i] = &model.Community{
			Id:                  uuid.New(),
			Name:                communityData.Name,
			Purpose:             communityData.Purpose,
			ImageUrl:            communityData.ImageUrl,
			NumberSubscriptions: 0,
			AuditFields: model.AuditFields{
				UpdatedBy: updatedBy,
			},
		}
	}

	if err := c.DaoPostgresql.Community.BulkCreateCommunities(communitiesModel); err != nil {
		return nil, &errors.BadRequestError.CommunityNotCreated
	}

	communities := make([]*schemas.Community, len(communitiesModel))
	for i, communityModel := range communitiesModel {
		communities[i] = &schemas.Community{
			Id:                  communityModel.Id,
			Name:                communityModel.Name,
			Purpose:             communityModel.Purpose,
			ImageUrl:            communityModel.ImageUrl,
			NumberSubscriptions: communityModel.NumberSubscriptions,
		}
	}

	return communities, nil
}

// Updates a community from a Postgresql DB given its ID and adapts it to a community schema.
func (c *Community) UpdatePostgresqlCommunity(
	id uuid.UUID,
	name *string,
	purpose *string,
	imageUrl *string,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	if updatedBy == "" {
		return nil, &errors.BadRequestError.InvalidUpdatedByValue
	}

	communityModel, err := c.DaoPostgresql.Community.UpdateCommunity(id, name, purpose, imageUrl, updatedBy)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &errors.ObjectNotFoundError.CommunityNotFound
		}
		return nil, &errors.BadRequestError.CommunityNotUpdated
	}

	return &schemas.Community{
		Id:                  communityModel.Id,
		Name:                communityModel.Name,
		Purpose:             communityModel.Purpose,
		ImageUrl:            communityModel.ImageUrl,
		NumberSubscriptions: communityModel.NumberSubscriptions,
	}, nil
}

// Soft deletes a community from postgresql DB.
func (c *Community) DeletePostgresqlCommunity(communityId uuid.UUID) *errors.Error {
	err := c.DaoPostgresql.Community.DeleteCommunity(communityId)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return &errors.ObjectNotFoundError.CommunityNotFound
		}
		return &errors.BadRequestError.CommunityNotSoftDeleted
	}

	return nil
}

// Bulk deletes communities from postgresql DB.
func (c *Community) BulkDeletePostgresqlCommunities(
	communityIds []uuid.UUID,
) *errors.Error {
	if err := c.DaoPostgresql.Community.BulkDeleteCommunities(communityIds); err != nil {
		return &errors.BadRequestError.CommunityNotSoftDeleted
	}

	return nil
}

// GetCommunityReport obtiene el reporte de comunidades para el dashboard admin
type CommunityReportParams struct {
	From    *time.Time
	To      *time.Time
	GroupBy string // "day", "week", "month"
}

type CommunityReportData struct {
	CommunityId   string
	CommunityName string
	Total         int
	// Métricas por estado
	ActiveMemberships    int
	ExpiredMemberships   int
	CancelledMemberships int
	// Métricas de engagement
	ActiveUsers       int // usuarios que han hecho reservas
	InactiveUsers     int // usuarios sin reservas
	TotalReservations int
	// Métricas de planes
	MonthlyPlans int
	AnnualPlans  int
	// Datos temporales
	Data []struct {
		Date  string `json:"date"`
		Count int    `json:"count"`
	}
}

func (c *Community) GetCommunityReport(params CommunityReportParams) (total int, communities []CommunityReportData, err error) {
	// Construir la consulta base para comunidades
	query := c.DaoPostgresql.Community.PostgresqlDB.Model(&model.Community{})

	var communityModels []model.Community
	if err := query.Find(&communityModels).Error; err != nil {
		return 0, nil, err
	}

	// Obtener membresías para cada comunidad con preload completo
	membershipsQuery := c.DaoPostgresql.Membership.PostgresqlDB.Model(&model.Membership{}).
		Preload("Community").
		Preload("User").
		Preload("Plan")

	// Filtrar por vigencia real de la membresía
	if params.From != nil {
		membershipsQuery = membershipsQuery.Where("end_date >= ?", *params.From)
	}
	if params.To != nil {
		membershipsQuery = membershipsQuery.Where("start_date <= ?", *params.To)
	}

	var memberships []model.Membership
	if err := membershipsQuery.Find(&memberships).Error; err != nil {
		return 0, nil, err
	}

	// Obtener reservas para calcular engagement
	reservationsQuery := c.DaoPostgresql.Reservation.PostgresqlDB.Model(&model.Reservation{}).
		Preload("Membership.Community")

	if params.From != nil {
		reservationsQuery = reservationsQuery.Where("reservation_time >= ?", *params.From)
	}
	if params.To != nil {
		reservationsQuery = reservationsQuery.Where("reservation_time <= ?", *params.To)
	}

	var reservations []model.Reservation
	if err := reservationsQuery.Find(&reservations).Error; err != nil {
		return 0, nil, err
	}

	// Procesar datos por comunidad
	communityData := make(map[string]*CommunityReportData)
	userActivity := make(map[string]map[uuid.UUID]bool) // communityId -> userId -> hasReservations

	// Inicializar datos de actividad de usuarios
	for _, membership := range memberships {
		communityId := membership.CommunityId.String()
		if userActivity[communityId] == nil {
			userActivity[communityId] = make(map[uuid.UUID]bool)
		}
	}

	// Procesar reservas para determinar usuarios activos
	for _, reservation := range reservations {
		if reservation.Membership != nil {
			communityId := reservation.Membership.CommunityId.String()
			userId := reservation.UserId
			userActivity[communityId][userId] = true
		}
	}

	// Procesar membresías
	for _, membership := range memberships {
		communityId := membership.CommunityId.String()
		communityName := membership.Community.Name

		if communityData[communityId] == nil {
			communityData[communityId] = &CommunityReportData{
				CommunityId:          communityId,
				CommunityName:        communityName,
				Total:                0,
				ActiveMemberships:    0,
				ExpiredMemberships:   0,
				CancelledMemberships: 0,
				ActiveUsers:          0,
				InactiveUsers:        0,
				TotalReservations:    0,
				MonthlyPlans:         0,
				AnnualPlans:          0,
				Data: make([]struct {
					Date  string `json:"date"`
					Count int    `json:"count"`
				}, 0),
			}
		}

		// Contar por estado
		switch membership.Status {
		case model.MembershipStatusActive:
			communityData[communityId].ActiveMemberships++
		case model.MembershipStatusExpired:
			communityData[communityId].ExpiredMemberships++
		case model.MembershipStatusCancelled:
			communityData[communityId].CancelledMemberships++
		}

		// Contar por tipo de plan
		switch membership.Plan.Type {
		case model.PlanTypeMonthly:
			communityData[communityId].MonthlyPlans++
		case model.PlanTypeAnual:
			communityData[communityId].AnnualPlans++
		}

		// Agrupar por fecha según el parámetro groupBy (usando start_date)
		var dateKey string
		switch params.GroupBy {
		case "month":
			dateKey = membership.StartDate.Format("2006-01")
		case "week":
			y, w := membership.StartDate.ISOWeek()
			dateKey = fmt.Sprintf("%d-W%02d", y, w)
		default:
			dateKey = membership.StartDate.Format("2006-01-02")
		}

		// Incrementar contadores
		communityData[communityId].Total++

		// Agregar dato por fecha
		found := false
		for i, data := range communityData[communityId].Data {
			if data.Date == dateKey {
				communityData[communityId].Data[i].Count++
				found = true
				break
			}
		}
		if !found {
			communityData[communityId].Data = append(communityData[communityId].Data, struct {
				Date  string `json:"date"`
				Count int    `json:"count"`
			}{
				Date:  dateKey,
				Count: 1,
			})
		}
	}

	// Calcular métricas de engagement y reservas por comunidad
	for _, reservation := range reservations {
		if reservation.Membership != nil {
			communityId := reservation.Membership.CommunityId.String()
			if communityData[communityId] != nil {
				communityData[communityId].TotalReservations++
			}
		}
	}

	// Calcular usuarios activos vs inactivos por comunidad
	for communityId, data := range communityData {
		activeUsers := 0
		inactiveUsers := 0

		// Contar usuarios únicos por comunidad
		userCount := make(map[uuid.UUID]bool)
		for _, membership := range memberships {
			if membership.CommunityId.String() == communityId {
				userCount[membership.UserId] = true
			}
		}

		// Determinar si cada usuario está activo
		for userId := range userCount {
			if userActivity[communityId][userId] {
				activeUsers++
			} else {
				inactiveUsers++
			}
		}

		data.ActiveUsers = activeUsers
		data.InactiveUsers = inactiveUsers
	}

	// Construir la respuesta
	for _, data := range communityData {
		communities = append(communities, *data)
		total += data.Total
	}

	return total, communities, nil
}
