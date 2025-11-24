import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import './ChangePassword.css'

export default function ChangePassword() {
  const navigate = useNavigate()
  const [oldPassword, setOldPassword] = useState('')
  const [newPassword, setNewPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)

    if (newPassword !== confirmPassword) {
      setError('两次输入的密码不一致哦~')
      return
    }

    if (newPassword.length < 6) {
      setError('密码长度至少为6位')
      return
    }

    setLoading(true)

    try {
      await axios.post('/api/user/change-password', {
        old_password: oldPassword,
        new_password: newPassword,
      })

      setSuccess(true)
      setTimeout(() => {
        navigate('/dashboard')
      }, 2000)
    } catch (error: any) {
      setError(error.response?.data?.message || '密码修改失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="change-password-page">
      <div className="change-password-container">
        <header className="change-password-header">
          <h1 className="change-password-title">🔐 修改密码</h1>
          <Button variant="outline" onClick={() => navigate('/dashboard')}>
            返回
          </Button>
        </header>

        <Card className="change-password-card">
          {success && (
            <div className="success-message">
              ✨ 密码修改成功！正在返回...
            </div>
          )}

          <form onSubmit={handleSubmit} className="change-password-form">
            {error && <div className="error-message">{error}</div>}

            <Input
              label="原密码"
              type="password"
              placeholder="请输入原密码"
              value={oldPassword}
              onChange={(e) => setOldPassword(e.target.value)}
              required
            />

            <Input
              label="新密码"
              type="password"
              placeholder="请输入新密码（至少6位）"
              value={newPassword}
              onChange={(e) => setNewPassword(e.target.value)}
              required
            />

            <Input
              label="确认新密码"
              type="password"
              placeholder="请再次输入新密码"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              required
            />

            <Button type="submit" fullWidth disabled={loading || success}>
              {loading ? '修改中...' : '修改密码'}
            </Button>
          </form>
        </Card>
      </div>
    </div>
  )
}


