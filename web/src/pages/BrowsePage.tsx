import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link, useParams, useNavigate } from 'react-router-dom'
import { getBrowse, postChown, postMk } from '../api/client'
import ResultsTable from '../components/ResultsTable'
import ConnectionStatusBadge from '../components/ConnectionStatusBadge'
import type { TableDetailResponse } from '../types'

export default function BrowsePage() {
  const params = useParams()
  const navigate = useNavigate()
  const queryClient = useQueryClient()
  const path = params['*'] ?? ''
  const segments = path.split('/').filter(Boolean)
  const conn = segments[0] ?? ''

  const [more, setMore] = useState(false)
  const [newName, setNewName] = useState('')
  const [chownTarget, setChownTarget] = useState<string | null>(null)
  const [ownerInput, setOwnerInput] = useState('')

  const { data, isLoading, error } = useQuery({
    queryKey: ['browse', path, more],
    queryFn: () => getBrowse(path, more),
    enabled: !!path,
  })

  const mkMutation = useMutation({
    mutationFn: (name: string) => postMk(`${path}/${name}`),
    onSuccess: () => {
      setNewName('')
      queryClient.invalidateQueries({ queryKey: ['browse', path] })
    },
  })

  const chownMutation = useMutation({
    mutationFn: ({ target, owner }: { target: string; owner: string }) =>
      postChown(`${path}/${target}`, owner),
    onSuccess: () => {
      setChownTarget(null)
      setOwnerInput('')
      queryClient.invalidateQueries({ queryKey: ['browse', path] })
    },
  })

  if (!path) return <div className="text-sm text-slate-500">no connection selected</div>
  if (isLoading) return <div className="text-sm text-slate-400">loading...</div>
  if (error) return <div className="text-sm text-red-400">{(error as Error).message}</div>
  if (!data) return null

  const breadcrumbs = segments.map((seg, i) => ({
    label: seg,
    href: `/browse/${segments.slice(0, i + 1).join('/')}`,
  }))

  return (
    <div>
      <div className="mb-3 flex items-center justify-between">
        <nav className="flex items-center gap-1 text-sm text-slate-400">
          {breadcrumbs.map((b, i) => (
            <span key={b.href}>
              {i > 0 && <span className="mx-1 text-slate-600">/</span>}
              <Link to={b.href} className="hover:text-emerald-400">
                {b.label}
              </Link>
            </span>
          ))}
        </nav>
        {conn && <ConnectionStatusBadge conn={conn} />}
      </div>

      {data.kind === 'table' ? (
        <TableDetail path={path} data={data} more={more} setMore={setMore} />
      ) : (
        <>
          <div className="mb-3 flex items-center gap-4">
            <label className="flex items-center gap-2 text-sm text-slate-400">
              <input type="checkbox" checked={more} onChange={(e) => setMore(e.target.checked)} />
              show more details
            </label>
            {(data.kind === 'databases' || data.kind === 'schemas') && (
              <form
                className="flex items-center gap-2"
                onSubmit={(e) => {
                  e.preventDefault()
                  if (newName.trim()) mkMutation.mutate(newName.trim())
                }}
              >
                <input
                  value={newName}
                  onChange={(e) => setNewName(e.target.value)}
                  placeholder={data.kind === 'databases' ? 'new database name' : 'new schema name'}
                  className="rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
                />
                <button
                  type="submit"
                  disabled={mkMutation.isPending}
                  className="rounded bg-emerald-700 px-3 py-1 text-sm font-medium text-white hover:bg-emerald-600"
                >
                  create
                </button>
              </form>
            )}
          </div>

          {mkMutation.data && (
            <div className="mb-3 text-sm text-emerald-400">{mkMutation.data.messages.join(', ')}</div>
          )}
          {mkMutation.error && (
            <div className="mb-3 text-sm text-red-400">{(mkMutation.error as Error).message}</div>
          )}

          <ResultsTable
            result={data.result}
            onRowClick={(row) => navigate(`/browse/${path}/${row[0]}`)}
          />

          {(data.kind === 'databases' || data.kind === 'schemas') && (
            <div className="mt-3 flex flex-wrap gap-2">
              {data.result.rows.map((row) => (
                <button
                  key={row[0]}
                  onClick={() => {
                    setChownTarget(row[0])
                    setOwnerInput('')
                  }}
                  className="rounded border border-slate-800 px-2 py-1 text-xs text-slate-400 hover:border-slate-600 hover:text-slate-200"
                >
                  chown {row[0]}
                </button>
              ))}
            </div>
          )}

          {chownTarget && (
            <form
              className="mt-3 flex items-center gap-2"
              onSubmit={(e) => {
                e.preventDefault()
                if (ownerInput.trim()) chownMutation.mutate({ target: chownTarget, owner: ownerInput.trim() })
              }}
            >
              <span className="text-sm text-slate-400">new owner for {chownTarget}:</span>
              <input
                value={ownerInput}
                onChange={(e) => setOwnerInput(e.target.value)}
                autoFocus
                className="rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
              />
              <button
                type="submit"
                disabled={chownMutation.isPending}
                className="rounded bg-emerald-700 px-3 py-1 text-sm font-medium text-white hover:bg-emerald-600"
              >
                set owner
              </button>
              <button
                type="button"
                onClick={() => setChownTarget(null)}
                className="text-sm text-slate-500 hover:text-slate-300"
              >
                cancel
              </button>
            </form>
          )}
          {chownMutation.error && (
            <div className="mt-2 text-sm text-red-400">{(chownMutation.error as Error).message}</div>
          )}
        </>
      )}
    </div>
  )
}

function TableDetail({
  path,
  data,
  more,
  setMore,
}: {
  path: string
  data: TableDetailResponse
  more: boolean
  setMore: (v: boolean) => void
}) {
  return (
    <div>
      <div className="mb-3 flex items-center gap-4">
        <label className="flex items-center gap-2 text-sm text-slate-400">
          <input type="checkbox" checked={more} onChange={(e) => setMore(e.target.checked)} />
          show indexes &amp; constraints
        </label>
        <Link
          to={`/search/${path}`}
          className="rounded bg-slate-800 px-3 py-1 text-sm text-slate-200 hover:bg-slate-700"
        >
          search
        </Link>
        <Link
          to={`/sql/${path}`}
          className="rounded bg-slate-800 px-3 py-1 text-sm text-slate-200 hover:bg-slate-700"
        >
          open in SQL editor
        </Link>
      </div>

      <h2 className="mb-1 text-sm font-medium text-slate-400">Columns</h2>
      <ResultsTable result={data.columns} />

      {data.indexes && (
        <>
          <h2 className="mb-1 mt-4 text-sm font-medium text-slate-400">Indexes</h2>
          <ResultsTable result={data.indexes} />
        </>
      )}
      {data.constraints && (
        <>
          <h2 className="mb-1 mt-4 text-sm font-medium text-slate-400">Constraints</h2>
          <ResultsTable result={data.constraints} />
        </>
      )}
    </div>
  )
}
