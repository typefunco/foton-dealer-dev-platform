// React Hook для работы с данными команды продаж
import { useState, useEffect, useCallback } from 'react';
import * as salesTeamApi from '../api/sales-team';
import type { SalesTeamMember, SalesTeamFilters } from '../api/sales-team';

interface UseSalesTeamOptions {
  initialFilters?: SalesTeamFilters;
  autoLoad?: boolean;
}

export function useSalesTeam(options: UseSalesTeamOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [members, setMembers] = useState<SalesTeamMember[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<SalesTeamFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка данных команды продаж
  const loadSalesTeam = useCallback(async (customFilters?: SalesTeamFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await salesTeamApi.getSalesTeamData(filtersToUse);
      
      setMembers(response.members);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load sales team data');
      console.error('Failed to load sales team data:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<SalesTeamFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadSalesTeam();
    }
  }, [autoLoad]);

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadSalesTeam(filters);
    }
  }, [filters]);

  return {
    members,
    loading,
    error,
    filters,
    pagination,
    loadSalesTeam,
    updateFilters,
    resetFilters,
  };
}
