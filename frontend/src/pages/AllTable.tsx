import React, { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAllDealerData } from '../hooks/useAllDealerData'
import { AllDealerData } from '../api/all-data'

// Используем реальный интерфейс AllDealerData из API

const AllTable: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<string>('all-russia')

  // Используем хук для получения реальных данных
  const { data: dealers, loading, error, loadDataForRegion } = useAllDealerData({
    region: selectedRegion === 'all-russia' ? undefined : selectedRegion,
    quarter: 'Q1',
    year: 2024,
    autoLoad: true
  })

  const regions = [
    { id: 'all-russia', name: 'All Russia' },
    { id: 'Central', name: 'Central' },
    { id: 'North West', name: 'North West' },
    { id: 'Volga', name: 'Volga' },
    { id: 'South', name: 'South' },
    { id: 'Ural', name: 'Ural' },
    { id: 'Siberia', name: 'Siberia' },
    { id: 'Far East', name: 'Far East' }
  ]

  // Обработчик смены региона
  const handleRegionChange = async (regionId: string) => {
    setSelectedRegion(regionId)
    try {
      await loadDataForRegion(regionId)
    } catch (error) {
      console.error('Error loading data for region:', error)
    }
  }

  // Вспомогательные функции для форматирования данных
  const formatNumber = (num?: number) => {
    if (num === undefined || num === null) return 'N/A'
    return num.toLocaleString()
  }

  const formatCurrency = (num?: number) => {
    if (num === undefined || num === null) return 'N/A'
    return `₽${num.toLocaleString()}`
  }

  const formatStockData = (hdt?: number, mdt?: number, ldt?: number) => {
    if (hdt === undefined || mdt === undefined || ldt === undefined) return 'N/A'
    return `${hdt}/${mdt}/${ldt}`
  }

  const getClassColor = (dealerClass?: string) => {
    if (!dealerClass) return 'text-gray-400'
    switch (dealerClass) {
      case 'A': return 'text-green-400'
      case 'B': return 'text-yellow-400'
      case 'C': return 'text-orange-400'
      case 'D': return 'text-red-400'
      default: return 'text-gray-400'
    }
  }

  const getChecklistColor = (score?: number) => {
    if (!score) return 'text-gray-600'
    if (score >= 90) return 'text-green-600'
    if (score >= 80) return 'text-yellow-600'
    if (score >= 70) return 'text-orange-600'
    return 'text-red-600'
  }

  const getDecisionColor = (decision?: string) => {
    if (!decision) return 'text-gray-600'
    switch (decision) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs development':
      case 'Needs Development': return 'text-yellow-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-gray-600'
    }
  }

  const getJointDecision = (dealer: AllDealerData) => {
    // Logic to determine joint decision based on all factors
    const dealerClass = dealer.dealership_class
    const checklist = dealer.check_list_score
    const salesTrainings = dealer.sales_trainings === 'Yes'
    const asTrainings = dealer.as_trainings === true
    
    if (dealerClass === 'D' || (checklist && checklist < 70)) {
      return 'Close Down'
    }
    if (dealerClass === 'A' && checklist && checklist >= 90 && salesTrainings && asTrainings) {
      return 'Planned Result'
    }
    if (dealerClass === 'C' || (checklist && checklist < 80)) {
      return 'Find New Candidate'
    }
    return 'Needs Development'
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 relative">
      {/* Header with Back Button */}
      <div className="relative pt-20 pb-16">
        {/* Back to Home Button */}
        <div className="absolute left-28 top-1/2 transform -translate-y-1/2 z-50">
          <Link
            to="/"
            className="w-12 h-12 hover:w-40 bg-white bg-opacity-20 hover:bg-opacity-30 rounded-xl flex items-center justify-center transition-all duration-700 ease-out backdrop-blur-sm group overflow-hidden"
            title="Back to Home"
          >
            <svg 
              className="w-6 h-6 text-white transition-all duration-700 ease-out flex-shrink-0" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path 
                strokeLinecap="round" 
                strokeLinejoin="round" 
                strokeWidth={2} 
                d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m5-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" 
              />
            </svg>
            <span className="text-white font-medium ml-3 opacity-0 group-hover:opacity-100 transition-all duration-700 ease-out whitespace-nowrap transform translate-x-2 group-hover:translate-x-0">
              Back to Home
            </span>
          </Link>
        </div>

        {/* Title */}
        <div className="text-center">
          <h1 className="text-5xl md:text-5xl font-bold text-white mb-3">
            COMPREHENSIVE
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            DEALER ANALYSIS
          </h2>
        </div>
      </div>

      {/* Region Navigation */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mb-6">
        <div className="flex flex-wrap justify-center gap-3">
          {regions.map((region) => (
            <button
              key={region.id}
              onClick={() => handleRegionChange(region.id)}
              disabled={loading}
              className={`px-6 py-3 rounded-lg font-medium transition-all duration-200 ${
                selectedRegion === region.id
                  ? 'bg-blue-400 text-white shadow-lg'
                  : 'bg-white text-blue-900 hover:bg-blue-50'
              } ${loading ? 'opacity-50 cursor-not-allowed' : ''}`}
            >
              {region.name}
            </button>
          ))}
        </div>
      </div>

      {/* All Data Table with Horizontal Scroll */}
      <div className="w-full px-4 sm:px-6 lg:px-8 pb-8">
        <div className="p-6">
          <h3 className="text-2xl font-bold text-white mb-6 text-center">
            Complete Dealer Data Overview
          </h3>
          
          {/* Loading State */}
          {loading && (
            <div className="text-center py-12">
              <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-white mx-auto mb-4"></div>
              <p className="text-white text-xl">Loading dealer data...</p>
            </div>
          )}
          
          {/* Error State */}
          {error && (
            <div className="text-center py-12">
              <div className="text-red-400 mb-4">
                <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h2 className="text-white text-xl font-bold mb-2">Error Loading Data</h2>
              <p className="text-blue-200 mb-4">{error}</p>
              <button
                onClick={() => window.location.reload()}
                className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition-colors"
              >
                Retry
              </button>
            </div>
          )}
          
          {/* Data Table */}
          {!loading && !error && (
            <div className="overflow-x-auto">
            <div className="min-w-max">
              <table className="w-full text-sm">
                <thead>
                  <tr className="bg-blue-800 bg-opacity-50">
                    {/* Common Fields */}
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      Dealer Name
                    </th>
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      City
                    </th>
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      Sales Manager
                    </th>
                    
                                         {/* Dealer Development Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Class
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Checklist
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Brands in Portfolio
                     </th>
                     
                     {/* Sales Team Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Sales Target
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       <div>Stock</div>
                       <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       <div>Buyout</div>
                       <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Foton Salesmen
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Sales Trainings
                     </th>
                     
                     {/* Performance Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Revenue Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Profit Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Margin %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Revenue Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Profits Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Margin %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Ranking
                     </th>
                     
                     {/* After Sales Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Recommended Stock %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Warranty Stock %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Foton Labor Hours %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Service Contract
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       AS Trainings
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Foton Warranty Hours
                     </th>
                     
                     {/* Joint Decision Field */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-red-600 bg-opacity-30">
                       Joint Decision
                     </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-blue-200 divide-opacity-30">
                  {dealers.map((dealer) => (
                    <tr key={dealer.dealer_id} className="hover:bg-blue-800 hover:bg-opacity-30">
                      {/* Common Fields */}
                      <td className="px-3 py-2 text-center">
                        <Link 
                          to={`/dealer/${dealer.dealer_id}`}
                          className="text-xs font-medium text-white hover:text-blue-200 transition-colors duration-200 cursor-pointer"
                        >
                          {dealer.dealer_name_ru}
                        </Link>
                      </td>
                      <td className="px-3 py-2 text-center">
                        <div className="text-xs text-white">{dealer.city}</div>
                      </td>
                      <td className="px-3 py-2 text-center">
                        <div className="text-xs text-white">{dealer.manager}</div>
                      </td>
                      
                      {/* Dealer Development Fields */}
                      <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                        <div className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getClassColor(dealer.dealership_class)}`}>
                          {dealer.dealership_class || 'N/A'}
                        </div>
                      </td>
                      <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                        <div className={`text-xs font-medium ${getChecklistColor(dealer.check_list_score)}`}>
                          {dealer.check_list_score || 'N/A'}
                        </div>
                      </td>
                      <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                        <div className="flex justify-center">
                          <div className="w-6 h-6 bg-blue-400 bg-opacity-80 rounded-full flex items-center justify-center border border-blue-300" title="FOTON">
                            <span className="text-xs font-bold text-white">F</span>
                          </div>
                        </div>
                      </td>
                       
                      {/* Sales Team Fields */}
                      <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.sales_target_plan)}/{formatNumber(dealer.sales_target_fact)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatStockData(dealer.stock_hdt, dealer.stock_mdt, dealer.stock_ldt)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatStockData(dealer.buyout_hdt, dealer.buyout_mdt, dealer.buyout_ldt)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.foton_sales_personnel)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                        <div className={`text-xs font-medium ${dealer.sales_trainings === 'Yes' ? 'text-green-400' : 'text-white'}`}>
                          {dealer.sales_trainings || 'N/A'}
                        </div>
                      </td>
                       
                      {/* Performance Fields */}
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatCurrency(dealer.sales_revenue_rub)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatCurrency(dealer.sales_profit_rub)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.sales_margin_percent)}%</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatCurrency(dealer.after_sales_revenue_rub)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatCurrency(dealer.after_sales_profit_rub)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.after_sales_margin_pct)}%</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.foton_rank)}</div>
                      </td>
                       
                      {/* After Sales Fields */}
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.recommended_stock)}%</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.warranty_stock)}%</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.foton_labor_hours)}%</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className="text-xs text-white">{formatNumber(dealer.service_contracts)}</div>
                      </td>
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className={`text-xs font-medium ${dealer.as_trainings ? 'text-green-400' : 'text-white'}`}>
                          {dealer.as_trainings ? 'Yes' : 'No'}
                        </div>
                      </td>
                      <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                        <div className="text-xs text-white">{dealer.csi || 'N/A'}</div>
                      </td>
                      
                      {/* Joint Decision Field */}
                      <td className="px-3 py-2 text-center bg-red-600 bg-opacity-30">
                        <div className={`text-xs font-medium ${getDecisionColor(getJointDecision(dealer))}`}>
                          {getJointDecision(dealer)}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
          )}
          
        </div>
      </div>
    </div>
  )
}

export default AllTable
