import React, { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useDealerDevData } from '../hooks/useDynamicData'

interface Dealer {
  id: string
  name: string
  city: string
  class: string
  checklist: number
  brandsInPortfolio: string[]
  brandsCount: number
  branding: boolean
  buySideBusiness: string[]
  dealerDevRecommendation: string
}

const SalesTeamTable: React.FC = () => {
  const location = useLocation()
  const [selectedRegion, setSelectedRegion] = useState<string>('Central')
  const [sortConfig, setSortConfig] = useState<{
    key: keyof Dealer | null
    direction: 'asc' | 'desc' | null
  }>({ key: null, direction: null })

  // Получаем параметры из навигации
  const navigationFilters = location.state?.filters || {}

  const { data: dealers, loading, error, updateParams } = useDealerDevData({
    region: navigationFilters.region || (selectedRegion === 'all-russia' ? undefined : selectedRegion) || 'Central',
    quarter: navigationFilters.quarter || 'Q1',
    year: navigationFilters.year || 2024
  })


  // Обработка изменения региона
  useEffect(() => {
    updateParams({ region: selectedRegion === 'all-russia' ? undefined : selectedRegion })
  }, [selectedRegion, updateParams])

  // Применяем параметры из навигации при загрузке
  useEffect(() => {
    if (navigationFilters.region) {
      setSelectedRegion(navigationFilters.region)
    }
    if (navigationFilters.quarter || navigationFilters.year) {
      updateParams({
        quarter: navigationFilters.quarter,
        year: navigationFilters.year
      })
    }
  }, [navigationFilters, updateParams])

  const handleSort = (key: keyof Dealer) => {
    let direction: 'asc' | 'desc' | null = 'asc'
    
    if (sortConfig.key === key) {
      if (sortConfig.direction === 'asc') {
        direction = 'desc'
      } else if (sortConfig.direction === 'desc') {
        direction = null
      }
    }
    
    setSortConfig({ key, direction })
  }

  const getSortedDealers = () => {
    if (!dealers || !sortConfig.key || !sortConfig.direction) {
      return dealers || []
    }

    return [...dealers].sort((a, b) => {
      const aValue = a[sortConfig.key!]
      const bValue = b[sortConfig.key!]
      
      if (sortConfig.key === 'branding') {
        const aBool = aValue as boolean
        const bBool = bValue as boolean
        
        if (sortConfig.direction === 'asc') {
          return aBool === bBool ? 0 : aBool ? -1 : 1 // true first
        } else {
          return aBool === bBool ? 0 : aBool ? 1 : -1 // false first
        }
      }
      
      if (sortConfig.key === 'dealerDevRecommendation') {
        const decisionOrder = { 
          'Planned Result': 4, 
          'Needs Development': 3, 
          'Find New Candidate': 2, 
          'Close Down': 1 
        }
        const aOrder = decisionOrder[aValue as keyof typeof decisionOrder]
        const bOrder = decisionOrder[bValue as keyof typeof decisionOrder]
        
        if (sortConfig.direction === 'asc') {
          return bOrder - aOrder // Planned Result first
        } else {
          return aOrder - bOrder // Close Down first
        }
      }
      
      if (aValue < bValue) {
        return sortConfig.direction === 'asc' ? -1 : 1
      }
      if (aValue > bValue) {
        return sortConfig.direction === 'asc' ? 1 : -1
      }
      return 0
    })
  }

  const getSortIcon = (key: keyof Dealer) => {
    if (sortConfig.key !== key) {
      return (
        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
        </svg>
      )
    }
    
    if (sortConfig.direction === 'asc') {
      return (
        <svg className="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
        </svg>
      )
    }
    
    if (sortConfig.direction === 'desc') {
      return (
        <svg className="w-4 h-4 text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      )
    }
    
    return (
      <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
      </svg>
    )
  }

  const regions = [
    { id: 'all-russia', name: 'All Russia' },
    { id: 'Central', name: 'Central' },
    { id: 'North West', name: 'North West' },
    { id: 'Volga', name: 'Volga' },
    { id: 'South', name: 'South' },
    { id: 'Kavkaz', name: 'Kavkaz' },
    { id: 'Ural', name: 'Ural' },
    { id: 'Siberia', name: 'Siberia' },
    { id: 'Far East', name: 'Far East' }
  ]

  // Данные теперь получаются из API через хук useSalesTeamData

  const getSalesDecisionColor = (decision: string) => {
    switch (decision) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs Development': return 'text-green-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-green-600'
    }
  }

  // Обработка состояний загрузки и ошибок
  if (loading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-16 w-16 border-b-2 border-white mx-auto mb-4"></div>
          <p className="text-white text-xl">Loading dealer development data...</p>
        </div>
      </div>
    )
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-center bg-white bg-opacity-10 backdrop-blur-sm rounded-lg p-8 max-w-md">
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
      </div>
    )
  }

  // Проверка на пустые данные
  if (!dealers || dealers.length === 0) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-center bg-white bg-opacity-10 backdrop-blur-sm rounded-lg p-8 max-w-md">
          <div className="text-blue-400 mb-4">
            <svg className="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
          </div>
          <h2 className="text-white text-xl font-bold mb-2">No Data Available</h2>
          <p className="text-blue-200 mb-4">No dealer development data found for the selected criteria.</p>
          <button
            onClick={() => window.location.reload()}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition-colors"
          >
            Refresh Page
          </button>
        </div>
      </div>
    )
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
             DEALER
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            DEVELOPMENT
          </h2>
        </div>
      </div>

      {/* Region Navigation */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mb-6">
        <div className="flex flex-wrap justify-center gap-3">
          {regions.map((region) => (
            <button
              key={region.id}
              onClick={() => setSelectedRegion(region.id)}
              className={`px-6 py-3 rounded-lg font-medium transition-all duration-200 ${
                selectedRegion === region.id
                  ? 'bg-blue-400 text-white shadow-lg'
                  : 'bg-white text-blue-900 hover:bg-blue-50'
              }`}
            >
              {region.name}
            </button>
          ))}
        </div>
      </div>

      {/* Dealers Table */}
      <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
        <table className="w-full">
          <thead>
            <tr>
              <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                Dealer Name
              </th>
              <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                City
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('class')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Class</span>
                  {getSortIcon('class')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('checklist')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Checklist Score</span>
                  {getSortIcon('checklist')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('brandsCount')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Brands Count</span>
                  {getSortIcon('brandsCount')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('branding')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Branding</span>
                  {getSortIcon('branding')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('dealerDevRecommendation')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Recommendation</span>
                  {getSortIcon('dealerDevRecommendation')}
                </div>
              </th>
            </tr>
          </thead>
          <tbody className="divide-y divide-blue-200 divide-opacity-30">
            {getSortedDealers().map((dealer) => (
              <tr key={dealer.id} className="hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200">
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <Link 
                    to={`/dealer/${dealer.id}`}
                    className="text-sm font-medium text-white hover:text-blue-200 transition-colors duration-200 cursor-pointer"
                  >
                    {dealer.name}
                  </Link>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.city}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.class}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.checklist}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.brandsCount}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className={`text-sm font-medium ${
                    dealer.branding ? 'text-green-600' : 'text-white'
                  }`}>
                    {dealer.branding ? 'Yes' : 'No'}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className={`text-sm font-medium ${getSalesDecisionColor(dealer.dealerDevRecommendation)}`}>
                    {dealer.dealerDevRecommendation}
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}

export default SalesTeamTable
