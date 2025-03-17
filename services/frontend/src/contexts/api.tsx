import { createContext, useContext } from "react";
import { useAuthContext } from "./auth";
import { GetMeResponse, User } from "@/models/user";
import { executeAPIRequest } from "@/services/api";
import { CreateOrganizationRequest, CreateOrganizationResponse, GetOrganizationMembershipsResponse, OrganizationMembership } from "@/models/organization";
import { APIError } from "@/services/errors";
import { useMutation, UseMutationResult, useQuery, useQueryClient, UseQueryResult } from "@tanstack/react-query";

interface APIContextType {
    useGetUserInfo(): UseQueryResult<User, APIError>;
    useGetOrganizationMemberships(): UseQueryResult<Array<OrganizationMembership>, APIError>;
    useCreateOrganization: UseMutationResult<CreateOrganizationResponse, APIError, CreateOrganizationRequest, unknown>;
}

const APIContext = createContext<APIContextType | null>(null);

let ongoingRefreshTokenRequest: Promise<string> | null = null;

function APIProvider({ children } : {children: React.ReactNode}) {
    const {accessToken, refreshUserAuthentication} = useAuthContext();
    const queryClient = useQueryClient()

    const baseRequestParams = {
        accessToken: accessToken,
        onTokenExpired: async () => {
            if (ongoingRefreshTokenRequest) {
                return ongoingRefreshTokenRequest;
            }
            ongoingRefreshTokenRequest = (async () => {
                try {
                    const response = await refreshUserAuthentication.mutateAsync(undefined);
                    return response.token;
                } finally {
                    ongoingRefreshTokenRequest = null;
                }
            })();
            return ongoingRefreshTokenRequest;
        }
    };

    const useGetUserInfo = (): UseQueryResult<User, APIError> => {
        return useQuery({
            queryFn: async () => {
                const response = await executeAPIRequest<null, GetMeResponse>({
                    endpoint: "users/me",
                    method: "GET",
                    ...baseRequestParams,
                });
                return response.user;
            },
            queryKey: ["users.me"],
            enabled: !!accessToken,
            retry: false,
        });
    };

    const useGetOrganizationMemberships = (): UseQueryResult<Array<OrganizationMembership>, APIError> => {
        return useQuery({
            queryFn: async () => {
                const response = await executeAPIRequest<null, GetOrganizationMembershipsResponse>({
                    endpoint: "users/me/organizations",
                    method: "GET",
                    ...baseRequestParams,
                });
                return response.data;
            },
            queryKey: ["users.me.organizations"],
            enabled: !!accessToken,
            retry: false,
        });
    };

    const useCreateOrganization = useMutation<CreateOrganizationResponse, APIError, CreateOrganizationRequest, unknown>({
      mutationFn: async (request) => {
        return executeAPIRequest<CreateOrganizationRequest, CreateOrganizationResponse>({
            endpoint: "organizations",
            method: "POST",
            body: request,
            ...baseRequestParams,
        });
      },
      onSuccess: async () => {
        await queryClient.invalidateQueries({queryKey: ["users.me.organizations"]});
      }
    });

    return (
        <APIContext.Provider
            value={{
                useGetUserInfo: useGetUserInfo,
                useGetOrganizationMemberships: useGetOrganizationMemberships,
                useCreateOrganization: useCreateOrganization,
            }}
        >
            {children}
        </APIContext.Provider>
    );
}

function useAPIContext(): APIContextType {
    const context = useContext(APIContext);
    if (!context) {
      throw new Error("useAPIContext() must be used within the APIProvider");
    }
    return context;
};

export {APIProvider, useAPIContext};