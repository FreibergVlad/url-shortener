import { GetMeResponse, User } from '@/models/user'
import { executeAPIRequest } from '@/services/api'
import { CreateOrganizationRequest, CreateOrganizationResponse, GetOrganizationMembershipsResponse, OrganizationMembership } from '@/models/organization'
import { AppError } from '@/services/errors'
import { useMutation, useQuery, useQueryClient, UseQueryResult } from '@tanstack/react-query'
import { ListShortURLResponse } from '@/models/shorturl'
import { APIContext } from '@/contexts/api'
import { useAuthContext } from '@/contexts/auth'

let ongoingRefreshTokenRequest: Promise<string> | null = null

export const APIProvider = ({ children }: { children: React.ReactNode }) => {
  const { accessToken, refreshUserAuthentication } = useAuthContext()
  const queryClient = useQueryClient()

  const baseRequestParams = {
    accessToken: accessToken,
    onTokenExpired: async () => {
      if (ongoingRefreshTokenRequest) {
        return ongoingRefreshTokenRequest
      }
      ongoingRefreshTokenRequest = (async () => {
        try {
          const response = await refreshUserAuthentication.mutateAsync(undefined)
          return response.token
        }
        finally {
          ongoingRefreshTokenRequest = null
        }
      })()
      return ongoingRefreshTokenRequest
    },
  }

  const useGetUserInfo = (): UseQueryResult<User, AppError> => {
    return useQuery({
      queryFn: async () => {
        const response = await executeAPIRequest<null, GetMeResponse>({
          endpoint: 'users/me',
          method: 'GET',
          ...baseRequestParams,
        })
        return response.user
      },
      queryKey: ['users', 'me'],
      enabled: !!accessToken,
      retry: false,
    })
  }

  const useGetOrganizationMemberships = (): UseQueryResult<OrganizationMembership[], AppError> => {
    return useQuery({
      queryFn: async () => {
        const response = await executeAPIRequest<null, GetOrganizationMembershipsResponse>({
          endpoint: 'users/me/organizations',
          method: 'GET',
          ...baseRequestParams,
        })
        return response.data
      },
      queryKey: ['users', 'me', 'organizations'],
      enabled: !!accessToken,
      retry: false,
    })
  }

  const useCreateOrganization = useMutation<CreateOrganizationResponse, AppError, CreateOrganizationRequest>({
    mutationFn: async (request) => {
      return executeAPIRequest<CreateOrganizationRequest, CreateOrganizationResponse>({
        endpoint: 'organizations',
        method: 'POST',
        body: request,
        ...baseRequestParams,
      })
    },
    onSuccess: async () => {
      await queryClient.invalidateQueries({ queryKey: ['users', 'me', 'organizations'] })
    },
  })

  const useListShortURL = (
    organizationId: string,
    pageNum: number,
    pageSize: number,
  ): UseQueryResult<ListShortURLResponse, AppError> => {
    return useQuery({
      queryFn: async () => {
        const response = await executeAPIRequest<null, GetOrganizationMembershipsResponse>({
          endpoint: `organizations/${organizationId}/short-urls`,
          method: 'GET',
          queryParams: { pageNum: String(pageNum), pageSize: String(pageSize) },
          ...baseRequestParams,
        })
        return response.data
      },
      queryKey: ['organizations', organizationId, 'short-urls', pageNum, pageSize],
      enabled: !!accessToken,
    })
  }

  return (
    <APIContext.Provider
      value={{
        useGetUserInfo: useGetUserInfo,
        useGetOrganizationMemberships: useGetOrganizationMemberships,
        useCreateOrganization: useCreateOrganization,
        useListShortURL: useListShortURL,
      }}
    >
      {children}
    </APIContext.Provider>
  )
}
