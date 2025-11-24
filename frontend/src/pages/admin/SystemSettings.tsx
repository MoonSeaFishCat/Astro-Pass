import { useState } from 'react'
import Card from '../../components/Card'
import Button from '../../components/Button'
import './SystemSettings.css'

export default function SystemSettings() {
  const [settings, setSettings] = useState({
    siteName: '星穹通行证',
    siteDescription: '统一身份认证通行证系统',
    allowRegistration: true,
    requireEmailVerification: false,
    sessionTimeout: 30,
    maxLoginAttempts: 5,
    lockoutDuration: 30,
  })

  const handleSave = () => {
    // TODO: 实现保存到后端
    alert('设置已保存（功能待实现）')
  }

  return (
    <div className="system-settings">
      <div className="page-header">
        <h2 className="page-title">⚙️ 系统设置</h2>
        <Button variant="primary" onClick={handleSave}>
          保存设置
        </Button>
      </div>

      <div className="settings-grid">
        <Card className="settings-card">
          <h3 className="card-title">基本信息</h3>
          <div className="form-group">
            <label>站点名称</label>
            <input
              type="text"
              value={settings.siteName}
              onChange={(e) => setSettings({ ...settings, siteName: e.target.value })}
            />
          </div>
          <div className="form-group">
            <label>站点描述</label>
            <textarea
              value={settings.siteDescription}
              onChange={(e) => setSettings({ ...settings, siteDescription: e.target.value })}
              rows={3}
            />
          </div>
        </Card>

        <Card className="settings-card">
          <h3 className="card-title">注册设置</h3>
          <div className="form-group">
            <label className="checkbox-label">
              <input
                type="checkbox"
                checked={settings.allowRegistration}
                onChange={(e) =>
                  setSettings({ ...settings, allowRegistration: e.target.checked })
                }
              />
              <span>允许用户注册</span>
            </label>
          </div>
          <div className="form-group">
            <label className="checkbox-label">
              <input
                type="checkbox"
                checked={settings.requireEmailVerification}
                onChange={(e) =>
                  setSettings({ ...settings, requireEmailVerification: e.target.checked })
                }
              />
              <span>要求邮箱验证</span>
            </label>
          </div>
        </Card>

        <Card className="settings-card">
          <h3 className="card-title">安全设置</h3>
          <div className="form-group">
            <label>会话超时（分钟）</label>
            <input
              type="number"
              value={settings.sessionTimeout}
              onChange={(e) =>
                setSettings({ ...settings, sessionTimeout: parseInt(e.target.value) })
              }
              min={5}
              max={1440}
            />
          </div>
          <div className="form-group">
            <label>最大登录尝试次数</label>
            <input
              type="number"
              value={settings.maxLoginAttempts}
              onChange={(e) =>
                setSettings({ ...settings, maxLoginAttempts: parseInt(e.target.value) })
              }
              min={3}
              max={10}
            />
          </div>
          <div className="form-group">
            <label>账户锁定时长（分钟）</label>
            <input
              type="number"
              value={settings.lockoutDuration}
              onChange={(e) =>
                setSettings({ ...settings, lockoutDuration: parseInt(e.target.value) })
              }
              min={5}
              max={1440}
            />
          </div>
        </Card>

        <Card className="settings-card">
          <h3 className="card-title">邮件设置</h3>
          <div className="form-group">
            <label>SMTP服务器</label>
            <input type="text" placeholder="smtp.example.com" disabled />
            <p className="form-hint">在配置文件中设置</p>
          </div>
          <div className="form-group">
            <label>发件人邮箱</label>
            <input type="email" placeholder="noreply@example.com" disabled />
            <p className="form-hint">在配置文件中设置</p>
          </div>
        </Card>
      </div>
    </div>
  )
}


