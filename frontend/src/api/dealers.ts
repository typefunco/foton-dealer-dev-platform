// API клиент для работы с дилерами
// Интеграция с backend Dealer API

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

export interface Dealer {
  dealer_id: number;
  dealer_name_ru: string;
  region: string;
  city: string;
  manager: string;
  created_at: string;
  updated_at: string;
}

export interface DealerCard {
  dealer_id: number;
  dealer_name_ru: string;
  region: string;
  city: string;
  manager: string;
  created_at: string;
  updated_at: string;
}

export interface DealersResponse {
  dealers: Dealer[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface DealerFilters {
  search?: string;
  region?: string;
  class?: string;
  page?: number;
  limit?: number;
  quarter?: string;
  year?: number;
}

/**
 * Получить список дилеров с фильтрацией
 */
export async function getDealers(filters?: DealerFilters): Promise<DealersResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.class) params.append('class', filters.class);
  if (filters?.quarter) params.append('quarter', filters.quarter);
  if (filters?.year) params.append('year', filters.year.toString());
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());

  const response = await fetch(`${API_BASE_URL}/dealers?${params.toString()}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch dealers: ${response.statusText}`);
  }
  
  const dealers = await response.json();
  
  // Преобразуем ответ бэкенда в формат, ожидаемый фронтендом
  return {
    dealers: Array.isArray(dealers) ? dealers : [],
    pagination: {
      page: filters?.page || 1,
      limit: filters?.limit || 10,
      total: Array.isArray(dealers) ? dealers.length : 0,
      totalPages: Math.ceil((Array.isArray(dealers) ? dealers.length : 0) / (filters?.limit || 10))
    }
  };
}

/**
 * Получить дилера по ID
 */
export async function getDealerById(id: string): Promise<Dealer> {
  const response = await fetch(`${API_BASE_URL}/dealers/${id}`);
  
  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('Dealer not found');
    }
    throw new Error(`Failed to fetch dealer: ${response.statusText}`);
  }
  
  return response.json();
}

/**
 * Получить полную карточку дилера
 */
export async function getDealerCard(id: string): Promise<DealerCard> {
  const response = await fetch(`${API_BASE_URL}/dealers/${id}/card`);
  
  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('Dealer card not found');
    }
    throw new Error(`Failed to fetch dealer card: ${response.statusText}`);
  }
  
  return response.json();
}

// Хелперы

/**
 * Список доступных классов дилеров
 */
export const DEALER_CLASSES = ['A', 'B', 'C', 'D'] as const;

/**
 * Список доступных рекомендаций
 */
export const DEALER_RECOMMENDATIONS = [
  'Planned Result',
  'Needs Development', 
  'Find New Candidate',
  'Close Down'
] as const;

export type DealerClass = typeof DEALER_CLASSES[number];
export type DealerRecommendation = typeof DEALER_RECOMMENDATIONS[number];
