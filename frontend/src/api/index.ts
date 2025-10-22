// Универсальный API сервис для всех типов данных
// Интеграция с backend API

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || (import.meta.env.DEV ? '' : '');

export { API_BASE_URL };

// Универсальная функция для API запросов с поддержкой localStorage токенов
export async function apiRequest<T>(
  endpoint: string, 
  options: RequestInit = {}
): Promise<T> {
  const url = `${API_BASE_URL}${endpoint}`;
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const defaultOptions: RequestInit = {
    headers: {
      'Content-Type': 'application/json',
      ...(token && { 'Authorization': `Bearer ${token}` }),
      ...options.headers,
    },
  };

  const response = await fetch(url, { ...defaultOptions, ...options });

  if (!response.ok) {
    const error = await response.json().catch(() => ({ error: 'Request failed' }));
    throw new Error(error.error || `HTTP ${response.status}`);
  }

  return response.json();
}

export interface SearchFilters {
  region?: string;
  quarter?: string;
  year?: number;
  search?: string;
  page?: number;
  limit?: number;
}

export interface ApiResponse<T> {
  data: T[];
  pagination?: {
    page: number;
    limit: number;
    total: number;
    totalPages: number;
  };
}

// Универсальная функция для выполнения API запросов
export async function fetchData<T>(
  endpoint: string, 
  filters: SearchFilters = {}
): Promise<T[]> {
  const params = new URLSearchParams();
  
  // Добавляем параметры фильтрации
  if (filters.region && filters.region !== 'all') {
    params.append('region', filters.region);
  }
  if (filters.quarter) {
    params.append('quarter', filters.quarter);
  }
  if (filters.year) {
    params.append('year', filters.year.toString());
  }
  if (filters.search) {
    params.append('search', filters.search);
  }
  if (filters.page) {
    params.append('page', filters.page.toString());
  }
  if (filters.limit) {
    params.append('limit', filters.limit.toString());
  }

  const url = `${API_BASE_URL}/${endpoint}?${params.toString()}`;
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  try {
    const response = await fetch(url, {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      },
    });
    
    if (!response.ok) {
      throw new Error(`Failed to fetch ${endpoint}: ${response.statusText}`);
    }
    
    const data = await response.json();
    
    // Возвращаем массив данных (бэкенд возвращает массив)
    return Array.isArray(data) ? data : [];
  } catch (error) {
    console.error(`Error fetching ${endpoint}:`, error);
    throw error;
  }
}

// Специфичные функции для каждого типа данных
export async function fetchDealers(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('dealers', filters);
}

export async function fetchDealerDev(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('dealerdev', filters);
}

export async function fetchSales(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('sales', filters);
}

export async function fetchPerformance(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('performance', filters);
}

export async function fetchAfterSales(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('aftersales', filters);
}

export async function fetchSalesTeam(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('sales', filters);
}

export async function fetchQuarterComparison(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('quarter-comparison', filters);
}

export async function fetchAllData(filters: SearchFilters = {}): Promise<any[]> {
  return fetchData('all-data', filters);
}

// Маппинг параметров на API endpoints
export const PARAMETER_ENDPOINTS = {
  'all': 'all-data',
  'total-performance': 'all-data', 
  'dealer-development': 'dealer_dev',
  'sales': 'sales',
  'after-sales': 'after_sales',
  'performance': 'performance',
  'sales-team': 'sales',
  'quarter-comparison': 'quarter-comparison'
} as const;

// Маппинг регионов для соответствия бэкенду
export const REGION_MAPPING = {
  'all': 'all-russia',
  'central': 'Central',
  'north-west': 'North West',
  'volga': 'Volga',
  'south': 'South',
  'n-caucasus': 'Kavkaz',
  'ural': 'Ural',
  'siberia': 'Siberia',
  'far-east': 'Far East'
} as const;

// Маппинг кварталов для соответствия бэкенду
export const QUARTER_MAPPING = {
  'q1': 'Q1',
  'q2': 'Q2', 
  'q3': 'Q3',
  'q4': 'Q4'
} as const;

// Экспорт динамического API
export * from './dynamic';
