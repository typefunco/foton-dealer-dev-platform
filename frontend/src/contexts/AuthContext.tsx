import { createContext, useContext, useState, useEffect, useCallback, ReactNode } from 'react'
import { login as apiLogin, logout as apiLogout, isAuthenticated, getUserFromToken, User } from '../api/auth'

interface AuthContextType {
  user: User | null
  isAuthenticated: boolean
  isAdmin: boolean
  isLoading: boolean
  login: (credentials: { login: string; password: string }) => Promise<User>
  logout: () => Promise<void>
  error: string | null
  clearError: () => void
}

const AuthContext = createContext<AuthContextType | undefined>(undefined)

interface AuthProviderProps {
  children: ReactNode
}

/**
 * Провайдер контекста аутентификации
 */
export function AuthProvider({ children }: AuthProviderProps) {
  const [user, setUser] = useState<User | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  // Проверяем аутентификацию при загрузке
  useEffect(() => {
    const checkAuth = () => {
      try {
        console.log('AuthContext: checking authentication...')
        
        // Небольшая задержка для того, чтобы cookie успел загрузиться
        setTimeout(() => {
          const hasToken = isAuthenticated()
          console.log('AuthContext: hasToken:', hasToken)
          
          if (hasToken) {
            const userData = getUserFromToken()
            console.log('AuthContext: userData from token:', userData)
            setUser(userData)
          } else {
            console.log('AuthContext: no token found')
            setUser(null)
          }
          setIsLoading(false)
        }, 100)
      } catch (err) {
        console.error('AuthContext: auth check failed:', err)
        setUser(null)
        setIsLoading(false)
      }
    }

    checkAuth()
  }, [])

  // Функция логина
  const handleLogin = useCallback(async (credentials: { login: string; password: string }): Promise<User> => {
    try {
      setIsLoading(true)
      setError(null)
      
      console.log('Starting login process...')
      const response = await apiLogin(credentials)
      console.log('Login response:', response)
      
      // Обновляем состояние пользователя
      setUser(response.user)
      console.log('User state updated:', response.user)
      
      // Небольшая задержка перед сбросом loading
      setTimeout(() => {
        setIsLoading(false)
      }, 100)
      
      // Возвращаем данные пользователя для использования в компоненте
      return response.user
    } catch (err) {
      const errorMessage = err instanceof Error ? err.message : 'Login failed'
      setError(errorMessage)
      console.error('Login error:', err)
      setIsLoading(false)
      throw err
    }
  }, [])

  // Функция логаута
  const handleLogout = useCallback(async () => {
    try {
      await apiLogout()
      setUser(null)
    } catch (err) {
      console.error('Logout failed:', err)
      // Даже если логаут на сервере не удался, очищаем локальное состояние
      setUser(null)
    }
  }, [])

  // Очистка ошибки
  const clearError = useCallback(() => {
    setError(null)
  }, [])

  const value: AuthContextType = {
    user,
    isAuthenticated: !!user,
    isAdmin: user?.is_admin || false,
    isLoading,
    login: handleLogin,
    logout: handleLogout,
    error,
    clearError,
  }

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  )
}

/**
 * Хук для использования контекста аутентификации
 */
export function useAuth(): AuthContextType {
  const context = useContext(AuthContext)
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider')
  }
  return context
}
