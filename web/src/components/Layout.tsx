import { Link, Outlet } from 'react-router-dom'

export default function Layout() {
  return (
    <div className="flex min-h-full flex-col">
      <header className="flex items-center gap-4 border-b border-slate-800 bg-slate-900 px-4 py-3">
        <Link to="/" className="text-lg font-semibold text-emerald-400">
          pg
        </Link>
        <span className="text-sm text-slate-500">web UI</span>
      </header>
      <main className="flex-1 overflow-auto p-4">
        <Outlet />
      </main>
    </div>
  )
}
