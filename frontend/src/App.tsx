import React, { useState } from 'react'
import { Routes, Route, Link } from 'react-router-dom'
import DealersTable from './pages/DealersTable'
import SalesTeamTable from './pages/SalesTeamTable'
import AfterSalesTable from './pages/AfterSalesTable'
import PerformanceTable from './pages/PerformanceTable'
import AllTable from './pages/AllTable'
import Login from './pages/Login'
import ForgotPassword from './pages/ForgotPassword'
import Admin from './pages/Admin'

const App: React.FC = () => {
  const [selectedCategory, setSelectedCategory] = useState<string>('')
  const [selectedFrom, setSelectedFrom] = useState<string>('')
  const [selectedTo, setSelectedTo] = useState<string>('')
  const [selectedYear, setSelectedYear] = useState<string>('')
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)
  const [showFromDropdown, setShowFromDropdown] = useState(false)
  const [showToDropdown, setShowToDropdown] = useState(false)
  const [showYearDropdown, setShowYearDropdown] = useState(false)
  const [showAdvancedModal, setShowAdvancedModal] = useState(false)
  
  // Advanced Settings states
  const [selectedRegion, setSelectedRegion] = useState<string>('')
  const [dealerType, setDealerType] = useState<'single' | 'multiple'>('multiple')
  const [selectedDealers, setSelectedDealers] = useState<string[]>([])
  const [advancedSettingsApplied, setAdvancedSettingsApplied] = useState(false)

  const categories = [
    { id: 'all', name: 'ALL', description: 'All Categories' },
    { id: 'after-sales', name: 'After Sales', description: 'After Sales Service' },
    { id: 'dealer-dev', name: 'Dealer Dev', description: 'Dealer Development' },
    { id: 'sales', name: 'Sales', description: 'Sales' },
    { id: 'performance', name: 'Performance', description: 'Performance' }
  ]

  const quarters = [
    { id: 'q1', name: 'Q1', description: 'Quarter 1' },
    { id: 'q2', name: 'Q2', description: 'Quarter 2' },
    { id: 'q3', name: 'Q3', description: 'Quarter 3' },
    { id: 'q4', name: 'Q4', description: 'Quarter 4' }
  ]

  const years = [
    { id: '2024', name: '2024', description: 'Year 2024' },
    { id: '2025', name: '2025', description: 'Year 2025' },
    { id: '2026', name: '2026', description: 'Year 2026' },
    { id: '2027', name: '2027', description: 'Year 2027' }
  ]

  const regions = [
    { id: 'central', name: 'Central', description: 'Central Region' },
    { id: 'volga', name: 'Volga', description: 'Volga Region' },
    { id: 'kavkaz', name: 'Kavkaz', description: 'Caucasus Region' },
    { id: 'south', name: 'South', description: 'Southern Region' },
    { id: 'north-west', name: 'North West', description: 'North West Region' },
    { id: 'far-east', name: 'Far East', description: 'Far East Region' },
    { id: 'ural', name: 'Ural', description: 'Ural Region' }
  ]

  const sampleDealers = [
    { id: 'dealer1', name: 'AutoDealer Moscow', region: 'Central' },
    { id: 'dealer2', name: 'AutoDealer St. Petersburg', region: 'North West' },
    { id: 'dealer3', name: 'AutoDealer Kazan', region: 'Volga' },
    { id: 'dealer4', name: 'AutoDealer Rostov', region: 'South' },
    { id: 'dealer5', name: 'AutoDealer Yekaterinburg', region: 'Ural' }
  ]

  const handleCategorySelect = (categoryId: string) => {
    setSelectedCategory(categoryId)
    setShowCategoryDropdown(false)
  }

  const handleFromSelect = (quarterId: string) => {
    setSelectedFrom(quarterId)
    setShowFromDropdown(false)
  }

  const handleToSelect = (quarterId: string) => {
    setSelectedTo(quarterId)
    setShowToDropdown(false)
  }

  const handleYearSelect = (yearId: string) => {
    setSelectedYear(yearId)
    setShowYearDropdown(false)
  }

  const closeAllDropdowns = () => {
    setShowCategoryDropdown(false)
    setShowFromDropdown(false)
    setShowToDropdown(false)
    setShowYearDropdown(false)
  }

  const handleAdvancedSettingsSave = () => {
    console.log('Advanced Settings:', {
      region: selectedRegion,
      dealerType,
      selectedDealers
    })
    setAdvancedSettingsApplied(true)
    setShowAdvancedModal(false)
  }

  const handleDealerToggle = (dealerId: string) => {
    if (dealerType === 'single') {
      setSelectedDealers([dealerId])
    } else {
      setSelectedDealers(prev => 
        prev.includes(dealerId) 
          ? prev.filter(id => id !== dealerId)
          : [...prev, dealerId]
      )
    }
  }

  const handleDealerTypeChange = (type: 'single' | 'multiple') => {
    setDealerType(type)
    // Clear selected dealers when switching types
    setSelectedDealers([])
  }

  const resetAdvancedSettings = () => {
    setSelectedRegion('')
    setDealerType('multiple')
    setSelectedDealers([])
    setAdvancedSettingsApplied(false)
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
              <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
                DEVELOPMENT PLATFORM
              </h2>
            </div>

            {/* Search Panel */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 flex items-center justify-center min-h-[400px]">
              <div className="bg-white rounded-2xl shadow-5xl p-10 w-full">
                <div className="flex flex-col lg:flex-row gap-6 items-center justify-center">
                  {/* Category Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowCategoryDropdown(!showCategoryDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedCategory ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedCategory ? categories.find(cat => cat.id === selectedCategory)?.name : 'Category'}
                      </span>
                      <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showCategoryDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                    
                    {showCategoryDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50 max-h-60 overflow-y-auto">
                        {categories.map((category) => (
                          <button
                            key={category.id}
                            onClick={() => handleCategorySelect(category.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{category.name}</div>
                            <div className="text-sm text-gray-500">{category.description}</div>
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
                      <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showYearDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
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
                            <div className="text-sm text-gray-500">{year.description}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* From Quarter Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowFromDropdown(!showFromDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedFrom ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedFrom ? quarters.find(q => q.id === selectedFrom)?.name : 'From'}
                      </span>
                      <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showFromDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                    
                    {showFromDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
                        {quarters.map((quarter) => (
                          <button
                            key={quarter.id}
                            onClick={() => handleFromSelect(quarter.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{quarter.name}</div>
                            <div className="text-sm text-gray-500">{quarter.description}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* To Quarter Dropdown */}
                  <div className="relative flex-1 min-w-0">
                    <button 
                      onClick={() => setShowToDropdown(!showToDropdown)}
                      className="w-full bg-white border-2 border-gray-200 rounded-xl px-4 py-4 text-left hover:border-blue-300 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center justify-between"
                    >
                      <span className={selectedTo ? 'text-gray-900' : 'text-gray-500'}>
                        {selectedTo ? quarters.find(q => q.id === selectedTo)?.name : 'To'}
                      </span>
                      <svg className={`w-5 h-5 text-gray-400 transition-transform duration-200 ${showToDropdown ? 'rotate-180' : ''}`} fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                      </svg>
                    </button>
                    
                    {showToDropdown && (
                      <div className="absolute top-full left-0 right-0 mt-2 bg-white border border-gray-200 rounded-xl shadow-lg z-50">
                        {quarters.map((quarter) => (
                          <button
                            key={quarter.id}
                            onClick={() => handleToSelect(quarter.id)}
                            className="w-full px-4 py-3 text-left hover:bg-blue-50 transition-colors duration-200 first:rounded-t-xl last:rounded-b-xl"
                          >
                            <div className="font-medium text-gray-900">{quarter.name}</div>
                            <div className="text-sm text-gray-500">{quarter.description}</div>
                          </button>
                        ))}
                      </div>
                    )}
                  </div>

                  {/* Advanced Settings Button */}
                  <button 
                    onClick={() => setShowAdvancedModal(true)}
                    className="bg-white border-2 border-gray-200 rounded-xl px-6 py-4 text-gray-700 hover:border-blue-300 hover:bg-blue-50 focus:border-blue-500 focus:outline-none transition-colors duration-200 flex items-center space-x-2"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    </svg>
                    <span className="font-medium">Advanced Settings</span>
                  </button>

                  {/* Find Results Button */}
                  <button 
                    onClick={() => {
                      console.log('Searching with:', {
                        category: selectedCategory,
                        from: selectedFrom,
                        to: selectedTo,
                        advancedSettings: advancedSettingsApplied ? {
                          region: selectedRegion,
                          dealerType,
                          selectedDealers
                        } : null
                      })
                      // Search logic will be here
                    }}
                    className="bg-orange-500 hover:bg-orange-600 text-white font-bold py-4 px-8 rounded-xl transition-colors duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                    </svg>
                    <span>Find results</span>
                  </button>
                </div>
              </div>
            </div>

            {/* Table Navigation Links */}
            <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-8">
              <div className="bg-white rounded-2xl shadow-xl p-8">
                <h3 className="text-2xl font-bold text-gray-900 mb-6 text-center">
                  Access Data Tables
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                  <Link to="/dealers" className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    Dealers Table
                  </Link>
                  <Link to="/sales-team" className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    Sales Team Table
                  </Link>
                  <Link to="/after-sales" className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    After Sales Table
                  </Link>
                  <Link to="/performance" className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    Performance Table
                  </Link>
                  <Link to="/all" className="bg-purple-600 hover:bg-purple-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    All Data Table
                  </Link>
                  <Link to="/admin" className="bg-green-600 hover:bg-green-700 text-white px-6 py-4 rounded-xl font-medium transition-all duration-200 transform hover:scale-105 text-center">
                    Admin Panel
                  </Link>
                </div>
              </div>
            </div>



            {/* Applied Advanced Settings Display */}
            {advancedSettingsApplied && (
              <div className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 -mt-5">
                <div className="bg-white rounded-2xl shadow-xl p-6 border-l-4 border-blue-500">
                  <div className="flex justify-between items-start mb-4">
                    <div className="flex items-center space-x-3">
                      <div className="w-10 h-10 bg-blue-100 rounded-full flex items-center justify-center">
                        <svg className="w-5 h-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                        </svg>
                      </div>
                      <h3 className="text-lg font-semibold text-gray-900">Applied Settings</h3>
                    </div>
                    <button
                      onClick={resetAdvancedSettings}
                      className="text-gray-400 hover:text-gray-600 transition-colors duration-200 p-2 hover:bg-gray-100 rounded-lg"
                                              title="Reset Settings"
                    >
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                      </svg>
                    </button>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {/* Region */}
                    {selectedRegion && (
                      <div className="bg-gray-50 rounded-lg p-4">
                        <div className="text-sm font-medium text-gray-500 mb-1">Region</div>
                        <div className="text-gray-900 font-medium">
                          {regions.find(r => r.id === selectedRegion)?.name}
                        </div>
                        <div className="text-sm text-gray-500">
                          {regions.find(r => r.id === selectedRegion)?.description}
                        </div>
                      </div>
                    )}

                    {/* Dealer Type */}
                    <div className="bg-gray-50 rounded-lg p-4">
                                              <div className="text-sm font-medium text-gray-500 mb-1">Dealer Type</div>
                        <div className="text-gray-900 font-medium">
                          {dealerType === 'single' ? 'Single Dealer' : 'Multiple Dealers'}
                        </div>
                        <div className="text-sm text-gray-500">
                          {dealerType === 'single' ? 'Select specific dealer' : 'Select multiple dealers'}
                        </div>
                    </div>

                    {/* Selected Dealers */}
                    {dealerType === 'single' && selectedDealers.length > 0 && (
                      <div className="bg-gray-50 rounded-lg p-4">
                        <div className="text-sm font-medium text-gray-500 mb-1">Selected Dealer</div>
                        <div className="text-gray-900 font-medium">
                          {sampleDealers.find(d => d.id === selectedDealers[0])?.name}
                        </div>
                        <div className="text-sm text-gray-500">
                          {sampleDealers.find(d => d.id === selectedDealers[0])?.region}
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            )}
          </>
        } />
        <Route path="/dealers" element={<DealersTable />} />
        <Route path="/sales-team" element={<SalesTeamTable />} />
        <Route path="/after-sales" element={<AfterSalesTable />} />
        <Route path="/performance" element={<PerformanceTable />} />
        <Route path="/all" element={<AllTable />} />
        <Route path="/login" element={<Login />} />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/admin" element={<Admin />} />
      </Routes>

      {/* Advanced Settings Modal */}
      {showAdvancedModal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
          {/* Backdrop with blur */}
          <div 
            className="absolute inset-0 bg-black bg-opacity-50 backdrop-blur-sm"
            onClick={() => setShowAdvancedModal(false)}
          />
          
          {/* Modal Content */}
          <div className="relative bg-white rounded-2xl shadow-2xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
            {/* Header */}
            <div className="sticky top-0 bg-white px-8 py-6 border-b border-gray-200 rounded-t-2xl">
              <div className="flex justify-between items-center">
                <h3 className="text-2xl font-bold text-gray-900">Advanced Settings</h3>
                <button
                  onClick={() => setShowAdvancedModal(false)}
                  className="text-gray-400 hover:text-gray-600 transition-colors duration-200 p-2 hover:bg-gray-100 rounded-lg"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>

            {/* Modal Body */}
            <div className="px-8 py-6 space-y-8">
              {/* Region Selection */}
              <div>
                <h4 className="text-lg font-semibold text-gray-900 mb-4">Region Selection</h4>
                <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
                  {regions.map((region) => (
                    <button
                      key={region.id}
                      onClick={() => setSelectedRegion(region.id)}
                      className={`p-4 border-2 rounded-xl text-left transition-all duration-200 ${
                        selectedRegion === region.id
                          ? 'border-blue-500 bg-blue-50 text-blue-900'
                          : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50'
                      }`}
                    >
                      <div className="font-medium">{region.name}</div>
                      <div className="text-sm text-gray-500">{region.description}</div>
                    </button>
                  ))}
                </div>
              </div>

              {/* Dealer Type Selection */}
              <div>
                <h4 className="text-lg font-semibold text-gray-900 mb-4">Dealer Selection Type</h4>
                <div className="flex space-x-4">
                  <button
                    onClick={() => handleDealerTypeChange('multiple')}
                    className={`flex-1 p-4 border-2 rounded-xl text-center transition-all duration-200 ${
                      dealerType === 'multiple'
                        ? 'border-blue-500 bg-blue-50 text-blue-900'
                        : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50'
                    }`}
                  >
                                          <div className="font-medium">Multiple Dealers</div>
                      <div className="text-sm text-gray-500">Select multiple dealers</div>
                  </button>
                  
                  <button
                    onClick={() => handleDealerTypeChange('single')}
                    className={`flex-1 p-4 border-2 rounded-xl text-center transition-all duration-200 ${
                      dealerType === 'single'
                        ? 'border-blue-500 bg-blue-50 text-blue-900'
                        : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50'
                    }`}
                  >
                                          <div className="font-medium">Single Dealer</div>
                      <div className="text-sm text-gray-500">Select one dealer</div>
                  </button>
                </div>
              </div>

              {/* Dealer Selection - Only show when "Single Dealer" is selected */}
              {selectedRegion && dealerType === 'single' && (
                <div>
                                      <h4 className="text-lg font-semibold text-gray-900 mb-4">
                      Dealer Selection
                    </h4>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                    {sampleDealers
                      .filter(dealer => dealer.region.toLowerCase().includes(selectedRegion))
                      .map((dealer) => (
                        <button
                          key={dealer.id}
                          onClick={() => handleDealerToggle(dealer.id)}
                          className={`p-4 border-2 rounded-xl text-left transition-all duration-200 ${
                            selectedDealers.includes(dealer.id)
                              ? 'border-blue-500 bg-blue-50 text-blue-900'
                              : 'border-gray-200 hover:border-blue-300 hover:bg-blue-50'
                          }`}
                        >
                          <div className="font-medium">{dealer.name}</div>
                          <div className="text-sm text-gray-500">{dealer.region}</div>
                        </button>
                      ))}
                  </div>
                  {sampleDealers.filter(dealer => dealer.region.toLowerCase().includes(selectedRegion)).length === 0 && (
                    <p className="text-gray-500 text-center py-4">No dealers found for selected region</p>
                  )}
                </div>
              )}
            </div>

            {/* Footer */}
            <div className="sticky bottom-0 bg-white px-8 py-6 border-t border-gray-200 rounded-b-2xl">
              <div className="flex justify-end space-x-3">
                                  <button
                    onClick={() => setShowAdvancedModal(false)}
                    className="px-6 py-2 text-gray-600 hover:text-gray-800 transition-colors duration-200"
                  >
                    Cancel
                  </button>
                                  <button
                    onClick={handleAdvancedSettingsSave}
                    className="px-6 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors duration-200"
                  >
                    Apply Settings
                  </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Close dropdowns when clicking outside */}
      {(showCategoryDropdown || showFromDropdown || showToDropdown || showYearDropdown) && (
        <div 
          className="fixed inset-0 z-40" 
          onClick={closeAllDropdowns}
        />
      )}
    </div>
  )
}

export default App
