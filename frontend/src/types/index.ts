// User типы
export interface User {
  id: number;
  name: string;
  username: string;
  password?: string;
  role: 'admin' | 'client';
  group_id: number;
  group_name: string;
}

export interface UpdateUserInput {
  name?: string;
  username?: string;
  password?: string;
  role?: 'admin' | 'client';
  group_id?: number | null;
  group_name?: string;
}

export interface CreateUserInput {
  name: string;
  username: string;
  password: string;
  role: 'admin' | 'client';
  group_id: number;
}

export interface ChangePasswordInput {
  password: string;
}

export type UserStatus = 'online' | 'offline';

export interface UserWithStatus extends User {
  status: UserStatus;
}

// Group типы
export interface Group {
  id: number;
  name: string;
}

export interface GroupWithUserCount extends Group {
  userCount: number;
}

export interface CreateGroupInput {
  name: string;
}

export interface UpdateGroupInput {
  name?: string;
}

// Whitelist типы
export interface WhitelistEntry {
  id: number;
  group_id: number;
  value: string;
  created_at: string;
}

// Settings типы
export interface Settings {
  id: number;
  key: string;
  value: string;
  updated_at: string;
}

// Auth типы
export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
  user: User;
}

export interface SignUpRequest {
  name: string;
  username: string;
  password: string;
  role?: 'admin' | 'client';
}

// API Response типы
export interface ApiError {
  message: string;
  error?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  pageSize: number;
}

// Dashboard Statistics
export interface DashboardStats {
  totalUsers: number;
  adminUsers: number;
  clientUsers: number;
  activeGroups: number;
  onlineUsers: number; 
}
