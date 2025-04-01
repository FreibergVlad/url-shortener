export interface User {
  id: string
  email: string
  fullName: string
}

export interface ShortUser {
  id: string
  email: string
}

export interface GetMeResponse {
  user: User
}

export interface CreateUserRequest {
  email: string
  password: string
  fullName: string
}

export interface CreateUserResponse {
  user: User
}
