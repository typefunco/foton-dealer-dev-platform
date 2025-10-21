// API клиент для работы с дилерами

export interface DealerListItem {
  id: number
  name: string
  region: string
  city: string
  manager: string
}

export interface DealerListParams {
  region?: string
  limit?: number
  offset?: number
}

const API_BASE_URL = import.meta.env.DEV ? '/api' : 'http://localhost:8080/api'

/**
 * Получает список дилеров для выпадающих меню
 */
export async function getDealersList(params: DealerListParams = {}): Promise<DealerListItem[]> {
  const searchParams = new URLSearchParams()
  
  if (params.region) {
    searchParams.set('region', params.region)
  }
  if (params.limit) {
    searchParams.set('limit', params.limit.toString())
  }
  if (params.offset) {
    searchParams.set('offset', params.offset.toString())
  }

  const url = `${API_BASE_URL}/dealers/list?${searchParams.toString()}`
  
  try {
    const response = await fetch(url, {
      credentials: 'include', // Включаем cookies для аутентификации
    })
    
    if (!response.ok) {
      throw new Error(`Failed to fetch dealers list: ${response.status} ${response.statusText}`)
    }
    
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Error fetching dealers list:', error)
    throw error
  }
}

/**
 * Получает всех дилеров (для случаев когда нужно загрузить всех сразу)
 */
export async function getAllDealers(): Promise<DealerListItem[]> {
  return getDealersList({ limit: 100 })
}

/**
 * Получает дилеров по региону
 */
export async function getDealersByRegion(region: string): Promise<DealerListItem[]> {
  return getDealersList({ region, limit: 100 })
}