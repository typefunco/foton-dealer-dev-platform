// React Hook для работы с данными After Sales
import { useState, useEffect, useCallback } from 'react';
import * as afterSalesApi from '../api/aftersales';
import type { AfterSalesDealer, AfterSalesFilters } from '../api/aftersales';

interface UseAfterSalesOptions {
  initialFilters?: AfterSalesFilters;
  autoLoad?: boolean;
}

export function useAfterSales(options: UseAfterSalesOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [dealers, setDealers] = useState<AfterSalesDealer[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<AfterSalesFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка данных After Sales
  const loadAfterSales = useCallback(async (customFilters?: AfterSalesFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await afterSalesApi.getAfterSalesData(filtersToUse);
      
      setDealers(response.dealers);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load after sales data');
      console.error('Failed to load after sales data:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<AfterSalesFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadAfterSales();
    }
  }, [autoLoad]);

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadAfterSales(filters);
    }
  }, [filters]);

  return {
    dealers,
    loading,
    error,
    filters,
    pagination,
    loadAfterSales,
    updateFilters,
    resetFilters,
  };
}
