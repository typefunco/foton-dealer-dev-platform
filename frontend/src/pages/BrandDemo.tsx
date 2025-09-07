import React from 'react'
import { Link } from 'react-router-dom'
import BrandLogos from '../components/BrandLogos'

interface DemoDealer {
  id: string
  name: string
  city: string
  brandsInPortfolio: string[]
  brandsCount: number
}

const BrandDemo: React.FC = () => {
  const demoDealers: DemoDealer[] = [
    {
      id: '1',
      name: 'Dealer 1',
      city: 'Moscow',
      brandsInPortfolio: ['FOTON', 'DONGFENG', 'GAZ', 'KAMAZ', 'SHACMAN'],
      brandsCount: 5
    },
    {
      id: '2',
      name: 'Dealer 2',
      city: 'St. Petersburg',
      brandsInPortfolio: ['FOTON', 'FAW'],
      brandsCount: 2
    },
    {
      id: '3',
      name: 'Dealer 3',
      city: 'Novosibirsk',
      brandsInPortfolio: ['FOTON', 'SANY', 'SITRAK', 'SOLLERS', 'VALDAI', 'ISUZU', 'CHENLONG', 'AMBERTRUCK'],
      brandsCount: 8
    }
  ]

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
          <h1 className="text-5xl md:text-6xl font-bold text-white mb-4">
            BRAND LOGOS
          </h1>
          <h2 className="text-3xl md:text-4xl font-bold text-blue-200">
            ADAPTIVE SIZES DEMO
          </h2>
        </div>
      </div>

      {/* Demo Cards */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pb-8">
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {demoDealers.map((dealer) => (
            <div
              key={dealer.id}
              className="bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-8 border border-white border-opacity-20 hover:bg-opacity-20 transition-all duration-300"
            >
              {/* Dealer Info */}
              <div className="text-center mb-6">
                <h3 className="text-2xl font-bold text-white mb-2">{dealer.name}</h3>
                <p className="text-blue-200 mb-2">{dealer.city}</p>
                <p className="text-sm text-gray-300">
                  {dealer.brandsCount} brand{dealer.brandsCount !== 1 ? 's' : ''} in portfolio
                </p>
              </div>

              {/* Brand Logos */}
              <div className="mb-6">
                <h4 className="text-lg font-semibold text-white mb-4 text-center">
                  Brands Portfolio
                </h4>
                <BrandLogos brands={dealer.brandsInPortfolio} className="justify-center" />
              </div>

              {/* Size Info */}
              <div className="text-center">
                <div className="bg-blue-800 bg-opacity-30 rounded-lg p-4">
                  <p className="text-sm text-blue-200">
                    Logo size: {
                      dealer.brandsCount <= 2 ? '112px (Large)' :
                      dealer.brandsCount <= 4 ? '96px (Medium)' :
                      dealer.brandsCount <= 6 ? '80px (Compact)' :
                      '72px (Small)'
                    }
                  </p>
                </div>
              </div>
            </div>
          ))}
        </div>

        {/* Explanation */}
        <div className="mt-12 bg-white bg-opacity-10 backdrop-blur-sm rounded-2xl p-8 border border-white border-opacity-20">
          <h3 className="text-2xl font-bold text-white mb-4 text-center">
            How It Works
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="bg-blue-600 bg-opacity-50 rounded-lg p-4 mb-3">
                <span className="text-2xl font-bold text-white">1-2</span>
              </div>
              <p className="text-sm text-blue-200">Brands</p>
              <p className="text-xs text-gray-300">112px (Large)</p>
            </div>
            <div className="text-center">
              <div className="bg-blue-600 bg-opacity-50 rounded-lg p-4 mb-3">
                <span className="text-2xl font-bold text-white">3-4</span>
              </div>
              <p className="text-sm text-blue-200">Brands</p>
              <p className="text-xs text-gray-300">96px (Medium)</p>
            </div>
            <div className="text-center">
              <div className="bg-blue-600 bg-opacity-50 rounded-lg p-4 mb-3">
                <span className="text-2xl font-bold text-white">5-6</span>
              </div>
              <p className="text-sm text-blue-200">Brands</p>
              <p className="text-xs text-gray-300">80px (Compact)</p>
            </div>
            <div className="text-center">
              <div className="bg-blue-600 bg-opacity-50 rounded-lg p-4 mb-3">
                <span className="text-2xl font-bold text-white">7+</span>
              </div>
              <p className="text-sm text-blue-200">Brands</p>
              <p className="text-xs text-gray-300">72px (Small)</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

export default BrandDemo
