import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import Loading from '../components/Loading'
import './OAuth2Clients.css'

interface OAuth2Client {
  id: number
  client_id: string
  client_name: string
  client_uri?: string
  logo_uri?: string
  status: string
  created_at: string
}

export default function OAuth2Clients() {
  const navigate = useNavigate()
  const [clients, setClients] = useState<OAuth2Client[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')
  const [showCreateForm, setShowCreateForm] = useState(false)
  const [creating, setCreating] = useState(false)
  const [formData, setFormData] = useState({
    client_name: '',
    client_uri: '',
    logo_uri: '',
    redirect_uris: '',
  })
  const [revoking, setRevoking] = useState<string | null>(null)

  useEffect(() => {
    fetchClients()
  }, [])

  const fetchClients = async () => {
    try {
      setLoading(true)
      const response = await axios.get('/api/oauth2/clients')
      setClients(response.data.data || [])
    } catch (error: any) {
      setError(error.response?.data?.message || '获取客户端列表失败')
    } finally {
      setLoading(false)
    }
  }

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault()
    setCreating(true)

    try {
      const redirectURIs = formData.redirect_uris
        .split(',')
        .map(uri => uri.trim())
        .filter(uri => uri.length > 0)

      const response = await axios.post('/api/oauth2/clients', {
        client_name: formData.client_name,
        client_uri: formData.client_uri || undefined,
        logo_uri: formData.logo_uri || undefined,
        redirect_uris: redirectURIs,
      })

      const data = response.data.data
      alert(`客户端创建成功！\nClient ID: ${data.client_id}\nClient Secret: ${data.client_secret}\n\n请妥善保管Client Secret，它只会显示一次！`)
      
      setShowCreateForm(false)
      setFormData({ client_name: '', client_uri: '', logo_uri: '', redirect_uris: '' })
      await fetchClients()
    } catch (error: any) {
      alert(error.response?.data?.message || '创建失败')
    } finally {
      setCreating(false)
    }
  }

  const handleRevokeClient = async (clientId: string) => {
    try {
      setRevoking(clientId)
      await axios.delete(`/api/oauth2/clients/${clientId}`)
      await fetchClients()
      alert('客户端已撤销')
    } catch (error: any) {
      alert(error.response?.data?.message || '撤销失败')
    } finally {
      setRevoking(null)
    }
  }

  if (loading) {
    return (
      <div className="oauth2-clients-page">
        <div className="oauth2-clients-container">
          <Loading text="加载中..." />
        </div>
      </div>
    )
  }

  return (
    <div className="oauth2-clients-page">
      <div className="oauth2-clients-container">
        <header className="oauth2-clients-header">
          <h1 className="oauth2-clients-title">🔑 OAuth2 客户端</h1>
          <div className="oauth2-clients-actions">
            <Button variant="outline" onClick={() => navigate('/dashboard')}>
              返回
            </Button>
            <Button variant="primary" onClick={() => setShowCreateForm(!showCreateForm)}>
              {showCreateForm ? '取消' : '创建客户端'}
            </Button>
          </div>
        </header>

        {showCreateForm && (
          <Card className="oauth2-clients-card" style={{ marginBottom: '24px' }}>
            <h2 className="card-subtitle">创建新的OAuth2客户端</h2>
            <form onSubmit={handleCreate} className="create-form">
              <Input
                label="客户端名称"
                type="text"
                placeholder="请输入客户端名称"
                value={formData.client_name}
                onChange={(e) => setFormData({ ...formData, client_name: e.target.value })}
                required
              />
              <Input
                label="客户端URI"
                type="url"
                placeholder="https://example.com"
                value={formData.client_uri}
                onChange={(e) => setFormData({ ...formData, client_uri: e.target.value })}
              />
              <Input
                label="Logo URI"
                type="url"
                placeholder="https://example.com/logo.png"
                value={formData.logo_uri}
                onChange={(e) => setFormData({ ...formData, logo_uri: e.target.value })}
              />
              <Input
                label="重定向URI（多个用逗号分隔）"
                type="text"
                placeholder="https://example.com/callback,https://example.com/callback2"
                value={formData.redirect_uris}
                onChange={(e) => setFormData({ ...formData, redirect_uris: e.target.value })}
                required
              />
              <Button type="submit" fullWidth disabled={creating}>
                {creating ? '创建中...' : '创建客户端'}
              </Button>
            </form>
          </Card>
        )}

        <Card className="oauth2-clients-card">
          {error && <div className="error-message">{error}</div>}

          {clients.length === 0 ? (
            <div className="empty-state">
              <p>您还没有创建任何OAuth2客户端</p>
              <p className="hint">创建OAuth2客户端以允许其他应用使用您的账户进行授权登录</p>
              <Button
                variant="primary"
                onClick={() => setShowCreateForm(true)}
                style={{ marginTop: '16px' }}
              >
                创建第一个客户端
              </Button>
            </div>
          ) : (
            <div className="clients-list">
              {clients.map((client) => (
                <div key={client.id} className="client-item">
                  <div className="client-info">
                    <div className="client-name">{client.client_name}</div>
                    <div className="client-id">Client ID: {client.client_id}</div>
                    {client.client_uri && (
                      <div className="client-uri">
                        <a href={client.client_uri} target="_blank" rel="noopener noreferrer">
                          {client.client_uri}
                        </a>
                      </div>
                    )}
                    <div className="client-status">
                      <span className={`status-badge status-${client.status}`}>
                        {client.status}
                      </span>
                    </div>
                  </div>
                  <div className="client-actions">
                    <Button
                      variant="outline"
                      onClick={() => {
                        if (confirm('确定要撤销这个客户端吗？')) {
                          handleRevokeClient(client.client_id)
                        }
                      }}
                      disabled={revoking === client.client_id || client.status === 'revoked'}
                    >
                      {revoking === client.client_id ? '撤销中...' : client.status === 'revoked' ? '已撤销' : '撤销'}
                    </Button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </Card>
      </div>
    </div>
  )
}

