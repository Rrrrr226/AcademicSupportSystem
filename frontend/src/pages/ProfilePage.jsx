import React, { useEffect, useState } from 'react';
import { Card, Button, Form, Input, message, Spin } from 'antd';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';

const ProfilePage = () => {
  const [userInfo, setUserInfo] = useState(null);
  const [loading, setLoading] = useState(true);
  const [pwdLoading, setPwdLoading] = useState(false);
  const [form] = Form.useForm();
  const navigate = useNavigate();

  useEffect(() => {
    // 获取个人信息
    setLoading(true);
    axios.get('/user/v1/info', {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    })
      .then(res => {
        setUserInfo(res.data.data);
      })
      .catch((err) => {
        console.error(err);
        message.error('获取个人信息失败');
      })
      .finally(() => setLoading(false));
  }, []);

  const onFinish = (values) => {
    setPwdLoading(true);
    axios.post('/user/v1/direct/modify', {
      userId: userInfo?.id,
      password: values.password
    }, {
      headers: { Authorization: 'Bearer ' + localStorage.getItem('token') }
    })
      .then(() => {
        message.success('密码修改成功，请重新登录');
        setTimeout(() => {
          localStorage.clear();
          navigate('/login');
        }, 1500);
      })
      .catch(err => {
        message.error(err.response?.data?.message || '密码修改失败');
      })
      .finally(() => setPwdLoading(false));
  };

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', padding: 40, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
      <Card style={{ borderRadius: 16, minWidth: 400 }}>
        <h2 style={{ color: '#1890ff', textAlign: 'center' }}>个人信息</h2>
        {loading ? <Spin /> : userInfo && (
          <div style={{ marginBottom: 32 }}>
            <p><b>用户名：</b>{userInfo.username}</p>
            <p><b>姓名：</b>{userInfo.name}</p>
            <p><b>邮箱：</b>{userInfo.email || '-'}</p>
            <p><b>手机号：</b>{userInfo.phone || '-'}</p>
          </div>
        )}
        <h3 style={{ color: '#1890ff' }}>修改密码</h3>
        <Form form={form} onFinish={onFinish} layout="vertical">
          <Form.Item name="password" label="新密码" rules={[{ required: true, message: '请输入新密码' }]}>
            <Input.Password placeholder="请输入新密码" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" loading={pwdLoading} block>修改密码</Button>
          </Form.Item>
        </Form>
        <Button type="link" onClick={() => navigate(-1)} style={{ marginTop: 8 }}>返回</Button>
      </Card>
    </div>
  );
};

export default ProfilePage;
