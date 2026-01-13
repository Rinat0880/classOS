import { api } from '../api/axios';
import type { DeviceStatus } from '../../types';
import { DEVICES_ENDPOINTS } from '../../constants';

export const devicesService = {
  async getAll(): Promise<DeviceStatus[]> {
    const response = await api.get(DEVICES_ENDPOINTS.GET_ALL);
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getOnline(): Promise<DeviceStatus[]> {
    const response = await api.get(DEVICES_ENDPOINTS.GET_ONLINE);
    const data = response.data?.data || response.data;
    return Array.isArray(data) ? data : [];
  },

  async getByName(name: string): Promise<DeviceStatus> {
    const response = await api.get(DEVICES_ENDPOINTS.GET_BY_NAME(name));
    return response.data?.data || response.data;
  },
};
