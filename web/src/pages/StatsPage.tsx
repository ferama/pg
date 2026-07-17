import { useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { getStats } from '../api/client'
import ConnectionStatusBadge from '../components/ConnectionStatusBadge'

export default function StatsPage() {
  const { conn = '' } = useParams()

  const { data, isLoading, error } = useQuery({
    queryKey: ['stats', conn],
    queryFn: () => getStats(conn),
    enabled: !!conn,
  })

  return (
    <div className="mx-auto max-w-md">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold">Stats: {conn}</h1>
        <ConnectionStatusBadge conn={conn} />
      </div>

      {isLoading && <div className="text-sm text-slate-400">loading...</div>}
      {error && <div className="text-sm text-red-400">{(error as Error).message}</div>}

      {data && (
        <div className="grid grid-cols-2 gap-3">
          <div className="rounded border border-slate-800 bg-slate-900 p-4">
            <div className="text-xs text-slate-500">MAX CONNECTIONS</div>
            <div className="text-2xl font-semibold text-slate-100">{data.maxConnections}</div>
          </div>
          <div className="rounded border border-slate-800 bg-slate-900 p-4">
            <div className="text-xs text-slate-500">CURRENT CONNECTIONS</div>
            <div className="text-2xl font-semibold text-slate-100">{data.currentConnections}</div>
          </div>
        </div>
      )}
    </div>
  )
}
