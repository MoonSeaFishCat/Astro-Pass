import { useState, useEffect } from 'react'
import { useNavigate } from 'react-router-dom'
import { QRCodeSVG } from 'qrcode.react'
import axios from 'axios'
import Card from '../components/Card'
import Button from '../components/Button'
import Input from '../components/Input'
import './MFASettings.css'

const API_BASE_URL = '/api'

export default function MFASettings() {
  const navigate = useNavigate()
  const [step, setStep] = useState<'generate' | 'verify' | 'enabled'>('generate')
  const [qrCodeURL, setQrCodeURL] = useState('')
  const [secret, setSecret] = useState('')
  const [code, setCode] = useState('')
  const [recoveryCodes, setRecoveryCodes] = useState<string[]>([])
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    generateTOTP()
  }, [])

  const generateTOTP = async () => {
    try {
      const response = await axios.post(`${API_BASE_URL}/mfa/generate`)
      const { data } = response.data
      setQrCodeURL(data.qr_code_url)
      setSecret(data.secret)
      setStep('verify')
    } catch (error: any) {
      setError(error.response?.data?.message || '生成TOTP密钥失败')
    }
  }

  const handleEnable = async (e: React.FormEvent) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const response = await axios.post(`${API_BASE_URL}/mfa/enable`, { code })
      const { data } = response.data
      setRecoveryCodes(data.recovery_codes)
      setStep('enabled')
    } catch (error: any) {
      setError(error.response?.data?.message || '启用MFA失败，请检查验证码')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="mfa-page">
      <div className="mfa-container">
        <header className="mfa-header">
          <h1 className="mfa-title">🔐 安全守护契约</h1>
          <Button variant="outline" onClick={() => navigate('/dashboard')}>
            返回
          </Button>
        </header>

        <Card className="mfa-card">
          {step === 'verify' && (
            <div className="mfa-content">
              <div className="mfa-step-header">
                <h2>步骤 1：扫描二维码</h2>
                <p className="mfa-description">
                  使用您的身份验证应用（如 Google Authenticator、Microsoft Authenticator）扫描下方二维码
                </p>
              </div>

              {qrCodeURL && (
                <div className="qr-code-container">
                  <QRCodeSVG value={qrCodeURL} size={200} />
                </div>
              )}

              <div className="secret-container">
                <p className="secret-label">或者手动输入密钥：</p>
                <code className="secret-code">{secret}</code>
              </div>

              <div className="mfa-step-header" style={{ marginTop: '32px' }}>
                <h2>步骤 2：验证代码</h2>
                <p className="mfa-description">
                  在您的身份验证应用中输入6位验证码以完成设置
                </p>
              </div>

              <form onSubmit={handleEnable} className="mfa-form">
                {error && <div className="error-message">{error}</div>}

                <Input
                  label="验证码"
                  type="text"
                  placeholder="请输入6位验证码"
                  value={code}
                  onChange={(e) => setCode(e.target.value.replace(/\D/g, '').slice(0, 6))}
                  required
                />

                <Button type="submit" fullWidth disabled={loading || code.length !== 6}>
                  {loading ? '启用中...' : '启用MFA'}
                </Button>
              </form>
            </div>
          )}

          {step === 'enabled' && (
            <div className="mfa-content">
              <div className="success-message">
                <h2>✨ MFA启用成功！</h2>
                <p>您的账户现在受到多因素认证保护</p>
              </div>

              <div className="recovery-codes-container">
                <h3>恢复码（请妥善保管）</h3>
                <p className="recovery-codes-hint">
                  如果丢失了身份验证设备，可以使用这些恢复码登录。每个恢复码只能使用一次。
                </p>
                <div className="recovery-codes-list">
                  {recoveryCodes.map((code, index) => (
                    <code key={index} className="recovery-code">
                      {code}
                    </code>
                  ))}
                </div>
                <Button
                  variant="secondary"
                  fullWidth
                  onClick={() => {
                    // 复制恢复码到剪贴板
                    navigator.clipboard.writeText(recoveryCodes.join('\n'))
                    alert('恢复码已复制到剪贴板')
                  }}
                >
                  复制所有恢复码
                </Button>
              </div>

              <Button
                variant="primary"
                fullWidth
                onClick={() => navigate('/dashboard')}
                style={{ marginTop: '24px' }}
              >
                完成
              </Button>
            </div>
          )}
        </Card>
      </div>
    </div>
  )
}


