import React, { useState } from 'react';
import { Button, Card, message } from 'antd';
import { getThirdPartyJumpUrl } from '../api';
import { useLocation } from 'react-router-dom';

const Login = () => {
  const [loading, setLoading] = useState(false);
  const location = useLocation();

  // HDUHelp 三方登录
  const handleThirdPartyLogin = async () => {
    setLoading(true);
    try {
      // 保存当前来源页面，用于登录成功后跳转回来
      const from = location.state?.from?.pathname || '/subjects';
      localStorage.setItem('loginFrom', from);

      const res = await getThirdPartyJumpUrl('/');
      
      if (res.data && res.data.data && res.data.data.url) {
        // 跳转到第三方授权页面
        window.location.href = res.data.data.url;
      } else {
        throw new Error('获取登录地址失败');
      }
    } catch (err) {
      console.error('获取三方登录地址失败:', err);
      message.error(err.response?.data?.message || err.message || '获取登录地址失败');
      setLoading(false);
    }
  };

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', display: 'flex', alignItems: 'center', justifyContent: 'center' }}>
      <Card style={{ width: 350, borderRadius: 16, textAlign: 'center' }}>
        <h2 style={{ color: '#1890ff', marginBottom: 32 }}>用户登录</h2>
        
        <Button 
          type="primary"
          size="large"
          block 
          loading={loading}
          onClick={handleThirdPartyLogin}
          style={{ 
            borderRadius: 8,
            height: 48,
            fontSize: 16
          }}
        >
          HDUHelp 统一身份认证登录
        </Button>
        
        <p style={{ marginTop: 24, color: '#999', fontSize: 12 }}>
          点击上方按钮，使用学校统一身份认证登录
        </p>
      </Card>
    </div>
  );
};

export default Login;