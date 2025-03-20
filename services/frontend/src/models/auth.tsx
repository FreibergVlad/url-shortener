export interface IssueAuthenticationTokenRequest {
  email: string
  password: string
}

export interface IssueAuthenticationTokenResponse {
  token: string
  refreshToken: string
}

export interface RefreshAuthenticationTokenRequest {
  refreshToken: string
}

export interface RefreshAuthenticationTokenResponse {
  token: string
}

export interface Role {
  id: string
  name: string
  description: string
}
