export interface QueryResult {
  columns: string[]
  rows: string[][]
  elapsedMs: number
  message?: string
}

export interface BrowseResponse {
  kind: 'databases' | 'schemas' | 'tables'
  result: QueryResult
}

export interface TableDetailResponse {
  kind: 'table'
  columns: QueryResult
  indexes?: QueryResult
  constraints?: QueryResult
}

export type AnyBrowseResponse = BrowseResponse | TableDetailResponse

export interface PingResponse {
  connected: boolean
}

export interface MkResponse {
  messages: string[]
}

export interface ChownResponse {
  messages: string[]
}

export interface HistoryItem {
  idx: number
  query: string
}

export interface AutocompleteResponse {
  keywords: string[]
  tables: string[]
  columns: string[]
}

export interface StatsResponse {
  maxConnections: string
  currentConnections: string
}

export interface AttrResult {
  attribute: string
  success: boolean
  message: string
}

export interface MessageResponse {
  message: string
}
