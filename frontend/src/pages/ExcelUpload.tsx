import React, { useRef, useState } from 'react';
import { useExcelUpload } from '../hooks/useExcel';

const ExcelUploadPage: React.FC = () => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const {
    uploadProgress,
    uploadStatus,
    error,
    result,
    preview,
    previewFile,
    uploadFile,
    resetState,
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
