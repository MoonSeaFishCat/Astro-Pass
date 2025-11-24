import { useState, useEffect } from 'react'
import axios from 'axios'
import Card from '../../components/Card'
import Button from '../../components/Button'
import Loading from '../../components/Loading'
import './RoleManagement.css'

interface Role {
  id: number
  name: string
  display_name: string
  description?: string
  permissions?: Permission[]
}

interface Permission {
  id: number
  name: string
  display_name: string
  resource: string
  action: string
}

export default function RoleManagement() {
  const [roles, setRoles] = useState<Role[]>([])
  const [permissions, setPermissions] = useState<Permission[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showEditModal, setShowEditModal] = useState(false)
  const [showAssignModal, setShowAssignModal] = useState(false)
  const [selectedRole, setSelectedRole] = useState<Role | null>(null)
  const [formData, setFormData] = useState({
    name: '',
    display_name: '',
    description: '',
  })

  useEffect(() => {
    fetchRoles()
    fetchPermissions()
  }, [])

  const fetchRoles = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/admin/roles')
      setRoles(response.data.data || [])
    } catch (err: any) {
      setError(err.response?.data?.message || '获取角色列表失败')
    } finally {
      setLoading(false)
    }
  }

  const fetchPermissions = async () => {
    try {
      const response = await axios.get('/api/admin/permissions')
      setPermissions(response.data.data || [])
    } catch (err: any) {
      console.error('获取权限列表失败:', err)
    }
  }

  const handleCreate = async () => {
    try {
      await axios.post('/api/permission/role', formData)
      setShowCreateModal(false)
      setFormData({ name: '', display_name: '', description: '' })
      fetchRoles()
    } catch (err: any) {
      alert(err.response?.data?.message || '创建失败')
    }
  }

  const handleUpdate = async () => {
    if (!selectedRole) return
    try {
      await axios.put(`/api/admin/roles/${selectedRole.id}`, {
        display_name: formData.display_name,
        description: formData.description,
      })
      setShowEditModal(false)
      setSelectedRole(null)
      setFormData({ name: '', display_name: '', description: '' })
      fetchRoles()
    } catch (err: any) {
      alert(err.response?.data?.message || '更新失败')
    }
  }

  const handleDelete = async (roleId: number) => {
    if (!confirm('确定要删除这个角色吗？')) return
    try {
      await axios.delete(`/api/admin/roles/${roleId}`)
      fetchRoles()
    } catch (err: any) {
      alert(err.response?.data?.message || '删除失败')
    }
  }

  const handleAssignPermission = async (roleName: string, resource: string, action: string) => {
    try {
      await axios.post(`/api/permission/role/${roleName}/permission`, {
        resource,
        action,
      })
      fetchRoles()
    } catch (err: any) {
      alert(err.response?.data?.message || '分配权限失败')
    }
  }

  const openEditModal = (role: Role) => {
    setSelectedRole(role)
    setFormData({
      name: role.name,
      display_name: role.display_name || '',
      description: role.description || '',
    })
    setShowEditModal(true)
  }

  const openAssignModal = (role: Role) => {
    setSelectedRole(role)
    setShowAssignModal(true)
  }

  if (loading && roles.length === 0) {
    return (
      <div className="role-management">
        <Loading text="加载中..." />
      </div>
    )
  }

  return (
    <div className="role-management">
      <div className="page-header">
        <h2 className="page-title">👑 角色管理</h2>
        <Button variant="primary" onClick={() => setShowCreateModal(true)}>
          ➕ 创建角色
        </Button>
      </div>

      {error && <div className="error-message">{error}</div>}

      <div className="roles-grid">
        {roles.map((role) => (
          <Card key={role.id} className="role-card">
            <div className="role-header">
              <h3 className="role-name">{role.display_name || role.name}</h3>
              <div className="role-actions">
                <Button
                  variant="secondary"
                  size="small"
                  onClick={() => openEditModal(role)}
                >
                  编辑
                </Button>
                <Button
                  variant="outline"
                  size="small"
                  onClick={() => openAssignModal(role)}
                >
                  分配权限
                </Button>
                <Button
                  variant="outline"
                  size="small"
                  onClick={() => handleDelete(role.id)}
                >
                  删除
                </Button>
              </div>
            </div>
            <div className="role-info">
              <div className="info-row">
                <span className="info-label">标识：</span>
                <span className="info-value">{role.name}</span>
              </div>
              {role.description && (
                <div className="info-row">
                  <span className="info-label">描述：</span>
                  <span className="info-value">{role.description}</span>
                </div>
              )}
            </div>
          </Card>
        ))}
      </div>

      {/* 创建角色模态框 */}
      {showCreateModal && (
        <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>创建角色</h3>
            <div className="form-group">
              <label>角色标识 *</label>
              <input
                type="text"
                value={formData.name}
                onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                placeholder="例如: editor"
              />
            </div>
            <div className="form-group">
              <label>显示名称</label>
              <input
                type="text"
                value={formData.display_name}
                onChange={(e) => setFormData({ ...formData, display_name: e.target.value })}
                placeholder="例如: 编辑者"
              />
            </div>
            <div className="form-group">
              <label>描述</label>
              <textarea
                value={formData.description}
                onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                placeholder="角色描述"
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

      {/* 编辑角色模态框 */}
      {showEditModal && selectedRole && (
        <div className="modal-overlay" onClick={() => setShowEditModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h3>编辑角色</h3>
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

      {/* 分配权限模态框 */}
      {showAssignModal && selectedRole && (
        <div className="modal-overlay" onClick={() => setShowAssignModal(false)}>
          <div className="modal-content large" onClick={(e) => e.stopPropagation()}>
            <h3>为角色 "{selectedRole.display_name || selectedRole.name}" 分配权限</h3>
            <div className="permissions-list">
              {permissions.map((permission) => (
                <div key={permission.id} className="permission-item">
                  <span className="permission-name">
                    {permission.display_name || permission.name}
                  </span>
                  <span className="permission-resource">
                    {permission.resource}:{permission.action}
                  </span>
                  <Button
                    variant="secondary"
                    size="small"
                    onClick={() =>
                      handleAssignPermission(
                        selectedRole.name,
                        permission.resource,
                        permission.action
                      )
                    }
                  >
                    分配
                  </Button>
                </div>
              ))}
            </div>
            <div className="modal-actions">
              <Button variant="outline" onClick={() => setShowAssignModal(false)}>
                关闭
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  )
}


