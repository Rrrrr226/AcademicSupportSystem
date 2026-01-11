import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Space, Popconfirm, Typography, message, Tag, Radio } from 'antd';
import { DeleteOutlined, PlusOutlined, EditOutlined, AppstoreOutlined } from '@ant-design/icons';
import { getFastgptAppList, createFastgptApp, updateFastgptApp, deleteFastgptApp } from '../../api';

const { Title } = Typography;

const FastGPTAppsTab = () => {
  const [fastgptApps, setFastgptApps] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [modalVisible, setModalVisible] = useState(false);
  const [editingFastgptApp, setEditingFastgptApp] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchFastgptApps(1, 10);
  }, []);

  const fetchFastgptApps = async (page = 1, pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getFastgptAppList(token, page, pageSize);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setFastgptApps(response.data.data?.apps || []);
        setPagination({
          current: page,
          pageSize: pageSize,
          total: response.data.data?.total || 0,
        });
      }
    } catch (error) {
      console.error('获取应用列表失败:', error);
      message.error(error.response?.data?.message || '获取应用列表失败');
    } finally {
      setLoading(false);
    }
  };

  const handleAddOrUpdateFastgptApp = async (values) => {
    const token = localStorage.getItem('adminToken');
    try {
      let response;
      if (editingFastgptApp) {
        response = await updateFastgptApp({
          id: editingFastgptApp.id,
          appName: values.appName,
          apiKey: values.apiKey,
          description: values.description
        }, token);
      } else {
        response = await createFastgptApp({
          appName: values.appName,
          apiKey: values.apiKey,
          description: values.description
        }, token);
      }

      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success(editingFastgptApp ? '更新成功' : '创建成功');
        setModalVisible(false);
        form.resetFields();
        setEditingFastgptApp(null);
        fetchFastgptApps(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '操作失败');
      }
    } catch (error) {
      console.error('操作失败:', error);
      message.error(error.response?.data?.message || '操作失败');
    }
  };

  const handleDeleteFastgptApp = async (id) => {
    const token = localStorage.getItem('adminToken');
    try {
      const response = await deleteFastgptApp(id, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchFastgptApps(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除应用失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const openAddModal = () => {
    setEditingFastgptApp(null);
    form.resetFields();
    setModalVisible(true);
  };

  const openEditModal = (record) => {
    setEditingFastgptApp(record);
    form.setFieldsValue({
      appName: record.appName,
      apiKey: record.apiKey,
      description: record.description
    });
    setModalVisible(true);
  };

  const handleTableChange = (pagination) => {
    fetchFastgptApps(pagination.current, pagination.pageSize);
  };

  const columns = [
    { title: '学科名称', dataIndex: 'appName', key: 'appName' },
    { 
      title: '密钥', 
      dataIndex: 'apiKey', 
      key: 'apiKey',
      render: (text) => (
        <Tag color="blue" style={{ fontFamily: 'monospace' }}>
          {text ? text.substring(0, 8) + '...' + text.substring(text.length - 4) : '-'}
        </Tag>
      )
    },
    { title: '描述', dataIndex: 'description', key: 'description', ellipsis: true },
    { 
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space>
          <Button type="link" icon={<EditOutlined />} onClick={() => openEditModal(record)}>编辑</Button>
          <Popconfirm title="确定删除吗？" onConfirm={() => handleDeleteFastgptApp(record.id)} okText="确定" cancelText="取消">
            <Button type="link" danger icon={<DeleteOutlined />}>删除</Button>
          </Popconfirm>
        </Space>
      )
    }
  ];

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>学科管理</Title>
        <Button
          type="primary"
          icon={<PlusOutlined />}
          onClick={openAddModal}
        >
          添加学科
        </Button>
      </div>

      <Table
        columns={columns}
        dataSource={fastgptApps}
        rowKey="id"
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
      />

      <Modal
        title={editingFastgptApp ? "编辑学科" : "添加学科"}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingFastgptApp(null);
        }}
        footer={null}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleAddOrUpdateFastgptApp}
        >
          <Form.Item
            name="appName"
            label="学科名称"
            rules={[{ required: true, message: '请输入学科名称' }]}
          >
            <Input placeholder="给学科起个名字" />
          </Form.Item>
          <Form.Item
            name="apiKey"
            label="密钥"
            rules={[{ required: true, message: '请输入密钥' }]}
          >
            <Input.Password placeholder="FastGPT API Key" />
          </Form.Item>
          <Form.Item
            name="description"
            label="描述"
          >
            <Input.TextArea placeholder="学科描述" />
          </Form.Item>
          <Form.Item>
            <Space style={{ width: '100%', justifyContent: 'flex-end' }}>
              <Button onClick={() => {
                setModalVisible(false);
                form.resetFields();
                setEditingFastgptApp(null);
              }}>
                取消
              </Button>
              <Button type="primary" htmlType="submit">
                {editingFastgptApp ? '更新' : '添加'}
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default FastGPTAppsTab;
