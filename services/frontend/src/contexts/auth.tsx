import { IssueAuthenticationTokenRequest, IssueAuthenticationTokenResponse, RefreshAuthenticationTokenRequest, RefreshAuthenticationTokenResponse } from "@/models/auth";
import { CreateUserRequest, CreateUserResponse } from "@/models/user";
import { executeAPIRequest } from "@/services/api";
import { APIError } from "@/services/errors";
import { useMutation, UseMutationResult } from "@tanstack/react-query";
import { useLocalStorage } from "@uidotdev/usehooks";
import { createContext, useContext } from "react";

const ACCESS_TOKEN_STORAGE_KEY = "access_token";
const REFRESH_TOKEN_STORAGE_KEY = "refresh_token";

interface AuthContextType {
    accessToken: string | undefined;
    setAccessToken(token: string): void;

    refreshToken: string | undefined;
    setRefreshToken(token: string): void;

    createUser: UseMutationResult<CreateUserResponse, APIError, CreateUserRequest, unknown>;
    authenticateUser: UseMutationResult<IssueAuthenticationTokenResponse, APIError, IssueAuthenticationTokenRequest, unknown>;
    refreshUserAuthentication: UseMutationResult<RefreshAuthenticationTokenResponse, APIError, undefined, unknown>;
    logout(): Promise<void>;
}

const AuthContext = createContext<AuthContextType | null>(null);

function AuthProvider({ children } : {children: React.ReactNode}) {
    const [accessToken, setAccessToken] = useLocalStorage<string | undefined>(ACCESS_TOKEN_STORAGE_KEY)
    const [refreshToken, setRefreshToken] = useLocalStorage<string | undefined>(REFRESH_TOKEN_STORAGE_KEY)

    const createUser = useMutation<CreateUserResponse, APIError, CreateUserRequest, unknown>({
      mutationFn: async (request) => {
        return executeAPIRequest<CreateUserRequest, CreateUserResponse>({
          endpoint: "users",
          method: "POST",
          body: request,
        });
      }
    });
 
    const authenticateUser = useMutation<IssueAuthenticationTokenResponse, APIError, IssueAuthenticationTokenRequest, unknown>({
      mutationFn: async (request) => {
        return executeAPIRequest<IssueAuthenticationTokenRequest, IssueAuthenticationTokenResponse>({
            endpoint: "tokens/issue",
            method: "POST",
            body: request,
        });
      },
      onSuccess: (response) => {
        setAccessToken(response.token);
        setRefreshToken(response.refreshToken);
      }
    });

    const refreshUserAuthentication = useMutation<RefreshAuthenticationTokenResponse, APIError, undefined, unknown>({
      mutationFn: async () => {
        if (!refreshToken) {
          throw new Error("refreshToken doesn't exist in refreshUserAuthentication");
        }
        return executeAPIRequest<RefreshAuthenticationTokenRequest, RefreshAuthenticationTokenResponse>({
            endpoint: "tokens/refresh",
            method: "POST",
            body: {refreshToken},
        });
      },
      onSuccess: (response) => {
        setAccessToken(response.token);
      },
      onError: async () => {await logout()},
    });

    const logout = async () => {
        setAccessToken(undefined);
        setRefreshToken(undefined);
        window.location.reload();
    };

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
    );
}

function useAuthContext(): AuthContextType {
    const context = useContext(AuthContext);
    if (!context) {
      throw new Error("useAuthContext() must be used within the AuthProvider");
    }
    return context;
};

export {AuthProvider, useAuthContext};