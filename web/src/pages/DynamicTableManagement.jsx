import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { 
  Card, 
  Table, 
  Button, 
  Space, 
  Modal, 
  Form, 
  Input, 
  Select, 
  Switch, 
  message, 
  Popconfirm,
  Tag,
  Tooltip,
  Row,
  Col
} from 'antd';
import messageUtils from '../utils/message';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  SettingOutlined,
  TableOutlined,
  EyeOutlined,
  DatabaseOutlined
} from '@ant-design/icons';
import { dynamicTableApi } from '../api/dynamic';

const { Option } = Select;
const { TextArea } = Input;

const DynamicTableManagement = () => {
  const navigate = useNavigate();
  const [tables, setTables] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingTable, setEditingTable] = useState(null);
  const [form] = Form.useForm();

  // 获取动态表列表
  const fetchTables = async () => {
    setLoading(true);
    try {
      const response = await dynamicTableApi.getTableList();
      if (response.data.success) {
        setTables(response.data.data.list || []);
      } else {
        message.error(response.data.message || '获取表列表失败');
      }
    } catch (error) {
      message.error('获取表列表失败');
      console.error('Error fetching tables:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTables();
  }, []);

  // 打开创建/编辑模态框
  const openModal = (table = null) => {
    setEditingTable(table);
    setModalVisible(true);
    if (table) {
      form.setFieldsValue(table);
    } else {
      form.resetFields();
    }
  };

  // 关闭模态框
  const closeModal = () => {
    setModalVisible(false);
    setEditingTable(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async (values) => {
    try {
      if (editingTable) {
        // 更新表
        const response = await dynamicTableApi.updateTable({
          ...values,
          id: editingTable.id
        });
        if (response.data.success) {
          messageUtils.success('更新表成功');
          fetchTables();
          closeModal();
        } else {
          message.error(response.data.message || '更新表失败');
        }
      } else {
        // 创建表
        const response = await dynamicTableApi.createTable(values);
        if (response.data.success) {
          messageUtils.success('创建表成功');
          fetchTables();
          closeModal();
        } else {
          message.error(response.data.message || '创建表失败');
        }
      }
    } catch (error) {
      message.error(editingTable ? '更新表失败' : '创建表失败');
      console.error('Error submitting form:', error);
    }
  };

  // 删除表
  const handleDelete = async (id) => {
    try {
      const response = await dynamicTableApi.deleteTable(id);
      if (response.data.success) {
        messageUtils.success('删除表成功');
        fetchTables();
      } else {
        message.error(response.data.message || '删除表失败');
      }
    } catch (error) {
      message.error('删除表失败');
      console.error('Error deleting table:', error);
    }
  };

  // 切换表状态
  const handleToggleStatus = async (id) => {
    try {
      const response = await dynamicTableApi.toggleTableStatus(id);
      if (response.data.success) {
        messageUtils.success('状态切换成功');
        fetchTables();
      } else {
        message.error(response.data.message || '状态切换失败');
      }
    } catch (error) {
      message.error('状态切换失败');
      console.error('Error toggling status:', error);
    }
  };

  // 表格列定义
  const columns = [
    {
      title: 'ID',
      dataIndex: 'id',
      key: 'id',
      width: 80,
    },
    {
      title: '表名',
      dataIndex: 'name',
      key: 'name',
      render: (text, record) => (
        <Space>
          <DatabaseOutlined />
          <span>{text}</span>
        </Space>
      ),
    },
    {
      title: '显示名称',
      dataIndex: 'display_name',
      key: 'display_name',
    },
    {
      title: '物理表名',
      dataIndex: 'table_name',
      key: 'table_name',
      render: (text) => <Tag color="blue">{text}</Tag>,
    },
    {
      title: '描述',
      dataIndex: 'description',
      key: 'description',
      ellipsis: {
        showTitle: false,
      },
      render: (text) => (
        <Tooltip placement="topLeft" title={text}>
          {text}
        </Tooltip>
      ),
    },
    {
      title: '状态',
      dataIndex: 'status',
      key: 'status',
      render: (status, record) => (
        <Switch
          checked={status === 1}
          onChange={() => handleToggleStatus(record.id)}
          checkedChildren="启用"
          unCheckedChildren="禁用"
        />
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'created_at',
      key: 'created_at',
      render: (text) => new Date(text).toLocaleString(),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Tooltip title="查看数据">
            <Button
              type="link"
              icon={<EyeOutlined />}
              onClick={() => {
                // 跳转到数据管理页面
                navigate(`/dynamic/data/${record.table_name}`);
              }}
            />
          </Tooltip>
          <Tooltip title="字段管理">
            <Button
              type="link"
              icon={<SettingOutlined />}
              onClick={() => {
                // 跳转到字段管理页面
                navigate(`/dynamic/fields/${record.id}`);
              }}
            />
          </Tooltip>
          <Tooltip title="编辑">
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => openModal(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定要删除这个表吗？"
            description="删除后将无法恢复，请谨慎操作。"
            onConfirm={() => handleDelete(record.id)}
            okText="确定"
            cancelText="取消"
          >
            <Tooltip title="删除">
              <Button
                type="link"
                danger
                icon={<DeleteOutlined />}
              />
            </Tooltip>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Row gutter={[16, 16]}>
        <Col span={24}>
          <Card
            title={
              <Space>
                <TableOutlined />
                <span>动态表管理</span>
              </Space>
            }
            extra={
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => openModal()}
              >
                创建表
              </Button>
            }
          >
            <Table
              columns={columns}
              dataSource={tables}
              rowKey="id"
              loading={loading}
              pagination={{
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total, range) =>
                  `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
              }}
            />
          </Card>
        </Col>
      </Row>

      {/* 创建/编辑表模态框 */}
      <Modal
        title={editingTable ? '编辑表' : '创建表'}
        open={modalVisible}
        onCancel={closeModal}
        footer={null}
        width={600}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Form.Item
            name="name"
            label="表名"
            rules={[
              { required: true, message: '请输入表名' },
              { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '表名只能包含字母、数字和下划线，且必须以字母开头' }
            ]}
          >
            <Input placeholder="请输入表名（英文）" />
          </Form.Item>

          <Form.Item
            name="displayName"
            label="显示名称"
            rules={[{ required: true, message: '请输入显示名称' }]}
          >
            <Input placeholder="请输入显示名称（中文）" />
          </Form.Item>

          <Form.Item
            name="description"
            label="描述"
          >
            <TextArea
              rows={3}
              placeholder="请输入表的描述信息"
            />
          </Form.Item>

          <Form.Item
            name="category"
            label="分类"
          >
            <Select placeholder="请选择表分类">
              <Option value="business">业务数据</Option>
              <Option value="system">系统数据</Option>
              <Option value="config">配置数据</Option>
              <Option value="log">日志数据</Option>
              <Option value="other">其他</Option>
            </Select>
          </Form.Item>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingTable ? '更新' : '创建'}
              </Button>
              <Button onClick={closeModal}>
                取消
              </Button>
            </Space>
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default DynamicTableManagement;