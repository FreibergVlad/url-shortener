import { IssueAuthenticationTokenRequest, IssueAuthenticationTokenResponse, RefreshAuthenticationTokenResponse } from '@/models/auth'
import { CreateUserRequest, CreateUserResponse } from '@/models/user'
import { AppError } from '@/services/errors'
import { UseMutationResult } from '@tanstack/react-query'
import { createContext, useContext } from 'react'

interface AuthContextType {
  accessToken: string | undefined
  setAccessToken(token: string): void

  refreshToken: string | undefined
  setRefreshToken(token: string): void

  createUser: UseMutationResult<CreateUserResponse, AppError, CreateUserRequest>
  authenticateUser: UseMutationResult<IssueAuthenticationTokenResponse, AppError, IssueAuthenticationTokenRequest>
  refreshUserAuthentication: UseMutationResult<RefreshAuthenticationTokenResponse, AppError, undefined>
  logout: () => void
}

export const AuthContext = createContext<AuthContextType | null>(null)

export const useAuthContext = (): AuthContextType => {
  const context = useContext(AuthContext)
  if (!context) {
    throw new Error('useAuthContext() must be used within the AuthProvider')
  }
  return context
}
