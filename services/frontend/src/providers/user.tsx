import { useAPIContext } from '@/contexts/api'
import { UserContext } from '@/contexts/user'
import { OrganizationMembership } from '@/models/organization'
import { useLocalStorage } from '@uidotdev/usehooks'
import { useEffect } from 'react'

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { useGetUserInfo, useGetOrganizationMemberships } = useAPIContext()

  const { data: user, ...userQuery } = useGetUserInfo()
  const { data: organizationMemberships, ...organizationMembershipsQuery } = useGetOrganizationMemberships()

  const [currentOrganizationMembership, setCurrentOrganizationMembership] = useLocalStorage<OrganizationMembership | undefined>('current_organization_membership')

  useEffect(() => {
    if (organizationMemberships && organizationMemberships.length > 0 && !currentOrganizationMembership) {
      setCurrentOrganizationMembership(organizationMemberships[0])
    }
  }, [organizationMemberships, currentOrganizationMembership, setCurrentOrganizationMembership])

  return (
    <UserContext.Provider
      value={{
        user: user,
        organizationMemberships: organizationMemberships,
        currentOrganizationMembership: currentOrganizationMembership,
        setCurrentOrganizationMembership: setCurrentOrganizationMembership,
        isLoading: userQuery.isLoading || organizationMembershipsQuery.isLoading,
      }}
    >
      {children}
    </UserContext.Provider>
  )
}
