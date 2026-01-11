import React, { useState, useEffect } from 'react';
import {
  Card, Table, Button, Modal, Form, Input, message, Space, Popconfirm, Select
} from 'antd';
import {
  PlusOutlined, DeleteOutlined, EditOutlined, SearchOutlined, ReloadOutlined
} from '@ant-design/icons';
import { useNavigate } from 'react-router-dom';
import {
  getUserSubjectList, addUserSubject, deleteUserSubject, updateUserSubject, getSubjectList
} from '../api';

const StudentSubjectsManagement = () => {
  const navigate = useNavigate();
  const [dataSource, setDataSource] = useState([]);
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingRecord, setEditingRecord] = useState(null);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [filters, setFilters] = useState({ staffId: '', subjectName: '' });
  const [form] = Form.useForm();

  useEffect(() => {
    const token = localStorage.getItem('adminToken') || localStorage.getItem('token');
    if (!token) {
      message.error('请先登录');
      navigate('/login');
      return;
    }
    // 如果 adminToken 不存在，从 token 复制
    if (!localStorage.getItem('adminToken') && localStorage.getItem('token')) {
      localStorage.setItem('adminToken', localStorage.getItem('token'));
    }
    fetchData();
    fetchSubjects();
  }, [navigate]);

  const fetchData = async (page = 1, pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getUserSubjectList(
        {
          page,
          pageSize,
          staffId: filters.staffId,
          subjectName: filters.subjectName
        },
        token
      );
      if (response.data?.code === 0 || response.data?.code === 200) {
        const data = response.data.data;
        setDataSource(data.user_subjects || []);
        setPagination({
          current: data.page,
          pageSize: data.page_size,
          total: data.total,
        });
      }
    } catch (error) {
      console.error('获取数据失败:', error);
      message.error(error.response?.data?.message || '获取数据失败');
    } finally {
      setLoading(false);
    }
  };

  const fetchSubjects = async () => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await getSubjectList(token, 1, 100);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setSubjects(response.data.data?.subjects || []);
      }
    } catch (error) {
      console.error('获取学科列表失败:', error);
    }
  };

  const handleAdd = () => {
    setEditingRecord(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    setEditingRecord(record);
    form.setFieldsValue({
      staff_id: record.staff_id,
      subject_name: record.subject_name,
    });
    setModalVisible(true);
  };

  const handleDelete = async (id) => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await deleteUserSubject(id, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchData(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const handleSubmit = async (values) => {
    const token = localStorage.getItem('adminToken');
    try {
      let response;
      if (editingRecord) {
        response = await updateUserSubject({
          id: editingRecord.id,
          staff_id: values.staff_id,
          subject_name: values.subject_name
        }, token);
      } else {
        response = await addUserSubject({
          staff_id: values.staff_id,
          subject_name: values.subject_name
        }, token);
      }

      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success(editingRecord ? '更新成功' : '添加成功');
        setModalVisible(false);
        form.resetFields();
        fetchData(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '操作失败');
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleSearch = () => {
    fetchData(1, pagination.pageSize);
  };

  const handleReset = () => {
    setFilters({ staffId: '', subjectName: '' });
    setTimeout(() => fetchData(1, pagination.pageSize), 0);
  };

  const handleTableChange = (newPagination) => {
    fetchData(newPagination.current, newPagination.pageSize);
  };

  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '学号',
      dataIndex: 'staff_id',
      key: 'staff_id',
      width: 150,
    },
    {
      title: '科目名称',
      dataIndex: 'subject_name',
      key: 'subject_name',
    },
    {
      title: '操作',
      key: 'action',
      width: 180,
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EditOutlined />}
            onClick={() => handleEdit(record)}
          >
            编辑
          </Button>
          <Popconfirm
            title="确定要删除此记录吗？"
            onConfirm={() => handleDelete(record.id)}
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
    <div style={{
      minHeight: '100vh',
      background: 'linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%)',
      padding: '40px 20px'
    }}>
      <Card
        title="学生科目管理"
        extra={
          <Button type="default" onClick={() => navigate('/admin/dashboard')}>
            返回管理后台
          </Button>
        }
        style={{
          maxWidth: 1400,
          margin: '0 auto',
          borderRadius: 16,
          boxShadow: '0 10px 25px rgba(0,0,0,0.08)',
        }}
      >
        {/* 搜索区域 */}
        <div style={{ marginBottom: 16 }}>
          <Space wrap>
            <Input
              placeholder="学号"
              value={filters.staffId}
              onChange={(e) => setFilters({ ...filters, staffId: e.target.value })}
              style={{ width: 200 }}
              allowClear
            />
            <Input
              placeholder="科目名称"
              value={filters.subjectName}
              onChange={(e) => setFilters({ ...filters, subjectName: e.target.value })}
              style={{ width: 200 }}
              allowClear
            />
            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              搜索
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              重置
            </Button>
            <Button type="primary" icon={<PlusOutlined />} onClick={handleAdd}>
              添加关联
            </Button>
          </Space>
        </div>

        {/* 表格 */}
        <Table
          columns={columns}
          dataSource={dataSource}
          rowKey="id"
          loading={loading}
          pagination={pagination}
          onChange={handleTableChange}
        />
      </Card>

      {/* 添加/编辑模态框 */}
      <Modal
        title={editingRecord ? '编辑学生科目' : '添加学生科目'}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingRecord(null);
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="staff_id"
            label="学号"
            rules={[{ required: true, message: '请输入学号' }]}
          >
            <Input placeholder="请输入学号" disabled={!!editingRecord} />
          </Form.Item>

          <Form.Item
            name="subject_name"
            label="科目名称"
            rules={[{ required: true, message: '请选择科目' }]}
          >
            <Select
              placeholder="请选择科目"
              showSearch
              filterOption={(input, option) =>
                option.children.toLowerCase().includes(input.toLowerCase())
              }
            >
              {subjects.map(subject => (
                <Select.Option key={subject.ID} value={subject.SubjectName}>
                  {subject.SubjectName}
                </Select.Option>
              ))}
            </Select>
          </Form.Item>

          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
                setEditingRecord(null);
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingRecord ? '确认更新' : '确认添加'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default StudentSubjectsManagement;
