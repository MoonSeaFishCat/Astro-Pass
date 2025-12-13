import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../../components/Card'
import Button from '../../components/Button'
import Input from '../../components/Input'
import Loading from '../../components/Loading'
import './SAMLManagement.css'

interface SAMLConfig {
  id: number
  entity_id: string
  type: string
  name: string
  description: string
  status: string
  idp_sso_service_url: string
  idp_slo_service_url: string
  sign_assertions: boolean
  encrypt_assertions: boolean
  sign_requests: boolean
  created_at: string
  updated_at: string
}

export default function SAMLManagement() {
  const navigate = useNavigate()
  const [configs, setConfigs] = useState<SAMLConfig[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const [showCreateModal, setShowCreateModal] = useState(false)
  const [showMetadataModal, setShowMetadataModal] = useState(false)
  const [selectedMetadata, setSelectedMetadata] = useState('')
  const [selectedEntityId, setSelectedEntityId] = useState('')
  
  // 表单状态
  const [formData, setFormData] = useState({
    entity_id: '',
    type: 'idp',
    name: '',
    description: '',
    status: 'active',
    sign_assertions: true,
    encrypt_assertions: false,
    sign_requests: false,
  })

  // 获取SAML配置列表
  const fetchConfigs = async () => {
    setLoading(true)
    try {
      const response = await axios.get('/api/admin/saml/configs')
      if (response.data.code === 200) {
        setConfigs(response.data.data.configs || [])
      }
    } catch (error: any) {
      setError(error.response?.data?.message || '获取SAML配置失败')
    } finally {
      setLoading(false)
    }
  }

  // 创建SAML配置
  const handleCreate = () => {
    setFormData({
      entity_id: '',
      type: 'idp',
      name: '',
      description: '',
      status: 'active',
      sign_assertions: true,
      encrypt_assertions: false,
      sign_requests: false,
    })
    setShowCreateModal(true)
  }

  // 保存SAML配置
  const handleSave = async () => {
    try {
      await axios.post('/api/admin/saml/configs', formData)
      alert('SAML配置创建成功')
      setShowCreateModal(false)
      fetchConfigs()
    } catch (error: any) {
      alert(error.response?.data?.message || '保存SAML配置失败')
    }
  }

  // 删除SAML配置
  const handleDelete = (config: SAMLConfig) => {
    if (!confirm(`确定要删除SAML配置 "${config.name}" 吗？`)) {
      return
    }

    axios.delete(`/api/admin/saml/configs/${config.id}`)
      .then(() => {
        alert('SAML配置删除成功')
        fetchConfigs()
      })
      .catch((error: any) => {
        alert(error.response?.data?.message || '删除SAML配置失败')
      })
  }

  // 查看元数据
  const handleViewMetadata = async (entityId: string) => {
    try {
      const response = await axios.get(`/api/saml/metadata?entity_id=${entityId}`)
      setSelectedMetadata(response.data)
      setSelectedEntityId(entityId)
      setShowMetadataModal(true)
    } catch (error: any) {
      alert(error.response?.data?.message || '获取元数据失败')
    }
  }

  // 复制元数据URL
  const copyMetadataURL = (entityId: string) => {
    const url = `${window.location.origin}/api/saml/metadata?entity_id=${entityId}`
    navigator.clipboard.writeText(url).then(() => {
      alert('元数据URL已复制到剪贴板')
    })
  }

  const formatDate = (dateString: string) => {
    const date = new Date(dateString)
    return date.toLocaleString('zh-CN')
  }

  const getStatusColor = (status: string) => {
    return status === 'active' ? 'status-active' : 'status-inactive'
  }

  const getTypeLabel = (type: string) => {
    return type === 'idp' ? 'IdP (身份提供者)' : 'SP (服务提供者)'
  }

  useEffect(() => {
    fetchConfigs()
  }, [])

  if (loading) {
    return (
      <div className="saml-management-page">
        <div className="saml-management-container">
          <Loading text="加载中..." />
        </div>
      </div>
    )
  }

  return (
    <div className="saml-management-page">
      <div className="saml-management-container">
        <header className="saml-management-header">
          <h1 className="saml-management-title">🔒 SAML配置管理</h1>
          <div className="saml-management-actions">
            <Button variant="outline" onClick={() => navigate('/admin')}>
              返回
            </Button>
            <Button variant="secondary" onClick={fetchConfigs} disabled={loading}>
              刷新
            </Button>
            <Button variant="primary" onClick={handleCreate}>
              新建配置
            </Button>
          </div>
        </header>

        <Card className="saml-management-card">
          <div className="saml-info">
            <p>
              SAML (Security Assertion Markup Language) 是一种用于在不同安全域之间交换认证和授权数据的标准。
              这里可以配置作为身份提供者(IdP)或服务提供者(SP)的SAML设置。
            </p>
          </div>

          {error && <div className="error-message">{error}</div>}

          {configs.length === 0 ? (
            <div className="empty-state">
              <p>暂无SAML配置</p>
              <Button variant="primary" onClick={handleCreate}>
                创建第一个配置
              </Button>
            </div>
          ) : (
            <div className="configs-list">
              {configs.map((config) => (
                <div key={config.id} className="config-item">
                  <div className="config-info">
                    <div className="config-header-info">
                      <span className="config-icon">🔒</span>
                      <div className="config-main-info">
                        <div className="config-name">{config.name}</div>
                        <div className="config-entity-id">实体ID: {config.entity_id}</div>
                        <div className="config-description">{config.description}</div>
                      </div>
                    </div>
                    <div className="config-details">
                      <div className="detail-item">
                        <span className="detail-label">类型：</span>
                        <span className="detail-value type-badge">{getTypeLabel(config.type)}</span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">状态：</span>
                        <span className={`detail-value ${getStatusColor(config.status)}`}>
                          {config.status === 'active' ? '启用' : '禁用'}
                        </span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">创建时间：</span>
                        <span className="detail-value">{formatDate(config.created_at)}</span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">安全设置：</span>
                        <div className="security-badges">
                          {config.sign_assertions && <span className="security-badge">签名断言</span>}
                          {config.encrypt_assertions && <span className="security-badge">加密断言</span>}
                          {config.sign_requests && <span className="security-badge">签名请求</span>}
                        </div>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">SSO服务URL：</span>
                        <span className="detail-value url-text">{config.idp_sso_service_url}</span>
                      </div>
                      <div className="detail-item">
                        <span className="detail-label">SLO服务URL：</span>
                        <span className="detail-value url-text">{config.idp_slo_service_url}</span>
                      </div>
                    </div>
                  </div>
                  <div className="config-actions">
                    <Button
                      variant="outline"
                      size="small"
                      onClick={() => handleViewMetadata(config.entity_id)}
                    >
                      查看元数据
                    </Button>
                    <Button
                      variant="outline"
                      size="small"
                      onClick={() => copyMetadataURL(config.entity_id)}
                    >
                      复制URL
                    </Button>
                    <Button
                      variant="outline"
                      size="small"
                      onClick={() => handleDelete(config)}
                    >
                      删除
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>

        {/* 创建配置模态框 */}
        {showCreateModal && (
          <div className="modal-overlay" onClick={() => setShowCreateModal(false)}>
            <div className="modal-content large-modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>新建SAML配置</h3>
                <button className="modal-close" onClick={() => setShowCreateModal(false)}>
                  ×
                </button>
              </div>
              <div className="modal-body">
                <div className="form-section">
                  <h4>基本信息</h4>
                  <div className="form-grid">
                    <Input
                      label="实体ID"
                      value={formData.entity_id}
                      onChange={(e) => setFormData({ ...formData, entity_id: e.target.value })}
                      placeholder="例如: https://your-domain.com/saml/metadata"
                      required
                    />
                    <div className="input-group">
                      <label className="input-label">类型 <span className="required">*</span></label>
                      <select
                        className="input"
                        value={formData.type}
                        onChange={(e) => setFormData({ ...formData, type: e.target.value })}
                      >
                        <option value="idp">身份提供者 (IdP)</option>
                        <option value="sp">服务提供者 (SP)</option>
                      </select>
                    </div>
                  </div>
                  <Input
                    label="配置名称"
                    value={formData.name}
                    onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                    placeholder="例如: 企业SAML IdP"
                    required
                  />
                  <div className="input-group">
                    <label className="input-label">描述</label>
                    <textarea
                      className="input textarea"
                      value={formData.description}
                      onChange={(e) => setFormData({ ...formData, description: e.target.value })}
                      placeholder="配置描述信息"
                      rows={3}
                    />
                  </div>
                  <div className="input-group">
                    <label className="input-label">状态</label>
                    <select
                      className="input"
                      value={formData.status}
                      onChange={(e) => setFormData({ ...formData, status: e.target.value })}
                    >
                      <option value="active">启用</option>
                      <option value="inactive">禁用</option>
                    </select>
                  </div>
                </div>

                <div className="form-section">
                  <h4>安全设置</h4>
                  <div className="checkbox-group">
                    <label className="checkbox-item">
                      <input
                        type="checkbox"
                        checked={formData.sign_assertions}
                        onChange={(e) => setFormData({ ...formData, sign_assertions: e.target.checked })}
                      />
                      <span className="checkbox-label">签名断言</span>
                    </label>
                    <label className="checkbox-item">
                      <input
                        type="checkbox"
                        checked={formData.encrypt_assertions}
                        onChange={(e) => setFormData({ ...formData, encrypt_assertions: e.target.checked })}
                      />
                      <span className="checkbox-label">加密断言</span>
                    </label>
                    <label className="checkbox-item">
                      <input
                        type="checkbox"
                        checked={formData.sign_requests}
                        onChange={(e) => setFormData({ ...formData, sign_requests: e.target.checked })}
                      />
                      <span className="checkbox-label">签名请求</span>
                    </label>
                  </div>
                </div>
              </div>
              <div className="modal-footer">
                <Button variant="outline" onClick={() => setShowCreateModal(false)}>
                  取消
                </Button>
                <Button variant="primary" onClick={handleSave}>
                  创建
                </Button>
              </div>
            </div>
          </div>
        )}

        {/* 元数据查看模态框 */}
        {showMetadataModal && (
          <div className="modal-overlay" onClick={() => setShowMetadataModal(false)}>
            <div className="modal-content large-modal" onClick={(e) => e.stopPropagation()}>
              <div className="modal-header">
                <h3>SAML元数据 - {selectedEntityId}</h3>
                <button className="modal-close" onClick={() => setShowMetadataModal(false)}>
                  ×
                </button>
              </div>
              <div className="modal-body">
                <div className="metadata-info">
                  <p className="metadata-url">
                    <strong>元数据URL:</strong> {window.location.origin}/api/saml/metadata?entity_id={selectedEntityId}
                  </p>
                  <div className="metadata-actions">
                    <Button
                      variant="outline"
                      onClick={() => {
                        navigator.clipboard.writeText(selectedMetadata)
                        alert('元数据已复制到剪贴板')
                      }}
                    >
                      复制元数据
                    </Button>
                    <Button
                      variant="outline"
                      onClick={() => {
                        const url = `${window.location.origin}/api/saml/metadata?entity_id=${selectedEntityId}`
                        navigator.clipboard.writeText(url)
                        alert('元数据URL已复制到剪贴板')
                      }}
                    >
                      复制URL
                    </Button>
                  </div>
                </div>
                <div className="metadata-content">
                  <pre className="metadata-xml">{selectedMetadata}</pre>
                </div>
              </div>
              <div className="modal-footer">
                <Button onClick={() => setShowMetadataModal(false)}>关闭</Button>
              </div>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}