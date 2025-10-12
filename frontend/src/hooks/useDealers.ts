// React Hook для работы с дилерами
import { useState, useEffect, useCallback } from 'react';
import * as dealerApi from '../api/dealers';
import type { Dealer, DealerCard, DealerFilters } from '../api/dealers';

// Расширяем интерфейс фильтров для поддержки новых параметров
interface ExtendedDealerFilters extends DealerFilters {
  quarter?: string;
  year?: number;
}

interface UseDealersOptions {
  initialFilters?: ExtendedDealerFilters;
  autoLoad?: boolean;
}

export function useDealers(options: UseDealersOptions = {}) {
  const { initialFilters = {}, autoLoad = true } = options;

  const [dealers, setDealers] = useState<Dealer[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<ExtendedDealerFilters>(initialFilters);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0,
  });

  // Загрузка дилеров
  const loadDealers = useCallback(async (customFilters?: ExtendedDealerFilters) => {
    try {
      setLoading(true);
      setError(null);
      
      const filtersToUse = customFilters || filters;
      const response = await dealerApi.getDealers(filtersToUse);
      
      setDealers(response.dealers);
      setPagination(response.pagination);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dealers');
      console.error('Failed to load dealers:', err);
    } finally {
      setLoading(false);
    }
  }, [filters]);

  // Получение дилера по ID
  const getDealerById = useCallback(async (id: string): Promise<Dealer | null> => {
    try {
      setLoading(true);
      setError(null);
      
      const dealer = await dealerApi.getDealerById(id);
      return dealer;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dealer');
      console.error('Failed to load dealer:', err);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Получение карточки дилера
  const getDealerCard = useCallback(async (id: string): Promise<DealerCard | null> => {
    try {
      setLoading(true);
      setError(null);
      
      const card = await dealerApi.getDealerCard(id);
      return card;
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load dealer card');
      console.error('Failed to load dealer card:', err);
      return null;
    } finally {
      setLoading(false);
    }
  }, []);

  // Обновление фильтров
  const updateFilters = useCallback((newFilters: Partial<ExtendedDealerFilters>) => {
    setFilters(prev => ({ ...prev, ...newFilters }));
  }, []);

  // Сброс фильтров
  const resetFilters = useCallback(() => {
    setFilters({});
  }, []);

  // Автоматическая загрузка при монтировании
  useEffect(() => {
    if (autoLoad) {
      loadDealers();
    }
  }, [autoLoad]);

  // Перезагрузка при изменении фильтров
  useEffect(() => {
    if (autoLoad && Object.keys(filters).length > 0) {
      loadDealers(filters);
    }
  }, [filters]);

  return {
    dealers,
    loading,
    error,
    filters,
    pagination,
    loadDealers,
    getDealerById,
    getDealerCard,
    updateFilters,
    resetFilters,
  };
}
