import React, { useState, useEffect } from 'react';
import { Layout, Menu, Typography, Space, Button, message } from 'antd';
import {
  UserOutlined, LogoutOutlined, TeamOutlined,
  FileExcelOutlined, BookOutlined, AppstoreOutlined, HomeOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { getManagerInfo } from '../api';

import ImportTab from './admin/ImportTab';
import StudentSubjectsTab from './admin/StudentSubjectsTab';
import SubjectsTab from './admin/SubjectsTab';
import ManagersTab from './admin/ManagersTab';
import FastGPTAppsTab from './admin/FastGPTAppsTab';

const { Header, Content, Sider } = Layout;
const { Title, Text } = Typography;


const AdminDashboard = () => {
  const navigate = useNavigate();
  const [currentUser, setCurrentUser] = useState(null);
  const [activeTab, setActiveTab] = useState(localStorage.getItem('adminActiveTab') || 'import');
  const [isAuthenticated, setIsAuthenticated] = useState(false);

  useEffect(() => {
    localStorage.setItem('adminActiveTab', activeTab);
  }, [activeTab]);

  useEffect(() => {
    const token = localStorage.getItem('adminToken') || localStorage.getItem('token');
    if (!token) {
      message.error('请先登录');
      navigate('/login', { replace: true });
      return;
    }
    // 如果 adminToken 不存在，从 token 复制
    if (!localStorage.getItem('adminToken') && localStorage.getItem('token')) {
      localStorage.setItem('adminToken', localStorage.getItem('token'));
    }
    setIsAuthenticated(true);
    fetchCurrentUser(token);
  }, [navigate]);

  const fetchCurrentUser = async (token) => {
    if (!token) token = localStorage.getItem('adminToken');
    try {
      const response = await getManagerInfo(token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setCurrentUser(response.data.data);
      }
    } catch (error) {
      console.error('获取用户信息失败:', error);
      if (error.response?.status === 401) {
        handleLogout();
      }
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('adminToken');
    localStorage.removeItem('adminRefreshToken');
    localStorage.removeItem('adminTokenExpireIn');
    localStorage.removeItem('adminLoginTime');
    navigate('/admin/login');
  };

  const renderContent = () => {
    switch (activeTab) {
      case 'import':
        return <ImportTab />;
      case 'student-subjects':
        return <StudentSubjectsTab />;
      case 'subjects':
        return <SubjectsTab />;
      case 'managers':
        return <ManagersTab currentUser={currentUser} />;
      case 'fastgpt-apps':
        return <FastGPTAppsTab />;
      default:
        return <ImportTab />;
    }
  };

  // 如果未认证，显示加载状态
  if (!isAuthenticated) {
    return (
      <div style={{ 
        minHeight: '100vh', 
        display: 'flex', 
        justifyContent: 'center', 
        alignItems: 'center' 
      }}>
        <Text>正在验证登录状态...</Text>
      </div>
    );
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{
        background: '#fff',
        padding: '0 24px',
        display: 'flex',
        justifyContent: 'space-between',
        alignItems: 'center',
        boxShadow: '0 2px 8px rgba(0,0,0,0.1)'
      }}>
        <Title level={4} style={{ margin: 0 }}>学业辅助系统 - 管理后台</Title>
        <Space>
          <Text>欢迎, {currentUser?.name || '管理员'}</Text>
          <Button icon={<HomeOutlined />} onClick={() => navigate('/subjects')}>进入学生端</Button>
          <Button icon={<LogoutOutlined />} onClick={handleLogout}>退出</Button>
        </Space>
      </Header>

      <Layout>
        <Sider width={200} style={{ background: '#fff' }}>
          <Menu
            mode="inline"
            selectedKeys={[activeTab]}
            style={{ height: '100%', borderRight: 0 }}
            onClick={({ key }) => setActiveTab(key)}
            items={[
              {
                key: 'import',
                icon: <FileExcelOutlined />,
                label: '导入学生科目',
              },
              {
                key: 'student-subjects',
                icon: <UserOutlined />,
                label: '学生科目管理',
              },
              {
                key: 'managers',
                icon: <TeamOutlined />,
                label: '管理员管理',
              },
              {
                key: 'fastgpt-apps',
                icon: <BookOutlined />,
                label: '学科管理',
              },
            ]}
          />
        </Sider>

        <Layout style={{ padding: '24px' }}>
          <Content style={{ background: '#fff', padding: 24, borderRadius: 8 }}>
            {renderContent()}
          </Content>
        </Layout>
      </Layout>
    </Layout>
  );
};

export default AdminDashboard;
