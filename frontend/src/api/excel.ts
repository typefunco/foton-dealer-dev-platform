import { API_BASE_URL } from './index';

export interface ExcelUploadResponse {
  status: string;
  message: string;
  tables_created: string[];
  rows_inserted: number;
  processing_time: number;
}

export interface ExcelTableMetadata {
  table_name: string;
  rows_count: number;
  created_at: string;
  columns: string[];
}

export interface ExcelTableData {
  data: Record<string, any>[];
  pagination: {
    limit: number;
    page: number;
    total_count: number;
    total_pages: number;
  };
}

export interface ExcelFilePreview {
  fileName: string;
  tableName: string;
  quarter: string;
  year: number;
  estimatedRows: number;
  sheets: string[];
}

// Загрузка Excel файла
export const uploadExcelFile = async (file: File): Promise<ExcelUploadResponse> => {
  const formData = new FormData();
  formData.append('file', file);

  // Получаем токен из localStorage
    
    // Получаем токен из localStorage
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(`${API_BASE_URL}/api/excel/upload`, {
    method: 'POST',
    headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      },
    body: formData,
  });

  if (!response.ok) {
    const error = await response.json();
    throw new Error(error.error || 'Failed to upload Excel file');
  }

  return response.json();
};

// Получение списка таблиц
export const getExcelTables = async (): Promise<ExcelTableMetadata[]> => {
  // Получаем токен из localStorage
    
    // Получаем токен из localStorage
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(`${API_BASE_URL}/api/excel/tables`, {
    headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch Excel tables');
  }

  return response.json();
};

// Получение метаданных таблицы
export const getExcelTableMetadata = async (tableName: string): Promise<ExcelTableMetadata> => {
  // Получаем токен из localStorage
    
    // Получаем токен из localStorage
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(`${API_BASE_URL}/api/excel/tables/${tableName}`, {
    headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      }
  });
  
  if (!response.ok) {
    throw new Error('Failed to fetch table metadata');
  }

  return response.json();
};

// Получение данных таблицы
export const getExcelTableData = async (
  tableName: string, 
  limit: number = 10, 
  offset: number = 0
): Promise<ExcelTableData> => {
  // Получаем токен из localStorage
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(
    `${API_BASE_URL}/excel/tables/${tableName}/data?limit=${limit}&offset=${offset}`,
    {
      headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      }
    }
  );
  
  if (!response.ok) {
    throw new Error('Failed to fetch table data');
  }

  return response.json();
};

// Удаление таблицы
export const deleteExcelTable = async (tableName: string): Promise<void> => {
  // Получаем токен из localStorage
    
    // Получаем токен из localStorage
  
  // Получаем токен из localStorage
  const token = localStorage.getItem('auth_token');
  
  const response = await fetch(`${API_BASE_URL}/api/excel/tables/${tableName}`, {
    method: 'DELETE',
    headers: {
        'Content-Type': 'application/json',
        ...(token && { 'Authorization': `Bearer ${token}` }),
      }
  });
  
  if (!response.ok) {
    throw new Error('Failed to delete table');
  }
};

// Предварительный просмотр файла (парсинг названия файла)
export const parseFileName = (fileName: string): ExcelFilePreview => {
  // Убираем расширение
  const name = fileName.replace(/\.(xlsx?|XLSX?)$/i, '');
  
  // Извлекаем квартал и год
  const quarterYearMatch = name.match(/q([1-4])[_\-\s]*(\d{4})/i);
  const quarterMatch = name.match(/q([1-4])/i);
  
  let quarter = 'Q1';
  let year = new Date().getFullYear();
  
  if (quarterYearMatch) {
    quarter = `Q${quarterYearMatch[1]}`;
    year = parseInt(quarterYearMatch[2]);
  } else if (quarterMatch) {
    quarter = `Q${quarterMatch[1]}`;
  }
  
  // Генерируем название таблицы
  const tableName = `dealer_net_${year}_${quarter.toLowerCase()}`;
  
  return {
    fileName,
    tableName,
    quarter,
    year,
    estimatedRows: 0, // Будет заполнено после загрузки
    sheets: [], // Будет заполнено после загрузки
  };
};
