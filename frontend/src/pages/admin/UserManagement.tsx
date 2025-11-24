import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../../components/Card'
import Button from '../../components/Button'
import Loading from '../../components/Loading'
import './UserManagement.css'

interface User {
  id: number
  username: string
  email: string
  nickname?: string
  status: string
  email_verified: boolean
  mfa_enabled: boolean
  created_at: string
  roles: Array<{ id: number; name: string; display_name: string }>
}

export default function UserManagement() {
  const navigate = useNavigate()
  const [users, setUsers] = useState<User[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [page, setPage] = useState(1)
  const [totalPages, setTotalPages] = useState(1)
  const [search, setSearch] = useState('')

  useEffect(() => {
    fetchUsers()
  }, [page, search])

  const fetchUsers = async () => {
    try {
      setLoading(true)
      const params = new URLSearchParams({
        page: page.toString(),
        page_size: '20',
      })
      if (search) {
        params.append('search', search)
      }
      const response = await axios.get(`/api/admin/users?${params}`)
      setUsers(response.data.data.users)
      setTotalPages(response.data.data.pagination.total_pages)
    } catch (err: any) {
      setError(err.response?.data?.message || '获取用户列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleDelete = async (userId: number) => {
    if (!confirm('确定要删除这个用户吗？')) return

    try {
      await axios.delete(`/api/admin/users/${userId}`)
      fetchUsers()
    } catch (err: any) {
      alert(err.response?.data?.message || '删除失败')
    }
  }

  const handleStatusChange = async (userId: number, newStatus: string) => {
    try {
      await axios.put(`/api/admin/users/${userId}`, { status: newStatus })
      fetchUsers()
    } catch (err: any) {
      alert(err.response?.data?.message || '更新失败')
    }
  }

  if (loading && users.length === 0) {
    return (
      <div className="user-management">
        <Loading text="加载中..." />
      </div>
    )
  }

  return (
    <div className="user-management">
      <div className="page-header">
        <h2 className="page-title">👥 用户管理</h2>
        <Button variant="primary" onClick={() => navigate('/admin/users/new')}>
          ➕ 添加用户
        </Button>
      </div>

      <Card className="search-card">
        <input
          type="text"
          placeholder="搜索用户名、邮箱或昵称..."
          value={search}
          onChange={(e) => {
            setSearch(e.target.value)
            setPage(1)
          }}
          className="search-input"
        />
      </Card>

      {error && <div className="error-message">{error}</div>}

      <Card className="users-table-card">
        <table className="users-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>用户名</th>
              <th>邮箱</th>
              <th>昵称</th>
              <th>状态</th>
              <th>角色</th>
              <th>MFA</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {users.map((user) => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>{user.nickname || '-'}</td>
                <td>
                  <select
                    value={user.status}
                    onChange={(e) => handleStatusChange(user.id, e.target.value)}
                    className="status-select"
                  >
                    <option value="active">活跃</option>
                    <option value="suspended">暂停</option>
                  </select>
                </td>
                <td>
                  <div className="roles-tags">
                    {user.roles.map((role) => (
                      <span key={role.id} className="role-tag">
                        {role.display_name || role.name}
                      </span>
                    ))}
                  </div>
                </td>
                <td>{user.mfa_enabled ? '✅' : '❌'}</td>
                <td>
                  <div className="action-buttons">
                    <Button
                      variant="secondary"
                      size="small"
                      onClick={() => navigate(`/admin/users/${user.id}`)}
                    >
                      编辑
                    </Button>
                    <Button
                      variant="outline"
                      size="small"
                      onClick={() => handleDelete(user.id)}
                    >
                      删除
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {users.length === 0 && !loading && (
          <div className="empty-state">暂无用户</div>
        )}

        {totalPages > 1 && (
          <div className="pagination">
            <Button
              variant="outline"
              onClick={() => setPage(page - 1)}
              disabled={page === 1}
            >
              上一页
            </Button>
            <span className="page-info">
              第 {page} / {totalPages} 页
            </span>
            <Button
              variant="outline"
              onClick={() => setPage(page + 1)}
              disabled={page >= totalPages}
            >
              下一页
            </Button>
          </div>
        )}
      </Card>
    </div>
  )
}


