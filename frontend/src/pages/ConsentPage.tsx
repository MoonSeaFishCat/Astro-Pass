import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import api from '../utils/api';
import './ConsentPage.css';

interface ScopeInfo {
  scope: string;
  description: string;
}

interface ClientInfo {
  client_name: string;
  client_uri: string;
  logo_uri: string;
  scopes: ScopeInfo[];
}

const ConsentPage: React.FC = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [clientInfo, setClientInfo] = useState<ClientInfo | null>(null);
  const [error, setError] = useState('');

  const clientId = searchParams.get('client_id');
  const scope = searchParams.get('scope');
  const redirectUri = searchParams.get('redirect_uri');
  const state = searchParams.get('state');
  const responseType = searchParams.get('response_type');

  useEffect(() => {
    if (!clientId || !scope) {
      setError('缺少必要参数');
      setLoading(false);
      return;
    }

    fetchClientInfo();
  }, [clientId, scope]);

  const fetchClientInfo = async () => {
    try {
      const response = await api.get('/oauth2/consent/info', {
        params: { client_id: clientId, scope: scope }
      });
      setClientInfo(response.data.data);
    } catch (err: any) {
      setError(err.response?.data?.message || '获取应用信息失败');
    } finally {
      setLoading(false);
    }
  };

  const handleApprove = async () => {
    try {
      setLoading(true);
      await api.post('/oauth2/consent/approve', {
        client_id: clientId,
        scope: scope
      });

      // 重定向回授权端点继续流程
      const params = new URLSearchParams({
        response_type: responseType || 'code',
        client_id: clientId || '',
        redirect_uri: redirectUri || '',
        scope: scope || '',
        state: state || '',
        consent: 'approved'
      });
      
      window.location.href = `/api/oauth2/authorize?${params.toString()}`;
    } catch (err: any) {
      setError(err.response?.data?.message || '授权失败');
      setLoading(false);
    }
  };

  const handleDeny = () => {
    // 重定向回应用并带上错误信息
    if (redirectUri) {
      const params = new URLSearchParams({
        error: 'access_denied',
        error_description: '用户拒绝授权',
        state: state || ''
      });
      window.location.href = `${redirectUri}?${params.toString()}`;
    } else {
      navigate('/dashboard');
    }
  };

  if (loading) {
    return (
      <div className="consent-page">
        <div className="consent-card">
          <div className="loading">加载中...</div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="consent-page">
        <div className="consent-card">
          <div className="error-message">{error}</div>
          <button onClick={() => navigate('/dashboard')} className="btn-secondary">
            返回
          </button>
        </div>
      </div>
    );
  }

  if (!clientInfo) {
    return null;
  }

  return (
    <div className="consent-page">
      <div className="consent-card">
        <div className="consent-header">
          {clientInfo.logo_uri && (
            <img src={clientInfo.logo_uri} alt={clientInfo.client_name} className="client-logo" />
          )}
          <h2 className="consent-title">授权请求</h2>
          <p className="consent-description">
            <strong>{clientInfo.client_name}</strong> 请求访问您的账户
          </p>
        </div>

        <div className="consent-body">
          <div className="permissions-section">
            <h3>该应用将能够：</h3>
            <ul className="permissions-list">
              {clientInfo.scopes.map((scopeInfo, index) => (
                <li key={index} className="permission-item">
                  <span className="permission-icon">✓</span>
                  <span className="permission-text">{scopeInfo.description}</span>
                </li>
              ))}
            </ul>
          </div>

          <div className="consent-info">
            <p className="info-text">
              授权后，您可以随时在账户设置中撤销此应用的访问权限。
            </p>
            {clientInfo.client_uri && (
              <p className="client-link">
                <a href={clientInfo.client_uri} target="_blank" rel="noopener noreferrer">
                  了解更多关于 {clientInfo.client_name}
                </a>
              </p>
            )}
          </div>
        </div>

        <div className="consent-actions">
          <button onClick={handleDeny} className="btn-deny" disabled={loading}>
            拒绝
          </button>
          <button onClick={handleApprove} className="btn-approve" disabled={loading}>
            {loading ? '授权中...' : '授权'}
          </button>
        </div>
      </div>
    </div>
  );
};

export default ConsentPage;
