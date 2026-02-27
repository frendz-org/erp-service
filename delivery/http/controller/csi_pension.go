package controller

import (
	"strconv"

	"erp-service/delivery/http/middleware"
	"erp-service/saving/participant"

	"github.com/gofiber/fiber/v2"
)

type CsiPensionUsecase interface {
	participant.CsiAmountSummaryReader
	participant.CsiLedgerHistoryReader
	participant.CsiBalanceOverTimeReader
}

type CsiPensionController struct {
	usecase CsiPensionUsecase
}

func NewCsiPensionController(usecase CsiPensionUsecase) *CsiPensionController {
	return &CsiPensionController{
		usecase: usecase,
	}
}

func (ctrl *CsiPensionController) GetCsiAmountSummary(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	summaries, err := ctrl.usecase.GetCsiAmountSummary(c.UserContext(), &participant.CsiAmountSummaryRequest{
		UserID: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to get CSI amount summary",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    summaries,
	})
}

func (ctrl *CsiPensionController) GetCsiLedgerHistory(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	page := queryInt(c, "page", 1)
	perPage := queryInt(c, "per_page", 12)
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 12
	}
	if perPage > 120 {
		perPage = 120
	}

	var yearFrom, yearTo *int
	if v := c.Query("year_from"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1900 || n > 2100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid year_from parameter",
			})
		}
		yearFrom = &n
	}
	if v := c.Query("year_to"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1900 || n > 2100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid year_to parameter",
			})
		}
		yearTo = &n
	}
	if yearFrom != nil && yearTo != nil && *yearFrom > *yearTo {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "year_from must be less than or equal to year_to",
		})
	}

	history, err := ctrl.usecase.GetCsiLedgerHistory(c.UserContext(), &participant.CsiLedgerHistoryRequest{
		UserID:   userID,
		Page:     page,
		PerPage:  perPage,
		YearFrom: yearFrom,
		YearTo:   yearTo,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to get CSI ledger history",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    history,
	})
}

func (ctrl *CsiPensionController) GetBalanceOverTime(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	granularity := c.Query("granularity", "yearly")
	if granularity != "monthly" && granularity != "yearly" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "granularity must be 'monthly' or 'yearly'",
		})
	}

	var yearFrom, yearTo *int
	if v := c.Query("year_from"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1900 || n > 2100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid year_from parameter",
			})
		}
		yearFrom = &n
	}
	if v := c.Query("year_to"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil || n < 1900 || n > 2100 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"error":   "invalid year_to parameter",
			})
		}
		yearTo = &n
	}
	if yearFrom != nil && yearTo != nil && *yearFrom > *yearTo {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"error":   "year_from must be less than or equal to year_to",
		})
	}

	results, err := ctrl.usecase.GetBalanceOverTime(c.UserContext(), &participant.BalanceOverTimeRequest{
		UserID:      userID,
		Granularity: granularity,
		YearFrom:    yearFrom,
		YearTo:      yearTo,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to get balance over time",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data":    results,
	})
}

func queryInt(c *fiber.Ctx, key string, defaultVal int) int {
	v := c.Query(key)
	if v == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return defaultVal
	}
	return n
}
