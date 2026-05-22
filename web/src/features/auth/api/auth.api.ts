import { apiRequest } from '@/shared/api'
import type { User } from '@/entities/user'

export interface AuthResponse {
  token: string
  user: User
}

export function registerRequest(email: string, password: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/auth/register', {
    method: 'POST',
    body: { email, password },
  })
}

export function loginRequest(email: string, password: string): Promise<AuthResponse> {
  return apiRequest<AuthResponse>('/auth/login', {
    method: 'POST',
    body: { email, password },
  })
}

export function meRequest(): Promise<User> {
  return apiRequest<User>('/me')
}
