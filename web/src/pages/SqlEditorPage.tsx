import { useMemo, useState } from 'react'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import CodeMirror from '@uiw/react-codemirror'
import { sql, PostgreSQL } from '@codemirror/lang-sql'
import { autocompletion, type CompletionContext, type CompletionSource } from '@codemirror/autocomplete'
import { getAutocomplete, postQuery } from '../api/client'
import ResultsTable from '../components/ResultsTable'
import RowDetailModal from '../components/RowDetailModal'
import HistoryPanel from '../components/HistoryPanel'
import ConnectionStatusBadge from '../components/ConnectionStatusBadge'

export default function SqlEditorPage() {
  const params = useParams()
  const path = params['*'] ?? ''
  const segments = path.split('/').filter(Boolean)
  const conn = segments[0] ?? ''
  const tableName = segments[3]

  const [query, setQuery] = useState(
    tableName ? `select *\nfrom ${tableName}\nlimit 10` : '',
  )
  const [selectedRow, setSelectedRow] = useState<string[] | null>(null)
  const [historyOpen, setHistoryOpen] = useState(false)

  const autocompleteQuery = useQuery({
    queryKey: ['autocomplete', path],
    queryFn: () => getAutocomplete(path),
    enabled: !!path,
  })

  const mutation = useMutation({
    mutationFn: (q: string) => postQuery(path, q),
  })

  const completionSource: CompletionSource = useMemo(() => {
    const data = autocompleteQuery.data
    return (context: CompletionContext) => {
      const word = context.matchBefore(/\w+/)
      if (!word || (word.from === word.to && !context.explicit)) return null
      if (!data) return null

      const options = [
        ...data.keywords.map((label) => ({ label, type: 'keyword' })),
        ...data.tables.map((label) => ({ label, type: 'class' })),
        ...data.columns.map((label) => ({ label, type: 'property' })),
      ]
      return { from: word.from, options }
    }
  }, [autocompleteQuery.data])

  const extensions = useMemo(
    () => [
      sql({ dialect: PostgreSQL }),
      autocompletion({ override: [completionSource] }),
    ],
    [completionSource],
  )

  const runQuery = () => {
    if (query.trim()) mutation.mutate(query)
  }

  return (
    <div className="flex h-full flex-col">
      <div className="mb-3 flex items-center justify-between">
        <h1 className="text-lg font-semibold">SQL Editor: {path}</h1>
        <div className="flex items-center gap-3">
          <ConnectionStatusBadge conn={conn} />
          <button
            onClick={() => setHistoryOpen(true)}
            className="rounded bg-slate-800 px-3 py-1 text-sm text-slate-200 hover:bg-slate-700"
          >
            history
          </button>
          <button
            onClick={runQuery}
            disabled={mutation.isPending}
            className="rounded bg-emerald-700 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-600"
          >
            run (Ctrl+Enter)
          </button>
        </div>
      </div>

      <div
        className="mb-4 overflow-hidden rounded border border-slate-800"
        onKeyDown={(e) => {
          if ((e.ctrlKey || e.metaKey) && e.key === 'Enter') {
            e.preventDefault()
            runQuery()
          }
        }}
      >
        <CodeMirror
          value={query}
          height="200px"
          theme="dark"
          extensions={extensions}
          onChange={setQuery}
        />
      </div>

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

      <HistoryPanel
        open={historyOpen}
        onClose={() => setHistoryOpen(false)}
        onSelect={(q) => {
          setQuery(q)
          setHistoryOpen(false)
        }}
      />
    </div>
  )
}
