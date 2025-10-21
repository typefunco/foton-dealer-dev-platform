import React, { useEffect, useState } from 'react';
import { useExcelTables, useExcelTableData } from '../hooks/useExcel';

const ExcelTablesPage: React.FC = () => {
  const { tables, loading, error, fetchTables, deleteTable } = useExcelTables();
  const [selectedTable, setSelectedTable] = useState<string | null>(null);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<string | null>(null);
  const { data: tableData, loading: dataLoading, fetchData } = useExcelTableData(selectedTable);

  useEffect(() => {
    fetchTables();
  }, [fetchTables]);

  useEffect(() => {
    if (selectedTable) {
      fetchData(10, 0);
    }
  }, [selectedTable, fetchData]);

  const handleDeleteTable = async (tableName: string) => {
    try {
      await deleteTable(tableName);
      setShowDeleteConfirm(null);
      if (selectedTable === tableName) {
        setSelectedTable(null);
      }
    } catch (err) {
      console.error('Failed to delete table:', err);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600 mx-auto mb-4"></div>
          <p className="text-gray-600">Загрузка таблиц...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">
            Excel таблицы
          </h1>
          <p className="text-gray-600">
            Просмотр и управление созданными таблицами
          </p>
        </div>

        {error && (
          <div className="bg-red-50 border border-red-200 rounded-lg p-4 mb-6">
            <p className="text-red-700">{error}</p>
          </div>
        )}

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Список таблиц */}
          <div className="lg:col-span-1">
            <div className="bg-white rounded-lg shadow-lg p-6">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">
                Созданные таблицы ({tables.length})
              </h2>
              
              {tables.length === 0 ? (
                <p className="text-gray-500 text-center py-8">
                  Нет созданных таблиц
                </p>
              ) : (
                <div className="space-y-3">
                  {tables.map((table) => (
                    <div
                      key={table.table_name}
                      className={`p-4 border rounded-lg cursor-pointer transition-colors ${
                        selectedTable === table.table_name
                          ? 'border-blue-500 bg-blue-50'
                          : 'border-gray-200 hover:border-gray-300'
                      }`}
                      onClick={() => setSelectedTable(table.table_name)}
                    >
                      <div className="flex justify-between items-start">
                        <div className="flex-1">
                          <h3 className="font-medium text-gray-900 text-sm">
                            {table.table_name}
                          </h3>
                          <p className="text-xs text-gray-500 mt-1">
                            {table.rows_count} строк • {table.columns.length} колонок
                          </p>
                          <p className="text-xs text-gray-400 mt-1">
                            {formatDate(table.created_at)}
                          </p>
                        </div>
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setShowDeleteConfirm(table.table_name);
                          }}
                          className="text-red-500 hover:text-red-700 text-xs"
                        >
                          Удалить
                        </button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>

          {/* Данные таблицы */}
          <div className="lg:col-span-2">
            {selectedTable ? (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <div className="flex justify-between items-center mb-6">
                  <h2 className="text-lg font-semibold text-gray-900">
                    Данные таблицы: {selectedTable}
                  </h2>
                  <button
                    onClick={() => setSelectedTable(null)}
                    className="text-gray-500 hover:text-gray-700"
                  >
                    <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>

                {dataLoading ? (
                  <div className="text-center py-8">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto mb-4"></div>
                    <p className="text-gray-600">Загрузка данных...</p>
                  </div>
                ) : tableData ? (
                  <div className="overflow-x-auto">
                    <table className="min-w-full divide-y divide-gray-200">
                      <thead className="bg-gray-50">
                        <tr>
                          {Object.keys(tableData.data[0] || {}).map((column) => (
                            <th
                              key={column}
                              className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                            >
                              {column}
                            </th>
                          ))}
                        </tr>
                      </thead>
                      <tbody className="bg-white divide-y divide-gray-200">
                        {tableData.data.map((row, index) => (
                          <tr key={index}>
                            {Object.values(row).map((value, cellIndex) => (
                              <td
                                key={cellIndex}
                                className="px-6 py-4 whitespace-nowrap text-sm text-gray-900"
                              >
                                {value === null ? (
                                  <span className="text-gray-400 italic">null</span>
                                ) : (
                                  String(value)
                                )}
                              </td>
                            ))}
                          </tr>
                        ))}
                      </tbody>
                    </table>
                    
                    <div className="mt-4 text-sm text-gray-500">
                      Показано {tableData.data.length} из {tableData.pagination.total_count} записей
                    </div>
                  </div>
                ) : (
                  <p className="text-gray-500 text-center py-8">
                    Нет данных для отображения
                  </p>
                )}
              </div>
            ) : (
              <div className="bg-white rounded-lg shadow-lg p-6">
                <div className="text-center py-12">
                  <div className="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                    <svg className="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                  </div>
                  <h3 className="text-lg font-medium text-gray-900 mb-2">
                    Выберите таблицу
                  </h3>
                  <p className="text-gray-500">
                    Выберите таблицу из списка слева для просмотра данных
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Модальное окно подтверждения удаления */}
        {showDeleteConfirm && (
          <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
            <div className="bg-white rounded-lg p-6 max-w-md w-full mx-4">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                Подтверждение удаления
              </h3>
              <p className="text-gray-600 mb-6">
                Вы уверены, что хотите удалить таблицу "{showDeleteConfirm}"? 
                Это действие нельзя отменить.
              </p>
              <div className="flex space-x-4">
                <button
                  onClick={() => handleDeleteTable(showDeleteConfirm)}
                  className="flex-1 bg-red-600 text-white px-4 py-2 rounded-lg hover:bg-red-700 transition-colors"
                >
                  Удалить
                </button>
                <button
                  onClick={() => setShowDeleteConfirm(null)}
                  className="flex-1 border border-gray-300 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  Отмена
                </button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};

export default ExcelTablesPage;
