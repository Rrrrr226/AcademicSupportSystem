import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { login } from '../api';
import { useNavigate } from 'react-router-dom';

const Login = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      const res = await login(values);
      message.success('登录成功');
      localStorage.setItem('token', res.data.data.accessToken);
      localStorage.setItem('staffId', res.data.data.staffId || values.username);
      navigate('/subjects');
    } catch (err) {
      message.error(err.response?.data?.message || '登录失败');
    }
    setLoading(false);
  };

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <Card style={{ width: 350, borderRadius: 16 }}>
        <h2 style={{ color: '#1890ff', textAlign: 'center' }}>用户登录</h2>
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading} style={{ borderRadius: 8 }}>
              登录
            </Button>
          </Form.Item>
          <Form.Item>
            <Button type="link" block onClick={() => navigate('/register')}>
              没有账号？去注册
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Login;