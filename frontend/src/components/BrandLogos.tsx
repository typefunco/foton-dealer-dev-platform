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
  // Проверяем, что brands не null и является массивом
  const safeBrands = brands || []
  
  // Определяем размер логотипов в зависимости от количества брендов
  const getLogoSize = (brandCount: number) => {
    if (brandCount <= 2) return 'w-28 h-28' // 112px для 1-2 брендов (было 80px)
    if (brandCount <= 4) return 'w-24 h-24' // 96px для 3-4 брендов (было 64px)
    if (brandCount <= 6) return 'w-20 h-20' // 80px для 5-6 брендов (было 56px)
    return 'w-18 h-18' // 72px для 7+ брендов (было 48px)
  }

  const logoSize = getLogoSize(safeBrands.length)

  // Список всех доступных брендов с их логотипами
  const allBrands: Record<string, BrandLogo> = {
    'FOTON': {
      name: 'FOTON',
      logo: '/brands/Foton.png',
      alt: 'FOTON Logo'
    },
    'DONGFENG': {
      name: 'DONGFENG',
      logo: '/brands/Dongfeng.png',
      alt: 'DONGFENG Logo'
    },
    'GAZ': {
      name: 'GAZ',
      logo: '/brands/GAZ.png',
      alt: 'GAZ Logo'
    },
    'KAMAZ': {
      name: 'KAMAZ',
      logo: '/brands/KAMAZ.png',
      alt: 'KAMAZ Logo'
    },
    'SHACMAN': {
      name: 'SHACMAN',
      logo: '/brands/Shacman.png',
      alt: 'SHACMAN Logo'
    },
    'ISUZU': {
      name: 'ISUZU',
      logo: '/brands/ISUZU.png',
      alt: 'ISUZU Logo'
    },
    'JAC': {
      name: 'JAC',
      logo: '/brands/Jac.png',
      alt: 'JAC Logo'
    },
    'MAZ': {
      name: 'MAZ',
      logo: '/brands/MAZ.png',
      alt: 'MAZ Logo'
    },
    'SANY': {
      name: 'SANY',
      logo: '/brands/SANY.png',
      alt: 'SANY Logo'
    },
    'SITRAK': {
      name: 'SITRAK',
      logo: '/brands/Sitrak.png',
      alt: 'SITRAK Logo'
    },
    'SOLLERS': {
      name: 'SOLLERS',
      logo: '/brands/SOLLERS.png',
      alt: 'SOLLERS Logo'
    },
    'VALDAI': {
      name: 'VALDAI',
      logo: '/brands/VALDAI.png',
      alt: 'VALDAI Logo'
    },
    'VALDAI-1': {
      name: 'VALDAI-1',
      logo: '/brands/VALDAI-1.png',
      alt: 'VALDAI-1 Logo'
    },
    'FAW': {
      name: 'FAW',
      logo: '/brands/Faw.png',
      alt: 'FAW Logo'
    },
    'CHENLONG': {
      name: 'CHENLONG',
      logo: '/brands/CHENLONG.png',
      alt: 'CHENLONG Logo'
    },
    'AMBERTRUCK': {
      name: 'AMBERTRUCK',
      logo: '/brands/AMBERTRUCK.png',
      alt: 'AMBERTRUCK Logo'
    }
  }

  // Фильтруем только те бренды, которые есть в портфеле дилера
  const availableBrands = safeBrands
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
          <div className={`${logoSize} bg-white bg-opacity-20 backdrop-blur-sm rounded-lg border border-white border-opacity-30 flex items-center justify-center p-2 hover:bg-opacity-30 transition-all duration-200`}>
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
            <span className="hidden text-xs font-bold text-white text-center leading-tight">
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
