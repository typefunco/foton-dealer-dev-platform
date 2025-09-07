import React from 'react'

interface BrandLogoProps {
  brand: string
  size?: 'sm' | 'md' | 'lg'
  className?: string
}

const BrandLogo: React.FC<BrandLogoProps> = ({ brand, size = 'sm', className = '' }) => {
  // Маппинг названий брендов на файлы
  const brandImageMap: { [key: string]: string } = {
    'FOTON': 'Foton.png',
    'DONGFENG': 'Dongfeng.png',
    'GAZ': 'GAZ.png',
    'KAMAZ': 'KAMAZ.png',
    'SHACMAN': 'Shacman.png',
    'ISUZU': 'ISUZU.png',
    'JAC': 'Jac.png',
    'MAZ': 'MAZ.png',
    'SANY': 'SANY.png',
    'SITRAK': 'Sitrak.png',
    'SOLLERS': 'SOLLERS.png',
    'VALDAI': 'VALDAI.png',
    'VALDAI-1': 'VALDAI-1.png',
    'FAW': 'Faw.png',
    'CHENLONG': 'CHENLONG.png',
    'AMBERTRUCK': 'AMBERTRUCK.png'
  }

  const sizeClasses = {
    sm: 'w-8 h-8',
    md: 'w-12 h-12',
    lg: 'w-16 h-16'
  }

  const imagePath = brandImageMap[brand.toUpperCase()]
  
  if (!imagePath) {
    // Fallback для неизвестных брендов
    return (
      <div className={`${sizeClasses[size]} bg-blue-400 bg-opacity-80 rounded-full flex items-center justify-center border border-blue-300 ${className}`}>
        <span className="text-xs font-bold text-white">{brand.charAt(0)}</span>
      </div>
    )
  }

  return (
    <div className={`${sizeClasses[size]} rounded-full overflow-hidden border border-blue-300 ${className}`}>
      <img
        src={`/brands/${imagePath}`}
        alt={brand}
        className="w-full h-full object-cover"
        title={brand}
      />
    </div>
  )
}

export default BrandLogo
