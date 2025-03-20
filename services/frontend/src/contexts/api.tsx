import { CreateOrganizationRequest, CreateOrganizationResponse, OrganizationMembership } from '@/models/organization'
import { ListShortURLResponse } from '@/models/shorturl'
import { User } from '@/models/user'
import { APIError } from '@/services/errors'
import { UseMutationResult, UseQueryResult } from '@tanstack/react-query'
import { createContext, useContext } from 'react'

interface APIContextType {
  useGetUserInfo(): UseQueryResult<User, APIError>
  useGetOrganizationMemberships(): UseQueryResult<OrganizationMembership[], APIError>
  useCreateOrganization: UseMutationResult<CreateOrganizationResponse, APIError, CreateOrganizationRequest>
  useListShortURL(organizationId: string, pageNum: number, pageSize: number): UseQueryResult<ListShortURLResponse, APIError>
}

export const APIContext = createContext<APIContextType | null>(null)

export const useAPIContext = (): APIContextType => {
  const context = useContext(APIContext)
  if (!context) {
    throw new Error('useAPIContext() must be used within the APIProvider')
  }
  return context
}
