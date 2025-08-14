import React, { useEffect, useState } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { permissionApi } from '../api';

// 权限守卫组件
const PermissionGuard = ({ children, requiredPermission, fallback = null }) => {
  const [hasPermission, setHasPermission] = useState(null);
  const [loading, setLoading] = useState(true);
  const location = useLocation();

  useEffect(() => {
    checkPermission();
  }, [location.pathname, requiredPermission]);

  const checkPermission = async () => {
    try {
      setLoading(true);
      
      // 如果没有指定权限要求，直接允许访问
      if (!requiredPermission) {
        setHasPermission(true);
        setLoading(false);
        return;
      }

      // 检查用户是否有访问权限
      const response = await permissionApi.checkPermission(
        requiredPermission.resource,
        requiredPermission.action
      );
      setHasPermission(response.data.hasPermission);
    } catch (error) {
      console.error('权限检查失败:', error);
      setHasPermission(false);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="permission-loading">
        <div className="permission-loading-spinner">
          <i className="fas fa-spinner"></i>
          <span>检查权限中...</span>
        </div>
      </div>
    );
  }

  if (!hasPermission) {
    if (fallback) {
      return fallback;
    }
    
    return (
      <div className="permission-denied">
        <div className="permission-denied-content">
          <i className="fas fa-lock fa-3x"></i>
          <h3>访问被拒绝</h3>
          <p>您没有访问此页面的权限</p>
          <button onClick={() => window.history.back()}>返回上一页</button>
        </div>
      </div>
    );
  }

  return children;
};

// 路由权限守卫高阶组件
export const withPermission = (Component, permission) => {
  return (props) => (
    <PermissionGuard requiredPermission={permission}>
      <Component {...props} />
    </PermissionGuard>
  );
};

// 菜单权限过滤Hook
export const useMenuPermissions = () => {
  const [userPermissions, setUserPermissions] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchUserPermissions();
  }, []);

  const fetchUserPermissions = async () => {
    try {
      const userId = localStorage.getItem('userId');
      if (!userId) {
        setLoading(false);
        return;
      }

      const response = await permissionApi.getUserPermissions(userId);
      setUserPermissions(response.data.data?.menus || []);
    } catch (error) {
      console.error('获取用户权限失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 检查是否有菜单权限
  const hasMenuPermission = (menuPath) => {
    return userPermissions.some(permission => 
      permission.path === menuPath || permission.resource === menuPath
    );
  };

  // 过滤有权限的菜单
  const filterMenusByPermission = (menus) => {
    return menus.filter(menu => {
      // 如果是父菜单，检查是否有子菜单有权限
      if (menu.children && menu.children.length > 0) {
        const filteredChildren = filterMenusByPermission(menu.children);
        return filteredChildren.length > 0;
      }
      
      // 检查单个菜单权限
      return hasMenuPermission(menu.path);
    }).map(menu => {
      // 递归过滤子菜单
      if (menu.children && menu.children.length > 0) {
        return {
          ...menu,
          children: filterMenusByPermission(menu.children)
        };
      }
      return menu;
    });
  };

  return {
    userPermissions,
    loading,
    hasMenuPermission,
    filterMenusByPermission,
    refreshPermissions: fetchUserPermissions
  };
};

// 按钮权限控制组件
export const PermissionButton = ({ 
  children, 
  permission, 
  fallback = null, 
  ...props 
}) => {
  const [hasPermission, setHasPermission] = useState(false);

  useEffect(() => {
    checkButtonPermission();
  }, [permission]);

  const checkButtonPermission = async () => {
    if (!permission) {
      setHasPermission(true);
      return;
    }

    try {
      const response = await permissionApi.checkPermission(
        permission.resource,
        permission.action
      );
      setHasPermission(response.data.hasPermission);
    } catch (error) {
      console.error('按钮权限检查失败:', error);
      setHasPermission(false);
    }
  };

  if (!hasPermission) {
    return fallback;
  }

  return (
    <button {...props}>
      {children}
    </button>
  );
};

export default PermissionGuard;