import { useAPIContext } from '@/contexts/api'
import { UserContext } from '@/contexts/user'

export const UserProvider = ({ children }: { children: React.ReactNode }) => {
  const { useGetUserInfo, useGetOrganizationMemberships } = useAPIContext()

  const { data: user, ...userQuery } = useGetUserInfo()
  const { data: organizationMemberships, ...organizationMembershipsQuery } = useGetOrganizationMemberships()

  return (
    <UserContext.Provider
      value={{
        user: user,
        organizationMemberships: organizationMemberships,
        isLoading: userQuery.isLoading || organizationMembershipsQuery.isLoading,
      }}
    >
      {children}
    </UserContext.Provider>
  )
}
