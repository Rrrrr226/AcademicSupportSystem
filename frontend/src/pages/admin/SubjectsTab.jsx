import React, { useState, useEffect } from 'react';
import { Table, Button, Modal, Form, Input, Space, Popconfirm, Typography, message } from 'antd';
import { DeleteOutlined, PlusOutlined, EditOutlined } from '@ant-design/icons';
import { getSubjectList, addSubject, deleteSubject, updateSubject } from '../../api';

const { Title } = Typography;

const SubjectsTab = () => {
  const [subjects, setSubjects] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({ current: 1, pageSize: 10, total: 0 });
  const [modalVisible, setModalVisible] = useState(false);
  const [editingSubject, setEditingSubject] = useState(null);
  const [form] = Form.useForm();

  useEffect(() => {
    fetchSubjects(1, 10);
  }, []);

  const fetchSubjects = async (page = 1, pageSize = 10) => {
    const token = localStorage.getItem('adminToken');
    setLoading(true);
    try {
      const response = await getSubjectList(token, page, pageSize); 
      if (response.data?.code === 0 || response.data?.code === 200) {
        setSubjects(response.data.data?.subjects || []);
        setPagination({
            current: response.data.data?.page || page,
            pageSize: response.data.data?.page_size || pageSize,
            total: response.data.data?.total || 0
        });
      }
    } catch (error) {
      console.error('获取学科列表失败:', error);
    } finally {
      setLoading(false);
    }
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
        setModalVisible(false);
        form.resetFields();
        setEditingSubject(null);
        fetchSubjects(pagination.current, pagination.pageSize);
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
        fetchSubjects(pagination.current, pagination.pageSize);
      } else {
        message.error(response.data?.message || '删除失败');
      }
    } catch (error) {
      console.error('删除学科失败:', error);
      message.error(error.response?.data?.message || '删除失败');
    }
  };

  const openAddModal = () => {
    setEditingSubject(null);
    form.resetFields();
    setModalVisible(true);
  };

  const openEditModal = (record) => {
    setEditingSubject(record);
    form.setFieldsValue({
      subjectName: record.subject_name,
      subjectLink: record.subject_link
    });
    setModalVisible(true);
  };

  const handleTableChange = (pagination) => {
    fetchSubjects(pagination.current, pagination.pageSize);
  };

  const columns = [
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
            onClick={() => openEditModal(record)}
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

  return (
    <div>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Title level={4} style={{ margin: 0 }}>学科列表</Title>
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
        dataSource={subjects}
        rowKey="id"
        pagination={pagination}
        loading={loading}
        onChange={handleTableChange}
      />

      <Modal
        title={editingSubject ? "编辑学科" : "添加学科"}
        open={modalVisible}
        onCancel={() => {
          setModalVisible(false);
          form.resetFields();
          setEditingSubject(null);
        }}
        footer={null}
      >
        <Form
          form={form}
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
                setModalVisible(false);
                form.resetFields();
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
    </div>
  );
};

export default SubjectsTab;
