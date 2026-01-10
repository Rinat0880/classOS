import { api } from '../api/axios';
import type { User } from '../../types';

export const usersService = {
  // Получить всех пользователей
  async getAll(): Promise<User[]> {
    const response = await api.get<User[]>('/api/users/');
    return response.data;
  },

  // Получить пользователя по ID
  async getById(id: number): Promise<User> {
    const response = await api.get<User>(`/api/users/${id}`);
    return response.data;
  },

  // Обновить пользователя
  async update(id: number, data: Partial<User>): Promise<User> {
    const response = await api.patch<User>(`/api/users/${id}`, data);
    return response.data;
  },

  // Удалить пользователя
  async delete(id: number): Promise<void> {
    await api.delete(`/api/users/${id}`);
  },

  // Изменить пароль
  async changePassword(id: number, newPassword: string): Promise<void> {
    await api.post(`/api/users/${id}/password`, { password: newPassword });
  },
};