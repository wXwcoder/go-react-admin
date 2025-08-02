import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
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
  InputNumber,
  message, 
  Popconfirm,
  Tag,
  Tooltip,
  Row,
  Col,
  Divider,
  Checkbox
} from 'antd';
import messageUtils from '../utils/message';
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  ArrowUpOutlined,
  ArrowDownOutlined,
  SettingOutlined,
  FieldStringOutlined,
  ArrowLeftOutlined,
  DatabaseOutlined
} from '@ant-design/icons';
import { dynamicFieldApi, dynamicTableApi } from '../api/dynamic';

const { Option } = Select;
const { TextArea } = Input;

const DynamicFieldManagement = () => {
  const { tableId } = useParams();
  const navigate = useNavigate();
  const [fields, setFields] = useState([]);
  const [tableInfo, setTableInfo] = useState(null);
  const [fieldTypes, setFieldTypes] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingField, setEditingField] = useState(null);
  const [form] = Form.useForm();

  // 获取字段类型列表
  const fetchFieldTypes = async () => {
    try {
      const response = await dynamicFieldApi.getFieldTypes();
      if (response.data.success) {
        setFieldTypes(response.data.data || []);
      }
    } catch (error) {
      console.error('Error fetching field types:', error);
    }
  };

  // 获取表信息
  const fetchTableInfo = async () => {
    try {
      const response = await dynamicTableApi.getTableByID(tableId);
      if (response.data.success) {
        setTableInfo(response.data.data);
      } else {
        message.error('获取表信息失败');
        navigate('/dynamic-tables');
      }
    } catch (error) {
      message.error('获取表信息失败');
      navigate('/dynamic-tables');
    }
  };

  // 返回数据管理页面
  const handleBackToDataManagement = () => {
    if (tableInfo) {
      navigate(`/dynamic/data/${tableInfo.table_name}`, { state: { refreshFields: true } });
    } else {
      navigate('/dynamic-tables');
    }
  };

  // 获取字段列表
  const fetchFields = async () => {
    setLoading(true);
    try {
      const response = await dynamicFieldApi.getFieldsByTableID(tableId);
      if (response.data.success) {
        setFields(response.data.data || []);
      } else {
        message.error(response.data.message || '获取字段列表失败');
      }
    } catch (error) {
      message.error('获取字段列表失败');
      console.error('Error fetching fields:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (tableId) {
      fetchTableInfo();
      fetchFields();
      fetchFieldTypes();
    }
  }, [tableId]);

  // 打开创建/编辑模态框
  const openModal = (field = null) => {
    setEditingField(field);
    setModalVisible(true);
    if (field) {
      // 处理选项数据格式
      let optionsText = '';
      if (field.options) {
        if (Array.isArray(field.options)) {
          optionsText = field.options.map(opt => 
            typeof opt === 'object' ? opt.value || opt.label : opt
          ).join('\n');
        } else if (typeof field.options === 'string') {
          optionsText = field.options;
        }
      }
      
      form.setFieldsValue({
        ...field,
        fieldType: field.fieldType || field.field_type,
        displayName: field.displayName || field.display_name,
        isRequired: field.isRequired !== undefined ? field.isRequired : field.required,
        isUnique: field.isUnique !== undefined ? field.isUnique : field.unique,
        isSearchable: field.isSearchable !== undefined ? field.isSearchable : field.searchable,
        isSortable: field.isSortable !== undefined ? field.isSortable : field.sortable,
        options: optionsText,
        validation: field.validation || {}
      });
    } else {
      form.resetFields();
      form.setFieldsValue({ tableId: parseInt(tableId) });
    }
  };

  // 关闭模态框
  const closeModal = () => {
    setModalVisible(false);
    setEditingField(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async (values) => {
    try {
      // 处理选项数据
      let processedOptions = [];
      if (values.options) {
        if (typeof values.options === 'string') {
          processedOptions = values.options.split('\n')
            .filter(opt => opt.trim())
            .map(opt => ({
              label: opt.trim(),
              value: opt.trim()
            }));
        } else if (Array.isArray(values.options)) {
          processedOptions = values.options;
        }
      }

      const processedValues = {
        ...values,
        table_id: parseInt(tableId),
        field_name: values.fieldName,
        display_name: values.displayName,
        field_type: values.fieldType,
        is_required: values.isRequired || false,
        is_unique: values.isUnique || false,
        is_searchable: values.isSearchable || false,
        is_sortable: values.isSortable || false,
        default_value: values.defaultValue || '',
        options: processedOptions,
        validation: values.validation || {}
      };

      // 移除不需要的字段
      delete processedValues.fieldType; // 使用field_type

      if (editingField) {
        // 更新字段
        const response = await dynamicFieldApi.updateField({
          ...processedValues,
          id: editingField.id
        });
        if (response.data.success) {
          messageUtils.success('更新字段成功');
          fetchFields();
          closeModal();
        } else {
          message.error(response.data.message || '更新字段失败');
        }
      } else {
        // 创建字段
        const response = await dynamicFieldApi.createField(processedValues);
        if (response.data.success) {
          messageUtils.success('创建字段成功');
          fetchFields();
          closeModal();
        } else {
          message.error(response.data.message || '创建字段失败');
        }
      }
    } catch (error) {
      message.error(editingField ? '更新字段失败' : '创建字段失败');
      console.error('Error submitting form:', error);
    }
  };

  // 删除字段
  const handleDelete = async (id) => {
    try {
      const response = await dynamicFieldApi.deleteField(id);
      if (response.data.success) {
        messageUtils.success('删除字段成功');
        fetchFields();
      } else {
        message.error(response.data.message || '删除字段失败');
      }
    } catch (error) {
      message.error('删除字段失败');
      console.error('Error deleting field:', error);
    }
  };

  // 切换字段状态
  const handleToggleStatus = async (id) => {
    try {
      const response = await dynamicFieldApi.toggleFieldStatus(id);
      if (response.data.success) {
        messageUtils.success('状态切换成功');
        fetchFields();
      } else {
        message.error(response.data.message || '状态切换失败');
      }
    } catch (error) {
      message.error('状态切换失败');
      console.error('Error toggling status:', error);
    }
  };

  // 更新字段排序
  const handleUpdateOrder = async (fieldId, direction) => {
    const currentIndex = fields.findIndex(f => f.id === fieldId);
    if (currentIndex === -1) return;

    const newIndex = direction === 'up' ? currentIndex - 1 : currentIndex + 1;
    if (newIndex < 0 || newIndex >= fields.length) return;

    const newFields = [...fields];
    [newFields[currentIndex], newFields[newIndex]] = [newFields[newIndex], newFields[currentIndex]];

    // 更新排序值
    const orderData = newFields.map((field, index) => ({
      id: field.id,
      sort_order: index + 1
    }));

    try {
      const response = await dynamicFieldApi.updateFieldOrder(tableId, { fields: orderData });
      if (response.data.success) {
        messageUtils.success('排序更新成功');
        fetchFields();
      } else {
        message.error(response.data.message || '排序更新失败');
      }
    } catch (error) {
      message.error('排序更新失败');
      console.error('Error updating order:', error);
    }
  };

  // 获取字段类型标签颜色
  const getFieldTypeColor = (type) => {
    const colorMap = {
      'string': 'blue',
      'text': 'cyan',
      'int': 'green',
      'float': 'orange',
      'date': 'purple',
      'datetime': 'magenta',
      'select': 'gold',
      'multiselect': 'lime',
      'file': 'red',
      'image': 'pink',
      'boolean': 'geekblue'
    };
    return colorMap[type] || 'default';
  };

  // 获取字段类型显示名称
  const getFieldTypeLabel = (type) => {
    const typeMap = {
      'string': '单行文本',
      'text': '多行文本',
      'int': '整数',
      'float': '小数',
      'date': '日期',
      'datetime': '日期时间',
      'select': '下拉选择',
      'multiselect': '多选下拉',
      'file': '文件上传',
      'image': '图片上传',
      'boolean': '布尔值'
    };
    return typeMap[type] || type;
  };

  // 表格列定义
  const columns = [
    {
      title: '排序',
      dataIndex: 'sort_order',
      key: 'sort_order',
      width: 80,
      render: (_, record, index) => (
        <Space>
          <Button
            type="text"
            size="small"
            icon={<ArrowUpOutlined />}
            disabled={index === 0}
            onClick={() => handleUpdateOrder(record.id, 'up')}
          />
          <Button
            type="text"
            size="small"
            icon={<ArrowDownOutlined />}
            disabled={index === fields.length - 1}
            onClick={() => handleUpdateOrder(record.id, 'down')}
          />
        </Space>
      ),
    },
    {
      title: '字段名',
      dataIndex: 'field_name',
      key: 'field_name',
      render: (text) => (
        <Space>
          <FieldStringOutlined />
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
      title: '字段类型',
      dataIndex: 'field_type',
      key: 'field_type',
      render: (type) => (
        <Tag color={getFieldTypeColor(type)}>{getFieldTypeLabel(type)}</Tag>
      ),
    },
    {
      title: '必填',
      dataIndex: 'is_required',
      key: 'is_required',
      render: (required) => (
        <Tag color={required ? 'red' : 'default'}>
          {required ? '必填' : '可选'}
        </Tag>
      ),
    },
    {
      title: '唯一',
      dataIndex: 'is_unique',
      key: 'is_unique',
      render: (unique) => (
        <Tag color={unique ? 'orange' : 'default'}>
          {unique ? '唯一' : '普通'}
        </Tag>
      ),
    },
    {
      title: '可搜索',
      dataIndex: 'is_searchable',
      key: 'is_searchable',
      render: (searchable) => searchable ? '✓' : '✗',
    },
    {
      title: '可排序',
      dataIndex: 'is_sortable',
      key: 'is_sortable',
      render: (sortable) => sortable ? '✓' : '✗',
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
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Tooltip title="编辑">
            <Button
              type="link"
              icon={<EditOutlined />}
              onClick={() => openModal(record)}
            />
          </Tooltip>
          <Popconfirm
            title="确定要删除这个字段吗？"
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
                <Button
                  type="text"
                  icon={<ArrowLeftOutlined />}
                  onClick={() => navigate('/dynamic-tables')}
                >
                  返回列表
                </Button>
                <Button
                  type="text"
                  icon={<DatabaseOutlined />}
                  onClick={handleBackToDataManagement}
                >
                  返回数据管理
                </Button>
                <Divider type="vertical" />
                <SettingOutlined />
                <span>字段管理</span>
                {tableInfo && (
                  <Tag color="blue">{tableInfo.displayName}</Tag>
                )}
              </Space>
            }
            extra={
              <Button
                type="primary"
                icon={<PlusOutlined />}
                onClick={() => openModal()}
              >
                添加字段
              </Button>
            }
          >
            <Table
              columns={columns}
              dataSource={fields}
              rowKey="id"
              loading={loading}
              pagination={false}
            />
          </Card>
        </Col>
      </Row>

      {/* 创建/编辑字段模态框 */}
      <Modal
        title={editingField ? '编辑字段' : '添加字段'}
        open={modalVisible}
        onCancel={closeModal}
        footer={null}
        width={800}
      >
        <Form
          form={form}
          layout="vertical"
          onFinish={handleSubmit}
        >
          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="fieldName"
                label="字段名"
                rules={[
                  { required: true, message: '请输入字段名' },
                  { pattern: /^[a-zA-Z][a-zA-Z0-9_]*$/, message: '字段名只能包含字母、数字和下划线，且必须以字母开头' }
                ]}
              >
                <Input placeholder="请输入字段名（英文）" />
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="displayName"
                label="显示名称"
                rules={[{ required: true, message: '请输入显示名称' }]}
              >
                <Input placeholder="请输入显示名称（中文）" />
              </Form.Item>
            </Col>
          </Row>

          <Row gutter={16}>
            <Col span={12}>
              <Form.Item
                name="fieldType"
                label="字段类型"
                rules={[{ required: true, message: '请选择字段类型' }]}
              >
                <Select placeholder="请选择字段类型">
                  {fieldTypes.map(type => (
                    <Option key={type.value} value={type.value}>
                      {type.label}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
            </Col>
            <Col span={12}>
              <Form.Item
                name="defaultValue"
                label="默认值"
              >
                <Input placeholder="请输入默认值" />
              </Form.Item>
            </Col>
          </Row>

          <Form.Item
            name="description"
            label="描述"
          >
            <TextArea
              rows={2}
              placeholder="请输入字段描述"
            />
          </Form.Item>

          <Form.Item
            noStyle
            shouldUpdate={(prevValues, currentValues) => 
              prevValues.fieldType !== currentValues.fieldType
            }
          >
            {({ getFieldValue }) => {
              const fieldType = getFieldValue('fieldType');
              if (fieldType === 'select' || fieldType === 'multiselect') {
                return (
                  <Form.Item
                    name="options"
                    label="选项列表"
                    rules={[{ required: true, message: '请输入选项列表' }]}
                  >
                    <TextArea
                      rows={4}
                      placeholder="请输入选项，每行一个选项"
                    />
                  </Form.Item>
                );
              }
              return null;
            }}
          </Form.Item>

          <Row gutter={16}>
            <Col span={6}>
              <Form.Item
                name="isRequired"
                valuePropName="checked"
              >
                <Checkbox>必填</Checkbox>
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item
                name="isUnique"
                valuePropName="checked"
              >
                <Checkbox>唯一</Checkbox>
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item
                name="isSearchable"
                valuePropName="checked"
              >
                <Checkbox>可搜索</Checkbox>
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item
                name="isSortable"
                valuePropName="checked"
              >
                <Checkbox>可排序</Checkbox>
              </Form.Item>
            </Col>
          </Row>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingField ? '更新' : '添加'}
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

export default DynamicFieldManagement;