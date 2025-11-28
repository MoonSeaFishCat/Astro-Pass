import { useState, useEffect } from 'react';
import api from '../../utils/api';
import './BackupManagement.css';

interface Backup {
  id: number;
  file_name: string;
  file_size: number;
  backup_type: string;
  status: string;
  message: string;
  created_at: string;
}

interface BackupStats {
  total_count: number;
  success_count: number;
  failed_count: number;
  total_size: number;
  last_backup: string;
}

interface BackupConfig {
  auto_enabled: boolean;
  schedule: string;
  retention_days: number;
  max_backups: number;
}

function BackupManagement() {
  const [backups, setBackups] = useState<Backup[]>([]);
  const [stats, setStats] = useState<BackupStats | null>(null);
  const [config, setConfig] = useState<BackupConfig>({
    auto_enabled: true,
    schedule: '0 2 * * *',
    retention_days: 30,
    max_backups: 10,
  });
  const [loading, setLoading] = useState(false);
  const [page, setPage] = useState(1);
  const [total, setTotal] = useState(0);
  const [showConfigModal, setShowConfigModal] = useState(false);

  useEffect(() => {
    fetchBackups();
    fetchStats();
    fetchConfig();
  }, [page]);

  const fetchBackups = async () => {
    try {
      const response = await api.get(`/admin/backup?page=${page}&page_size=10`);
      setBackups(response.data.data.backups || []);
      setTotal(response.data.data.total || 0);
    } catch (error) {
      console.error('获取备份列表失败:', error);
    }
  };

  const fetchStats = async () => {
    try {
      const response = await api.get('/admin/backup/stats');
      setStats(response.data.data);
    } catch (error) {
      console.error('获取统计信息失败:', error);
    }
  };

  const fetchConfig = async () => {
    try {
      const response = await api.get('/admin/config/backup');
      setConfig(response.data.data);
    } catch (error) {
      console.error('获取备份配置失败:', error);
    }
  };

  const createBackup = async () => {
    if (loading) return;
    
    setLoading(true);
    try {
      await api.post('/admin/backup');
      alert('备份创建成功！');
      fetchBackups();
      fetchStats();
    } catch (error: any) {
      alert('备份创建失败: ' + (error.response?.data?.message || error.message));
    } finally {
      setLoading(false);
    }
  };

  const deleteBackup = async (id: number) => {
    if (!confirm('确定要删除这个备份吗？')) return;

    try {
      await api.delete(`/admin/backup/${id}`);
      alert('删除成功！');
      fetchBackups();
      fetchStats();
    } catch (error: any) {
      alert('删除失败: ' + (error.response?.data?.message || error.message));
    }
  };

  const restoreBackup = async (id: number) => {
    if (!confirm('确定要恢复这个备份吗？这将覆盖当前数据！')) return;

    setLoading(true);
    try {
      await api.post(`/admin/backup/${id}/restore`);
      alert('恢复成功！');
    } catch (error: any) {
      alert('恢复失败: ' + (error.response?.data?.message || error.message));
    } finally {
      setLoading(false);
    }
  };

  const downloadBackup = (id: number) => {
    window.open(`${api.defaults.baseURL}/admin/backup/${id}/download`, '_blank');
  };

  const cleanOldBackups = async () => {
    if (!confirm(`确定要清理超过${config.retention_days}天的旧备份吗？`)) return;

    try {
      await api.post(`/admin/backup/clean?days=${config.retention_days}`);
      alert('清理成功！');
      fetchBackups();
      fetchStats();
    } catch (error: any) {
      alert('清理失败: ' + (error.response?.data?.message || error.message));
    }
  };

  const saveConfig = async () => {
    try {
      await api.put('/admin/config/backup', config);
      alert('配置保存成功！');
      setShowConfigModal(false);
    } catch (error: any) {
      alert('配置保存失败: ' + (error.response?.data?.message || error.message));
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes === 0) return '0 B';
    const k = 1024;
    const sizes = ['B', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i];
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('zh-CN');
  };

  return (
    <div className="backup-management">
      <div className="page-header">
        <h1>备份管理</h1>
        <div className="header-actions">
          <button onClick={() => setShowConfigModal(true)} className="btn-secondary">
            备份设置
          </button>
          <button onClick={createBackup} disabled={loading} className="btn-primary">
            {loading ? '创建中...' : '创建备份'}
          </button>
        </div>
      </div>

      {/* 统计信息 */}
      {stats && (
        <div className="stats-grid">
          <div className="stat-card">
            <div className="stat-label">总备份数</div>
            <div className="stat-value">{stats.total_count}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">成功</div>
            <div className="stat-value success">{stats.success_count}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">失败</div>
            <div className="stat-value error">{stats.failed_count}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">总大小</div>
            <div className="stat-value">{formatFileSize(stats.total_size)}</div>
          </div>
          <div className="stat-card">
            <div className="stat-label">最后备份</div>
            <div className="stat-value small">{formatDate(stats.last_backup)}</div>
          </div>
        </div>
      )}

      {/* 备份列表 */}
      <div className="backup-list">
        <div className="list-header">
          <h2>备份列表</h2>
          <button onClick={cleanOldBackups} className="btn-warning">
            清理旧备份
          </button>
        </div>

        <table className="backup-table">
          <thead>
            <tr>
              <th>文件名</th>
              <th>大小</th>
              <th>类型</th>
              <th>状态</th>
              <th>创建时间</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {backups.map((backup) => (
              <tr key={backup.id}>
                <td>{backup.file_name}</td>
                <td>{formatFileSize(backup.file_size)}</td>
                <td>
                  <span className={`badge ${backup.backup_type}`}>
                    {backup.backup_type === 'auto' ? '自动' : '手动'}
                  </span>
                </td>
                <td>
                  <span className={`badge ${backup.status}`}>
                    {backup.status === 'success' ? '成功' : backup.status === 'failed' ? '失败' : '进行中'}
                  </span>
                </td>
                <td>{formatDate(backup.created_at)}</td>
                <td>
                  <div className="action-buttons">
                    <button
                      onClick={() => downloadBackup(backup.id)}
                      className="btn-icon"
                      title="下载"
                    >
                      📥
                    </button>
                    <button
                      onClick={() => restoreBackup(backup.id)}
                      className="btn-icon"
                      title="恢复"
                      disabled={loading || backup.status !== 'success'}
                    >
                      🔄
                    </button>
                    <button
                      onClick={() => deleteBackup(backup.id)}
                      className="btn-icon danger"
                      title="删除"
                    >
                      🗑️
                    </button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>

        {/* 分页 */}
        <div className="pagination">
          <button
            onClick={() => setPage(page - 1)}
            disabled={page === 1}
            className="btn-secondary"
          >
            上一页
          </button>
          <span>第 {page} 页 / 共 {Math.ceil(total / 10)} 页</span>
          <button
            onClick={() => setPage(page + 1)}
            disabled={page >= Math.ceil(total / 10)}
            className="btn-secondary"
          >
            下一页
          </button>
        </div>
      </div>

      {/* 配置模态框 */}
      {showConfigModal && (
        <div className="modal-overlay" onClick={() => setShowConfigModal(false)}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <h2>备份设置</h2>
            <div className="form-group">
              <label>
                <input
                  type="checkbox"
                  checked={config.auto_enabled}
                  onChange={(e) => setConfig({ ...config, auto_enabled: e.target.checked })}
                />
                启用自动备份
              </label>
            </div>
            <div className="form-group">
              <label>备份时间（Cron表达式）</label>
              <input
                type="text"
                value={config.schedule}
                onChange={(e) => setConfig({ ...config, schedule: e.target.value })}
                placeholder="0 2 * * *"
              />
              <small>默认：每天凌晨2点</small>
            </div>
            <div className="form-group">
              <label>保留天数</label>
              <input
                type="number"
                value={config.retention_days}
                onChange={(e) => setConfig({ ...config, retention_days: parseInt(e.target.value) })}
                min="1"
              />
            </div>
            <div className="form-group">
              <label>最大备份数量</label>
              <input
                type="number"
                value={config.max_backups}
                onChange={(e) => setConfig({ ...config, max_backups: parseInt(e.target.value) })}
                min="1"
              />
            </div>
            <div className="modal-actions">
              <button onClick={() => setShowConfigModal(false)} className="btn-secondary">
                取消
              </button>
              <button onClick={saveConfig} className="btn-primary">
                保存
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}

export default BackupManagement;
