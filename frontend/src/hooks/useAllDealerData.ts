import { useState, useEffect, useCallback } from 'react'
import { getAllDealerData, AllDealerData, AllDataParams } from '../api/all-data'

interface UseAllDealerDataOptions {
  region?: string
  quarter?: string
  year?: number
  autoLoad?: boolean
}

interface UseAllDealerDataReturn {
  data: AllDealerData[]
  loading: boolean
  error: string | null
  refetch: () => Promise<void>
  loadDataForRegion: (region: string) => Promise<void>
  clearError: () => void
}

/**
 * Хук для работы со всеми данными дилеров
 */
export function useAllDealerData(options: UseAllDealerDataOptions = {}): UseAllDealerDataReturn {
  const { region, quarter, year, autoLoad = true } = options
  
  const [data, setData] = useState<AllDealerData[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const loadData = useCallback(async (params: AllDataParams = {}) => {
    try {
      setLoading(true)
      setError(null)
      
      const dealerData = await getAllDealerData(params)
      setData(dealerData)
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Unknown error occurred'
      setError(errorMessage)
      console.error('Error loading all dealer data:', err)
    } finally {
      setLoading(false)
    }
  }, [])

  const loadDataForRegion = useCallback(async (regionName: string) => {
    const params: AllDataParams = {
      region: regionName === 'all-russia' ? undefined : regionName,
      quarter: quarter || 'Q1',
      year: year || 2024
    }
    await loadData(params)
  }, [loadData, quarter, year])

  const refetch = useCallback(async () => {
    const params: AllDataParams = {
      region: region === 'all-russia' ? undefined : region,
      quarter: quarter || 'Q1',
      year: year || 2024
    }
    await loadData(params)
  }, [loadData, region, quarter, year])

  const clearError = useCallback(() => {
    setError(null)
  }, [])

  // Автоматическая загрузка при монтировании компонента
  useEffect(() => {
    if (autoLoad) {
      const params: AllDataParams = {
        region: region === 'all-russia' ? undefined : region,
        quarter: quarter || 'Q1',
        year: year || 2024
      }
      loadData(params)
    }
  }, [autoLoad, region, quarter, year, loadData])

  return {
    data,
    loading,
    error,
    refetch,
    loadDataForRegion,
    clearError
  }
}

/**
 * Хук для получения всех данных дилеров для региона Central
 */
export function useAllDealerDataForCentral(): UseAllDealerDataReturn {
  return useAllDealerData({ region: 'Central', quarter: 'Q1', year: 2024, autoLoad: true })
}
