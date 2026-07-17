import { Route, Routes } from 'react-router-dom'
import Layout from './components/Layout'
import ConnectionsPage from './pages/ConnectionsPage'
import BrowsePage from './pages/BrowsePage'
import SearchPage from './pages/SearchPage'
import SqlEditorPage from './pages/SqlEditorPage'
import StatsPage from './pages/StatsPage'
import UsersPage from './pages/UsersPage'

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<ConnectionsPage />} />
        <Route path="/browse/*" element={<BrowsePage />} />
        <Route path="/search/*" element={<SearchPage />} />
        <Route path="/sql/*" element={<SqlEditorPage />} />
        <Route path="/stats/:conn" element={<StatsPage />} />
        <Route path="/users/:conn" element={<UsersPage />} />
      </Route>
    </Routes>
  )
}
