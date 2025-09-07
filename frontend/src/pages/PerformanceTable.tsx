import React, { useState } from 'react'
import { Link } from 'react-router-dom'

interface PerformanceDealer {
  id: string
  name: string
  city: string
  srRub: string
  salesProfit: number
  salesMargin: number
  autoSalesRevenue: string
  rap: string
  autoSalesProfitsRap: string
  autoSalesMargin: number
  marketingInvestment: number
  ranking: number
  autoSalesDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down'
}

const PerformanceTable: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<string>('center')
  const [sortConfig, setSortConfig] = useState<{
    key: keyof PerformanceDealer | null
    direction: 'asc' | 'desc' | null
  }>({ key: null, direction: null })

  const handleSort = (key: keyof PerformanceDealer) => {
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

  const regions = [
    { id: 'all-russia', name: 'All Russia' },
    { id: 'center', name: 'Center' },
    { id: 'north-west', name: 'North West' },
    { id: 'volga', name: 'Volga' },
    { id: 'south', name: 'South' },
    { id: 'ural', name: 'Ural' },
    { id: 'siberia', name: 'Siberia' },
    { id: 'far-east', name: 'Far East' }
  ]

  const dealers: PerformanceDealer[] = [
    {
      id: '1',
      name: 'AvtoFurgon',
      city: 'Moscow',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 2.5,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '2',
      name: 'Avtokub',
      city: 'Moscow',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 3.2,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '3',
      name: 'Avto-M',
      city: 'Moscow',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 1.8,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '4',
      name: 'BTS Belgorod',
      city: 'Moscow',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 4.1,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '5',
      name: 'BTS Smolensk',
      city: 'Noginsk',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 2.9,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '6',
      name: 'Centr Trak Grupp',
      city: 'Solnechnogorsk',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 3.7,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '7',
      name: 'Ecomtekh',
      city: 'Ecomtekh',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 2.3,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '8',
      name: 'GAS 36',
      city: 'Yaroslavl',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 1.5,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '9',
      name: 'Global Truck Sales',
      city: 'Ryazan',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 2.1,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '10',
      name: 'Gus Tekhcentr',
      city: 'Tambov',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 1.9,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '11',
      name: 'KomDorAvto',
      city: 'Tula',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 1.2,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    },
    {
      id: '12',
      name: 'Major Trak Centr',
      city: 'Lipeck',
      srRub: '5 555 555',
      salesProfit: 5,
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      rap: 'Gold',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      marketingInvestment: 2.8,
      ranking: 5,
      autoSalesDecision: 'Needs development'
    }
  ]

  const getSortedDealers = () => {
    if (!sortConfig.key || !sortConfig.direction) {
      return dealers
    }

    return [...dealers].sort((a, b) => {
      const aValue = a[sortConfig.key!]
      const bValue = b[sortConfig.key!]
      
      if (sortConfig.key === 'autoSalesDecision') {
        const decisionOrder = { 
          'Planned Result': 4, 
          'Needs development': 3, 
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
      
      // Для строковых значений с числами (srRub, autoSalesRevenue, autoSalesProfitsRap)
      if (sortConfig.key === 'srRub' || sortConfig.key === 'autoSalesRevenue' || sortConfig.key === 'autoSalesProfitsRap') {
        const aNum = parseFloat((aValue as string).replace(/[^\d.-]/g, ''))
        const bNum = parseFloat((bValue as string).replace(/[^\d.-]/g, ''))
        
        if (sortConfig.direction === 'asc') {
          return aNum - bNum
        } else {
          return bNum - aNum
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

  const getSortIcon = (key: keyof PerformanceDealer) => {
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

  const getAutoSalesDecisionColor = (decision: string) => {
    switch (decision) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs development': return 'text-green-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-green-600'
    }
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
            PERFORMANCE
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

      {/* Performance Table */}
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
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('srRub')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Revenue Rub</span>
                  {getSortIcon('srRub')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('salesProfit')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Profit %</span>
                  {getSortIcon('salesProfit')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('salesMargin')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Margin %</span>
                  {getSortIcon('salesMargin')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('autoSalesRevenue')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>After Sales Revenue Rub</span>
                  {getSortIcon('autoSalesRevenue')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('autoSalesProfitsRap')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>After Sales Profits Rub</span>
                  {getSortIcon('autoSalesProfitsRap')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('autoSalesMargin')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>After Sales Margin %</span>
                  {getSortIcon('autoSalesMargin')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('marketingInvestment')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Marketing Investments (M Rub)</span>
                  {getSortIcon('marketingInvestment')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('ranking')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Foton Rank</span>
                  {getSortIcon('ranking')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-800 hover:bg-opacity-30 transition-colors duration-200"
                onClick={() => handleSort('autoSalesDecision')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Performance Decision</span>
                  {getSortIcon('autoSalesDecision')}
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
                  <div className="text-sm text-white">{dealer.srRub}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.salesProfit}%</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.salesMargin}%</div>
                </td>
                                 <td className="px-6 py-4 whitespace-nowrap text-center">
                   <div className="text-sm text-white">{dealer.autoSalesRevenue}</div>
                 </td>
                 <td className="px-6 py-4 whitespace-nowrap text-center">
                   <div className="text-sm text-white">{dealer.autoSalesProfitsRap}</div>
                 </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.autoSalesMargin}%</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.marketingInvestment} M Rub</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.ranking}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className={`text-sm font-medium ${getAutoSalesDecisionColor(dealer.autoSalesDecision)}`}>
                    {dealer.autoSalesDecision}
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

export default PerformanceTable
