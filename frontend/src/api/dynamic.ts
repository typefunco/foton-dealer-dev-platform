import axios from 'axios'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || ''

// Создаем экземпляр axios для API запросов
const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

// Добавляем interceptor для автоматического добавления токена
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('auth_token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

// Добавляем interceptor для обработки ошибок авторизации
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      // Токен невалиден или истек - очищаем localStorage и перенаправляем на логин
      console.warn('Unauthorized: token expired or invalid')
      localStorage.removeItem('auth_token')
      // Перенаправляем на страницу логина, сохранив текущий путь для редиректа после логина
      const currentPath = window.location.pathname + window.location.search
      window.location.href = `/login?redirect=${encodeURIComponent(currentPath)}`
    }
    return Promise.reject(error)
  }
)

// Типы для динамического API
export interface DynamicTableParams {
  year?: number
  quarter?: string
  regions?: string[]
  region?: string // для обратной совместимости
  dealer_ids?: number[]
  limit?: number
  offset?: number
  sort_by?: string
  sort_order?: 'asc' | 'desc'
}

export interface DynamicDataResponse<T = any> {
  tableType: string
  year: number
  quarter: string
  regions: string[]
  dealer_ids: number[]
  data: T[]
  count: number
}

// Типы таблиц
export type TableType = 'dealer_dev' | 'sales' | 'after_sales' | 'performance' | 'sales_team'

// Универсальная функция для получения данных из таблиц
export const getDynamicData = async <T = any>(
  tableType: TableType,
  params: DynamicTableParams = {}
): Promise<DynamicDataResponse<T>> => {
  const searchParams = new URLSearchParams()
  
  // Добавляем параметры в URL
  Object.entries(params).forEach(([key, value]) => {
    if (value !== undefined && value !== null && value !== '') {
      if (Array.isArray(value)) {
        // Для массивов объединяем через запятую
        searchParams.append(key, value.join(','))
      } else {
        searchParams.append(key, value.toString())
      }
    }
  })
  
  const queryString = searchParams.toString()
  const url = `${API_BASE_URL}/api/${tableType}${queryString ? `?${queryString}` : ''}`
  
  const response = await api.get<DynamicDataResponse<T>>(url)
  return response.data
}

// Специализированные функции для каждого типа таблицы

// Dealer Development
export const getDealerDevData = (params: DynamicTableParams = {}) => 
  getDynamicData('dealer_dev', params)

// Sales
export const getSalesData = (params: DynamicTableParams = {}) => 
  getDynamicData('sales', params)

// After Sales
export const getAfterSalesData = (params: DynamicTableParams = {}) => 
  getDynamicData('after_sales', params)

// Performance
export const getPerformanceData = (params: DynamicTableParams = {}) => 
  getDynamicData('performance', params)

// Sales Team
export const getSalesTeamData = (params: DynamicTableParams = {}) => 
  getDynamicData('sales_team', params)

// Утилиты для работы с параметрами
export const buildDynamicParams = (filters: {
  region?: string
  regions?: string[]
  quarter?: string
  year?: number
  dealers?: string[]
  limit?: number
  offset?: number
  sortBy?: string
  sortOrder?: 'asc' | 'desc'
}): DynamicTableParams => {
  const params: DynamicTableParams = {}
  
  // Поддержка множественных регионов
  if (filters.regions && filters.regions.length > 0) {
    params.regions = filters.regions
  } else if (filters.region) {
    params.region = filters.region // обратная совместимость
  }
  
  if (filters.quarter) params.quarter = filters.quarter
  if (filters.year) params.year = filters.year
  
  // Преобразуем строковые ID дилеров в числовые
  if (filters.dealers && filters.dealers.length > 0) {
    params.dealer_ids = filters.dealers.map(id => parseInt(id)).filter(id => !isNaN(id))
  }
  
  if (filters.limit) params.limit = filters.limit
  if (filters.offset) params.offset = filters.offset
  if (filters.sortBy) params.sort_by = filters.sortBy
  if (filters.sortOrder) params.sort_order = filters.sortOrder
  
  return params
}

// Маппинг типов таблиц для фронтенда
export const TABLE_TYPE_MAPPING = {
  'dealer-development': 'dealer_dev' as TableType,
  'sales': 'sales' as TableType,
  'after-sales': 'after_sales' as TableType,
  'performance': 'performance' as TableType,
  'sales-team': 'sales_team' as TableType,
} as const

// Получение типа таблицы по ключу фронтенда
export const getTableTypeFromKey = (key: string): TableType | null => {
  return TABLE_TYPE_MAPPING[key as keyof typeof TABLE_TYPE_MAPPING] || null
}
