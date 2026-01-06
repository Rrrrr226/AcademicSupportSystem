import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Space, Popconfirm, Typography, message } from 'antd';
import { DeleteOutlined, PlusOutlined, EditOutlined, SearchOutlined, ReloadOutlined } from '@ant-design/icons';
import { getUserSubjectList, addUserSubject, deleteUserSubject, updateUserSubject } from '../../api';

const { Title } = Typography;

const StudentSubjectsTab = () => {
  const [userSubjects, setUserSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [filters, setFilters] = useState({ staffId: '', subjectName: '' });
  const [modalVisible, setModalVisible] = useState(false);
  const [editingUserSubject, setEditingUserSubject] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchUserSubjects(1, 10, filters);
  }, []);

  const fetchUserSubjects = async (page = 1, pageSize = 10, searchFilters = {}) => {
    setLoading(true);
    const token = localStorage.getItem('adminToken');
    try {
      const response = await getUserSubjectList({
        page,
        pageSize,
        ...searchFilters
      }, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setUserSubjects(response.data.data?.user_subjects || []);
        setPagination({
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
      setLoading(false);
    }
  };

  const handleAddUserSubject = () => {
    setEditingUserSubject(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEditUserSubject = (record) => {
    setEditingUserSubject(record);
    form.setFieldsValue({
      staffId: record.staff_id,
      subjectName: record.subject_name
    });
    setModalVisible(true);
  };

  const handleDeleteUserSubject = async (id) => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await deleteUserSubject(id, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchUserSubjects(pagination.current, pagination.pageSize, filters);
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
      const values = await form.validateFields();
      const token = localStorage.getItem('adminToken');
      
      if (editingUserSubject) {
        const response = await updateUserSubject({
          id: editingUserSubject.id,
          ...values
        }, token);
        if (response.data?.code === 0 || response.data?.code === 200) {
          message.success('更新成功');
          setModalVisible(false);
          fetchUserSubjects(pagination.current, pagination.pageSize, filters);
        } else {
          message.error(response.data?.message || '更新失败');
        }
      } else {
        const response = await addUserSubject(values, token);
        if (response.data?.code === 0 || response.data?.code === 200) {
          message.success('添加成功');
          setModalVisible(false);
          fetchUserSubjects(pagination.current, pagination.pageSize, filters);
        } else {
          message.error(response.data?.message || '添加失败');
        }
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleSearch = () => {
    fetchUserSubjects(1, pagination.pageSize, filters);
  };

  const handleReset = () => {
    setFilters({ staffId: '', subjectName: '' });
    fetchUserSubjects(1, pagination.pageSize, { staffId: '', subjectName: '' });
  };

  const handleTableChange = (pagination) => {
    fetchUserSubjects(pagination.current, pagination.pageSize, filters);
  };

  const columns = [
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
  ];

  return (
    <div>
      <div style={{ marginBottom: 16 }}>
        <Title level={4}>学生科目管理 (共 {userSubjects.length} 条)</Title>
        <div style={{ display: 'flex', gap: 16, marginBottom: 16 }}>
          <Input
            placeholder="学号"
            value={filters.staffId}
            onChange={(e) => setFilters({ ...filters, staffId: e.target.value })}
            style={{ width: 200 }}
            prefix={<SearchOutlined />}
          />
          <Input
            placeholder="科目名称"
            value={filters.subjectName}
            onChange={(e) => setFilters({ ...filters, subjectName: e.target.value })}
            style={{ width: 200 }}
            prefix={<SearchOutlined />}
          />
          <Button type="primary" onClick={handleSearch} icon={<SearchOutlined />}>
            搜索
          </Button>
          <Button onClick={handleReset} icon={<ReloadOutlined />}>
            重置
          </Button>
          <Button type="primary" onClick={handleAddUserSubject} icon={<PlusOutlined />}>
            添加
          </Button>
        </div>
      </div>

      <Table
        columns={columns}
        dataSource={userSubjects}
        rowKey="id"
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
      />

      <Modal
        title={editingUserSubject ? "编辑学生科目" : "添加学生科目"}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingUserSubject(null);
        }}
        footer={null}
      >
        <Form
          form={form}
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
                setModalVisible(false);
                form.resetFields();
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
    </div>
  );
};

export default StudentSubjectsTab;
