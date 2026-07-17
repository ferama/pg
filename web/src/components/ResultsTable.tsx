import type { QueryResult } from '../types'

interface Props {
  result: QueryResult
  onRowClick?: (row: string[]) => void
}

export default function ResultsTable({ result, onRowClick }: Props) {
  if (result.message) {
    return (
      <div className="text-sm text-slate-400">
        {result.message}
        {result.elapsedMs > 0 && <span className="ml-2 text-slate-500">[{result.elapsedMs}ms]</span>}
      </div>
    )
  }

  if (result.rows.length === 0) {
    return <div className="text-sm text-slate-500">no rows</div>
  }

  return (
    <div className="overflow-auto rounded border border-slate-800">
      <div className="border-b border-slate-800 bg-slate-900 px-3 py-1.5 text-xs text-slate-500">
        {result.rows.length} row{result.rows.length === 1 ? '' : 's'} · {result.elapsedMs}ms
      </div>
      <table className="w-full text-left text-sm">
        <thead className="bg-slate-900 text-slate-400">
          <tr>
            {result.columns.map((c) => (
              <th key={c} className="whitespace-nowrap px-3 py-2 font-medium">
                {c}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {result.rows.map((row, i) => (
            <tr
              key={i}
              onClick={() => onRowClick?.(row)}
              className={
                'border-t border-slate-800 ' +
                (onRowClick ? 'cursor-pointer hover:bg-slate-800/60' : '')
              }
            >
              {row.map((cell, j) => (
                <td key={j} className="max-w-xs truncate whitespace-nowrap px-3 py-1.5 text-slate-300">
                  {cell}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
