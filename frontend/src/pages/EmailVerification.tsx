import { useState, useEffect } from 'react'
import { useNavigate, useSearchParams } from 'react-router-dom'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Loading from '../components/Loading'
import { useAuthStore } from '../stores/authStore'
import './EmailVerification.css'

export default function EmailVerification() {
  const navigate = useNavigate()
  const [searchParams] = useSearchParams()
  const { user } = useAuthStore()
  const [loading, setLoading] = useState(false)
  const [status, setStatus] = useState<'idle' | 'success' | 'error'>('idle')
  const [message, setMessage] = useState('')
  const [sending, setSending] = useState(false)

  const token = searchParams.get('token')

  useEffect(() => {
    if (token) {
      verifyEmail(token)
    }
  }, [token])

  const verifyEmail = async (verifyToken: string) => {
    try {
      setLoading(true)
      await axios.post('/api/email-verification/verify', { token: verifyToken })
      setStatus('success')
      setMessage('邮箱验证成功！')
    } catch (err: any) {
      setStatus('error')
      setMessage(err.response?.data?.message || '验证失败')
    } finally {
      setLoading(false)
    }
  }

  const sendVerificationEmail = async () => {
    if (!user?.email) {
      setMessage('请先设置邮箱地址')
      return
    }

    try {
      setSending(true)
      await axios.post('/api/email-verification/send', { email: user.email })
      setMessage('验证邮件已发送，请查收您的邮箱')
      setStatus('success')
    } catch (err: any) {
      setMessage(err.response?.data?.message || '发送失败')
      setStatus('error')
    } finally {
      setSending(false)
    }
  }

  if (loading) {
    return (
      <div className="email-verification-page">
        <Loading text="验证中..." />
      </div>
    )
  }

  return (
    <div className="email-verification-page">
      <Card className="verification-card">
        <div className="verification-icon">
          {status === 'success' ? '✅' : status === 'error' ? '❌' : '📧'}
        </div>
        <h2 className="verification-title">
          {status === 'success'
            ? '验证成功'
            : status === 'error'
            ? '验证失败'
            : '邮箱验证'}
        </h2>

        {token ? (
          <div className="verification-content">
            {status === 'success' && (
              <>
                <p className="verification-message">{message}</p>
                <Button variant="primary" onClick={() => navigate('/dashboard')}>
                  返回首页
                </Button>
              </>
            )}
            {status === 'error' && (
              <>
                <p className="verification-message error">{message}</p>
                <Button variant="primary" onClick={() => navigate('/dashboard')}>
                  返回首页
                </Button>
              </>
            )}
          </div>
        ) : (
          <div className="verification-content">
            <p className="verification-message">
              {user?.email_verified
                ? '您的邮箱已验证'
                : '请验证您的邮箱地址以激活账户'}
            </p>
            {user?.email && (
              <div className="email-info">
                <p>当前邮箱：{user.email}</p>
              </div>
            )}
            {!user?.email_verified && (
              <Button
                variant="primary"
                onClick={sendVerificationEmail}
                disabled={sending}
                fullWidth
              >
                {sending ? '发送中...' : '发送验证邮件'}
              </Button>
            )}
            {message && (
              <p className={`verification-message ${status === 'error' ? 'error' : ''}`}>
                {message}
              </p>
            )}
            <Button variant="outline" onClick={() => navigate('/dashboard')} fullWidth>
              返回首页
            </Button>
          </div>
        )}
      </Card>
    </div>
  )
}


