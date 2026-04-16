import { ApiError, ApiErrorCode } from "@/shared/api/types";

type Primitive = string | number | boolean | null | undefined;
type QueryParams = Record<string, Primitive>;

type HttpClientConfig = {
  baseUrl: string;
  getAccessToken: () => string | null;
  refreshAccessToken?: () => Promise<string | null>;
  onAuthFailure?: () => void;
};

type RequestOptions<TBody> = {
  method?: "GET" | "POST" | "PATCH" | "PUT" | "DELETE";
  body?: TBody;
  headers?: Record<string, string>;
  query?: QueryParams;
  auth?: boolean;
  retryOnUnauthorized?: boolean;
  signal?: AbortSignal;
};

function buildUrl(baseUrl: string, path: string, query?: QueryParams): string {
  const url = new URL(path, baseUrl);

  if (!query) {
    return url.toString();
  }

  Object.entries(query).forEach(([key, value]) => {
    if (value === undefined || value === null || value === "") {
      return;
    }
    url.searchParams.set(key, String(value));
  });

  return url.toString();
}

function toApiError(status: number, payload: unknown): ApiError {
  if (payload && typeof payload === "object") {
    if ("error" in payload) {
      const rawError = (payload as { error: unknown }).error;
      if (rawError && typeof rawError === "object") {
        const code = (rawError as { code?: ApiErrorCode }).code ?? "UNKNOWN";
        const message = (rawError as { message?: string }).message ?? "Request failed";
        return new ApiError(status, code, message, payload);
      }
      if (typeof rawError === "string") {
        return new ApiError(status, "UNKNOWN", rawError, payload);
      }
    }
  }

  return new ApiError(status, "UNKNOWN", `Request failed with status ${status}`, payload);
}

async function parseResponsePayload(response: Response): Promise<unknown> {
  const contentType = response.headers.get("content-type") ?? "";
  if (!contentType.includes("application/json")) {
    const text = await response.text();
    return text || null;
  }

  return response.json();
}

export function createHttpClient(config: HttpClientConfig) {
  async function request<TResponse, TBody = unknown>(
    path: string,
    options: RequestOptions<TBody> = {}
  ): Promise<TResponse> {
    const {
      method = "GET",
      body,
      headers,
      query,
      auth = true,
      retryOnUnauthorized = true,
      signal
    } = options;

    const url = buildUrl(config.baseUrl, path, query);
    const defaultHeaders: Record<string, string> = {
      Accept: "application/json"
    };
    const accessToken = config.getAccessToken();
    if (auth && accessToken) {
      defaultHeaders.Authorization = `Bearer ${accessToken}`;
    }
    let requestBody: BodyInit | undefined;
    if (body !== undefined && body !== null) {
      if (body instanceof FormData) {
        requestBody = body;
      } else {
        defaultHeaders["Content-Type"] = "application/json";
        requestBody = JSON.stringify(body);
      }
    }

    const response = await fetch(url, {
      method,
      headers: {
        ...defaultHeaders,
        ...headers
      },
      body: requestBody,
      signal
    });

    if (response.status === 401 && retryOnUnauthorized && auth && config.refreshAccessToken) {
      const refreshed = await config.refreshAccessToken();
      if (!refreshed) {
        config.onAuthFailure?.();
      } else {
        return request<TResponse, TBody>(path, {
          ...options,
          retryOnUnauthorized: false
        });
      }
    }

    if (response.status === 204) {
      return undefined as TResponse;
    }

    const payload = await parseResponsePayload(response);

    if (!response.ok) {
      throw toApiError(response.status, payload);
    }

    return payload as TResponse;
  }

  return {
    request
  };
}
