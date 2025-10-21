// API клиент для работы с данными производительности
// Интеграция с backend Performance API

const API_BASE_URL = import.meta.env.DEV ? '/api' : 'http://localhost:8080/api'

export interface PerformanceDealer {
  id: string;
  name: string;
  city: string;
  srRub: string;
  salesProfit: number;
  salesMargin: number;
  autoSalesRevenue: string;
  rap: string;
  autoSalesProfitsRap: string;
  autoSalesMargin: number;
  marketingInvestment: number;
  ranking: number;
  autoSalesDecision: 'Needs Development' | 'Planned Result' | 'Find New Candidate' | 'Close Down';
}

export interface PerformanceResponse {
  dealers: PerformanceDealer[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface PerformanceFilters {
  search?: string;
  region?: string;
  decision?: string;
  page?: number;
  limit?: number;
}

/**
 * Получить данные производительности дилеров
 */
export async function getPerformanceData(filters?: PerformanceFilters): Promise<PerformanceResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.decision) params.append('decision', filters.decision);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());
  
  // Добавляем обязательные параметры quarter и year
  params.append('quarter', 'Q1');
  params.append('year', '2024');

  const response = await fetch(`${API_BASE_URL}/performance?${params.toString()}`, {
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    throw new Error(`Failed to fetch performance data: ${response.statusText}`);
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

// Хелперы

/**
 * Список доступных решений по производительности
 */
export const PERFORMANCE_DECISIONS = [
  'Planned Result',
  'Needs Development',
  'Find New Candidate',
  'Close Down'
] as const;

export type PerformanceDecision = typeof PERFORMANCE_DECISIONS[number];
