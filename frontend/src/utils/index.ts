// Utilities for the application

export const formatCurrency = (amount: number, currency = 'RUB'): string => {
  return new Intl.NumberFormat('ru-RU', {
    style: 'currency',
    currency,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0,
  }).format(amount)
}

export const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat('ru-RU', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }).format(date)
}

export const formatRelativeTime = (date: Date): string => {
  const now = new Date()
  const diffInSeconds = Math.floor((now.getTime() - date.getTime()) / 1000)
  
  if (diffInSeconds < 60) {
    return 'just now'
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60)
          return `${minutes} ${minutes === 1 ? 'minute' : 'minutes'} ago`
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600)
          return `${hours} ${hours === 1 ? 'hour' : 'hours'} ago`
  } else {
    const days = Math.floor(diffInSeconds / 86400)
          return `${days} ${days === 1 ? 'day' : 'days'} ago`
  }
}

export const classNames = (...classes: (string | boolean | undefined)[]): string => {
  return classes.filter(Boolean).join(' ')
}

export const generateId = (): string => {
  return Math.random().toString(36).substr(2, 9)
}

/**
 * Маппинг названий бизнесов с бэкенда к именам файлов иконок
 * Преобразует различные варианты названий бизнесов в корректные пути к иконкам
 */
const businessIconMap: Record<string, string> = {
  // Construction
  'construction': 'construction.png',
  
  // Timber
  'timber': 'Timber.png',
  'wood': 'Timber.png',
  
  // Development
  'development': 'development.png',
  
  // Real Estate
  'real estate': 'real_estate.png',
  'real_estate': 'real_estate.png',
  'realestate': 'real_estate.png',
  
  // Spare Parts
  'spare parts': 'Spare_parts.png',
  'spare_parts': 'Spare_parts.png',
  'spareparts': 'Spare_parts.png',
  
  // Transport & Logistics
  'transport & logistics': 'transport_company.png',
  'transport_and_logistics': 'transport_company.png',
  'transport': 'transport_company.png',
  'logistics': 'transport_company.png',
  'transport &': 'transport_company.png',
  '& logistics': 'transport_company.png',
  
  // Car Salons
  'car salons': 'Car salons.png',
  'car_salons': 'Car salons.png',
  
  // Rental
  'rental': 'Rental.png',
  
  // Industrial Plant
  'industrial plant': 'manufacture.png',
  'industrial_plant': 'manufacture.png',
  'industrial': 'manufacture.png',
  'plant': 'manufacture.png',
  
  // Mining
  'mining': 'mining.png',
  
  // Auto Accessories
  'auto accessories': 'Spare_parts.png',
  'auto_accessories': 'Spare_parts.png',
  
  // Bodybuilding
  'bodybuilding': 'bodubuilding.png',
  'body building': 'bodubuilding.png',
  'body_building': 'bodubuilding.png',
  
  // Agriculture Machines
  'agriculture machines': 'agriculture.png',
  'agriculture_machines': 'agriculture.png',
  'agriculture': 'agriculture.png',
  'agricultural': 'agriculture.png',
  
  // Construction Tools
  'construction tools': 'Construction tools.png',
  'construction_tools': 'Construction tools.png',
  
  // Fuel Stations
  'fuel stations': 'fuel station.png',
  'fuel_stations': 'fuel station.png',
  'fuel station': 'fuel station.png',
  'gas station': 'fuel station.png',
  'gas stations': 'fuel station.png',
  
  // Garbage
  'garbage': 'Garbage.png',
  'waste': 'Garbage.png',
  'waste management': 'Garbage.png',
  
  // Bank & Leasing
  'bank&leasing': 'bank_finance.png',
  'bank & leasing': 'bank_finance.png',
  'bank_and_leasing': 'bank_finance.png',
  'bank': 'bank_finance.png',
  'leasing': 'bank_finance.png',
  'finance': 'bank_finance.png',
  
  // Coolant Equipment (используем manufacture.png для промышленного оборудования)
  'coolant equipment': 'manufacture.png',
  'coolant_equipment': 'manufacture.png',
  'coolant': 'manufacture.png',
  'equipment': 'manufacture.png',
  
  // Express Delivery
  'express delivery': 'transport_company.png',
  'express_delivery': 'transport_company.png',
  'express': 'transport_company.png',
  'delivery': 'transport_company.png',
  
  // Ecom
  'ecom': 'ecom.png',
  'e-commerce': 'ecom.png',
  'ecommerce': 'ecom.png',
}

/**
 * Конвертирует название бизнеса с бэкенда в путь к иконке
 * @param businessName - Название бизнеса с бэкенда (может быть в разных форматах)
 * @returns Путь к иконке или null, если иконка не найдена
 */
export const getBusinessIconPath = (businessName: string): string | null => {
  if (!businessName) return null
  
  // Нормализуем название: убираем пробелы по краям, приводим к нижнему регистру
  // Также нормализуем символ & для совместимости
  let normalized = businessName.trim().toLowerCase()
  
  // Нормализуем различные варианты написания символа &
  normalized = normalized.replace(/&/g, '&').replace(/and/g, '&')
  
  // Прямой поиск в маппинге
  if (businessIconMap[normalized]) {
    return `/icons/${businessIconMap[normalized]}`
  }
  
  // Специальная обработка для различения "development" и "real estate"
  // Приоритет: сначала проверяем "real estate" (более специфичное название)
  if (normalized.includes('real') && normalized.includes('estate')) {
    // Это точно "real estate" - возвращаем соответствующую иконку
    const iconPath = `/icons/${businessIconMap['real estate'] || businessIconMap['real_estate'] || businessIconMap['realestate'] || 'real_estate.png'}`
    console.log(`[getBusinessIconPath] "real estate" matched for "${businessName}" -> ${iconPath}`)
    return iconPath
  }
  
  // Если это точно "development" (без слов "real" и "estate")
  if (normalized === 'development' || (normalized.includes('development') && !normalized.includes('real') && !normalized.includes('estate'))) {
    const iconPath = `/icons/${businessIconMap['development'] || 'development.png'}`
    console.log(`[getBusinessIconPath] "development" matched for "${businessName}" -> ${iconPath}`)
    return iconPath
  }
  
  // Поиск по частичному совпадению (для остальных случаев)
  // Сначала ищем самые длинные совпадения (более точные)
  const sortedKeys = Object.keys(businessIconMap).sort((a, b) => b.length - a.length)
  
  for (const key of sortedKeys) {
    // Пропускаем "development" и "real estate", так как они уже обработаны выше
    if (key === 'development' || key === 'real estate' || key === 'real_estate' || key === 'realestate') {
      continue
    }
    
    // Проверяем точное совпадение или включение
    if (normalized === key || normalized.includes(key) || key.includes(normalized)) {
      return `/icons/${businessIconMap[key]}`
    }
    
    // Проверка по словам для многословных названий (например, "Transport & Logistics")
    const normalizedWords = normalized.split(/\s+/)
    const keyWords = key.split(/\s+/)
    if (normalizedWords.some(word => keyWords.some(kw => kw.includes(word) || word.includes(kw)))) {
      return `/icons/${businessIconMap[key]}`
    }
  }
  
  // Если не найдено, возвращаем null (можно будет использовать дефолтную иконку)
  return null
}

/**
 * Получает иконки для массива названий бизнесов
 * @param businesses - Массив названий бизнесов
 * @returns Массив объектов с названием бизнеса и путем к иконке
 */
export const getBusinessIcons = (businesses: string[]): Array<{ name: string; iconPath: string | null }> => {
  if (!businesses || businesses.length === 0) return []
  
  return businesses.map(business => ({
    name: business,
    iconPath: getBusinessIconPath(business)
  }))
}

/**
 * Транслитерация русских букв в английские
 * Преобразует русские названия брендов (МАЗ, ГАЗ, КАМАЗ, ЛИАЗ, КАВЗ, ПАЗ) в английские эквиваленты
 */
const transliterateRussian = (text: string): string => {
  const russianToLatin: Record<string, string> = {
    'А': 'A', 'Б': 'B', 'В': 'V', 'Г': 'G', 'Д': 'D',
    'Е': 'E', 'Ё': 'YO', 'Ж': 'ZH', 'З': 'Z', 'И': 'I',
    'Й': 'Y', 'К': 'K', 'Л': 'L', 'М': 'M', 'Н': 'N',
    'О': 'O', 'П': 'P', 'Р': 'R', 'С': 'S', 'Т': 'T',
    'У': 'U', 'Ф': 'F', 'Х': 'KH', 'Ц': 'TS', 'Ч': 'CH',
    'Ш': 'SH', 'Щ': 'SCH', 'Ъ': '', 'Ы': 'Y', 'Ь': '',
    'Э': 'E', 'Ю': 'YU', 'Я': 'YA',
    'а': 'a', 'б': 'b', 'в': 'v', 'г': 'g', 'д': 'd',
    'е': 'e', 'ё': 'yo', 'ж': 'zh', 'з': 'z', 'и': 'i',
    'й': 'y', 'к': 'k', 'л': 'l', 'м': 'm', 'н': 'n',
    'о': 'o', 'п': 'p', 'р': 'r', 'с': 's', 'т': 't',
    'у': 'u', 'ф': 'f', 'х': 'kh', 'ц': 'ts', 'ч': 'ch',
    'ш': 'sh', 'щ': 'sch', 'ъ': '', 'ы': 'y', 'ь': '',
    'э': 'e', 'ю': 'yu', 'я': 'ya',
  }

  return text
    .split('')
    .map(char => russianToLatin[char] || char)
    .join('')
}

/**
 * Маппинг названий брендов с бэкенда к именам файлов иконок
 * Преобразует различные варианты названий брендов в корректные пути к иконкам
 */
const brandIconMap: Record<string, string> = {
  // Foton
  'foton': 'Foton.png',
  
  // KAMAZ (может прийти как КАМАЗ на русском)
  'kamaz': 'KAMAZ.png',
  
  // Sitrak
  'sitrak': 'Sitrak.png',
  
  // GAZ (может прийти как ГАЗ на русском)
  'gaz': 'GAZ.png',
  
  // MAZ (может прийти как МАЗ на русском)
  'maz': 'MAZ.png',
  
  // LIAZ (может прийти как ЛИАЗ на русском)
  'liaz': 'liaz.png',
  
  // KAVZ (может прийти как КАВЗ на русском)
  'kavz': 'kavz.png',
  
  // PAZ (может прийти как ПАЗ на русском)
  'paz': 'paz.png',
  
  // Liugong (на английском)
  'liugong': 'liugong.png',
  'lugong': 'liugong.png',
  
  // Dongfeng
  'dongfeng': 'Dongfeng.png',
  
  // Jac
  'jac': 'Jac.png',
  
  // Shacman
  'shacman': 'Shacman.png',
  
  // SANY
  'sany': 'SANY.png',
  
  // Sollers
  'sollers': 'SOLLERS.png',
  
  // VALDAI
  'valdai': 'VALDAI.png',
  'valdai-1': 'VALDAI-1.png',
  
  // ISUZU
  'isuzu': 'ISUZU.png',
  
  // FAW
  'faw': 'Faw.png',
  
  // CHENLONG
  'chenlong': 'CHENLONG.png',
  
  // AMBERTRUCK
  'ambertruck': 'AMBERTRUCK.png',
  
  // Новые бренды
  'forland': 'forland.png',
  'howo': 'howo.png',
  'lovol': 'lovol.png',
  'sunward': 'sunward.png',
}

/**
 * Получает путь к иконке бренда по его названию
 * @param brandName - Название бренда с бэкенда
 * @returns Путь к иконке или null, если иконка не найдена
 */
export const getBrandIconPath = (brandName: string): string | null => {
  if (!brandName) return null
  
  // Нормализуем название: убираем пробелы по краям, приводим к нижнему регистру
  let normalized = brandName.trim().toLowerCase()
  
  // Проверяем, есть ли русские буквы в названии
  const hasRussianLetters = /[а-яёА-ЯЁ]/.test(normalized)
  
  // Если есть русские буквы, транслитерируем их в английские
  if (hasRussianLetters) {
    const originalName = normalized
    const transliterated = transliterateRussian(normalized).toLowerCase()
    normalized = transliterated
    console.log(`[getBrandIconPath] Transliterated "${brandName}" (${originalName}) -> "${normalized}"`)
  }
  
  // Прямой поиск в маппинге
  if (brandIconMap[normalized]) {
    return `/brands/${brandIconMap[normalized]}`
  }
  
  // Поиск по частичному совпадению
  const sortedKeys = Object.keys(brandIconMap).sort((a, b) => b.length - a.length)
  
  for (const key of sortedKeys) {
    if (normalized === key || normalized.includes(key) || key.includes(normalized)) {
      return `/brands/${brandIconMap[key]}`
    }
  }
  
  return null
}

/**
 * Получает иконки для массива названий брендов
 * @param brands - Массив названий брендов
 * @returns Массив объектов с названием бренда и путем к иконке
 */
export const getBrandIcons = (brands: string[]): Array<{ name: string; iconPath: string | null }> => {
  if (!brands || brands.length === 0) return []
  
  return brands.map(brand => ({
    name: brand,
    iconPath: getBrandIconPath(brand)
  }))
}
