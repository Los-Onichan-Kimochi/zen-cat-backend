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

	// ===== PUBLIC ENDPOINTS (No authentication required) =====
	healthCheck := a.Echo.Group("/health-check")
	healthCheck.GET("/", a.HealthCheck)

	// Login endpoints (public)
	a.Echo.POST("/login/", a.Login)
	a.Echo.POST("/register/", a.Register)

	// ===== PROTECTED ENDPOINTS (JWT Authentication required) =====

	// Protected auth endpoints
	a.Echo.GET("/me/", a.GetCurrentUser, a.JWTMiddleware)

	auth := a.Echo.Group("/auth")
	auth.Use(a.JWTMiddleware) // Apply JWT middleware to all auth routes
	auth.POST("/refresh/", a.RefreshToken)

	// Community endpoints (all protected)
	community := a.Echo.Group("/community")
	community.Use(a.JWTMiddleware) // Apply JWT middleware to all community routes
	community.GET("/:communityId/", a.GetCommunity)
	community.GET("/", a.FetchCommunities)
	community.POST("/", a.CreateCommunity)
	community.PATCH("/:communityId/", a.UpdateCommunity)
	community.DELETE("/:communityId/", a.DeleteCommunity)
	community.POST("/bulk/", a.BulkCreateCommunities)

	// Professional endpoints (all protected)
	professional := a.Echo.Group("/professional")
	professional.Use(a.JWTMiddleware) // Apply JWT middleware to all professional routes
	professional.GET("/:professionalId/", a.GetProfessional)
	professional.GET("/", a.FetchProfessionals)
	professional.POST("/", a.CreateProfessional)
	professional.PATCH("/:professionalId/", a.UpdateProfessional)
	professional.DELETE("/:professionalId/", a.DeleteProfessional)
	professional.POST("/bulk-create/", a.BulkCreateProfessionals)
	professional.DELETE("/bulk-delete/", a.BulkDeleteProfessionals)

	// Local endpoints (all protected)
	local := a.Echo.Group("/local")
	local.Use(a.JWTMiddleware) // Apply JWT middleware to all local routes
	local.GET("/:localId/", a.GetLocal)
	local.GET("/", a.FetchLocals)
	local.POST("/", a.CreateLocal)
	local.PATCH("/:localId/", a.UpdateLocal)
	local.DELETE("/:localId/", a.DeleteLocal)

	// Plan endpoints (all protected)
	plan := a.Echo.Group("/plan")
	plan.Use(a.JWTMiddleware) // Apply JWT middleware to all plan routes
	plan.GET("/:planId/", a.GetPlan)
	plan.GET("/", a.FetchPlans)
	plan.POST("/", a.CreatePlan)
	plan.PATCH("/:planId/", a.UpdatePlan)
	plan.DELETE("/:planId/", a.DeletePlan)

	// User endpoints (all protected)
	user := a.Echo.Group("/user")
	user.Use(a.JWTMiddleware) // Apply JWT middleware to all user routes
	user.GET("/:userId/", a.GetUser)
	user.GET("/", a.FetchUsers)
	user.POST("/", a.CreateUser)
	user.PATCH("/:userId/", a.UpdateUser)
	user.DELETE("/:userId/", a.DeleteUser)
	user.POST("/bulk-create/", a.BulkCreateUsers)
	user.DELETE("/bulk-delete/", a.BulkDeleteUsers)

	// Service Endpoints (all protected)
	service := a.Echo.Group("/service")
	service.Use(a.JWTMiddleware) // Apply JWT middleware to all service routes
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
	session.POST("/check-conflicts/", a.CheckSessionConflicts)
	session.POST("/availability/", a.GetDayAvailability)

	// Reservation endpoints (read-only)
	reservation := a.Echo.Group("/reservation")
	reservation.GET("/:reservationId/", a.GetReservation)
	reservation.GET("/", a.FetchReservations)

	// CommunityPlan endpoints (all protected)
	communityPlan := a.Echo.Group("/community-plan")
	communityPlan.Use(a.JWTMiddleware) // Apply JWT middleware to all community-plan routes
	communityPlan.POST("/", a.CreateCommunityPlan)
	communityPlan.GET("/:communityId/:planId/", a.GetCommunityPlan)
	communityPlan.DELETE("/:communityId/:planId/", a.DeleteCommunityPlan)
	communityPlan.POST("/bulk/", a.BulkCreateCommunityPlans)
	communityPlan.GET("/", a.FetchCommunityPlans)
	communityPlan.DELETE("/bulk/", a.BulkDeleteCommunityPlans)

	// CommunityService endpoints (all protected)
	communityService := a.Echo.Group("/community-service")
	communityService.Use(a.JWTMiddleware) // Apply JWT middleware to all community-service routes
	communityService.POST("/", a.CreateCommunityService)
	communityService.GET("/:communityId/:serviceId/", a.GetCommunityService)
	communityService.DELETE("/:communityId/:serviceId/", a.DeleteCommunityService)
	communityService.POST("/bulk/", a.BulkCreateCommunityServices)
	communityService.GET("/", a.FetchCommunityServices)
	communityService.DELETE("/bulk/", a.BulkDeleteCommunityServices)

	// ServiceLocal endpoints
	serviceLocal := a.Echo.Group("/service-local")
	serviceLocal.POST("/", a.CreateServiceLocal)
	serviceLocal.GET("/:serviceId/:localId/", a.GetServiceLocal)
	serviceLocal.DELETE("/:serviceId/:localId/", a.DeleteServiceLocal)
	serviceLocal.POST("/bulk/", a.BulkCreateServiceLocals)
	serviceLocal.GET("/", a.FetchServiceLocals)
	serviceLocal.DELETE("/bulk/", a.BulkDeleteServiceLocals)

	// ServiceProfessional endpoints
	serviceProfessional := a.Echo.Group("/service-professional")
	serviceProfessional.POST("/", a.CreateServiceProfessional)
	serviceProfessional.GET("/:serviceId/:professionalId/", a.GetServiceProfessional)
	serviceProfessional.DELETE("/:serviceId/:professionalId/", a.DeleteServiceProfessional)
	serviceProfessional.POST("/bulk/", a.BulkCreateServiceProfessionals)
	serviceProfessional.GET("/", a.FetchServiceProfessionals)
	serviceProfessional.DELETE("/bulk/", a.BulkDeleteServiceProfessionals)

	// Start the server
	a.Logger.Infoln(fmt.Sprintf("AstroCat server running on port %s", a.EnvSettings.MainPort))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", a.EnvSettings.MainPort)))
}
