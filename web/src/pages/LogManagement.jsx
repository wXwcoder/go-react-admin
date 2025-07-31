import React, { useState, useEffect } from 'react';
import { Table, Button, Input, Select, DatePicker, Space, Card, Tag, Modal, Descriptions } from 'antd';
import { SearchOutlined, ReloadOutlined, EyeOutlined } from '@ant-design/icons';
import dayjs from 'dayjs';
import api from '../api';
import '../assets/styles/management.css';

const { RangePicker } = DatePicker;
const { Option } = Select;

const LogManagement = () => {
  const [logs, setLogs] = useState([]);
  const [loading, setLoading] = useState(false);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchParams, setSearchParams] = useState({
    username: '',
    method: '',
    statusCode: '',
    dateRange: [],
  });
  const [selectedLog, setSelectedLog] = useState(null);
  const [modalVisible, setModalVisible] = useState(false);

  // 获取日志列表
  const fetchLogs = async (params = {}) => {
    setLoading(true);
    try {
      const response = await api.get('/log/list', {
        params: {
          page: pagination.current,
          pageSize: pagination.pageSize,
          ...params,
        },
      });
      
      if (response.data && response.data.logs) {
        setLogs(response.data.logs);
        setPagination({
          ...pagination,
          total: response.data.total || response.data.logs.length,
        });
      }
    } catch (error) {
      console.error('获取日志列表失败:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  // 处理搜索
  const handleSearch = () => {
    const params = {};
    if (searchParams.username) params.username = searchParams.username;
    if (searchParams.method) params.method = searchParams.method;
    if (searchParams.statusCode) params.statusCode = searchParams.statusCode;
    if (searchParams.dateRange && searchParams.dateRange.length === 2) {
      params.startDate = searchParams.dateRange[0].format('YYYY-MM-DD');
      params.endDate = searchParams.dateRange[1].format('YYYY-MM-DD');
    }
    
    setPagination({ ...pagination, current: 1 });
    fetchLogs(params);
  };

  // 重置搜索
  const handleReset = () => {
    setSearchParams({
      username: '',
      method: '',
      statusCode: '',
      dateRange: [],
    });
    setPagination({ ...pagination, current: 1 });
    fetchLogs();
  };

  // 处理表格分页
  const handleTableChange = (newPagination) => {
    setPagination(newPagination);
    fetchLogs({
      page: newPagination.current,
      pageSize: newPagination.pageSize,
    });
  };

  // 查看日志详情
  const handleViewDetail = (record) => {
    setSelectedLog(record);
    setModalVisible(true);
  };

  // 获取状态码颜色
  const getStatusColor = (status) => {
    if (status >= 200 && status < 300) return 'success';
    if (status >= 300 && status < 400) return 'warning';
    if (status >= 400) return 'error';
    return 'default';
  };

  // 获取HTTP方法颜色
  const getMethodColor = (method) => {
    const colors = {
      GET: 'green',
      POST: 'blue',
      PUT: 'orange',
      DELETE: 'red',
      PATCH: 'purple',
    };
    return colors[method] || 'default';
  };

  const columns = [
    {
      title: '用户',
      dataIndex: 'username',
      key: 'username',
      width: 120,
    },
    {
      title: 'IP地址',
      dataIndex: 'ip',
      key: 'ip',
      width: 120,
    },
    {
      title: '请求方法',
      dataIndex: 'method',
      key: 'method',
      width: 80,
      render: (method) => (
        <Tag color={getMethodColor(method)}>{method}</Tag>
      ),
    },
    {
      title: '请求路径',
      dataIndex: 'path',
      key: 'path',
      width: 200,
      ellipsis: true,
    },
    {
      title: '状态码',
      dataIndex: 'status_code',
      key: 'status_code',
      width: 80,
      render: (status) => (
        <Tag color={getStatusColor(status)}>{status}</Tag>
      ),
    },
    {
      title: '响应时间',
      dataIndex: 'response_time',
      key: 'response_time',
      width: 100,
      render: (time) => `${time}ms`,
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      width: 180,
      render: (date) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'),
    },
    {
      title: '操作',
      key: 'action',
      width: 80,
      fixed: 'right',
      render: (_, record) => (
        <Space>
          <Button
            type="link"
            icon={<EyeOutlined />}
            onClick={() => handleViewDetail(record)}
          >
            详情
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="management-container">
      <Card className="management-card">
        <div className="page-header">
          <h2>日志管理</h2>
          <p>查看和管理系统操作日志</p>
        </div>

        {/* 搜索区域 */}
        <div className="search-section">
          <Space wrap>
            <Input
              placeholder="用户名"
              value={searchParams.username}
              onChange={(e) => setSearchParams({ ...searchParams, username: e.target.value })}
              style={{ width: 150 }}
            />
            <Select
              placeholder="请求方法"
              value={searchParams.method}
              onChange={(value) => setSearchParams({ ...searchParams, method: value })}
              style={{ width: 100 }}
              allowClear
            >
              <Option value="GET">GET</Option>
              <Option value="POST">POST</Option>
              <Option value="PUT">PUT</Option>
              <Option value="DELETE">DELETE</Option>
              <Option value="PATCH">PATCH</Option>
            </Select>
            <Select
              placeholder="状态码"
              value={searchParams.statusCode}
              onChange={(value) => setSearchParams({ ...searchParams, statusCode: value })}
              style={{ width: 100 }}
              allowClear
            >
              <Option value="200">200</Option>
              <Option value="201">201</Option>
              <Option value="400">400</Option>
              <Option value="401">401</Option>
              <Option value="403">403</Option>
              <Option value="404">404</Option>
              <Option value="500">500</Option>
            </Select>
            <RangePicker
              value={searchParams.dateRange}
              onChange={(dates) => setSearchParams({ ...searchParams, dateRange: dates })}
              style={{ width: 240 }}
            />
            <Button type="primary" icon={<SearchOutlined />} onClick={handleSearch}>
              搜索
            </Button>
            <Button icon={<ReloadOutlined />} onClick={handleReset}>
              重置
            </Button>
          </Space>
        </div>

        {/* 表格区域 */}
        <div className="table-container">
          <Table
            columns={columns}
            dataSource={logs}
            rowKey="id"
            loading={loading}
            pagination={pagination}
            onChange={handleTableChange}
            scroll={{ x: 1200 }}
            size="middle"
            className="management-table"
          />
        </div>
      </Card>

      {/* 日志详情模态框 */}
      <Modal
        title="日志详情"
        open={modalVisible}
        onCancel={() => setModalVisible(false)}
        footer={null}
        width={600}
      >
        {selectedLog && (
          <Descriptions column={1} bordered>
            <Descriptions.Item label="用户">
              {selectedLog.username || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="IP地址">
              {selectedLog.ip || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="请求方法">
              <Tag color={getMethodColor(selectedLog.method)}>
                {selectedLog.method}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="请求路径">
              {selectedLog.path}
            </Descriptions.Item>
            <Descriptions.Item label="状态码">
              <Tag color={getStatusColor(selectedLog.status_code)}>
                {selectedLog.status_code}
              </Tag>
            </Descriptions.Item>
            <Descriptions.Item label="响应时间">
              {selectedLog.response_time}ms
            </Descriptions.Item>
            <Descriptions.Item label="User Agent">
              {selectedLog.user_agent || '-'}
            </Descriptions.Item>
            <Descriptions.Item label="创建时间">
              {dayjs(selectedLog.created_at).format('YYYY-MM-DD HH:mm:ss')}
            </Descriptions.Item>
          </Descriptions>
        )}
      </Modal>
    </div>
  );
};

export default LogManagement;