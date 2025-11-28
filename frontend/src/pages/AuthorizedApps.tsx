import React, { useEffect, useState } from 'react';
import api from '../utils/api';
import './AuthorizedApps.css';

interface Consent {
  id: number;
  client_id: string;
  scope: string;
  created_at: string;
  expires_at: string;
}

const AuthorizedApps: React.FC = () => {
  const [consents, setConsents] = useState<Consent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    fetchConsents();
  }, []);

  const fetchConsents = async () => {
    try {
      const response = await api.get('/oauth2/consent/list');
      setConsents(response.data.data || []);
    } catch (err: any) {
      setError(err.response?.data?.message || '获取授权列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleRevoke = async (clientId: string) => {
    if (!confirm('确定要撤销此应用的授权吗？')) {
      return;
    }

    try {
      await api.delete(`/oauth2/consent/${clientId}`);
      setConsents(consents.filter(c => c.client_id !== clientId));
    } catch (err: any) {
      alert(err.response?.data?.message || '撤销授权失败');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('zh-CN', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  };

  const parseScopeDescriptions = (scopeString: string) => {
    const scopeMap: { [key: string]: string } = {
      'openid': '基本身份信息',
      'profile': '个人资料',
      'email': '邮箱地址',
      'phone': '手机号码',
      'address': '地址信息'
    };

    return scopeString.split(' ').map(s => scopeMap[s] || s).join('、');
  };

  if (loading) {
    return (
      <div className="authorized-apps-page">
        <div className="loading">加载中...</div>
      </div>
    );
  }

  return (
    <div className="authorized-apps-page">
      <div className="page-header">
        <h1>已授权的应用</h1>
        <p className="page-description">
          管理已授权访问您账户的第三方应用
        </p>
      </div>

      {error && (
        <div className="error-banner">
          {error}
        </div>
      )}

      {consents.length === 0 ? (
        <div className="empty-state">
          <div className="empty-icon">🔐</div>
          <h3>暂无授权应用</h3>
          <p>您还没有授权任何第三方应用访问您的账户</p>
        </div>
      ) : (
        <div className="consents-list">
          {consents.map((consent) => (
            <div key={consent.id} className="consent-card">
              <div className="consent-info">
                <div className="consent-header">
                  <h3 className="client-name">{consent.client_id}</h3>
                  <span className="consent-date">
                    授权于 {formatDate(consent.created_at)}
                  </span>
                </div>
                <div className="consent-details">
                  <div className="detail-item">
                    <span className="detail-label">权限范围：</span>
                    <span className="detail-value">
                      {parseScopeDescriptions(consent.scope)}
                    </span>
                  </div>
                  <div className="detail-item">
                    <span className="detail-label">有效期至：</span>
                    <span className="detail-value">
                      {formatDate(consent.expires_at)}
                    </span>
                  </div>
                </div>
              </div>
              <div className="consent-actions">
                <button
                  onClick={() => handleRevoke(consent.client_id)}
                  className="btn-revoke"
                >
                  撤销授权
                </button>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

export default AuthorizedApps;
