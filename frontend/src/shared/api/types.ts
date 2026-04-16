export type ApiErrorCode =
  | "INVALID_REQUEST"
  | "UNAUTHORIZED"
  | "FORBIDDEN"
  | "NOT_FOUND"
  | "CONFLICT"
  | "INTERNAL_ERROR"
  | "UNKNOWN";

export class ApiError extends Error {
  readonly status: number;
  readonly code: ApiErrorCode;
  readonly payload: unknown;

  constructor(status: number, code: ApiErrorCode, message: string, payload: unknown) {
    super(message);
    this.name = "ApiError";
    this.status = status;
    this.code = code;
    this.payload = payload;
  }
}

export type PaginationQuery = {
  limit?: number;
  offset?: number;
};

export type ListResponse<T> = {
  items: T[];
  total_count: number;
  limit: number;
  offset: number;
};

export type TokenPair = {
  access_token: string;
  refresh_token: string;
};
