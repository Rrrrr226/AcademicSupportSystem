import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import Login from './pages/Login';
import Subjects from './pages/Subjects';
import ProfilePage from './pages/ProfilePage';
import LoginCallback from './pages/LoginCallback';
import AdminDashboard from './pages/AdminDashboard';

const theme = {
  token: {
    colorPrimary: '#1890ff',
    borderRadius: 12,
  },
};

function App() {
  return (
    <ConfigProvider locale={zhCN} theme={theme}>
      <Router>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/login/callback" element={<LoginCallback />} />
          <Route path="/subjects" element={<Subjects />} />
          <Route path="/profile" element={<ProfilePage />} />
          
          {/* 管理员路由 */}
          <Route path="/admin/dashboard" element={<AdminDashboard />} />
          
          <Route path="*" element={<Navigate to="/login" />} />
        </Routes>
      </Router>
    </ConfigProvider>
  );
}

export default App;