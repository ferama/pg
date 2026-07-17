import { useQuery, useQueryClient } from '@tanstack/react-query'
import { deleteHistory, getHistory } from '../api/client'

interface Props {
  open: boolean
  onClose: () => void
  onSelect: (query: string) => void
}

export default function HistoryPanel({ open, onClose, onSelect }: Props) {
  const queryClient = useQueryClient()
  const { data } = useQuery({
    queryKey: ['history'],
    queryFn: getHistory,
    enabled: open,
  })

  if (!open) return null

  const handleDelete = async (idx: number) => {
    await deleteHistory(idx)
    queryClient.invalidateQueries({ queryKey: ['history'] })
  }

  return (
    <div className="fixed inset-0 z-40 flex justify-end bg-black/50" onClick={onClose}>
      <div
        className="h-full w-full max-w-md overflow-auto border-l border-slate-700 bg-slate-900 p-4"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-sm font-semibold text-emerald-400">Query History</h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-200">
            close
          </button>
        </div>
        {(!data || data.length === 0) && (
          <div className="text-sm text-slate-500">no history yet</div>
        )}
        <ul className="space-y-1">
          {data?.map((item) => (
            <li
              key={item.idx}
              className="group flex items-start justify-between gap-2 rounded px-2 py-1.5 hover:bg-slate-800"
            >
              <button
                onClick={() => onSelect(item.query)}
                className="flex-1 truncate whitespace-pre-wrap text-left text-xs text-slate-300"
                title={item.query}
              >
                {item.query}
              </button>
              <button
                onClick={() => handleDelete(item.idx)}
                className="text-xs text-slate-500 opacity-0 hover:text-red-400 group-hover:opacity-100"
              >
                delete
              </button>
            </li>
          ))}
        </ul>
      </div>
    </div>
  )
}
