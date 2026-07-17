import { useQuery } from '@tanstack/react-query'
import { getPing } from '../api/client'

export default function ConnectionStatusBadge({ conn }: { conn: string }) {
  const { data } = useQuery({
    queryKey: ['ping', conn],
    queryFn: () => getPing(conn),
    refetchInterval: 3000,
    enabled: !!conn,
  })

  const connected = data?.connected ?? false

  return (
    <span
      className={
        'inline-flex items-center gap-1.5 rounded-full px-2.5 py-0.5 text-xs font-medium ' +
        (connected ? 'bg-emerald-900 text-emerald-300' : 'bg-red-900 text-red-300')
      }
    >
      <span
        className={'h-1.5 w-1.5 rounded-full ' + (connected ? 'bg-emerald-400' : 'bg-red-400')}
      />
      {conn} · {connected ? 'connected' : 'disconnected'}
    </span>
  )
}
