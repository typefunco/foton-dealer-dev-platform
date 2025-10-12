// React Hook для работы с данными производительности
import { useState, useEffect, useCallback } from 'react';
import * as performanceApi from '../api/performance';
import type { PerformanceDealer, PerformanceFilters } from '../api/performance';

interface UsePerformanceOptions {
  initialFilters?: PerformanceFilters;
  autoLoad?: boolean;
}

export function usePerformance(options: UsePerformanceOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [dealers, setDealers] = useState<PerformanceDealer[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<PerformanceFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка данных производительности
  const loadPerformance = useCallback(async (customFilters?: PerformanceFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await performanceApi.getPerformanceData(filtersToUse);
      
      setDealers(response.dealers);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load performance data');
      console.error('Failed to load performance data:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<PerformanceFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadPerformance();
    }
  }, [autoLoad]);

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadPerformance(filters);
    }
  }, [filters]);

  return {
    dealers,
    loading,
    error,
    filters,
    pagination,
    loadPerformance,
    updateFilters,
    resetFilters,
  };
}
