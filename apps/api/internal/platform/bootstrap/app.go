package bootstrap

import (
	"context"
	"fmt"

	analyticsservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/service"
	analyticsusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/usecase"
	analyticsrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/infrastructure/repository"
	analyticshandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/handler"
	analyticsroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/routes"
	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	authusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/usecase"
	authhash "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/infrastructure/hash"
	authrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/infrastructure/repository"
	authtoken "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/infrastructure/token"
	authhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/handler"
	authroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/routes"
	consentshandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/consents/interfaces/http/handler"
	consentsroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/consents/interfaces/http/routes"
	deviceservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/service"
	deviceusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/usecase"
	devicerepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/infrastructure/repository"
	devicehandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/handler"
	deviceroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/routes"
	healthhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/health/interfaces/http/handler"
	healthroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/health/interfaces/http/routes"
	openfinanceservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	ofusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/usecase"
	openfinanceusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/usecase"
	openfinanceprovider "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/infrastructure/provider"
	openfinancerepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/infrastructure/repository"
	openfinancehandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/handler"
	openfinanceroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/routes"
	userservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	userusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/usecase"
	userrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/infrastructure/repository"
	userhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/handler"
	userroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/routes"
	platformauth "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/auth"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/cache"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
	platformcrypto "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/crypto"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/database"
	httplayer "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/http"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/logger"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/observability"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharedmiddleware "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
)

type App struct {
	Config              *config.Config
	Logger              *logger.Logger
	Server              *httplayer.Server
	OpenFinanceUseCases *ofusecase.UseCases
	close               func()
}

func Build(ctx context.Context) (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	appLogger := logger.New(cfg.LogLevel)
	observabilityShutdown, err := observability.Initialize(cfg, appLogger)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize observability: %w", err)
	}

	pool, err := database.Connect(ctx, cfg)
	if err != nil {
		observabilityShutdown(ctx)
		return nil, fmt.Errorf("failed to connect postgres: %w", err)
	}

	redisClient, err := cache.Connect(ctx, cfg)
	if err != nil {
		observabilityShutdown(ctx)
		pool.Close()
		return nil, fmt.Errorf("failed to connect redis: %w", err)
	}

	cipherService, err := platformcrypto.NewCipherService(cfg.EncryptionKey)
	if err != nil {
		observabilityShutdown(ctx)
		pool.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("failed to initialize crypto service: %w", err)
	}

	jwtService := platformauth.NewJWTService(cfg)
	requestValidator := platformvalidator.New()
	server := httplayer.NewServer(cfg, appLogger)
	observability.AttachGin(server.Engine(), cfg)

	userRepo := userrepository.NewPostgresUserRepository(pool)
	deviceRepo := devicerepository.NewPostgresDeviceRepository(pool)
	refreshTokenRepo := authrepository.NewPostgresRefreshTokenRepository(pool)
	analyticsRepo := analyticsrepository.NewPostgresAnalyticsRepository(pool)
	ofInstitutionRepo := openfinancerepository.NewPostgresInstitutionRepository(pool)
	ofConsentRepo := openfinancerepository.NewPostgresConsentRepository(pool)
	ofAuthorizationRepo := openfinancerepository.NewPostgresAuthorizationRepository(pool)
	ofTokenRepo := openfinancerepository.NewPostgresTokenRepository(pool)
	ofConnectionRepo := openfinancerepository.NewPostgresConnectionRepository(pool)
	ofSyncJobRepo := openfinancerepository.NewPostgresSyncJobRepository(pool)

	userService := userservice.NewUserService(userRepo)
	deviceService := deviceservice.NewDeviceService(deviceRepo)
	passwordHasher := authhash.NewPasswordHasher(platformcrypto.NewArgon2idPasswordHasher())
	refreshTokenManager := authhash.NewRefreshTokenManager()
	accessTokenService := authtoken.NewAccessTokenService(jwtService)
	authService := authservice.NewAuthService(
		userRepo,
		refreshTokenRepo,
		deviceRepo,
		passwordHasher,
		refreshTokenManager,
		accessTokenService,
		cfg.RefreshTokenTTL,
	)
	analyticsService := analyticsservice.NewService(analyticsRepo)
	analyticsUseCases := analyticsusecase.New(analyticsService)
	openFinanceService := openfinanceservice.NewService(
		ofInstitutionRepo,
		ofConsentRepo,
		ofAuthorizationRepo,
		ofTokenRepo,
		ofConnectionRepo,
		ofSyncJobRepo,
		openfinanceprovider.NewMockProvider(),
		cipherService,
	)
	openFinanceUseCases := openfinanceusecase.New(openFinanceService)

	healthHandler := healthhandler.NewHandler(cfg, []healthhandler.ReadinessProbe{
		database.NewReadinessProbe(pool),
		cache.NewReadinessProbe(redisClient),
	})
	authHTTPHandler := authhandler.NewHandler(
		authusecase.NewRegisterUseCase(authService),
		authusecase.NewLoginUseCase(authService),
		authusecase.NewRefreshUseCase(authService),
		authusecase.NewLogoutUseCase(authService),
		authusecase.NewLogoutAllUseCase(authService),
		requestValidator,
	)
	userHTTPHandler := userhandler.NewHandler(
		userusecase.NewGetCurrentUserUseCase(userService),
		userusecase.NewUpdateCurrentUserUseCase(userService),
		requestValidator,
	)
	deviceHTTPHandler := devicehandler.NewHandler(
		deviceusecase.NewListDevicesUseCase(deviceService),
		deviceusecase.NewDeleteDeviceUseCase(deviceService),
	)
	consentHTTPHandler := consentshandler.NewHandler()
	analyticsHTTPHandler := analyticshandler.NewHandler(analyticsUseCases, requestValidator)
	openFinanceHTTPHandler := openfinancehandler.NewHandler(openFinanceUseCases, requestValidator)

	healthroutes.Register(server.Engine(), healthHandler)

	publicV1 := server.Engine().Group("/v1")
	protectedV1 := server.Engine().Group("/v1")
	internalGroup := server.Engine().Group("/internal")
	protectedV1.Use(sharedmiddleware.RequireAuthentication(jwtService))
	internalGroup.Use(sharedmiddleware.RequireInternalSecret(cfg.CronSecret))

	authroutes.Register(publicV1, protectedV1, authHTTPHandler)
	userroutes.Register(protectedV1, userHTTPHandler)
	deviceroutes.Register(protectedV1, deviceHTTPHandler)
	consentsroutes.Register(protectedV1, consentHTTPHandler)
	analyticsroutes.Register(protectedV1, analyticsHTTPHandler)
	openfinanceroutes.Register(publicV1, protectedV1, openFinanceHTTPHandler)
	openfinanceroutes.RegisterInternal(internalGroup, openFinanceHTTPHandler)

	if appError := openFinanceUseCases.EnsureInstitutions(ctx); appError != nil {
		observabilityShutdown(ctx)
		pool.Close()
		_ = redisClient.Close()
		return nil, fmt.Errorf("failed to seed open finance institutions: %w", appError)
	}

	return &App{
		Config:              &cfg,
		Logger:              appLogger,
		Server:              server,
		OpenFinanceUseCases: openFinanceUseCases,
		close: func() {
			observabilityShutdown(context.Background())
			pool.Close()
			_ = redisClient.Close()
		},
	}, nil
}

func (app *App) Close() {
	if app == nil || app.close == nil {
		return
	}

	app.close()
}
