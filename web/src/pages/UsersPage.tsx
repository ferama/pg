import { useState } from 'react'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { useParams } from 'react-router-dom'
import { deleteUser, getUsers, postUser, postUserAttr } from '../api/client'
import ConnectionStatusBadge from '../components/ConnectionStatusBadge'
import ConfirmDialog from '../components/ConfirmDialog'
import type { AttrResult } from '../types'

export default function UsersPage() {
  const { conn = '' } = useParams()
  const queryClient = useQueryClient()

  const [newUsername, setNewUsername] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [attrTarget, setAttrTarget] = useState<string | null>(null)
  const [deleteTarget, setDeleteTarget] = useState<string | null>(null)
  const [attrResults, setAttrResults] = useState<AttrResult[] | null>(null)

  const usersQuery = useQuery({
    queryKey: ['users', conn],
    queryFn: () => getUsers(conn),
    enabled: !!conn,
  })

  const invalidateUsers = () => queryClient.invalidateQueries({ queryKey: ['users', conn] })

  const addMutation = useMutation({
    mutationFn: () => postUser(conn, newUsername.trim(), newPassword),
    onSuccess: () => {
      setNewUsername('')
      setNewPassword('')
      invalidateUsers()
    },
  })

  const deleteMutation = useMutation({
    mutationFn: (username: string) => deleteUser(username, conn),
    onSuccess: () => {
      setDeleteTarget(null)
      invalidateUsers()
    },
  })

  const usernameIdx = usersQuery.data?.columns.findIndex((c) => c.toUpperCase() === 'USERNAME') ?? 0

  return (
    <div className="mx-auto max-w-3xl">
      <div className="mb-4 flex items-center justify-between">
        <h1 className="text-lg font-semibold">Users: {conn}</h1>
        <ConnectionStatusBadge conn={conn} />
      </div>

      <form
        className="mb-6 flex flex-wrap items-end gap-2"
        onSubmit={(e) => {
          e.preventDefault()
          if (newUsername.trim()) addMutation.mutate()
        }}
      >
        <div className="flex flex-col gap-1">
          <label className="text-xs text-slate-500">username</label>
          <input
            value={newUsername}
            onChange={(e) => setNewUsername(e.target.value)}
            className="rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
          />
        </div>
        <div className="flex flex-col gap-1">
          <label className="text-xs text-slate-500">password (optional)</label>
          <input
            type="password"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            className="rounded border border-slate-700 bg-slate-900 px-2 py-1 text-sm"
          />
        </div>
        <button
          type="submit"
          disabled={addMutation.isPending}
          className="rounded bg-emerald-700 px-4 py-1.5 text-sm font-medium text-white hover:bg-emerald-600"
        >
          add user
        </button>
      </form>
      {addMutation.error && <div className="mb-3 text-sm text-red-400">{(addMutation.error as Error).message}</div>}
      {addMutation.data && <div className="mb-3 text-sm text-emerald-400">{addMutation.data.message}</div>}

      {usersQuery.isLoading && <div className="text-sm text-slate-400">loading...</div>}
      {usersQuery.error && <div className="text-sm text-red-400">{(usersQuery.error as Error).message}</div>}

      {usersQuery.data && (
        <div className="overflow-auto rounded border border-slate-800">
          <table className="w-full text-left text-sm">
            <thead className="bg-slate-900 text-slate-400">
              <tr>
                {usersQuery.data.columns.map((c) => (
                  <th key={c} className="whitespace-nowrap px-3 py-2 font-medium">
                    {c}
                  </th>
                ))}
                <th className="px-3 py-2 font-medium">actions</th>
              </tr>
            </thead>
            <tbody>
              {usersQuery.data.rows.map((row, i) => (
                <tr key={i} className="border-t border-slate-800">
                  {row.map((cell, j) => (
                    <td key={j} className="whitespace-nowrap px-3 py-1.5 text-slate-300">
                      {cell}
                    </td>
                  ))}
                  <td className="whitespace-nowrap px-3 py-1.5">
                    <button
                      onClick={() => {
                        setAttrTarget(row[usernameIdx])
                        setAttrResults(null)
                      }}
                      className="mr-2 text-xs text-slate-400 hover:text-emerald-400"
                    >
                      attr
                    </button>
                    <button
                      onClick={() => setDeleteTarget(row[usernameIdx])}
                      className="text-xs text-slate-400 hover:text-red-400"
                    >
                      delete
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {attrTarget && (
        <AttrForm
          conn={conn}
          username={attrTarget}
          onClose={() => setAttrTarget(null)}
          onDone={(results) => {
            setAttrResults(results)
            invalidateUsers()
          }}
        />
      )}
      {attrResults && (
        <div className="mt-4 space-y-1 rounded border border-slate-800 bg-slate-900 p-3">
          {attrResults.map((r) => (
            <div key={r.attribute} className={r.success ? 'text-sm text-emerald-400' : 'text-sm text-red-400'}>
              {r.attribute}: {r.message}
            </div>
          ))}
        </div>
      )}

      {deleteTarget && (
        <ConfirmDialog
          message={`I'm going to drop user '${deleteTarget}'. Proceed?`}
          onCancel={() => setDeleteTarget(null)}
          onConfirm={() => deleteMutation.mutate(deleteTarget)}
        />
      )}
      {deleteMutation.error && (
        <div className="mt-3 text-sm text-red-400">{(deleteMutation.error as Error).message}</div>
      )}
    </div>
  )
}

function AttrForm({
  conn,
  username,
  onClose,
  onDone,
}: {
  conn: string
  username: string
  onClose: () => void
  onDone: (results: AttrResult[]) => void
}) {
  const [password, setPassword] = useState('')
  const [superuser, setSuperuser] = useState(false)
  const [nosuperuser, setNosuperuser] = useState(false)
  const [createdb, setCreatedb] = useState(false)
  const [nocreatedb, setNocreatedb] = useState(false)

  const mutation = useMutation({
    mutationFn: () =>
      postUserAttr(username, {
        conn,
        password: password || undefined,
        superuser: superuser || undefined,
        nosuperuser: nosuperuser || undefined,
        createdb: createdb || undefined,
        nocreatedb: nocreatedb || undefined,
      }),
    onSuccess: (results) => {
      onDone(results)
      onClose()
    },
  })

  return (
    <form
      className="mt-4 flex flex-wrap items-end gap-3 rounded border border-slate-800 bg-slate-900 p-3"
      onSubmit={(e) => {
        e.preventDefault()
        mutation.mutate()
      }}
    >
      <span className="text-sm font-medium text-slate-300">attrs for {username}</span>
      <input
        type="password"
        placeholder="new password"
        value={password}
        onChange={(e) => setPassword(e.target.value)}
        className="rounded border border-slate-700 bg-slate-950 px-2 py-1 text-sm"
      />
      <label className="flex items-center gap-1 text-xs text-slate-400">
        <input type="checkbox" checked={superuser} onChange={(e) => setSuperuser(e.target.checked)} />
        superuser
      </label>
      <label className="flex items-center gap-1 text-xs text-slate-400">
        <input type="checkbox" checked={nosuperuser} onChange={(e) => setNosuperuser(e.target.checked)} />
        nosuperuser
      </label>
      <label className="flex items-center gap-1 text-xs text-slate-400">
        <input type="checkbox" checked={createdb} onChange={(e) => setCreatedb(e.target.checked)} />
        createdb
      </label>
      <label className="flex items-center gap-1 text-xs text-slate-400">
        <input type="checkbox" checked={nocreatedb} onChange={(e) => setNocreatedb(e.target.checked)} />
        nocreatedb
      </label>
      <button
        type="submit"
        disabled={mutation.isPending}
        className="rounded bg-emerald-700 px-3 py-1 text-sm font-medium text-white hover:bg-emerald-600"
      >
        apply
      </button>
      <button type="button" onClick={onClose} className="text-sm text-slate-500 hover:text-slate-300">
        cancel
      </button>
    </form>
  )
}
