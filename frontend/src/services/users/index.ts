import { api } from '../api/axios';
import type { User } from '../../types';

export const usersService = {
  async getAll(): Promise<User[]> {
    const response = await api.get('/api/users/');
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getById(id: number): Promise<User> {
    const response = await api.get(`/api/users/${id}/`);
    return response.data?.data || response.data;
  },

  async update(id: number, data: Partial<User>): Promise<User> {
    const response = await api.patch(`/api/users/${id}`, data);
    return response.data?.data || response.data;
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/api/users/${id}`);
  },

  async changePassword(id: number, newPassword: string): Promise<void> {
    await api.post(`/api/users/${id}/password`, { password: newPassword });
  },
};
