// API Base URL - будет заменен через env переменные
export const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

// Auth endpoints
export const AUTH_ENDPOINTS = {
  SIGN_UP: '/auth/sign-up',
  SIGN_IN: '/auth/sign-in',
} as const;

// Groups endpoints
export const GROUPS_ENDPOINTS = {
  GET_ALL: '/api/groups',
  CREATE: '/api/groups',
  GET_BY_ID: (id: number) => `/api/groups/${id}`,
  UPDATE: (id: number) => `/api/groups/${id}`,
  DELETE: (id: number) => `/api/groups/${id}`,
  
  // Users в группе
  GET_USERS: (groupId: number) => `/api/groups/${groupId}/users`,
  CREATE_USER: (groupId: number) => `/api/groups/${groupId}/users`,
  
  // Whitelist
  GET_WHITELIST: (groupId: number) => `/api/groups/${groupId}/whitelist`,
  CREATE_WHITELIST: (groupId: number) => `/api/groups/${groupId}/whitelist`,
  GET_WHITELIST_BY_ID: (groupId: number, whitelistId: number) => 
    `/api/groups/${groupId}/whitelist/${whitelistId}`,
  UPDATE_WHITELIST: (groupId: number, whitelistId: number) => 
    `/api/groups/${groupId}/whitelist/${whitelistId}`,
  DELETE_WHITELIST: (groupId: number, whitelistId: number) => 
    `/api/groups/${groupId}/whitelist/${whitelistId}`,
} as const;

// Users endpoints
export const USERS_ENDPOINTS = {
  GET_BY_ID: (id: number) => `/api/users/${id}`,
  UPDATE: (id: number) => `/api/users/${id}`,
  DELETE: (id: number) => `/api/users/${id}`,
  CHANGE_PASSWORD: (id: number) => `/api/users/${id}/password`,
} as const;

// Admin endpoints
export const ADMIN_ENDPOINTS = {
  SYNC_AD: '/api/admin/sync',
  CHECK_AD_STATUS: '/api/admin/ad/status',
} as const;

// Local storage keys
export const STORAGE_KEYS = {
  ACCESS_TOKEN: 'access_token',
  USER_DATA: 'user_data',
} as const;

// Routes
export const ROUTES = {
  LOGIN: '/login',
  DASHBOARD: '/',
  GROUPS: '/groups',
  USERS: '/users',
  SETTINGS: '/settings',
} as const;
