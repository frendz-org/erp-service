package router

import (
	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"
	"erp-service/iam/auth"

	"github.com/gofiber/fiber/v2"
)

func SetupMasterdataRoutes(api fiber.Router, cfg *config.Config, mc *controller.MasterdataController, blacklistStore ...auth.TokenBlacklistStore) {
	masterdata := api.Group("/masterdata")

	publicRoutes := masterdata.Group("")
	{
		publicRoutes.Get("/categories", mc.ListCategories)
		publicRoutes.Get("/categories/code/:code", mc.GetCategoryByCode)
		publicRoutes.Get("/categories/:id", mc.GetCategoryByID)
		publicRoutes.Get("/categories/:id/children", mc.GetCategoryChildren)

		publicRoutes.Get("/items", mc.ListItems)
		publicRoutes.Get("/items/tree/:categoryCode", mc.GetItemTree)
		publicRoutes.Get("/items/by-parent/:categoryCode/:parentCode", mc.ListItemsByParent)
		publicRoutes.Get("/items/default/:categoryCode", mc.GetDefaultItem)
		publicRoutes.Get("/items/code/:categoryCode/:itemCode", mc.GetItemByCode)
		publicRoutes.Get("/items/:id", mc.GetItemByID)
		publicRoutes.Get("/items/:id/children", mc.GetItemChildrenByCode)

		publicRoutes.Post("/validate", mc.ValidateCode)
		publicRoutes.Post("/validate/batch", mc.ValidateCodes)
	}

	protectedRoutes := masterdata.Group("")
	protectedRoutes.Use(middleware.JWTAuth(cfg, blacklistStore...))
	protectedRoutes.Use(middleware.RequirePlatformAdmin())
	{

		protectedRoutes.Post("/categories", mc.CreateCategory)
		protectedRoutes.Put("/categories/:id", mc.UpdateCategory)
		protectedRoutes.Delete("/categories/:id", mc.DeleteCategory)

		protectedRoutes.Post("/items", mc.CreateItem)
		protectedRoutes.Put("/items/:id", mc.UpdateItem)
		protectedRoutes.Delete("/items/:id", mc.DeleteItem)
	}
}
