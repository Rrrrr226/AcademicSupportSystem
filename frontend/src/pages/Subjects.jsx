import React, { useEffect, useState } from 'react';
import { Card, Button, Spin, message, Row, Col, Typography, Empty, Space } from 'antd';
import { BookOutlined, UserOutlined, ArrowRightOutlined, DashboardOutlined } from '@ant-design/icons';
import { getSubjectLink } from '../api/subjects';
import { useNavigate } from 'react-router-dom';

const { Title, Text } = Typography;

const Subjects = () => {
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isManager, setIsManager] = useState(false);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    const staffId = localStorage.getItem('staffId');
    const userInfo = localStorage.getItem('userInfo');
    
    if (!token || !staffId) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }

    // 检查是否是管理员
    try {
      const user = userInfo ? JSON.parse(userInfo) : {};
      setIsManager(user.isManager || false);
    } catch (e) {
      console.error('解析用户信息失败:', e);
    }

    getSubjectLink(staffId)
      .then(res => {
        let subjects = res.data?.data?.subjects;
        if (Array.isArray(subjects)) {
          subjects = subjects.map(item => ({
            subject_name: item.subject_name || item.name || item.SubjectName,
            subject_link: item.subject_link || item.link || item.SubjectLink,
            app_id: item.app_id
          }));
        } else {
          subjects = [];
        }
        setSubjects(subjects);
      })
      .catch((error) => {
        console.error('获取学科失败:', error);
        message.error(error.response?.data?.message || '获取学科失败');
      })
      .finally(() => setLoading(false));
  }, [navigate]);

  const handleSubjectClick = (item) => {
    if (item.app_id) {
      navigate('/chat', { state: { appId: item.app_id, title: item.subject_name } });
    } else if (item.subject_link) {
      window.open(item.subject_link, '_blank');
    } else {
      message.warning('无效链接');
    }
  };

  return (
    <div style={{ 
      minHeight: '100vh', 
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)', 
      padding: '40px 20px',
      display: 'flex',
      justifyContent: 'center',
      alignItems: 'flex-start'
    }}>
      <Card 
        style={{ 
          width: '100%', 
          maxWidth: 1000, 
          borderRadius: 16, 
          boxShadow: '0 10px 25px rgba(0,0,0,0.08)',
          minHeight: '80vh',
          display: 'flex',
          flexDirection: 'column'
        }}
        bodyStyle={{ flex: 1, display: 'flex', flexDirection: 'column', padding: 40 }}
      >
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: 32 }}>
          <div style={{ display: 'flex', alignItems: 'center' }}>
             <div style={{ 
               width: 48, height: 48, borderRadius: 12, background: '#e6f7ff', 
               display: 'flex', alignItems: 'center', justifyContent: 'center', marginRight: 16 
             }}>
               <BookOutlined style={{ fontSize: 24, color: '#1890ff' }} />
             </div>
             <div>
               <Title level={3} style={{ margin: 0 }}>我的学科</Title>
               <Text type="secondary">查看和访问您的课程资源</Text>
             </div>
          </div>
          <Space>
            {isManager && (
              <Button 
                type="primary" 
                shape="round" 
                icon={<DashboardOutlined />} 
                onClick={() => {
                  // 确保 adminToken 存在
                  if (!localStorage.getItem('adminToken')) {
                    localStorage.setItem('adminToken', localStorage.getItem('token'));
                  }
                  navigate('/admin/dashboard');
                }}
                size="large"
              >
                管理后台
              </Button>
            )}
            <Button 
              type="default" 
              shape="round" 
              icon={<UserOutlined />} 
              onClick={() => navigate('/profile')}
              size="large"
            >
              个人中心
            </Button>
          </Space>
        </div>

        {loading ? (
          <div style={{ flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <Spin size="large" tip="加载中..." />
          </div>
        ) : subjects.length === 0 ? (
          <div style={{ flex: 1, display: 'flex', justifyContent: 'center', alignItems: 'center' }}>
            <Empty description="暂无学科信息" image={Empty.PRESENTED_IMAGE_SIMPLE} />
          </div>
        ) : (
          <Row gutter={[24, 24]}>
            {subjects.map((item, idx) => (
              <Col xs={24} sm={12} md={8} lg={6} key={idx}>
                <Card
                  hoverable
                  style={{ 
                    borderRadius: 12, 
                    border: '1px solid #f0f0f0',
                    height: '100%',
                    transition: 'all 0.3s'
                  }}
                  bodyStyle={{ padding: 24, display: 'flex', flexDirection: 'column', alignItems: 'center', textAlign: 'center', height: '100%', justifyContent: 'center' }}
                  onClick={() => handleSubjectClick(item)}
                >
                  <div style={{ 
                    width: 64, height: 64, borderRadius: '50%', background: item.subject_link ? '#e6f7ff' : '#f5f5f5', 
                    display: 'flex', alignItems: 'center', justifyContent: 'center', marginBottom: 16
                  }}>
                    <BookOutlined style={{ fontSize: 28, color: item.subject_link ? '#1890ff' : '#bfbfbf' }} />
                  </div>
                  <Title level={5} style={{ marginBottom: 8, width: '100%', overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}>
                    {item.subject_name || '未命名学科'}
                  </Title>
                  {item.subject_link ? (
                    <Text type="secondary" style={{ fontSize: 12 }}>点击访问课程主页 <ArrowRightOutlined /></Text>
                  ) : (
                    <Text type="secondary" style={{ fontSize: 12 }}>暂无链接</Text>
                  )}
                </Card>
              </Col>
            ))}
          </Row>
        )}
      </Card>
    </div>
  );
};

export default Subjects;