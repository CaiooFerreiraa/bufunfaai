export interface AccountSnapshot {
  readonly id: string;
  readonly connection_id: string;
  readonly institution_id: string;
  readonly institution_name: string;
  readonly item_id: string;
  readonly type: string;
  readonly subtype?: string;
  readonly name: string;
  readonly marketing_name?: string;
  readonly number?: string;
  readonly currency_code: string;
  readonly balance: number;
  readonly bank_transfer_number?: string;
  readonly credit_brand?: string;
  readonly available_credit_limit?: number;
}

export interface FinancialTransaction {
  readonly id: string;
  readonly account_id: string;
  readonly connection_id: string;
  readonly institution_name: string;
  readonly account_name: string;
  readonly description: string;
  readonly category: string;
  readonly type: string;
  readonly status: string;
  readonly currency_code: string;
  readonly amount: number;
  readonly signed_amount: number;
  readonly date: string;
  readonly merchant_name?: string;
}

export interface CategoryBreakdown {
  readonly category: string;
  readonly amount: number;
  readonly percent: number;
}

export interface FinancialOverview {
  readonly accounts: AccountSnapshot[];
  readonly recent_transactions: FinancialTransaction[];
  readonly expense_breakdown: CategoryBreakdown[];
  readonly total_available: number;
  readonly credit_card_balance: number;
  readonly month_income: number;
  readonly month_expenses: number;
  readonly connected_accounts: number;
  readonly connections_with_data: number;
}

export interface FinancialTransactionFeed {
  readonly transactions: FinancialTransaction[];
  readonly expense_breakdown: CategoryBreakdown[];
  readonly month_expense_total: number;
  readonly previous_month_expense_total: number;
  readonly month_income_total: number;
}
