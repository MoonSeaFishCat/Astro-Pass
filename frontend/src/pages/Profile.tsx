import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '../stores/authStore'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import './Profile.css'

export default function Profile() {
  const navigate = useNavigate()
  const { user } = useAuthStore()

  const [nickname, setNickname] = useState(user?.nickname || '')
  const [loading, setLoading] = useState(false)
  const [message, setMessage] = useState('')

  if (!user) {
    navigate('/login')
    return null
  }

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setMessage('')

    try {
      const response = await axios.put('/api/user/profile', { nickname })
      const { data } = response.data
      
      // 更新store中的用户信息
      useAuthStore.setState({ user: data.user })
      
      setMessage('资料更新成功！')
      setTimeout(() => {
        navigate('/dashboard')
      }, 1500)
    } catch (error: any) {
      setMessage(error.response?.data?.message || '更新失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="profile-page">
      <div className="profile-container">
        <header className="profile-header">
          <h1 className="profile-title">📋 个人资料</h1>
          <Button variant="outline" onClick={() => navigate('/dashboard')}>
            返回
          </Button>
        </header>

        <Card className="profile-card">
          <form onSubmit={handleSubmit} className="profile-form">
            {message && (
              <div className={`message ${message.includes('成功') ? 'message-success' : 'message-error'}`}>
                {message}
              </div>
            )}

            <div className="form-group">
              <label className="form-label">用户名</label>
              <input
                type="text"
                value={user.username}
                className="form-input"
                disabled
              />
              <p className="form-hint">用户名不可修改</p>
            </div>

            <div className="form-group">
              <label className="form-label">邮箱</label>
              <input
                type="email"
                value={user.email}
                className="form-input"
                disabled
              />
              <p className="form-hint">邮箱不可修改</p>
            </div>

            <div className="form-group">
              <label className="form-label">昵称</label>
              <input
                type="text"
                value={nickname}
                onChange={(e) => setNickname(e.target.value)}
                className="form-input"
                placeholder="请输入昵称"
              />
            </div>

            <div className="form-actions">
              <Button type="submit" fullWidth disabled={loading}>
                {loading ? '保存中...' : '保存更改'}
              </Button>
            </div>
          </form>
        </Card>
      </div>
    </div>
  )
}

