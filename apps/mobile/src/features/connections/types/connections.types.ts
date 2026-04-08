export interface InstitutionItem {
  readonly id: string;
  readonly brand_name: string;
  readonly display_name: string;
  readonly authorisation_server_url: string;
  readonly resources_base_url: string;
  readonly status: string;
  readonly supports_data_sharing: boolean;
  readonly supports_payments: boolean;
}

export interface ConnectionItem {
  readonly id: string;
  readonly user_id: string;
  readonly institution_id: string;
  readonly consent_id: string;
  readonly status: string;
  readonly first_sync_at?: string;
  readonly last_sync_at?: string;
  readonly last_successful_sync_at?: string;
  readonly last_error_code?: string;
  readonly last_error_message_redacted?: string;
  readonly created_at: string;
  readonly updated_at: string;
}

export interface SyncJobItem {
  readonly id: string;
  readonly connection_id: string;
  readonly resource_type: string;
  readonly status: string;
  readonly attempt_count: number;
  readonly scheduled_at?: string;
  readonly started_at?: string;
  readonly finished_at?: string;
  readonly error_code?: string;
  readonly error_message_redacted?: string;
}

export interface SyncStatusData {
  readonly connection: ConnectionItem;
  readonly jobs: SyncJobItem[];
}
