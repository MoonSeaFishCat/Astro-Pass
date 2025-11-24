import { useState, useEffect } from 'react'
import axios from 'axios'
import Card from '../../components/Card'
import Button from '../../components/Button'
import Loading from '../../components/Loading'
import './PermissionManagement.css'

interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
  description?: string
}

export default function PermissionManagement() {
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [selectedPermission, setSelectedPermission] = useState<Permission | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    display_name: '',
    resource: '',
    action: '',
    description: '',
  })

  useEffect(() => {
    fetchPermissions()
  }, [])

  const fetchPermissions = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/admin/permissions')
      setPermissions(response.data.data || [])
    } catch (err: any) {
      setError(err.response?.data?.message || '获取权限列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async () => {
    try {
      await axios.post('/api/permission/permission', formData)
      setShowCreateModal(false)
      setFormData({
        name: '',
        display_name: '',
        resource: '',
        action: '',
        description: '',
      })
      fetchPermissions()
    } catch (err: any) {
      alert(err.response?.data?.message || '创建失败')
    }
  }

  const handleUpdate = async () => {
    if (!selectedPermission) return
    try {
      await axios.put(`/api/admin/permissions/${selectedPermission.id}`, {
        display_name: formData.display_name,
        description: formData.description,
      })
      setShowEditModal(false)
      setSelectedPermission(null)
      setFormData({
        name: '',
        display_name: '',
        resource: '',
        action: '',
        description: '',
      })
      fetchPermissions()
    } catch (err: any) {
      alert(err.response?.data?.message || '更新失败')
    }
  }

  const handleDelete = async (permissionId: number) => {
    if (!confirm('确定要删除这个权限吗？')) return
    try {
      await axios.delete(`/api/admin/permissions/${permissionId}`)
      fetchPermissions()
    } catch (err: any) {
      alert(err.response?.data?.message || '删除失败')
    }
  }

  const openEditModal = (permission: Permission) => {
    setSelectedPermission(permission)
    setFormData({
      name: permission.name,
      display_name: permission.display_name || '',
      resource: permission.resource,
      action: permission.action,
      description: permission.description || '',
    })
    setShowEditModal(true)
  }

  if (loading && permissions.length === 0) {
    return (
      <div className="permission-management">
        <Loading text="加载中..." />
      </div>
    )
  }

  return (
    <div className="permission-management">
      <div className="page-header">
        <h2 className="page-title">🔐 权限管理</h2>
        <Button variant="primary" onClick={() => setShowCreateModal(true)}>
          ➕ 创建权限
        </Button>
      </div>

      {error && <div className="error-message">{error}</div>}

      <Card className="permissions-table-card">
        <table className="permissions-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>资源</th>
              <th>操作</th>
              <th>描述</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {permissions.map((permission) => (
              <tr key={permission.id}>
                <td>{permission.id}</td>
                <td>{permission.display_name || permission.name}</td>
                <td>
                  <span className="resource-badge">{permission.resource}</span>
                </td>
                <td>
                  <span className="action-badge">{permission.action}</span>
                </td>
                <td>{permission.description || '-'}</td>
                <td>
                  <div className="action-buttons">
                    <Button
                      variant="secondary"
                      size="small"
                      onClick={() => openEditModal(permission)}
                    >
                      编辑
                    </Button>
                    <Button
                      variant="outline"
                      size="small"
                      onClick={() => handleDelete(permission.id)}
                    >
                      删除
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {permissions.length === 0 && !loading && (
          <div className="empty-state">暂无权限</div>
        )}
      </Card>

      {/* 创建权限模态框 */}
      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>创建权限</h3>
            <div className="form-group">
              <label>权限标识 *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="例如: user:read"
              />
            </div>
            <div className="form-group">
              <label>显示名称</label>
              <input
                type="text"
                value={formData.display_name}
                onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
                placeholder="例如: 查看用户"
              />
            </div>
            <div className="form-group">
              <label>资源 *</label>
              <input
                type="text"
                value={formData.resource}
                onChange={(e) => setFormData({ ...formData, resource: e.target.value })}
                placeholder="例如: user"
              />
            </div>
            <div className="form-group">
              <label>操作 *</label>
              <input
                type="text"
                value={formData.action}
                onChange={(e) => setFormData({ ...formData, action: e.target.value })}
                placeholder="例如: read"
              />
            </div>
            <div className="form-group">
              <label>描述</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="权限描述"
                rows={3}
              />
            </div>
            <div className="modal-actions">
              <Button variant="outline" onClick={() => setShowCreateModal(false)}>
                取消
              </Button>
              <Button variant="primary" onClick={handleCreate}>
                创建
              </Button>
            </div>
          </div>
        </div>
      )}

      {/* 编辑权限模态框 */}
      {showEditModal && selectedPermission && (
        <div className="modal-overlay" onClick={() => setShowEditModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>编辑权限</h3>
            <div className="form-group">
              <label>显示名称</label>
              <input
                type="text"
                value={formData.display_name}
                onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
              />
            </div>
            <div className="form-group">
              <label>描述</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                rows={3}
              />
            </div>
            <div className="modal-actions">
              <Button variant="outline" onClick={() => setShowEditModal(false)}>
                取消
              </Button>
              <Button variant="primary" onClick={handleUpdate}>
                保存
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}


