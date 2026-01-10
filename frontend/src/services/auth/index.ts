import { api } from '../api/axios';
import { AUTH_ENDPOINTS, STORAGE_KEYS } from '../../constants';
import type { LoginRequest, LoginResponse, SignUpRequest, User } from '../../types';

export const authService = {
  // Вход
  async signIn(credentials: LoginRequest): Promise<LoginResponse> {
    const response = await api.post<LoginResponse>(AUTH_ENDPOINTS.SIGN_IN, credentials);
    
    // Сохраняем токен и данные пользователя
    if (response.data.token) {
      localStorage.setItem(STORAGE_KEYS.ACCESS_TOKEN, response.data.token);
      localStorage.setItem(STORAGE_KEYS.USER_DATA, JSON.stringify(response.data.user));
    }
    
    return response.data;
  },

  // Регистрация
  async signUp(data: SignUpRequest): Promise<User> {
    const response = await api.post<User>(AUTH_ENDPOINTS.SIGN_UP, data);
    return response.data;
  },

  // Выход
  logout(): void {
    localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
    localStorage.removeItem(STORAGE_KEYS.USER_DATA);
  },

  // Проверка авторизации
  isAuthenticated(): boolean {
    return !!localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  },

  // Получить текущего пользователя из localStorage
  getCurrentUser(): User | null {
    const userData = localStorage.getItem(STORAGE_KEYS.USER_DATA);
    return userData ? JSON.parse(userData) : null;
  },

  // Получить токен
  getToken(): string | null {
    return localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  },
};
