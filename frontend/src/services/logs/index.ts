import { api } from '../api/axios';
import type { UserLog, LogsFilter } from '../../types';
import { LOGS_ENDPOINTS } from '../../constants';

export const logsService = {
  async getAll(filters?: LogsFilter): Promise<UserLog[]> {
    const response = await api.get(LOGS_ENDPOINTS.GET_ALL, { params: filters });
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getByUsername(username: string, limit?: number): Promise<UserLog[]> {
    const response = await api.get(LOGS_ENDPOINTS.GET_BY_USERNAME(username.toLowerCase()), {
      params: { limit },
    });
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getByDevice(device: string, limit?: number): Promise<UserLog[]> {
    const response = await api.get(LOGS_ENDPOINTS.GET_BY_DEVICE(device), {
      params: { limit },
    });
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },
};
