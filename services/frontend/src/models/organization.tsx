import { Role } from './auth'

export interface ShortOrganization {
  id: string
  slug: string
}

export interface Organization {
  id: string
  name: string
  slug: string
}

export interface OrganizationMembership {
  organization: ShortOrganization
  role: Role
  createdAt: string
}

export interface GetOrganizationMembershipsResponse {
  data: OrganizationMembership[]
}

export interface CreateOrganizationRequest {
  name: string
  slug: string
}

export interface CreateOrganizationResponse {
  organization: Organization
}
