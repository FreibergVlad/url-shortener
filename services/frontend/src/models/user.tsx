export interface User {
    id: string;
    email: string;
    firstName: string;
    lastName: string;
}

export interface GetMeResponse {
    user: User;
}

export interface CreateUserRequest {
    email: string;
    password: string;
}

export interface CreateUserResponse {
    user: User;
}