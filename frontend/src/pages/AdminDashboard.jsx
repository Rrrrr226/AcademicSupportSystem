import React, { useState, useEffect } from 'react';
import {
  Layout, Menu, Card, Table, Button, Modal, Form, Input, Upload, message,
  Typography, Space, Popconfirm, Tabs, Statistic, Row, Col, Alert
} from 'antd';
import {
  UserOutlined, UploadOutlined, LogoutOutlined, TeamOutlined,
  FileExcelOutlined, PlusOutlined, DeleteOutlined, BookOutlined, EditOutlined, 
  DownloadOutlined, SearchOutlined, ReloadOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import {
  getManagerInfo, getManagerList, addManager, deleteManager, importStudentSubjects,
  getSubjectList, addSubject, deleteSubject, updateSubject, downloadImportTemplate,
  getUserSubjectList, addUserSubject, deleteUserSubject, updateUserSubject
} from '../api';

const { Header, Content, Sider } = Layout;
const { Title, Text } = Typography;
const { Dragger } = Upload;

const AdminDashboard = () => {
  const navigate = useNavigate();
  const [currentUser, setCurrentUser] = useState(null);
  const [managers, setManagers] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [userSubjects, setUserSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [subjectLoading, setSubjectLoading] = useState(false);
  const [userSubjectLoading, setUserSubjectLoading] = useState(false);
  const [addModalVisible, setAddModalVisible] = useState(false);
  const [subjectModalVisible, setSubjectModalVisible] = useState(false);
  const [userSubjectModalVisible, setUserSubjectModalVisible] = useState(false);
  const [editingSubject, setEditingSubject] = useState(null);
  const [editingUserSubject, setEditingUserSubject] = useState(null);
  const [importResult, setImportResult] = useState(null);
  const [activeTab, setActiveTab] = useState('import');
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [userSubjectPagination, setUserSubjectPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [managerPagination, setManagerPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [subjectPagination, setSubjectPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [userSubjectFilters, setUserSubjectFilters] = useState({ staffId: '', subjectName: '' });
  const [form] = Form.useForm();
  const [subjectForm] = Form.useForm();
  const [userSubjectForm] = Form.useForm();

  useEffect(() => {
    const token = localStorage.getItem('adminToken') || localStorage.getItem('token');
    if (!token) {
      message.error('请先登录');
      navigate('/login', { replace: true });
      return;
    }
    // 如果 adminToken 不存在，从 token 复制
    if (!localStorage.getItem('adminToken') && localStorage.getItem('token')) {
      localStorage.setItem('adminToken', localStorage.getItem('token'));
    }
    setIsAuthenticated(true);
    fetchCurrentUser(token);
    fetchManagers(1, 10);
    fetchSubjects(1, 10);
  }, [navigate]);

  useEffect(() => {
    if (activeTab === 'student-subjects') {
      fetchUserSubjects(userSubjectPagination.current, userSubjectPagination.pageSize, userSubjectFilters);
    } else if (activeTab === 'subjects') {
      fetchSubjects(subjectPagination.current, subjectPagination.pageSize);
    } else if (activeTab === 'managers') {
      fetchManagers(managerPagination.current, managerPagination.pageSize);
    }
  }, [activeTab]);

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

  const fetchManagers = async (page = 1, pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getManagerList(token, page, pageSize);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setManagers(response.data.data?.managers || []);
        setManagerPagination({
          current: response.data.data?.page || page,
          pageSize: response.data.data?.page_size || pageSize,
          total: response.data.data?.total || 0
        });
      }
    } catch (error) {
      console.error('获取管理员列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const fetchSubjects = async (page = 1,pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setSubjectLoading(true);
    try {
      const response = await getSubjectList(token, page, pageSize); 
      if (response.data?.code === 0 || response.data?.code === 200) {
        setSubjects(response.data.data?.subjects || []);
        setSubjectPagination({
            current: response.data.data?.page || page,
            pageSize: response.data.data?.page_size || pageSize,
            total: response.data.data?.total || 0
        });
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
          subject_id: editingSubject.id,
          subject_name: values.subjectName,
          subject_link: values.subjectLink
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
        fetchSubjects(subjectPagination.current, subjectPagination.pageSize);
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
        fetchSubjects(subjectPagination.current, subjectPagination.pageSize);
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
      subjectName: record.subject_name,
      subjectLink: record.subject_link
    });
    setSubjectModalVisible(true);
  };

  // User Subject Handlers
  const fetchUserSubjects = async (page = 1, pageSize = 10, filters = {}) => {
    console.log('fetchUserSubjects called', { page, pageSize, filters });
    setUserSubjectLoading(true);
    const token = localStorage.getItem('adminToken');
    try {
      const response = await getUserSubjectList({
        page,
        pageSize,
        ...filters
      }, token);
      console.log('getUserSubjectList response:', response);
      if (response.data?.code === 0 || response.data?.code === 200) {
        console.log('Setting userSubjects:', response.data.data?.user_subjects);
        setUserSubjects(response.data.data?.user_subjects || []);
        setUserSubjectPagination({
          current: response.data.data?.page || page,
          pageSize: response.data.data?.page_size || pageSize,
          total: response.data.data?.total || 0
        });
      } else {
        message.error(response.data?.message || '获取学生科目列表失败');
      }
    } catch (error) {
      console.error('获取学生科目列表失败:', error);
      message.error(error.response?.data?.message || '获取学生科目列表失败');
    } finally {
      setUserSubjectLoading(false);
    }
  };

  const handleAddUserSubject = () => {
    setEditingUserSubject(null);
    userSubjectForm.resetFields();
    setUserSubjectModalVisible(true);
  };

  const handleEditUserSubject = (record) => {
    setEditingUserSubject(record);
    userSubjectForm.setFieldsValue({
      staffId: record.staff_id,
      subjectName: record.subject_name
    });
    setUserSubjectModalVisible(true);
  };

  const handleDeleteUserSubject = async (id) => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await deleteUserSubject(id, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchUserSubjects(userSubjectPagination.current, userSubjectPagination.pageSize, userSubjectFilters);
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const handleUserSubjectSubmit = async () => {
    try {
      const values = await userSubjectForm.validateFields();
      const token = localStorage.getItem('adminToken');
      
      if (editingUserSubject) {
        const response = await updateUserSubject({
          id: editingUserSubject.id,
          ...values
        }, token);
        if (response.data?.code === 0 || response.data?.code === 200) {
          message.success('更新成功');
          setUserSubjectModalVisible(false);
          fetchUserSubjects(userSubjectPagination.current, userSubjectPagination.pageSize, userSubjectFilters);
        } else {
          message.error(response.data?.message || '更新失败');
        }
      } else {
        const response = await addUserSubject(values, token);
        if (response.data?.code === 0 || response.data?.code === 200) {
          message.success('添加成功');
          setUserSubjectModalVisible(false);
          fetchUserSubjects(userSubjectPagination.current, userSubjectPagination.pageSize, userSubjectFilters);
        } else {
          message.error(response.data?.message || '添加失败');
        }
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleUserSubjectSearch = () => {
    fetchUserSubjects(1, userSubjectPagination.pageSize, userSubjectFilters);
  };

  const handleUserSubjectReset = () => {
    setUserSubjectFilters({ staffId: '', subjectName: '' });
    fetchUserSubjects(1, userSubjectPagination.pageSize, { staffId: '', subjectName: '' });
  };

  const handleUserSubjectTableChange = (pagination) => {
    fetchUserSubjects(pagination.current, pagination.pageSize, userSubjectFilters);
  };

  const handleManagerTableChange = (pagination) => {
    fetchManagers(pagination.current, pagination.pageSize);
  };

  const handleSubjectTableChange = (pagination) => {
    fetchSubjects(pagination.current, pagination.pageSize);
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
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '学科名称',
      dataIndex: 'subject_name',
      key: 'subject_name',
    },
    {
      title: '学科链接',
      dataIndex: 'subject_link',
      key: 'subject_link',
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
            onConfirm={() => handleDeleteSubject(record.id)}
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
                key: 'student-subjects',
                icon: <UserOutlined />,
                label: '学生科目管理',
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
                  style={{ marginBottom: 16 }}
                />

                <div style={{ marginBottom: 24 }}>
                  <Button 
                    type="default" 
                    icon={<DownloadOutlined />}
                    onClick={() => window.open(downloadImportTemplate(), '_blank')}
                  >
                    下载导入模板
                  </Button>
                </div>

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
                  rowKey="id"
                  pagination={subjectPagination}
                  loading={subjectLoading}
                  onChange={handleSubjectTableChange}
                />
              </div>
            )}

            {activeTab === 'student-subjects' && (
              <div>
                <div style={{ marginBottom: 16 }}>
                  <Title level={4}>学生科目管理 (共 {userSubjects.length} 条)</Title>
                  <div style={{ display: 'flex', gap: 16, marginBottom: 16 }}>
                    <Input
                      placeholder="学号"
                      value={userSubjectFilters.staffId}
                      onChange={(e) => setUserSubjectFilters({ ...userSubjectFilters, staffId: e.target.value })}
                      style={{ width: 200 }}
                      prefix={<SearchOutlined />}
                    />
                    <Input
                      placeholder="科目名称"
                      value={userSubjectFilters.subjectName}
                      onChange={(e) => setUserSubjectFilters({ ...userSubjectFilters, subjectName: e.target.value })}
                      style={{ width: 200 }}
                      prefix={<SearchOutlined />}
                    />
                    <Button type="primary" onClick={handleUserSubjectSearch} icon={<SearchOutlined />}>
                      搜索
                    </Button>
                    <Button onClick={handleUserSubjectReset} icon={<ReloadOutlined />}>
                      重置
                    </Button>
                    <Button type="primary" onClick={handleAddUserSubject} icon={<PlusOutlined />}>
                      添加
                    </Button>
                  </div>
                </div>

                <Table
                  columns={[
                    {
                      title: 'ID',
                      dataIndex: 'id',
                      key: 'id',
                      width: 200,
                      ellipsis: true,
                    },
                    {
                      title: '学号',
                      dataIndex: 'staff_id',
                      key: 'staff_id',
                      width: 120,
                    },
                    {
                      title: '科目名称',
                      dataIndex: 'subject_name',
                      key: 'subject_name',
                      width: 150,
                    },
                    {
                      title: '创建时间',
                      dataIndex: 'created_at',
                      key: 'created_at',
                      width: 180,
                      render: (text) => text ? new Date(text).toLocaleString() : '-',
                    },
                    {
                      title: '操作',
                      key: 'action',
                      render: (_, record) => (
                        <Space size="middle">
                          <Button
                            type="link"
                            icon={<EditOutlined />}
                            onClick={() => handleEditUserSubject(record)}
                          >
                            编辑
                          </Button>
                          <Popconfirm
                            title="确定删除吗？"
                            onConfirm={() => handleDeleteUserSubject(record.id)}
                            okText="确定"
                            cancelText="取消"
                          >
                            <Button type="link" danger icon={<DeleteOutlined />}>
                              删除
                            </Button>
                          </Popconfirm>
                        </Space>
                      ),
                    },
                  ]}
                  dataSource={userSubjects}
                  rowKey="id"
                  loading={userSubjectLoading}
                  pagination={userSubjectPagination}
                  onChange={handleUserSubjectTableChange}
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
                  pagination={managerPagination}
                  onChange={handleManagerTableChange}
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

      {/* 添加/编辑学生科目模态框 */}
      <Modal
        title={editingUserSubject ? "编辑学生科目" : "添加学生科目"}
        open={userSubjectModalVisible}
        onCancel={() => {
          setUserSubjectModalVisible(false);
          userSubjectForm.resetFields();
          setEditingUserSubject(null);
        }}
        footer={null}
      >
        <Form
          form={userSubjectForm}
          layout="vertical"
          onFinish={handleUserSubjectSubmit}
        >
          <Form.Item
            name="staffId"
            label="学号"
            rules={[{ required: true, message: '请输入学号' }]}
          >
            <Input placeholder="请输入学号" disabled={!!editingUserSubject} />
          </Form.Item>
          <Form.Item
            name="subjectName"
            label="科目名称"
            rules={[{ required: true, message: '请输入科目名称' }]}
          >
            <Input placeholder="请输入科目名称" />
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setUserSubjectModalVisible(false);
                userSubjectForm.resetFields();
                setEditingUserSubject(null);
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingUserSubject ? '更新' : '添加'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </Layout>
  );
};

export default AdminDashboard;
