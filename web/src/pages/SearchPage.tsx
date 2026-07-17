import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { getSearch } from '../api/client'
import ResultsTable from '../components/ResultsTable'
import RowDetailModal from '../components/RowDetailModal'

export default function SearchPage() {
  const params = useParams()
  const path = params['*'] ?? ''

  const [columns, setColumns] = useState('')
  const [filters, setFilters] = useState('')
  const [limit, setLimit] = useState(10)
  const [selectedRow, setSelectedRow] = useState<string[] | null>(null)

  const mutation = useMutation({
    mutationFn: () => {
      const cols = columns.split(',').map((c) => c.trim()).filter(Boolean)
      const fils = filters.split(',').map((f) => f.trim()).filter(Boolean)
      return getSearch(path, cols, fils, limit)
    },
  })

  return (
    <div>
      <h1 className="mb-3 text-lg font-semibold">Search: {path}</h1>

      <form
        className="mb-4 flex flex-wrap items-end gap-3"
        onSubmit={(e) => {
          e.preventDefault()
          mutation.mutate()
        }}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-slate-500">columns (comma separated, empty = all)</label>
          <input
            value={columns}
            onChange={(e) => setColumns(e.target.value)}
            placeholder="age,sex,city"
            className="w-56 rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs text-slate-500">filters (comma separated)</label>
          <input
            value={filters}
            onChange={(e) => setFilters(e.target.value)}
            placeholder="sex=M,name like test%"
            className="w-64 rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs text-slate-500">limit</label>
          <input
            type="number"
            value={limit}
            onChange={(e) => setLimit(Number(e.target.value))}
            className="w-24 rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
          />
        </div>
        <button
          type="submit"
          disabled={mutation.isPending}
          className="rounded bg-emerald-700 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-600"
        >
          run
        </button>
      </form>

      {mutation.error && <div className="mb-3 text-sm text-red-400">{(mutation.error as Error).message}</div>}
      {mutation.data && (
        <ResultsTable result={mutation.data} onRowClick={(row) => setSelectedRow(row)} />
      )}

      {selectedRow && mutation.data && (
        <RowDetailModal
          columns={mutation.data.columns}
          row={selectedRow}
          onClose={() => setSelectedRow(null)}
        />
      )}
    </div>
  )
}
