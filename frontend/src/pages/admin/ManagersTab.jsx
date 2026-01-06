import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Space, Popconfirm, Typography, message } from 'antd';
import { DeleteOutlined, PlusOutlined } from '@ant-design/icons';
import { getManagerList, addManager, deleteManager } from '../../api';

const { Title } = Typography;

const ManagersTab = ({ currentUser }) => {
  const [managers, setManagers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [addModalVisible, setAddModalVisible] = useState(false);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchManagers(1, 10);
  }, []);

  const fetchManagers = async (page = 1, pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getManagerList(token, page, pageSize);
      if (response.data?.code === 0 || response.data?.code === 200) {
        setManagers(response.data.data?.managers || []);
        setPagination({
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

  const handleAddManager = async (values) => {
    const token = localStorage.getItem('adminToken');
    try {
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
      const response = await deleteManager(staffId, token);
      if (response.data?.code === 0 || response.data?.code === 200) {
        message.success('删除成功');
        fetchManagers(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除管理员失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const handleTableChange = (pagination) => {
    fetchManagers(pagination.current, pagination.pageSize);
  };

  const columns = [
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

  return (
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
        columns={columns}
        dataSource={managers}
        rowKey="id"
        loading={loading}
        pagination={pagination}
        onChange={handleTableChange}
      />

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
    </div>
  );
};

export default ManagersTab;
