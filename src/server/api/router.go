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

	// Login endpoints (public)
	a.Echo.POST("/login/", a.Login)
	a.Echo.POST("/register/", a.Register)
	a.Echo.POST("/forgot-password/", a.ForgotPassword)
	a.Echo.POST("/login/google/", a.GoogleLogin)

	// Contact endpoints (public)
	a.Echo.POST("/contact", a.ContactMessage)

	// Public Community endpoints (for browsing)
	a.Echo.GET("/community/", a.FetchCommunities)
	a.Echo.GET("/community/:communityId/", a.GetCommunity)
	a.Echo.GET("/community/:communityId/image/", a.GetCommunityWithImage)

	// Public Service endpoints (for browsing)
	a.Echo.GET("/service/", a.FetchServices)
	a.Echo.GET("/service/:serviceId/", a.GetService)
	a.Echo.GET("/service/:serviceId/image/", a.GetServiceWithImage)

	// Public Plan endpoints (for browsing)
	a.Echo.GET("/plan/", a.FetchPlans)
	a.Echo.GET("/plan/:planId/", a.GetPlan)

	// Public CommunityService endpoints (for browsing)
	a.Echo.GET("/community-service/", a.FetchCommunityServices)
	a.Echo.GET("/community-service/:communityId/", a.GetServicesByCommunityId)
	a.Echo.GET("/community-service/id/:id/", a.GetCommunityServiceById)

	// Public CommunityPlan endpoints (for browsing)
	a.Echo.GET("/community-plan/", a.FetchCommunityPlans)

	// Public Local endpoints (for browsing)
	a.Echo.GET("/local/", a.FetchLocals)
	a.Echo.GET("/local/:localId/", a.GetLocal)
	a.Echo.GET("/local/:localId/image/", a.GetLocalWithImage)

	// Public Professional endpoints (for browsing)
	a.Echo.GET("/professional/", a.FetchProfessionals)
	a.Echo.GET("/professional/:professionalId/", a.GetProfessional)
	a.Echo.GET("/professional/:professionalId/image/", a.GetProfessionalWithImage)

	// Public ServiceLocal endpoints (for browsing)
	a.Echo.GET("/service-local/", a.FetchServiceLocals)

	// Public ServiceProfessional endpoints (for browsing)
	a.Echo.GET("/service-professional/", a.FetchServiceProfessionals)

	// Public Session endpoints (for browsing)
	a.Echo.GET("/session/", a.FetchSessions)
	a.Echo.GET("/session/:sessionId/", a.GetSession)

	// ===== PROTECTED ENDPOINTS (JWT Authentication required) =====

	// Protected auth endpoints
	a.Echo.GET("/me/", a.GetCurrentUser, mw.JWTMiddleware)

	auth := a.Echo.Group("/auth")
	auth.Use(mw.JWTMiddleware) // Apply JWT middleware to all auth routes
	auth.POST("/refresh/", a.RefreshToken)
	auth.POST("/logout/", a.Logout)

	// ===== ADMIN ONLY ENDPOINTS (Administrator role required) =====

	// Community endpoints (admin only - write operations)
	community := a.Echo.Group("/community")
	community.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	community.POST("/", a.CreateCommunity)
	community.PATCH("/:communityId/", a.UpdateCommunity)
	community.DELETE("/:communityId/", a.DeleteCommunity)
	community.POST("/bulk-create/", a.BulkCreateCommunities)
	community.DELETE("/bulk-delete/", a.BulkDeleteCommunities)

	// Professional endpoints (admin only - write operations)
	professional := a.Echo.Group("/professional")
	professional.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	professional.POST("/", a.CreateProfessional)
	professional.PATCH("/:professionalId/", a.UpdateProfessional)
	professional.DELETE("/:professionalId/", a.DeleteProfessional)
	professional.POST("/bulk-create/", a.BulkCreateProfessionals)
	professional.DELETE("/bulk-delete/", a.BulkDeleteProfessionals)

	// Local endpoints (admin only - write operations)
	local := a.Echo.Group("/local")
	local.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	local.POST("/", a.CreateLocal)
	local.PATCH("/:localId/", a.UpdateLocal)
	local.DELETE("/:localId/", a.DeleteLocal)
	local.POST("/bulk-create/", a.BulkCreateLocals)
	local.DELETE("/bulk-delete/", a.BulkDeleteLocals)

	// Plan endpoints (admin only - write operations)
	plan := a.Echo.Group("/plan")
	plan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	plan.POST("/", a.CreatePlan)
	plan.PATCH("/:planId/", a.UpdatePlan)
	plan.DELETE("/:planId/", a.DeletePlan)
	plan.POST("/bulk-create/", a.BulkCreatePlans)
	plan.DELETE("/bulk-delete/", a.BulkDeletePlans)

	// User endpoints (admin only)
	user := a.Echo.Group("/user")
	user.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	user.GET("/:userId/", a.GetUser)
	user.GET("/:userId/image/", a.GetUserWithImage)
	user.GET("/", a.FetchUsers)
	user.GET("/exists", a.CheckUserExists)
	user.POST("/", a.CreateUser)
	user.PATCH("/:userId/", a.UpdateUser)
	user.DELETE("/:userId/", a.DeleteUser)
	user.POST("/bulk-create/", a.BulkCreateUsers)
	user.DELETE("/bulk-delete/", a.BulkDeleteUsers)
	user.POST("/change-password/", a.ChangePassword)
	user.PATCH("/:userId/role/", a.ChangeUserRole)
	user.GET("/stats/", a.GetUserStats)

	// Service endpoints (admin only - write operations)
	service := a.Echo.Group("/service")
	service.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	service.POST("/", a.CreateService)
	service.PATCH("/:serviceId/", a.UpdateService)
	service.DELETE("/:serviceId/", a.DeleteService)
	service.DELETE("/bulk-delete/", a.BulkDeleteServices)

	// Session endpoints (admin only - write operations)
	session := a.Echo.Group("/session")
	session.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	session.POST("/", a.CreateSession)
	session.PATCH("/:sessionId/", a.UpdateSession)
	session.DELETE("/:sessionId/", a.DeleteSession)
	session.POST("/bulk/", a.BulkCreateSessions)
	session.DELETE("/bulk-delete/", a.BulkDeleteSessions)
	session.POST("/check-conflicts/", a.CheckSessionConflicts)
	session.POST("/availability/", a.GetDayAvailability)

	// CommunityPlan endpoints (admin only - write operations)
	communityPlan := a.Echo.Group("/community-plan")
	communityPlan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	communityPlan.POST("/", a.CreateCommunityPlan)
	communityPlan.GET("/:communityId/:planId/", a.GetCommunityPlan)
	communityPlan.DELETE("/:communityId/:planId/", a.DeleteCommunityPlan)
	communityPlan.POST("/bulk-create/", a.BulkCreateCommunityPlans)
	communityPlan.DELETE("/bulk-delete/", a.BulkDeleteCommunityPlans)

	// CommunityService endpoints (admin only - write operations)
	communityService := a.Echo.Group("/community-service")
	communityService.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	communityService.POST("/", a.CreateCommunityService)
	communityService.GET("/:communityId/:serviceId/", a.GetCommunityService)
	communityService.DELETE("/:communityId/:serviceId/", a.DeleteCommunityService)
	communityService.POST("/bulk-create/", a.BulkCreateCommunityServices)
	communityService.DELETE("/bulk-delete/", a.BulkDeleteCommunityServices)

	// ServiceLocal endpoints (admin only - write operations)
	serviceLocal := a.Echo.Group("/service-local")
	serviceLocal.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	serviceLocal.POST("/", a.CreateServiceLocal)
	serviceLocal.GET("/:serviceId/:localId/", a.GetServiceLocal)
	serviceLocal.GET("/:serviceId/", a.GetLocalsByServiceId)
	serviceLocal.DELETE("/:serviceId/:localId/", a.DeleteServiceLocal)
	serviceLocal.POST("/bulk/", a.BulkCreateServiceLocals)
	serviceLocal.DELETE("/bulk/", a.BulkDeleteServiceLocals)

	// ServiceProfessional endpoints (admin only - write operations)
	serviceProfessional := a.Echo.Group("/service-professional")
	serviceProfessional.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	serviceProfessional.POST("/", a.CreateServiceProfessional)
	serviceProfessional.GET("/:serviceId/:professionalId/", a.GetServiceProfessional)
	serviceProfessional.GET("/:serviceId/", a.GetProfessionalsByServiceId)
	serviceProfessional.DELETE("/:serviceId/:professionalId/", a.DeleteServiceProfessional)
	serviceProfessional.POST("/bulk/", a.BulkCreateServiceProfessionals)
	serviceProfessional.DELETE("/bulk/", a.BulkDeleteServiceProfessionals)

	// AuditLog endpoints (admin only)
	auditLog := a.Echo.Group("/audit-log")
	auditLog.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	auditLog.GET("/", a.GetAuditLogs)
	auditLog.GET("/:auditLogId/", a.GetAuditLogById)
	auditLog.GET("/stats/", a.GetAuditStats)
	auditLog.DELETE("/cleanup/", a.DeleteOldAuditLogs)

	// ErrorLog endpoints (admin only)
	errorLog := a.Echo.Group("/error-log")
	errorLog.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	errorLog.GET("/", a.GetErrorLogs)
	errorLog.GET("/:auditLogId/", a.GetErrorLogById)
	errorLog.GET("/stats/", a.GetErrorStats)

	// Reports endpoints (admin only)
	reports := a.Echo.Group("/reports")
	reports.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware)
	reports.GET("/services", a.GetServiceReport)
	reports.GET("/communities", a.GetCommunityReport)

	// ===== CLIENT ENDPOINTS (Client role required) =====

	// Onboarding endpoints (client only)
	onboarding := a.Echo.Group("/onboarding")
	onboarding.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware) // Apply JWT + Client middleware
	onboarding.GET("/:onboardingId/", a.GetOnboarding)
	onboarding.GET("/", a.FetchOnboardings)
	onboarding.GET("/user/:userId/", a.GetOnboardingByUserId)
	onboarding.POST("/user/:userId/", a.CreateOnboardingForUser)
	onboarding.PATCH("/:onboardingId/", a.UpdateOnboarding)
	onboarding.PATCH("/user/:userId/", a.UpdateOnboardingByUserId)
	onboarding.DELETE("/:onboardingId/", a.DeleteOnboarding)
	onboarding.DELETE("/user/:userId/", a.DeleteOnboardingByUserId)

	// Membership endpoints (client only)
	membership := a.Echo.Group("/membership")
	membership.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware) // Apply JWT + Client middleware
	membership.GET("/:membershipId/", a.GetMembership)
	membership.GET("/", a.FetchMemberships)
	membership.GET("/user/:userId/", a.GetMembershipsByUserId)
	membership.GET("/community/:communityId/", a.GetMembershipsByCommunityId)
	membership.GET("/community/:communityId/users", a.GetUsersByCommunityId)
	membership.GET("/user/:userId/community/:communityId", a.GetMembershipByUserAndCommunity)
	membership.POST("/", a.CreateMembership)
	membership.POST("/user/:userId/", a.CreateMembershipForUser)
	membership.PATCH("/:membershipId/", a.UpdateMembership)
	membership.DELETE("/:membershipId/", a.DeleteMembership)

	// Reservation endpoints (client only)
	reservation := a.Echo.Group("/reservation")
	reservation.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware) // Apply JWT + Client middleware
	reservation.GET("/:reservationId/", a.GetReservation)
	reservation.GET("/:communityId/:userId/", a.GetReservationsByCommunityIdByUserId)
	reservation.GET("/", a.FetchReservations)
	reservation.POST("/", a.CreateReservation)
	reservation.PATCH("/:reservationId/", a.UpdateReservation)
	reservation.DELETE("/:reservationId/", a.DeleteReservation)
	reservation.DELETE("/bulk-delete/", a.BulkDeleteReservations)
}

func (a *Api) RunApi(envSettings *schemas.EnvSettings) {
	a.RegisterRoutes(envSettings)

	// Start the server
	a.Logger.Infoln(fmt.Sprintf("AstroCat server running on port %s", a.EnvSettings.MainPort))
	a.Logger.Fatal(a.Echo.Start(fmt.Sprintf(":%s", a.EnvSettings.MainPort)))
}
