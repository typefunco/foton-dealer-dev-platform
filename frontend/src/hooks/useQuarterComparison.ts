// React Hook для работы с данными сравнения кварталов
import { useState, useEffect, useCallback } from 'react';
import * as quarterComparisonApi from '../api/quarter-comparison';
import type { QuarterComparisonData, QuarterComparisonFilters } from '../api/quarter-comparison';

interface UseQuarterComparisonOptions {
  initialFilters?: QuarterComparisonFilters;
  autoLoad?: boolean;
}

export function useQuarterComparison(options: UseQuarterComparisonOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [data, setData] = useState<QuarterComparisonData[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<QuarterComparisonFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка данных сравнения кварталов
  const loadQuarterComparison = useCallback(async (customFilters?: QuarterComparisonFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await quarterComparisonApi.getQuarterComparison(filtersToUse);
      
      setData(response.data);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load quarter comparison data');
      console.error('Failed to load quarter comparison data:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<QuarterComparisonFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadQuarterComparison();
    }
  }, [autoLoad]);

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadQuarterComparison(filters);
    }
  }, [filters]);

  return {
    data,
    loading,
    error,
    filters,
    pagination,
    loadQuarterComparison,
    updateFilters,
    resetFilters,
  };
}
