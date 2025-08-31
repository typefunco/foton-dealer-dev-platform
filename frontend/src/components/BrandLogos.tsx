import React from 'react'

interface BrandLogo {
  name: string
  logo: string
  alt: string
}

interface BrandLogosProps {
  brands: string[]
  className?: string
}

const BrandLogos: React.FC<BrandLogosProps> = ({ brands, className = '' }) => {
  // Список всех доступных брендов с их логотипами
  const allBrands: Record<string, BrandLogo> = {
    'Shacman': {
      name: 'Shacman',
      logo: '/brands/shacman.svg',
      alt: 'Shacman Logo'
    },
    'Dongfeng': {
      name: 'Dongfeng',
      logo: '/brands/dongfeng.svg',
      alt: 'Dongfeng Logo'
    },
    'Sitrak': {
      name: 'Sitrak',
      logo: '/brands/sitrak.svg',
      alt: 'Sitrak Logo'
    },
    'FAW': {
      name: 'FAW',
      logo: '/brands/faw.svg',
      alt: 'FAW Logo'
    },
    'JAC': {
      name: 'JAC',
      logo: '/brands/jac.svg',
      alt: 'JAC Logo'
    },
    'Isuzu': {
      name: 'Isuzu',
      logo: '/brands/isuzu.svg',
      alt: 'Isuzu Logo'
    },
    'Kamaz': {
      name: 'Kamaz',
      logo: '/brands/kamaz.svg',
      alt: 'Kamaz Logo'
    },
    'MAZ': {
      name: 'MAZ',
      logo: '/brands/maz.svg',
      alt: 'MAZ Logo'
    },
    'GAZ': {
      name: 'GAZ',
      logo: '/brands/gaz.svg',
      alt: 'GAZ Logo'
    },
    'Valdai': {
      name: 'Valdai',
      logo: '/brands/valdai.svg',
      alt: 'Valdai Logo'
    },
    'Foton': {
      name: 'Foton',
      logo: '/brands/foton.svg',
      alt: 'Foton Logo'
    },
    'SDAC': {
      name: 'SDAC',
      logo: '/brands/sdac.svg',
      alt: 'SDAC Logo'
    },
    'Ambertruck': {
      name: 'Ambertruck',
      logo: '/brands/ambertruck.svg',
      alt: 'Ambertruck Logo'
    },
    'Sollers': {
      name: 'Sollers',
      logo: '/brands/sollers.svg',
      alt: 'Sollers Logo'
    },
    'Chenglong': {
      name: 'Chenglong',
      logo: '/brands/chenglong.svg',
      alt: 'Chenglong Logo'
    },
    'Sany': {
      name: 'Sany',
      logo: '/brands/sany.svg',
      alt: 'Sany Logo'
    }
  }

  // Фильтруем только те бренды, которые есть в портфеле дилера
  const availableBrands = brands
    .filter(brand => allBrands[brand])
    .map(brand => allBrands[brand])

  if (availableBrands.length === 0) {
    return (
      <div className={`text-gray-400 text-sm ${className}`}>
        No brand information available
      </div>
    )
  }

  return (
    <div className={`flex flex-wrap gap-3 items-center ${className}`}>
      {availableBrands.map((brand) => (
        <div
          key={brand.name}
          className="relative group cursor-pointer"
          title={brand.name}
        >
          <div className="w-12 h-12 bg-white rounded-lg shadow-md flex items-center justify-center p-2 hover:shadow-lg transition-shadow duration-200">
            <img
              src={brand.logo}
              alt={brand.alt}
              className="w-full h-full object-contain"
              onError={(e) => {
                // Fallback для отсутствующих логотипов
                const target = e.target as HTMLImageElement
                target.style.display = 'none'
                target.nextElementSibling?.classList.remove('hidden')
              }}
            />
            {/* Fallback текст если логотип не загрузился */}
            <span className="hidden text-xs font-bold text-gray-600 text-center leading-tight">
              {brand.name}
            </span>
          </div>
          
          {/* Tooltip */}
          <div className="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-1 bg-black text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
            {brand.name}
            <div className="absolute top-full left-1/2 transform -translate-x-1/2 w-0 h-0 border-l-4 border-r-4 border-t-4 border-transparent border-t-black"></div>
          </div>
        </div>
      ))}
    </div>
  )
}

export default BrandLogos
