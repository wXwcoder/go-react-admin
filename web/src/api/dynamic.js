import api from './index';

// 动态表管理API
export const dynamicTableApi = {
  // 创建动态表
  createTable: (data) => api.post('/dynamicTable/createTable', data),
  // 获取动态表列表
  getTableList: (params) => api.get('/dynamicTable/getTableList', { params }),
  // 根据ID获取动态表
  getTableByID: (id) => api.get(`/dynamicTable/getTable/${id}`),
  // 更新动态表
  updateTable: (data) => api.put(`/dynamicTable/updateTable/${data.id}`, data),
  // 删除动态表
  deleteTable: (id) => api.delete(`/dynamicTable/deleteTable/${id}`),
  // 切换表状态
  toggleTableStatus: (id) => api.patch(`/dynamicTable/toggleStatus/${id}`),
  // 获取表结构信息
  getTableSchema: (tableName) => api.get(`/dynamicTable/getSchema/${tableName}`),
  // 验证表名
  validateTableName: (name) => api.get('/dynamicTable/validateTableName', { params: { name } }),
};

// 动态字段管理API
export const dynamicFieldApi = {
  // 创建动态字段
  createField: (data) => api.post('/dynamicField/createField', data),
  // 根据表ID获取字段列表
  getFieldsByTableID: (tableId) => api.get(`/dynamicField/getFields/${tableId}`),
  // 根据ID获取字段
  getFieldByID: (id) => api.get(`/dynamicField/getField/${id}`),
  // 更新字段
  updateField: (data) => api.put(`/dynamicField/updateField/${data.id}`, data),
  // 删除字段
  deleteField: (id) => api.delete(`/dynamicField/deleteField/${id}`),
  // 更新字段排序
  updateFieldOrder: (tableId, data) => api.patch(`/dynamicField/updateOrder/${tableId}`, data),
  // 切换字段状态
  toggleFieldStatus: (id) => api.patch(`/dynamicField/toggleStatus/${id}`),
  // 批量创建字段
  batchCreateFields: (data) => api.post('/dynamicField/batchCreate', data),
  // 获取支持的字段类型
  getFieldTypes: () => api.get('/dynamicField/getFieldTypes'),
};

// 动态数据管理API
export const dynamicDataApi = {
  // 创建动态数据
  createData: (tableName, data) => api.post(`/dynamicData/${tableName}/create`, data),
  // 获取动态数据列表
  getDataList: (tableName, params) => api.get(`/dynamicData/${tableName}/list`, { params }),
  // 根据ID获取动态数据
  getDataByID: (tableName, id) => api.get(`/dynamicData/${tableName}/get/${id}`),
  // 更新动态数据
  updateData: (tableName, id, data) => api.put(`/dynamicData/${tableName}/update/${id}`, data),
  // 删除动态数据
  deleteData: (tableName, id) => api.delete(`/dynamicData/${tableName}/delete/${id}`),
  // 批量删除动态数据
  batchDeleteData: (tableName, ids) => api.delete(`/dynamicData/${tableName}/batchDelete`, { data: { ids } }),
  // 高级查询动态数据
  queryData: (tableName, query) => api.post(`/dynamicData/${tableName}/query`, query),
  // 获取数据统计
  getDataStatistics: (tableName, params) => api.post(`/dynamicData/${tableName}/statistics`, params),
};

// 动态数据导入导出API
export const dynamicImportExportApi = {
  // 导出数据
  exportData: (tableName, params) => api.get(`/dynamicData/${tableName}/export`, { 
    params,
    responseType: 'blob' // 用于下载文件
  }),
  // 导入数据
  importData: (tableName, file) => {
    const formData = new FormData();
    formData.append('file', file);
    return api.post(`/dynamicData/${tableName}/import`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
  // 获取导入模板
  getImportTemplate: (tableName) => api.get(`/dynamicData/${tableName}/template`, {
    responseType: 'blob'
  }),
  // 获取导入导出记录
  getImportExportLogs: (params) => api.get('/dynamicData/logs', { params }),
};

// 动态数据权限API
export const dynamicPermissionApi = {
  // 获取表权限
  getTablePermissions: (tableId) => api.get(`/dynamicPermission/table/${tableId}`),
  // 设置表权限
  setTablePermissions: (data) => api.post('/dynamicPermission/table', data),
  // 检查表权限
  checkTablePermission: (tableName, action) => api.get('/dynamicPermission/check', {
    params: { tableName, action }
  }),
  // 获取用户表权限列表
  getUserTablePermissions: (userId) => api.get(`/dynamicPermission/user/${userId}/tables`),
};

// 动态数据视图API
export const dynamicViewApi = {
  // 创建数据视图
  createView: (data) => api.post('/dynamicView/create', data),
  // 获取视图列表
  getViewList: (tableId) => api.get(`/dynamicView/list/${tableId}`),
  // 获取视图详情
  getViewByID: (id) => api.get(`/dynamicView/get/${id}`),
  // 更新视图
  updateView: (data) => api.put(`/dynamicView/update/${data.id}`, data),
  // 删除视图
  deleteView: (id) => api.delete(`/dynamicView/delete/${id}`),
  // 应用视图
  applyView: (viewId, params) => api.post(`/dynamicView/apply/${viewId}`, params),
};