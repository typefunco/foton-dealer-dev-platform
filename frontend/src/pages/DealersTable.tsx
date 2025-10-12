import React, { useState, useEffect } from 'react'
import { Link, useLocation } from 'react-router-dom'
import { useDealerDevData } from '../hooks/useDynamicData'
import type { Dealer } from '../api/dealers'

const DealersTable: React.FC = () => {
  const location = useLocation()
  const [selectedRegion, setSelectedRegion] = useState<string>('all-russia')
  const [sortConfig, setSortConfig] = useState<{
    key: keyof Dealer | null
    direction: 'asc' | 'desc' | null
  }>({ key: null, direction: null })

  // Получаем параметры из навигации
  const navigationFilters = location.state?.filters || {}

  const { data: dealers, loading, error, updateParams } = useDealerDevData({
    region: navigationFilters.region || (selectedRegion === 'all-russia' ? undefined : selectedRegion),
    quarter: navigationFilters.quarter,
    year: navigationFilters.year
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
      
      if (aValue < bValue) {
        return sortConfig.direction === 'asc' ? -1 : 1
      }
      if (aValue > bValue) {
        return sortConfig.direction === 'asc' ? 1 : -1
      }
      return 0
    })
  }

  const regions = [
    { id: 'all-russia', name: 'All Russia' },
    { id: 'Central', name: 'Central' },
    { id: 'Ural', name: 'Ural' },
    { id: 'Volga', name: 'Volga' },
    { id: 'Kavkaz', name: 'Kavkaz' },
    { id: 'North West', name: 'North West' },
    { id: 'Siberia', name: 'Siberia' },
    { id: 'South', name: 'South' }
  ]

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
             DEALER DEVELOPMENT
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            ANALYSIS
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
        {loading && (
          <div className="flex justify-center items-center py-8">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white"></div>
            <span className="ml-4 text-white text-lg">Loading dealers...</span>
          </div>
        )}
        
        {error && (
          <div className="bg-red-500 bg-opacity-20 border border-red-500 rounded-lg p-4 mb-6">
            <div className="flex items-center">
              <svg className="w-5 h-5 text-red-400 mr-2" fill="currentColor" viewBox="0 0 20 20">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd" />
              </svg>
              <span className="text-red-200">Error: {error}</span>
            </div>
          </div>
        )}
        
        {!loading && !error && (
          <div className="overflow-x-auto">
            <table className="w-full">
            <thead>
              <tr>
                <th 
                  className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                  onClick={() => handleSort('dealer_name_ru')}
                >
                  <div className="flex items-center justify-center space-x-1">
                    <span>Dealer Name</span>
                    {getSortIcon('dealer_name_ru')}
                  </div>
                </th>
                <th 
                  className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                  onClick={() => handleSort('city')}
                >
                  <div className="flex items-center justify-center space-x-1">
                    <span>City</span>
                    {getSortIcon('city')}
                  </div>
                </th>
                <th 
                  className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                  onClick={() => handleSort('region')}
                >
                  <div className="flex items-center justify-center space-x-1">
                    <span>Region</span>
                    {getSortIcon('region')}
                  </div>
                </th>
                <th 
                  className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                  onClick={() => handleSort('manager')}
                >
                  <div className="flex items-center justify-center space-x-1">
                    <span>Manager</span>
                    {getSortIcon('manager')}
                  </div>
                </th>
              </tr>
            </thead>
            <tbody className="divide-y divide-blue-200 divide-opacity-30">
              {getSortedDealers().map((dealer) => (
                <tr key={dealer.dealer_id} className="hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200">
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <Link 
                      to={`/dealer/${dealer.dealer_id}`}
                      className="text-sm font-medium text-white hover:text-blue-200 transition-colors duration-200 cursor-pointer"
                    >
                      {dealer.dealer_name_ru}
                    </Link>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="text-sm text-white">{dealer.city}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="text-sm text-white">{dealer.region}</div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="text-sm text-white">{dealer.manager}</div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
          </div>
        )}
      </div>
    </div>
  )
}

export default DealersTable
