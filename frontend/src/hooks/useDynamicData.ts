import { useState, useEffect, useCallback } from 'react'
import { 
  getDynamicData, 
  DynamicTableParams, 
  DynamicDataResponse, 
  TableType,
  buildDynamicParams 
} from '../api/dynamic'

// Экспортируем buildDynamicParams для использования в компонентах
export { buildDynamicParams }

interface UseDynamicDataOptions {
  tableType: TableType
  params?: DynamicTableParams
  enabled?: boolean
}

interface UseDynamicDataReturn<T = any> {
  data: T[] | null
  response: DynamicDataResponse<T> | null
  loading: boolean
  error: string | null
  refetch: () => void
  updateParams: (newParams: DynamicTableParams) => void
}

export const useDynamicData = <T = any>({
  tableType,
  params = {},
  enabled = true
}: UseDynamicDataOptions): UseDynamicDataReturn<T> => {
  const [data, setData] = useState<T[] | null>(null)
  const [response, setResponse] = useState<DynamicDataResponse<T> | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [currentParams, setCurrentParams] = useState<DynamicTableParams>(params)

  const fetchData = useCallback(async () => {
    if (!enabled) return

    setLoading(true)
    setError(null)

    try {
      // Добавляем тайм-аут для запроса
      const timeoutPromise = new Promise((_, reject) => {
        setTimeout(() => reject(new Error('Request timeout after 10 seconds')), 10000)
      })

      const dataPromise = getDynamicData<T>(tableType, currentParams)
      
      const result = await Promise.race([dataPromise, timeoutPromise]) as any
      
      setResponse(result)
      setData(result.data)
      
      // Проверяем, есть ли данные
      if (!result.data || result.data.length === 0) {
        console.warn(`No data found for ${tableType} with params:`, currentParams)
      }
      
    } catch (err) {
      console.error(`Error fetching ${tableType} data:`, err)
      let errorMessage = 'Unknown error occurred'
      
      if (err instanceof Error) {
        errorMessage = err.message
      } else if (typeof err === 'object' && err !== null) {
        // Обработка ошибок axios
        if ('response' in err) {
          const axiosError = err as any
          errorMessage = `Request failed with status code ${axiosError.response?.status}`
          if (axiosError.response?.data?.error) {
            errorMessage += `: ${axiosError.response.data.error}`
          }
        } else if ('message' in err) {
          errorMessage = (err as any).message
        }
      }
      
      setError(errorMessage)
    } finally {
      setLoading(false)
    }
  }, [tableType, currentParams, enabled])

  // Мемоизируем refetch для предотвращения ненужных перерендеров
  const refetch = useCallback(() => {
    fetchData()
  }, [fetchData])

  const updateParams = useCallback((newParams: DynamicTableParams) => {
    setCurrentParams(prev => {
      // Проверяем, действительно ли параметры изменились
      const hasChanged = Object.keys(newParams).some(key => 
        prev[key as keyof DynamicTableParams] !== newParams[key as keyof DynamicTableParams]
      )
      return hasChanged ? { ...prev, ...newParams } : prev
    })
  }, [fetchData]) // Добавляем fetchData в зависимости

  // Единственный эффект для загрузки данных
  useEffect(() => {
    if (enabled) {
      fetchData()
    }
  }, [currentParams, enabled, fetchData]) // Добавляем fetchData в зависимости

  // Cleanup эффект для предотвращения утечек памяти
  useEffect(() => {
    return () => {
      // Сбрасываем состояние при размонтировании компонента
      setLoading(false)
      setError(null)
    }
  }, [])

  return {
    data,
    response,
    loading,
    error,
    refetch,
    updateParams
  }
}

// Специализированные хуки для каждого типа таблицы

export const useDealerDevData = (params?: DynamicTableParams, enabled?: boolean) =>
  useDynamicData({ tableType: 'dealer_dev', params, enabled })

export const useSalesData = (params?: DynamicTableParams, enabled?: boolean) =>
  useDynamicData({ tableType: 'sales', params, enabled })

export const useAfterSalesData = (params?: DynamicTableParams, enabled?: boolean) =>
  useDynamicData({ tableType: 'after_sales', params, enabled })

export const usePerformanceData = (params?: DynamicTableParams, enabled?: boolean) =>
  useDynamicData({ tableType: 'performance', params, enabled })

export const useSalesTeamData = (params?: DynamicTableParams, enabled?: boolean) =>
  useDynamicData({ tableType: 'sales_team', params, enabled })

// Хук для работы с фильтрами из App.tsx
export const useDynamicFilters = () => {
  const [filters, setFilters] = useState({
    region: '',
    quarter: '',
    year: 2024,
    dealers: [] as string[],
    limit: 100,
    offset: 0,
    sortBy: '',
    sortOrder: 'asc' as 'asc' | 'desc'
  })

  const updateFilters = useCallback((newFilters: Partial<typeof filters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }))
  }, [])

  const resetFilters = useCallback(() => {
    setFilters({
      region: '',
      quarter: '',
      year: 2024,
      dealers: [],
      limit: 100,
      offset: 0,
      sortBy: '',
      sortOrder: 'asc'
    })
  }, [])

  const getApiParams = useCallback(() => {
    return buildDynamicParams(filters)
  }, [filters])

  return {
    filters,
    updateFilters,
    resetFilters,
    getApiParams
  }
}
