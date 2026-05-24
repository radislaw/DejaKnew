import { API_BASE_URL } from '@/shared/config'

// The auth token lives here so the client can attach it to every request.
// The auth store (a feature) calls setAuthToken — shared never imports up.
let authToken: string | null = null

export function setAuthToken(token: string | null): void {
  authToken = token
}

export class ApiError extends Error {
  constructor(
    public readonly status: number,
    message: string,
  ) {
    super(message)
    this.name = 'ApiError'
  }
}

interface RequestOptions {
  method?: string
  body?: unknown
}

export async function apiRequest<T>(path: string, options: RequestOptions = {}): Promise<T> {
  const headers: Record<string, string> = { 'Content-Type': 'application/json' }
  if (authToken) headers.Authorization = `Bearer ${authToken}`

  const response = await fetch(`${API_BASE_URL}${path}`, {
    method: options.method ?? 'GET',
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  })

  if (response.status === 204) return undefined as T

  const data = await response.json().catch(() => null)
  if (!response.ok) {
    const message =
      data && typeof data.error === 'string'
        ? data.error
        : `Request failed (${response.status})`
    throw new ApiError(response.status, message)
  }
  return data as T
}

export async function apiBlobRequest(path: string): Promise<{ blob: Blob; filename: string }> {
  const headers: Record<string, string> = {}
  if (authToken) headers.Authorization = `Bearer ${authToken}`

  const response = await fetch(`${API_BASE_URL}${path}`, { headers })

  if (!response.ok) {
    const data = await response.json().catch(() => null)
    const message =
      data && typeof data.error === 'string'
        ? data.error
        : `Request failed (${response.status})`
    throw new ApiError(response.status, message)
  }

  const contentDisposition = response.headers.get('Content-Disposition') ?? ''
  const match = contentDisposition.match(/filename=(.+)/)
  const filename = match ? match[1].trim() : 'export'

  const blob = await response.blob()
  return { blob, filename }
}
