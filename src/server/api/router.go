package api

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "onichankimochi.com/astro_cat_backend/src/server/api/docs" // Import generated swagger docs
	auditMiddleware "onichankimochi.com/astro_cat_backend/src/server/api/middleware"
	"onichankimochi.com/astro_cat_backend/src/server/schemas"
)

// HealthCheck godoc
// @Summary 			Health Check
// @Description 		Check the health status of the API
// @Tags 				Health
// @Accept 				json
// @Produce 			json
// @Success 			200 {object} map[string]string "API is healthy"
// @Router 				/health-check/ [get]
func (a *Api) HealthCheck(c echo.Context) error {
	return c.JSON(200, map[string]string{
		"status": "OK",
	})
}

func (a *Api) RegisterRoutes(envSettings *schemas.EnvSettings) {
	// CORS
	corsConfig := middleware.CORSConfig{
		AllowOrigins:     []string{"*"}, // TODO: allow only FrontOffice and UserFront origins
		AllowCredentials: true,
	}
	a.Echo.Use(middleware.CORSWithConfig(corsConfig))

	// Add audit middleware (disabled for tests)
	mw := auditMiddleware.NewMiddleware(a.Logger, a.BllController, a.EnvSettings, a.Echo)
	if !a.EnvSettings.DisableAuthForTests {
		a.Echo.Use(mw.AuditMiddleware)
	}

	if envSettings.EnableSwagger {
		a.Echo.GET("/swagger/*", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("server")))
	}

	// ===== PUBLIC ENDPOINTS (No authentication required) =====
	healthCheck := a.Echo.Group("/health-check")
	healthCheck.GET("/", a.HealthCheck)

	// Authentication endpoints (public)
	a.Echo.POST("/login/", a.Login)
	a.Echo.POST("/register/", a.Register)
	a.Echo.POST("/forgot-password/", a.ForgotPassword)
	a.Echo.POST("/login/google/", a.GoogleLogin)

	// Contact endpoints (public)
	a.Echo.POST("/contact", a.ContactMessage)

	// Public browsing endpoints (for both authenticated and unauthenticated users)
	// Communities
	a.Echo.GET("/community/", a.FetchCommunities)
	a.Echo.GET("/community/:communityId/", a.GetCommunity)
	a.Echo.GET("/community/:communityId/image/", a.GetCommunityWithImage)

	// Services
	a.Echo.GET("/service/", a.FetchServices)
	a.Echo.GET("/service/:serviceId/", a.GetService)
	a.Echo.GET("/service/:serviceId/image/", a.GetServiceWithImage)

	// Plans
	a.Echo.GET("/plan/", a.FetchPlans)
	a.Echo.GET("/plan/:planId/", a.GetPlan)

	// Locals
	a.Echo.GET("/local/", a.FetchLocals)
	a.Echo.GET("/local/:localId/", a.GetLocal)
	a.Echo.GET("/local/:localId/image/", a.GetLocalWithImage)

	// Professionals
	a.Echo.GET("/professional/", a.FetchProfessionals)
	a.Echo.GET("/professional/:professionalId/", a.GetProfessional)
	a.Echo.GET("/professional/:professionalId/image/", a.GetProfessionalWithImage)

	// Sessions
	a.Echo.GET("/session/", a.FetchSessions)
	a.Echo.GET("/session/:sessionId/", a.GetSession)

	// Community Services
	a.Echo.GET("/community-service/", a.FetchCommunityServices)
	a.Echo.GET("/community-service/:communityId/", a.GetServicesByCommunityId)
	a.Echo.GET("/community-service/id/:id/", a.GetCommunityServiceById)

	// Community Plans
	a.Echo.GET("/community-plan/", a.FetchCommunityPlans)

	// Service Locals
	a.Echo.GET("/service-local/", a.FetchServiceLocals)

	// Service Professionals
	a.Echo.GET("/service-professional/", a.FetchServiceProfessionals)

	// ===== AUTHENTICATED ENDPOINTS (JWT required, any authenticated user) =====

	// Current user info
	a.Echo.GET("/me/", a.GetCurrentUser, mw.JWTMiddleware)

	// Auth management (authenticated users)
	auth := a.Echo.Group("/auth")
	auth.Use(mw.JWTMiddleware)
	auth.POST("/refresh/", a.RefreshToken)
	auth.POST("/logout/", a.Logout)

	// ===== ADMIN + CLIENT MIXED ENDPOINTS (Both roles can access) =====

	// Session availability and conflicts (both admin and client need this)
	sessionMixed := a.Echo.Group("/session")
	sessionMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware) // Admin or Client required
	sessionMixed.POST("/check-conflicts/", a.CheckSessionConflicts)
	sessionMixed.POST("/availability/", a.GetDayAvailability)

	// Community Plan associations (both admin and client need to read)
	communityPlanMixed := a.Echo.Group("/community-plan")
	communityPlanMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware) // Admin or Client required
	communityPlanMixed.GET("/:communityId/:planId/", a.GetCommunityPlan)

	// Community Service associations (both admin and client need to read)
	communityServiceMixed := a.Echo.Group("/community-service")
	communityServiceMixed.Use(
		mw.JWTMiddleware,
		mw.AdminOrClientMiddleware,
	) // Admin or Client required
	communityServiceMixed.GET("/:communityId/:serviceId/", a.GetCommunityService)

	// Service Local associations (both admin and client need to read)
	serviceLocalMixed := a.Echo.Group("/service-local")
	serviceLocalMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware) // Admin or Client required
	serviceLocalMixed.GET("/:serviceId/:localId/", a.GetServiceLocal)

	// Service Professional associations (both admin and client need to read)
	serviceProfessionalMixed := a.Echo.Group("/service-professional")
	serviceProfessionalMixed.Use(
		mw.JWTMiddleware,
		mw.AdminOrClientMiddleware,
	) // Admin or Client required
	serviceProfessionalMixed.GET("/:serviceId/:professionalId/", a.GetServiceProfessional)

	// Membership endpoints that both admin and client need
	membershipMixed := a.Echo.Group("/membership")
	membershipMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware) // Admin or Client required
	membershipMixed.GET("/community/:communityId/users", a.GetUsersByCommunityId)
	membershipMixed.GET("/user/:userId/community/:communityId", a.GetMembershipByUserAndCommunity)

	// Reservation endpoints that both admin and client need
	reservationMixed := a.Echo.Group("/reservation")
	reservationMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware) // Admin or Client required
	reservationMixed.GET("/:reservationId/", a.GetReservation)
	reservationMixed.GET("/", a.FetchReservations)
	reservationMixed.GET("/:communityId/:userId/", a.GetReservationsByCommunityIdByUserId)
	reservationMixed.POST("/", a.CreateReservation)
	reservationMixed.PATCH("/:reservationId/", a.UpdateReservation)
	reservationMixed.DELETE("/:reservationId/", a.DeleteReservation)
	reservationMixed.DELETE("/bulk-delete/", a.BulkDeleteReservations)

	// ===== ADMIN ONLY ENDPOINTS (Administrator role required) =====

	// Community management (admin only)
	community := a.Echo.Group("/community")
	community.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	community.POST("/", a.CreateCommunity)
	community.PATCH("/:communityId/", a.UpdateCommunity)
	community.DELETE("/:communityId/", a.DeleteCommunity)
	community.POST("/bulk-create/", a.BulkCreateCommunities)
	community.DELETE("/bulk-delete/", a.BulkDeleteCommunities)

	// Professional management (admin only)
	professional := a.Echo.Group("/professional")
	professional.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	professional.POST("/", a.CreateProfessional)
	professional.PATCH("/:professionalId/", a.UpdateProfessional)
	professional.DELETE("/:professionalId/", a.DeleteProfessional)
	professional.POST("/bulk-create/", a.BulkCreateProfessionals)
	professional.DELETE("/bulk-delete/", a.BulkDeleteProfessionals)

	// Local management (admin only)
	local := a.Echo.Group("/local")
	local.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	local.POST("/", a.CreateLocal)
	local.PATCH("/:localId/", a.UpdateLocal)
	local.DELETE("/:localId/", a.DeleteLocal)
	local.POST("/bulk-create/", a.BulkCreateLocals)
	local.DELETE("/bulk-delete/", a.BulkDeleteLocals)

	// Plan management (admin only)
	plan := a.Echo.Group("/plan")
	plan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	plan.POST("/", a.CreatePlan)
	plan.PATCH("/:planId/", a.UpdatePlan)
	plan.DELETE("/:planId/", a.DeletePlan)
	plan.POST("/bulk-create/", a.BulkCreatePlans)
	plan.DELETE("/bulk-delete/", a.BulkDeletePlans)

	// User management (admin only)
	user := a.Echo.Group("/user")
	user.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	user.GET("/:userId/", a.GetUser)
	user.GET("/:userId/image/", a.GetUserWithImage)
	user.GET("/", a.FetchUsers)
	user.GET("/exists", a.CheckUserExists)
	user.POST("/", a.CreateUser)
	user.PATCH("/:userId/", a.UpdateUser)
	user.DELETE("/:userId/", a.DeleteUser)
	user.POST("/bulk-create/", a.BulkCreateUsers)
	user.DELETE("/bulk-delete/", a.BulkDeleteUsers)
	user.PATCH("/:userId/role/", a.ChangeUserRole)
	user.GET("/stats/", a.GetUserStats)

	// User management (admin and client)
	userMixed := a.Echo.Group("/user")
	userMixed.Use(mw.JWTMiddleware, mw.AdminOrClientMiddleware)
	userMixed.POST("/change-password/", a.ChangePassword)

	// Service management (admin only)
	service := a.Echo.Group("/service")
	service.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	service.POST("/", a.CreateService)
	service.PATCH("/:serviceId/", a.UpdateService)
	service.DELETE("/:serviceId/", a.DeleteService)
	service.DELETE("/bulk-delete/", a.BulkDeleteServices)

	// Session management (admin only)
	session := a.Echo.Group("/session")
	session.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	session.POST("/", a.CreateSession)
	session.PATCH("/:sessionId/", a.UpdateSession)
	session.DELETE("/:sessionId/", a.DeleteSession)
	session.POST("/bulk/", a.BulkCreateSessions)
	session.DELETE("/bulk-delete/", a.BulkDeleteSessions)

	// Community Plan management (admin only)
	communityPlan := a.Echo.Group("/community-plan")
	communityPlan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	communityPlan.POST("/", a.CreateCommunityPlan)
	communityPlan.DELETE("/:communityId/:planId/", a.DeleteCommunityPlan)
	communityPlan.POST("/bulk-create/", a.BulkCreateCommunityPlans)
	communityPlan.DELETE("/bulk-delete/", a.BulkDeleteCommunityPlans)

	// Community Service management (admin only)
	communityService := a.Echo.Group("/community-service")
	communityService.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	communityService.POST("/", a.CreateCommunityService)
	communityService.DELETE("/:communityId/:serviceId/", a.DeleteCommunityService)
	communityService.POST("/bulk-create/", a.BulkCreateCommunityServices)
	communityService.DELETE("/bulk-delete/", a.BulkDeleteCommunityServices)

	// Service Local management (admin only)
	serviceLocal := a.Echo.Group("/service-local")
	serviceLocal.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	serviceLocal.POST("/", a.CreateServiceLocal)
	serviceLocal.DELETE("/:serviceId/:localId/", a.DeleteServiceLocal)
	serviceLocal.POST("/bulk/", a.BulkCreateServiceLocals)
	serviceLocal.DELETE("/bulk/", a.BulkDeleteServiceLocals)

	// Service Professional management (admin only)
	serviceProfessional := a.Echo.Group("/service-professional")
	serviceProfessional.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	serviceProfessional.POST("/", a.CreateServiceProfessional)
	serviceProfessional.DELETE("/:serviceId/:professionalId/", a.DeleteServiceProfessional)
	serviceProfessional.POST("/bulk/", a.BulkCreateServiceProfessionals)
	serviceProfessional.DELETE("/bulk/", a.BulkDeleteServiceProfessionals)

	// Audit Log management (admin only)
	auditLog := a.Echo.Group("/audit-log")
	auditLog.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	auditLog.GET("/", a.GetAuditLogs)
	auditLog.GET("/:auditLogId/", a.GetAuditLogById)
	auditLog.GET("/stats/", a.GetAuditStats)
	auditLog.DELETE("/cleanup/", a.DeleteOldAuditLogs)

	// Error Log management (admin only)
	errorLog := a.Echo.Group("/error-log")
	errorLog.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	errorLog.GET("/", a.GetErrorLogs)
	errorLog.GET("/:auditLogId/", a.GetErrorLogById)
	errorLog.GET("/stats/", a.GetErrorStats)

	// Reports (admin only)
	reports := a.Echo.Group("/reports")
	reports.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	reports.GET("/services", a.GetServiceReport)
	reports.GET("/communities", a.GetCommunityReport)

	// ===== CLIENT ONLY ENDPOINTS (Client role required) =====

	// Onboarding (client only)
	onboarding := a.Echo.Group("/onboarding")
	onboarding.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware)
	onboarding.GET("/:onboardingId/", a.GetOnboarding)
	onboarding.GET("/", a.FetchOnboardings)
	onboarding.GET("/user/:userId/", a.GetOnboardingByUserId)
	onboarding.POST("/user/:userId/", a.CreateOnboardingForUser)
	onboarding.PATCH("/:onboardingId/", a.UpdateOnboarding)
	onboarding.PATCH("/user/:userId/", a.UpdateOnboardingByUserId)
	onboarding.DELETE("/:onboardingId/", a.DeleteOnboarding)
	onboarding.DELETE("/user/:userId/", a.DeleteOnboardingByUserId)

	// Membership management (client only)
	membership := a.Echo.Group("/membership")
	membership.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware)
	membership.GET("/:membershipId/", a.GetMembership)
	membership.GET("/", a.FetchMemberships)
	membership.GET("/user/:userId/", a.GetMembershipsByUserId)
	membership.GET("/community/:communityId/", a.GetMembershipsByCommunityId)
	membership.POST("/", a.CreateMembership)
	membership.POST("/user/:userId/", a.CreateMembershipForUser)
	membership.PATCH("/:membershipId/", a.UpdateMembership)
	membership.DELETE("/:membershipId/", a.DeleteMembership)
}

func (a *Api) RunApi(envSettings *schemas.EnvSettings) {
	a.RegisterRoutes(envSettings)

	// Start the server
	a.Logger.Infoln(fmt.Sprintf("AstroCat server running on port %s", a.EnvSettings.MainPort))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", a.EnvSettings.MainPort)))
}
