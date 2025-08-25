import React, { useState } from 'react'

interface User {
  id: string
  email: string
  firstName: string
  lastName: string
  region: string
  position: string
  createdAt: string
  status: 'active' | 'inactive'
}

interface RegionStats {
  region: string
  userCount: number
  users: User[]
}

const Admin: React.FC = () => {
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [newUser, setNewUser] = useState({
    email: '',
    firstName: '',
    lastName: '',
    region: '',
    position: ''
  })
  const [generatedCredentials, setGeneratedCredentials] = useState<{
    email: string
    password: string
  } | null>(null)
  
  // Search and filter states
  const [searchTerm, setSearchTerm] = useState('')
  const [selectedRegion, setSelectedRegion] = useState('')
  const [selectedPosition, setSelectedPosition] = useState('')
  
  // Pagination states
  const [currentPage, setCurrentPage] = useState(1)
  const usersPerPage = 10

  // Mock data for regions and users
  const regions = [
    'Central',
    'Caucasus', 
    'Volga',
    'Ural',
    'Siberia',
    'Far East',
    'North-West',
    'South'
  ]

  const mockUsers: User[] = [
    {
      id: '1',
      email: 'ivan.petrov@company.com',
      firstName: 'Ivan',
      lastName: 'Petrov',
      region: 'Central',
      position: 'Sales Manager',
      createdAt: '2024-01-15',
      status: 'active'
    },
    {
      id: '2',
      email: 'maria.sidorova@company.com',
      firstName: 'Maria',
      lastName: 'Sidorova',
      region: 'Caucasus',
      position: 'Regional Director',
      createdAt: '2024-02-01',
      status: 'active'
    },
    {
      id: '3',
      email: 'alex.kuznetsov@company.com',
      firstName: 'Alex',
      lastName: 'Kuznetsov',
      region: 'Volga',
      position: 'Account Manager',
      createdAt: '2024-01-20',
      status: 'active'
    },
    {
      id: '4',
      email: 'elena.ivanova@company.com',
      firstName: 'Elena',
      lastName: 'Ivanova',
      region: 'Ural',
      position: 'Sales Representative',
      createdAt: '2024-02-10',
      status: 'active'
    },
    {
      id: '5',
      email: 'dmitry.sokolov@company.com',
      firstName: 'Dmitry',
      lastName: 'Sokolov',
      region: 'Siberia',
      position: 'Regional Manager',
      createdAt: '2024-01-25',
      status: 'active'
    },
    {
      id: '6',
      email: 'anna.kozlov@company.com',
      firstName: 'Anna',
      lastName: 'Kozlov',
      region: 'Far East',
      position: 'Sales Director',
      createdAt: '2024-02-05',
      status: 'active'
    },
    {
      id: '7',
      email: 'sergey.morozov@company.com',
      firstName: 'Sergey',
      lastName: 'Morozov',
      region: 'North-West',
      position: 'Account Executive',
      createdAt: '2024-01-30',
      status: 'active'
    },
    {
      id: '8',
      email: 'olga.volkova@company.com',
      firstName: 'Olga',
      lastName: 'Volkova',
      region: 'South',
      position: 'Sales Manager',
      createdAt: '2024-02-15',
      status: 'active'
    }
  ]

  // Calculate statistics by region
  const regionStats: RegionStats[] = regions.map(region => {
    const usersInRegion = mockUsers.filter(user => user.region === region)
    return {
      region,
      userCount: usersInRegion.length,
      users: usersInRegion
    }
  })

  const totalUsers = mockUsers.length

  // Filter users based on search and filters
  const filteredUsers = mockUsers.filter(user => {
    const matchesSearch = searchTerm === '' || 
      user.firstName.toLowerCase().includes(searchTerm.toLowerCase()) ||
      user.lastName.toLowerCase().includes(searchTerm.toLowerCase())
    
    const matchesRegion = selectedRegion === '' || user.region === selectedRegion
    const matchesPosition = selectedPosition === '' || user.position === selectedPosition
    
    return matchesSearch && matchesRegion && matchesPosition
  })

  // Predefined positions for new accounts
  const predefinedPositions = [
    'Sales Manager',
    'Regional Director', 
    'Regional Manager',
    'Sales Director',
    'Account Manager',
    'Sales Representative',
    'Account Executive',
    'Senior Sales Manager',
    'Head of Sales',
    'Business Development Manager'
  ]

  // Get unique positions for filter (from existing users)
  const positions = [...new Set(mockUsers.map(user => user.position))]
  
  // Pagination logic
  const totalPages = Math.ceil(filteredUsers.length / usersPerPage)
  const startIndex = (currentPage - 1) * usersPerPage
  const endIndex = startIndex + usersPerPage
  const currentUsers = filteredUsers.slice(startIndex, endIndex)
  
  // Reset to first page when filters change
  React.useEffect(() => {
    setCurrentPage(1)
  }, [searchTerm, selectedRegion, selectedPosition])

  const generatePassword = () => {
    const chars = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*'
    let password = ''
    for (let i = 0; i < 12; i++) {
      password += chars.charAt(Math.floor(Math.random() * chars.length))
    }
    return password
  }

  const handleCreateAccount = (e: React.FormEvent) => {
    e.preventDefault()
    
    // Generate credentials
    const password = generatePassword()
    const credentials = {
      email: newUser.email,
      password: password
    }
    
    setGeneratedCredentials(credentials)
    
    // Reset form
    setNewUser({
      email: '',
      firstName: '',
      lastName: '',
      region: '',
      position: ''
    })
    
    setShowCreateForm(false)
  }

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    const { name, value } = e.target
    setNewUser(prev => ({
      ...prev,
      [name]: value
    }))
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
    // You could add a toast notification here
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-900 via-blue-800 to-blue-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Header */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold text-white mb-4">
            ADMIN PANEL
          </h1>
          <p className="text-xl text-gray-300">
            FOTON DEALER DEVELOPMENT PLATFORM
          </p>
        </div>

        {/* Total Users Overview */}
        <div className="bg-white rounded-2xl shadow-xl p-6 mb-8">
          <div className="text-center">
            <h2 className="text-3xl font-bold text-gray-900 mb-2">
              Total Users: {totalUsers}
            </h2>
            <p className="text-gray-600">
              Active accounts across all regions
            </p>
          </div>
        </div>

        {/* Create Account Section */}
        <div className="bg-white rounded-2xl shadow-xl p-6 mb-8">
          <div className="flex justify-between items-center mb-6">
            <h3 className="text-2xl font-bold text-gray-900">
              Create New Account
            </h3>
            <button
              onClick={() => setShowCreateForm(!showCreateForm)}
              className="bg-blue-600 hover:bg-blue-700 text-white px-6 py-2 rounded-xl transition-colors duration-200"
            >
              {showCreateForm ? 'Cancel' : 'Create Account'}
            </button>
          </div>

          {showCreateForm && (
            <form onSubmit={handleCreateAccount} className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Email Address *
                  </label>
                  <input
                    type="email"
                    name="email"
                    required
                    value={newUser.email}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="user@company.com"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Region *
                  </label>
                  <select
                    name="region"
                    required
                    value={newUser.region}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">Select Region</option>
                    {regions.map(region => (
                      <option key={region} value={region}>{region}</option>
                    ))}
                  </select>
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    First Name *
                  </label>
                  <input
                    type="text"
                    name="firstName"
                    required
                    value={newUser.firstName}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="First Name"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Last Name *
                  </label>
                  <input
                    type="text"
                    name="lastName"
                    required
                    value={newUser.lastName}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                    placeholder="Last Name"
                  />
                </div>
                <div className="md:col-span-2">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Position *
                  </label>
                  <select
                    name="position"
                    required
                    value={newUser.position}
                    onChange={handleInputChange}
                    className="w-full px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500"
                  >
                    <option value="">Select Position</option>
                    {predefinedPositions.map(position => (
                      <option key={position} value={position}>{position}</option>
                    ))}
                  </select>
                </div>
              </div>
              <div className="flex justify-end">
                <button
                  type="submit"
                  className="bg-green-600 hover:bg-green-700 text-white px-8 py-3 rounded-xl transition-colors duration-200 font-medium"
                >
                  Create Account
                </button>
              </div>
            </form>
          )}

          {/* Generated Credentials Display */}
          {generatedCredentials && (
            <div className="mt-6 p-4 bg-green-50 border border-green-200 rounded-xl">
              <h4 className="text-lg font-semibold text-green-800 mb-3">
                Account Created Successfully!
              </h4>
              <div className="space-y-2">
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-green-700">Email:</span>
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-green-800">{generatedCredentials.email}</span>
                    <button
                      onClick={() => copyToClipboard(generatedCredentials.email)}
                      className="text-green-600 hover:text-green-800"
                      title="Copy to clipboard"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                    </button>
                  </div>
                </div>
                <div className="flex items-center justify-between">
                  <span className="text-sm font-medium text-green-700">Password:</span>
                  <div className="flex items-center space-x-2">
                    <span className="text-sm text-green-800 font-mono">{generatedCredentials.password}</span>
                    <button
                      onClick={() => copyToClipboard(generatedCredentials.password)}
                      className="text-green-600 hover:text-green-800"
                      title="Copy to clipboard"
                    >
                      <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
                      </svg>
                    </button>
                  </div>
                </div>
              </div>
              <div className="mt-4 text-xs text-green-600">
                Please save these credentials securely. The password cannot be recovered later.
              </div>
            </div>
          )}
        </div>

        {/* Regional Statistics */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {regionStats.map((stat) => (
            <div key={stat.region} className="bg-white rounded-2xl shadow-xl p-6">
              <div className="text-center">
                <h3 className="text-xl font-bold text-gray-900 mb-2">
                  {stat.region}
                </h3>
                <div className="text-3xl font-bold text-blue-600 mb-2">
                  {stat.userCount}
                </div>
                <p className="text-sm text-gray-600">
                  {stat.userCount === 1 ? 'User' : 'Users'}
                </p>
              </div>
            </div>
          ))}
        </div>



        {/* Detailed User List */}
        <div className="bg-white rounded-2xl shadow-xl p-6">
          <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between mb-6">
            <h3 className="text-2xl font-bold text-gray-900 mb-4 lg:mb-0">
              All Users ({filteredUsers.length})
            </h3>
            
            {/* Filters */}
            <div className="flex flex-col sm:flex-row gap-3">
              {/* Search by Name */}
              <div className="min-w-[200px]">
                <input
                  type="text"
                  placeholder="Search by name..."
                  value={searchTerm}
                  onChange={(e) => setSearchTerm(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                />
              </div>
              
              {/* Region Filter */}
              <div className="min-w-[150px]">
                <select
                  value={selectedRegion}
                  onChange={(e) => setSelectedRegion(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                >
                  <option value="">All Regions</option>
                  {regions.map(region => (
                    <option key={region} value={region}>{region}</option>
                  ))}
                </select>
              </div>
              
              {/* Position Filter */}
              <div className="min-w-[150px]">
                <select
                  value={selectedPosition}
                  onChange={(e) => setSelectedPosition(e.target.value)}
                  className="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
                >
                  <option value="">All Positions</option>
                  {positions.map(position => (
                    <option key={position} value={position}>{position}</option>
                  ))}
                </select>
              </div>
              
              {/* Clear Filters */}
              <button
                onClick={() => {
                  setSearchTerm('')
                  setSelectedRegion('')
                  setSelectedPosition('')
                }}
                className="px-4 py-2 bg-gray-500 hover:bg-gray-600 text-white rounded-lg transition-colors duration-200 text-sm"
              >
                Clear
              </button>
            </div>
          </div>
          
          <div className="overflow-x-auto">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                    User
                  </th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Region
                  </th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Position
                  </th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Created
                  </th>
                  <th className="px-6 py-3 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Status
                  </th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                {currentUsers.map((user) => (
                  <tr key={user.id} className="hover:bg-gray-50">
                    <td className="px-6 py-4 whitespace-nowrap text-center">
                      <div>
                        <div className="text-sm font-medium text-gray-900">
                          {user.firstName} {user.lastName}
                        </div>
                        <div className="text-sm text-gray-500">
                          {user.email}
                        </div>
                      </div>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center">
                      <span className="inline-flex px-2 py-1 text-xs font-semibold rounded-full bg-blue-100 text-blue-800">
                        {user.region}
                      </span>
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900 text-center">
                      {user.position}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500 text-center">
                      {new Date(user.createdAt).toLocaleDateString()}
                    </td>
                    <td className="px-6 py-4 whitespace-nowrap text-center">
                      <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${
                        user.status === 'active' 
                          ? 'bg-green-100 text-green-800' 
                          : 'bg-red-100 text-red-800'
                      }`}>
                        {user.status}
                      </span>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          
          {/* Pagination */}
          {totalPages > 1 && (
            <div className="mt-6 flex items-center justify-between">
              <div className="text-sm text-gray-700">
                Showing {startIndex + 1} to {Math.min(endIndex, filteredUsers.length)} of {filteredUsers.length} results
              </div>
              
              <div className="flex items-center space-x-2">
                <button
                  onClick={() => setCurrentPage(prev => Math.max(prev - 1, 1))}
                  disabled={currentPage === 1}
                  className="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  Previous
                </button>
                
                <div className="flex items-center space-x-1">
                  {Array.from({ length: totalPages }, (_, i) => i + 1).map(page => (
                    <button
                      key={page}
                      onClick={() => setCurrentPage(page)}
                      className={`px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-200 ${
                        currentPage === page
                          ? 'bg-blue-600 text-white'
                          : 'text-gray-700 bg-white border border-gray-300 hover:bg-gray-50'
                      }`}
                    >
                      {page}
                    </button>
                  ))}
                </div>
                
                <button
                  onClick={() => setCurrentPage(prev => Math.min(prev + 1, totalPages))}
                  disabled={currentPage === totalPages}
                  className="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  Next
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

export default Admin
