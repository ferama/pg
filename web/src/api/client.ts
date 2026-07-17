import type {
  AnyBrowseResponse,
  AttrResult,
  AutocompleteResponse,
  ChownResponse,
  HistoryItem,
  MessageResponse,
  MkResponse,
  PingResponse,
  QueryResult,
  StatsResponse,
} from '../types'

async function request<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(path, {
    headers: init?.body ? { 'Content-Type': 'application/json' } : undefined,
    ...init,
  })
  if (!res.ok) {
    let message = res.statusText
    try {
      const body = await res.json()
      if (body?.error) message = body.error
    } catch {
      // ignore non-JSON error bodies
    }
    throw new Error(message)
  }
  if (res.status === 204) return undefined as T
  return res.json() as Promise<T>
}

export function getConnections(): Promise<string[]> {
  return request('/api/connections')
}

export function getPing(conn: string): Promise<PingResponse> {
  return request(`/api/ping?conn=${encodeURIComponent(conn)}`)
}

export function getBrowse(path: string, more: boolean): Promise<AnyBrowseResponse> {
  const params = new URLSearchParams({ path, more: more ? 'true' : 'false' })
  return request(`/api/browse?${params}`)
}

export function postMk(path: string): Promise<MkResponse> {
  return request('/api/mk', { method: 'POST', body: JSON.stringify({ path }) })
}

export function postChown(path: string, owner: string): Promise<ChownResponse> {
  return request('/api/chown', { method: 'POST', body: JSON.stringify({ path, owner }) })
}

export function getSearch(
  path: string,
  columns: string[],
  filters: string[],
  limit: number,
): Promise<QueryResult> {
  const params = new URLSearchParams({ path, limit: String(limit) })
  for (const c of columns) params.append('columns', c)
  for (const f of filters) params.append('filters', f)
  return request(`/api/search?${params}`)
}

export function postQuery(path: string, query: string): Promise<QueryResult> {
  return request('/api/query', { method: 'POST', body: JSON.stringify({ path, query }) })
}

export function getHistory(): Promise<HistoryItem[]> {
  return request('/api/history')
}

export function deleteHistory(idx: number): Promise<void> {
  return request(`/api/history/${idx}`, { method: 'DELETE' })
}

export function getAutocomplete(path: string): Promise<AutocompleteResponse> {
  return request(`/api/autocomplete?path=${encodeURIComponent(path)}`)
}

export function getStats(conn: string): Promise<StatsResponse> {
  return request(`/api/stats?conn=${encodeURIComponent(conn)}`)
}

export function getUsers(conn: string): Promise<QueryResult> {
  return request(`/api/users?conn=${encodeURIComponent(conn)}`)
}

export function postUser(conn: string, username: string, password: string): Promise<MessageResponse> {
  return request('/api/users', { method: 'POST', body: JSON.stringify({ conn, username, password }) })
}

export function deleteUser(username: string, conn: string): Promise<MessageResponse> {
  const params = new URLSearchParams({ conn, confirm: 'true' })
  return request(`/api/users/${encodeURIComponent(username)}?${params}`, { method: 'DELETE' })
}

export interface UserAttrRequest {
  conn: string
  password?: string
  superuser?: boolean
  nosuperuser?: boolean
  createdb?: boolean
  nocreatedb?: boolean
}

export function postUserAttr(username: string, body: UserAttrRequest): Promise<AttrResult[]> {
  return request(`/api/users/${encodeURIComponent(username)}/attr`, {
    method: 'POST',
    body: JSON.stringify(body),
  })
}
