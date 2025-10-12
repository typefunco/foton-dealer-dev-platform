// React Hook для работы с пользователями
import { useState, useEffect, useCallback } from 'react';
import * as userApi from '../api/users';
import type { User, UserFilters, CreateUserRequest, UpdateUserRequest, UserCredentials } from '../api/users';

interface UseUsersOptions {
  initialFilters?: UserFilters;
  autoLoad?: boolean;
}

export function useUsers(options: UseUsersOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<UserFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка пользователей
  const loadUsers = useCallback(async (customFilters?: UserFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await userApi.getUsers(filtersToUse);
      
      setUsers(response.users);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load users');
      console.error('Failed to load users:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Создание пользователя
  const createUser = useCallback(async (data: CreateUserRequest): Promise<UserCredentials | null> => {
    try {
      setLoading(true);
      setError(null);
      
      const response = await userApi.createUser(data);
      
      // Перезагрузить список пользователей
      await loadUsers();
      
      return response.credentials || null;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to create user');
      console.error('Failed to create user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [loadUsers]);

  // Обновление пользователя
  const updateUser = useCallback(async (id: string, data: UpdateUserRequest) => {
    try {
      setLoading(true);
      setError(null);
      
      await userApi.updateUser(id, data);
      
      // Перезагрузить список пользователей
      await loadUsers();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update user');
      console.error('Failed to update user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [loadUsers]);

  // Удаление пользователя
  const deleteUser = useCallback(async (id: string) => {
    try {
      setLoading(true);
      setError(null);
      
      await userApi.deleteUser(id);
      
      // Перезагрузить список пользователей
      await loadUsers();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete user');
      console.error('Failed to delete user:', err);
      throw err;
    } finally {
      setLoading(false);
    }
  }, [loadUsers]);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<UserFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadUsers();
    }
  }, [autoLoad]); // Загружаем только один раз при монтировании

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadUsers(filters);
    }
  }, [filters]); // Зависимость только от filters

  return {
    users,
    loading,
    error,
    filters,
    pagination,
    loadUsers,
    createUser,
    updateUser,
    deleteUser,
    updateFilters,
    resetFilters,
  };
}

// Hook для получения статистики пользователей
export function useUserStats() {
  const [stats, setStats] = useState<userApi.UserStatsResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const loadStats = useCallback(async () => {
    try {
      setLoading(true);
      setError(null);
      
      const data = await userApi.getUserStats();
      setStats(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load stats');
      console.error('Failed to load stats:', err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    loadStats();
  }, [loadStats]);

  return {
    stats,
    loading,
    error,
    reload: loadStats,
  };
}

