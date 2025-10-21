import React, { useState, useEffect } from 'react'
import { Routes, Route, useNavigate, useLocation } from 'react-router-dom'
import SalesTeamTable from './pages/SalesTeamTable'
import SalesTable from './pages/SalesTable'
import AfterSalesTable from './pages/AfterSalesTable'
import PerformanceTable from './pages/PerformanceTable'
import AllTable from './pages/AllTable'
import QuarterComparison from './pages/QuarterComparison'
import DealerCard from './pages/DealerCard'
import BrandDemo from './pages/BrandDemo'
import Login from './pages/Login'
import ForgotPassword from './pages/ForgotPassword'
import Admin from './pages/Admin'
import ExcelUploadPage from './pages/ExcelUpload'
import ExcelTablesPage from './pages/ExcelTables'
import { REGION_MAPPING, QUARTER_MAPPING, buildDynamicParams } from './api/index'
import { useDealers } from './hooks/useDealers'
import { DealerListItem } from './api/dealers'

const App: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  
  // Обновленная логика согласно требованиям заказчика
  const [selectedRegion, setSelectedRegion] = useState<string>('')
  const [selectedDealers, setSelectedDealers] = useState<string[]>([])
  const [selectedParameters, setSelectedParameters] = useState<string>('')
  const [selectedPeriod, setSelectedPeriod] = useState<string>('')
  const [selectedYear, setSelectedYear] = useState<string>('2024')
  
  // Состояния для выпадающих меню
  const [showRegionDropdown, setShowRegionDropdown] = useState(false)
  const [showParametersDropdown, setShowParametersDropdown] = useState(false)
  const [showPeriodDropdown, setShowPeriodDropdown] = useState(false)
  const [showYearDropdown, setShowYearDropdown] = useState(false)

  // Состояние для поиска дилеров
  const [dealerSearchQuery, setDealerSearchQuery] = useState<string>('')
  
  // Состояние для модального окна выбора дилеров
  const [showDealersModal, setShowDealersModal] = useState(false)

  // Используем хук для получения дилеров
  const { dealers, loading: dealersLoading, error: dealersError, loadDealersByRegion } = useDealers({ 
    region: selectedRegion, 
    autoLoad: true 
  })

  // Читаем параметры из URL при загрузке страницы
  useEffect(() => {
    const urlParams = new URLSearchParams(location.search)
    
    // Восстанавливаем состояние формы из URL параметров
    const region = urlParams.get('region')
    const quarter = urlParams.get('quarter')
    const year = urlParams.get('year')
    const dealers = urlParams.get('dealers')
    
    if (region) {
      // Находим соответствующий ID региона
      const regionEntry = Object.entries(REGION_MAPPING).find(([_, value]) => value === region)
      if (regionEntry) {
        setSelectedRegion(regionEntry[0])
      }
    }
    
    if (quarter) {
      // Находим соответствующий ID квартала
      const quarterEntry = Object.entries(QUARTER_MAPPING).find(([_, value]) => value === quarter)
      if (quarterEntry) {
        setSelectedPeriod(quarterEntry[0])
      }
    }
    
    if (year) {
      setSelectedYear(year)
    }
    
    if (dealers) {
      setSelectedDealers(dealers.split(',').filter(id => id.trim() !== ''))
    }
  }, [location.search])

  // Правильный порядок регионов согласно требованиям
  const regions = [
    { id: 'all', name: 'All' },
    { id: 'central', name: 'Central' },
    { id: 'north-west', name: 'North-West' },
    { id: 'volga', name: 'Volga' },
    { id: 'south', name: 'South' },
    { id: 'n-caucasus', name: 'N. Caucasus' },
    { id: 'ural', name: 'Ural' },
    { id: 'siberia', name: 'Siberia' },
    { id: 'far-east', name: 'Far East' }
  ]

  // Категории параметров согласно требованиям
  const parameters = [
    { id: 'all', name: 'All Data' },
    { id: 'dealer-development', name: 'Dealer Development' },
    { id: 'sales', name: 'Sales' },
    { id: 'after-sales', name: 'After Sales' },
    { id: 'performance', name: 'Performance' }
  ]

  const quarters = [
    { id: 'q1', name: 'Q1' },
    { id: 'q2', name: 'Q2' },
    { id: 'q3', name: 'Q3' },
    { id: 'q4', name: 'Q4' }
  ]

  const years = [
    { id: '2024', name: '2024' },
    { id: '2025', name: '2025' },
    { id: '2026', name: '2026' },
    { id: '2027', name: '2027' }
  ]


  // Обработчики выбора
  const handleRegionSelect = async (regionId: string) => {
    setSelectedRegion(regionId)
    setSelectedDealers([]) // Сбрасываем выбранных дилеров при смене региона
    setDealerSearchQuery('') // Сбрасываем поиск дилеров при смене региона
    setShowRegionDropdown(false)
    
    // Загружаем дилеров для выбранного региона
    try {
      await loadDealersByRegion(regionId)
    } catch (error) {
      console.error('Error loading dealers for region:', error)
    }
  }

  const handleDealerToggle = (dealerId: string) => {
    setSelectedDealers(prev => 
      prev.includes(dealerId) 
        ? prev.filter(id => id !== dealerId)
        : [...prev, dealerId]
    )
  }

  const handleParametersSelect = (paramId: string) => {
    setSelectedParameters(paramId)
    setShowParametersDropdown(false)
  }

  const handlePeriodSelect = (quarterId: string) => {
    setSelectedPeriod(quarterId)
    setShowPeriodDropdown(false)
  }

  const handleYearSelect = (yearId: string) => {
    setSelectedYear(yearId)
    setShowYearDropdown(false)
  }

  const closeAllDropdowns = () => {
    setShowRegionDropdown(false)
    setShowParametersDropdown(false)
    setShowPeriodDropdown(false)
    setShowYearDropdown(false)
  }

  const handleFindResults = () => {
    // Проверяем обязательные параметры
    if (!selectedParameters) {
      alert('Please select parameters (Dealer Development, Sales, AfterSales, etc.)')
      return
    }

    if (!selectedPeriod) {
      alert('Please select period (Q1, Q2, Q3, Q4)')
      return
    }

    // Формируем параметры для нового динамического API
    const dynamicParams = buildDynamicParams({
      region: selectedRegion ? REGION_MAPPING[selectedRegion as keyof typeof REGION_MAPPING] || selectedRegion : 'all-russia',
      quarter: QUARTER_MAPPING[selectedPeriod as keyof typeof QUARTER_MAPPING] || selectedPeriod,
      year: parseInt(selectedYear),
      dealers: selectedDealers.length > 0 ? selectedDealers : undefined
    })

    console.log('Searching with dynamic API params:', dynamicParams)

    // Определяем маршрут на основе выбранных параметров
    const routeMapping: { [key: string]: string } = {
      'all': '/all',
      'dealer-development': '/dealer-dev',
      'sales': '/sales',
      'after-sales': '/after-sales',
      'performance': '/performance'
    }

    const targetRoute = routeMapping[selectedParameters] || '/all'
    
    // Формируем URL с параметрами
    const searchParams = new URLSearchParams()
    
    // Добавляем параметры в URL
    if (dynamicParams.region) {
      searchParams.set('region', dynamicParams.region)
    }
    if (dynamicParams.quarter) {
      searchParams.set('quarter', dynamicParams.quarter)
    }
    if (dynamicParams.year) {
      searchParams.set('year', dynamicParams.year.toString())
    }
    if (dynamicParams.dealer_ids && dynamicParams.dealer_ids.length > 0) {
      searchParams.set('dealers', dynamicParams.dealer_ids.join(','))
    }
    
    // Переходим на соответствующую страницу с параметрами в URL
    const urlWithParams = `${targetRoute}?${searchParams.toString()}`
    navigate(urlWithParams)
  }

  // Получаем доступных дилеров для выбранного региона
  const getAvailableDealers = (): DealerListItem[] => {
    if (selectedRegion === 'all' || !selectedRegion) {
      return dealers
    }
    return dealers.filter(dealer => dealer.region === REGION_MAPPING[selectedRegion as keyof typeof REGION_MAPPING])
  }

  const availableDealers = getAvailableDealers()

  const getSelectedDealersText = () => {
    if (selectedDealers.length === 0) {
      return 'Select Dealers'
    }
    if (selectedDealers.length === 1) {
      const dealer = dealers.find(d => d.id.toString() === selectedDealers[0])
      return dealer ? dealer.name : 'Select Dealers'
    }
    if (selectedDealers.length === availableDealers.length && availableDealers.length > 0) {
      return 'All Dealers'
    }
    return `${selectedDealers.length} Dealers Selected`
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700">
      {/* Routes */}
      <Routes>
        <Route path="/" element={
          <>
            {/* Header */}
            <div className="text-center pt-20 pb-16">
              <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
                FOTON DEALER
              </h1>
              <h2 className="text-3xl md:text-4xl font-bold text-blue-200 mb-8">
                DEVELOPMENT PLATFORM
              </h2>
              
            </div>

            {/* Search Panel */}
            <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 flex items-center justify-center min-h-[400px]">
              <div className="bg-white rounded-2xl shadow-5xl p-10 w-full">
                <div className="flex flex-col lg:flex-row gap-6 items-center justify-center">
                  {/* Region Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowRegionDropdown(!showRegionDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedRegion ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedRegion ? regions.find(r => r.id === selectedRegion)?.name : 'Region'}
                      </span>
                      <div className="flex items-center space-x-2">
                        {selectedRegion && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedRegion('')
                              setSelectedDealers([])
                            }}
                            className="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-50"
                            title="Clear region selection"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        )}
                        <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showRegionDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </div>
                    </button>
                    
                    {showRegionDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50 max-h-60 overflow-y-auto">
                        {regions.map((region) => (
                          <button
                            key={region.id}
                            onClick={() => handleRegionSelect(region.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{region.name}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Dealers Selection Button */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowDealersModal(true)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedDealers.length > 0 ? 'text-gray-900' : 'text-gray-500'}>
                        {getSelectedDealersText()}
                      </span>
                      <div className="flex items-center space-x-2">
                        {selectedDealers.length > 0 && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedDealers([])
                            }}
                            className="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-50"
                            title="Clear selection"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        )}
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                        </svg>
                      </div>
                    </button>
                  </div>

                  {/* Parameters Dropdown */}
                  <div className="relative flex-2 min-w-0">
                    <button 
                      onClick={() => setShowParametersDropdown(!showParametersDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedParameters ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedParameters ? parameters.find(p => p.id === selectedParameters)?.name : 'Parameters'}
                      </span>
                      <div className="flex items-center space-x-2">
                        {selectedParameters && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedParameters('')
                            }}
                            className="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-50"
                            title="Clear parameter selection"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        )}
                        <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showParametersDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </div>
                    </button>
                    
                    {showParametersDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
                        {parameters.map((param) => (
                          <button
                            key={param.id}
                            onClick={() => handleParametersSelect(param.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{param.name}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Year Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowYearDropdown(!showYearDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedYear ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedYear ? years.find(year => year.id === selectedYear)?.name : 'Year'}
                      </span>
                      <div className="flex items-center space-x-2">
                        {selectedYear !== '2024' && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedYear('2024')
                            }}
                            className="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-50"
                            title="Reset to default year"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        )}
                        <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showYearDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </div>
                    </button>
                    
                    {showYearDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
                        {years.map((year) => (
                          <button
                            key={year.id}
                            onClick={() => handleYearSelect(year.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{year.name}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Period Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowPeriodDropdown(!showPeriodDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedPeriod ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedPeriod ? quarters.find(q => q.id === selectedPeriod)?.name : 'Period'}
                      </span>
                      <div className="flex items-center space-x-2">
                        {selectedPeriod && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation()
                              setSelectedPeriod('')
                            }}
                            className="text-red-500 hover:text-red-700 p-1 rounded-full hover:bg-red-50"
                            title="Clear period selection"
                          >
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                            </svg>
                          </button>
                        )}
                        <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showPeriodDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </div>
                    </button>
                    
                    {showPeriodDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
                        {quarters.map((quarter) => (
                          <button
                            key={quarter.id}
                            onClick={() => handlePeriodSelect(quarter.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{quarter.name}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Reset All Button */}
                  <button 
                    onClick={() => {
                      setSelectedRegion('')
                      setSelectedDealers([])
                      setSelectedParameters('')
                      setSelectedPeriod('')
                      setSelectedYear('2024')
                      setDealerSearchQuery('')
                      // Закрываем модальное окно если оно открыто
                      if (showDealersModal) {
                        setShowDealersModal(false)
                      }
                    }}
                    className="bg-gray-500 hover:bg-gray-600 text-white font-bold py-4 px-6 rounded-xl transition-colors duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl"
                    title="Reset all selections"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    <span>Reset All</span>
                  </button>

                  {/* Find Results Button */}
                  <button 
                    onClick={handleFindResults}
                    className="bg-orange-500 hover:bg-orange-600 text-white font-bold py-4 px-8 rounded-xl transition-colors duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <span>Find Results</span>
                  </button>
                </div>
              </div>
            </div>
          </>
        } />
        <Route path="/dealer-dev" element={<SalesTeamTable />} />
        <Route path="/sales" element={<SalesTable />} />
        <Route path="/after-sales" element={<AfterSalesTable />} />
        <Route path="/performance" element={<PerformanceTable />} />
        <Route path="/all" element={<AllTable />} />
        <Route path="/quarter-comparison" element={<QuarterComparison />} />
        <Route path="/brand-demo" element={<BrandDemo />} />
        <Route path="/dealer/:dealerId" element={<DealerCard />} />
        <Route path="/login" element={<Login />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/admin" element={<Admin />} />
        <Route path="/excel-upload" element={<ExcelUploadPage />} />
        <Route path="/excel-tables" element={<ExcelTablesPage />} />
      </Routes>

      {/* Dealers Selection Modal */}
      {showDealersModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          {/* Backdrop */}
          <div 
            className="absolute inset-0 bg-black bg-opacity-50 backdrop-blur-sm"
            onClick={() => {
              setShowDealersModal(false)
              setDealerSearchQuery('') // Сбрасываем поиск при закрытии
            }}
          />
          
          {/* Modal Content */}
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
            {/* Header */}
            <div className="px-8 py-6 border-b border-gray-200">
              <div className="flex justify-between items-center">
                <div>
                  <h3 className="text-2xl font-bold text-gray-900">Select Dealers</h3>
                  <p className="text-gray-600 mt-1">
                    {selectedRegion ? `Available dealers in ${regions.find(r => r.id === selectedRegion)?.name} region` : 'All available dealers'}
                  </p>
                </div>
                <button
                  onClick={() => setShowDealersModal(false)}
                  className="text-gray-400 hover:text-gray-600 transition-colors duration-200 p-2 hover:bg-gray-100 rounded-lg"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            {/* Search and Controls */}
            <div className="px-8 py-4 border-b border-gray-200 bg-gray-50">
              <div className="flex items-center space-x-4">
                <div className="flex-1">
                  <input
                    type="text"
                    placeholder="Search dealers by name..."
                    value={dealerSearchQuery}
                    onChange={(e) => setDealerSearchQuery(e.target.value)}
                    className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:border-blue-500 focus:outline-none"
                  />
                </div>
                <button
                  onClick={() => {
                    if (selectedDealers.length === availableDealers.length) {
                      setSelectedDealers([])
                    } else {
                      setSelectedDealers(availableDealers.map(d => d.id.toString()))
                    }
                  }}
                  className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
                >
                  {selectedDealers.length === availableDealers.length ? 'Deselect All' : 'Select All'}
                </button>
                <button
                  onClick={() => setSelectedDealers([])}
                  className="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors duration-200"
                >
                  Clear Selection
                </button>
                <button
                  onClick={() => setDealerSearchQuery('')}
                  className="px-4 py-2 bg-gray-400 text-white rounded-lg hover:bg-gray-500 transition-colors duration-200"
                  disabled={!dealerSearchQuery.trim()}
                >
                  Clear Search
                </button>
              </div>
            </div>

            {/* Dealers List */}
            <div className="flex-1 overflow-y-auto max-h-96">
              {dealersLoading ? (
                <div className="px-8 py-12 text-center text-gray-500">
                  <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
                  <p className="text-lg font-medium">Loading dealers...</p>
                </div>
              ) : dealersError ? (
                <div className="px-8 py-12 text-center text-gray-500">
                  <svg className="w-16 h-16 mx-auto text-red-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <p className="text-lg font-medium text-red-600">Error loading dealers</p>
                  <p className="text-sm text-red-500">{dealersError}</p>
                </div>
              ) : availableDealers.length > 0 ? (
                <div className="divide-y divide-gray-200">
                  {availableDealers
                    .filter(dealer => 
                      !dealerSearchQuery.trim() || 
                      dealer.name.toLowerCase().includes(dealerSearchQuery.toLowerCase())
                    )
                    .map((dealer) => (
                      <div
                        key={dealer.id}
                        className={`px-8 py-4 hover:bg-gray-50 transition-colors duration-200 cursor-pointer ${
                          selectedDealers.includes(dealer.id.toString()) ? 'bg-blue-50 border-l-4 border-l-blue-500' : ''
                        }`}
                        onClick={() => handleDealerToggle(dealer.id.toString())}
                      >
                        <div className="flex items-center justify-between">
                          <div className="flex-1">
                            <div className="font-medium text-gray-900 text-lg">{dealer.name}</div>
                            <div className="text-sm text-gray-500 mt-1">
                              {dealer.city} • {dealer.region} • {dealer.manager}
                            </div>
                          </div>
                          <div className="flex items-center space-x-3">
                            {selectedDealers.includes(dealer.id.toString()) && (
                              <div className="w-6 h-6 bg-blue-600 rounded-full flex items-center justify-center">
                                <svg className="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                                </svg>
                              </div>
                            )}
                            <div className={`w-6 h-6 border-2 rounded-full ${
                              selectedDealers.includes(dealer.id.toString()) 
                                ? 'border-blue-600 bg-blue-600' 
                                : 'border-gray-300'
                            }`}>
                              {selectedDealers.includes(dealer.id.toString()) && (
                                <div className="w-2 h-2 bg-white rounded-full m-auto mt-1"></div>
                              )}
                            </div>
                          </div>
                        </div>
                      </div>
                    ))}
                </div>
              ) : (
                <div className="px-8 py-12 text-center text-gray-500">
                  <svg className="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                  <p className="text-lg font-medium">No dealers found</p>
                  <p className="text-sm">Try selecting a different region or clearing your search</p>
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="px-8 py-6 border-t border-gray-200 bg-gray-50">
              <div className="flex justify-between items-center">
                <div className="text-sm text-gray-600">
                  {selectedDealers.length > 0 ? (
                    <span>
                      <span className="font-medium text-gray-900">{selectedDealers.length}</span> dealer{selectedDealers.length !== 1 ? 's' : ''} selected
                    </span>
                  ) : (
                    <span>No dealers selected</span>
                  )}
                </div>
                <div className="flex space-x-3">
                  <button
                    onClick={() => {
                      setShowDealersModal(false)
                      setDealerSearchQuery('') // Сбрасываем поиск при отмене
                    }}
                    className="px-6 py-2 text-gray-600 hover:text-gray-800 transition-colors duration-200"
                  >
                    Cancel
                  </button>
                  <button
                    onClick={() => {
                      setShowDealersModal(false)
                      setDealerSearchQuery('') // Сбрасываем поиск при применении
                    }}
                    className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
                  >
                    Apply Selection
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Close dropdowns when clicking outside */}
      {(showRegionDropdown || showParametersDropdown || showPeriodDropdown || showYearDropdown) && (
        <div 
          className="fixed inset-0 z-40" 
          onClick={closeAllDropdowns}
        />
      )}
    </div>
  )
}

export default App
