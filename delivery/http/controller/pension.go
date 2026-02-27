package controller

import (
	"math"
	"strings"
	"time"

	"erp-service/delivery/http/middleware"
	"erp-service/entity"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PensionController struct {
	platformDB *gorm.DB
	// prodDB     *gorm.DB
}

func NewPensionController(platformDB *gorm.DB) *PensionController {
	return &PensionController{
		platformDB: platformDB,
		// prodDB:     prodDB,
	}
}

type pensionTransaction struct {
	ID              int       `gorm:"column:id;primaryKey"`
	TransactionDate time.Time `gorm:"column:transaction_date"`
	Code            string    `gorm:"column:code"`
	EmpNo           string    `gorm:"column:emp_no"`
	Amount          float64   `gorm:"column:amount"`
	Transaction     string    `gorm:"column:transaction"`
	Status          *bool     `gorm:"column:status"`
	EmpNoOld        *string   `gorm:"column:emp_no_old"`
	ROI             *float64  `gorm:"column:roi"`
}

func (pensionTransaction) TableName() string {
	return "pension_transactions"
}

func (ctrl *PensionController) GetAmountSummary(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"success": false,
			"error":   "unauthorized",
		})
	}

	// Get user full name and employee number from platform DB
	fullName, empNo, err := ctrl.resolveUser(c, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to resolve user",
		})
	}

	// Query pension summary from prod DB
	totalSaldo, growthPercentage, err := ctrl.calculatePensionSummary(empNo)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "failed to calculate pension summary",
		})
	}

	return c.JSON(fiber.Map{
		"success": true,
		"data": fiber.Map{
			"user":             fullName,
			"totalSaldo":       totalSaldo,
			"growthPercentage": growthPercentage,
		},
	})
}

func (ctrl *PensionController) resolveUser(c *fiber.Ctx, userID uuid.UUID) (string, string, error) {
	// Get full name from user_profiles
	var profile entity.UserProfile
	if err := ctrl.platformDB.WithContext(c.UserContext()).
		Where("user_id = ?", userID).
		First(&profile).Error; err != nil {
		return "", "", err
	}

	fullName := strings.TrimSpace(profile.FirstName + " " + profile.LastName)

	// Get employee_number from any participant record for this user
	var participant entity.Participant
	err := ctrl.platformDB.WithContext(c.UserContext()).
		Where("user_id = ? AND deleted_at IS NULL", userID).
		First(&participant).Error
	if err != nil {
		// No participant found — empNo stays empty (aggregate mode)
		return fullName, "", nil
	}

	empNo := ""
	if participant.EmployeeNumber != nil {
		empNo = *participant.EmployeeNumber
	}

	return fullName, empNo, nil
}

func (ctrl *PensionController) calculatePensionSummary(empNo string) (int64, float64, error) {
	// Step 1: Find latest transaction date
	var latestDate *time.Time
	q := ctrl.platformDB.Model(&pensionTransaction{})
	if empNo != "" {
		q = q.Where("emp_no = ?", empNo)
	}
	if err := q.Select("MAX(transaction_date)").Scan(&latestDate).Error; err != nil {
		return 0, 0, err
	}

	if latestDate == nil {
		return 0, 0, nil
	}

	// Step 2: Calculate total saldo (cumulative up to latest date)
	var totalSaldoRaw float64
	q2 := ctrl.platformDB.Model(&pensionTransaction{}).
		Where("transaction_date <= ?", *latestDate)
	if empNo != "" {
		q2 = q2.Where("emp_no = ?", empNo)
	}
	if err := q2.Select("COALESCE(SUM(amount), 0)").Scan(&totalSaldoRaw).Error; err != nil {
		return 0, 0, err
	}

	// Step 3: Calculate last month saldo
	oneMonthBefore := latestDate.AddDate(0, -1, 0)
	var lastMonthSaldo float64
	q3 := ctrl.platformDB.Model(&pensionTransaction{}).
		Where("transaction_date <= ?", oneMonthBefore)
	if empNo != "" {
		q3 = q3.Where("emp_no = ?", empNo)
	}
	if err := q3.Select("COALESCE(SUM(amount), 0)").Scan(&lastMonthSaldo).Error; err != nil {
		return 0, 0, err
	}

	// Step 4: Calculate growth percentage
	var growthPercentage float64
	if lastMonthSaldo == 0 {
		if totalSaldoRaw > 0 {
			growthPercentage = 100.0
		}
	} else {
		growthPercentage = ((totalSaldoRaw - lastMonthSaldo) / lastMonthSaldo) * 100
	}
	growthPercentage = math.Round(growthPercentage*100) / 100

	totalSaldo := int64(math.Round(totalSaldoRaw))

	return totalSaldo, growthPercentage, nil
}
