import React, { useState } from 'react'
import { Link } from 'react-router-dom'

interface Dealer {
  id: string
  name: string
  city: string
  salesManager: string
  salesTarget: string
  stockHdtMdtLdt: string
  buyoutHdtMdtLdt: string
  fotonSalesmen: number
  salesTrainings: boolean
  salesDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down'
}

const SalesTeamTable: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<string>('center')
  const [sortConfig, setSortConfig] = useState<{
    key: keyof Dealer | null
    direction: 'asc' | 'desc' | null
  }>({ key: null, direction: null })

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
    if (!sortConfig.key || !sortConfig.direction) {
      return dealers
    }

    return [...dealers].sort((a, b) => {
      const aValue = a[sortConfig.key!]
      const bValue = b[sortConfig.key!]
      
      if (sortConfig.key === 'salesTrainings') {
        const aBool = aValue as boolean
        const bBool = bValue as boolean
        
        if (sortConfig.direction === 'asc') {
          return aBool === bBool ? 0 : aBool ? -1 : 1 // true first
        } else {
          return aBool === bBool ? 0 : aBool ? 1 : -1 // false first
        }
      }
      
      if (sortConfig.key === 'salesDecision') {
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
    { id: 'center', name: 'Center' },
    { id: 'north-west', name: 'North West' },
    { id: 'volga', name: 'Volga' },
    { id: 'south', name: 'South' },
    { id: 'ural', name: 'Ural' },
    { id: 'siberia', name: 'Siberia' },
    { id: 'far-east', name: 'Far East' }
  ]

  const dealers: Dealer[] = [
    {
      id: '1',
      name: 'AvtoFurgon',
      city: 'Moscow',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '2',
      name: 'Avtokub',
      city: 'Moscow',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    },
    {
      id: '3',
      name: 'Avto-M',
      city: 'Moscow',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '4',
      name: 'BTS Belgorod',
      city: 'Moscow',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    },
    {
      id: '5',
      name: 'BTS Smolensk',
      city: 'Noginsk',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '6',
      name: 'Centr Trak Grupp',
      city: 'Solnechnogorsk',
      salesManager: 'Kozeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    },
    {
      id: '7',
      name: 'Ecomtekh',
      city: 'Ecomtekh',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '8',
      name: 'GAS 36',
      city: 'Yaroslavl',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    },
    {
      id: '9',
      name: 'Global Truck Sales',
      city: 'Ryazan',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '10',
      name: 'Gus Tekhcentr',
      city: 'Tambov',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    },
    {
      id: '11',
      name: 'KomDorAvto',
      city: 'Tula',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development'
    },
    {
      id: '12',
      name: 'Major Trak Centr',
      city: 'Lipeck',
      salesManager: 'Avdeev',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development'
    }
  ]

  const getSalesDecisionColor = (decision: string) => {
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
             SALES
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
                onClick={() => handleSort('salesManager')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Manager</span>
                  {getSortIcon('salesManager')}
                </div>
              </th>
              <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                Sales Target
              </th>
              <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                <div>Stock</div>
                <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
              </th>
              <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                <div>Buyout</div>
                <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('fotonSalesmen')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Foton Salesmen</span>
                  {getSortIcon('fotonSalesmen')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('salesTrainings')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Trainings</span>
                  {getSortIcon('salesTrainings')}
                </div>
              </th>
              <th 
                className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider cursor-pointer hover:bg-blue-700 hover:bg-opacity-80 hover:shadow-lg hover:shadow-blue-500/50 transition-all duration-300 rounded-3xl mx-2"
                onClick={() => handleSort('salesDecision')}
              >
                <div className="flex items-center justify-center space-x-1">
                  <span>Sales Decision</span>
                  {getSortIcon('salesDecision')}
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
                  <div className="text-sm text-white">{dealer.salesManager}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.salesTarget}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.stockHdtMdtLdt}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.buyoutHdtMdtLdt}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className="text-sm text-white">{dealer.fotonSalesmen}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className={`text-sm font-medium ${
                    dealer.salesTrainings ? 'text-green-600' : 'text-white'
                  }`}>
                    {dealer.salesTrainings ? 'Yes' : 'No'}
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-center">
                  <div className={`text-sm font-medium ${getSalesDecisionColor(dealer.salesDecision)}`}>
                    {dealer.salesDecision}
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
