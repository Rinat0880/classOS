import { api } from '../api/axios';
import type { Group, CreateGroupInput, UpdateGroupInput, User, WhitelistEntry } from '../../types';

export const groupsService = {
  async getAll(): Promise<Group[]> {
    const response = await api.get('/api/groups/');
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getById(id: number): Promise<Group> {
    const response = await api.get(`/api/groups/${id}`);
    return response.data?.data || response.data;
  },

  async create(data: CreateGroupInput): Promise<Group> {
    const response = await api.post('/api/groups', data);
    return response.data?.data || response.data;
  },

  async update(id: number, data: UpdateGroupInput): Promise<Group> {
    const response = await api.patch(`/api/groups/${id}`, data);
    return response.data?.data || response.data;
  },

  async delete(id: number): Promise<void> {
    await api.delete(`/api/groups/${id}`);
  },

  async getUsers(groupId: number): Promise<User[]> {
    const response = await api.get(`/api/groups/${groupId}/users`);
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getWhitelist(groupId: number): Promise<WhitelistEntry[]> {
    const response = await api.get(`/api/groups/${groupId}/whitelist`);
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },
};