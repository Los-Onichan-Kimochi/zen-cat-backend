package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	_ "onichankimochi.com/astro_cat_backend/src/server/api/docs" // Import generated swagger docs
	auditMiddleware "onichankimochi.com/astro_cat_backend/src/server/api/middleware"
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

	// ===== PROTECTED ENDPOINTS (JWT Authentication required) =====

	// Protected auth endpoints
	a.Echo.GET("/me/", a.GetCurrentUser, mw.JWTMiddleware)

	auth := a.Echo.Group("/auth")
	auth.Use(mw.JWTMiddleware) // Apply JWT middleware to all auth routes
	auth.POST("/refresh/", a.RefreshToken)
	auth.POST("/logout/", a.Logout)

	// ===== ADMIN ONLY ENDPOINTS (Administrator role required) =====

	// Community endpoints (admin only)
	community := a.Echo.Group("/community")
	community.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	community.GET("/:communityId/", a.GetCommunity)
	community.GET("/:communityId/image/", a.GetCommunityWithImage)
	community.GET("/", a.FetchCommunities)
	community.POST("/", a.CreateCommunity)
	community.PATCH("/:communityId/", a.UpdateCommunity)
	community.DELETE("/:communityId/", a.DeleteCommunity)
	community.POST("/bulk-create/", a.BulkCreateCommunities)
	community.DELETE("/bulk-delete/", a.BulkDeleteCommunities)

	// Professional endpoints (admin only)
	professional := a.Echo.Group("/professional")
	professional.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	professional.GET("/:professionalId/", a.GetProfessional)
	professional.GET("/", a.FetchProfessionals)
	professional.POST("/", a.CreateProfessional)
	professional.PATCH("/:professionalId/", a.UpdateProfessional)
	professional.DELETE("/:professionalId/", a.DeleteProfessional)
	professional.POST("/bulk-create/", a.BulkCreateProfessionals)
	professional.DELETE("/bulk-delete/", a.BulkDeleteProfessionals)

	// Local endpoints (admin only)
	local := a.Echo.Group("/local")
	local.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	local.GET("/:localId/", a.GetLocal)
	local.GET("/", a.FetchLocals)
	local.POST("/", a.CreateLocal)
	local.PATCH("/:localId/", a.UpdateLocal)
	local.DELETE("/:localId/", a.DeleteLocal)
	local.POST("/bulk-create/", a.BulkCreateLocals)
	local.DELETE("/bulk-delete/", a.BulkDeleteLocals)

	// Plan endpoints (admin only)
	plan := a.Echo.Group("/plan")
	plan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	plan.GET("/:planId/", a.GetPlan)
	plan.GET("/", a.FetchPlans)
	plan.POST("/", a.CreatePlan)
	plan.PATCH("/:planId/", a.UpdatePlan)
	plan.DELETE("/:planId/", a.DeletePlan)
	plan.POST("/bulk-create/", a.BulkCreatePlans)
	plan.DELETE("/bulk-delete/", a.BulkDeletePlans)

	// User endpoints (admin only)
	user := a.Echo.Group("/user")
	user.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	user.GET("/:userId/", a.GetUser)
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

	// Service endpoints (admin only)
	service := a.Echo.Group("/service")
	service.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	service.GET("/:serviceId/", a.GetService)
	service.GET("/", a.FetchServices)
	service.POST("/", a.CreateService)
	service.PATCH("/:serviceId/", a.UpdateService)
	service.DELETE("/:serviceId/", a.DeleteService)
	service.DELETE("/bulk-delete/", a.BulkDeleteServices)

	// Session endpoints (admin only)
	session := a.Echo.Group("/session")
	session.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	session.GET("/:sessionId/", a.GetSession)
	session.GET("/", a.FetchSessions)
	session.POST("/", a.CreateSession)
	session.PATCH("/:sessionId/", a.UpdateSession)
	session.DELETE("/:sessionId/", a.DeleteSession)
	session.POST("/bulk/", a.BulkCreateSessions)
	session.DELETE("/bulk-delete/", a.BulkDeleteSessions)
	session.POST("/check-conflicts/", a.CheckSessionConflicts)
	session.POST("/availability/", a.GetDayAvailability)

	// CommunityPlan endpoints (admin only)
	communityPlan := a.Echo.Group("/community-plan")
	communityPlan.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	communityPlan.POST("/", a.CreateCommunityPlan)
	communityPlan.GET("/:communityId/:planId/", a.GetCommunityPlan)
	communityPlan.DELETE("/:communityId/:planId/", a.DeleteCommunityPlan)
	communityPlan.POST("/bulk-create/", a.BulkCreateCommunityPlans)
	communityPlan.GET("/", a.FetchCommunityPlans)
	communityPlan.DELETE("/bulk-delete/", a.BulkDeleteCommunityPlans)

	// CommunityService endpoints (admin only)
	communityService := a.Echo.Group("/community-service")
	communityService.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	communityService.POST("/", a.CreateCommunityService)
	communityService.GET("/:communityId/:serviceId/", a.GetCommunityService)
	communityService.GET("/:communityId/", a.GetServicesByCommunityId)
	communityService.DELETE("/:communityId/:serviceId/", a.DeleteCommunityService)
	communityService.POST("/bulk-create/", a.BulkCreateCommunityServices)
	communityService.GET("/", a.FetchCommunityServices)
	communityService.DELETE("/bulk-delete/", a.BulkDeleteCommunityServices)

	// ServiceLocal endpoints (admin only)
	serviceLocal := a.Echo.Group("/service-local")
	serviceLocal.Use(mw.JWTMiddleware, mw.AdminOnlyMiddleware) // Apply JWT + Admin middleware
	serviceLocal.POST("/", a.CreateServiceLocal)
	serviceLocal.GET("/:serviceId/:localId/", a.GetServiceLocal)
	serviceLocal.DELETE("/:serviceId/:localId/", a.DeleteServiceLocal)
	serviceLocal.POST("/bulk/", a.BulkCreateServiceLocals)
	serviceLocal.GET("/", a.FetchServiceLocals)
	serviceLocal.DELETE("/bulk/", a.BulkDeleteServiceLocals)

	// ServiceProfessional endpoints (admin only)
	serviceProfessional := a.Echo.Group("/service-professional")
	serviceProfessional.Use(
		mw.JWTMiddleware,
		mw.AdminOnlyMiddleware,
	) // Apply JWT + Admin middleware
	serviceProfessional.POST("/", a.CreateServiceProfessional)
	serviceProfessional.GET("/:serviceId/:professionalId/", a.GetServiceProfessional)
	serviceProfessional.DELETE("/:serviceId/:professionalId/", a.DeleteServiceProfessional)
	serviceProfessional.POST("/bulk/", a.BulkCreateServiceProfessionals)
	serviceProfessional.GET("/", a.FetchServiceProfessionals)
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
	membership.POST("/", a.CreateMembership)
	membership.POST("/user/:userId/", a.CreateMembershipForUser)
	membership.PATCH("/:membershipId/", a.UpdateMembership)
	membership.DELETE("/:membershipId/", a.DeleteMembership)

	// Reservation endpoints (client only)
	reservation := a.Echo.Group("/reservation")
	reservation.Use(mw.JWTMiddleware, mw.ClientOnlyMiddleware) // Apply JWT + Client middleware
	reservation.GET("/:reservationId/", a.GetReservation)
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
