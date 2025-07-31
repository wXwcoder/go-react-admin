/**
 * 图标常量定义
 * 统一管理项目中的图标，确保MenuManagement和Sidebar使用一致的图标系统
 */

// Font Awesome 图标映射
export const ICONS = {
  // 系统图标
  DASHBOARD: 'fas fa-tachometer-alt',
  USER: 'fas fa-user',
  SETTING: 'fas fa-cog',
  MENU: 'fas fa-bars',
  ROLE: 'fas fa-users',
  PERMISSION: 'fas fa-lock',
  TENANT: 'fas fa-building',
  API: 'fas fa-plug',
  LOG: 'fas fa-file-alt',
  HOME: 'fas fa-home',
  
  // 导航图标
  CHEVRON_DOWN: 'fas fa-chevron-down',
  CHEVRON_RIGHT: 'fas fa-angle-right',
  CHEVRON_LEFT: 'fas fa-angle-left',
  
  // 通用图标
  QUESTION: 'fas fa-question',
  EDIT: 'fas fa-edit',
  DELETE: 'fas fa-trash',
  ADD: 'fas fa-plus',
  SEARCH: 'fas fa-search',
  REFRESH: 'fas fa-sync',
  EXPORT: 'fas fa-download',
  IMPORT: 'fas fa-upload'
};

// 菜单图标配置
export const MENU_ICONS = [
  { key: 'DASHBOARD', label: '仪表盘', icon: ICONS.DASHBOARD },
  { key: 'USER', label: '用户', icon: ICONS.USER },
  { key: 'SETTING', label: '设置', icon: ICONS.SETTING },
  { key: 'MENU', label: '菜单', icon: ICONS.MENU },
  { key: 'ROLE', label: '角色', icon: ICONS.ROLE },
  { key: 'PERMISSION', label: '权限', icon: ICONS.PERMISSION },
  { key: 'TENANT', label: '租户', icon: ICONS.TENANT },
  { key: 'API', label: 'API', icon: ICONS.API },
  { key: 'LOG', label: '日志', icon: ICONS.LOG },
  { key: 'HOME', label: '首页', icon: ICONS.HOME }
];

// 获取图标类名
export const getIconClass = (iconKey) => {
  return ICONS[iconKey] || ICONS.QUESTION;
};

// 获取图标配置
export const getIconConfig = (iconKey) => {
  return MENU_ICONS.find(item => item.key === iconKey) || 
         MENU_ICONS.find(item => item.icon === iconKey) ||
         { key: 'QUESTION', label: '未知', icon: ICONS.QUESTION };
};

// 获取图标显示文本
export const getIconLabel = (iconClass) => {
  const config = MENU_ICONS.find(item => item.icon === iconClass);
  return config ? config.label : '未知图标';
};