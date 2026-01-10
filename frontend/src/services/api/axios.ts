import axios, { type AxiosError, type InternalAxiosRequestConfig } from 'axios';
import { API_BASE_URL, STORAGE_KEYS, ROUTES } from '../../constants';


// Создаем инстанс axios
export const api = axios.create({
  baseURL: API_BASE_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor - добавляем токен к каждому запросу
api.interceptors.request.use(
  (config: InternalAxiosRequestConfig) => {
    const token = localStorage.getItem(STORAGE_KEYS.ACCESS_TOKEN);
    
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    return config;
  },
  (error: AxiosError) => {
    return Promise.reject(error);
  }
);

// Response interceptor - обработка ошибок
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    // Если 401 - разлогиниваем пользователя
    if (error.response?.status === 401) {
      localStorage.removeItem(STORAGE_KEYS.ACCESS_TOKEN);
      localStorage.removeItem(STORAGE_KEYS.USER_DATA);
      
      // Редирект на страницу логина
      if (window.location.pathname !== ROUTES.LOGIN) {
        window.location.href = ROUTES.LOGIN;
      }
    }
    
    return Promise.reject(error);
  }
);

export default api;
