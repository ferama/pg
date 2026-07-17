interface Props {
  columns: string[]
  row: string[]
  onClose: () => void
}

export default function RowDetailModal({ columns, row, onClose }: Props) {
  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4"
      onClick={onClose}
    >
      <div
        className="max-h-[80vh] w-full max-w-2xl overflow-auto rounded border border-slate-700 bg-slate-900 p-4 shadow-xl"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-sm font-semibold text-emerald-400">Row Detail</h2>
          <button onClick={onClose} className="text-slate-400 hover:text-slate-200">
            close
          </button>
        </div>
        <dl className="divide-y divide-slate-800 text-sm">
          {columns.map((col, i) => (
            <div key={col} className="grid grid-cols-3 gap-2 py-2">
              <dt className="font-medium text-slate-400">{col}</dt>
              <dd className="col-span-2 break-all text-slate-200">{row[i]}</dd>
            </div>
          ))}
        </dl>
      </div>
    </div>
  )
}
