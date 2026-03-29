import axios, { AxiosError, AxiosRequestConfig } from "axios";

export type ApiEnvelope<T> = {
  data: T;
  partial?: boolean;
  degradedFields?: string[];
};

type GoFrameResponse<T> = {
  code: number;
  message?: string;
  data: T;
};

class ApiResponseError extends Error {
  constructor(public code: number, message?: string) {
    super(message ?? "API request failed");
    this.name = "ApiResponseError";
  }
}

export type ApiHooks = {
  getToken?: () => string | null;
  onUnauthorized?: () => void;
  onDegraded?: (fields: string[]) => void;
};

export type ApiRequestConfig = AxiosRequestConfig & {
  silentDegraded?: boolean;
};

export type CreateApiClientOptions = {
  baseURL?: string;
  timeoutMs?: number;
  hooks?: ApiHooks;
};

export class ApiClient {
  private hooks: ApiHooks = {};

  private readonly http;

  constructor(options?: CreateApiClientOptions) {
    this.hooks = options?.hooks ?? {};
    this.http = axios.create({
      baseURL: options?.baseURL ?? process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8000",
      withCredentials: true,
      timeout: options?.timeoutMs ?? 15000
    });

    this.http.interceptors.request.use((config) => {
      const token = this.hooks.getToken?.();
      if (token) {
        const headers = config.headers ?? {};
        (headers as Record<string, string>).Authorization = `Bearer ${token}`;
        config.headers = headers;
      }
      return config;
    });

    this.http.interceptors.response.use(
      (response) => response,
      (error: AxiosError) => {
        if (error.response?.status === 401) {
          this.hooks.onUnauthorized?.();
        }
        return Promise.reject(error);
      }
    );
  }

  configure(hooks: Partial<ApiHooks>) {
    this.hooks = {
      ...this.hooks,
      ...hooks
    };
  }

  async get<T>(url: string, config?: ApiRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "GET", url });
  }

  async post<T>(url: string, data?: unknown, config?: ApiRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "POST", url, data });
  }

  async put<T>(url: string, data?: unknown, config?: ApiRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "PUT", url, data });
  }

  async patch<T>(url: string, data?: unknown, config?: ApiRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "PATCH", url, data });
  }

  async delete<T>(url: string, config?: ApiRequestConfig): Promise<T> {
    return this.request<T>({ ...config, method: "DELETE", url });
  }

  private async request<T>(config: ApiRequestConfig): Promise<T> {
    const response = await this.http.request<ApiEnvelope<T> | GoFrameResponse<ApiEnvelope<T> | T> | T>(config);
    const { result, partial, degradedFields } = this.unwrapPayload<T>(response.data);

    if (partial && !config.silentDegraded) {
      this.hooks.onDegraded?.(degradedFields ?? []);
    }

    return result;
  }

  private unwrapPayload<T>(payload: unknown): { result: T; partial: boolean; degradedFields?: string[] } {
    let partial = false;
    let degradedFields: string[] | undefined;
    let value = payload;

    if (this.isGoFrameResponse(value)) {
      const { code, message, data } = value;
      if (code !== 0) {
        throw new ApiResponseError(code, message);
      }
      value = data;
    }

    if (value && typeof value === "object" && "data" in value) {
      const envelope = value as ApiEnvelope<T>;
      partial = partial || Boolean(envelope.partial);
      if (envelope.degradedFields?.length) {
        degradedFields = envelope.degradedFields;
      }
      value = envelope.data;
    }

    if (!partial && value && typeof value === "object") {
      const maybePartial = (value as Record<string, unknown>).partial;
      if (typeof maybePartial === "boolean") {
        partial = partial || maybePartial;
      }
      if (!degradedFields) {
        const maybeFields = (value as Record<string, unknown>).degradedFields;
        if (Array.isArray(maybeFields)) {
          degradedFields = maybeFields.map((field) => String(field));
        }
      }
    }

    return { result: value as T, partial, degradedFields };
  }

  private isGoFrameResponse<T>(value: unknown): value is GoFrameResponse<T> {
    return (
      Boolean(value) &&
      typeof value === "object" &&
      "code" in (value as Record<string, unknown>) &&
      "data" in (value as Record<string, unknown>) &&
      typeof (value as Record<string, unknown>).code === "number"
    );
  }
}

export function createApiClient(options?: CreateApiClientOptions) {
  return new ApiClient(options);
}

// Keep a default singleton for backward compatibility.
export const apiClient = createApiClient();
