package presenter

import (
	"encoding/json"
	"math"
	"time"

	ofdto "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/dto"
	ofservice "github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/application/service"
	"github.com/bufunfaai/bufunfaai/apps/api/internal/modules/openfinance/domain/entity"
)

func InstitutionOutput(institution entity.Institution) ofdto.InstitutionOutput {
	return ofdto.InstitutionOutput{
		ID:                     institution.ID,
		DirectoryOrgID:         institution.DirectoryOrgID,
		BrandName:              institution.BrandName,
		DisplayName:            institution.DisplayName,
		AuthorisationServerURL: institution.AuthorisationServerURL,
		ResourcesBaseURL:       institution.ResourcesBaseURL,
		LogoURL:                institution.LogoURL,
		Status:                 institution.Status,
		SupportsDataSharing:    institution.SupportsDataSharing,
		SupportsPayments:       institution.SupportsPayments,
	}
}

func ConsentOutput(consent entity.Consent) ofdto.ConsentOutput {
	return ofdto.ConsentOutput{
		ID:                consent.ID,
		UserID:            consent.UserID,
		InstitutionID:     consent.InstitutionID,
		ExternalConsentID: consent.ExternalConsentID,
		Status:            consent.Status,
		Purpose:           consent.Purpose,
		Permissions:       parsePermissions(consent.PermissionsJSON),
		ExpirationAt:      formatTimePointer(consent.ExpirationAt),
		RevokedAt:         formatTimePointer(consent.RevokedAt),
		AuthorisedAt:      formatTimePointer(consent.AuthorisedAt),
		RedirectURI:       consent.RedirectURI,
		CreatedAt:         consent.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         consent.UpdatedAt.Format(time.RFC3339),
	}
}

func ConnectionOutput(connection entity.Connection) ofdto.ConnectionOutput {
	return ofdto.ConnectionOutput{
		ID:                   connection.ID,
		UserID:               connection.UserID,
		InstitutionID:        connection.InstitutionID,
		ConsentID:            connection.ConsentID,
		Status:               connection.Status,
		FirstSyncAt:          formatTimePointer(connection.FirstSyncAt),
		LastSyncAt:           formatTimePointer(connection.LastSyncAt),
		LastSuccessfulSyncAt: formatTimePointer(connection.LastSuccessfulSyncAt),
		LastErrorCode:        connection.LastErrorCode,
		LastErrorMessage:     connection.LastErrorMessageRedacted,
		CreatedAt:            connection.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            connection.UpdatedAt.Format(time.RFC3339),
	}
}

func SyncJobOutput(job entity.SyncJob) ofdto.SyncJobOutput {
	return ofdto.SyncJobOutput{
		ID:                   job.ID,
		ConnectionID:         job.ConnectionID,
		ResourceType:         job.ResourceType,
		Status:               job.Status,
		AttemptCount:         job.AttemptCount,
		ScheduledAt:          formatTimePointer(job.ScheduledAt),
		StartedAt:            formatTimePointer(job.StartedAt),
		FinishedAt:           formatTimePointer(job.FinishedAt),
		ErrorCode:            job.ErrorCode,
		ErrorMessageRedacted: job.ErrorMessageRedacted,
	}
}

func AccountSnapshotOutput(account ofservice.AccountSnapshot) ofdto.AccountSnapshotOutput {
	return ofdto.AccountSnapshotOutput{
		ID:                   account.ID,
		ConnectionID:         account.ConnectionID,
		InstitutionID:        account.InstitutionID,
		InstitutionName:      account.InstitutionName,
		ItemID:               account.ItemID,
		Type:                 account.Type,
		Subtype:              account.Subtype,
		Name:                 account.Name,
		MarketingName:        account.MarketingName,
		Number:               account.Number,
		CurrencyCode:         account.CurrencyCode,
		Balance:              account.Balance,
		BankTransferNumber:   account.BankTransferNumber,
		CreditBrand:          account.CreditBrand,
		AvailableCreditLimit: account.AvailableCreditLimit,
	}
}

func TransactionSnapshotOutput(transaction ofservice.TransactionSnapshot) ofdto.TransactionSnapshotOutput {
	return ofdto.TransactionSnapshotOutput{
		ID:              transaction.ID,
		AccountID:       transaction.AccountID,
		ConnectionID:    transaction.ConnectionID,
		InstitutionName: transaction.InstitutionName,
		AccountName:     transaction.AccountName,
		Description:     transaction.Description,
		Category:        transaction.Category,
		Type:            transaction.Type,
		Status:          transaction.Status,
		CurrencyCode:    transaction.CurrencyCode,
		Amount:          transaction.Amount,
		SignedAmount:    normalizedTransactionAmount(transaction),
		Date:            transaction.Date.Format(time.RFC3339),
		MerchantName:    transaction.MerchantName,
	}
}

func CategoryBreakdownOutput(category ofservice.CategoryBreakdown) ofdto.CategoryBreakdownOutput {
	return ofdto.CategoryBreakdownOutput{
		Category: category.Category,
		Amount:   category.Amount,
		Percent:  category.Percent,
	}
}

func OverviewOutput(overview ofservice.Overview) ofdto.OverviewOutput {
	accounts := make([]ofdto.AccountSnapshotOutput, 0, len(overview.Accounts))
	for _, account := range overview.Accounts {
		accounts = append(accounts, AccountSnapshotOutput(account))
	}

	transactions := make([]ofdto.TransactionSnapshotOutput, 0, len(overview.RecentTransactions))
	for _, transaction := range overview.RecentTransactions {
		transactions = append(transactions, TransactionSnapshotOutput(transaction))
	}

	breakdown := make([]ofdto.CategoryBreakdownOutput, 0, len(overview.ExpenseBreakdown))
	for _, category := range overview.ExpenseBreakdown {
		breakdown = append(breakdown, CategoryBreakdownOutput(category))
	}

	return ofdto.OverviewOutput{
		Accounts:            accounts,
		RecentTransactions:  transactions,
		ExpenseBreakdown:    breakdown,
		TotalAvailable:      overview.TotalAvailable,
		CreditCardBalance:   overview.CreditCardBalance,
		MonthIncome:         overview.MonthIncome,
		MonthExpenses:       overview.MonthExpenses,
		ConnectedAccounts:   overview.ConnectedAccounts,
		ConnectionsWithData: overview.ConnectionsWithData,
	}
}

func TransactionFeedOutput(feed ofservice.TransactionFeed) ofdto.TransactionFeedOutput {
	transactions := make([]ofdto.TransactionSnapshotOutput, 0, len(feed.Transactions))
	for _, transaction := range feed.Transactions {
		transactions = append(transactions, TransactionSnapshotOutput(transaction))
	}

	breakdown := make([]ofdto.CategoryBreakdownOutput, 0, len(feed.ExpenseBreakdown))
	for _, category := range feed.ExpenseBreakdown {
		breakdown = append(breakdown, CategoryBreakdownOutput(category))
	}

	return ofdto.TransactionFeedOutput{
		Transactions:              transactions,
		ExpenseBreakdown:          breakdown,
		MonthExpenseTotal:         feed.MonthExpenseTotal,
		PreviousMonthExpenseTotal: feed.PreviousMonthExpenseTotal,
		MonthIncomeTotal:          feed.MonthIncomeTotal,
	}
}

func parsePermissions(rawPermissions string) []string {
	if rawPermissions == "" {
		return []string{}
	}

	permissions := make([]string, 0)
	if err := json.Unmarshal([]byte(rawPermissions), &permissions); err != nil {
		return []string{}
	}

	return permissions
}

func formatTimePointer(value *time.Time) string {
	if value == nil {
		return ""
	}

	return value.Format(time.RFC3339)
}

func normalizedTransactionAmount(transaction ofservice.TransactionSnapshot) float64 {
	amount := math.Abs(transaction.Amount)
	if transaction.Type == "DEBIT" {
		return -amount
	}
	return amount
}
