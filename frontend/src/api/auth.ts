// API клиент для аутентификации

export interface LoginRequest {
  login: string
  password: string
}

export interface LoginResponse {
  token: string
  user: {
    login: string
    is_admin: boolean
    role: string
  }
}

export interface User {
  login: string
  is_admin: boolean
  role: string
}

const API_BASE_URL = import.meta.env.DEV ? '' : 'http://localhost:8080'

/**
 * Выполняет логин пользователя
 */
export async function login(credentials: LoginRequest): Promise<LoginResponse> {
  const response = await fetch(`${API_BASE_URL}/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    credentials: 'include', // Включаем cookies
    body: JSON.stringify(credentials),
  })

  if (!response.ok) {
    const error = await response.json()
    throw new Error(error.error || 'Login failed')
  }

  const result = await response.json()
  
  // Сохраняем токен в localStorage
  if (result.token) {
    localStorage.setItem('auth_token', result.token)
    console.log('Token saved to localStorage')
  }
  
  return result
}

/**
 * Выполняет логаут пользователя
 */
export async function logout(): Promise<void> {
  // Удаляем токен из localStorage
  localStorage.removeItem('auth_token')
  console.log('Token removed from localStorage')
  
  // Удаляем cookie на фронтенде (если есть)
  document.cookie = 'auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;'
}

/**
 * Проверяет, аутентифицирован ли пользователь
 */
export function isAuthenticated(): boolean {
  // Сначала проверяем localStorage
  const localStorageToken = localStorage.getItem('auth_token')
  console.log('isAuthenticated: localStorage token:', localStorageToken ? 'exists' : 'null')
  
  if (localStorageToken) {
    return true
  }
  
  // Fallback на cookie (для совместимости)
  const cookieToken = getCookie('auth_token')
  console.log('isAuthenticated: cookie token:', cookieToken ? 'exists' : 'null')
  
  return cookieToken !== null
}

/**
 * Получает токен из localStorage или cookie
 */
export function getAuthToken(): string | null {
  // Сначала проверяем localStorage
  const localStorageToken = localStorage.getItem('auth_token')
  if (localStorageToken) {
    return localStorageToken
  }
  
  // Fallback на cookie
  return getCookie('auth_token')
}

/**
 * Получает cookie по имени
 */
function getCookie(name: string): string | null {
  console.log(`getCookie: looking for cookie "${name}"`)
  const value = `; ${document.cookie}`
  console.log(`getCookie: document.cookie = "${document.cookie}"`)
  const parts = value.split(`; ${name}=`)
  console.log(`getCookie: parts =`, parts)
  if (parts.length === 2) {
    const result = parts.pop()?.split(';').shift() || null
    console.log(`getCookie: found cookie "${name}" = "${result}"`)
    return result
  }
  console.log(`getCookie: cookie "${name}" not found`)
  return null
}

/**
 * Проверяет, является ли пользователь администратором
 */
export function isAdmin(): boolean {
  const token = getAuthToken()
  if (!token) return false
  
  try {
    // Декодируем JWT токен (только для проверки роли)
    const payload = JSON.parse(atob(token.split('.')[1]))
    return payload.is_admin === true
  } catch {
    return false
  }
}

/**
 * Получает информацию о пользователе из токена
 */
export function getUserFromToken(): User | null {
  const token = getAuthToken()
  if (!token) return null
  
  try {
    const payload = JSON.parse(atob(token.split('.')[1]))
    return {
      login: payload.login,
      is_admin: payload.is_admin,
      role: payload.role,
    }
  } catch {
    return null
  }
}
