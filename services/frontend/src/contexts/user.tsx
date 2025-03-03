import { User } from "@/models/user";
import { createContext, useContext } from "react";
import { useAPIContext } from "./api";
import { OrganizationMembership } from "@/models/organization";

interface UserContextType {
    user: User | undefined;
    organizationMemberships: Array<OrganizationMembership> | undefined
    isLoading: boolean | undefined
}

const UserContext = createContext<UserContextType | null>(null);

function UserProvider({ children } : {children: React.ReactNode}) {
    const {useGetUserInfo, useGetOrganizationMemberships} = useAPIContext();

    const {data: user, ...userQuery} = useGetUserInfo();
    const {data: organizationMemberships, ...organizationMembershipsQuery} = useGetOrganizationMemberships();

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
    );
}

function useUserContext(): UserContextType {
    const context = useContext(UserContext);
    if (!context) {
      throw new Error("useUserContext() must be used within the UserProvider");
    }
    return context;
};

export {UserProvider, useUserContext};