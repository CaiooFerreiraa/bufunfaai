package integration

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"

	analyticsservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/service"
	analyticsusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/application/usecase"
	analyticshandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/handler"
	analyticsroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/analytics/interfaces/http/routes"
	authservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/service"
	authusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/application/usecase"
	authentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/entity"
	authrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/domain/repository"
	authhash "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/infrastructure/hash"
	authtoken "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/infrastructure/token"
	authhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/handler"
	authroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/auth/interfaces/http/routes"
	consentshandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/consents/interfaces/http/handler"
	consentsroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/consents/interfaces/http/routes"
	deviceservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/service"
	deviceusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/application/usecase"
	deviceentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/entity"
	devicerepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/domain/repository"
	devicehandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/handler"
	deviceroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/devices/interfaces/http/routes"
	healthhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/health/interfaces/http/handler"
	healthroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/health/interfaces/http/routes"
	openfinanceservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	openfinanceusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/usecase"
	ofentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
	ofrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/repository"
	openfinanceprovider "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/infrastructure/provider"
	openfinancehandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/handler"
	openfinanceroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/routes"
	userservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/service"
	userusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/application/usecase"
	userentity "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/entity"
	userrepository "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/domain/repository"
	userhandler "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/handler"
	userroutes "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/users/interfaces/http/routes"
	platformauth "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/auth"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/config"
	platformcrypto "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/crypto"
	httplayer "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/http"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/platform/logger"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharedmiddleware "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
)

func newTestEngine() *gin.Engine {
	cfg := config.Config{
		AppName:             "finance-api",
		Environment:         "test",
		Port:                "8080",
		HTTPReadTimeout:     15 * time.Second,
		HTTPWriteTimeout:    15 * time.Second,
		HTTPShutdownTimeout: 5 * time.Second,
		CORSAllowedOrigins:  []string{"http://localhost:8081"},
		AccessTokenSecret:   "test-secret",
		AccessTokenIssuer:   "finance-api",
		AccessTokenAudience: "finance-mobile",
		AccessTokenTTL:      15 * time.Minute,
		RefreshTokenTTL:     30 * 24 * time.Hour,
		CronSecret:          "test-cron-secret",
		WorkerBatchSize:     10,
		LogLevel:            "error",
	}

	server := httplayer.NewServer(cfg, logger.New("error"))
	jwtService := platformauth.NewJWTService(cfg)
	requestValidator := platformvalidator.New()

	userRepo := newInMemoryUserRepository()
	tokenAuthenticator := &testTokenAuthenticator{jwtService: jwtService}
	deviceRepo := newInMemoryDeviceRepository()
	refreshRepo := newInMemoryRefreshTokenRepository()
	analyticsRepo := newInMemoryAnalyticsRepository()
	ofInstitutionRepo := newInMemoryInstitutionRepository()
	ofConsentRepo := newInMemoryConsentRepository()
	ofAuthorizationRepo := newInMemoryAuthorizationRepository()
	ofTokenRepo := newInMemoryTokenRepository()
	ofConnectionRepo := newInMemoryConnectionRepository()
	ofSyncJobRepo := newInMemorySyncJobRepository()
	cipherService, _ := platformcrypto.NewCipherService("0123456789ABCDEF0123456789ABCDEF")

	authService := authservice.NewAuthService(
		userRepo,
		refreshRepo,
		deviceRepo,
		authhash.NewPasswordHasher(platformcrypto.NewArgon2idPasswordHasher()),
		authhash.NewRefreshTokenManager(),
		authtoken.NewAccessTokenService(jwtService),
		cfg.RefreshTokenTTL,
	)
	userService := userservice.NewUserService(userRepo)
	deviceService := deviceservice.NewDeviceService(deviceRepo)
	analyticsService := analyticsservice.NewService(analyticsRepo)
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
	analyticsUseCases := analyticsusecase.New(analyticsService)
	_ = openFinanceUseCases.EnsureInstitutions(context.Background())

	healthroutes.Register(server.Engine(), healthhandler.NewHandler(cfg, nil))

	publicV1 := server.Engine().Group("/v1")
	protectedV1 := server.Engine().Group("/v1")
	internalGroup := server.Engine().Group("/internal")
	protectedV1.Use(sharedmiddleware.RequireAuthentication(tokenAuthenticator))
	internalGroup.Use(sharedmiddleware.RequireInternalSecret(cfg.CronSecret))

	authroutes.Register(
		publicV1,
		protectedV1,
		authhandler.NewHandler(
			authusecase.NewRegisterUseCase(authService),
			authusecase.NewLoginUseCase(authService),
			authusecase.NewRefreshUseCase(authService),
			authusecase.NewLogoutUseCase(authService),
			authusecase.NewLogoutAllUseCase(authService),
			requestValidator,
		),
	)
	userroutes.Register(
		protectedV1,
		userhandler.NewHandler(
			userusecase.NewGetCurrentUserUseCase(userService),
			userusecase.NewUpdateCurrentUserUseCase(userService),
			requestValidator,
		),
	)
	deviceroutes.Register(
		protectedV1,
		devicehandler.NewHandler(
			deviceusecase.NewListDevicesUseCase(deviceService),
			deviceusecase.NewDeleteDeviceUseCase(deviceService),
		),
	)
	consentsroutes.Register(protectedV1, consentshandler.NewHandler())
	analyticsroutes.Register(protectedV1, analyticshandler.NewHandler(analyticsUseCases, requestValidator))
	openfinanceroutes.Register(publicV1, protectedV1, openfinancehandler.NewHandler(openFinanceUseCases, requestValidator))
	openfinanceroutes.RegisterInternal(internalGroup, openfinancehandler.NewHandler(openFinanceUseCases, requestValidator))

	return server.Engine()
}

type testTokenAuthenticator struct {
	jwtService *platformauth.JWTService
}

func (authenticator *testTokenAuthenticator) Authenticate(_ context.Context, bearerToken string) (platformauth.AuthenticatedIdentity, error) {
	claims, err := authenticator.jwtService.Parse(bearerToken)
	if err != nil {
		return platformauth.AuthenticatedIdentity{}, err
	}

	return platformauth.AuthenticatedIdentity{
		LocalUserID: claims.Subject,
		SessionID:   claims.SessionID,
	}, nil
}

type inMemoryUserRepository struct {
	mutex   sync.RWMutex
	byID    map[string]userentity.User
	byEmail map[string]string
}

func newInMemoryUserRepository() *inMemoryUserRepository {
	return &inMemoryUserRepository{
		byID:    make(map[string]userentity.User),
		byEmail: make(map[string]string),
	}
}

func (repository *inMemoryUserRepository) Create(_ context.Context, user userentity.User) (userentity.User, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if _, exists := repository.byEmail[user.Email]; exists {
		return userentity.User{}, errors.New("email already exists")
	}

	repository.byID[user.ID] = user
	repository.byEmail[user.Email] = user.ID
	return user, nil
}

func (repository *inMemoryUserRepository) GetByEmail(_ context.Context, email string) (userentity.User, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	userID, ok := repository.byEmail[strings.ToLower(strings.TrimSpace(email))]
	if !ok {
		return userentity.User{}, userservice.ErrUserNotFound
	}

	return repository.byID[userID], nil
}

func (repository *inMemoryUserRepository) GetByID(_ context.Context, userID string) (userentity.User, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	user, ok := repository.byID[userID]
	if !ok {
		return userentity.User{}, userservice.ErrUserNotFound
	}

	return user, nil
}

func (repository *inMemoryUserRepository) UpdateProfile(_ context.Context, userID string, fullName string, phone string) (userentity.User, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	user, ok := repository.byID[userID]
	if !ok {
		return userentity.User{}, userservice.ErrUserNotFound
	}

	user.FullName = fullName
	user.Phone = phone
	user.UpdatedAt = time.Now().UTC()
	repository.byID[userID] = user
	return user, nil
}

func (repository *inMemoryUserRepository) UpdateLastLogin(_ context.Context, userID string) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	user, ok := repository.byID[userID]
	if !ok {
		return userservice.ErrUserNotFound
	}

	now := time.Now().UTC()
	user.LastLoginAt = &now
	user.UpdatedAt = now
	repository.byID[userID] = user
	return nil
}

var _ userrepository.UserRepository = (*inMemoryUserRepository)(nil)

type inMemoryDeviceRepository struct {
	mutex   sync.RWMutex
	devices map[string]deviceentity.Device
}

func newInMemoryDeviceRepository() *inMemoryDeviceRepository {
	return &inMemoryDeviceRepository{
		devices: make(map[string]deviceentity.Device),
	}
}

func (repository *inMemoryDeviceRepository) Upsert(_ context.Context, device deviceentity.Device) (deviceentity.Device, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for existingID, existingDevice := range repository.devices {
		if existingDevice.UserID == device.UserID && existingDevice.FingerprintHash != "" && existingDevice.FingerprintHash == device.FingerprintHash {
			existingDevice.DeviceName = device.DeviceName
			existingDevice.Platform = device.Platform
			existingDevice.AppVersion = device.AppVersion
			existingDevice.LastSeenAt = device.LastSeenAt
			repository.devices[existingID] = existingDevice
			return existingDevice, nil
		}
	}

	repository.devices[device.ID] = device
	return device, nil
}

func (repository *inMemoryDeviceRepository) ListByUserID(_ context.Context, userID string) ([]deviceentity.Device, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	devices := make([]deviceentity.Device, 0)
	for _, device := range repository.devices {
		if device.UserID == userID {
			devices = append(devices, device)
		}
	}

	return devices, nil
}

func (repository *inMemoryDeviceRepository) Delete(_ context.Context, userID string, deviceID string) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	device, ok := repository.devices[deviceID]
	if !ok || device.UserID != userID {
		return deviceservice.ErrDeviceNotFound
	}

	delete(repository.devices, deviceID)
	return nil
}

var _ devicerepository.DeviceRepository = (*inMemoryDeviceRepository)(nil)

type inMemoryRefreshTokenRepository struct {
	mutex  sync.RWMutex
	tokens map[string]authentity.RefreshToken
}

func newInMemoryRefreshTokenRepository() *inMemoryRefreshTokenRepository {
	return &inMemoryRefreshTokenRepository{
		tokens: make(map[string]authentity.RefreshToken),
	}
}

func (repository *inMemoryRefreshTokenRepository) Create(_ context.Context, token authentity.RefreshToken) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.tokens[token.TokenHash] = token
	return nil
}

func (repository *inMemoryRefreshTokenRepository) GetByTokenHash(_ context.Context, tokenHash string) (authentity.RefreshToken, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	token, ok := repository.tokens[tokenHash]
	if !ok {
		return authentity.RefreshToken{}, authservice.ErrRefreshTokenNotFound
	}

	return token, nil
}

func (repository *inMemoryRefreshTokenRepository) Revoke(_ context.Context, tokenID string, revokedAt time.Time) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for tokenHash, token := range repository.tokens {
		if token.ID == tokenID {
			token.RevokedAt = &revokedAt
			repository.tokens[tokenHash] = token
			return nil
		}
	}

	return authservice.ErrRefreshTokenNotFound
}

func (repository *inMemoryRefreshTokenRepository) RevokeByUserID(_ context.Context, userID string, revokedAt time.Time) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for tokenHash, token := range repository.tokens {
		if token.UserID == userID && token.RevokedAt == nil {
			token.RevokedAt = &revokedAt
			repository.tokens[tokenHash] = token
		}
	}

	return nil
}

func (repository *inMemoryRefreshTokenRepository) RevokeByUserIDAndDeviceID(_ context.Context, userID string, deviceID string, revokedAt time.Time) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for tokenHash, token := range repository.tokens {
		if token.UserID == userID && token.DeviceID != nil && *token.DeviceID == deviceID && token.RevokedAt == nil {
			token.RevokedAt = &revokedAt
			repository.tokens[tokenHash] = token
		}
	}

	return nil
}

var _ authrepository.RefreshTokenRepository = (*inMemoryRefreshTokenRepository)(nil)

type inMemoryInstitutionRepository struct {
	mutex        sync.RWMutex
	institutions map[string]ofentity.Institution
}

func newInMemoryInstitutionRepository() *inMemoryInstitutionRepository {
	return &inMemoryInstitutionRepository{institutions: make(map[string]ofentity.Institution)}
}

func (repository *inMemoryInstitutionRepository) SaveMany(_ context.Context, institutions []ofentity.Institution) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	for _, institution := range institutions {
		repository.institutions[institution.ID] = institution
	}

	return nil
}

func (repository *inMemoryInstitutionRepository) List(_ context.Context) ([]ofentity.Institution, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]ofentity.Institution, 0, len(repository.institutions))
	for _, institution := range repository.institutions {
		items = append(items, institution)
	}

	return items, nil
}

func (repository *inMemoryInstitutionRepository) GetByID(_ context.Context, institutionID string) (ofentity.Institution, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	institution, ok := repository.institutions[institutionID]
	if !ok {
		return ofentity.Institution{}, openfinanceservice.ErrInstitutionNotFound
	}

	return institution, nil
}

var _ ofrepository.InstitutionRepository = (*inMemoryInstitutionRepository)(nil)

type inMemoryConsentRepository struct {
	mutex     sync.RWMutex
	consents  map[string]ofentity.Consent
	stateToID map[string]string
}

func newInMemoryConsentRepository() *inMemoryConsentRepository {
	return &inMemoryConsentRepository{
		consents:  make(map[string]ofentity.Consent),
		stateToID: make(map[string]string),
	}
}

func (repository *inMemoryConsentRepository) Create(_ context.Context, consent ofentity.Consent) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.consents[consent.ID] = consent
	repository.stateToID[consent.State] = consent.ID
	return nil
}

func (repository *inMemoryConsentRepository) GetByID(_ context.Context, consentID string) (ofentity.Consent, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	consent, ok := repository.consents[consentID]
	if !ok {
		return ofentity.Consent{}, openfinanceservice.ErrConsentNotFound
	}

	return consent, nil
}

func (repository *inMemoryConsentRepository) GetByState(_ context.Context, state string) (ofentity.Consent, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	consentID, ok := repository.stateToID[state]
	if !ok {
		return ofentity.Consent{}, openfinanceservice.ErrConsentNotFound
	}

	return repository.consents[consentID], nil
}

func (repository *inMemoryConsentRepository) Update(_ context.Context, consent ofentity.Consent) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.consents[consent.ID] = consent
	repository.stateToID[consent.State] = consent.ID
	return nil
}

var _ ofrepository.ConsentRepository = (*inMemoryConsentRepository)(nil)

type inMemoryAuthorizationRepository struct {
	mutex          sync.RWMutex
	authorizations []ofentity.Authorization
}

func newInMemoryAuthorizationRepository() *inMemoryAuthorizationRepository {
	return &inMemoryAuthorizationRepository{authorizations: make([]ofentity.Authorization, 0)}
}

func (repository *inMemoryAuthorizationRepository) Create(_ context.Context, authorization ofentity.Authorization) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.authorizations = append(repository.authorizations, authorization)
	return nil
}

var _ ofrepository.AuthorizationRepository = (*inMemoryAuthorizationRepository)(nil)

type inMemoryTokenRepository struct {
	mutex  sync.RWMutex
	tokens map[string]ofentity.Token
}

func newInMemoryTokenRepository() *inMemoryTokenRepository {
	return &inMemoryTokenRepository{tokens: make(map[string]ofentity.Token)}
}

func (repository *inMemoryTokenRepository) UpsertByConsentID(_ context.Context, token ofentity.Token) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.tokens[token.ConsentID] = token
	return nil
}

func (repository *inMemoryTokenRepository) GetByConsentID(_ context.Context, consentID string) (ofentity.Token, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	token, ok := repository.tokens[consentID]
	if !ok {
		return ofentity.Token{}, openfinanceservice.ErrTokenNotFound
	}

	return token, nil
}

func (repository *inMemoryTokenRepository) RevokeByConsentID(_ context.Context, consentID string) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	token, ok := repository.tokens[consentID]
	if !ok {
		return nil
	}

	now := time.Now().UTC()
	token.RevokedAt = &now
	repository.tokens[consentID] = token
	return nil
}

var _ ofrepository.TokenRepository = (*inMemoryTokenRepository)(nil)

type inMemoryConnectionRepository struct {
	mutex           sync.RWMutex
	connections     map[string]ofentity.Connection
	consentToConnID map[string]string
}

func newInMemoryConnectionRepository() *inMemoryConnectionRepository {
	return &inMemoryConnectionRepository{
		connections:     make(map[string]ofentity.Connection),
		consentToConnID: make(map[string]string),
	}
}

func (repository *inMemoryConnectionRepository) CreateOrUpdate(_ context.Context, connection ofentity.Connection) (ofentity.Connection, error) {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	if existingID, ok := repository.consentToConnID[connection.ConsentID]; ok {
		existing := repository.connections[existingID]
		existing.Status = connection.Status
		existing.FirstSyncAt = coalesceTimePointer(existing.FirstSyncAt, connection.FirstSyncAt)
		existing.LastSyncAt = connection.LastSyncAt
		existing.LastSuccessfulSyncAt = connection.LastSuccessfulSyncAt
		existing.LastErrorCode = connection.LastErrorCode
		existing.LastErrorMessageRedacted = connection.LastErrorMessageRedacted
		existing.UpdatedAt = connection.UpdatedAt
		repository.connections[existingID] = existing
		return existing, nil
	}

	repository.connections[connection.ID] = connection
	repository.consentToConnID[connection.ConsentID] = connection.ID
	return connection, nil
}

func (repository *inMemoryConnectionRepository) ListByUserID(_ context.Context, userID string) ([]ofentity.Connection, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]ofentity.Connection, 0)
	for _, connection := range repository.connections {
		if connection.UserID == userID {
			items = append(items, connection)
		}
	}

	return items, nil
}

func (repository *inMemoryConnectionRepository) ListActive(_ context.Context, limit int) ([]ofentity.Connection, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	items := make([]ofentity.Connection, 0)
	for _, connection := range repository.connections {
		if connection.Status == ofentity.ConnectionStatusActive || connection.Status == ofentity.ConnectionStatusSyncError {
			items = append(items, connection)
			if len(items) == limit {
				break
			}
		}
	}

	return items, nil
}

func (repository *inMemoryConnectionRepository) GetByID(_ context.Context, connectionID string) (ofentity.Connection, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	connection, ok := repository.connections[connectionID]
	if !ok {
		return ofentity.Connection{}, openfinanceservice.ErrConnectionNotFound
	}

	return connection, nil
}

func (repository *inMemoryConnectionRepository) GetByConsentID(_ context.Context, consentID string) (ofentity.Connection, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	connectionID, ok := repository.consentToConnID[consentID]
	if !ok {
		return ofentity.Connection{}, openfinanceservice.ErrConnectionNotFound
	}

	return repository.connections[connectionID], nil
}

func (repository *inMemoryConnectionRepository) Update(_ context.Context, connection ofentity.Connection) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.connections[connection.ID] = connection
	repository.consentToConnID[connection.ConsentID] = connection.ID
	return nil
}

var _ ofrepository.ConnectionRepository = (*inMemoryConnectionRepository)(nil)

type inMemorySyncJobRepository struct {
	mutex sync.RWMutex
	jobs  map[string][]ofentity.SyncJob
}

func newInMemorySyncJobRepository() *inMemorySyncJobRepository {
	return &inMemorySyncJobRepository{jobs: make(map[string][]ofentity.SyncJob)}
}

func (repository *inMemorySyncJobRepository) ReplaceForConnection(_ context.Context, connectionID string, jobs []ofentity.SyncJob) error {
	repository.mutex.Lock()
	defer repository.mutex.Unlock()

	repository.jobs[connectionID] = jobs
	return nil
}

func (repository *inMemorySyncJobRepository) ListByConnectionID(_ context.Context, connectionID string) ([]ofentity.SyncJob, error) {
	repository.mutex.RLock()
	defer repository.mutex.RUnlock()

	return repository.jobs[connectionID], nil
}

var _ ofrepository.SyncJobRepository = (*inMemorySyncJobRepository)(nil)

func coalesceTimePointer(left *time.Time, right *time.Time) *time.Time {
	if left != nil {
		return left
	}

	return right
}
