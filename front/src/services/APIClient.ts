interface ClientOptions {
  method?: string;
  body?: any;
  queryParams?: Record<string, string | number | boolean>;
  headers?: Headers;
}

interface LoginBody {
  email?: string;
  username?: string;
  password: string;
}

export class APIClient {
  private baseUrl: string = "/api/v1";
  private accessToken: string = "";
  private refreshToken: string = "";
  private expiryDate: string = "";

  constructor() {
    // Access tokens should be managed in Pinia
  }

  public setTokens(accessToken: string, refreshToken: string, expiryDate: string): void {
    this.accessToken = accessToken;
    this.refreshToken = refreshToken;
    this.expiryDate = expiryDate;
  }

  public clearTokens(): void {
    this.accessToken = "";
    this.refreshToken = "";
    this.expiryDate = "";
  }

  private buildHeaders(headers: Headers): Headers {
    if (this.accessToken) {
      headers.set("Authorization", `Bearer ${this.accessToken}`);
    }
    return headers;
  }

  private async request(
    endpoint: string,
    { method = "GET", body = null, queryParams = {}, headers = new Headers() }: ClientOptions = {}
  ): Promise<any> {
    headers = this.buildHeaders(headers);

    const queryString = Object.keys(queryParams)
      .map((key) => `${encodeURIComponent(key)}=${encodeURIComponent(queryParams[key])}`)
      .join("&");

    const url = `${this.baseUrl}${endpoint}${queryString ? `?${queryString}` : ""}`;

    if (body && method !== "GET") {
      headers.set("Content-Type", "application/json");
    }

    const response = await fetch(url, {
      method,
      headers,
      body: body ? JSON.stringify(body) : null,
    });

    if (response.status === 401) {
      const newAccessToken = await this.getNewAccessToken();
      if (newAccessToken) {
        this.setTokens(newAccessToken, this.refreshToken, this.expiryDate);
        headers.set("Authorization", `Bearer ${newAccessToken}`);

        // Retry the request with the new access token
        const retryResponse = await fetch(url, {
          method,
          headers,
          body: body ? JSON.stringify(body) : null,
        });

        return retryResponse.json();
      } else {
        throw new Error("Unauthorized: Please log in again.");
      }
    }

    if (!response.ok) {
      throw new Error(`API request failed: ${response.status}`);
    }

    return response.json();
  }

  private async getNewAccessToken(): Promise<string | null> {
    const response = await this.request("/auth/refresh", {
      method: "POST",
      headers: new Headers({ "Refresh_Token": this.refreshToken }),
    });

    return response?.access_token || null;
  }

  public async login(credential: string, password: string): Promise<any> {
    const body: LoginBody = {
      [credential.includes("@") ? "email" : "username"]: credential,
      password,
    };

    return await this.request("/auth/login", {
      method: "POST",
      body,
    });
  }

  public async logout(): Promise<void> {
    await this.request("/auth/logout", {
      method: "POST",
      headers: new Headers({ "Refresh-Token": this.refreshToken }),
    });
    this.clearTokens();
  }
}
