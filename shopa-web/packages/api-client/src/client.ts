import axios, { AxiosError, AxiosRequestConfig } from "axios";

export type ApiEnvelope<T> = {
  data: T;
  partial?: boolean;
  degradedFields?: string[];
};

type ApiHooks = {
  getToken?: () => string | null;
  onUnauthorized?: () => void;
  onDegraded?: (fields: string[]) => void;
};

export type ApiRequestConfig = AxiosRequestConfig & {
  silentDegraded?: boolean;
};

const baseURL = process.env.NEXT_PUBLIC_API_BASE_URL ?? "http://localhost:8000";

class ApiClient {
  private hooks: ApiHooks = {};

  private readonly http = axios.create({
    baseURL,
    withCredentials: true,
    timeout: 15000
  });

  constructor() {
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
    const response = await this.http.request<ApiEnvelope<T> | T>(config);
    const payload = response.data;

    if (payload && typeof payload === "object" && "data" in payload) {
      const envelope = payload as ApiEnvelope<T>;
      if (envelope.partial && !config.silentDegraded) {
        this.hooks.onDegraded?.(envelope.degradedFields ?? []);
      }
      return envelope.data;
    }

    return payload as T;
  }
}

export const apiClient = new ApiClient();
