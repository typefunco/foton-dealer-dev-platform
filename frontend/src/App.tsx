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
import { AuthGuard, AdminGuard } from './components/AuthGuard'
import { AuthProvider, useAuth } from './contexts/AuthContext'
import { REGION_MAPPING, buildDynamicParams } from './api/index'
import { useDealers } from './hooks/useDealers'
import { DealerListItem } from './api/dealers'

const AppContent: React.FC = () => {
  const navigate = useNavigate()
  const location = useLocation()
  const { user, logout, isAuthenticated } = useAuth()
  
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
  const { dealers, loading: dealersLoading, error: dealersError } = useDealers({ 
    region: selectedRegion, 
    autoLoad: true 
  })

  // Читаем параметры из URL при загрузке страницы
  useEffect(() => {
    const urlParams = new URLSearchParams(location.search)
    const region = urlParams.get('region')
    const parameters = urlParams.get('parameters')
    const period = urlParams.get('period')
    const year = urlParams.get('year')
    const dealers = urlParams.get('dealers')

    if (region) setSelectedRegion(region)
    if (parameters) setSelectedParameters(parameters)
    if (period) setSelectedPeriod(period)
    if (year) setSelectedYear(year)
    if (dealers) {
      setSelectedDealers(dealers.split(',').filter(id => id.trim() !== ''))
    }
  }, [location.search])

  // Регионы
  const regions = [
    { id: 'Central', name: 'Central' },
    { id: 'North West', name: 'North West' },
    { id: 'Volga', name: 'Volga' },
    { id: 'South', name: 'South' },
    { id: 'Ural', name: 'Ural' },
    { id: 'Siberia', name: 'Siberia' },
    { id: 'Far East', name: 'Far East' },
    { id: 'all', name: 'All Russia' },
  ]

  // Параметры
  const parameters = [
    { id: 'dealer-dev', name: 'Dealer Development' },
    { id: 'sales', name: 'Sales Team' },
    { id: 'after-sales', name: 'After Sales' },
    { id: 'performance', name: 'Performance' },
    { id: 'all', name: 'All Data' },
  ]

  // Периоды
  const periods = [
    { id: 'Q1', name: 'Q1' },
    { id: 'Q2', name: 'Q2' },
    { id: 'Q3', name: 'Q3' },
    { id: 'Q4', name: 'Q4' },
  ]

  // Годы
  const years = [
    { id: '2024', name: '2024' },
    { id: '2025', name: '2025' },
  ]

  // Маппинг параметров на роуты
  const routeMapping: Record<string, string> = {
    'dealer-dev': '/dealer-dev',
    'sales': '/sales',
    'after-sales': '/after-sales',
    'performance': '/performance',
    'all': '/all',
  }

  // Функция для закрытия всех выпадающих меню
  const closeAllDropdowns = () => {
    setShowRegionDropdown(false)
    setShowParametersDropdown(false)
    setShowPeriodDropdown(false)
    setShowYearDropdown(false)
  }

  // Функция для обработки поиска
  const handleSearch = () => {
    console.log('handleSearch called with:', {
      region: selectedRegion,
      quarter: selectedPeriod,
      year: selectedYear,
      dealers: selectedDealers,
      parameters: selectedParameters
    })

    const searchParams = buildDynamicParams({
      region: selectedRegion,
      quarter: selectedPeriod,
      year: parseInt(selectedYear),
      dealers: selectedDealers.length > 0 ? selectedDealers : undefined
    })

    console.log('buildDynamicParams result:', searchParams)

    // Преобразуем объект в строку URL параметров
    const urlParams = new URLSearchParams()
    Object.entries(searchParams).forEach(([key, value]) => {
      if (value !== undefined && value !== null) {
        if (Array.isArray(value)) {
          value.forEach(item => urlParams.append(key, item.toString()))
        } else {
          urlParams.append(key, value.toString())
        }
      }
    })

    console.log('URLSearchParams result:', urlParams.toString())

    const targetRoute = routeMapping[selectedParameters] || '/all'
    const urlWithParams = `${targetRoute}?${urlParams.toString()}`
    
    console.log('Final URL:', urlWithParams)
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
    <>
      {/* Routes */}
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        
        {/* Protected Routes */}
        <Route path="/" element={
          <AuthGuard>
            <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700">
              {/* Header */}
              <div className="bg-white bg-opacity-10 backdrop-blur-sm border-b border-white border-opacity-20">
                <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
                  <div className="flex justify-between items-center py-4">
                    <div className="flex items-center space-x-4">
                      <h1 className="text-2xl font-bold text-white">Dealer Development Platform</h1>
                      {user && (
                        <div className="text-blue-200">
                          Добро пожаловать, {user.login} {user.is_admin && '(Админ)'}
                        </div>
                      )}
                    </div>
                    <div className="flex items-center space-x-4">
                      {isAuthenticated && (
                        <button
                          onClick={logout}
                          className="bg-red-600 hover:bg-red-700 text-white px-4 py-2 rounded-lg transition-colors duration-200"
                        >
                          Выйти
                        </button>
                      )}
                    </div>
                  </div>
                </div>
              </div>

              {/* Main Content */}
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
                    {/* Parameters Dropdown */}
                    <div className="relative flex-1 min-w-0">
                      <button 
                        onClick={() => {
                          console.log('Parameters button clicked, current state:', showParametersDropdown)
                          setShowParametersDropdown(!showParametersDropdown)
                        }}
                        className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                      >
                        <span className={selectedParameters ? 'text-gray-900' : 'text-gray-500'}>
                          {selectedParameters ? parameters.find(p => p.id === selectedParameters)?.name : 'Parameters'}
                        </span>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                      {showParametersDropdown && (
                        <div className="absolute z-50 w-full mt-2 bg-white border border-gray-200 rounded-xl shadow-lg max-h-60 overflow-y-auto">
                          {parameters.map((parameter) => (
                            <button
                              key={parameter.id}
                              onClick={() => {
                                setSelectedParameters(parameter.id)
                                setShowParametersDropdown(false)
                              }}
                              className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200"
                            >
                              {parameter.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Region Dropdown */}
                    <div className="relative flex-1 min-w-0">
                      <button 
                        onClick={() => {
                          console.log('Region button clicked, current state:', showRegionDropdown)
                          setShowRegionDropdown(!showRegionDropdown)
                        }}
                        className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                      >
                        <span className={selectedRegion ? 'text-gray-900' : 'text-gray-500'}>
                          {selectedRegion ? regions.find(r => r.id === selectedRegion)?.name : 'Region'}
                        </span>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                      {showRegionDropdown && (
                        <div className="absolute z-50 w-full mt-2 bg-white border border-gray-200 rounded-xl shadow-lg max-h-60 overflow-y-auto">
                          {regions.map((region) => (
                            <button
                              key={region.id}
                              onClick={() => {
                                setSelectedRegion(region.id)
                                setShowRegionDropdown(false)
                              }}
                              className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200"
                            >
                              {region.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Period Dropdown */}
                    <div className="relative flex-1 min-w-0">
                      <button 
                        onClick={() => {
                          console.log('Period button clicked, current state:', showPeriodDropdown)
                          setShowPeriodDropdown(!showPeriodDropdown)
                        }}
                        className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                      >
                        <span className={selectedPeriod ? 'text-gray-900' : 'text-gray-500'}>
                          {selectedPeriod ? periods.find(p => p.id === selectedPeriod)?.name : 'Period'}
                        </span>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                      {showPeriodDropdown && (
                        <div className="absolute z-50 w-full mt-2 bg-white border border-gray-200 rounded-xl shadow-lg max-h-60 overflow-y-auto">
                          {periods.map((period) => (
                            <button
                              key={period.id}
                              onClick={() => {
                                setSelectedPeriod(period.id)
                                setShowPeriodDropdown(false)
                              }}
                              className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200"
                            >
                              {period.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Year Dropdown */}
                    <div className="relative flex-1 min-w-0">
                      <button 
                        onClick={() => {
                          console.log('Year button clicked, current state:', showYearDropdown)
                          setShowYearDropdown(!showYearDropdown)
                        }}
                        className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                      >
                        <span className={selectedYear ? 'text-gray-900' : 'text-gray-500'}>
                          {selectedYear ? years.find(y => y.id === selectedYear)?.name : 'Year'}
                        </span>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                      {showYearDropdown && (
                        <div className="absolute z-50 w-full mt-2 bg-white border border-gray-200 rounded-xl shadow-lg max-h-60 overflow-y-auto">
                          {years.map((year) => (
                            <button
                              key={year.id}
                              onClick={() => {
                                setSelectedYear(year.id)
                                setShowYearDropdown(false)
                              }}
                              className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200"
                            >
                              {year.name}
                            </button>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Dealers Selection Button */}
                    <div className="flex-1 min-w-0">
                      <button
                        onClick={() => setShowDealersModal(true)}
                        className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                      >
                        <span className={selectedDealers.length > 0 ? 'text-gray-900' : 'text-gray-500'}>
                          {getSelectedDealersText()}
                        </span>
                        <svg className="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                        </svg>
                      </button>
                    </div>

                    {/* Search Button */}
                    <div className="flex-shrink-0">
                      <button
                        onClick={handleSearch}
                        className="bg-blue-600 hover:bg-blue-700 text-white px-8 py-4 rounded-xl font-medium transition-colors duration-200 flex items-center space-x-2"
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                        </svg>
                        <span>Find Results</span>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </AuthGuard>
        } />
          
        {/* Protected Routes */}
        <Route path="/dealer-dev" element={
          <AuthGuard>
            <SalesTeamTable />
          </AuthGuard>
        } />
        <Route path="/sales" element={
          <AuthGuard>
            <SalesTable />
          </AuthGuard>
        } />
        <Route path="/after-sales" element={
          <AuthGuard>
            <AfterSalesTable />
          </AuthGuard>
        } />
        <Route path="/performance" element={
          <AuthGuard>
            <PerformanceTable />
          </AuthGuard>
        } />
        <Route path="/all" element={
          <AuthGuard>
            <AllTable />
          </AuthGuard>
        } />
        <Route path="/quarter-comparison" element={
          <AuthGuard>
            <QuarterComparison />
          </AuthGuard>
        } />
        <Route path="/brand-demo" element={
          <AuthGuard>
            <BrandDemo />
          </AuthGuard>
        } />
        <Route path="/dealer/:dealerId" element={
          <AuthGuard>
            <DealerCard />
          </AuthGuard>
        } />
          
        {/* Admin Routes */}
        <Route path="/admin" element={
          <AdminGuard>
            <Admin />
          </AdminGuard>
        } />
        <Route path="/excel-upload" element={
          <AdminGuard>
            <ExcelUploadPage />
          </AdminGuard>
        } />
        <Route path="/excel-tables" element={
          <AdminGuard>
            <ExcelTablesPage />
          </AdminGuard>
        } />
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
              </div>
            </div>

            {/* Dealers List */}
            <div className="px-8 py-6 max-h-96 overflow-y-auto">
              {dealersLoading ? (
                <div className="text-center py-8">
                  <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
                  <p className="text-gray-600 mt-2">Loading dealers...</p>
                </div>
              ) : dealersError ? (
                <div className="text-center py-8">
                  <p className="text-red-600">Error loading dealers: {dealersError}</p>
                </div>
              ) : availableDealers.length > 0 ? (
                <div className="space-y-2">
                  {availableDealers
                    .filter(dealer => 
                      dealerSearchQuery === '' || 
                      dealer.name.toLowerCase().includes(dealerSearchQuery.toLowerCase())
                    )
                    .map((dealer) => (
                      <label key={dealer.id} className="flex items-center space-x-3 p-3 hover:bg-gray-50 rounded-lg cursor-pointer">
                        <input
                          type="checkbox"
                          checked={selectedDealers.includes(dealer.id.toString())}
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedDealers([...selectedDealers, dealer.id.toString()])
                            } else {
                              setSelectedDealers(selectedDealers.filter(id => id !== dealer.id.toString()))
                            }
                          }}
                          className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                        />
                        <div className="flex-1">
                          <p className="text-sm font-medium text-gray-900">{dealer.name}</p>
                          <p className="text-sm text-gray-500">{dealer.city}, {dealer.region}</p>
                        </div>
                      </label>
                    ))}
                </div>
              ) : (
                <div className="text-center py-8">
                  <p className="text-gray-600">No dealers found</p>
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="px-8 py-4 border-t border-gray-200 bg-gray-50">
              <div className="flex justify-between items-center">
                <div className="text-sm text-gray-600">
                  {selectedDealers.length} dealer(s) selected
                </div>
                <div className="flex space-x-3">
                  <button
                    onClick={() => {
                      setShowDealersModal(false)
                      setDealerSearchQuery('') // Сбрасываем поиск при закрытии
                    }}
                    className="px-6 py-2 bg-gray-300 text-gray-700 rounded-lg hover:bg-gray-400 transition-colors duration-200"
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
    </>
  )
}

const App: React.FC = () => {
  return (
    <AuthProvider>
      <AppContent />
    </AuthProvider>
  )
}

export default App