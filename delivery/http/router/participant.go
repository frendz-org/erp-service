package router

import (
	"time"

	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func selfRegRateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			if uid, ok := c.Locals("userID").(string); ok && uid != "" {
				return "self-reg:" + uid
			}
			return "self-reg:" + c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"error":   "too many requests",
				"code":    "ERR_TOO_MANY_REQUESTS",
			})
		},
	})
}

func SetupParticipantRoutes(api fiber.Router, ctrl *controller.ParticipantController, jwtMiddleware fiber.Handler, frendzSavingMW fiber.Handler, cfg *config.Config) {
	isDev := cfg.IsDevelopment()

	selfReg := api.Group("/participants")
	selfReg.Use(jwtMiddleware)
	if isDev {
		selfReg.Post("/self-register", ctrl.SelfRegister)
	} else {
		selfReg.Post("/self-register", selfRegRateLimit(), ctrl.SelfRegister)
	}

	participants := api.Group("/participants")
	participants.Use(jwtMiddleware)
	participants.Use(middleware.ExtractTenantContext())
	participants.Use(frendzSavingMW)

	creatorMW := middleware.RequireProductRole("PARTICIPANT_CREATOR")
	approverMW := middleware.RequireProductRole("PARTICIPANT_APPROVER")
	anyRoleMW := middleware.RequireProductRole("PARTICIPANT_CREATOR", "PARTICIPANT_APPROVER")

	participants.Post("/", creatorMW, ctrl.Create)
	participants.Get("/", anyRoleMW, ctrl.List)
	participants.Get("/:id", anyRoleMW, ctrl.Get)

	participants.Put("/:id/personal-data", creatorMW, ctrl.UpdatePersonalData)

	participants.Put("/:id/identities", creatorMW, ctrl.SaveIdentity)
	participants.Delete("/:id/identities/:identityId", creatorMW, ctrl.DeleteIdentity)

	participants.Put("/:id/addresses", creatorMW, ctrl.SaveAddresses)
	participants.Delete("/:id/addresses/:addressId", creatorMW, ctrl.DeleteAddress)

	participants.Put("/:id/bank-accounts", creatorMW, ctrl.SaveBankAccount)
	participants.Delete("/:id/bank-accounts/:accountId", creatorMW, ctrl.DeleteBankAccount)

	participants.Put("/:id/family-members", creatorMW, ctrl.SaveFamilyMembers)
	participants.Delete("/:id/family-members/:memberId", creatorMW, ctrl.DeleteFamilyMember)

	participants.Put("/:id/employment", creatorMW, ctrl.SaveEmployment)
	participants.Put("/:id/pension", creatorMW, ctrl.SavePension)

	participants.Put("/:id/beneficiaries", creatorMW, ctrl.SaveBeneficiaries)
	participants.Delete("/:id/beneficiaries/:beneficiaryId", creatorMW, ctrl.DeleteBeneficiary)

	participants.Post("/:id/files", creatorMW, ctrl.UploadFile)
	participants.Get("/:id/status-history", anyRoleMW, ctrl.GetStatusHistory)

	participants.Post("/:id/submit", creatorMW, ctrl.Submit)
	participants.Post("/:id/approve", approverMW, ctrl.Approve)
	participants.Post("/:id/reject", approverMW, ctrl.Reject)
	participants.Delete("/:id", approverMW, ctrl.Delete)
}
