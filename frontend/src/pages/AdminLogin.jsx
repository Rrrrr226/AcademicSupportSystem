import React, { useState } from 'react';
import { Form, Input, Button, Card, message, Typography } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import { managerLogin } from '../api';

const { Title } = Typography;

const AdminLogin = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const response = await managerLogin(values.username, values.password);
      console.log('登录响应:', response.data);
      
      if (response.data?.code === 200 && response.data?.data) {
        const { token, refreshToken, expireIn } = response.data.data;
        
        // 存储管理员 token（与普通用户区分）
        localStorage.setItem('adminToken', token);
        localStorage.setItem('adminRefreshToken', refreshToken);
        localStorage.setItem('adminTokenExpireIn', String(expireIn));
        localStorage.setItem('adminLoginTime', Date.now().toString());
        
        message.success('登录成功');
        // 使用 replace 避免返回到登录页
        setTimeout(() => {
          navigate('/admin/dashboard', { replace: true });
        }, 100);
      } else {
        message.error(response.data?.message || '登录失败');
      }
    } catch (error) {
      console.error('登录错误:', error);
      message.error(error.response?.data?.message || '登录失败，请检查用户名和密码');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{
      minHeight: '100vh',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'center',
      background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)'
    }}>
      <Card
        style={{
          width: 400,
          boxShadow: '0 4px 12px rgba(0, 0, 0, 0.15)',
          borderRadius: 8
        }}
      >
        <div style={{ textAlign: 'center', marginBottom: 24 }}>
          <Title level={3} style={{ marginBottom: 8 }}>管理员登录</Title>
          <Typography.Text type="secondary">学业辅助系统管理后台</Typography.Text>
        </div>
        
        <Form
          name="admin_login"
          onFinish={onFinish}
          autoComplete="off"
          size="large"
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: '请输入用户名' }]}
          >
            <Input
              prefix={<UserOutlined />}
              placeholder="用户名"
            />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: '请输入密码' }]}
          >
            <Input.Password
              prefix={<LockOutlined />}
              placeholder="密码"
            />
          </Form.Item>

          <Form.Item>
            <Button
              type="primary"
              htmlType="submit"
              loading={loading}
              block
              style={{ height: 44 }}
            >
              登录
            </Button>
          </Form.Item>
        </Form>

        <div style={{ textAlign: 'center' }}>
          <Button type="link" onClick={() => navigate('/login')}>
            返回用户登录
          </Button>
        </div>
      </Card>
    </div>
  );
};

export default AdminLogin;
