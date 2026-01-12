import { api } from '../api/axios';
import { AUTH_ENDPOINTS, STORAGE_KEYS } from '../../constants';
import type { LoginRequest, LoginResponse, User } from '../../types';

export const authService = {
  async signIn(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post(AUTH_ENDPOINTS.SIGN_IN, credentials);
    const data = response.data?.data || response.data;
    
    if (data.token) {
      localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, data.token);
      localStorage.setItem(STORAGE_KEYS.USER_DATA, JSON.stringify(data.user));
    }
    
    return data;
  },

  logout(): void {
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER_DATA);
  },

  isAuthenticated(): boolean {
    return !!localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  },

  getCurrentUser(): User | null {
    const userData = localStorage.getItem(STORAGE_KEYS.USER_DATA);
    return userData ? JSON.parse(userData) : null;
  },

  getToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  },
};
