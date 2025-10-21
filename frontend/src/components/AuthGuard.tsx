import React from 'react'
import { Navigate, useLocation } from 'react-router-dom'
import { useAuth } from '../contexts/AuthContext'

interface AuthGuardProps {
  children: React.ReactNode
  requireAdmin?: boolean
}

/**
 * Компонент для защиты маршрутов, требующих аутентификации
 */
export function AuthGuard({ children, requireAdmin = false }: AuthGuardProps) {
  const { isAuthenticated, isAdmin, isLoading, user } = useAuth()
  const location = useLocation()

  console.log('AuthGuard check:', { 
    isAuthenticated, 
    isAdmin, 
    isLoading, 
    user: user?.login,
    location: location.pathname,
    timestamp: new Date().toISOString()
  })

  // Показываем загрузку пока проверяем аутентификацию
  if (isLoading) {
    console.log('AuthGuard: showing loading')
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-white mx-auto mb-4"></div>
          <p className="text-white text-lg">Проверка аутентификации...</p>
        </div>
      </div>
    )
  }

  // Если пользователь не аутентифицирован, перенаправляем на логин
  if (!isAuthenticated || !user) {
    console.log('AuthGuard: redirecting to login', { isAuthenticated, user, currentPath: location.pathname })
    return <Navigate to="/login" state={{ from: { pathname: location.pathname } }} replace />
  }

  // Если требуется админ, но пользователь не админ
  if (requireAdmin && !isAdmin) {
    console.log('AuthGuard: admin access denied')
    return (
      <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700 flex items-center justify-center">
        <div className="text-center">
          <div className="text-red-400 text-6xl mb-4">🚫</div>
          <h1 className="text-white text-2xl font-bold mb-2">Доступ запрещен</h1>
          <p className="text-blue-200 mb-4">У вас нет прав для доступа к этой странице</p>
          <button
            onClick={() => window.history.back()}
            className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-lg transition-colors"
          >
            Назад
          </button>
        </div>
      </div>
    )
  }

  console.log('AuthGuard: access granted')
  return <>{children}</>
}

/**
 * Компонент для защиты админских маршрутов
 */
export function AdminGuard({ children }: { children: React.ReactNode }) {
  return <AuthGuard requireAdmin={true}>{children}</AuthGuard>
}
