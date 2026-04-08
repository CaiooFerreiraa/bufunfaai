package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	ofdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/dto"
	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	ofusecase "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/usecase"
	ofpresenter "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/interfaces/http/presenter"
	platformvalidator "github.com/bufunfaai/bufunfaai/apps/api/internal/platform/validator"
	sharederrors "github.com/bufunfaai/bufunfaai/apps/api/internal/shared/errors"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/middleware"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/shared/response"
)

type Handler struct {
	useCases  *ofusecase.UseCases
	validator *platformvalidator.Validator
}

func NewHandler(useCases *ofusecase.UseCases, validator *platformvalidator.Validator) *Handler {
	return &Handler{useCases: useCases, validator: validator}
}

func (handler *Handler) ListInstitutions(context *gin.Context) {
	institutions, appError := handler.useCases.ListInstitutions(context.Request.Context())
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]ofdto.InstitutionOutput, 0, len(institutions))
	for _, institution := range institutions {
		payload = append(payload, ofpresenter.InstitutionOutput(institution))
	}

	response.OK(context, gin.H{"institutions": payload})
}

func (handler *Handler) GetInstitution(context *gin.Context) {
	institution, appError := handler.useCases.GetInstitution(context.Request.Context(), context.Param("id"))
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"institution": ofpresenter.InstitutionOutput(institution)})
}

func (handler *Handler) CreateConsent(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	var request ofdto.CreateConsentRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	consent, appError := handler.useCases.CreateConsent(context.Request.Context(), userID, request)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.Created(context, gin.H{"consent": ofpresenter.ConsentOutput(consent)})
}

func (handler *Handler) GetConsent(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	consent, appError := handler.useCases.GetConsent(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"consent": ofpresenter.ConsentOutput(consent)})
}

func (handler *Handler) AuthorizeConsent(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	authorizationURL, appError := handler.useCases.AuthorizeConsent(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, ofdto.AuthorizationURLOutput{
		ConsentID:        context.Param("id"),
		AuthorizationURL: authorizationURL,
	})
}

func (handler *Handler) CreateConnectToken(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	connectToken, appError := handler.useCases.CreateConnectToken(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, ofdto.ConnectTokenOutput{
		ConsentID:           context.Param("id"),
		ConnectToken:        connectToken.ConnectToken,
		SelectedConnectorID: connectToken.SelectedConnectorID,
	})
}

func (handler *Handler) CompleteConsent(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	var request ofdto.CompleteConsentRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	consent, connection, appError := handler.useCases.CompleteConsent(context.Request.Context(), context.Param("id"), userID, request.ItemID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, ofdto.CallbackResultOutput{
		Consent:    ofpresenter.ConsentOutput(consent),
		Connection: ofpresenter.ConnectionOutput(connection),
	})
}

func (handler *Handler) Callback(context *gin.Context) {
	state := context.Query("state")
	code := context.Query("code")

	if state == "" || code == "" {
		var request ofdto.CallbackRequest
		if err := context.ShouldBindJSON(&request); err != nil {
			response.Error(context, sharederrors.New("INVALID_PAYLOAD", "payload invalido", http.StatusBadRequest))
			return
		}
		if appError := handler.validator.Validate(request); appError != nil {
			response.Error(context, appError)
			return
		}

		state = request.State
		code = request.Code
	}

	consent, connection, appError := handler.useCases.HandleCallback(context.Request.Context(), state, code)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, ofdto.CallbackResultOutput{
		Consent:    ofpresenter.ConsentOutput(consent),
		Connection: ofpresenter.ConnectionOutput(connection),
	})
}

func (handler *Handler) RevokeConsent(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	appError := handler.useCases.RevokeConsent(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"revoked": true})
}

func (handler *Handler) ListConnections(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	connections, appError := handler.useCases.ListConnections(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]ofdto.ConnectionOutput, 0, len(connections))
	for _, connection := range connections {
		payload = append(payload, ofpresenter.ConnectionOutput(connection))
	}

	response.OK(context, gin.H{"connections": payload})
}

func (handler *Handler) ListAccounts(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	accounts, appError := handler.useCases.ListAccountSnapshots(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]ofdto.AccountSnapshotOutput, 0, len(accounts))
	for _, account := range accounts {
		payload = append(payload, ofpresenter.AccountSnapshotOutput(account))
	}

	response.OK(context, gin.H{"accounts": payload})
}

func (handler *Handler) Overview(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	overview, appError := handler.useCases.GetOverview(context.Request.Context(), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"overview": ofpresenter.OverviewOutput(overview)})
}

func (handler *Handler) ListTransactions(context *gin.Context) {
	userID := middleware.CurrentUserID(context)

	var request ofdto.TransactionsQuery
	if err := context.ShouldBindQuery(&request); err != nil {
		response.Error(context, sharederrors.New("INVALID_QUERY", "consulta invalida", http.StatusBadRequest))
		return
	}

	if appError := handler.validator.Validate(request); appError != nil {
		response.Error(context, appError)
		return
	}

	query := ofservice.ProviderTransactionQuery{
		PageSize: 500,
	}

	if request.Limit > 0 {
		query.PageSize = request.Limit
	}

	if request.From != "" {
		from, err := time.Parse("2006-01-02", request.From)
		if err != nil {
			response.Error(context, sharederrors.New("INVALID_QUERY", "periodo inicial invalido", http.StatusBadRequest))
			return
		}
		query.From = &from
	}

	if request.To != "" {
		to, err := time.Parse("2006-01-02", request.To)
		if err != nil {
			response.Error(context, sharederrors.New("INVALID_QUERY", "periodo final invalido", http.StatusBadRequest))
			return
		}
		query.To = &to
	}

	feed, appError := handler.useCases.ListTransactions(context.Request.Context(), userID, query)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"feed": ofpresenter.TransactionFeedOutput(feed)})
}

func (handler *Handler) GetConnection(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	connection, appError := handler.useCases.GetConnection(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"connection": ofpresenter.ConnectionOutput(connection)})
}

func (handler *Handler) DeleteConnection(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	connection, appError := handler.useCases.GetConnection(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	appError = handler.useCases.RevokeConsent(context.Request.Context(), connection.ConsentID, userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, gin.H{"deleted": true})
}

func (handler *Handler) SyncConnection(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	jobs, appError := handler.useCases.SyncConnection(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]ofdto.SyncJobOutput, 0, len(jobs))
	for _, job := range jobs {
		payload = append(payload, ofpresenter.SyncJobOutput(job))
	}

	response.OK(context, gin.H{"jobs": payload})
}

func (handler *Handler) SyncStatus(context *gin.Context) {
	userID := middleware.CurrentUserID(context)
	connection, jobs, appError := handler.useCases.SyncStatus(context.Request.Context(), context.Param("id"), userID)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	payload := make([]ofdto.SyncJobOutput, 0, len(jobs))
	for _, job := range jobs {
		payload = append(payload, ofpresenter.SyncJobOutput(job))
	}

	response.OK(context, ofdto.SyncStatusOutput{
		Connection: ofpresenter.ConnectionOutput(connection),
		Jobs:       payload,
	})
}

func (handler *Handler) ReconcileConnections(context *gin.Context) {
	limit := 25
	if rawLimit := context.Query("limit"); rawLimit != "" {
		parsedLimit, err := strconv.Atoi(rawLimit)
		if err != nil || parsedLimit <= 0 {
			response.Error(context, sharederrors.New("INVALID_LIMIT", "limite invalido", http.StatusBadRequest))
			return
		}

		limit = parsedLimit
	}

	result, appError := handler.useCases.ReconcileConnections(context.Request.Context(), limit)
	if appError != nil {
		response.Error(context, appError)
		return
	}

	response.OK(context, ofdto.ReconciliationResultOutput{
		Processed:   result.Processed,
		Successful:  result.Successful,
		Failed:      result.Failed,
		JobsCreated: result.JobsCreated,
	})
}
