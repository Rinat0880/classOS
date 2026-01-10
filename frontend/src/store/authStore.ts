import { create } from 'zustand';
import type { User } from '../types';
import { authService } from '../services/auth';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  setUser: (user: User | null) => void;
  logout: () => void;
  initialize: () => void;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,

  // Установить пользователя
  setUser: (user) => set({ user, isAuthenticated: !!user }),

  // Выход
  logout: () => {
    authService.logout();
    set({ user: null, isAuthenticated: false });
  },

  // Инициализация при загрузке приложения
  initialize: () => {
    const user = authService.getCurrentUser();
    const isAuthenticated = authService.isAuthenticated();
    set({ user, isAuthenticated });
  },
}));
