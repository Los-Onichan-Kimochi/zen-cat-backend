package controller

import (
	"time"

	"github.com/google/uuid"
	"onichankimochi.com/astro_cat_backend/src/logging"
	bllAdapter "onichankimochi.com/astro_cat_backend/src/server/bll/adapter"
	errors "onichankimochi.com/astro_cat_backend/src/server/errors"
	schemas "onichankimochi.com/astro_cat_backend/src/server/schemas"
)

type Community struct {
	logger      logging.Logger
	Adapter     *bllAdapter.AdapterCollection
	EnvSettings *schemas.EnvSettings
}

// Create Community controller
func NewCommunityController(
	logger logging.Logger,
	adapter *bllAdapter.AdapterCollection,
	envSettings *schemas.EnvSettings,
) *Community {
	return &Community{
		logger:      logger,
		Adapter:     adapter,
		EnvSettings: envSettings,
	}
}

// Gets a community.
func (c *Community) GetCommunity(communityId uuid.UUID) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.GetPostgresqlCommunity(communityId)
}

// Fetch all communities.
func (c *Community) FetchCommunities() (*schemas.Communities, *errors.Error) {
	communities, err := c.Adapter.Community.FetchPostgresqlCommunities()
	if err != nil {
		return nil, err
	}

	return &schemas.Communities{Communities: communities}, nil
}

// Creates a community.
func (c *Community) CreateCommunity(
	createCommunityData schemas.CreateCommunityRequest,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.CreatePostgresqlCommunity(
		createCommunityData.Name,
		createCommunityData.Purpose,
		createCommunityData.ImageUrl,
		updatedBy,
	)
}

// Updates a community.
func (c *Community) UpdateCommunity(
	communityId uuid.UUID,
	updateCommunityData schemas.UpdateCommunityRequest,
	updatedBy string,
) (*schemas.Community, *errors.Error) {
	return c.Adapter.Community.UpdatePostgresqlCommunity(
		communityId,
		updateCommunityData.Name,
		updateCommunityData.Purpose,
		updateCommunityData.ImageUrl,
		updatedBy,
	)
}

// Soft deletes a community.
func (c *Community) DeleteCommunity(communityId uuid.UUID) *errors.Error {
	return c.Adapter.Community.DeletePostgresqlCommunity(communityId)
}

// Creates multiple communities
func (c *Community) BulkCreateCommunities(
	createCommunitiesData []*schemas.CreateCommunityRequest,
	updatedBy string,
) ([]*schemas.Community, *errors.Error) {
	return c.Adapter.Community.BulkCreatePostgresqlCommunities(createCommunitiesData, updatedBy)
}

// Bulk deletes communities.
func (c *Community) BulkDeleteCommunities(
	bulkDeleteCommunityData schemas.BulkDeleteCommunityRequest,
) *errors.Error {
	return c.Adapter.Community.BulkDeletePostgresqlCommunities(
		bulkDeleteCommunityData.Communities,
	)
}

// GetCommunityReport obtiene el reporte de comunidades para el dashboard admin
type CommunityReportResponse struct {
	Total       int                              `json:"totalMemberships"`
	Communities []bllAdapter.CommunityReportData `json:"communities"`
	// Métricas agregadas
	Summary struct {
		TotalActiveMemberships    int `json:"totalActiveMemberships"`
		TotalExpiredMemberships   int `json:"totalExpiredMemberships"`
		TotalCancelledMemberships int `json:"totalCancelledMemberships"`
		TotalActiveUsers          int `json:"totalActiveUsers"`
		TotalInactiveUsers        int `json:"totalInactiveUsers"`
		TotalReservations         int `json:"totalReservations"`
		TotalMonthlyPlans         int `json:"totalMonthlyPlans"`
		TotalAnnualPlans          int `json:"totalAnnualPlans"`
	} `json:"summary"`
}

func (c *Community) GetCommunityReport(from, to *time.Time, groupBy string) (*CommunityReportResponse, *errors.Error) {
	params := bllAdapter.CommunityReportParams{
		From:    from,
		To:      to,
		GroupBy: groupBy,
	}
	total, communities, err := c.Adapter.Community.GetCommunityReport(params)
	if err != nil {
		return nil, &errors.InternalServerError.Default
	}

	// Calcular métricas agregadas
	response := &CommunityReportResponse{
		Total:       total,
		Communities: communities,
	}

	for _, community := range communities {
		response.Summary.TotalActiveMemberships += community.ActiveMemberships
		response.Summary.TotalExpiredMemberships += community.ExpiredMemberships
		response.Summary.TotalCancelledMemberships += community.CancelledMemberships
		response.Summary.TotalActiveUsers += community.ActiveUsers
		response.Summary.TotalInactiveUsers += community.InactiveUsers
		response.Summary.TotalReservations += community.TotalReservations
		response.Summary.TotalMonthlyPlans += community.MonthlyPlans
		response.Summary.TotalAnnualPlans += community.AnnualPlans
	}

	return response, nil
}
