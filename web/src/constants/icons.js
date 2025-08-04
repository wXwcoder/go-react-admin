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
  
  // 后台管理相关图标
  ADMIN: 'fas fa-user-shield',
  SYSTEM: 'fas fa-server',
  MONITOR: 'fas fa-desktop',
  DATABASE: 'fas fa-database',
  BACKUP: 'fas fa-archive',
  SECURITY: 'fas fa-shield-alt',
  CONFIG: 'fas fa-sliders-h',
  
  // 计算机相关图标
  COMPUTER: 'fas fa-laptop',
  SERVER: 'fas fa-server',
  NETWORK: 'fas fa-network-wired',
  STORAGE: 'fas fa-hdd',
  MEMORY: 'fas fa-memory',
  CPU: 'fas fa-microchip',
  
  // 运营相关图标
  ANALYTICS: 'fas fa-chart-line',
  REPORT: 'fas fa-chart-bar',
  STATISTICS: 'fas fa-chart-pie',
  MARKETING: 'fas fa-bullhorn',
  CAMPAIGN: 'fas fa-megaphone',
  CUSTOMER: 'fas fa-user-tie',
  SALES: 'fas fa-shopping-cart',
  PRODUCT: 'fas fa-box',
  ORDER: 'fas fa-clipboard-list',
  
  // 运维相关图标
  DEPLOY: 'fas fa-rocket',
  MAINTENANCE: 'fas fa-tools',
  ALERT: 'fas fa-exclamation-triangle',
  MONITORING: 'fas fa-heartbeat',
  AUTOMATION: 'fas fa-robot',
  PIPELINE: 'fas fa-code-branch',
  
  // 云服务相关图标
  CLOUD: 'fas fa-cloud',
  CLOUD_UPLOAD: 'fas fa-cloud-upload-alt',
  CLOUD_DOWNLOAD: 'fas fa-cloud-download-alt',
  AWS: 'fab fa-aws',
  AZURE: 'fab fa-microsoft',
  KUBERNETES: 'fas fa-dharmachakra',
  CONTAINER: 'fab fa-docker',
  CDN: 'fas fa-globe',
  
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
  IMPORT: 'fas fa-upload',
  
  // 状态图标
  SUCCESS: 'fas fa-check-circle',
  WARNING: 'fas fa-exclamation-circle',
  ERROR: 'fas fa-times-circle',
  INFO: 'fas fa-info-circle',
  
  // 关于页面图标
  ABOUT: 'fas fa-info-circle',
  INFO_CIRCLE: 'fas fa-info-circle',
  QUESTION_CIRCLE: 'fas fa-question-circle',
  HELP: 'fas fa-question-circle',
  BOOK: 'fas fa-book',
  DOCUMENTATION: 'fas fa-file-alt',
  GITHUB: 'fab fa-github',
  VERSION: 'fas fa-code-branch',
  LICENSE: 'fas fa-certificate',
  TEAM: 'fas fa-users-cog',
  CONTACT: 'fas fa-address-book',
  ABOUT: 'fas fa-info-circle',
  
  // 操作图标
  VIEW: 'fas fa-eye',
  COPY: 'fas fa-copy',
  DOWNLOAD: 'fas fa-download',
  UPLOAD: 'fas fa-upload',
  SHARE: 'fas fa-share-alt',
  FILTER: 'fas fa-filter',
  SORT: 'fas fa-sort',
  PRINT: 'fas fa-print',
  SAVE: 'fas fa-save',
  
  // 通信图标
  EMAIL: 'fas fa-envelope',
  MESSAGE: 'fas fa-comment',
  PHONE: 'fas fa-phone',
  VIDEO: 'fas fa-video',
  
  // 文件图标
  FILE: 'fas fa-file',
  FOLDER: 'fas fa-folder',
  FOLDER_OPEN: 'fas fa-folder-open',
  ZIP: 'fas fa-file-archive',
  IMAGE: 'fas fa-image',
  PDF: 'fas fa-file-pdf',
  EXCEL: 'fas fa-file-excel',
  WORD: 'fas fa-file-word'
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
  { key: 'HOME', label: '首页', icon: ICONS.HOME },
  
  // 后台管理相关图标
  { key: 'ADMIN', label: '管理员', icon: ICONS.ADMIN },
  { key: 'SYSTEM', label: '系统', icon: ICONS.SYSTEM },
  { key: 'MONITOR', label: '监控', icon: ICONS.MONITOR },
  { key: 'DATABASE', label: '数据库', icon: ICONS.DATABASE },
  { key: 'BACKUP', label: '备份', icon: ICONS.BACKUP },
  { key: 'SECURITY', label: '安全', icon: ICONS.SECURITY },
  { key: 'CONFIG', label: '配置', icon: ICONS.CONFIG },
  
  // 计算机相关图标
  { key: 'COMPUTER', label: '电脑', icon: ICONS.COMPUTER },
  { key: 'SERVER', label: '服务器', icon: ICONS.SERVER },
  { key: 'NETWORK', label: '网络', icon: ICONS.NETWORK },
  { key: 'STORAGE', label: '存储', icon: ICONS.STORAGE },
  { key: 'MEMORY', label: '内存', icon: ICONS.MEMORY },
  { key: 'CPU', label: '处理器', icon: ICONS.CPU },
  
  // 运营相关图标
  { key: 'ANALYTICS', label: '分析', icon: ICONS.ANALYTICS },
  { key: 'REPORT', label: '报告', icon: ICONS.REPORT },
  { key: 'STATISTICS', label: '统计', icon: ICONS.STATISTICS },
  { key: 'MARKETING', label: '营销', icon: ICONS.MARKETING },
  { key: 'CAMPAIGN', label: '活动', icon: ICONS.CAMPAIGN },
  { key: 'CUSTOMER', label: '客户', icon: ICONS.CUSTOMER },
  { key: 'SALES', label: '销售', icon: ICONS.SALES },
  { key: 'PRODUCT', label: '产品', icon: ICONS.PRODUCT },
  { key: 'ORDER', label: '订单', icon: ICONS.ORDER },
  
  // 运维相关图标
  { key: 'DEPLOY', label: '部署', icon: ICONS.DEPLOY },
  { key: 'MAINTENANCE', label: '维护', icon: ICONS.MAINTENANCE },
  { key: 'ALERT', label: '告警', icon: ICONS.ALERT },
  { key: 'MONITORING', label: '监控', icon: ICONS.MONITORING },
  { key: 'AUTOMATION', label: '自动化', icon: ICONS.AUTOMATION },
  { key: 'PIPELINE', label: '流水线', icon: ICONS.PIPELINE },
  
  // 云服务相关图标
  { key: 'CLOUD', label: '云', icon: ICONS.CLOUD },
  { key: 'CLOUD_UPLOAD', label: '云上传', icon: ICONS.CLOUD_UPLOAD },
  { key: 'CLOUD_DOWNLOAD', label: '云下载', icon: ICONS.CLOUD_DOWNLOAD },
  { key: 'KUBERNETES', label: 'Kubernetes', icon: ICONS.KUBERNETES },
  { key: 'CONTAINER', label: '容器', icon: ICONS.CONTAINER },
  { key: 'CDN', label: 'CDN', icon: ICONS.CDN },
  
  // 状态图标
  { key: 'SUCCESS', label: '成功', icon: ICONS.SUCCESS },
  { key: 'WARNING', label: '警告', icon: ICONS.WARNING },
  { key: 'ERROR', label: '错误', icon: ICONS.ERROR },
  { key: 'INFO', label: '信息', icon: ICONS.INFO },
  
  // 操作图标
  { key: 'VIEW', label: '查看', icon: ICONS.VIEW },
  { key: 'COPY', label: '复制', icon: ICONS.COPY },
  { key: 'DOWNLOAD', label: '下载', icon: ICONS.DOWNLOAD },
  { key: 'UPLOAD', label: '上传', icon: ICONS.UPLOAD },
  { key: 'SHARE', label: '分享', icon: ICONS.SHARE },
  { key: 'FILTER', label: '筛选', icon: ICONS.FILTER },
  { key: 'SORT', label: '排序', icon: ICONS.SORT },
  { key: 'PRINT', label: '打印', icon: ICONS.PRINT },
  { key: 'SAVE', label: '保存', icon: ICONS.SAVE },
  
  // 通信图标
  { key: 'EMAIL', label: '邮件', icon: ICONS.EMAIL },
  { key: 'MESSAGE', label: '消息', icon: ICONS.MESSAGE },
  { key: 'PHONE', label: '电话', icon: ICONS.PHONE },
  { key: 'VIDEO', label: '视频', icon: ICONS.VIDEO },
  
  // 文件图标
  { key: 'FILE', label: '文件', icon: ICONS.FILE },
  { key: 'FOLDER', label: '文件夹', icon: ICONS.FOLDER },
  { key: 'ZIP', label: '压缩包', icon: ICONS.ZIP },
  { key: 'IMAGE', label: '图片', icon: ICONS.IMAGE },
  { key: 'PDF', label: 'PDF', icon: ICONS.PDF },
  { key: 'EXCEL', label: 'Excel', icon: ICONS.EXCEL },
  { key: 'WORD', label: 'Word', icon: ICONS.WORD },
  
  // 关于页面图标
  { key: 'ABOUT', label: '关于', icon: ICONS.ABOUT },
  { key: 'INFO_CIRCLE', label: '信息', icon: ICONS.INFO_CIRCLE },
  { key: 'QUESTION_CIRCLE', label: '问题', icon: ICONS.QUESTION_CIRCLE },
  { key: 'HELP', label: '帮助', icon: ICONS.HELP },
  { key: 'BOOK', label: '书籍', icon: ICONS.BOOK },
  { key: 'DOCUMENTATION', label: '文档', icon: ICONS.DOCUMENTATION },
  { key: 'GITHUB', label: 'GitHub', icon: ICONS.GITHUB },
  { key: 'VERSION', label: '版本', icon: ICONS.VERSION },
  { key: 'LICENSE', label: '许可证', icon: ICONS.LICENSE },
  { key: 'TEAM', label: '团队', icon: ICONS.TEAM },
  { key: 'CONTACT', label: '联系', icon: ICONS.CONTACT }
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