import { BrowserRouter, Routes, Route, Navigate } from "react-router-dom";
import { AuthProvider } from "./lib/auth";

import AuthLayout from "./components/layout/AuthLayout";
import PlatformLayout from "./components/layout/PlatformLayout";

import LoginPage from "./pages/auth/LoginPage";
import RegisterPage from "./pages/auth/RegisterPage";
import DashboardPage from "./pages/dashboard/DashboardPage";
import ModuleCataloguePage from "./pages/modules/ModuleCataloguePage";
import ModuleOverviewPage from "./pages/modules/ModuleOverviewPage";
import SlidePlayerPage from "./pages/modules/SlidePlayerPage";
import MyLearningPage from "./pages/my-learning/MyLearningPage";
import CertificatesPage from "./pages/certificates/CertificatesPage";
import ProfilePage from "./pages/profile/ProfilePage";
import NotificationsPage from "./pages/notifications/NotificationsPage";

import AdminDashboardPage from "./pages/admin/AdminDashboardPage";
import ModuleManagementPage from "./pages/admin/ModuleManagementPage";
import ModuleBuilderPage from "./pages/admin/ModuleBuilderPage";
import NewModulePage from "./pages/admin/NewModulePage";
import CategoryManagementPage from "./pages/admin/CategoryManagementPage";
import UserManagementPage from "./pages/admin/UserManagementPage";
import EnrollmentManagementPage from "./pages/admin/EnrollmentManagementPage";
import ReportsPage from "./pages/admin/ReportsPage";
import AuditLogPage from "./pages/admin/AuditLogPage";
import SettingsPage from "./pages/admin/SettingsPage";

function App() {
  return (
    <BrowserRouter>
      <AuthProvider>
        <Routes>
          {/* Auth routes */}
          <Route element={<AuthLayout />}>
            <Route path="/login" element={<LoginPage />} />
            <Route path="/register" element={<RegisterPage />} />
          </Route>

          {/* Platform routes (authenticated) */}
          <Route element={<PlatformLayout />}>
            <Route path="/dashboard" element={<DashboardPage />} />
            <Route path="/modules" element={<ModuleCataloguePage />} />
            <Route path="/modules/:slug" element={<ModuleOverviewPage />} />
            <Route path="/my-learning" element={<MyLearningPage />} />
            <Route path="/certificates" element={<CertificatesPage />} />
            <Route path="/profile" element={<ProfilePage />} />
            <Route path="/notifications" element={<NotificationsPage />} />

            {/* Admin routes */}
            <Route path="/admin" element={<AdminDashboardPage />} />
            <Route path="/admin/modules" element={<ModuleManagementPage />} />
            <Route path="/admin/modules/new" element={<NewModulePage />} />
            <Route path="/admin/modules/:id/builder" element={<ModuleBuilderPage />} />
            <Route path="/admin/categories" element={<CategoryManagementPage />} />
            <Route path="/admin/users" element={<UserManagementPage />} />
            <Route path="/admin/enrollments" element={<EnrollmentManagementPage />} />
            <Route path="/admin/reports" element={<ReportsPage />} />
            <Route path="/admin/audit" element={<AuditLogPage />} />
            <Route path="/admin/settings" element={<SettingsPage />} />
          </Route>

          {/* Slide player (full screen, outside normal layout) */}
          <Route path="/modules/:slug/learn" element={<SlidePlayerPage />} />
          <Route path="/modules/:slug/learn/:slideId" element={<SlidePlayerPage />} />

          {/* Redirects */}
          <Route path="/" element={<Navigate to="/dashboard" replace />} />
          <Route path="*" element={<Navigate to="/dashboard" replace />} />
        </Routes>
      </AuthProvider>
    </BrowserRouter>
  );
}

export default App;
