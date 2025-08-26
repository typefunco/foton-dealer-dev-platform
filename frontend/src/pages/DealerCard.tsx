import React, { useState, useEffect } from 'react'
import { useParams, Link } from 'react-router-dom'
import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, RadialBarChart, RadialBar } from 'recharts'

interface DealerData {
  id: string
  name: string
  city: string
  salesManager: string
  class: string
  checklist: number
  brandsInPortfolio: string[]
  salesTarget: string
  stockHdtMdtLdt: string
  buyoutHdtMdtLdt: string
  fotonSalesmen: string
  salesTrainings: boolean
  srRub: string
  salesProfit: number
  salesMargin: number
  autoSalesRevenue: string
  autoSalesProfitsRap: string
  autoSalesMargin: number
  ranking: number
  recommendedStock: number
  warrantyStock: number
  fotonLaborHours: number
  serviceContract: boolean
  asTrainings: boolean
  csi: number
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
      brandsInPortfolio: ['Foton', 'Honda', 'Toyota'],
      salesTarget: '1000',
      stockHdtMdtLdt: 'HDT: 50, MDT: 30, LDT: 20',
      buyoutHdtMdtLdt: 'HDT: 45, MDT: 25, LDT: 15',
      fotonSalesmen: '15',
      salesTrainings: true,
      srRub: '15,000,000',
      salesProfit: 25,
      salesMargin: 18,
      autoSalesRevenue: '8,000,000',
      autoSalesProfitsRap: '2,400,000',
      autoSalesMargin: 30,
      ranking: 1,
      recommendedStock: 85,
      warrantyStock: 70,
      fotonLaborHours: 92,
      serviceContract: true,
      asTrainings: true,
      csi: 95
    }

    setDealer(mockDealer)
  }, [dealerId])

  // Данные для pie charts
  const stockData = [
    { name: 'HDT', value: 50, color: '#3B82F6' },
    { name: 'MDT', value: 30, color: '#10B981' },
    { name: 'LDT', value: 20, color: '#8B5CF6' }
  ]

  const buyoutData = [
    { name: 'HDT', value: 45, color: '#EF4444' },
    { name: 'MDT', value: 25, color: '#F97316' },
    { name: 'LDT', value: 15, color: '#EAB308' }
  ]

  // Данные для performance metrics
  const performanceData = [
    { name: 'Recommended Stock', value: 85, fill: '#10B981' },
    { name: 'Warranty Stock', value: 70, fill: '#3B82F6' },
    { name: 'Foton Labor Hours', value: 92, fill: '#8B5CF6' }
  ]

  // Данные для sales performance
  const salesData = [
    { name: 'Sales Target', value: 1000, color: '#3B82F6' },
    { name: 'Sales Revenue', value: 15000, color: '#10B981' },
    { name: 'Auto Sales', value: 8000, color: '#8B5CF6' }
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
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <div className="text-blue-200 text-sm font-medium">Class</div>
            <div className="text-2xl font-bold text-white">{dealer.class}</div>
          </div>
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <div className="text-blue-200 text-sm font-medium">Checklist</div>
            <div className="text-2xl font-bold text-white">{dealer.checklist}</div>
          </div>
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <div className="text-blue-200 text-sm font-medium">Ranking</div>
            <div className="text-2xl font-bold text-white">#{dealer.ranking}</div>
          </div>
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <div className="text-blue-200 text-sm font-medium">CSI Score</div>
            <div className="text-2xl font-bold text-white">{dealer.csi}%</div>
          </div>
        </div>

        {/* Charts Section */}
        <div className="space-y-8">
          <h2 className="text-2xl font-bold text-white text-center mb-8">
            Comprehensive Performance Analytics
          </h2>
          
          {/* Sales Performance Bar Chart */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Sales Performance Overview</h3>
            <ResponsiveContainer width="100%" height={300}>
              <BarChart data={salesData} margin={{ top: 20, right: 30, left: 20, bottom: 5 }}>
                <CartesianGrid strokeDasharray="3 3" stroke="rgba(255,255,255,0.1)" />
                <XAxis 
                  dataKey="name" 
                  stroke="white"
                  fontSize={12}
                  tick={{ fill: 'white' }}
                />
                <YAxis 
                  stroke="white"
                  fontSize={12}
                  tick={{ fill: 'white' }}
                  tickFormatter={(value) => `${value / 1000}k`}
                />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: 'rgba(0,0,0,0.8)', 
                    border: 'none', 
                    borderRadius: '8px',
                    color: 'white'
                  }}
                />
                <Bar dataKey="value" radius={[4, 4, 0, 0]}>
                  {salesData.map((entry, index) => (
                    <Cell key={`cell-${index}`} fill={entry.color} />
                  ))}
                </Bar>
              </BarChart>
            </ResponsiveContainer>
          </div>

          {/* Stock & Inventory Pie Charts */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Stock & Inventory Distribution</h3>
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
              {/* Stock Distribution */}
              <div className="text-center">
                <h4 className="text-lg font-medium text-white mb-4">Stock Distribution</h4>
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
                        backgroundColor: 'rgba(0,0,0,0.8)', 
                        border: 'none', 
                        borderRadius: '8px',
                        color: 'white'
                      }}
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
                <h4 className="text-lg font-medium text-white mb-4">Buyout Distribution</h4>
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
                        backgroundColor: 'rgba(0,0,0,0.8)', 
                        border: 'none', 
                        borderRadius: '8px',
                        color: 'white'
                      }}
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
          </div>

          {/* Performance Metrics Radial Bar Chart */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Performance Metrics Radar</h3>
            <ResponsiveContainer width="100%" height={300}>
                             <RadialBarChart cx="50%" cy="50%" innerRadius="20%" outerRadius="80%" data={performanceData}>
                 <RadialBar 
                   label={{ fill: 'white', position: 'insideEnd', fontSize: 12 }}
                   background 
                   dataKey="value"
                 />
                <Tooltip 
                  contentStyle={{ 
                    backgroundColor: 'rgba(0,0,0,0.8)', 
                    border: 'none', 
                    borderRadius: '8px',
                    color: 'white'
                  }}
                />
              </RadialBarChart>
            </ResponsiveContainer>
          </div>

          {/* Training & Contract Status */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Training & Contract Status</h3>
            <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
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
              
              <div className="text-center">
                <div className={`w-20 h-20 mx-auto rounded-full flex items-center justify-center mb-3 ${
                  dealer.serviceContract ? 'bg-green-500' : 'bg-red-500'
                }`}>
                  <svg className="w-10 h-10 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </div>
                <div className="text-white font-medium">Service Contract</div>
                <div className={`text-sm ${dealer.serviceContract ? 'text-green-200' : 'text-red-200'}`}>
                  {dealer.serviceContract ? 'Active' : 'Inactive'}
                </div>
              </div>
            </div>
          </div>

          {/* Brands Portfolio Visualization */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Brands Portfolio</h3>
            <div className="flex justify-center space-x-6">
              {dealer.brandsInPortfolio.map((brand, index) => (
                <div key={index} className="text-center">
                  <div className="w-20 h-20 bg-gradient-to-br from-blue-400 to-blue-600 rounded-full flex items-center justify-center mb-3">
                    <span className="text-white font-bold text-lg">{brand.charAt(0)}</span>
                  </div>
                  <div className="text-white font-medium">{brand}</div>
                </div>
              ))}
            </div>
          </div>

          {/* CSI Score Pie Chart */}
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20">
            <h3 className="text-xl font-semibold text-white mb-4">Customer Satisfaction Index (CSI)</h3>
            <div className="flex justify-center">
              <ResponsiveContainer width={200} height={200}>
                <PieChart>
                  <Pie
                    data={[
                      { name: 'CSI Score', value: dealer.csi, color: '#10B981' },
                      { name: 'Remaining', value: 100 - dealer.csi, color: '#6B7280' }
                    ]}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    paddingAngle={0}
                    dataKey="value"
                  >
                    <Cell fill="#10B981" />
                    <Cell fill="#6B7280" />
                  </Pie>
                  <Tooltip 
                    contentStyle={{ 
                      backgroundColor: 'rgba(0,0,0,0.8)', 
                      border: 'none', 
                      borderRadius: '8px',
                      color: 'white'
                    }}
                  />
                </PieChart>
              </ResponsiveContainer>
            </div>
            <div className="text-center mt-4">
              <div className="text-white text-3xl font-bold">{dealer.csi}%</div>
              <div className="text-white text-sm">Excellent Performance</div>
            </div>
          </div>
        </div>

        {/* Detailed Data Table */}
        <div className="mt-12">
          <h2 className="text-2xl font-bold text-white text-center mb-8">Complete Performance Data</h2>
          <div className="bg-white bg-opacity-10 backdrop-blur-sm rounded-xl p-6 border border-white border-opacity-20 overflow-x-auto">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-white border-opacity-20">
                  <th className="text-left text-white font-semibold py-3">Metric</th>
                  <th className="text-center text-white font-semibold py-3">Value</th>
                  <th className="text-left text-white font-semibold py-3">Description</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-white divide-opacity-20">
                <tr>
                  <td className="text-blue-200 py-3">Sales Target</td>
                  <td className="text-center text-white py-3">{dealer.salesTarget}</td>
                  <td className="text-gray-300 py-3">Target sales volume for the period</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Sales Revenue (Rub)</td>
                  <td className="text-center text-white py-3">{dealer.srRub}</td>
                  <td className="text-gray-300 py-3">Total sales revenue in Russian Rubles</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Sales Profit %</td>
                  <td className="text-center text-white py-3">{dealer.salesProfit}%</td>
                  <td className="text-gray-300 py-3">Profit percentage from sales</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Sales Margin %</td>
                  <td className="text-center text-white py-3">{dealer.salesMargin}%</td>
                  <td className="text-gray-300 py-3">Sales margin percentage</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Auto Sales Revenue</td>
                  <td className="text-center text-white py-3">{dealer.autoSalesRevenue}</td>
                  <td className="text-gray-300 py-3">Revenue from automotive sales</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Auto Sales Profits</td>
                  <td className="text-center text-white py-3">{dealer.autoSalesProfitsRap}</td>
                  <td className="text-gray-300 py-3">Profits from automotive sales</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Auto Sales Margin %</td>
                  <td className="text-center text-white py-3">{dealer.autoSalesMargin}%</td>
                  <td className="text-gray-300 py-3">Margin percentage from auto sales</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Ranking</td>
                  <td className="text-center text-white py-3">#{dealer.ranking}</td>
                  <td className="text-gray-300 py-3">Current ranking position</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Recommended Stock %</td>
                  <td className="text-center text-white py-3">{dealer.recommendedStock}%</td>
                  <td className="text-gray-300 py-3">Recommended stock level percentage</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Warranty Stock %</td>
                  <td className="text-center text-white py-3">{dealer.warrantyStock}%</td>
                  <td className="text-gray-300 py-3">Warranty stock level percentage</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">Foton Labor Hours %</td>
                  <td className="text-center text-white py-3">{dealer.fotonLaborHours}%</td>
                  <td className="text-gray-300 py-3">Foton labor hours utilization</td>
                </tr>
                <tr>
                  <td className="text-blue-200 py-3">CSI Score</td>
                  <td className="text-center text-white py-3">{dealer.csi}%</td>
                  <td className="text-gray-300 py-3">Customer Satisfaction Index</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  )
}

export default DealerCard
