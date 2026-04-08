export interface ApiResponse<TData> {
  readonly success: boolean;
  readonly data: TData;
  readonly meta: {
    readonly request_id: string;
  };
}

export interface ApiErrorResponse {
  readonly success: false;
  readonly error: {
    readonly code: string;
    readonly message: string;
  };
  readonly meta: {
    readonly request_id: string;
  };
}
