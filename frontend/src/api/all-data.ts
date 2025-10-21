// API клиент для получения всех данных дилеров

export interface AllDealerData {
  dealer_id: number
  dealer_name_ru: string
  city: string
  region: string
  manager: string
  period: string
  
  // Dealer Development данные
  check_list_score?: number
  dealership_class?: string
  branding?: boolean
  marketing_investments?: number
  dd_recommendation?: string
  
  // Sales данные
  stock_hdt?: number
  stock_mdt?: number
  stock_ldt?: number
  buyout_hdt?: number
  buyout_mdt?: number
  buyout_ldt?: number
  foton_sales_personnel?: number
  sales_target_plan?: number
  sales_target_fact?: number
  service_contracts_sales?: number
  sales_trainings?: string
  sales_recommendation?: string
  
  // AfterSales данные
  recommended_stock?: number
  warranty_stock?: number
  foton_labor_hours?: number
  foton_warranty_hours?: number
  service_contracts?: number
  as_trainings?: boolean
  csi?: string
  as_decision?: string
  
  // Performance данные
  sales_revenue_rub?: number
  sales_profit_rub?: number
  sales_margin_percent?: number
  after_sales_revenue_rub?: number
  after_sales_profit_rub?: number
  after_sales_margin_pct?: number
  marketing_investment?: number
  foton_rank?: number
  performance_decision?: string
}

export interface AllDataParams {
  region?: string
  quarter?: string
  year?: number
}

const API_BASE_URL = 'http://localhost:8080/api'

/**
 * Получает все данные дилеров (комплексные данные из всех таблиц)
 */
export async function getAllDealerData(params: AllDataParams = {}): Promise<AllDealerData[]> {
  const searchParams = new URLSearchParams()
  
  if (params.region) {
    searchParams.set('region', params.region)
  }
  if (params.quarter) {
    searchParams.set('quarter', params.quarter)
  }
  if (params.year) {
    searchParams.set('year', params.year.toString())
  }

  const url = `${API_BASE_URL}/all-data?${searchParams.toString()}`
  
  try {
    const response = await fetch(url)
    
    if (!response.ok) {
      throw new Error(`Failed to fetch all dealer data: ${response.status} ${response.statusText}`)
    }
    
    const data = await response.json()
    return data
  } catch (error) {
    console.error('Error fetching all dealer data:', error)
    throw error
  }
}

/**
 * Получает все данные дилеров для региона Central по умолчанию
 */
export async function getAllDealerDataForCentral(): Promise<AllDealerData[]> {
  return getAllDealerData({ region: 'Central', quarter: 'Q1', year: 2024 })
}
