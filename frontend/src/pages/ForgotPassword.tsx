import { useState } from 'react'
import { Link } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import './ForgotPassword.css'

export default function ForgotPassword() {
  const [email, setEmail] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)
  const [success, setSuccess] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setSuccess(false)
    setLoading(true)

    try {
      await axios.post('/api/auth/forgot-password', { email })
      setSuccess(true)
    } catch (error: any) {
      setError(error.response?.data?.message || '发送失败，请稍后再试')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="forgot-password-page">
      <div className="forgot-password-container">
        <div className="forgot-password-header">
          <h1 className="forgot-password-title">🔑 找回密码</h1>
          <p className="forgot-password-subtitle">
            请输入您的邮箱地址，我们将发送密码重置链接
          </p>
        </div>

        <Card className="forgot-password-card">
          {success ? (
            <div className="success-content">
              <div className="success-icon">✨</div>
              <h2>邮件已发送</h2>
              <p>
                如果该邮箱存在，重置链接已发送到 <strong>{email}</strong>
              </p>
              <p className="hint">
                请检查您的邮箱（包括垃圾邮件文件夹），点击重置链接来设置新密码。
              </p>
              <Link to="/login">
                <Button variant="primary" fullWidth style={{ marginTop: '24px' }}>
                  返回登录
                </Button>
              </Link>
            </div>
          ) : (
            <form onSubmit={handleSubmit} className="forgot-password-form">
              {error && <div className="error-message">{error}</div>}

              <Input
                label="邮箱地址"
                type="email"
                placeholder="请输入注册时使用的邮箱"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
              />

              <Button type="submit" fullWidth disabled={loading}>
                {loading ? '发送中...' : '发送重置链接'}
              </Button>

              <div className="forgot-password-footer">
                <Link to="/login" className="link">
                  返回登录 →
                </Link>
              </div>
            </form>
          )}
        </Card>
      </div>
    </div>
  )
}


