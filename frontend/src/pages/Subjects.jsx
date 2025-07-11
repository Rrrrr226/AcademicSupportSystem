import React, { useEffect, useState } from 'react';
import { Card, Button, Spin, message, Row, Col } from 'antd';
import { getSubjectLinks } from '../api';
import { useNavigate } from 'react-router-dom';

const Subjects = () => {
  const [links, setLinks] = useState([]);
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
    getSubjectLinks(staffId, token)
      .then(res => {
        setLinks(res.data.data || []);
      })
      .catch(() => {
        message.error('获取学科失败');
      })
      .finally(() => setLoading(false));
  }, [navigate]);

  return (
    <div style={{ minHeight: '100vh', background: '#e6f7ff', padding: 40 }}>
      <Card style={{ borderRadius: 16, maxWidth: 800, margin: '0 auto' }}>
        <h2 style={{ color: '#1890ff', textAlign: 'center' }}>我的学科</h2>
        {loading ? (
          <Spin />
        ) : (
          <Row gutter={[16, 16]}>
            {links.length === 0 ? (
              <Col span={24} style={{ textAlign: 'center', color: '#999' }}>暂无学科</Col>
            ) : (
              links.map((link, idx) => (
                <Col xs={24} sm={12} md={8} key={idx}>
                  <Button
                    type="primary"
                    block
                    style={{ marginBottom: 12, borderRadius: 8, background: '#40a9ff', border: 'none' }}
                    onClick={() => window.open(link, '_blank')}
                  >
                    跳转到学科{idx + 1}
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