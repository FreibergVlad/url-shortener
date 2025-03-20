import { AuthContext } from '@/contexts/auth'
import { IssueAuthenticationTokenRequest, IssueAuthenticationTokenResponse, RefreshAuthenticationTokenRequest, RefreshAuthenticationTokenResponse } from '@/models/auth'
import { CreateUserRequest, CreateUserResponse } from '@/models/user'
import { executeAPIRequest } from '@/services/api'
import { AppError } from '@/services/errors'
import { useMutation } from '@tanstack/react-query'
import { useLocalStorage } from '@uidotdev/usehooks'

const ACCESS_TOKEN_STORAGE_KEY = 'access_token'
const REFRESH_TOKEN_STORAGE_KEY = 'refresh_token'

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [accessToken, setAccessToken] = useLocalStorage<string | undefined>(ACCESS_TOKEN_STORAGE_KEY)
  const [refreshToken, setRefreshToken] = useLocalStorage<string | undefined>(REFRESH_TOKEN_STORAGE_KEY)

  const createUser = useMutation<CreateUserResponse, AppError, CreateUserRequest>({
    mutationFn: async (request) => {
      return executeAPIRequest<CreateUserRequest, CreateUserResponse>({
        endpoint: 'users',
        method: 'POST',
        body: request,
      })
    },
  })

  const authenticateUser = useMutation<IssueAuthenticationTokenResponse, AppError, IssueAuthenticationTokenRequest>({
    mutationFn: async (request) => {
      return executeAPIRequest<IssueAuthenticationTokenRequest, IssueAuthenticationTokenResponse>({
        endpoint: 'tokens/issue',
        method: 'POST',
        body: request,
      })
    },
    onSuccess: (response) => {
      setAccessToken(response.token)
      setRefreshToken(response.refreshToken)
    },
  })

  const refreshUserAuthentication = useMutation<RefreshAuthenticationTokenResponse, AppError, undefined>({
    mutationFn: async () => {
      if (!refreshToken) {
        throw new Error('refreshToken doesn\'t exist in refreshUserAuthentication')
      }
      return executeAPIRequest<RefreshAuthenticationTokenRequest, RefreshAuthenticationTokenResponse>({
        endpoint: 'tokens/refresh',
        method: 'POST',
        body: { refreshToken },
      })
    },
    onSuccess: (response) => {
      setAccessToken(response.token)
    },
    onError: () => { logout() },
  })

  const logout = () => {
    setAccessToken(undefined)
    setRefreshToken(undefined)
    window.location.reload()
  }

  return (
    <AuthContext.Provider
      value={{
        accessToken: accessToken,
        setAccessToken: setAccessToken,
        refreshToken: refreshToken,
        setRefreshToken: setRefreshToken,
        createUser: createUser,
        authenticateUser: authenticateUser,
        refreshUserAuthentication: refreshUserAuthentication,
        logout: logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  )
}
