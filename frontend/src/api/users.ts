// API клиент для работы с пользователями
// Интеграция с backend User Management API

const API_BASE_URL = import.meta.env.DEV ? '/api' : 'http://localhost:8080/api'

export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  region: string;
  position: string;
  createdAt: string;
  status: 'active' | 'inactive';
}

export interface CreateUserRequest {
  email: string;
  firstName: string;
  lastName: string;
  region: string;
  position: string;
}

export interface UpdateUserRequest {
  email?: string;
  firstName?: string;
  lastName?: string;
  region?: string;
  position?: string;
  status?: 'active' | 'inactive';
}

export interface UserCredentials {
  email: string;
  password: string;
}

export interface CreateUserResponse {
  user: User;
  credentials?: UserCredentials;
}

export interface UsersResponse {
  users: User[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface RegionStats {
  region: string;
  userCount: number;
  users: User[];
}

export interface UserStatsResponse {
  totalUsers: number;
  regionStats: RegionStats[];
}

export interface UserFilters {
  search?: string;
  region?: string;
  position?: string;
  page?: number;
  limit?: number;
}

/**
 * Получить список пользователей с фильтрацией
 */
export async function getUsers(filters?: UserFilters): Promise<UsersResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.position) params.append('position', filters.position);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());

  const response = await fetch(`${API_BASE_URL}/users?${params.toString()}`, {
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    throw new Error(`Failed to fetch users: ${response.statusText}`);
  }
  
  return response.json();
}

/**
 * Получить пользователя по ID
 */
export async function getUserById(id: string): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/users/${id}`, {
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('User not found');
    }
    throw new Error(`Failed to fetch user: ${response.statusText}`);
  }
  
  return response.json();
}

/**
 * Создать нового пользователя
 */
export async function createUser(data: CreateUserRequest): Promise<CreateUserResponse> {
  const response = await fetch(`${API_BASE_URL}/users`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Включаем cookies для аутентификации
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to create user');
  }
  
  return response.json();
}

/**
 * Обновить пользователя
 */
export async function updateUser(id: string, data: UpdateUserRequest): Promise<User> {
  const response = await fetch(`${API_BASE_URL}/users/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Включаем cookies для аутентификации
    body: JSON.stringify(data),
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to update user');
  }
  
  return response.json();
}

/**
 * Удалить пользователя
 */
export async function deleteUser(id: string): Promise<void> {
  const response = await fetch(`${API_BASE_URL}/users/${id}`, {
    method: 'DELETE',
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to delete user');
  }
}

/**
 * Получить статистику пользователей
 */
export async function getUserStats(): Promise<UserStatsResponse> {
  const response = await fetch(`${API_BASE_URL}/users/stats`, {
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    throw new Error(`Failed to fetch user stats: ${response.statusText}`);
  }
  
  return response.json();
}

// Хелперы

/**
 * Список доступных регионов
 */
export const REGIONS = [
  'Central',
  'Caucasus',
  'Volga',
  'Ural',
  'Siberia',
  'Far East',
  'North-West',
  'South',
] as const;

/**
 * Список доступных должностей
 */
export const POSITIONS = [
  'Sales Manager',
  'Regional Director',
  'Regional Manager',
  'Sales Director',
  'Account Manager',
  'Sales Representative',
  'Account Executive',
  'Senior Sales Manager',
  'Head of Sales',
  'Business Development Manager',
] as const;

export type Region = typeof REGIONS[number];
export type Position = typeof POSITIONS[number];

