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
      const result = await getDynamicData<T>(tableType, currentParams)
      setResponse(result)
      setData(result.data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
      setError(errorMessage)
      console.error(`Error fetching ${tableType} data:`, err)
    } finally {
      setLoading(false)
    }
  }, [tableType, currentParams, enabled])

  const refetch = useCallback(() => {
    fetchData()
  }, [fetchData])

  const updateParams = useCallback((newParams: DynamicTableParams) => {
    setCurrentParams(prev => ({ ...prev, ...newParams }))
  }, [])

  useEffect(() => {
    fetchData()
  }, [fetchData])

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
