import { OrganizationMembership } from '@/models/organization'
import { User } from '@/models/user'
import { createContext, useContext } from 'react'

interface UserContextType {
  user: User | undefined
  organizationMemberships: OrganizationMembership[] | undefined
  isLoading: boolean | undefined
}

export const UserContext = createContext<UserContextType | null>(null)

export const useUserContext = (): UserContextType => {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useUserContext() must be used within the UserProvider')
  }
  return context
}
