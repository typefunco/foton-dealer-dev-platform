import { useState, useCallback } from 'react';
import { 
  uploadExcelFile, 
  getExcelTables, 
  getExcelTableData, 
  deleteExcelTable,
  ExcelUploadResponse,
  ExcelTableMetadata,
  ExcelTableData,
  ExcelFilePreview,
  parseFileName
} from '../api/excel';

export interface UseExcelUploadState {
  isUploading: boolean;
  uploadProgress: number;
  uploadStatus: 'idle' | 'preview' | 'uploading' | 'success' | 'error';
  error: string | null;
  result: ExcelUploadResponse | null;
  preview: ExcelFilePreview | null;
}

export const useExcelUpload = () => {
  const [state, setState] = useState<UseExcelUploadState>({
    isUploading: false,
    uploadProgress: 0,
    uploadStatus: 'idle',
    error: null,
    result: null,
    preview: null,
  });

  const resetState = useCallback(() => {
    setState({
      isUploading: false,
      uploadProgress: 0,
      uploadStatus: 'idle',
      error: null,
      result: null,
      preview: null,
    });
  }, []);

  const previewFile = useCallback((file: File) => {
    const preview = parseFileName(file.name);
    setState(prev => ({
      ...prev,
      uploadStatus: 'preview',
      preview,
      error: null,
    }));
  }, []);

  const uploadFile = useCallback(async (file: File) => {
    setState(prev => ({
      ...prev,
      isUploading: true,
      uploadStatus: 'uploading',
      uploadProgress: 0,
      error: null,
    }));

    try {
      // Симуляция прогресса загрузки
      const progressInterval = setInterval(() => {
        setState(prev => ({
          ...prev,
          uploadProgress: Math.min(prev.uploadProgress + Math.random() * 20, 90),
        }));
      }, 200);

      const result = await uploadExcelFile(file);

      clearInterval(progressInterval);

      setState(prev => ({
        ...prev,
        isUploading: false,
        uploadProgress: 100,
        uploadStatus: 'success',
        result,
      }));
    } catch (error) {
      setState(prev => ({
        ...prev,
        isUploading: false,
        uploadStatus: 'error',
        error: error instanceof Error ? error.message : 'Upload failed',
      }));
    }
  }, []);

  return {
    ...state,
    previewFile,
    uploadFile,
    resetState,
  };
};

export const useExcelTables = () => {
  const [tables, setTables] = useState<ExcelTableMetadata[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchTables = useCallback(async () => {
    setLoading(true);
    setError(null);
    
    try {
      const data = await getExcelTables();
      setTables(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch tables');
    } finally {
      setLoading(false);
    }
  }, []);

  const deleteTable = useCallback(async (tableName: string) => {
    try {
      await deleteExcelTable(tableName);
      await fetchTables(); // Обновляем список
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to delete table');
    }
  }, [fetchTables]);

  return {
    tables,
    loading,
    error,
    fetchTables,
    deleteTable,
  };
};

export const useExcelTableData = (tableName: string | null) => {
  const [data, setData] = useState<ExcelTableData | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const fetchData = useCallback(async (limit: number = 10, offset: number = 0) => {
    if (!tableName) return;
    
    setLoading(true);
    setError(null);
    
    try {
      const result = await getExcelTableData(tableName, limit, offset);
      setData(result);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch table data');
    } finally {
      setLoading(false);
    }
  }, [tableName]);

  return {
    data,
    loading,
    error,
    fetchData,
  };
};
