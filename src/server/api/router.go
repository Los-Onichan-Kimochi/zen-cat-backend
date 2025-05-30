package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "onichankimochi.com/astro_cat_backend/src/server/api/docs" // Import generated swagger docs
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// HealthCheck 			godoc
// @Summary 			Health Check
// @Description 		Verify connection in swagger
// @Tags 				Health Check
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} 	string "ok"
// @Router 				/health-check/ [get]
func (a *Api) HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, "Works well !!")
}

func (a *Api) RunApi(envSettings *schemas.EnvSettings) {
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // TODO: allow only FrontOffice and UserFront origins
		AllowCredentials: true,
	}
	a.Echo.Use(middleware.CORSWithConfig(corsConfig))

	if envSettings.EnableSwagger {
		a.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("server")))
	}

	healthCheck := a.Echo.Group("/health-check")
	healthCheck.GET("/", a.HealthCheck)

	// Community endpoints
	community := a.Echo.Group("/community")
	community.GET("/:communityId/", a.GetCommunity)
	community.GET("/", a.FetchCommunities)
	community.POST("/", a.CreateCommunity)
	community.PATCH("/:communityId/", a.UpdateCommunity)
	community.DELETE("/:communityId/", a.DeleteCommunity)
	community.POST("/bulk/", a.BulkCreateCommunities)

	// Professional endpoints
	professional := a.Echo.Group("/professional")
	professional.GET("/:professionalId/", a.GetProfessional)
	professional.GET("/", a.FetchProfessionals)
	professional.POST("/", a.CreateProfessional)
	professional.PATCH("/:professionalId/", a.UpdateProfessional)
	professional.DELETE("/:professionalId/", a.DeleteProfessional)
	professional.POST("/bulk-create/", a.BulkCreateProfessionals)
	professional.DELETE("/bulk-delete/", a.BulkDeleteProfessionals)

	// Local endpoints
	local := a.Echo.Group("/local")
	local.GET("/:localId/", a.GetLocal)
	local.GET("/", a.FetchLocals)
	local.POST("/", a.CreateLocal)
	local.PATCH("/:localId/", a.UpdateLocal)
	local.DELETE("/:localId/", a.DeleteLocal)

	// Plan endpoints
	plan := a.Echo.Group("/plan")
	plan.GET("/:planId/", a.GetPlan)
	plan.GET("/", a.FetchPlans)
	plan.POST("/", a.CreatePlan)
	plan.PATCH("/:planId/", a.UpdatePlan)
	plan.DELETE("/:planId/", a.DeletePlan)

	// User endpoints
	user := a.Echo.Group("/user")
	user.GET("/:userId/", a.GetUser)
	user.GET("/", a.FetchUsers)
	user.POST("/", a.CreateUser)
	user.PATCH("/:userId/", a.UpdateUser)
	user.DELETE("/:userId/", a.DeleteUser)

	// Service Endpoints
	service := a.Echo.Group("/service")
	service.GET("/:serviceId/", a.GetService)
	service.GET("/", a.FetchServices)
	service.POST("/", a.CreateService)
	service.PATCH("/:serviceId/", a.UpdateService)
	service.DELETE("/:serviceId/", a.DeleteService)

	// Session endpoints
	session := a.Echo.Group("/session")
	session.GET("/:sessionId/", a.GetSession)
	session.GET("/", a.FetchSessions)
	session.POST("/", a.CreateSession)
	session.PATCH("/:sessionId/", a.UpdateSession)
	session.DELETE("/:sessionId/", a.DeleteSession)
	session.POST("/bulk/", a.BulkCreateSessions)
	session.DELETE("/bulk-delete/", a.BulkDeleteSessions)

	// Reservation endpoints (read-only)
	reservation := a.Echo.Group("/reservation")
	reservation.GET("/:reservationId/", a.GetReservation)
	reservation.GET("/", a.FetchReservations)

	// CommunityPlan endpoints
	communityPlan := a.Echo.Group("/community-plan")
	communityPlan.POST("/", a.CreateCommunityPlan)
	communityPlan.GET("/:communityId/:planId/", a.GetCommunityPlan)
	communityPlan.DELETE("/:communityId/:planId/", a.DeleteCommunityPlan)
	communityPlan.POST("/bulk/", a.BulkCreateCommunityPlans)
	communityPlan.GET("/", a.FetchCommunityPlans)
	communityPlan.DELETE("/bulk/", a.BulkDeleteCommunityPlans)

	// CommunityService endpoints
	communityService := a.Echo.Group("/community-service")
	communityService.POST("/", a.CreateCommunityService)
	communityService.GET("/:communityId/:serviceId/", a.GetCommunityService)
	communityService.DELETE("/:communityId/:serviceId/", a.DeleteCommunityService)
	communityService.POST("/bulk/", a.BulkCreateCommunityServices)
	communityService.GET("/", a.FetchCommunityServices)
	communityService.DELETE("/bulk/", a.BulkDeleteCommunityServices)

	// Start the server
	a.Logger.Infoln(fmt.Sprintf("AstroCat server running on port %s", a.EnvSettings.MainPort))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", a.EnvSettings.MainPort)))
}
