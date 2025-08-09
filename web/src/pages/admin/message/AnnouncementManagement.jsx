import React, { useState, useEffect } from 'react';
import { Card, Button, Table, Space, Tag, Modal, Form, Input, Select, DatePicker, message, Popconfirm, Row, Col, Statistic } from 'antd';
import { PlusOutlined, EditOutlined, DeleteOutlined, EyeOutlined, SendOutlined, StopOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import { announcementApi } from '../../../api/admin';

const { Option } = Select;
const { TextArea } = Input;
const { RangePicker } = DatePicker;

const AnnouncementManagement = () => {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [data, setData] = useState([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [modalVisible, setModalVisible] = useState(false);
  const [modalType, setModalType] = useState('create');
  const [editingRecord, setEditingRecord] = useState(null);
  const [stats, setStats] = useState({
    total: 0,
    published: 0,
    draft: 0,
  });

  // 获取公告列表
  const fetchData = async (params = {}) => {
    setLoading(true);
    try {
      const response = await announcementApi.getAnnouncementList({
        page: pagination.current,
        page_size: pagination.pageSize,
        ...params,
      });
      
      setData(response.data.list);
      setPagination({
        ...pagination,
        total: response.data.total,
      });
    } catch (error) {
      message.error('获取公告列表失败');
    } finally {
      setLoading(false);
    }
  };

  // 获取统计信息
  const fetchStats = async () => {
    try {
      const response = await announcementApi.getAnnouncementStats();
      setStats(response.data);
    } catch (error) {
      console.error('获取统计信息失败:', error);
    }
  };

  useEffect(() => {
    fetchData();
    fetchStats();
  }, [pagination.current, pagination.pageSize]);

  const handleTableChange = (newPagination) => {
    setPagination({
      ...pagination,
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    });
  };

  const handleSearch = (values) => {
    setPagination({ ...pagination, current: 1 });
    fetchData(values);
  };

  const handleCreate = () => {
    setModalType('create');
    setEditingRecord(null);
    form.resetFields();
    setModalVisible(true);
  };

  const handleEdit = (record) => {
    setModalType('edit');
    setEditingRecord(record);
    form.setFieldsValue({
      ...record,
      expire_time: record.expire_time ? dayjs(record.expire_time) : null,
    });
    setModalVisible(true);
  };

  const handleView = (record) => {
    Modal.info({
      title: '公告详情',
      width: 600,
      content: (
        <div style={{ marginTop: 16 }}>
          <p><strong>标题：</strong>{record.title}</p>
          <p><strong>类型：</strong>
            <Tag color={getTypeColor(record.type)}>{getTypeText(record.type)}</Tag>
          </p>
          <p><strong>状态：</strong>
            <Tag color={getStatusColor(record.status)}>{getStatusText(record.status)}</Tag>
          </p>
          <p><strong>优先级：</strong>
            <Tag color={getPriorityColor(record.priority)}>{getPriorityText(record.priority)}</Tag>
          </p>
          <p><strong>内容：</strong></p>
          <div dangerouslySetInnerHTML={{ __html: record.content }} />
          {record.expire_time && (
            <p><strong>过期时间：</strong>{dayjs(record.expire_time).format('YYYY-MM-DD HH:mm:ss')}</p>
          )}
          <p><strong>阅读统计：</strong>已读 {record.read_count} 次</p>
          <p><strong>创建时间：</strong>{dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss')}</p>
        </div>
      ),
    });
  };

  const handleSubmit = async (values) => {
    try {
      if (modalType === 'create') {
        await announcementApi.createAnnouncement({
          ...values,
          expire_time: values.expire_time ? values.expire_time.toISOString() : null,
        });
        message.success('创建成功');
      } else {
        await announcementApi.updateAnnouncement(editingRecord.id, {
          ...values,
          expire_time: values.expire_time ? values.expire_time.toISOString() : null,
        });
        message.success('更新成功');
      }
      
      setModalVisible(false);
      fetchData();
      fetchStats();
    } catch (error) {
      message.error(modalType === 'create' ? '创建失败' : '更新失败');
    }
  };

  const handleDelete = async (id) => {
    try {
      await announcementApi.deleteAnnouncement(id);
      message.success('删除成功');
      fetchData();
      fetchStats();
    } catch (error) {
      message.error('删除失败');
    }
  };

  const handlePublish = async (id) => {
    try {
      await announcementApi.publishAnnouncement(id);
      message.success('发布成功');
      fetchData();
      fetchStats();
    } catch (error) {
      message.error('发布失败');
    }
  };

  const handleRevoke = async (id) => {
    try {
      await announcementApi.revokeAnnouncement(id);
      message.success('撤回成功');
      fetchData();
      fetchStats();
    } catch (error) {
      message.error('撤回失败');
    }
  };

  const getTypeColor = (type) => {
    const colors = {
      system: 'blue',
      notice: 'green',
      maintenance: 'orange',
      update: 'purple',
    };
    return colors[type] || 'default';
  };

  const getTypeText = (type) => {
    const texts = {
      system: '系统公告',
      notice: '通知公告',
      maintenance: '维护公告',
      update: '更新公告',
    };
    return texts[type] || type;
  };

  const getStatusColor = (status) => {
    const colors = {
      draft: 'default',
      published: 'success',
      revoked: 'error',
    };
    return colors[status] || 'default';
  };

  const getStatusText = (status) => {
    const texts = {
      draft: '草稿',
      published: '已发布',
      revoked: '已撤回',
    };
    return texts[status] || status;
  };

  const getPriorityColor = (priority) => {
    const colors = {
      low: 'success',
      medium: 'warning',
      high: 'error',
    };
    return colors[priority] || 'default';
  };

  const getPriorityText = (priority) => {
    const texts = {
      low: '低',
      medium: '中',
      high: '高',
    };
    return texts[priority] || priority;
  };

  const columns = [
    {
      title: '标题',
      dataIndex: 'title',
      key: 'title',
      ellipsis: true,
    },
    {
      title: '类型',
      dataIndex: 'type',
      key: 'type',
      render: (type) => (
        <Tag color={getTypeColor(type)}>{getTypeText(type)}</Tag>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status) => (
        <Tag color={getStatusColor(status)}>{getStatusText(status)}</Tag>
      ),
    },
    {
      title: '优先级',
      dataIndex: 'priority',
      key: 'priority',
      render: (priority) => (
        <Tag color={getPriorityColor(priority)}>{getPriorityText(priority)}</Tag>
      ),
    },
    {
      title: '阅读数',
      dataIndex: 'read_count',
      key: 'read_count',
      render: (count) => <span>{count} 次</span>,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleView(record)}
          >
            查看
          </Button>
          {record.status === 'draft' && (
            <>
              <Button
                type="link"
                icon={<EditOutlined />}
                onClick={() => handleEdit(record)}
              >
                编辑
              </Button>
              <Button
                type="link"
                icon={<SendOutlined />}
                onClick={() => handlePublish(record.id)}
              >
                发布
              </Button>
            </>
          )}
          {record.status === 'published' && (
            <Button
              type="link"
              icon={<StopOutlined />}
              onClick={() => handleRevoke(record.id)}
            >
              撤回
            </Button>
          )}
          {record.status !== 'published' && (
            <Popconfirm
              title="确定要删除这条公告吗？"
              onConfirm={() => handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="link" danger icon={<DeleteOutlined />}>
                删除
              </Button>
            </Popconfirm>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Row gutter={16} style={{ marginBottom: 16 }}>
        <Col span={6}>
          <Statistic title="总公告数" value={stats.total} />
        </Col>
        <Col span={6}>
          <Statistic title="已发布" value={stats.published} />
        </Col>
        <Col span={6}>
          <Statistic title="草稿" value={stats.draft} />
        </Col>
      </Row>

      <Card
        title="公告管理"
        extra={
          <Button type="primary" icon={<PlusOutlined />} onClick={handleCreate}>
            创建公告
          </Button>
        }
      >
        <Table
          columns={columns}
          dataSource={data}
          rowKey="id"
          loading={loading}
          pagination={pagination}
          onChange={handleTableChange}
        />
      </Card>

      <Modal
        title={modalType === 'create' ? '创建公告' : '编辑公告'}
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="title"
            label="标题"
            rules={[{ required: true, message: '请输入公告标题' }]}
          >
            <Input maxLength={255} placeholder="请输入公告标题" />
          </Form.Item>

          <Form.Item
            name="content"
            label="内容"
            rules={[{ required: true, message: '请输入公告内容' }]}
          >
            <TextArea rows={4} placeholder="请输入公告内容" />
          </Form.Item>

          <Form.Item
            name="type"
            label="类型"
            rules={[{ required: true, message: '请选择公告类型' }]}
          >
            <Select placeholder="请选择公告类型">
              <Option value="system">系统公告</Option>
              <Option value="notice">通知公告</Option>
              <Option value="maintenance">维护公告</Option>
              <Option value="update">更新公告</Option>
            </Select>
          </Form.Item>

          <Form.Item
            name="priority"
            label="优先级"
            initialValue="medium"
          >
            <Select>
              <Option value="low">低</Option>
              <Option value="medium">中</Option>
              <Option value="high">高</Option>
            </Select>
          </Form.Item>

          <Form.Item name="expire_time" label="过期时间">
            <DatePicker
              showTime
              placeholder="选择过期时间"
              style={{ width: '100%' }}
            />
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                保存
              </Button>
              <Button onClick={() => setModalVisible(false)}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default AnnouncementManagement;