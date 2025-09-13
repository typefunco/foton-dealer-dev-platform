import React, { useState } from 'react'
import { Link } from 'react-router-dom'

interface AllDealer {
  id: string
  // Common fields
  name: string
  city: string
  salesManager: string
  
  // Dealer Development fields
  class: 'A' | 'B' | 'C' | 'D'
  checklist: number
  brandsInPortfolio: string[]
  dealerDevRecommendation: 'Planned Result' | 'Needs Development' | 'Find New Candidate' | 'Close Down'
  
  // Sales Team fields
  salesTarget: string
  stockHdtMdtLdt: string
  buyoutHdtMdtLdt: string
  fotonSalesmen: number
  salesTrainings: boolean
  salesDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down'
  
  // Performance fields
  srRub: string
  salesProfit: string
  salesMargin: number
  autoSalesRevenue: string
  autoSalesProfitsRap: string
  autoSalesMargin: number
  ranking: number
  autoSalesDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down'
  
  // After Sales fields
  rStockPercent: number
  wStockPercent: number
  flhPercent: number
  serviceContract: string
  asTrainings: boolean
  csi: string
  asDecision: 'Needs development' | 'Planned Result' | 'Find New Candidate' | 'Close Down'
}

const AllTable: React.FC = () => {
  const [selectedRegion, setSelectedRegion] = useState<string>('all-russia')

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

  const dealers: AllDealer[] = [
    {
      id: '1',
      name: 'AvtoFurgon',
      city: 'Moscow',
      salesManager: 'Kozeev',
      class: 'B',
      checklist: 80,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Needs Development',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Needs development',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Needs development'
    },
    {
      id: '2',
      name: 'Avtokub',
      city: 'Moscow',
      salesManager: 'Kozeev',
      class: 'B',
      checklist: 85,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Needs Development',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Needs development',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Needs development',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Needs development'
    },
    {
      id: '3',
      name: 'Avto-M',
      city: 'Moscow',
      salesManager: 'Kozeev',
      class: 'B',
      checklist: 82,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Needs Development',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Needs development',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Needs development',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Needs development'
    },
    {
      id: '4',
      name: 'BTS Belgorod',
      city: 'Moscow',
      salesManager: 'Kozeev',
      class: 'A',
      checklist: 92,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Planned Result',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Planned Result',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Planned Result',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Planned Result'
    },
    {
      id: '5',
      name: 'BTS Smolensk',
      city: 'Noginsk',
      salesManager: 'Kozeev',
      class: 'A',
      checklist: 95,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Planned Result',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Planned Result',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Planned Result',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Planned Result'
    },
    {
      id: '6',
      name: 'Centr Trak Grupp',
      city: 'Solnechnogorsk',
      salesManager: 'Kozeev',
      class: 'A',
      checklist: 96,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Planned Result',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Planned Result',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Planned Result',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Planned Result'
    },
    {
      id: '7',
      name: 'Ecomtekh',
      city: 'Ecomtekh',
      salesManager: 'Avdeev',
      class: 'A',
      checklist: 91,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Planned Result',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Planned Result',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Planned Result',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Planned Result'
    },
    {
      id: '8',
      name: 'GAS 36',
      city: 'Yaroslavl',
      salesManager: 'Avdeev',
      class: 'C',
      checklist: 76,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Find New Candidate',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Find New Candidate',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Find New Candidate',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Find New Candidate'
    },
    {
      id: '9',
      name: 'Global Truck Sales',
      city: 'Ryazan',
      salesManager: 'Avdeev',
      class: 'C',
      checklist: 73,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Find New Candidate',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Find New Candidate',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Find New Candidate',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Find New Candidate'
    },
    {
      id: '10',
      name: 'Gus Tekhcentr',
      city: 'Tambov',
      salesManager: 'Avdeev',
      class: 'C',
      checklist: 72,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Find New Candidate',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Find New Candidate',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Find New Candidate',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Find New Candidate'
    },
    {
      id: '11',
      name: 'KomDorAvto',
      city: 'Tula',
      salesManager: 'Avdeev',
      class: 'D',
      checklist: 68,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Close Down',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: true,
      salesDecision: 'Close Down',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Close Down',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: true,
      csi: '1 500',
      asDecision: 'Close Down'
    },
    {
      id: '12',
      name: 'Major Trak Centr',
      city: 'Lipeck',
      salesManager: 'Avdeev',
      class: 'D',
      checklist: 66,
      brandsInPortfolio: ['FOTON'],
      dealerDevRecommendation: 'Close Down',
      salesTarget: '40/100',
      stockHdtMdtLdt: '5/2/3',
      buyoutHdtMdtLdt: '5/2/3',
      fotonSalesmen: 5,
      salesTrainings: false,
      salesDecision: 'Close Down',
      srRub: '5 555 555',
      salesProfit: '5 000 000',
      salesMargin: 5,
      autoSalesRevenue: '5 555 555',
      autoSalesProfitsRap: '5 555 555',
      autoSalesMargin: 5,
      ranking: 5,
      autoSalesDecision: 'Close Down',
      rStockPercent: 5,
      wStockPercent: 5,
      flhPercent: 5,
      serviceContract: 'Gold',
      asTrainings: false,
      csi: '800',
      asDecision: 'Close Down'
    }
  ]

  const getClassColor = (dealerClass: string) => {
    switch (dealerClass) {
      case 'A': return 'text-green-400'
      case 'B': return 'text-yellow-400'
      case 'C': return 'text-orange-400'
      case 'D': return 'text-red-400'
      default: return 'text-gray-400'
    }
  }

  const getChecklistColor = (score: number) => {
    if (score >= 90) return 'text-green-600'
    if (score >= 80) return 'text-yellow-600'
    if (score >= 70) return 'text-orange-600'
    return 'text-red-600'
  }

  const getDecisionColor = (decision: string) => {
    switch (decision) {
      case 'Planned Result': return 'text-green-600'
      case 'Needs development':
      case 'Needs Development': return 'text-yellow-600'
      case 'Find New Candidate': return 'text-orange-600'
      case 'Close Down': return 'text-red-600'
      default: return 'text-gray-600'
    }
  }

  const getJointDecision = (dealer: AllDealer) => {
    // Logic to determine joint decision based on all factors
    if (dealer.class === 'D' || dealer.checklist < 70) {
      return 'Close Down'
    }
    if (dealer.class === 'A' && dealer.checklist >= 90 && dealer.salesTrainings && dealer.asTrainings) {
      return 'Planned Result'
    }
    if (dealer.class === 'C' || dealer.checklist < 80) {
      return 'Find New Candidate'
    }
    return 'Needs Development'
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
            COMPREHENSIVE
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            DEALER ANALYSIS
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

      {/* All Data Table with Horizontal Scroll */}
      <div className="w-full px-4 sm:px-6 lg:px-8 pb-8">
        <div className="p-6">
          <h3 className="text-2xl font-bold text-white mb-6 text-center">
            Complete Dealer Data Overview
          </h3>
          
          <div className="overflow-x-auto">
            <div className="min-w-max">
              <table className="w-full text-sm">
                <thead>
                  <tr className="bg-blue-800 bg-opacity-50">
                    {/* Common Fields */}
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      Dealer Name
                    </th>
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      City
                    </th>
                    <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider">
                      Sales Manager
                    </th>
                    
                                         {/* Dealer Development Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Class
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Checklist
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-blue-600 bg-opacity-30">
                       Brands in Portfolio
                     </th>
                     
                     {/* Sales Team Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Sales Target
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       <div>Stock</div>
                       <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       <div>Buyout</div>
                       <div className="text-xs font-normal text-blue-200">hdt/mdt/ldt</div>
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Foton Salesmen
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-green-600 bg-opacity-30">
                       Sales Trainings
                     </th>
                     
                     {/* Performance Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Revenue Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Profit Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Sales Margin %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Revenue Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Profits Rub
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       After Sales Margin %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-yellow-600 bg-opacity-30">
                       Ranking
                     </th>
                     
                     {/* After Sales Fields */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Recommended Stock %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Warranty Stock %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Foton Labor Hours %
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Service Contract
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       AS Trainings
                     </th>
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-purple-600 bg-opacity-30">
                       Foton Warranty Hours
                     </th>
                     
                     {/* Joint Decision Field */}
                     <th className="px-3 py-3 text-center text-xs font-bold text-white uppercase tracking-wider bg-red-600 bg-opacity-30">
                       Joint Decision
                     </th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-blue-200 divide-opacity-30">
                  {dealers.map((dealer) => (
                    <tr key={dealer.id} className="hover:bg-blue-800 hover:bg-opacity-30">
                      {/* Common Fields */}
                      <td className="px-3 py-2 text-center">
                        <Link 
                          to={`/dealer/${dealer.id}`}
                          className="text-xs font-medium text-white hover:text-blue-200 transition-colors duration-200 cursor-pointer"
                        >
                          {dealer.name}
                        </Link>
                      </td>
                      <td className="px-3 py-2 text-center">
                        <div className="text-xs text-white">{dealer.city}</div>
                      </td>
                      <td className="px-3 py-2 text-center">
                        <div className="text-xs text-white">{dealer.salesManager}</div>
                      </td>
                      
                                             {/* Dealer Development Fields */}
                       <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                         <div className={`inline-flex items-center px-2 py-1 rounded-full text-xs font-medium ${getClassColor(dealer.class)}`}>
                           {dealer.class}
                         </div>
                       </td>
                       <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                         <div className={`text-xs font-medium ${getChecklistColor(dealer.checklist)}`}>
                           {dealer.checklist}
                         </div>
                       </td>
                       <td className="px-3 py-2 text-center bg-blue-600 bg-opacity-30">
                         <div className="flex justify-center">
                           {dealer.brandsInPortfolio.map((brand, index) => (
                             <div
                               key={index}
                               className="w-6 h-6 bg-blue-400 bg-opacity-80 rounded-full flex items-center justify-center border border-blue-300"
                               title={brand}
                             >
                               <span className="text-xs font-bold text-white">F</span>
                             </div>
                           ))}
                         </div>
                       </td>
                       
                       {/* Sales Team Fields */}
                       <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.salesTarget}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.stockHdtMdtLdt}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.buyoutHdtMdtLdt}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.fotonSalesmen}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-green-600 bg-opacity-30">
                         <div className={`text-xs font-medium ${dealer.salesTrainings ? 'text-green-400' : 'text-white'}`}>
                           {dealer.salesTrainings ? 'Yes' : 'No'}
                         </div>
                       </td>
                       
                       {/* Performance Fields */}
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.srRub}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.salesProfit}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.salesMargin}%</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.autoSalesRevenue}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.autoSalesProfitsRap}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.autoSalesMargin}%</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-yellow-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.ranking}</div>
                       </td>
                       
                       {/* After Sales Fields */}
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.rStockPercent}%</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.wStockPercent}%</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.flhPercent}%</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.serviceContract}</div>
                       </td>
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className={`text-xs font-medium ${dealer.asTrainings ? 'text-green-400' : 'text-white'}`}>
                           {dealer.asTrainings ? 'Yes' : 'No'}
                         </div>
                       </td>
                       <td className="px-3 py-2 text-center bg-purple-600 bg-opacity-30">
                         <div className="text-xs text-white">{dealer.csi}</div>
                       </td>
                      
                       {/* Joint Decision Field */}
                       <td className="px-3 py-2 text-center bg-red-600 bg-opacity-30">
                         <div className={`text-xs font-medium ${getDecisionColor(getJointDecision(dealer))}`}>
                           {getJointDecision(dealer)}
                         </div>
                       </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
          
          
        </div>
      </div>
    </div>
  )
}

export default AllTable
