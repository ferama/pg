import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getConnections } from '../api/client'
import ConnectionStatusBadge from '../components/ConnectionStatusBadge'

export default function ConnectionsPage() {
  const { data, isLoading, error } = useQuery({
    queryKey: ['connections'],
    queryFn: getConnections,
  })

  if (isLoading) return <div className="text-sm text-slate-400">loading...</div>
  if (error) return <div className="text-sm text-red-400">{(error as Error).message}</div>

  return (
    <div className="mx-auto max-w-2xl">
      <h1 className="mb-4 text-xl font-semibold">Connections</h1>
      <ul className="space-y-2">
        {data?.map((conn) => (
          <li
            key={conn}
            className="flex items-center justify-between rounded border border-slate-800 bg-slate-900 px-4 py-3"
          >
            <Link to={`/browse/${conn}`} className="font-medium text-slate-200 hover:text-emerald-400">
              {conn}
            </Link>
            <div className="flex items-center gap-3">
              <ConnectionStatusBadge conn={conn} />
              <Link to={`/stats/${conn}`} className="text-xs text-slate-400 hover:text-emerald-400">
                stats
              </Link>
              <Link to={`/users/${conn}`} className="text-xs text-slate-400 hover:text-emerald-400">
                users
              </Link>
            </div>
          </li>
        ))}
        {data?.length === 0 && (
          <li className="text-sm text-slate-500">
            no connections configured — add some to ~/.pg/conf.yaml
          </li>
        )}
      </ul>
    </div>
  )
}
