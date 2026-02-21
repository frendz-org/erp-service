package router

import (
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func SetupMemberRoutes(api fiber.Router, ctrl *controller.MemberController, jwtMiddleware fiber.Handler, frendzSavingMW fiber.Handler) {
	members := api.Group("/members")
	members.Use(jwtMiddleware)
	members.Use(middleware.ExtractTenantContext())
	members.Use(frendzSavingMW)

	members.Post("/register", ctrl.Register)

	adminMW := middleware.RequireProductRole("TENANT_PRODUCT_ADMIN")

	members.Get("/", adminMW, ctrl.List)
	members.Get("/:memberId", adminMW, ctrl.Get)
	members.Post("/:memberId/approve", adminMW, ctrl.Approve)
	members.Post("/:memberId/reject", adminMW, ctrl.Reject)
	members.Put("/:memberId/role", adminMW, ctrl.ChangeRole)
	members.Post("/:memberId/deactivate", adminMW, ctrl.Deactivate)
}
