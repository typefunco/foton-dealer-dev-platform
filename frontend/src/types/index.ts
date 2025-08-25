// Main types for the application

export interface User {
  id: string
  firstName: string
  lastName: string
  email: string
  phone: string
  company: string
  position: string
  city: string
  bio: string
  avatar?: string
}

export interface Stat {
  title: string
  value: string
  change: string
  changeType: 'positive' | 'negative'
  icon: React.ComponentType<{ className?: string }>
  color: string
}

export interface Activity {
  id: number
  action: string
  time: string
  type: 'success' | 'warning' | 'info' | 'error'
}

export interface Feature {
  icon: React.ComponentType<{ className?: string }>
  title: string
  description: string
  color: string
}

export interface NavigationItem {
  name: string
  href: string
  icon: React.ComponentType<{ className?: string }>
}
