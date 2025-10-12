// API клиент для работы с данными сравнения кварталов
// Интеграция с backend Quarter Comparison API

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || '/api';

export interface QuarterComparisonData {
  dealerId: string;
  dealerName: string;
  region: string;
  currentQuarter: {
    quarter: string;
    year: number;
    sales: number;
    profit: number;
    margin: number;
    units: number;
  };
  previousQuarter: {
    quarter: string;
    year: number;
    sales: number;
    profit: number;
    margin: number;
    units: number;
  };
  comparison: {
    salesChange: number;
    profitChange: number;
    marginChange: number;
    unitsChange: number;
  };
  trend: 'up' | 'down' | 'stable';
}

export interface QuarterComparisonResponse {
  data: QuarterComparisonData[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface QuarterComparisonFilters {
  search?: string;
  region?: string;
  trend?: string;
  quarter?: string;
  year?: number;
  page?: number;
  limit?: number;
}

/**
 * Получить данные сравнения кварталов
 */
export async function getQuarterComparison(filters?: QuarterComparisonFilters): Promise<QuarterComparisonResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.trend) params.append('trend', filters.trend);
  if (filters?.quarter) params.append('quarter', filters.quarter);
  if (filters?.year) params.append('year', filters.year.toString());
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());
  
  // Добавляем обязательные параметры quarter и year
  params.append('quarter', filters?.quarter || 'Q1');
  params.append('year', (filters?.year || 2024).toString());

  const response = await fetch(`${API_BASE_URL}/quarter-comparison?${params.toString()}`);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch quarter comparison data: ${response.statusText}`);
  }
  
  const data = await response.json();
  
  // Преобразуем ответ бэкенда в формат, ожидаемый фронтендом
  return {
    data: Array.isArray(data) ? data : [],
    pagination: {
      page: filters?.page || 1,
      limit: filters?.limit || 10,
      total: Array.isArray(data) ? data.length : 0,
      totalPages: Math.ceil((Array.isArray(data) ? data.length : 0) / (filters?.limit || 10))
    }
  };
}

// Хелперы

/**
 * Список доступных трендов
 */
export const QUARTER_TRENDS = ['up', 'down', 'stable'] as const;

/**
 * Список доступных кварталов
 */
export const QUARTERS = ['Q1', 'Q2', 'Q3', 'Q4'] as const;

export type QuarterTrend = typeof QUARTER_TRENDS[number];
export type Quarter = typeof QUARTERS[number];
