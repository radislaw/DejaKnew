import { apiRequest } from '@/shared/api'
import type { User } from '@/entities/user'

export function updateProfileRequest(email: string, locale: string): Promise<User> {
  return apiRequest<User>('/profile', {
    method: 'PATCH',
    body: { email, locale },
  })
}

export function changePasswordRequest(current: string, newPassword: string): Promise<void> {
  return apiRequest<void>('/profile/password', {
    method: 'PUT',
    body: { current, new: newPassword },
  })
}

export function deleteAccountRequest(): Promise<void> {
  return apiRequest<void>('/profile', { method: 'DELETE' })
}
