import React, { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { Spin, message, Result, Button } from 'antd';
import { thirdPartyCallback } from '../api';

/**
 * HDUHelp 三方登录回调页面
 * 处理从第三方授权服务器返回的 code 和 state
 */
const LoginCallback = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const handleCallback = async () => {
      // 从 URL 获取授权码和状态参数
      const code = searchParams.get('code');
      const state = searchParams.get('state');
      const ticket = searchParams.get('ticket');

      if (!code) {
        setError('授权失败：未获取到授权码');
        setLoading(false);
        return;
      }

      try {
        // 用 code 换取 token
        const res = await thirdPartyCallback({ code, state, ticket });
        
        if (res.data && res.data.data && res.data.data.token) {
          // 登录成功，保存 token
          localStorage.setItem('token', res.data.data.token.trim());
          
          // 如果返回了用户信息，也保存
          if (res.data.data.staffId) {
            localStorage.setItem('staffId', res.data.data.staffId);
          }
          
          // 检查是否是管理员
          const isManager = res.data.data.isManager || false;
          
          // 保存用户信息（包含 isManager 状态）
          const userInfo = {
            ...(res.data.data.userInfo || {}),
            isManager: isManager,
            staffId: res.data.data.staffId,
          };
          localStorage.setItem('userInfo', JSON.stringify(userInfo));

          message.success('登录成功');
          
          // 管理员也保存到 adminToken
          if (isManager) {
            localStorage.setItem('adminToken', res.data.data.token.trim());
            if (res.data.data.refreshToken) {
              localStorage.setItem('adminRefreshToken', res.data.data.refreshToken);
            }
            navigate('/admin/dashboard', { replace: true });
          } else {
            // 普通用户 - 跳转到用户页面
            const redirectPath = localStorage.getItem('loginFrom') || '/subjects';
            localStorage.removeItem('loginFrom');
            navigate(redirectPath, { replace: true });
          }
        } else {
          throw new Error('登录返回数据异常');
        }
      } catch (err) {
        console.error('三方登录失败:', err);
        const errorMsg = err.response?.data?.message || err.message || '登录失败，请重试';
        setError(errorMsg);
        message.error(errorMsg);
      } finally {
        setLoading(false);
      }
    };

    handleCallback();
  }, [searchParams, navigate]);

  if (loading) {
    return (
      <div style={{ 
        minHeight: '100vh', 
        display: 'flex', 
        flexDirection: 'column',
        alignItems: 'center', 
        justifyContent: 'center',
        background: '#e6f7ff'
      }}>
        <Spin size="large" />
        <p style={{ marginTop: 16, color: '#666' }}>正在登录，请稍候...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div style={{ 
        minHeight: '100vh', 
        display: 'flex', 
        alignItems: 'center', 
        justifyContent: 'center',
        background: '#e6f7ff'
      }}>
        <Result
          status="error"
          title="登录失败"
          subTitle={error}
          extra={[
            <Button type="primary" key="retry" onClick={() => navigate('/login')}>
              返回登录
            </Button>,
          ]}
        />
      </div>
    );
  }

  return null;
};

export default LoginCallback;
