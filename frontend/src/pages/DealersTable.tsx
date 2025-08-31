import React, { useState } from 'react'
import { Link } from 'react-router-dom'

interface Dealer {
  id: string
  name: string
  city: string
  class: 'A' | 'B' | 'C' | 'D'
  checklist: number
  brandsInPortfolio: string[]
  branding: boolean
  marketingInvestments: number // в миллионах рублей
  buySideBusiness: string[]
  dealerDevRecommendation: 'Planned Result' | 'Needs Development' | 'Find New Candidate' | 'Close Down'
}

const DealersTable: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<string>('all-russia')
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
      
      if (sortConfig.key === 'class') {
        const classOrder = { 'A': 4, 'B': 3, 'C': 2, 'D': 1 }
        const aOrder = classOrder[aValue as 'A' | 'B' | 'C' | 'D']
        const bOrder = classOrder[bValue as 'A' | 'B' | 'C' | 'D']
        
        if (sortConfig.direction === 'asc') {
          return bOrder - aOrder // A, B, C, D
        } else {
          return aOrder - bOrder // D, C, B, A
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

  const regions = [
    { id: 'all-russia', name: 'All Russia' },
    { id: 'central', name: 'Central' },
    { id: 'ural', name: 'Ural' },
    { id: 'volga', name: 'Volga' },
    { id: 'kavkaz', name: 'Kavkaz' },
    { id: 'north-west', name: 'North West' },
    { id: 'siberia', name: 'Siberia' },
    { id: 'south', name: 'South' }
  ]

  const dealers: Dealer[] = [
    {
      id: '1',
      name: 'AvtoFurgon',
      city: 'Moscow',
      class: 'B',
      checklist: 80,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 2.5,
      buySideBusiness: ['Logistics', 'Warehousing'],
      dealerDevRecommendation: 'Needs Development'
    },
    {
      id: '2',
      name: 'Avtokub',
      city: 'Moscow',
      class: 'B',
      checklist: 85,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 1.8,
      buySideBusiness: ['Transport'],
      dealerDevRecommendation: 'Needs Development'
    },
    {
      id: '3',
      name: 'Avto-M',
      city: 'Moscow',
      class: 'B',
      checklist: 82,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 3.2,
      buySideBusiness: ['Logistics', 'Retail'],
      dealerDevRecommendation: 'Needs Development'
    },
    {
      id: '4',
      name: 'BTS Belgorod',
      city: 'Moscow',
      class: 'A',
      checklist: 92,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 5.0,
      buySideBusiness: ['Logistics', 'Warehousing', 'Retail'],
      dealerDevRecommendation: 'Planned Result'
    },
    {
      id: '5',
      name: 'BTS Smolensk',
      city: 'Noginsk',
      class: 'A',
      checklist: 95,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 4.8,
      buySideBusiness: ['Logistics', 'Retail'],
      dealerDevRecommendation: 'Planned Result'
    },
    {
      id: '6',
      name: 'Centr Trak Grupp',
      city: 'Solnechnogorsk',
      class: 'A',
      checklist: 96,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 6.2,
      buySideBusiness: ['Logistics', 'Warehousing', 'Retail', 'Service'],
      dealerDevRecommendation: 'Planned Result'
    },
    {
      id: '7',
      name: 'Ecomtekh',
      city: 'Ecomtekh',
      class: 'A',
      checklist: 91,
      brandsInPortfolio: ['FOTON'],
      branding: true,
      marketingInvestments: 4.5,
      buySideBusiness: ['Logistics', 'Retail'],
      dealerDevRecommendation: 'Planned Result'
    },
    {
      id: '8',
      name: 'GAS 36',
      city: 'Yaroslavl',
      class: 'C',
      checklist: 76,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 1.2,
      buySideBusiness: ['Transport'],
      dealerDevRecommendation: 'Find New Candidate'
    },
    {
      id: '9',
      name: 'Global Truck Sales',
      city: 'Ryazan',
      class: 'C',
      checklist: 73,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 0.8,
      buySideBusiness: ['Transport'],
      dealerDevRecommendation: 'Find New Candidate'
    },
    {
      id: '10',
      name: 'Gus Tekhcentr',
      city: 'Tambov',
      class: 'C',
      checklist: 72,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 0.6,
      buySideBusiness: ['Transport'],
      dealerDevRecommendation: 'Find New Candidate'
    },
    {
      id: '11',
      name: 'KomDorAvto',
      city: 'Tula',
      class: 'D',
      checklist: 68,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 0.3,
      buySideBusiness: [],
      dealerDevRecommendation: 'Close Down'
    },
    {
      id: '12',
      name: 'Major Trak Centr',
      city: 'Lipeck',
      class: 'D',
      checklist: 66,
      brandsInPortfolio: ['FOTON'],
      branding: false,
      marketingInvestments: 0.2,
      buySideBusiness: [],
      dealerDevRecommendation: 'Close Down'
    }
  ]

  const getClassColor = (dealerClass: string) => {
    switch (dealerClass) {
      case 'A': return 'border-green-500 text-green-600'
      case 'B': return 'border-yellow-500 text-yellow-600'
      case 'C': return 'border-orange-500 text-orange-600'
      case 'D': return 'border-red-500 text-red-600'
      default: return 'border-gray-300 text-gray-600'
    }
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

  const getChecklistColor = (score: number) => {
    if (score >= 90) return 'text-green-600'
    if (score >= 80) return 'text-yellow-600'
    if (score >= 70) return 'text-orange-600'
    return 'text-red-600'
  }

  const getRecommendationColor = (recommendation: string) => {
    switch (recommendation) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs Development': return 'text-yellow-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-gray-600'
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
        <div className="overflow-x-auto">
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
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Checklist
                </th>
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Brands in Portfolio
                </th>
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Branding
                </th>
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Marketing Investments (M RUB)
                </th>
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Buy Side Business
                </th>
                <th className="px-6 py-4 text-center text-sm font-bold text-white uppercase tracking-wider">
                  Dealer Dev Recommendation
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
                    <div className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium border-2 ${getClassColor(dealer.class)}`}>
                      {dealer.class}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className={`text-sm font-medium ${getChecklistColor(dealer.checklist)}`}>
                      {dealer.checklist}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                                            <div className="flex justify-center space-x-2">
                          {dealer.brandsInPortfolio.map((brand, index) => (
                            <div
                              key={index}
                              className="w-8 h-8 bg-blue-400 bg-opacity-80 rounded-full flex items-center justify-center border border-blue-300"
                              title={brand}
                            >
                              <span className="text-xs font-bold text-white">F</span>
                            </div>
                          ))}
                        </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className={`text-sm font-medium ${dealer.branding ? 'text-green-600' : 'text-red-600'}`}>
                      {dealer.branding ? 'Yes' : 'No'}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="text-sm text-white font-medium">
                      {dealer.marketingInvestments} M RUB
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className="flex justify-center space-x-2">
                      {dealer.buySideBusiness.length > 0 ? (
                        dealer.buySideBusiness.map((business, index) => (
                          <div
                            key={index}
                            className="w-8 h-8 bg-purple-400 bg-opacity-80 rounded-full flex items-center justify-center border border-purple-300"
                            title={business}
                          >
                            <span className="text-xs font-bold text-white">B</span>
                          </div>
                        ))
                      ) : (
                        <div className="text-sm text-gray-400">-</div>
                      )}
                    </div>
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-center">
                    <div className={`text-sm font-medium ${getRecommendationColor(dealer.dealerDevRecommendation)}`}>
                      {dealer.dealerDevRecommendation}
                    </div>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  )
}

export default DealersTable
