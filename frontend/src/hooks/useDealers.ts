import { useState, useEffect, useCallback } from 'react'
import { getDealersList, DealerListItem, DealerListParams } from '../api/dealers'

interface UseDealersOptions {
  region?: string
  autoLoad?: boolean
}

interface UseDealersReturn {
  dealers: DealerListItem[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  loadDealersByRegion: (region: string) => Promise<void>
  clearError: () => void
}

/**
 * Хук для работы со списком дилеров
 */
export function useDealers(options: UseDealersOptions = {}): UseDealersReturn {
  const { region, autoLoad = true } = options
  
  const [dealers, setDealers] = useState<DealerListItem[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadDealers = useCallback(async (params: DealerListParams = {}) => {
    try {
      setLoading(true)
      setError(null)
      
      const data = await getDealersList(params)
      setDealers(data)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
      setError(errorMessage)
      console.error('Error loading dealers:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  const loadDealersByRegion = useCallback(async (regionName: string) => {
    if (regionName === 'all' || !regionName) {
      await loadDealers()
    } else {
      await loadDealers({ region: regionName })
    }
  }, [loadDealers])

  const refetch = useCallback(async () => {
    if (region) {
      await loadDealersByRegion(region)
    } else {
      await loadDealers()
    }
  }, [region, loadDealersByRegion, loadDealers])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  // Автоматическая загрузка при монтировании компонента
  useEffect(() => {
    if (autoLoad) {
      if (region) {
        loadDealersByRegion(region)
      } else {
        loadDealers()
      }
    }
  }, [autoLoad, region, loadDealersByRegion, loadDealers])

  return {
    dealers,
    loading,
    error,
    refetch,
    loadDealersByRegion,
    clearError
  }
}

/**
 * Хук для получения всех дилеров (упрощенная версия)
 */
export function useAllDealers(): UseDealersReturn {
  return useDealers({ autoLoad: true })
}

/**
 * Хук для получения дилеров по региону
 */
export function useDealersByRegion(region: string): UseDealersReturn {
  return useDealers({ region, autoLoad: true })
}