import { defineStore } from "pinia";
import { useCookies } from "@vueuse/integrations/useCookies";
import { APIClient } from "../services/APIClient";

export const useApiClientStore = defineStore("auth", {
  state: () => ({
    accessToken: "",
    refreshToken: "",
    expiryDate: "",
    apiClient: new APIClient(),
  }),

  actions: {
    initFromStorage() {
      this.accessToken = sessionStorage.getItem("access_token") || "";
      const cookies = useCookies();
      this.refreshToken = cookies.get("refresh_token") || "";
      this.expiryDate = sessionStorage.getItem("access_token_expiry") || "";
      this.apiClient.setTokens(this.accessToken, this.refreshToken, this.expiryDate);
    },

    setAccessToken(token: string) {
      this.accessToken = token;
      sessionStorage.setItem("access_token", token);
      this.apiClient.setTokens(token, this.refreshToken, this.expiryDate);
    },

    setRefreshToken(token: string) {
      this.refreshToken = token;
      const cookies = useCookies();
      cookies.set("refresh_token", token, { secure: true, sameSite: "strict" });
      this.apiClient.setTokens(this.accessToken, token, this.expiryDate);
    },

    setExpiryDate(expiry: string) {
      this.expiryDate = expiry;
      sessionStorage.setItem("access_token_expiry", expiry);
      this.apiClient.setTokens(this.accessToken, this.refreshToken, expiry);
    },

    /**
     * Logs a user in
     * @param credential The username or email of the user
     * @param password The password of the user
     */
    async login(credential: string, password: string): Promise<{ success: boolean; status: string }> {
      try {
        const response = await this.apiClient.login(credential, password);

        if (response && response.access_token) {
          this.setAccessToken(response.access_token);
          this.setExpiryDate(response.expiry_date);
          this.setRefreshToken(response.refresh_token);

          return { success: true, status: "Logged in successfully" };
        } else {
          return { success: false, status: "Failed to login: Invalid credentials" };
        }
      } catch (e: any) {
        return { success: false, status: `Failed to log in: ${e.message || e}` };
      }
    },

    async logout() {
      await this.apiClient.logout();
      this.$reset();
      sessionStorage.removeItem("access_token");
      sessionStorage.removeItem("access_token_expiry");
      const cookies = useCookies();
      cookies.remove("refresh_token");
    },

    isAuthenticated(): boolean {
      const now = new Date().getTime();
      const expiry = new Date(this.expiryDate).getTime();
      return now < expiry;
    },
  },

  persist: true,
});
