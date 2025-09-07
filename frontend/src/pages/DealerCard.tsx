import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { PieChart, Pie, Cell, Tooltip, ResponsiveContainer } from 'recharts'
import BrandLogos from '../components/BrandLogos'

interface DealerData {
  id: string
  name: string
  city: string
  salesManager: string
  class: string
  checklist: number
  brandsInPortfolio: string[]
  brandsCount: number
  salesTarget: string
  stockHdtMdtLdt: string
  buyoutHdtMdtLdt: string
  fotonSalesmen: string
  salesTrainings: boolean
  srRub: string
  salesProfit: number
  salesMargin: number
  afterSalesRevenue: string
  afterSalesProfitsRap: string
  afterSalesMargin: number
  ranking: number
  recommendedStock: number
  warrantyStock: number
  fotonLaborHours: number
  serviceContract: boolean
  asTrainings: boolean
  csi: number
  branding: boolean
  buySideBusiness: string[]
  marketingInvestment: number
}

const DealerCard: React.FC = () => {
  const { dealerId } = useParams<{ dealerId: string }>()
  const [dealer, setDealer] = useState<DealerData | null>(null)

  // Mock data - в реальном приложении это будет загружаться с сервера
  useEffect(() => {
    const mockDealer: DealerData = {
      id: 'dealer1',
      name: 'AutoDealer Moscow',
      city: 'Moscow',
      salesManager: 'Ivan Petrov',
      class: 'A',
      checklist: 92,
      brandsInPortfolio: ['FOTON', 'SHACMAN', 'DONGFENG', 'KAMAZ', 'GAZ'],
      brandsCount: 5,
      salesTarget: '100',
      stockHdtMdtLdt: 'HDT: 50, MDT: 30, LDT: 20',
      buyoutHdtMdtLdt: 'HDT: 45, MDT: 25, LDT: 15',
      fotonSalesmen: '15',
      salesTrainings: true,
      srRub: '15,000,000',
      salesProfit: 25,
      salesMargin: 18,
      afterSalesRevenue: '8,000,000',
      afterSalesProfitsRap: '2,400,000',
      afterSalesMargin: 30,
      ranking: 1,
      recommendedStock: 85,
      warrantyStock: 70,
      fotonLaborHours: 92,
      serviceContract: true,
      asTrainings: true,
      csi: 95,
      branding: true,
      buySideBusiness: ['Logistics', 'Warehousing', 'Retail'],
      marketingInvestment: 3200
    }

    setDealer(mockDealer)
  }, [dealerId])

  // Данные для pie charts - единые цвета для HDT, MDT, LDT
  const stockData = [
    { name: 'HDT', value: 50, color: '#3B82F6' },
    { name: 'MDT', value: 30, color: '#10B981' },
    { name: 'LDT', value: 20, color: '#8B5CF6' }
  ]

  const buyoutData = [
    { name: 'HDT', value: 45, color: '#3B82F6' }, // Тот же цвет что и для стоков
    { name: 'MDT', value: 25, color: '#10B981' }, // Тот же цвет что и для стоков
    { name: 'LDT', value: 15, color: '#8B5CF6' }  // Тот же цвет что и для стоков
  ]



  // Данные для sales target pie chart - план по году и выполнение в квартале
  const salesTargetData = [
    { name: 'Completed', value: 32, color: '#10B981' },
    { name: 'Remaining', value: 68, color: '#EF4444' }
  ]

  if (!dealer) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-white text-xl">Loading dealer data...</div>
      </div>
    )
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700">
      {/* Header */}
      <div className="bg-white bg-opacity-10 backdrop-blur-sm border-b border-white border-opacity-20">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
          <div className="flex items-center justify-between">
            <div>
              <Link 
                to="/" 
                className="text-blue-200 hover:text-white transition-colors duration-200 flex items-center space-x-2"
                title="Back to Home"
              >
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 19l-7-7m0 0l7-7m-7 7h18" />
                </svg>
                <span>Back to Home</span>
              </Link>
              <h1 className="text-3xl font-bold text-white mt-2">{dealer.name}</h1>
              <p className="text-blue-200 text-lg">{dealer.city} • {dealer.salesManager}</p>
            </div>
          </div>
        </div>
      </div>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Basic Info Cards */}
        <div className="space-y-6 mb-8">
          {/* Dealer Development Group */}
          <div>
            <h3 className="text-lg font-semibold text-white mb-4 text-center">Dealer Development</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
                <div className="text-blue-200 text-sm font-medium">Class</div>
                <div className="text-2xl font-bold text-white">{dealer.class}</div>
              </div>
              <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
                <div className="text-blue-200 text-sm font-medium">Checklist</div>
                <div className="text-2xl font-bold text-white">{dealer.checklist}</div>
              </div>
              <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
                <div className="text-blue-200 text-sm font-medium">Branding</div>
                <div className={`text-2xl font-bold ${dealer.branding ? 'text-green-400' : 'text-red-400'}`}>
                  {dealer.branding ? 'YES' : 'NO'}
                </div>
              </div>
              <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
                <div className="text-blue-200 text-sm font-medium">Byside Businesses</div>
                <div className="flex justify-center space-x-2 mt-2">
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
              </div>
            </div>
          </div>
        </div>

        {/* Brands Portfolio Visualization */}
        <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20 mb-8">
          <h3 className="text-xl font-semibold text-white mb-4">Brands Portfolio</h3>
          <BrandLogos brands={dealer.brandsInPortfolio} className="justify-center" />
        </div>

        {/* Charts Section */}
        <div className="space-y-8">
          <h2 className="text-2xl font-bold text-white text-center mb-8">
            Sales
          </h2>
          
          {/* Sales Target Pie Chart */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Sales Target Performance</h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 items-center">
              <div className="text-center">
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={salesTargetData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, value }) => `${name}: ${value}`}
                      outerRadius={100}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {salesTargetData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'rgba(0,0,0,0.9)', 
                        border: 'none', 
                        borderRadius: '8px',
                        color: 'white',
                        fontSize: '14px',
                        fontWeight: '500'
                      }}
                      labelStyle={{ color: 'white' }}
                      itemStyle={{ color: 'white' }}
                    />
                  </PieChart>
                </ResponsiveContainer>
              </div>
              <div className="text-center lg:text-left">
                <div className="space-y-4">
                  <div>
                    <h4 className="text-lg font-medium text-white mb-2">Annual Target in Units</h4>
                    <p className="text-3xl font-bold text-white">100</p>
                  </div>
                  <div>
                    <h4 className="text-lg font-medium text-white mb-2">Delivered in Units</h4>
                    <p className="text-3xl font-bold text-green-400">32</p>
                  </div>
                  <div>
                    <h4 className="text-lg font-bold text-white mb-2">Remaining in Units</h4>
                    <p className="text-3xl font-bold text-red-400">68</p>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Stock & Buyout by Segments */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Stock & Buyout by Segments</h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Stock Distribution */}
              <div className="text-center">
                <h4 className="text-lg font-medium text-white mb-4">Stock</h4>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={stockData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name} ${percent ? (percent * 100).toFixed(0) : 0}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {stockData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'rgba(0,0,0,0.9)', 
                        border: 'none', 
                        borderRadius: '8px',
                        color: 'white',
                        fontSize: '14px',
                        fontWeight: '500'
                      }}
                      labelStyle={{ color: 'white' }}
                      itemStyle={{ color: 'white' }}
                    />
                  </PieChart>
                </ResponsiveContainer>
                <div className="mt-4 space-y-2">
                  {stockData.map((item, index) => (
                    <div key={index} className="flex items-center justify-center space-x-2">
                      <div className="w-4 h-4 rounded-full" style={{ backgroundColor: item.color }}></div>
                      <span className="text-white text-sm">{item.name}: {item.value}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* Buyout Distribution */}
              <div className="text-center">
                <h4 className="text-lg font-medium text-white mb-4">Buyout</h4>
                <ResponsiveContainer width="100%" height={300}>
                  <PieChart>
                    <Pie
                      data={buyoutData}
                      cx="50%"
                      cy="50%"
                      labelLine={false}
                      label={({ name, percent }) => `${name} ${percent ? (percent * 100).toFixed(0) : 0}%`}
                      outerRadius={80}
                      fill="#8884d8"
                      dataKey="value"
                    >
                      {buyoutData.map((entry, index) => (
                        <Cell key={`cell-${index}`} fill={entry.color} />
                      ))}
                    </Pie>
                    <Tooltip 
                      contentStyle={{ 
                        backgroundColor: 'rgba(0,0,0,0.9)', 
                        border: 'none', 
                        borderRadius: '8px',
                        color: 'white',
                        fontSize: '14px',
                        fontWeight: '500'
                      }}
                      labelStyle={{ color: 'white' }}
                      itemStyle={{ color: 'white' }}
                    />
                  </PieChart>
                </ResponsiveContainer>
                <div className="mt-4 space-y-2">
                  {buyoutData.map((item, index) => (
                    <div key={index} className="flex items-center justify-center space-x-2">
                      <div className="w-4 h-4 rounded-full" style={{ backgroundColor: item.color }}></div>
                      <span className="text-white text-sm">{item.name}: {item.value}</span>
                    </div>
                  ))}
                </div>
              </div>
            </div>
            
            {/* Unified Color Legend */}
            <div className="mt-6 pt-6 border-t border-white border-opacity-20">
              <h5 className="text-lg font-medium text-white mb-3 text-center">Color Legend</h5>
              <div className="flex justify-center space-x-8">
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 rounded-full bg-blue-500"></div>
                  <span className="text-white text-sm">HDT (Heavy Duty Truck)</span>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 rounded-full bg-green-500"></div>
                  <span className="text-white text-sm">MDT (Medium Duty Truck)</span>
                </div>
                <div className="flex items-center space-x-2">
                  <div className="w-4 h-4 rounded-full bg-purple-500"></div>
                  <span className="text-white text-sm">LDT (Light Duty Truck)</span>
                </div>
              </div>
            </div>
          </div>

          {/* Sales Performance Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* Sales Revenue */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Sales Revenue</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.srRub}
                </div>
                <div className="text-white text-sm">
                  Total sales revenue
                </div>
              </div>
            </div>

            {/* Sales Profit */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Sales Profit</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.salesProfit} M Rub
                </div>
                <div className="text-white text-sm">
                  Absolute profit value
                </div>
              </div>
            </div>

            {/* Sales Margin */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Sales Margin</h3>
              <div className="text-center">
                <div className={`text-3xl font-bold mb-2 ${dealer.salesMargin > 0 ? 'text-green-400' : 'text-red-400'}`}>
                  {dealer.salesMargin}%
                </div>
                <div className="text-white text-sm">
                  Gross margin percentage
                </div>
              </div>
            </div>
          </div>

          {/* After Sales Analytics */}
          <h2 className="text-2xl font-bold text-white text-center mb-8">
            After Sales
          </h2>

          {/* Stock and Labor Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* Recommended Stock */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Recommended Stock</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.recommendedStock}%
                </div>
                <div className="text-white text-sm">
                  Recommended stock level
                </div>
              </div>
            </div>

            {/* Warranty Stock */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Warranty Stock</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.warrantyStock}%
                </div>
                <div className="text-white text-sm">
                  Warranty stock level
                </div>
              </div>
            </div>

            {/* Foton Labor Hours */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">Foton Labor Hours</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.fotonLaborHours}%
                </div>
                <div className="text-white text-sm">
                  Foton labor utilization
                </div>
              </div>
            </div>

            {/* CSI */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">CSI</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.csi}%
                </div>
                <div className="text-white text-sm">
                  Customer satisfaction index
                </div>
              </div>
            </div>
          </div>

          {/* After Sales Performance Metrics */}
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            {/* After Sales Revenue */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">After Sales Revenue</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.afterSalesRevenue}
                </div>
                <div className="text-white text-sm">
                  After-sales service revenue
                </div>
              </div>
            </div>

            {/* After Sales Profit */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">After Sales Profit</h3>
              <div className="text-center">
                <div className="text-3xl font-bold text-white mb-2">
                  {dealer.afterSalesProfitsRap}
                </div>
                <div className="text-white text-sm">
                  After-sales profit value
                </div>
              </div>
            </div>

            {/* After Sales Margin */}
            <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
              <h3 className="text-lg font-semibold text-white mb-3">After Sales Margin</h3>
              <div className="text-center">
                <div className={`text-3xl font-bold mb-2 ${dealer.afterSalesMargin > 0 ? 'text-green-400' : 'text-red-400'}`}>
                  {dealer.afterSalesMargin}%
                </div>
                <div className="text-white text-sm">
                  After-sales margin percentage
                </div>
              </div>
            </div>
          </div>

          {/* Training & Contract Status */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Sales and AfterSales Trainings</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="text-center">
                <div className={`w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-3 ${
                  dealer.salesTrainings ? 'bg-green-500' : 'bg-red-500'
                }`}>
                  <svg className="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.246 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                  </svg>
                </div>
                <div className="text-white font-medium">Sales Trainings</div>
                <div className={`text-sm ${dealer.salesTrainings ? 'text-green-200' : 'text-red-200'}`}>
                  {dealer.salesTrainings ? 'Completed' : 'Not Completed'}
                </div>
              </div>
              
              <div className="text-center">
                <div className={`w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-3 ${
                  dealer.asTrainings ? 'bg-green-500' : 'bg-red-500'
                }`}>
                  <svg className="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
                  </svg>
                </div>
                <div className="text-white font-medium">AS Trainings</div>
                <div className={`text-sm ${dealer.asTrainings ? 'text-green-200' : 'text-red-200'}`}>
                  {dealer.asTrainings ? 'Completed' : 'Not Completed'}
                </div>
              </div>
            </div>
          </div>

          {/* Marketing Investment */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Marketing Investment</h3>
            <div className="text-center">
              <div className="text-3xl font-bold text-white mb-2">
                3,200,000 Rub
              </div>
              <div className="text-white text-sm">
                Marketing investment value
              </div>
            </div>
          </div>

        </div>
      </div>
    </div>
  )
}

export default DealerCard
