import React, { useState, useEffect } from 'react';
import {
  Layout, Menu, Card, Table, Button, Modal, Form, Input, Upload, message,
  Typography, Space, Popconfirm, Tabs, Statistic, Row, Col, Alert
} from 'antd';
import {
  UserOutlined, UploadOutlined, LogoutOutlined, TeamOutlined,
  FileExcelOutlined, PlusOutlined, DeleteOutlined, BookOutlined, EditOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import {
  getManagerInfo, getManagerList, addManager, deleteManager, importStudentSubjects,
  getSubjectList, addSubject, deleteSubject, updateSubject
} from '../api';

const { Header, Content, Sider } = Layout;
const { Title, Text } = Typography;
const { Dragger } = Upload;

const AdminDashboard = () => {
  const navigate = useNavigate();
  const [currentUser, setCurrentUser] = useState(null);
  const [managers, setManagers] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [subjectLoading, setSubjectLoading] = useState(false);
  const [addModalVisible, setAddModalVisible] = useState(false);
  const [subjectModalVisible, setSubjectModalVisible] = useState(false);
  const [editingSubject, setEditingSubject] = useState(null);
  const [importResult, setImportResult] = useState(null);
  const [activeTab, setActiveTab] = useState('import');
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [form] = Form.useForm();
  const [subjectForm] = Form.useForm();

  useEffect(() => {
    const token = localStorage.getItem('adminToken');
    if (!token) {
      message.error('请先登录');
      navigate('/admin/login', { replace: true });
      return;
    }
    setIsAuthenticated(true);
    fetchCurrentUser(token);
    fetchManagers(token);
    fetchSubjects(token);
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

  const fetchManagers = async (token) => {
    if (!token) token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getManagerList(token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setManagers(response.data.data?.managers || []);
      }
    } catch (error) {
      console.error('获取管理员列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchSubjects = async (token) => {
    if (!token) token = localStorage.getItem('adminToken');
    setSubjectLoading(true);
    try {
      const response = await getSubjectList(token, 1, 100); // 获取前100个，后续可加分页
      if (response.data?.code === 0 || response.data?.code === 200) {
        setSubjects(response.data.data?.subjects || []);
      }
    } catch (error) {
      console.error('获取学科列表失败:', error);
    } finally {
      setSubjectLoading(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('adminToken');
    localStorage.removeItem('adminRefreshToken');
    localStorage.removeItem('adminTokenExpireIn');
    localStorage.removeItem('adminLoginTime');
    navigate('/admin/login');
  };

  const handleAddManager = async (values) => {
    const token = localStorage.getItem('adminToken');
    try {
      // 只传 staffId
      const response = await addManager(values.staffId, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('添加管理员成功');
        setAddModalVisible(false);
        form.resetFields();
        fetchManagers();
      } else {
        message.error(response.data?.message || '添加失败');
      }
    } catch (error) {
      console.error('添加管理员失败:', error);
      message.error(error.response?.data?.message || '添加失败');
    }
  };

  const handleDeleteManager = async (staffId) => {
    const token = localStorage.getItem('adminToken');
    try {
      // 使用 staffId 而不是 managerId
      const response = await deleteManager(staffId, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchManagers();
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除管理员失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const handleImportExcel = async (file) => {
    const token = localStorage.getItem('adminToken');
    setImportResult(null);
    try {
      const response = await importStudentSubjects(file, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setImportResult(response.data.data);
        message.success('导入完成');
      } else {
        message.error(response.data?.message || '导入失败');
      }
    } catch (error) {
      console.error('导入失败:', error);
      message.error(error.response?.data?.message || '导入失败');
    }
    return false; // 阻止默认上传行为
  };

  const handleAddOrUpdateSubject = async (values) => {
    const token = localStorage.getItem('adminToken');
    try {
      let response;
      if (editingSubject) {
        response = await updateSubject({
          SubjectId: editingSubject.ID,
          SubjectName: values.subjectName,
          SubjectLink: values.subjectLink
        }, token);
      } else {
        response = await addSubject({
          subject_name: values.subjectName,
          subject_link: values.subjectLink
        }, token);
      }

      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success(editingSubject ? '更新成功' : '添加成功');
        setSubjectModalVisible(false);
        subjectForm.resetFields();
        setEditingSubject(null);
        fetchSubjects();
      } else {
        message.error(response.data?.message || '操作失败');
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleDeleteSubject = async (id) => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await deleteSubject(id, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchSubjects();
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除学科失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const openAddSubjectModal = () => {
    setEditingSubject(null);
    subjectForm.resetFields();
    setSubjectModalVisible(true);
  };

  const openEditSubjectModal = (record) => {
    setEditingSubject(record);
    subjectForm.setFieldsValue({
      subjectName: record.SubjectName,
      subjectLink: record.SubjectLink
    });
    setSubjectModalVisible(true);
  };

  const uploadProps = {
    name: 'file',
    multiple: false,
    accept: '.xlsx,.xls',
    beforeUpload: handleImportExcel,
    showUploadList: false,
  };

  const subjectColumns = [
    {
      title: 'ID',
      dataIndex: 'ID',
      key: 'ID',
      width: 80,
    },
    {
      title: '学科名称',
      dataIndex: 'SubjectName',
      key: 'SubjectName',
    },
    {
      title: '学科链接',
      dataIndex: 'SubjectLink',
      key: 'SubjectLink',
      render: (text) => <a href={text} target="_blank" rel="noopener noreferrer">{text}</a>
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => openEditSubjectModal(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此学科吗？"
            onConfirm={() => handleDeleteSubject(record.ID)}
            okText="确定"
            cancelText="取消"
          >
            <Button
              type="link"
              danger
              icon={<DeleteOutlined />}
            >
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  const managerColumns = [
    {
      title: '学号/工号',
      dataIndex: 'staffId',
      key: 'staffId',
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Popconfirm
          title="确定要删除此管理员吗？"
          onConfirm={() => handleDeleteManager(record.staffId)}
          okText="确定"
          cancelText="取消"
          disabled={record.staffId === currentUser?.staffId}
        >
          <Button
            type="link"
            danger
            icon={<DeleteOutlined />}
            disabled={record.staffId === currentUser?.staffId}
          >
            删除
          </Button>
        </Popconfirm>
      ),
    },
  ];

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
                key: 'subjects',
                icon: <BookOutlined />,
                label: '学科管理',
              },
              {
                key: 'managers',
                icon: <TeamOutlined />,
                label: '管理员管理',
              },
            ]}
          />
        </Sider>

        <Layout style={{ padding: '24px' }}>
          <Content style={{ background: '#fff', padding: 24, borderRadius: 8 }}>
            {activeTab === 'import' && (
              <div>
                <Title level={4}>导入学生科目名单</Title>
                <Alert
                  message="Excel 格式要求"
                  description={
                    <div>
                      <p>Excel 文件需要包含以下列：</p>
                      <ul>
                        <li><strong>学号</strong>（或 staff_id、StaffId）- 学生学号</li>
                        <li><strong>科目名称</strong>（或 科目、subject_name、SubjectName）- 科目名称（需要与系统中已有科目名称一致）</li>
                      </ul>
                      <p>每行代表一个学生-科目的对应关系，同一学生可以有多行对应不同科目。</p>
                    </div>
                  }
                  type="info"
                  showIcon
                  style={{ marginBottom: 24 }}
                />

                <Dragger {...uploadProps} style={{ marginBottom: 24 }}>
                  <p className="ant-upload-drag-icon">
                    <UploadOutlined style={{ fontSize: 48, color: '#1890ff' }} />
                  </p>
                  <p className="ant-upload-text">点击或拖拽 Excel 文件到此区域上传</p>
                  <p className="ant-upload-hint">
                    支持 .xlsx 和 .xls 格式
                  </p>
                </Dragger>

                {importResult && (
                  <Card title="导入结果" style={{ marginTop: 24 }}>
                    <Row gutter={16}>
                      <Col span={8}>
                        <Statistic title="总记录数" value={importResult.total} />
                      </Col>
                      <Col span={8}>
                        <Statistic
                          title="成功数"
                          value={importResult.successCount}
                          valueStyle={{ color: '#3f8600' }}
                        />
                      </Col>
                      <Col span={8}>
                        <Statistic
                          title="失败数"
                          value={importResult.failCount}
                          valueStyle={{ color: importResult.failCount > 0 ? '#cf1322' : undefined }}
                        />
                      </Col>
                    </Row>
                    {importResult.errors && importResult.errors.length > 0 && (
                      <div style={{ marginTop: 16 }}>
                        <Title level={5}>错误详情：</Title>
                        <ul>
                          {importResult.errors.slice(0, 10).map((err, idx) => (
                            <li key={idx} style={{ color: '#cf1322' }}>{err}</li>
                          ))}
                          {importResult.errors.length > 10 && (
                            <li>...还有 {importResult.errors.length - 10} 个错误</li>
                          )}
                        </ul>
                      </div>
                    )}
                  </Card>
                )}
              </div>
            )}

            {activeTab === 'subjects' && (
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                  <Title level={4} style={{ margin: 0 }}>学科列表</Title>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={openAddSubjectModal}
                  >
                    添加学科
                  </Button>
                </div>

                <Table
                  columns={subjectColumns}
                  dataSource={subjects}
                  rowKey="ID"
                  loading={subjectLoading}
                />
              </div>
            )}

            {activeTab === 'managers' && (
              <div>
                <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
                  <Title level={4} style={{ margin: 0 }}>管理员列表</Title>
                  <Button
                    type="primary"
                    icon={<PlusOutlined />}
                    onClick={() => setAddModalVisible(true)}
                  >
                    添加管理员
                  </Button>
                </div>

                <Table
                  columns={managerColumns}
                  dataSource={managers}
                  rowKey="id"
                  loading={loading}
                />
              </div>
            )}
          </Content>
        </Layout>
      </Layout>

      {/* 添加/编辑学科模态框 */}
      <Modal
        title={editingSubject ? "编辑学科" : "添加学科"}
        open={subjectModalVisible}
        onCancel={() => {
          setSubjectModalVisible(false);
          subjectForm.resetFields();
          setEditingSubject(null);
        }}
        footer={null}
      >
        <Form
          form={subjectForm}
          layout="vertical"
          onFinish={handleAddOrUpdateSubject}
        >
          <Form.Item
            name="subjectName"
            label="学科名称"
            rules={[{ required: true, message: '请输入学科名称' }]}
          >
            <Input placeholder="请输入学科名称" />
          </Form.Item>
          <Form.Item
            name="subjectLink"
            label="学科链接"
            rules={[{ required: true, message: '请输入学科链接' }]}
          >
            <Input placeholder="请输入学科链接" />
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setSubjectModalVisible(false);
                subjectForm.resetFields();
                setEditingSubject(null);
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingSubject ? "确认更新" : "确认添加"}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>

      {/* 添加管理员模态框 */}
      <Modal
        title="添加管理员"
        open={addModalVisible}
        onCancel={() => {
          setAddModalVisible(false);
          form.resetFields();
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleAddManager}
        >
          <Form.Item
            name="staffId"
            label="学号/工号"
            rules={[{ required: true, message: '请输入学号/工号' }]}
          >
            <Input placeholder="请输入学号/工号" />
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setAddModalVisible(false);
                form.resetFields();
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                确认添加
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};

export default AdminDashboard;
