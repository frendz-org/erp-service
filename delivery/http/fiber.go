package rest

import (
	"context"
	"encoding/json"
	"erp-service/config"
	"erp-service/delivery/http/controller"
	"erp-service/delivery/http/dto/response"
	"erp-service/delivery/http/middleware"
	"erp-service/delivery/http/router"
	"erp-service/delivery/worker"
	"erp-service/files"
	"erp-service/iam/auth"
	"erp-service/iam/product"
	"erp-service/iam/role"
	"erp-service/iam/user"
	"erp-service/impl/mailer"
	implminio "erp-service/impl/minio"
	"erp-service/impl/postgres"
	implredis "erp-service/impl/redis"
	"erp-service/infrastructure"
	"erp-service/masterdata"
	apperrors "erp-service/pkg/errors"
	"erp-service/pkg/logger"
	"erp-service/saving/member"
	"erp-service/saving/participant"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"go.uber.org/zap"
)

type Server struct {
	app          *fiber.App
	config       *config.Config
	logger       *zap.Logger
	fileWorker   *worker.Worker
	workerCancel context.CancelFunc
}

func NewServer(cfg *config.Config) *Server {
	zapLogger, _ := logger.NewZapLoggerWithConfig(cfg.Log, cfg.App.Environment)
	auditLogger := logger.NewAuditLogger(zapLogger, logger.AuditConfig{
		Enabled: cfg.Log.AuditEnabled,
	})

	app := fiber.New(fiber.Config{
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		AppName:      cfg.App.Name,
		BodyLimit:    6 * 1024 * 1024,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: createErrorHandler(cfg, zapLogger),
	})

	postgresDB, err := infrastructure.NewPostgres(cfg.Infra.Postgres, zapLogger)
	if err != nil {
		log.Fatal("failed to connect to postgres:", err)
	}

	postgresProdDB, err := infrastructure.NewPostgresProd(zapLogger)
	if err != nil {
		log.Fatal("failed to connect to postgres prod:", err)
	}

	redisClient, err := infrastructure.NewRedis(cfg.Infra.Redis)
	if err != nil {
		log.Fatal("failed to connect to redis:", err)
	}
	inMemoryStore := implredis.NewRedis(redisClient)

	txManager := postgres.NewTransactionManager(postgresDB)

	authUserRepo := postgres.NewUserRepository(postgresDB)
	userProfileRepo := postgres.NewUserProfileRepository(postgresDB)
	userAuthMethodRepo := postgres.NewUserAuthMethodRepository(postgresDB)
	userSecurityStateRepo := postgres.NewUserSecurityStateRepository(postgresDB)
	tenantRepo := postgres.NewTenantRepository(postgresDB)
	roleRepo := postgres.NewRoleRepository(postgresDB)
	refreshTokenRepo := postgres.NewRefreshTokenRepository(postgresDB)
	userRoleRepo := postgres.NewUserRoleRepository(postgresDB)
	productRepo := postgres.NewProductRepository(postgresDB)
	permissionRepo := postgres.NewPermissionRepository(postgresDB)
	rolePermissionRepo := postgres.NewRolePermissionRepository(postgresDB)
	userSessionRepo := postgres.NewUserSessionRepository(postgresDB)
	userTenantRegRepo := postgres.NewUserTenantRegistrationRepository(postgresDB)
	productsByTenantRepo := postgres.NewProductsByTenantRepository(postgresDB)

	masterdataCategoryRepo := postgres.NewMasterdataCategoryRepository(postgresDB)
	masterdataItemRepo := postgres.NewMasterdataItemRepository(postgresDB)

	productRegConfigRepo := postgres.NewProductRegistrationConfigRepository(postgresDB)

	memberRepo := postgres.NewMemberRepository(postgresDB)
	csiEmployeeRepo := postgres.NewCsiEmployeeRepository(postgresDB)
	participantRepo := postgres.NewParticipantRepository(postgresDB)
	participantIdentityRepo := postgres.NewParticipantIdentityRepository(postgresDB)
	participantAddressRepo := postgres.NewParticipantAddressRepository(postgresDB)
	participantBankAccountRepo := postgres.NewParticipantBankAccountRepository(postgresDB)
	participantFamilyMemberRepo := postgres.NewParticipantFamilyMemberRepository(postgresDB)
	participantEmploymentRepo := postgres.NewParticipantEmploymentRepository(postgresDB)
	participantPensionRepo := postgres.NewParticipantPensionRepository(postgresDB)
	participantBeneficiaryRepo := postgres.NewParticipantBeneficiaryRepository(postgresDB)
	participantStatusHistoryRepo := postgres.NewParticipantStatusHistoryRepository(postgresDB)
	fileRepo := postgres.NewFileRepository(postgresDB)

	minioClient, err := infrastructure.NewMinIOClient(cfg)
	if err != nil {
		log.Fatal("failed to connect to minio:", err)
	}
	fileStorage := implminio.NewFileStorage(minioClient)

	emailService := mailer.NewEmailService(&cfg.Email)

	masterdataUsecase := masterdata.NewUsecase(
		cfg,
		masterdataCategoryRepo,
		masterdataItemRepo,
		inMemoryStore,
	)
	authUsecase := auth.NewUsecase(
		txManager,
		cfg,
		authUserRepo,
		userProfileRepo,
		userAuthMethodRepo,
		userSecurityStateRepo,
		tenantRepo,
		roleRepo,
		refreshTokenRepo,
		userRoleRepo,
		productRepo,
		permissionRepo,
		emailService,
		inMemoryStore,
		userSessionRepo,
		userTenantRegRepo,
		productsByTenantRepo,
		auditLogger,
		masterdataUsecase,
	)
	roleUsecase := role.NewUsecase(
		txManager,
		cfg,
		tenantRepo,
		roleRepo,
		rolePermissionRepo,
	)
	userUsecase := user.NewUsecase(
		txManager,
		cfg,
		authUserRepo,
		userProfileRepo,
		userAuthMethodRepo,
		userSecurityStateRepo,
		tenantRepo,
		roleRepo,
		userRoleRepo,
		masterdataUsecase,
	)

	memberUsecase := member.NewUsecase(
		cfg,
		txManager,
		userTenantRegRepo,
		userRoleRepo,
		productRepo,
		roleRepo,
		productRegConfigRepo,
		userProfileRepo,
		authUserRepo,
		memberRepo,
		csiEmployeeRepo,
		tenantRepo,
		masterdataUsecase,
	)
	participantUsecase := participant.NewUsecase(
		cfg,
		zapLogger,
		txManager,
		participantRepo,
		participantIdentityRepo,
		participantAddressRepo,
		participantBankAccountRepo,
		participantFamilyMemberRepo,
		participantEmploymentRepo,
		participantPensionRepo,
		participantBeneficiaryRepo,
		participantStatusHistoryRepo,
		fileStorage,
		fileRepo,
		tenantRepo,
		productRepo,
		productRegConfigRepo,
		userTenantRegRepo,
		userProfileRepo,
		masterdataUsecase,
		csiEmployeeRepo,
	)

	healthController := controller.NewHealthController(cfg)
	authController := controller.NewRegistrationController(cfg, authUsecase)
	roleController := controller.NewRoleController(cfg, roleUsecase)
	userController := controller.NewUserController(cfg, userUsecase)
	masterdataController := controller.NewMasterdataController(cfg, masterdataUsecase)
	memberController := controller.NewMemberController(memberUsecase)
	participantController := controller.NewParticipantController(participantUsecase)
	devController := controller.NewDevController(postgresDB, inMemoryStore)
	pensionController := controller.NewPensionController(postgresDB, postgresProdDB)

	fileCleanupUC := files.NewUsecase(fileRepo, fileStorage, txManager, zapLogger, files.DefaultConfig())
	fileWorker := worker.NewWorker(fileCleanupUC, zapLogger)

	server := &Server{
		app:        app,
		config:     cfg,
		logger:     zapLogger,
		fileWorker: fileWorker,
	}

	mw := middleware.New(cfg, zapLogger)
	mw.Setup(app)

	api := app.Group("/api")
	if cfg.IsProduction() {
		api.Use(limiter.New(limiter.Config{
			Max:               10,
			Expiration:        1 * time.Minute,
			LimiterMiddleware: limiter.SlidingWindow{},
			KeyGenerator: func(c *fiber.Ctx) string {
				return c.IP()
			},
			LimitReached: func(c *fiber.Ctx) error {
				return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
					"success": false,
					"error":   "too many requests, please try again later",
				})
			},
		}))
	}
	v1 := api.Group("/v1")

	router.SetupHealthRoutes(v1, healthController)
	router.SetupDocsRoutes(v1)
	router.SetupMasterdataRoutes(v1, cfg, masterdataController, inMemoryStore)

	iam := v1.Group("/iam")
	router.SetupAuthRoutes(iam, cfg, authController, inMemoryStore)
	router.SetupRoleRoutes(iam, cfg, roleController, inMemoryStore)
	router.SetupUserRoutes(iam, cfg, userController, inMemoryStore)

	jwtMiddleware := middleware.JWTAuth(cfg, inMemoryStore)

	productUsecase := product.NewUsecase(productRepo, inMemoryStore)
	frendzSavingMW := middleware.ExtractFrendzSavingProduct(productUsecase)

	saving := v1.Group("/saving")
	router.SetupParticipantRoutes(saving, participantController, jwtMiddleware, frendzSavingMW, cfg)
	router.SetupMemberRoutes(saving, memberController, jwtMiddleware, frendzSavingMW)

	router.SetupPensionRoutes(v1, pensionController, jwtMiddleware)

	router.SetupDevRoutes(v1, cfg, devController)

	return server
}

func (s *Server) App() *fiber.App {
	return s.app
}

func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting server on %s\n", addr)
	return s.app.Listen(addr)
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.app.ShutdownWithContext(ctx)
}

func (s *Server) StartWorker(ctx context.Context) {
	workerCtx, cancel := context.WithCancel(ctx)
	s.workerCancel = cancel
	s.fileWorker.Start(workerCtx)
}

func (s *Server) StopWorker() {
	if s.workerCancel != nil {
		s.workerCancel()
	}
	s.fileWorker.Stop()
}

func createErrorHandler(cfg *config.Config, zapLogger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		requestID := middleware.GetRequestID(c)
		includeDebug := cfg.IsDevelopment()

		var appErr *apperrors.AppError
		if apperrors.As(err, &appErr) {

			logFields := make([]zap.Field, 0, 10)
			logFields = append(logFields,
				zap.String("request_id", requestID),
				zap.String("code", appErr.Code),
				zap.String("message", appErr.Message),
				zap.String("kind", appErr.Kind.String()),
				zap.String("file", appErr.File),
				zap.Int("line", appErr.Line),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
			if appErr.Op != "" {
				logFields = append(logFields, zap.String("op", appErr.Op))
			}
			if appErr.Err != nil {
				logFields = append(logFields, zap.Error(appErr.Err))
			}

			if appErr.HTTPStatus >= 500 {
				zapLogger.Error(appErr.Message, logFields...)
			} else if appErr.HTTPStatus >= 400 {
				zapLogger.Warn(appErr.Message, logFields...)
			}

			resp := response.APIResponse{
				Success:   false,
				Error:     appErr.Code,
				Message:   appErr.Message,
				RequestID: requestID,
			}

			if appErr.Code == apperrors.CodeValidation && appErr.Details != nil {
				if fieldErrors, ok := appErr.Details["fields"].([]apperrors.FieldError); ok {
					resp.Errors = make([]response.FieldError, len(fieldErrors))
					for i, fe := range fieldErrors {
						resp.Errors[i] = response.FieldError{
							Field:   fe.Field,
							Message: fe.Message,
						}
					}
				}
			}

			if includeDebug && appErr.Op != "" {
				resp.Debug = &response.DebugInfo{
					Cause: appErr.Op,
				}
			}

			return c.Status(appErr.HTTPStatus).JSON(resp)
		}

		var fiberErr *fiber.Error
		if errors.As(err, &fiberErr) {
			zapLogger.Warn("Fiber error",
				zap.String("request_id", requestID),
				zap.Int("status", fiberErr.Code),
				zap.String("message", fiberErr.Message),
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
			)
			return c.Status(fiberErr.Code).JSON(response.APIResponse{
				Success:   false,
				Error:     "FIBER_ERROR",
				Message:   fiberErr.Message,
				RequestID: requestID,
			})
		}

		zapLogger.Error("Unexpected error",
			zap.String("request_id", requestID),
			zap.Error(err),
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
		)

		resp := response.APIResponse{
			Success:   false,
			Error:     "INTERNAL_SERVER_ERROR",
			Message:   "an unexpected error occurred",
			RequestID: requestID,
		}

		if includeDebug {
			resp.Debug = &response.DebugInfo{
				Cause: err.Error(),
			}
		}

		return c.Status(fiber.StatusInternalServerError).JSON(resp)
	}
}
