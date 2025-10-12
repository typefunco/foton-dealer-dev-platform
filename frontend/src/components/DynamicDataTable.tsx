import React from 'react'
import { useDynamicData, buildDynamicParams } from '../hooks/useDynamicData'
import { TableType } from '../api/dynamic'

interface DynamicDataTableProps {
  tableType: TableType
  filters: {
    region?: string
    quarter?: string
    year?: number
    dealers?: string[]
    limit?: number
    offset?: number
    sortBy?: string
    sortOrder?: 'asc' | 'desc'
  }
  columns: Array<{
    key: string
    title: string
    render?: (value: any, record: any) => React.ReactNode
  }>
}

export const DynamicDataTable: React.FC<DynamicDataTableProps> = ({
  tableType,
  filters,
  columns
}) => {
  const apiParams = buildDynamicParams(filters)
  
  const { data, response, loading, error, refetch } = useDynamicData({
    tableType,
    params: apiParams,
    enabled: true
  })

  if (loading) {
    return (
      <div className="flex items-center justify-center p-8">
        <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
        <span className="ml-2 text-gray-600">Loading...</span>
      </div>
    )
  }

  if (error) {
    return (
      <div className="bg-red-50 border border-red-200 rounded-lg p-4">
        <div className="flex items-center">
          <svg className="w-5 h-5 text-red-400 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span className="text-red-800 font-medium">Error loading data</span>
        </div>
        <p className="text-red-700 mt-1">{error}</p>
        <button
          onClick={refetch}
          className="mt-2 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
        >
          Retry
        </button>
      </div>
    )
  }

  if (!data || data.length === 0) {
    return (
      <div className="text-center py-8 text-gray-500">
        <svg className="w-16 h-16 mx-auto text-gray-300 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p className="text-lg font-medium">No data found</p>
        <p className="text-sm">Try adjusting your filters</p>
      </div>
    )
  }

  return (
    <div className="bg-white rounded-lg shadow-lg overflow-hidden">
      {/* Header with metadata */}
      {response && (
        <div className="bg-gray-50 px-6 py-4 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div>
              <h3 className="text-lg font-semibold text-gray-900 capitalize">
                {tableType.replace('_', ' ')} Data
              </h3>
              <p className="text-sm text-gray-600">
                {response.quarter} {response.year} • {response.regions.join(', ')} • {response.count} records
              </p>
            </div>
            <button
              onClick={refetch}
              className="px-4 py-2 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors"
            >
              Refresh
            </button>
          </div>
        </div>
      )}

      {/* Table */}
      <div className="overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              {columns.map((column) => (
                <th
                  key={column.key}
                  className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                >
                  {column.title}
                </th>
              ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {data.map((record: any, index: number) => (
              <tr key={record.id || index} className="hover:bg-gray-50">
                {columns.map((column) => (
                  <td key={column.key} className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                    {column.render 
                      ? column.render(record[column.key], record)
                      : record[column.key]
                    }
                  </td>
                ))}
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* Footer with pagination info */}
      {response && (
        <div className="bg-gray-50 px-6 py-3 border-t border-gray-200">
          <div className="flex items-center justify-between text-sm text-gray-600">
            <span>
              Showing {data.length} of {response.count} records
            </span>
            <span>
              Table: {response.tableType}
            </span>
          </div>
        </div>
      )}
    </div>
  )
}

// Пример использования компонента
export const ExampleUsage: React.FC = () => {
  const filters = {
    region: 'central',
    quarter: 'Q1',
    year: 2024,
    limit: 50
  }

  const columns = [
    { key: 'id', title: 'ID' },
    { key: 'name', title: 'Name' },
    { key: 'city', title: 'City' },
    { 
      key: 'class', 
      title: 'Class',
      render: (value: string) => (
        <span className={`px-2 py-1 text-xs font-medium rounded-full ${
          value === 'A' ? 'bg-green-100 text-green-800' :
          value === 'B' ? 'bg-blue-100 text-blue-800' :
          value === 'C' ? 'bg-yellow-100 text-yellow-800' :
          'bg-gray-100 text-gray-800'
        }`}>
          {value}
        </span>
      )
    },
    { key: 'checklist', title: 'Checklist Score' }
  ]

  return (
    <div className="p-6">
      <h2 className="text-2xl font-bold mb-4">Dealer Development Data</h2>
      <DynamicDataTable
        tableType="dealer_dev"
        filters={filters}
        columns={columns}
      />
    </div>
  )
}
