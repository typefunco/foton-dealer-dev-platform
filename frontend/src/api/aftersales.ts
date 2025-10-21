// API клиент для работы с данными After Sales
// Интеграция с backend After Sales API

const API_BASE_URL = import.meta.env.DEV ? '/api' : 'http://localhost:8080/api'

export interface AfterSalesDealer {
  id: string;
  name: string;
  city: string;
  rStockPercent: number;
  wStockPercent: number;
  flhPercent: number;
  flhSharePercent: string;
  warrantyHours: number;
  serviceContractsHours: number;
  asTrainings: boolean;
  csi: string;
  asDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down';
  sparePartsSalesQ3: string;
  sparePartsSalesYtd: string;
}

export interface AfterSalesResponse {
  dealers: AfterSalesDealer[];
  pagination: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

export interface AfterSalesFilters {
  search?: string;
  region?: string;
  decision?: string;
  page?: number;
  limit?: number;
}

/**
 * Получить данные After Sales дилеров
 */
export async function getAfterSalesData(filters?: AfterSalesFilters): Promise<AfterSalesResponse> {
  const params = new URLSearchParams();
  
  if (filters?.search) params.append('search', filters.search);
  if (filters?.region) params.append('region', filters.region);
  if (filters?.decision) params.append('decision', filters.decision);
  if (filters?.page) params.append('page', filters.page.toString());
  if (filters?.limit) params.append('limit', filters.limit.toString());
  
  // Добавляем обязательные параметры quarter и year
  params.append('quarter', 'Q1');
  params.append('year', '2024');

  const response = await fetch(`${API_BASE_URL}/aftersales?${params.toString()}`, {
    credentials: 'include', // Включаем cookies для аутентификации
  });
  
  if (!response.ok) {
    throw new Error(`Failed to fetch after sales data: ${response.statusText}`);
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
 * Список доступных решений по After Sales
 */
export const AFTER_SALES_DECISIONS = [
  'Planned Result',
  'Needs development',
  'Find New Candidate',
  'Close Down'
] as const;

export type AfterSalesDecision = typeof AFTER_SALES_DECISIONS[number];
