import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Loading from '../components/Loading'
import './Permissions.css'

interface Role {
  id: number
  name: string
  display_name: string
  description?: string
}

export default function Permissions() {
  const navigate = useNavigate()
  const [roles, setRoles] = useState<Role[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    fetchRoles()
  }, [])

  const fetchRoles = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/permission/roles')
      setRoles(response.data.data || [])
    } catch (error: any) {
      setError(error.response?.data?.message || '获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="permissions-page">
        <div className="permissions-container">
          <Loading text="加载中..." />
        </div>
      </div>
    )
  }

  return (
    <div className="permissions-page">
      <div className="permissions-container">
        <header className="permissions-header">
          <h1 className="permissions-title">👑 我的角色</h1>
          <Button variant="outline" onClick={() => navigate('/dashboard')}>
            返回
          </Button>
        </header>

        <Card className="permissions-card">
          {error && <div className="error-message">{error}</div>}

          {roles.length === 0 ? (
            <div className="empty-state">
              <p>您还没有分配任何角色</p>
              <p className="hint">请联系管理员为您分配角色</p>
            </div>
          ) : (
            <div className="roles-list">
              {roles.map((role) => (
                <div key={role.id} className="role-item">
                  <div className="role-icon">👑</div>
                  <div className="role-info">
                    <div className="role-name">{role.display_name || role.name}</div>
                    {role.description && (
                      <div className="role-description">{role.description}</div>
                    )}
                    <div className="role-id">角色标识: {role.name}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        <Card className="permissions-card" style={{ marginTop: '24px' }}>
          <h2 className="card-subtitle">📋 权限说明</h2>
          <div className="permissions-info">
            <p>角色决定了您在系统中的权限范围。不同的角色拥有不同的操作权限。</p>
            <p className="hint">如果您需要更多权限，请联系系统管理员。</p>
          </div>
        </Card>
      </div>
    </div>
  )
}


