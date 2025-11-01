import React, { useState, useEffect } from 'react'
import { Link, useLocation, useSearchParams } from 'react-router-dom'
import { useDealerDevData } from '../hooks/useDynamicData'
import { getBusinessIcons, getBrandIcons } from '../utils'

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
  const [searchParams] = useSearchParams()
  const [selectedRegion, setSelectedRegion] = useState<string>('Central')
  const [sortConfig, setSortConfig] = useState<{
    key: keyof Dealer | null
    direction: 'asc' | 'desc' | null
  }>({ key: null, direction: null })

  // Получаем параметры из URL
  const regionFromUrl = searchParams.get('region') || 'Central'
  const quarterFromUrl = searchParams.get('quarter') || ''
  const yearFromUrl = parseInt(searchParams.get('year') || '0')

  // Получаем параметры из навигации (если есть)
  const navigationFilters = location.state?.filters || {}

  const { data: dealers, loading, error, updateParams } = useDealerDevData({
    region: regionFromUrl === 'all-russia' ? undefined : regionFromUrl,
    quarter: quarterFromUrl || navigationFilters.quarter,
    year: yearFromUrl || navigationFilters.year
  })

  // Отладка: проверяем данные, которые приходят с бэкенда
  useEffect(() => {
    if (dealers && dealers.length > 0) {
      console.log('=== DEBUG: Dealers Data ===')
      console.log('Total dealers:', dealers.length)
      console.log('First dealer sample:', dealers[0])
      console.log('First dealer buySideBusiness:', dealers[0]?.buySideBusiness)
      console.log('Type of buySideBusiness:', typeof dealers[0]?.buySideBusiness)
      console.log('Is array:', Array.isArray(dealers[0]?.buySideBusiness))
      
      if (dealers[0]?.buySideBusiness && dealers[0].buySideBusiness.length > 0) {
        console.log('Business icons result:', getBusinessIcons(dealers[0].buySideBusiness))
      } else {
        console.warn('buySideBusiness is empty or undefined for first dealer')
      }
      
      // Проверяем, есть ли хотя бы один дилер с бизнесами
      const dealersWithBusinesses = dealers.filter(d => d.buySideBusiness && d.buySideBusiness.length > 0)
      console.log('Dealers with businesses:', dealersWithBusinesses.length)
      if (dealersWithBusinesses.length > 0) {
        console.log('Sample dealer with businesses:', dealersWithBusinesses[0])
      }
      console.log('=== END DEBUG ===')
    }
  }, [dealers])


  // Обработка изменения региона
  useEffect(() => {
    updateParams({ 
      region: selectedRegion === 'all-russia' ? undefined : selectedRegion,
      quarter: quarterFromUrl || navigationFilters.quarter,
      year: yearFromUrl || navigationFilters.year
    })
  }, [selectedRegion, updateParams, quarterFromUrl, yearFromUrl, navigationFilters])

  // Инициализируем регион из URL при загрузке
  useEffect(() => {
    if (regionFromUrl && regionFromUrl !== selectedRegion) {
      setSelectedRegion(regionFromUrl)
    }
  }, [regionFromUrl])

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
      
      if (sortConfig.key === 'brandsInPortfolio') {
        const aCount = (a.brandsInPortfolio?.length || 0)
        const bCount = (b.brandsInPortfolio?.length || 0)
        
        if (sortConfig.direction === 'asc') {
          return bCount - aCount // Больше к меньшему
        } else {
          return aCount - bCount // Меньше к большему
        }
      }
      
      if (sortConfig.key === 'buySideBusiness') {
        const aCount = (a.buySideBusiness?.length || 0)
        const bCount = (b.buySideBusiness?.length || 0)
        
        if (sortConfig.direction === 'asc') {
          return bCount - aCount // Больше к меньшему
        } else {
          return aCount - bCount // Меньше к большему
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
                onClick={() => handleSort('brandsInPortfolio')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Brands Portfolio</span>
                  {getSortIcon('brandsInPortfolio')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('buySideBusiness')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Business Portfolio</span>
                  {getSortIcon('buySideBusiness')}
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
                  <div className="flex items-center justify-center gap-1.5 flex-wrap">
                    {dealer.brandsInPortfolio && dealer.brandsInPortfolio.length > 0 ? (
                      getBrandIcons(dealer.brandsInPortfolio).map((brand, index) => (
                        <div
                          key={index}
                          className="relative group"
                          title={brand.name}
                        >
                          {brand.iconPath ? (
                            <img
                              src={brand.iconPath}
                              alt={brand.name}
                              className="w-[40px] h-[40px] object-contain rounded-lg hover:scale-125 transition-transform duration-200 bg-white bg-opacity-10 p-0.5"
                              onError={(e) => {
                                const target = e.target as HTMLImageElement
                                target.style.display = 'none'
                                if (target.nextSibling) {
                                  (target.nextSibling as HTMLElement).style.display = 'inline-flex'
                                }
                              }}
                            />
                          ) : null}
                          {!brand.iconPath && (
                            <div className="w-[40px] h-[40px] bg-blue-500 bg-opacity-50 rounded-lg flex items-center justify-center text-white text-xs font-semibold">
                              {brand.name.charAt(0).toUpperCase()}
                            </div>
                          )}
                          <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                            {brand.name}
                          </div>
                        </div>
                      ))
                    ) : (
                      <span className="text-gray-400 text-sm">-</span>
                    )}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="flex items-center justify-center gap-1.5 flex-wrap">
                    {dealer.buySideBusiness && dealer.buySideBusiness.length > 0 ? (
                      getBusinessIcons(dealer.buySideBusiness).map((business, index) => {
                        // Отладка для первого бизнеса первого дилера
                        if (index === 0 && dealers && dealers[0]?.id === dealer.id) {
                          console.log('Rendering business icon:', business)
                        }
                        return (
                        <div
                          key={index}
                          className="relative group"
                          title={business.name}
                        >
                          {business.iconPath ? (
                            <img
                              src={business.iconPath}
                              alt={business.name}
                              className="w-[40px] h-[40px] object-contain rounded-lg hover:scale-125 transition-transform duration-200 bg-white bg-opacity-10 p-0.5"
                              onError={(e) => {
                                // Fallback если иконка не загрузилась
                                const target = e.target as HTMLImageElement
                                target.style.display = 'none'
                                if (target.nextSibling) {
                                  (target.nextSibling as HTMLElement).style.display = 'inline-flex'
                                }
                              }}
                            />
                          ) : null}
                          {/* Fallback если иконка не найдена */}
                          {!business.iconPath && (
                            <div className="w-[40px] h-[40px] bg-blue-500 bg-opacity-50 rounded-lg flex items-center justify-center text-white text-xs font-semibold">
                              {business.name.charAt(0).toUpperCase()}
                            </div>
                          )}
                          {/* Tooltip с названием бизнеса */}
                          <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-2 py-1 bg-gray-900 text-white text-xs rounded opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                            {business.name}
                          </div>
                        </div>
                        )
                      })
                    ) : (
                      <span className="text-gray-400 text-sm" title={`Debug: buySideBusiness = ${JSON.stringify(dealer.buySideBusiness)}`}>-</span>
                    )}
                  </div>
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
