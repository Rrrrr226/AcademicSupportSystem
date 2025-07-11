import React, { useState } from 'react';
import { Form, Input, Button, Card, message } from 'antd';
import { register } from '../api';
import { useNavigate } from 'react-router-dom';

const Register = () => {
  const [loading, setLoading] = useState(false);
  const navigate = useNavigate();

  const onFinish = async (values) => {
    setLoading(true);
    try {
      await register(values);
      message.success('注册成功，请登录');
      navigate('/login');
    } catch (err) {
      message.error(err.response?.data?.message || '注册失败');
    }
    setLoading(false);
  };

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <Card style={{ width: 350, borderRadius: 16 }}>
        <h2 style={{ color: '#1890ff', textAlign: 'center' }}>用户注册</h2>
        <Form layout="vertical" onFinish={onFinish}>
          <Form.Item name="username" label="用户名" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱">
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block loading={loading} style={{ borderRadius: 8 }}>
              注册
            </Button>
          </Form.Item>
          <Form.Item>
            <Button type="link" block onClick={() => navigate('/login')}>
              已有账号？去登录
            </Button>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default Register;