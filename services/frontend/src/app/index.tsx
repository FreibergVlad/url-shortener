import Loader from '@/components/loader'
import { useUserContext } from '@/contexts/user'
import { AuthenticatedWithOrganizationRouter, AuthenticatedWithoutOrganizationRouter, UnauthenticatedRouter } from '@/routers'

export default function App() {
  const { user, organizationMemberships, isLoading } = useUserContext()

  if (isLoading) {
    return <Loader />
  }

  if (!user) {
    return <UnauthenticatedRouter />
  }

  if (!organizationMemberships?.length) {
    return <AuthenticatedWithoutOrganizationRouter />
  }

  return <AuthenticatedWithOrganizationRouter />
}
