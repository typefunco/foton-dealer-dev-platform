import React, { useRef, useState } from 'react';
import { useExcelUpload } from '../hooks/useExcel';

const ExcelUploadPage: React.FC = () => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const {
    uploadProgress,
    uploadStatus,
    uploadType,
    error,
    result,
    preview,
    previewFile,
    uploadFile,
    resetState,
    setUploadType,
  } = useExcelUpload();

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      previewFile(file);
    }
  };

  const handleUpload = () => {
    if (selectedFile) {
      uploadFile(selectedFile);
    }
  };

  const handleReset = () => {
    resetState();
    setSelectedFile(null);
    if (fileInputRef.current) {
      fileInputRef.current.value = '';
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="bg-white rounded-lg shadow-lg p-8">
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">
              Загрузка Excel файлов
            </h1>
            <p className="text-gray-600">
              Загрузите Excel файл для создания таблицы в базе данных
            </p>
          </div>

          {/* Общая информация */}
          {uploadStatus === 'idle' && (
            <div className="mb-6 bg-purple-50 border border-purple-200 rounded-lg p-5">
              <h3 className="font-semibold text-purple-900 mb-3 flex items-center">
                <svg className="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                Общая информация
              </h3>
              <div className="text-sm text-purple-800 space-y-2">
                <div className="bg-red-50 border border-red-300 rounded p-3 mb-3">
                  <p className="font-semibold text-red-900 mb-2">⚠️ Критически важно:</p>
                  <ul className="list-disc list-inside space-y-1 text-xs">
                    <li><strong>Названия городов должны точно совпадать</strong> в обоих файлах (данные дилеров и бренды)</li>
                    <li><strong>Примеры ошибок:</strong> Sankt-Petersburg ≠ Saint-Petersburg, Вологда ≠ Vologda</li>
                    <li><strong>Названия дилеров должны быть идентичны</strong> в обоих таблицах</li>
                    <li>Проверьте точное совпадение заглавных и строчных букв</li>
                  </ul>
                </div>
                <p className="font-medium">📌 Порядок загрузки данных:</p>
                <ol className="list-decimal list-inside space-y-1 ml-2">
                  <li><strong>Сначала</strong> загрузите "Данные дилеров" — это создаст таблицу для указанного квартала</li>
                  <li><strong>Затем</strong> загрузите "Бренды и бизнесы" — заполнит информацию о брендах и побочных бизнесах</li>
                </ol>
                <p className="font-medium mt-3">⚠️ Важно:</p>
                <ul className="list-disc list-inside space-y-1 ml-2">
                  <li>Название файла должно содержать регион, год и квартал, например: <code className="bg-white px-1 rounded">NW_2025_Q3.xlsx</code> или <code className="bg-white px-1 rounded">Central-Byside-Businesses_2025_Q3.xlsx</code></li>
                  <li>Для брендов и бизнесов обязательно наличие колонок "Dealer Name" и "Dealer City"</li>
                  <li>Первая строка файла всегда должна содержать заголовки колонок</li>
                </ul>
              </div>
            </div>
          )}

          {/* Выбор типа загрузки */}
          {uploadStatus === 'idle' && (
            <div className="mb-6">
              <label className="block text-sm font-medium text-gray-700 mb-3">
                Тип загрузки
              </label>
              <div className="grid grid-cols-2 gap-4">
                <button
                  onClick={() => setUploadType('dealer_data')}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    uploadType === 'dealer_data'
                      ? 'border-blue-500 bg-blue-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="text-center">
                    <svg className="w-8 h-8 mx-auto mb-2 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                    </svg>
                    <p className="font-semibold text-gray-900">Данные дилеров</p>
                    <p className="text-xs text-gray-500">Sales, After Sales, Performance</p>
                  </div>
                </button>
                <button
                  onClick={() => setUploadType('brands')}
                  className={`p-4 rounded-lg border-2 transition-all ${
                    uploadType === 'brands'
                      ? 'border-green-500 bg-green-50'
                      : 'border-gray-200 hover:border-gray-300'
                  }`}
                >
                  <div className="text-center">
                    <svg className="w-8 h-8 mx-auto mb-2 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                    </svg>
                    <p className="font-semibold text-gray-900">Бренды и бизнесы</p>
                    <p className="text-xs text-gray-500">Brands & By-side Businesses</p>
                  </div>
                </button>
              </div>
              
              {/* Документация для "Данные дилеров" */}
              {uploadType === 'dealer_data' && (
                <div className="mt-4 space-y-3">
                  <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
                    <h4 className="font-semibold text-blue-900 mb-2">📋 Требования к файлу с данными дилеров</h4>
                    <div className="text-xs text-blue-800 space-y-2">
                      <div>
                        <strong>Формат названия файла:</strong> 
                        <br />
                        <code className="bg-white px-2 py-1 rounded block mt-1">
                          NW_2025_Q3.xlsx<br/>
                          Central_2025_Q3.xlsx<br/>
                          FarEast_2025_Q3.xlsx
                        </code>
                      </div>
                      <div>
                        <strong>Структура файла:</strong>
                        <ul className="list-disc list-inside mt-1 space-y-1">
                          <li><strong>Строка 1:</strong> Пустая строка</li>
                          <li><strong>Строка 2:</strong> Заголовки колонок</li>
                          <li><strong>Строка 3+:</strong> Данные дилеров</li>
                        </ul>
                      </div>
                      <div>
                        <strong>Регионы файлов:</strong>
                        <ul className="list-disc list-inside mt-1 space-y-1">
                          <li>NW (North-West)</li>
                          <li>Central</li>
                          <li>FarEast</li>
                          <li>Volga</li>
                          <li>South</li>
                          <li>Ural</li>
                          <li>Siberia</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {/* Документация для "Бренды и бизнесы" */}
              {uploadType === 'brands' && (
                <div className="mt-4 space-y-3">
                  <div className="bg-green-50 border border-green-200 rounded-lg p-4">
                    <h4 className="font-semibold text-green-900 mb-2">📋 Требования к файлу с брендами и бизнесами</h4>
                    <div className="text-xs text-green-800 space-y-2">
                      <div>
                        <strong>Формат названия файла:</strong> 
                        <br />
                        <code className="bg-white px-2 py-1 rounded block mt-1">
                          NW-Byside-Businesses_2025_Q3.xlsx<br/>
                          Central-Byside-Businesses_2025_Q3.xlsx<br/>
                          FarEast-Byside-Businesses_2025_Q3.xlsx
                        </code>
                      </div>
                      <div>
                        <strong>Обязательные колонки (7 колонок):</strong>
                        <ol className="list-decimal list-inside mt-1 space-y-1">
                          <li><strong>Dealer Name</strong> - Название дилера</li>
                          <li><strong>Manager</strong> - Менеджер</li>
                          <li><strong>Dealer City</strong> - Город дилера</li>
                          <li><strong>Brands</strong> - Бренды (через запятую)</li>
                          <li><strong>Foton Sales Personnel</strong> - Количество продавцов Foton</li>
                          <li><strong>Sales Target</strong> - Целевые продажи</li>
                          <li><strong>By-side Businesses</strong> - Побочные бизнесы (через запятую)</li>
                        </ol>
                        <div className="bg-yellow-50 border border-yellow-300 rounded p-2 mt-2 text-xs">
                          <strong>💡 Примечание:</strong> Все колонки обязательны, но значения могут быть пустыми
                        </div>
                      </div>
                      <div>
                        <strong>Структура:</strong> Первая строка - заголовки колонок, со второй - данные
                      </div>
                      <div className="bg-red-50 border border-red-300 rounded p-2 mt-2">
                        <strong>⚠️ Важно:</strong> Таблица с данными дилеров должна быть создана через "Данные дилеров" перед загрузкой брендов. Соответствие дилеров идет по <strong>Dealer Name + Dealer City</strong>.
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Загрузка файла */}
          {uploadStatus === 'idle' && (
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center hover:border-gray-400 transition-colors">
              <div className="space-y-4">
                <div className="mx-auto w-12 h-12 bg-blue-100 rounded-full flex items-center justify-center">
                  <svg className="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                  </svg>
                </div>
                <div>
                  <p className="text-lg font-medium text-gray-900">Выберите Excel файл</p>
                  <p className="text-sm text-gray-500">Поддерживаются файлы .xlsx и .xls</p>
                </div>
                <input
                  ref={fileInputRef}
                  type="file"
                  accept=".xlsx,.xls"
                  onChange={handleFileSelect}
                  className="hidden"
                />
                <button
                  onClick={() => fileInputRef.current?.click()}
                  className="bg-blue-600 text-white px-6 py-2 rounded-lg hover:bg-blue-700 transition-colors"
                >
                  Выбрать файл
                </button>
              </div>
            </div>
          )}

          {/* Предварительный просмотр */}
          {uploadStatus === 'preview' && preview && (
            <div className="space-y-6">
              <div className="bg-blue-50 border border-blue-200 rounded-lg p-6">
                <h3 className="text-lg font-semibold text-blue-900 mb-4">
                  Предварительный просмотр
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Название файла</label>
                    <p className="mt-1 text-sm text-gray-900">{preview.fileName}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Размер файла</label>
                    <p className="mt-1 text-sm text-gray-900">
                      {selectedFile && formatFileSize(selectedFile.size)}
                    </p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Квартал</label>
                    <p className="mt-1 text-sm text-gray-900">{preview.quarter}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Год</label>
                    <p className="mt-1 text-sm text-gray-900">{preview.year}</p>
                  </div>
                  <div className="md:col-span-2">
                    <label className="block text-sm font-medium text-gray-700">Название таблицы</label>
                    <p className="mt-1 text-sm text-gray-900 font-mono bg-gray-100 px-3 py-2 rounded">
                      {preview.tableName}
                    </p>
                  </div>
                </div>
              </div>

              <div className="flex space-x-4">
                <button
                  onClick={handleUpload}
                  className="flex-1 bg-green-600 text-white px-6 py-3 rounded-lg hover:bg-green-700 transition-colors font-medium"
                >
                  Создать таблицу
                </button>
                <button
                  onClick={handleReset}
                  className="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
                >
                  Отмена
                </button>
              </div>
            </div>
          )}

          {/* Процесс загрузки */}
          {uploadStatus === 'uploading' && (
            <div className="space-y-6">
              <div className="text-center">
                <div className="mx-auto w-16 h-16 bg-blue-100 rounded-full flex items-center justify-center mb-4">
                  <svg className="w-8 h-8 text-blue-600 animate-spin" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                </div>
                <h3 className="text-lg font-semibold text-gray-900 mb-2">
                  Обработка файла...
                </h3>
                <p className="text-gray-600">
                  Пожалуйста, подождите, пока файл обрабатывается
                </p>
              </div>

              <div className="bg-gray-200 rounded-full h-3">
                <div
                  className="bg-blue-600 h-3 rounded-full transition-all duration-300 ease-out"
                  style={{ width: `${uploadProgress}%` }}
                />
              </div>
              <p className="text-center text-sm text-gray-600">
                {Math.round(uploadProgress)}% завершено
              </p>
            </div>
          )}

          {/* Результат успешной загрузки */}
          {uploadStatus === 'success' && result && (
            <div className="space-y-6">
              <div className="bg-green-50 border border-green-200 rounded-lg p-6">
                <div className="flex items-center mb-4">
                  <div className="mx-auto w-12 h-12 bg-green-100 rounded-full flex items-center justify-center">
                    <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                  </div>
                </div>
                <h3 className="text-lg font-semibold text-green-900 text-center mb-4">
                  Файл успешно обработан!
                </h3>
                
                {/* Результаты для dealer_data */}
                {uploadType === 'dealer_data' && 'tables_created' in result && (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Созданные таблицы</label>
                      <p className="mt-1 text-sm text-gray-900">
                        {result.tables_created.join(', ')}
                      </p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Строк вставлено</label>
                      <p className="mt-1 text-sm text-gray-900">{result.rows_inserted}</p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Время обработки</label>
                      <p className="mt-1 text-sm text-gray-900">
                        {(result.processing_time / 1000000000).toFixed(2)} сек
                      </p>
                    </div>
                    <div>
                      <label className="block text-sm font-medium text-gray-700">Статус</label>
                      <p className="mt-1 text-sm text-green-600 font-medium">Успешно</p>
                    </div>
                  </div>
                )}

                {/* Результаты для brands */}
                {uploadType === 'brands' && 'updated_count' in result && (
                  <div className="space-y-4">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700">Обновлено дилеров</label>
                        <p className="mt-1 text-lg font-bold text-green-600">{result.updated_count}</p>
                      </div>
                      <div>
                        <label className="block text-sm font-medium text-gray-700">Время обработки</label>
                        <p className="mt-1 text-sm text-gray-900">{result.processing_time}</p>
                      </div>
                    </div>
                    {result.not_found_dealers.length > 0 && (
                      <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
                        <label className="block text-sm font-medium text-yellow-800 mb-2">
                          Дилеры не найдены ({result.not_found_dealers.length})
                        </label>
                        <ul className="text-xs text-yellow-700 space-y-1">
                          {result.not_found_dealers.map((dealer, idx) => (
                            <li key={idx}>• {dealer}</li>
                          ))}
                        </ul>
                      </div>
                    )}
                  </div>
                )}
              </div>

              <div className="flex justify-center">
                <button
                  onClick={handleReset}
                  className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors font-medium"
                >
                  Загрузить другой файл
                </button>
              </div>
            </div>
          )}

          {/* Ошибка */}
          {uploadStatus === 'error' && error && (
            <div className="space-y-6">
              <div className="bg-red-50 border border-red-200 rounded-lg p-6">
                <div className="flex items-center mb-4">
                  <div className="mx-auto w-12 h-12 bg-red-100 rounded-full flex items-center justify-center">
                    <svg className="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </div>
                </div>
                <h3 className="text-lg font-semibold text-red-900 text-center mb-4">
                  Ошибка при обработке файла
                </h3>
                <p className="text-red-700 text-center">{error}</p>
              </div>

              <div className="flex justify-center">
                <button
                  onClick={handleReset}
                  className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition-colors font-medium"
                >
                  Попробовать снова
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default ExcelUploadPage;
