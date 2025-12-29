import React, { useEffect, useState } from 'react';
import { Card, Button, Spin, message, Row, Col } from 'antd';
import { getSubjectLink } from '../api/subjects';
import { useNavigate } from 'react-router-dom';

const Subjects = () => {
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    const token = localStorage.getItem('token');
    const staffId = localStorage.getItem('staffId');
    if (!token || !staffId) {
      message.warning('请先登录');
      navigate('/login');
      return;
    }
    getSubjectLink(staffId)
      .then(res => {
        let subjects = res.data?.data?.subjects;
        if (Array.isArray(subjects)) {
          subjects = subjects.map(item => ({
            subject_name: item.subject_name || item.name || item.SubjectName,
            subject_link: item.subject_link || item.link || item.SubjectLink,
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

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', padding: 40, position: 'relative' }}>
      <div style={{ position: 'absolute', top: 24, right: 60, zIndex: 10 }}>
        <Button type="link" onClick={() => navigate('/profile')} style={{ fontWeight: 500, fontSize: 16 }}>个人信息</Button>
      </div>
      <Card style={{ borderRadius: 16, maxWidth: 800, margin: '0 auto' }}>
        <h2 style={{ color: '#1890ff', textAlign: 'center' }}>我的学科</h2>
        {loading ? (
          <Spin />
        ) : (
          <Row gutter={[16, 16]}>
            {subjects.length === 0 ? (
              <Col span={24} style={{ textAlign: 'center', color: '#999' }}>暂无学科</Col>
            ) : (
              subjects.map((item, idx) => (
                <Col xs={24} sm={12} md={8} key={idx}>
                  <Button
                    type="primary"
                    block
                    style={{ marginBottom: 12, borderRadius: 8, background: '#40a9ff', border: 'none' }}
                    onClick={() => item.subject_link ? window.open(item.subject_link, '_blank') : message.warning('无效链接')}
                    disabled={!item.subject_link}
                  >
                    {item.subject_name || '未命名学科'}
                  </Button>
                </Col>
              ))
            )}
          </Row>
        )}
      </Card>
    </div>
  );
};

export default Subjects;