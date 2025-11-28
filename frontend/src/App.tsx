import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom'
import { useAuthStore } from './stores/authStore'
import ErrorBoundary from './components/ErrorBoundary'
import Login from './pages/Login'
import Register from './pages/Register'
import UserDashboard from './pages/UserDashboard'
import Profile from './pages/Profile'
import MFASettings from './pages/MFASettings'
import ChangePassword from './pages/ChangePassword'
import ForgotPassword from './pages/ForgotPassword'
import ResetPassword from './pages/ResetPassword'
import Sessions from './pages/Sessions'
import Permissions from './pages/Permissions'
import AuditLogs from './pages/AuditLogs'
import OAuth2Clients from './pages/OAuth2Clients'
import Notifications from './pages/Notifications'
import EmailVerification from './pages/EmailVerification'
import ConsentPage from './pages/ConsentPage'
import AuthorizedApps from './pages/AuthorizedApps'
import AdminLayout from './layouts/AdminLayout'
import AdminDashboard from './pages/admin/AdminDashboard'
import UserManagement from './pages/admin/UserManagement'
import RoleManagement from './pages/admin/RoleManagement'
import PermissionManagement from './pages/admin/PermissionManagement'
import SystemSettings from './pages/admin/SystemSettings'
import BackupManagement from './pages/admin/BackupManagement'
import './App.css'

function PrivateRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated } = useAuthStore()
  return isAuthenticated ? <>{children}</> : <Navigate to="/login" />
}

function AdminRoute({ children }: { children: React.ReactNode }) {
  const { isAuthenticated, user } = useAuthStore()
  if (!isAuthenticated) {
    return <Navigate to="/login" />
  }
  const isAdmin = user?.roles?.some((role: any) => role.name === 'admin') || false
  return isAdmin ? <>{children}</> : <Navigate to="/dashboard" />
}

function App() {
  return (
    <ErrorBoundary>
      <BrowserRouter>
        <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route
          path="/dashboard"
          element={
            <PrivateRoute>
              <UserDashboard />
            </PrivateRoute>
          }
        />
        <Route
          path="/profile"
          element={
            <PrivateRoute>
              <Profile />
            </PrivateRoute>
          }
        />
        <Route
          path="/mfa"
          element={
            <PrivateRoute>
              <MFASettings />
            </PrivateRoute>
          }
        />
        <Route
          path="/change-password"
          element={
            <PrivateRoute>
              <ChangePassword />
            </PrivateRoute>
          }
        />
        <Route path="/forgot-password" element={<ForgotPassword />} />
        <Route path="/reset-password" element={<ResetPassword />} />
        <Route
          path="/notifications"
          element={
            <PrivateRoute>
              <Notifications />
            </PrivateRoute>
          }
        />
        <Route
          path="/verify-email"
          element={
            <PrivateRoute>
              <EmailVerification />
            </PrivateRoute>
          }
        />
        <Route
          path="/sessions"
          element={
            <PrivateRoute>
              <Sessions />
            </PrivateRoute>
          }
        />
        <Route
          path="/permissions"
          element={
            <PrivateRoute>
              <Permissions />
            </PrivateRoute>
          }
        />
        <Route
          path="/audit-logs"
          element={
            <PrivateRoute>
              <AuditLogs />
            </PrivateRoute>
          }
        />
        <Route
          path="/oauth2-clients"
          element={
            <PrivateRoute>
              <OAuth2Clients />
            </PrivateRoute>
          }
        />
        <Route
          path="/authorized-apps"
          element={
            <PrivateRoute>
              <AuthorizedApps />
            </PrivateRoute>
          }
        />
        <Route path="/oauth2/consent" element={<ConsentPage />} />
        {/* 管理员后台路由 */}
        <Route
          path="/admin"
          element={
            <AdminRoute>
              <AdminLayout>
                <AdminDashboard />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/users"
          element={
            <AdminRoute>
              <AdminLayout>
                <UserManagement />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/roles"
          element={
            <AdminRoute>
              <AdminLayout>
                <RoleManagement />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/permissions"
          element={
            <AdminRoute>
              <AdminLayout>
                <PermissionManagement />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/audit-logs"
          element={
            <AdminRoute>
              <AdminLayout>
                <AuditLogs />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/oauth2-clients"
          element={
            <AdminRoute>
              <AdminLayout>
                <OAuth2Clients />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/settings"
          element={
            <AdminRoute>
              <AdminLayout>
                <SystemSettings />
              </AdminLayout>
            </AdminRoute>
          }
        />
        <Route
          path="/admin/backup"
          element={
            <AdminRoute>
              <AdminLayout>
                <BackupManagement />
              </AdminLayout>
            </AdminRoute>
          }
        />

        <Route path="/" element={<Navigate to="/dashboard" />} />
        </Routes>
      </BrowserRouter>
    </ErrorBoundary>
  )
}

export default App

