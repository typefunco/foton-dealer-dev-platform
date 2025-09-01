import React, { useState } from 'react'
import { Link } from 'react-router-dom'

interface QuarterData {
  quarter: string
  year: string
  metrics: {
    // Dealer Development
    averageClass: string
    averageChecklist: number
    classDistribution: { [key: string]: number }
    
    // Sales Team
    averageSalesTarget: string
    averageFotonSalesmen: number
    salesTrainingsPercentage: number
    salesTrainingsDistribution: { [key: string]: number }
    
    // Performance
    averageSalesRevenue: string
    averageSalesProfit: number
    averageSalesMargin: number
    averageRanking: number
    marketingInvestment: number
    autoSalesRevenue: number
    autoSalesProfit: number
    autoSalesMargin: number
    
    // After Sales
    averageRStockPercent: number
    averageWStockPercent: number
    averageFlhPercent: number
    asTrainingsPercentage: number
    asTrainingsDistribution: { [key: string]: number }
    csiPercentage: number
    
    // Decisions
    decisionDistribution: { [key: string]: number }
  }
}

const QuarterComparison: React.FC = () => {
  const [quarter1, setQuarter1] = useState<string>('q1')
  const [quarter2, setQuarter2] = useState<string>('q2')
  const [year1, setYear1] = useState<string>('2025')
  const [year2, setYear2] = useState<string>('2025')
  
  // Applied filters state
  const [appliedQuarter1, setAppliedQuarter1] = useState<string>('q1')
  const [appliedQuarter2, setAppliedQuarter2] = useState<string>('q2')
  const [appliedYear1, setAppliedYear1] = useState<string>('2025')
  const [appliedYear2, setAppliedYear2] = useState<string>('2025')

  const quarters = [
    { id: 'q1', name: 'Q1' },
    { id: 'q2', name: 'Q2' },
    { id: 'q3', name: 'Q3' },
    { id: 'q4', name: 'Q4' }
  ]

  const years = [
    { id: '2024', name: '2024' },
    { id: '2025', name: '2025' },
    { id: '2026', name: '2026' }
  ]

  // Mock data for two quarters
  const quarter1Data: QuarterData = {
    quarter: 'q1',
    year: year1,
    metrics: {
      averageClass: 'B',
      averageChecklist: 82.5,
      classDistribution: { 'A': 35, 'B': 40, 'C': 20, 'D': 5 },
      averageSalesTarget: '45/100',
      averageFotonSalesmen: 4.8,
      salesTrainingsPercentage: 65,
      salesTrainingsDistribution: { 'Yes': 65, 'No': 35 },
      averageSalesRevenue: '5,200,000',
      averageSalesProfit: 2.8,
      averageSalesMargin: 4.2,
      averageRanking: 6.2,
      marketingInvestment: 850,
      autoSalesRevenue: 3200,
      autoSalesProfit: 180,
      autoSalesMargin: 5.6,
      averageRStockPercent: 78,
      averageWStockPercent: 82,
      averageFlhPercent: 75,
      asTrainingsPercentage: 70,
      asTrainingsDistribution: { 'Yes': 70, 'No': 30 },
      csiPercentage: 68,
      decisionDistribution: { 'Planned Result': 30, 'Needs Development': 45, 'Find New Candidate': 20, 'Close Down': 5 }
    }
  }

  const quarter2Data: QuarterData = {
    quarter: 'q2',
    year: year2,
    metrics: {
      averageClass: 'A',
      averageChecklist: 87.2,
      classDistribution: { 'A': 45, 'B': 35, 'C': 15, 'D': 5 },
      averageSalesTarget: '52/100',
      averageFotonSalesmen: 5.2,
      salesTrainingsPercentage: 78,
      salesTrainingsDistribution: { 'Yes': 78, 'No': 22 },
      averageSalesRevenue: '6,100,000',
      averageSalesProfit: 3.5,
      averageSalesMargin: 4.8,
      averageRanking: 5.8,
      marketingInvestment: 920,
      autoSalesRevenue: 3800,
      autoSalesProfit: 220,
      autoSalesMargin: 5.8,
      averageRStockPercent: 82,
      averageWStockPercent: 86,
      averageFlhPercent: 79,
      asTrainingsPercentage: 75,
      asTrainingsDistribution: { 'Yes': 75, 'No': 25 },
      csiPercentage: 72,
      decisionDistribution: { 'Planned Result': 40, 'Needs Development': 40, 'Find New Candidate': 15, 'Close Down': 5 }
    }
  }

  // Function to get data based on quarter and year
  const getQuarterData = (quarter: string, year: string): QuarterData => {
    if (year === '2026') {
      // 2026 data - improved performance
      return {
        quarter,
        year,
        metrics: {
          averageClass: 'A',
          averageChecklist: 89.5,
          classDistribution: { 'A': 50, 'B': 30, 'C': 15, 'D': 5 },
          averageSalesTarget: '58/100',
          averageFotonSalesmen: 5.8,
          salesTrainingsPercentage: 85,
          salesTrainingsDistribution: { 'Yes': 85, 'No': 15 },
          averageSalesRevenue: '7,200,000',
          averageSalesProfit: 4.2,
          averageSalesMargin: 5.5,
          averageRanking: 4.8,
          marketingInvestment: 1050,
          autoSalesRevenue: 4500,
          autoSalesProfit: 280,
          autoSalesMargin: 6.2,
          averageRStockPercent: 88,
          averageWStockPercent: 92,
          averageFlhPercent: 84,
          asTrainingsPercentage: 82,
          asTrainingsDistribution: { 'Yes': 82, 'No': 18 },
          csiPercentage: 78,
          decisionDistribution: { 'Planned Result': 50, 'Needs Development': 35, 'Find New Candidate': 10, 'Close Down': 5 }
        }
      }
    }
    
    // Default 2025 data
    if (quarter === 'q1') {
      return quarter1Data
    } else {
      return quarter2Data
    }
  }

  // Get current data based on selected quarters and years
  const currentQuarter1Data = getQuarterData(appliedQuarter1, appliedYear1)
  const currentQuarter2Data = getQuarterData(appliedQuarter2, appliedYear2)

  const getClassColor = (dealerClass: string) => {
    switch (dealerClass) {
      case 'A': return 'text-green-400'
      case 'B': return 'text-yellow-400'
      case 'C': return 'text-orange-400'
      case 'D': return 'text-red-400'
      default: return 'text-gray-400'
    }
  }

  const getDecisionColor = (decision: string) => {
    switch (decision) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs Development': return 'text-yellow-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-gray-600'
    }
  }

  const handleApplyFilters = () => {
    setAppliedQuarter1(quarter1)
    setAppliedQuarter2(quarter2)
    setAppliedYear1(year1)
    setAppliedYear2(year2)
  }

  const renderPieChart = (data: { [key: string]: number }, colors: string[], title: string) => {
    const total = Object.values(data).reduce((sum, value) => sum + value, 0)
    let currentAngle = 0

    return (
      <div className="relative w-48 h-48">
        <svg className="w-full h-full transform -rotate-90">
          {Object.entries(data).map(([key, value], index) => {
            const percentage = (value / total) * 100
            const angle = (percentage / 100) * 360
            const x1 = 96 + 84 * Math.cos((currentAngle * Math.PI) / 180)
            const y1 = 96 + 84 * Math.sin((currentAngle * Math.PI) / 180)
            const x2 = 96 + 84 * Math.cos(((currentAngle + angle) * Math.PI) / 180)
            const y2 = 96 + 84 * Math.sin(((currentAngle + angle) * Math.PI) / 180)
            
            const largeArcFlag = angle > 180 ? 1 : 0
            
            const pathData = [
              `M 96 96`,
              `L ${x1} ${y1}`,
              `A 84 84 0 ${largeArcFlag} 1 ${x2} ${y2}`,
              'Z'
            ].join(' ')

            currentAngle += angle

            return (
              <path
                key={key}
                d={pathData}
                fill={colors[index % colors.length]}
                className="transition-all duration-300 hover:opacity-80"
              />
            )
          })}
        </svg>
        <div className="absolute inset-0 flex items-center justify-center">
          <div className="text-center">
            <div className="text-lg font-bold text-white">{total}</div>
            <div className="text-xs text-gray-300">{title}</div>
          </div>
        </div>
      </div>
    )
  }

  const renderBarChart = (value1: number, value2: number, label: string, unit: string = '') => {
    const maxValue = Math.max(value1, value2)
    const percentage1 = (value1 / maxValue) * 100
    const percentage2 = (value2 / maxValue) * 100

    return (
      <div className="space-y-2">
        <div className="text-sm font-medium text-white">{label}</div>
        <div className="space-y-1">
          <div className="flex items-center space-x-2">
            <div className="w-16 text-xs text-gray-300">Q1</div>
            <div className="flex-1 bg-gray-700 rounded-full h-3">
              <div 
                className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                style={{ width: `${percentage1}%` }}
              />
            </div>
            <div className="w-12 text-xs text-white">{value1}{unit}</div>
          </div>
          <div className="flex items-center space-x-2">
            <div className="w-16 text-xs text-gray-300">Q2</div>
            <div className="flex-1 bg-gray-700 rounded-full h-3">
              <div 
                className="bg-green-500 h-3 rounded-full transition-all duration-500"
                style={{ width: `${percentage2}%` }}
              />
            </div>
            <div className="w-12 text-xs text-white">{value2}{unit}</div>
          </div>
        </div>
      </div>
    )
  }

  const renderMetricCard = (title: string, value1: string | number, value2: string | number, unit: string = '', isPercentage: boolean = false) => {
    const num1 = typeof value1 === 'string' ? parseFloat(value1.replace(/[^\d.-]/g, '')) : value1
    const num2 = typeof value2 === 'string' ? parseFloat(value2.replace(/[^\d.-]/g, '')) : value2
    const change = num2 - num1
    const changePercent = num1 !== 0 ? (change / num1) * 100 : 0

    return (
      <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-4 border border-white border-opacity-20">
        <div className="text-sm font-medium text-gray-300 mb-2">{title}</div>
        <div className="flex items-center justify-between">
          <div className="text-2xl font-bold text-white">
            {value1}{unit}
          </div>
          <div className="text-2xl font-bold text-white">
            {value2}{unit}
          </div>
        </div>
        <div className={`text-sm font-medium mt-1 ${change >= 0 ? 'text-green-400' : 'text-red-400'}`}>
          {change >= 0 ? '+' : ''}{change}{unit} ({changePercent >= 0 ? '+' : ''}{changePercent.toFixed(1)}%)
        </div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700">
      {/* Header */}
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
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
            COMPARISON BY QUARTER
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            TOTAL ANALYTICS
          </h2>
        </div>

        {/* Quarter Selection */}
        <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 mt-8">
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
            <div className="flex flex-col lg:flex-row gap-6 items-center justify-center">
              {/* Quarter 1 Selection */}
              <div className="flex items-center space-x-4">
                <span className="text-white font-medium">Quarter 1:</span>
                <select
                  value={quarter1}
                  onChange={(e) => setQuarter1(e.target.value)}
                  className="bg-white bg-opacity-20 text-white border border-white border-opacity-30 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-300"
                >
                  {quarters.map((q) => (
                    <option key={q.id} value={q.id} className="bg-blue-800">
                      {q.name}
                    </option>
                  ))}
                </select>
                <span className="text-white font-medium">Year:</span>
                <select
                  value={year1}
                  onChange={(e) => setYear1(e.target.value)}
                  className="bg-white bg-opacity-20 text-white border border-white border-opacity-30 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-300"
                >
                  {years.map((y) => (
                    <option key={y.id} value={y.id} className="bg-blue-800">
                      {y.name}
                    </option>
                  ))}
                </select>
              </div>

              {/* Quarter 2 Selection */}
              <div className="flex items-center space-x-4">
                <span className="text-white font-medium">Quarter 2:</span>
                <select
                  value={quarter2}
                  onChange={(e) => setQuarter2(e.target.value)}
                  className="bg-white bg-opacity-20 text-white border border-white border-opacity-30 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-300"
                >
                  {quarters.map((q) => (
                    <option key={q.id} value={q.id} className="bg-blue-800">
                      {q.name}
                    </option>
                  ))}
                </select>
                <span className="text-white font-medium">Year:</span>
                <select
                  value={year2}
                  onChange={(e) => setYear2(e.target.value)}
                  className="bg-white bg-opacity-20 text-white border border-white border-opacity-30 rounded-lg px-4 py-2 focus:outline-none focus:border-blue-300"
                >
                  {years.map((y) => (
                    <option key={y.id} value={y.id} className="bg-blue-800">
                      {y.name}
                    </option>
                  ))}
                </select>
              </div>
            </div>
            
            {/* Apply Button */}
            <div className="flex justify-center mt-6">
              <button
                onClick={handleApplyFilters}
                className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-8 rounded-xl transition-colors duration-200 flex items-center space-x-2 shadow-lg hover:shadow-xl"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                </svg>
                <span>Go</span>
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Comparison Content */}
      <div className="max-w-8xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Left Side - Quarter 1 */}
          <div className="space-y-6">
            <div className="text-center">
              <h3 className="text-3xl font-bold text-white mb-2">
                {quarters.find(q => q.id === appliedQuarter1)?.name} {appliedYear1}
              </h3>
              <div className="w-32 h-1 bg-blue-400 mx-auto rounded-full"></div>
            </div>

            {/* Dealer Development Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-blue-500 rounded-full mr-3"></div>
                Dealer Development
              </h4>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Average Class */}
                <div className="text-center">
                  <div className="text-sm text-gray-300 mb-2">Average Class</div>
                  <div className={`text-3xl font-bold ${getClassColor(currentQuarter1Data.metrics.averageClass)}`}>
                    {currentQuarter1Data.metrics.averageClass}
                  </div>
                </div>

                {/* Average Checklist */}
                <div className="text-center">
                  <div className="text-sm text-gray-300 mb-2">Average Checklist</div>
                  <div className="text-3xl font-bold text-white">
                    {currentQuarter1Data.metrics.averageChecklist}
                  </div>
                </div>

                {/* Class Distribution */}
                <div className="md:col-span-2">
                  <div className="text-sm text-gray-300 mb-3 text-center">By Class</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      currentQuarter1Data.metrics.classDistribution,
                      ['#10B981', '#F59E0B', '#F97316', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    {Object.entries(currentQuarter1Data.metrics.classDistribution).map(([key, value]) => (
                      <div key={key} className="flex items-center space-x-1">
                        <div className={`w-3 h-3 rounded-full ${getClassColor(key).replace('text-', 'bg-')}`}></div>
                        <span className="text-xs text-gray-300">{key}: {value}%</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            {/* Sales Team Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-green-500 rounded-full mr-3"></div>
                Sales Team
              </h4>
              
              <div className="space-y-4">
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Sales Target (Current/Yearly)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${parseFloat(currentQuarter1Data.metrics.averageSalesTarget.split('/')[0])}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter1Data.metrics.averageSalesTarget}</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Sales Trainings</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      currentQuarter1Data.metrics.salesTrainingsDistribution,
                      ['#10B981', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                      <span className="text-xs text-gray-300">Yes: {currentQuarter1Data.metrics.salesTrainingsDistribution['Yes']}%</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <span className="text-xs text-gray-300">No: {currentQuarter1Data.metrics.salesTrainingsDistribution['No']}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* After Sales Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-purple-500 rounded-full mr-3"></div>
                After Sales
              </h4>
              
              <div className="space-y-4">
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Recommended Stock (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter1Data.metrics.averageRStockPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter1Data.metrics.averageRStockPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Warranty Stock (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter1Data.metrics.averageWStockPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter1Data.metrics.averageWStockPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Foton Labor Hours (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter1Data.metrics.averageFlhPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter1Data.metrics.averageFlhPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">AS Trainings</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      currentQuarter1Data.metrics.asTrainingsDistribution,
                      ['#10B981', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                      <span className="text-xs text-gray-300">Yes: {currentQuarter1Data.metrics.asTrainingsDistribution['Yes']}%</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <span className="text-xs text-gray-300">No: {currentQuarter1Data.metrics.asTrainingsDistribution['No']}%</span>
                    </div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">CSI (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-blue-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter1Data.metrics.csiPercentage}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter1Data.metrics.csiPercentage}%</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Performance Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-yellow-500 rounded-full mr-3"></div>
                Performance
              </h4>
              
              <div className="space-y-4">
                {renderMetricCard(
                  'Sales Revenue',
                  currentQuarter1Data.metrics.averageSalesRevenue,
                  currentQuarter1Data.metrics.averageSalesRevenue,
                  ' ₽'
                )}
                
                {renderMetricCard(
                  'Sales Profit',
                  currentQuarter1Data.metrics.averageSalesProfit,
                  currentQuarter1Data.metrics.averageSalesProfit,
                  ' M Rub'
                )}
                
                {renderMetricCard(
                  'Sales Margin',
                  currentQuarter1Data.metrics.averageSalesMargin,
                  currentQuarter1Data.metrics.averageSalesMargin,
                  '%',
                  true
                )}
                
                {renderMetricCard(
                  'After Sales Revenue',
                  currentQuarter1Data.metrics.autoSalesRevenue,
                  currentQuarter1Data.metrics.autoSalesRevenue,
                  ' Rub'
                )}
                
                {renderMetricCard(
                  'After Sales Profit',
                  currentQuarter1Data.metrics.autoSalesProfit,
                  currentQuarter1Data.metrics.autoSalesProfit,
                  ' Rub'
                )}
                
                {renderMetricCard(
                  'After Sales Margin',
                  currentQuarter1Data.metrics.autoSalesMargin,
                  currentQuarter1Data.metrics.autoSalesMargin,
                  '%',
                  true
                )}
                
                {renderMetricCard(
                  'Marketing Investment',
                  currentQuarter1Data.metrics.marketingInvestment,
                  currentQuarter1Data.metrics.marketingInvestment,
                  ' Rub'
                )}
              </div>
            </div>

            {/* Decision Distribution */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-red-500 rounded-full mr-3"></div>
                Decision Chart
              </h4>
              
              <div className="flex justify-center">
                {renderPieChart(
                  currentQuarter1Data.metrics.decisionDistribution,
                  ['#10B981', '#F59E0B', '#F97316', '#EF4444'],
                  'Decisions'
                )}
              </div>
              <div className="flex flex-wrap justify-center gap-4 mt-3">
                {Object.entries(currentQuarter1Data.metrics.decisionDistribution).map(([key, value]) => (
                  <div key={key} className="flex items-center space-x-1">
                    <div className={`w-3 h-3 rounded-full ${getDecisionColor(key).replace('text-', 'bg-')}`}></div>
                    <span className="text-xs text-gray-300">{key}: {value}%</span>
                  </div>
                ))}
              </div>
            </div>
          </div>

          {/* Right Side - Quarter 2 */}
          <div className="space-y-6">
            <div className="text-center">
              <h3 className="text-3xl font-bold text-white mb-2">
                {quarters.find(q => q.id === appliedQuarter2)?.name} {appliedYear2}
              </h3>
              <div className="w-32 h-1 bg-green-400 mx-auto rounded-full"></div>
            </div>

            {/* Dealer Development Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-blue-500 rounded-full mr-3"></div>
                Dealer Development
              </h4>
              
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Average Class */}
                <div className="text-center">
                  <div className="text-sm text-gray-300 mb-2">Average Class</div>
                  <div className={`text-3xl font-bold ${getClassColor(quarter2Data.metrics.averageClass)}`}>
                    {quarter2Data.metrics.averageClass}
                  </div>
                </div>

                {/* Average Checklist */}
                <div className="text-center">
                  <div className="text-sm text-gray-300 mb-2">Average Checklist</div>
                  <div className="text-3xl font-bold text-white">
                    {quarter2Data.metrics.averageChecklist}
                  </div>
                </div>

                {/* Class Distribution */}
                <div className="md:col-span-2">
                  <div className="text-sm text-gray-300 mb-3 text-center">By Class</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      quarter2Data.metrics.classDistribution,
                      ['#10B981', '#F59E0B', '#F97316', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    {Object.entries(quarter2Data.metrics.classDistribution).map(([key, value]) => (
                      <div key={key} className="flex items-center space-x-1">
                        <div className={`w-3 h-3 rounded-full ${getClassColor(key).replace('text-', 'bg-')}`}></div>
                        <span className="text-xs text-gray-300">{key}: {value}%</span>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            </div>

            {/* Sales Team Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-green-500 rounded-full mr-3"></div>
                Sales Team
              </h4>
              
              <div className="space-y-4">
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Sales Target (Current/Yearly)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-green-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${parseFloat(quarter2Data.metrics.averageSalesTarget.split('/')[0])}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{quarter2Data.metrics.averageSalesTarget}</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Sales Trainings</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      quarter2Data.metrics.salesTrainingsDistribution,
                      ['#10B981', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                      <span className="text-xs text-gray-300">Yes: {quarter2Data.metrics.salesTrainingsDistribution['Yes']}%</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <span className="text-xs text-gray-300">No: {quarter2Data.metrics.salesTrainingsDistribution['No']}%</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* After Sales Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-purple-500 rounded-full mr-3"></div>
                After Sales
              </h4>
              
              <div className="space-y-4">
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Recommended Stock (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-green-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter2Data.metrics.averageRStockPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter2Data.metrics.averageRStockPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Warranty Stock (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-green-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter2Data.metrics.averageWStockPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter2Data.metrics.averageWStockPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">Average Foton Labor Hours (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-green-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter2Data.metrics.averageFlhPercent}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter2Data.metrics.averageFlhPercent}%</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">AS Trainings</div>
                  <div className="flex justify-center">
                    {renderPieChart(
                      currentQuarter2Data.metrics.asTrainingsDistribution,
                      ['#10B981', '#EF4444'],
                      'Dealers'
                    )}
                  </div>
                  <div className="flex justify-center space-x-4 mt-3">
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-green-500"></div>
                      <span className="text-xs text-gray-300">Yes: {currentQuarter2Data.metrics.asTrainingsDistribution['Yes']}%</span>
                    </div>
                    <div className="flex items-center space-x-1">
                      <div className="w-3 h-3 rounded-full bg-red-500"></div>
                      <span className="text-xs text-gray-300">No: {currentQuarter2Data.metrics.asTrainingsDistribution['No']}%</span>
                    </div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <div className="text-sm font-medium text-white">CSI (%)</div>
                  <div className="flex items-center space-x-2">
                    <div className="flex-1 bg-gray-700 rounded-full h-3">
                      <div 
                        className="bg-green-500 h-3 rounded-full transition-all duration-500"
                        style={{ width: `${currentQuarter2Data.metrics.csiPercentage}%` }}
                      />
                    </div>
                    <div className="w-12 text-xs text-white">{currentQuarter2Data.metrics.csiPercentage}%</div>
                  </div>
                </div>
              </div>
            </div>

            {/* Performance Metrics */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-yellow-500 rounded-full mr-3"></div>
                Performance
              </h4>
              
              <div className="space-y-4">
                {renderMetricCard(
                  'Sales Revenue',
                  currentQuarter2Data.metrics.averageSalesRevenue,
                  currentQuarter2Data.metrics.averageSalesRevenue,
                  ' ₽'
                )}
                
                {renderMetricCard(
                  'Sales Profit',
                  currentQuarter2Data.metrics.averageSalesProfit,
                  currentQuarter2Data.metrics.averageSalesProfit,
                  ' M Rub'
                )}
                
                {renderMetricCard(
                  'Sales Margin',
                  currentQuarter2Data.metrics.averageSalesMargin,
                  currentQuarter2Data.metrics.averageSalesMargin,
                  '%',
                  true
                )}
                
                {renderMetricCard(
                  'After Sales Revenue',
                  currentQuarter2Data.metrics.autoSalesRevenue,
                  currentQuarter2Data.metrics.autoSalesRevenue,
                  ' Rub'
                )}
                
                {renderMetricCard(
                  'After Sales Profit',
                  currentQuarter2Data.metrics.autoSalesProfit,
                  currentQuarter2Data.metrics.autoSalesProfit,
                  ' Rub'
                )}
                
                {renderMetricCard(
                  'After Sales Margin',
                  currentQuarter2Data.metrics.autoSalesMargin,
                  currentQuarter2Data.metrics.autoSalesMargin,
                  '%',
                  true
                )}
                
                {renderMetricCard(
                  'Marketing Investment',
                  currentQuarter2Data.metrics.marketingInvestment,
                  currentQuarter2Data.metrics.marketingInvestment,
                  ' Rub'
                )}
              </div>
            </div>

            {/* Decision Distribution */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-6 border border-white border-opacity-20">
              <h4 className="text-xl font-bold text-white mb-4 flex items-center">
                <div className="w-4 h-4 bg-red-500 rounded-full mr-3"></div>
                Decision Chart
              </h4>
              
              <div className="flex justify-center">
                {renderPieChart(
                  quarter2Data.metrics.decisionDistribution,
                  ['#10B981', '#F59E0B', '#F97316', '#EF4444'],
                  'Decisions'
                )}
              </div>
              <div className="flex flex-wrap justify-center gap-4 mt-3">
                {Object.entries(quarter2Data.metrics.decisionDistribution).map(([key, value]) => (
                  <div key={key} className="flex items-center space-x-1">
                    <div className={`w-3 h-3 rounded-full ${getDecisionColor(key).replace('text-', 'bg-')}`}></div>
                    <span className="text-xs text-gray-300">{key}: {value}%</span>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default QuarterComparison
