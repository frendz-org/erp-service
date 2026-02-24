package router

import (
	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"
	"erp-service/iam/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authController *controller.AuthController, blacklistStore auth.TokenBlacklistStore) {
	api.Post("/auth/refresh-token", authController.RefreshToken)

	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuth(cfg, blacklistStore))
	auth.Post("/logout", authController.Logout)
	auth.Post("/logout-all", authController.LogoutAll)

	registrations := api.Group("/registrations")
	registrations.Post("", authController.InitiateRegistration)
	registrations.Post("/:id/verify-otp", authController.VerifyRegistrationOTP)
	registrations.Post("/:id/resend-otp", authController.ResendRegistrationOTP)
	registrations.Get("/:id/status", authController.GetRegistrationStatus)
	registrations.Post("/:id/set-password", authController.SetPassword)
	registrations.Post("/:id/complete-profile", authController.CompleteProfileRegistration)

	login := api.Group("/login")
	login.Post("", authController.InitiateLogin)
	login.Post("/:id/verify-otp", authController.VerifyLoginOTP)
	login.Post("/:id/resend-otp", authController.ResendLoginOTP)
	login.Get("/:id/status", authController.GetLoginStatus)
}
