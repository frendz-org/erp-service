package router

import (
	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"
	"erp-service/iam/auth/contract"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(api fiber.Router, cfg *config.Config, authController *controller.AuthController, blacklistStore contract.TokenBlacklistStore) {
	auth := api.Group("/auth")
	auth.Use(middleware.JWTAuth(cfg, blacklistStore))
	auth.Post("/logout", authController.Logout)
	auth.Post("/logout-all", authController.LogoutAll)

	refreshToken := api.Group("/auth")
	refreshToken.Post("/refresh-token", authController.RefreshToken)

	registrations := api.Group("/registrations")
	registrations.Post("", authController.InitiateRegistration)
	registrations.Post("/:id/verify-otp", authController.VerifyRegistrationOTP)
	registrations.Post("/:id/resend-otp", authController.ResendRegistrationOTP)
	registrations.Get("/:id/status", authController.GetRegistrationStatus)
	registrations.Post("/:id/set-password", authController.SetPassword)
	registrations.Post("/:id/complete-profile", authController.CompleteProfileRegistration)
	registrations.Post("/:id/complete", authController.CompleteRegistration)

	login := api.Group("/login")
	login.Post("", authController.InitiateLogin)
	login.Post("/:id/verify-otp", authController.VerifyLoginOTP)
	login.Post("/:id/resend-otp", authController.ResendLoginOTP)
	login.Get("/:id/status", authController.GetLoginStatus)
}
