import React, { useState, useEffect } from 'react';
import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { 
  Card, 
  Table, 
  Button, 
  Space, 
  Modal, 
  Form, 
  Input, 
  Select, 
  DatePicker, 
  InputNumber,
  Switch,
  Upload,
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
  SearchOutlined,
  ReloadOutlined,
  ExportOutlined,
  ImportOutlined,
  ArrowLeftOutlined,
  DatabaseOutlined,
  UploadOutlined
} from '@ant-design/icons';
import { dynamicDataApi, dynamicFieldApi, dynamicTableApi } from '../api/dynamic';
import dayjs from 'dayjs';

const { Option } = Select;
const { TextArea } = Input;
const { RangePicker } = DatePicker;

const DynamicDataManagement = () => {
  const { tableName } = useParams();
  const navigate = useNavigate();
  const location = useLocation();
  const [data, setData] = useState([]);
  const [fields, setFields] = useState([]);
  const [tableInfo, setTableInfo] = useState(null);
  const [loading, setLoading] = useState(false);
  const [modalVisible, setModalVisible] = useState(false);
  const [editingData, setEditingData] = useState(null);
  const [searchVisible, setSearchVisible] = useState(false);
  const [selectedRowKeys, setSelectedRowKeys] = useState([]);
  const [pagination, setPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });
  const [searchParams, setSearchParams] = useState({});
  const [form] = Form.useForm();
  const [searchForm] = Form.useForm();

  // 获取表信息
  const fetchTableInfo = async () => {
    if (!tableName) {
      message.error('表名参数缺失');
      navigate('/dynamic-tables');
      return null;
    }
    
    try {
      console.log('Fetching table info for:', tableName);
      const response = await dynamicTableApi.getTableSchema(tableName);
      console.log('Table info response:', response);
      
      if (response.data && response.data.data && response.data.data.table_info) {
        setTableInfo(response.data.data.table_info);
        return response.data.data.table_info;
      } else {
        message.error(`表 "${tableName}" 不存在`);
        navigate('/dynamic-tables');
        return null;
      }
    } catch (error) {
      message.error(`获取表信息失败: ${error.message || '网络错误'}`);
      console.error('Error fetching table info:', error);
      navigate('/dynamic-tables');
      return null;
    }
  };

  // 获取字段列表
  const fetchFields = async (tableId) => {
    try {
      const response = await dynamicFieldApi.getFieldsByTableID(tableId);
      if (response.data.success) {
        console.log('All table fields:', response.data.data);
        const activeFields = response.data.data.filter(field => field.status === 1);
        console.log('Active fields:', activeFields);
        setFields(activeFields);
        return activeFields;
      }
    } catch (error) {
      console.error('Error fetching fields:', error);
    }
    return [];
  };

  // 获取数据列表
  const fetchData = async (page = 1, pageSize = 10, search = {}) => {
    if (!tableName) {
      message.error('表名参数缺失');
      return;
    }
    
    setLoading(true);
    try {
      const params = {
        page,
        pageSize,
        ...search,
      };
      
      console.log('Fetching data with params:', { tableName, params });
      const response = await dynamicDataApi.getDataList(tableName, params);
      console.log('API Response:', response);
      
      if (response.data.success) {
        // 确保数据结构正确
        const listData = response.data.data?.list || response.data.data || [];
        console.log('Processed data list:', listData);
        
        // 确保数据有正确的字段结构
        const processedData = Array.isArray(listData) ? listData : [listData];
        setData(processedData);
        setPagination({
          current: page,
          pageSize,
          total: response.data.data?.total || processedData.length || 0,
        });
      } else {
        message.error(response.data.message || '获取数据失败');
        setData([]);
        setPagination({
          current: page,
          pageSize,
          total: 0,
        });
      }
    } catch (error) {
      message.error(`获取数据失败: ${error.message || '网络错误'}`);
      console.error('Error fetching data:', error);
      setData([]);
      setPagination({
        current: page,
        pageSize,
        total: 0,
      });
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    const initData = async () => {
      const table = await fetchTableInfo();
      if (table) {
        const fieldList = await fetchFields(table.id);
        if (fieldList.length > 0) {
          fetchData();
        } else {
          message.warning('当前表没有可用字段');
        }
      }
    };
    
    if (tableName) {
      initData();
    } else {
      message.error('表名参数缺失');
      navigate('/dynamic-tables');
    }
  }, [tableName]);

  // 监听路由变化，当从字段管理页面返回时强制刷新字段数据
  useEffect(() => {
    if (location.state?.refreshFields && tableInfo) {
      console.log('检测到字段变更，强制刷新字段数据...');
      fetchFields(tableInfo.id);
      // 清除状态，避免重复刷新
      window.history.replaceState({}, document.title);
    }
  }, [location, tableInfo]);

  // 监听路由变化，当从字段管理页面返回时刷新字段数据
  useEffect(() => {
    const handleFocus = () => {
      console.log('页面重新获得焦点，刷新字段数据...');
      if (tableName && tableInfo) {
        fetchFields(tableInfo.id);
      }
    };

    window.addEventListener('focus', handleFocus);
    return () => {
      window.removeEventListener('focus', handleFocus);
    };
  }, [tableName, tableInfo]);

  // 监听页面可见性变化，当页面重新激活时刷新字段数据
  useEffect(() => {
    const handleVisibilityChange = async () => {
      if (!document.hidden && tableName) {
        console.log('页面重新激活，刷新字段数据...');
        const table = await fetchTableInfo();
        if (table) {
          await fetchFields(table.id);
          fetchData(pagination.current, pagination.pageSize, searchParams);
        }
      }
    };

    document.addEventListener('visibilitychange', handleVisibilityChange);
    
    return () => {
      document.removeEventListener('visibilitychange', handleVisibilityChange);
    };
  }, [tableName, pagination.current, pagination.pageSize, searchParams]);

  // 打开创建/编辑模态框
  const openModal = async (record = null) => {
    setEditingData(record);
    
    let currentFields = fields;
    
    // 重新获取字段数据，确保显示最新的字段
    if (tableInfo) {
      try {
        const response = await dynamicFieldApi.getFieldsByTableID(tableInfo.id);
        if (response.data.success) {
          const activeFields = response.data.data.filter(field => field.status === 1);
          console.log('Updated active fields:', activeFields);
          setFields(activeFields);
          currentFields = activeFields; // 使用最新的字段数据
        }
      } catch (error) {
        console.error('Error refreshing fields:', error);
      }
    }
    
    setModalVisible(true);
    
    // 使用setTimeout确保状态更新后再设置表单值
    setTimeout(() => {
      if (record) {
        // 处理日期字段
        const formData = { ...record };
        currentFields.forEach(field => {
          if ((field.field_type === 'date' || field.field_type === 'datetime') && formData[field.field_name]) {
            formData[field.field_name] = dayjs(formData[field.field_name]);
          }
        });
        form.setFieldsValue(formData);
      } else {
        form.resetFields();
      }
    }, 100);
  };

  // 关闭模态框
  const closeModal = () => {
    setModalVisible(false);
    setEditingData(null);
    form.resetFields();
  };

  // 提交表单
  const handleSubmit = async (values) => {
    try {
      // 处理日期字段和数字字段
      const processedValues = { ...values };
      fields.forEach(field => {
        if ((field.field_type === 'date' || field.field_type === 'datetime') && processedValues[field.field_name]) {
          processedValues[field.field_name] = processedValues[field.field_name].format(
            field.field_type === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss'
          );
        } else if ((field.field_type === 'int' || field.field_type === 'select') && processedValues[field.field_name] !== undefined && processedValues[field.field_name] !== null && processedValues[field.field_name] !== '') {
          // 将字符串转换为整数，处理status等整数字段
          let intValue;
          if (field.field_name === 'status') {
            // status字段特殊处理：active->1, inactive->2
            const statusMap = { 'active': 1, 'inactive': 2, '启用': 1, '禁用': 2 };
            const value = processedValues[field.field_name];
            if (typeof value === 'string') {
              intValue = statusMap[value.toLowerCase()] || parseInt(value);
            } else {
              intValue = parseInt(value);
            }
          } else {
            intValue = parseInt(processedValues[field.field_name]);
          }
          if (!isNaN(intValue)) {
            processedValues[field.field_name] = intValue;
          }
        } else if (field.field_type === 'float' && processedValues[field.field_name] !== undefined && processedValues[field.field_name] !== null && processedValues[field.field_name] !== '') {
          // 将字符串转换为浮点数
          const floatValue = parseFloat(processedValues[field.field_name]);
          if (!isNaN(floatValue)) {
            processedValues[field.field_name] = floatValue;
          }
        }
      });

      if (editingData) {
        // 更新数据
        const response = await dynamicDataApi.updateData(tableName, editingData.id, processedValues);
        if (response.data.success) {
          messageUtils.success('更新数据成功');
          fetchData(pagination.current, pagination.pageSize, searchParams);
          closeModal();
        } else {
          message.error(response.data.message || '更新数据失败');
        }
      } else {
        // 创建数据
        const response = await dynamicDataApi.createData(tableName, processedValues);
        if (response.data.success) {
          messageUtils.success('创建数据成功');
          fetchData(pagination.current, pagination.pageSize, searchParams);
          closeModal();
        } else {
          message.error(response.data.message || '创建数据失败');
        }
      }
    } catch (error) {
      message.error(editingData ? '更新数据失败' : '创建数据失败');
      console.error('Error submitting form:', error);
    }
  };

  // 删除数据
  const handleDelete = async (id) => {
    try {
      const response = await dynamicDataApi.deleteData(tableName, id);
      if (response.data.success) {
        messageUtils.success('删除数据成功');
        fetchData(pagination.current, pagination.pageSize, searchParams);
      } else {
        message.error(response.data.message || '删除数据失败');
      }
    } catch (error) {
      message.error('删除数据失败');
      console.error('Error deleting data:', error);
    }
  };

  // 批量删除
  const handleBatchDelete = async () => {
    if (selectedRowKeys.length === 0) {
      message.warning('请选择要删除的数据');
      return;
    }

    try {
      const response = await dynamicDataApi.batchDeleteData(tableName, selectedRowKeys);
      if (response.data.success) {
        messageUtils.success('批量删除成功');
        setSelectedRowKeys([]);
        fetchData(pagination.current, pagination.pageSize, searchParams);
      } else {
        message.error(response.data.message || '批量删除失败');
      }
    } catch (error) {
      message.error('批量删除失败');
      console.error('Error batch deleting:', error);
    }
  };

  // 搜索
  const handleSearch = (values) => {
    const searchData = {};
    Object.keys(values).forEach(key => {
      if (values[key] !== undefined && values[key] !== null && values[key] !== '') {
        searchData[key] = values[key];
      }
    });
    setSearchParams(searchData);
    fetchData(1, pagination.pageSize, searchData);
  };

  // 重置搜索
  const handleResetSearch = () => {
    searchForm.resetFields();
    setSearchParams({});
    fetchData(1, pagination.pageSize, {});
  };

  // 渲染表单项
  const renderFormItem = (field) => {
    const { field_name, display_name, field_type, is_required, options, default_value } = field;
    
    const rules = [];
    if (is_required) {
      rules.push({ required: true, message: `请输入${display_name}` });
    }

    switch (field_type) {
      case 'string':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value}
          >
            <Input placeholder={`请输入${display_name}`} />
          </Form.Item>
        );
      
      case 'text':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value}
          >
            <TextArea rows={3} placeholder={`请输入${display_name}`} />
          </Form.Item>
        );
      
      case 'int':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value ? parseInt(default_value) : undefined}
          >
            <InputNumber 
              style={{ width: '100%' }} 
              placeholder={`请输入${display_name}`} 
            />
          </Form.Item>
        );
      
      case 'float':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value ? parseFloat(default_value) : undefined}
          >
            <InputNumber 
              style={{ width: '100%' }} 
              step={0.01}
              placeholder={`请输入${display_name}`} 
            />
          </Form.Item>
        );
      
      case 'date':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
          >
            <DatePicker 
              style={{ width: '100%' }} 
              placeholder={`请选择${display_name}`} 
            />
          </Form.Item>
        );
      
      case 'datetime':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
          >
            <DatePicker 
              showTime
              style={{ width: '100%' }} 
              placeholder={`请选择${display_name}`} 
            />
          </Form.Item>
        );
      
      case 'select':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value}
          >
            <Select placeholder={`请选择${display_name}`}>
            {options && options.map((option, index) => {
              const optionValue = typeof option === 'object' ? option.value : option;
              const optionLabel = typeof option === 'object' ? option.label : option;
              return <Option key={`${optionValue}_${index}`} value={optionValue}>{optionLabel}</Option>;
            })}
          </Select>
          </Form.Item>
        );
      
      case 'multiselect':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
          >
            <Select 
              mode="multiple" 
              placeholder={`请选择${display_name}`}
            >
              {options && options.map((option, index) => {
                const optionValue = typeof option === 'object' ? option.value : option;
                const optionLabel = typeof option === 'object' ? option.label : option;
                return <Option key={`${optionValue}_${index}`} value={optionValue}>{optionLabel}</Option>;
              })}
            </Select>
          </Form.Item>
        );
      
      case 'boolean':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            valuePropName="checked"
            initialValue={default_value === 'true'}
          >
            <Switch />
          </Form.Item>
        );
      
      case 'file':
      case 'image':
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
          >
            <Upload>
              <Button icon={<UploadOutlined />}>上传文件</Button>
            </Upload>
          </Form.Item>
        );
      
      default:
        return (
          <Form.Item
            key={field_name}
            name={field_name}
            label={display_name}
            rules={rules}
            initialValue={default_value}
          >
            <Input placeholder={`请输入${display_name}`} />
          </Form.Item>
        );
    }
  };

  // 生成表格列
  const generateColumns = () => {
    console.log('Generating columns with fields:', fields);
    console.log('Current data:', data);
    
    const columns = [
      {
        title: 'ID',
        dataIndex: 'id',
        key: 'id',
        width: 80,
        fixed: 'left',
      },
    ];

    fields.forEach(field => {
      columns.push({
        title: field.display_name,
        dataIndex: field.field_name,
        key: field.field_name,
        ellipsis: true,
        render: (value) => {
          if (field.field_type === 'boolean') {
            return value ? '是' : '否';
          }
          if (field.field_type === 'date' || field.field_type === 'datetime') {
            return value ? dayjs(value).format(
              field.field_type === 'date' ? 'YYYY-MM-DD' : 'YYYY-MM-DD HH:mm:ss'
            ) : '';
          }
          if (field.field_type === 'select' && value) {
            // 处理select类型，可能是对象或字符串
            if (typeof value === 'object' && value !== null) {
              return value.label || value.value || String(value);
            }
            return String(value);
          }
          if (field.field_type === 'multiselect' && Array.isArray(value)) {
            return value.map(v => {
              const displayValue = typeof v === 'object' ? (v.label || v.value) : v;
              return <Tag key={displayValue}>{displayValue}</Tag>;
            });
          }
          return value !== null && value !== undefined ? String(value) : '';
        },
      });
    });

    columns.push({
      title: '操作',
      key: 'action',
      fixed: 'right',
      width: 120,
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
            title="确定要删除这条数据吗？"
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
    });

    return columns;
  };

  const rowSelection = {
    selectedRowKeys,
    onChange: setSelectedRowKeys,
  };

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
                  返回
                </Button>
                <Divider type="vertical" />
                <DatabaseOutlined />
                <span>数据管理</span>
                {tableInfo && (
                  <Tag color="blue">{tableInfo.display_name}</Tag>
                )}
              </Space>
            }
            extra={
              <Space>
                <Button
                  icon={<SearchOutlined />}
                  onClick={() => setSearchVisible(!searchVisible)}
                >
                  {searchVisible ? '隐藏搜索' : '显示搜索'}
                </Button>
                <Button
                  icon={<ReloadOutlined />}
                  onClick={async () => {
                    // 重新获取表信息和字段数据
                    const table = await fetchTableInfo();
                    if (table) {
                      await fetchFields(table.id);
                    }
                    fetchData(1, pagination.pageSize || 10, searchParams);
                  }}
                >
                  刷新
                </Button>
                <Button
                  icon={<ExportOutlined />}
                  onClick={() => message.info('导出功能开发中')}
                >
                  导出
                </Button>
                <Button
                  icon={<ImportOutlined />}
                  onClick={() => message.info('导入功能开发中')}
                >
                  导入
                </Button>
                <Button
                  type="primary"
                  icon={<PlusOutlined />}
                  onClick={() => openModal()}
                >
                  新增数据
                </Button>
              </Space>
            }
          >
            {/* 搜索区域 */}
            {searchVisible && (
              <Card size="small" style={{ marginBottom: 16 }}>
                <Form
                  form={searchForm}
                  layout="inline"
                  onFinish={handleSearch}
                >
                  {fields.filter(field => field.is_searchable).map(field => (
                    <Form.Item
                      key={field.field_name}
                      name={field.field_name}
                      label={field.display_name}
                    >
                      {field.field_type === 'select' ? (
                        <Select 
                          placeholder={`请选择${field.display_name}`}
                          allowClear
                          style={{ width: 150 }}
                        >
                          {field.options && field.options.map((option, index) => {
                            const optionValue = typeof option === 'object' ? option.value : option;
                            const optionLabel = typeof option === 'object' ? option.label : option;
                            return <Option key={`${optionValue}_${index}`} value={optionValue}>{optionLabel}</Option>;
                          })}
                        </Select>
                      ) : (
                        <Input 
                          placeholder={`请输入${field.display_name}`}
                          style={{ width: 150 }}
                        />
                      )}
                    </Form.Item>
                  ))}
                  <Form.Item>
                    <Space>
                      <Button type="primary" htmlType="submit">
                        搜索
                      </Button>
                      <Button onClick={handleResetSearch}>
                        重置
                      </Button>
                    </Space>
                  </Form.Item>
                </Form>
              </Card>
            )}

            {/* 调试信息 */}
            <div style={{ marginBottom: 16, padding: 8, background: '#f0f0f0', borderRadius: 4 }}>
              <div>数据条数: {data.length} | 字段数量: {fields.length} | 加载状态: {loading ? '加载中' : '完成'}</div>
              {data.length > 0 && (
                <div style={{ fontSize: 12, color: '#666' }}>
                  数据示例: {JSON.stringify(data[0]).substring(0, 200)}...
                </div>
              )}
              {fields.length > 0 && (
                <div style={{ fontSize: 12, color: '#666' }}>
                  字段: {fields.map(f => f.field_name).join(', ')}
                </div>
              )}
            </div>

            {/* 批量操作 */}
            {selectedRowKeys.length > 0 && (
              <div style={{ marginBottom: 16 }}>
                <Space>
                  <span>已选择 {selectedRowKeys.length} 项</span>
                  <Popconfirm
                    title="确定要删除选中的数据吗？"
                    onConfirm={handleBatchDelete}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button danger>批量删除</Button>
                  </Popconfirm>
                </Space>
              </div>
            )}

            <Table
              columns={generateColumns()}
              dataSource={data}
              rowKey="id"
              loading={loading}
              rowSelection={rowSelection}
              scroll={{ x: 'max-content' }}
              pagination={{
                ...pagination,
                showSizeChanger: true,
                showQuickJumper: true,
                showTotal: (total, range) =>
                  `第 ${range[0]}-${range[1]} 条/共 ${total} 条`,
                onChange: (page, pageSize) => {
                  fetchData(page, pageSize, searchParams);
                },
              }}
            />
          </Card>
        </Col>
      </Row>

      {/* 创建/编辑数据模态框 */}
      <Modal
        title={editingData ? '编辑数据' : '新增数据'}
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
            {fields.map((field, index) => (
              <Col span={index % 2 === 0 ? 12 : 12} key={field.field_name}>
                {renderFormItem(field)}
              </Col>
            ))}
          </Row>

          <Form.Item>
            <Space>
              <Button type="primary" htmlType="submit">
                {editingData ? '更新' : '创建'}
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

export default DynamicDataManagement;