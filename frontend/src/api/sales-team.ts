// API клиент для работы с данными команды продаж
// Интеграция с backend Sales Team API

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

export interface SalesTeamMember {
  id: string;
  name: string;
  position: string;
  region: string;
  email: string;
  phone: string;
  performance: {
    salesTarget: number;
    salesAchieved: number;
    achievementRate: number;
  };
  status: 'active' | 'inactive';
}

export interface SalesTeamResponse {
  members: SalesTeamMember[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface SalesTeamFilters {
  search?: string;
  region?: string;
  position?: string;
  status?: string;
  page?: number;
  limit?: number;
}

/**
 * Получить данные команды продаж
 */
export async function getSalesTeamData(filters?: SalesTeamFilters): Promise<SalesTeamResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.position) params.append('position', filters.position);
  if (filters?.status) params.append('status', filters.status);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());
  
  // Добавляем обязательные параметры quarter и year
  params.append('quarter', 'Q1');
  params.append('year', '2024');

  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(`${API_BASE_URL}/api/sales?${params.toString()}`, {
    headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      }
  });
  
  if (!response.ok) {
    throw new Error(`Failed to fetch sales team data: ${response.statusText}`);
  }
  
  const members = await response.json();
  
  // Преобразуем ответ бэкенда в формат, ожидаемый фронтендом
  return {
    members: Array.isArray(members) ? members : [],
    pagination: {
      page: filters?.page || 1,
      limit: filters?.limit || 10,
      total: Array.isArray(members) ? members.length : 0,
      totalPages: Math.ceil((Array.isArray(members) ? members.length : 0) / (filters?.limit || 10))
    }
  };
}

// Хелперы

/**
 * Список доступных статусов
 */
export const SALES_TEAM_STATUSES = ['active', 'inactive'] as const;

export type SalesTeamStatus = typeof SALES_TEAM_STATUSES[number];
