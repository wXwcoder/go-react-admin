import React, { Suspense, lazy } from 'react';
import { Navigate } from 'react-router-dom';
import PermissionGuard from '../components/PermissionGuard';

// 懒加载页面组件
const Dashboard = lazy(() => import('../pages/Dashboard'));
const UserManagement = lazy(() => import('../pages/UserManagement'));
const RoleManagement = lazy(() => import('../pages/RoleManagement'));
const MenuManagement = lazy(() => import('../pages/MenuManagement'));
const ApiManagement = lazy(() => import('../pages/ApiManagement'));
const PermissionManagement = lazy(() => import('../pages/PermissionManagement'));
const LogManagement = lazy(() => import('../pages/LogManagement'));
const UserProfile = lazy(() => import('../pages/UserProfile'));

// 动态数据管理页面
const DynamicTableManagement = lazy(() => import('../pages/DynamicTableManagement'));
const DynamicFieldManagement = lazy(() => import('../pages/DynamicFieldManagement'));
const DynamicDataManagement = lazy(() => import('../pages/DynamicDataManagement'));

// 其他页面
const About = lazy(() => import('../pages/About'));

// 第三方客户系统页面
const CustomerLogin = lazy(() => import('../pages/CustomerLogin'));
const CustomerRegister = lazy(() => import('../pages/CustomerRegister'));
const CustomerProfile = lazy(() => import('../pages/CustomerProfile'));
const CustomerMessages = lazy(() => import('../pages/CustomerMessages'));
const CustomerManagement = lazy(() => import('../pages/CustomerManagement'));

// 加载组件包装器
const LoadingWrapper = ({ children }) => (
  <Suspense fallback={
    <div className="page-loading">
      <div className="loading-spinner">
        <i className="fas fa-spinner fa-spin"></i>
        <span>页面加载中...</span>
      </div>
    </div>
  }>
    {children}
  </Suspense>
);

// 权限包装器
const PermissionWrapper = ({ children, permission }) => (
  <PermissionGuard requiredPermission={permission}>
    <LoadingWrapper>
      {children}
    </LoadingWrapper>
  </PermissionGuard>
);

// 静态路由配置
export const staticRoutes = [
  {
    path: '/',
    element: <Navigate to="/dashboard" replace />,
    meta: {
      title: '首页',
      hidden: true
    }
  },
  {
    path: '/dashboard',
    element: (
      <PermissionWrapper permission={{ resource: 'dashboard', action: 'read' }}>
        <Dashboard />
      </PermissionWrapper>
    ),
    meta: {
      title: '仪表板',
      icon: 'fas fa-tachometer-alt',
      keepAlive: true
    }
  },
  // 系统管理路由
  {
    path: '/system',
    element: <Navigate to="/system/users" replace />,
    meta: {
      title: '系统管理',
      hidden: true
    }
  },
  {
    path: '/system/users',
    element: (
      <PermissionWrapper permission={{ resource: 'user', action: 'read' }}>
        <UserManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '用户管理',
      icon: 'fas fa-users',
      keepAlive: true
    }
  },
  {
    path: '/system/roles',
    element: (
      <PermissionWrapper permission={{ resource: 'role', action: 'read' }}>
        <RoleManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '角色管理',
      icon: 'fas fa-user-tag',
      keepAlive: true
    }
  },
  {
    path: '/system/menus',
    element: (
      <PermissionWrapper permission={{ resource: 'menu', action: 'read' }}>
        <MenuManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '菜单管理',
      icon: 'fas fa-list',
      keepAlive: true
    }
  },
  {
    path: '/system/apis',
    element: (
      <PermissionWrapper permission={{ resource: 'api', action: 'read' }}>
        <ApiManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: 'API管理',
      icon: 'fas fa-code',
      keepAlive: true
    }
  },
  {
    path: '/system/permissions',
    element: (
      <PermissionWrapper permission={{ resource: 'permission', action: 'read' }}>
        <PermissionManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '权限管理',
      icon: 'fas fa-shield-alt',
      keepAlive: true
    }
  },
  {
    path: '/system/logs',
    element: (
      <PermissionWrapper permission={{ resource: 'log', action: 'read' }}>
        <LogManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '日志管理',
      icon: 'fas fa-history',
      keepAlive: true
    }
  },
  
  // 第三方客户管理路由
  {
    path: '/customers',
    element: (
      <PermissionWrapper permission={{ resource: 'customer', action: 'read' }}>
        <CustomerManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '客户管理',
      icon: 'fas fa-user-friends',
      keepAlive: true
    }
  },
  
  // 动态数据管理路由
  {
    path: '/dynamic',
    element: <Navigate to="/dynamic/tables" replace />,
    meta: {
      title: '动态数据管理',
      hidden: true
    }
  },
  {
    path: '/dynamic/tables',
    element: (
      <PermissionWrapper permission={{ resource: 'dynamic_table', action: 'read' }}>
        <DynamicTableManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '动态表管理',
      icon: 'fas fa-table',
      keepAlive: true
    }
  },
  {
    path: '/dynamic/fields/:tableId',
    element: (
      <PermissionWrapper permission={{ resource: 'dynamic_field', action: 'read' }}>
        <DynamicFieldManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '字段管理',
      icon: 'fas fa-columns',
      keepAlive: false,
      hidden: true
    }
  },
  {
    path: '/dynamic/data/:tableName',
    element: (
      <PermissionWrapper permission={{ resource: 'dynamic_data', action: 'read' }}>
        <DynamicDataManagement />
      </PermissionWrapper>
    ),
    meta: {
      title: '数据管理',
      icon: 'fas fa-database',
      keepAlive: false,
      hidden: true
    }
  },
  
  // 兼容旧路径的重定向
  {
    path: '/dynamic-tables',
    element: <Navigate to="/dynamic/tables" replace />,
    meta: {
      title: '动态表管理',
      hidden: true
    }
  },
  {
    path: '/dynamic-fields/:tableId',
    element: <Navigate to="/dynamic/fields/:tableId" replace />,
    meta: {
      title: '字段管理',
      hidden: true
    }
  },
  {
    path: '/dynamic-data/:tableName',
    element: <Navigate to="/dynamic/data/:tableName" replace />,
    meta: {
      title: '数据管理',
      hidden: true
    }
  },
  // 兼容旧路径的重定向
  {
    path: '/users',
    element: <Navigate to="/system/users" replace />,
    meta: {
      title: '用户管理',
      hidden: true
    }
  },
  {
    path: '/roles',
    element: <Navigate to="/system/roles" replace />,
    meta: {
      title: '角色管理',
      hidden: true
    }
  },
  {
    path: '/menus',
    element: <Navigate to="/system/menus" replace />,
    meta: {
      title: '菜单管理',
      hidden: true
    }
  },
  {
    path: '/apis',
    element: <Navigate to="/system/apis" replace />,
    meta: {
      title: 'API管理',
      hidden: true
    }
  },
  {
    path: '/permissions',
    element: <Navigate to="/system/permissions" replace />,
    meta: {
      title: '权限管理',
      hidden: true
    }
  },
  {
    path: '/logs',
    element: <Navigate to="/system/logs" replace />,
    meta: {
      title: '日志管理',
      hidden: true
    }
  },
  {
    path: '/profile',
    element: (
      <PermissionWrapper permission={{ resource: 'profile', action: 'read' }}>
        <UserProfile />
      </PermissionWrapper>
    ),
    meta: {
      title: '个人信息',
      icon: 'fas fa-user',
      keepAlive: true,
      hidden: true
    }
  },
  {
    path: '/about',
    element: (
      <PermissionWrapper permission={{ resource: 'about', action: 'read' }}>
        <About />
      </PermissionWrapper>
    ),
    meta: {
      title: '关于我们',
      icon: 'fas fa-info-circle',
      keepAlive: true
    }
  }
];

// 动态路由生成器
export const generateDynamicRoutes = (menuData) => {
  const routes = [];

  const processMenu = (menu, parentPath = '') => {
    const fullPath = parentPath + menu.path;
    
    if (menu.component) {
      // 动态导入组件
      const Component = lazy(() => 
        import(`../pages/${menu.component}`).catch(() => 
          import('../pages/NotFound')
        )
      );

      routes.push({
        path: fullPath,
        element: (
          <PermissionWrapper 
            permission={{ 
              resource: menu.resource || menu.path.replace('/', ''), 
              action: 'read' 
            }}
          >
            <Component />
          </PermissionWrapper>
        ),
        meta: {
          title: menu.title || menu.name,
          icon: menu.icon,
          keepAlive: menu.keepAlive !== false,
          hidden: menu.hidden === true,
          menuId: menu.id,
          parentId: menu.parentId
        }
      });
    }

    // 处理子菜单
    if (menu.children && menu.children.length > 0) {
      menu.children.forEach(child => {
        processMenu(child, fullPath);
      });
    }
  };

  menuData.forEach(menu => {
    processMenu(menu);
  });

  return routes;
};

// 路由元信息处理
export const getRouteMetaByPath = (path, routes) => {
  const route = routes.find(r => r.path === path);
  return route?.meta || {};
};

// 面包屑生成
export const generateBreadcrumb = (path, routes, menuData) => {
  const breadcrumb = [];
  const pathSegments = path.split('/').filter(Boolean);
  
  let currentPath = '';
  pathSegments.forEach(segment => {
    currentPath += `/${segment}`;
    const route = routes.find(r => r.path === currentPath);
    if (route && route.meta) {
      breadcrumb.push({
        path: currentPath,
        title: route.meta.title,
        icon: route.meta.icon
      });
    }
  });

  return breadcrumb;
};

// 路由权限检查
export const checkRoutePermission = async (path, userPermissions) => {
  // 查找对应路由的权限要求
  const route = staticRoutes.find(r => r.path === path);
  if (!route || !route.meta?.permission) {
    return true; // 无权限要求的路由默认允许访问
  }

  const requiredPermission = route.meta.permission;
  return userPermissions.some(permission => 
    permission.resource === requiredPermission.resource &&
    permission.action === requiredPermission.action
  );
};

// 获取用户可访问的路由
export const getAccessibleRoutes = (routes, userPermissions) => {
  return routes.filter(route => {
    if (!route.meta?.permission) {
      return true;
    }
    
    const { resource, action } = route.meta.permission;
    return userPermissions.some(permission => 
      permission.resource === resource && permission.action === action
    );
  });
};

export default {
  staticRoutes,
  generateDynamicRoutes,
  getRouteMetaByPath,
  generateBreadcrumb,
  checkRoutePermission,
  getAccessibleRoutes
};