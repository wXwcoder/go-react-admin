import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useMenu } from '../store/MenuContext';
import { ICONS } from '../constants/icons';
import '../assets/styles/Sidebar.css';

const Sidebar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    sidebarCollapsed,
    toggleSidebar,
    openTab,
    getFilteredMenus,
    buildMenuTree,
    fetchUserMenus
  } = useMenu();

  const [expandedMenus, setExpandedMenus] = useState(new Set(['system'])); // 默认展开系统管理

  useEffect(() => {
    // 获取当前用户的菜单权限
    fetchUserMenus();
  }, []); // 空依赖数组，只在组件挂载时执行一次

  // 切换菜单展开状态
  const toggleMenuExpand = (menuId) => {
    const newExpanded = new Set(expandedMenus);
    if (newExpanded.has(menuId)) {
      newExpanded.delete(menuId);
    } else {
      newExpanded.add(menuId);
    }
    setExpandedMenus(newExpanded);
  };

  // 处理菜单点击
  const handleMenuClick = (menu) => {
    if (menu.type === 'group' || (menu.children && menu.children.length > 0)) {
      // 菜单组或有子菜单的项目，切换展开状态
      toggleMenuExpand(menu.name || menu.id);
    } else {
      // 普通菜单项，导航到对应页面
      openTab(menu);
      navigate(menu.path);
    }
  };

  // 渲染菜单项
  const renderMenuItem = (menu, level = 0) => {
    const isActive = location.pathname === menu.path;
    const isExpanded = expandedMenus.has(menu.name || menu.id);
    const hasChildren = menu.children && menu.children.length > 0;
    const isGroup = menu.type === 'group';

    return (
      <div key={menu.id} className="menu-item-container">
        <div
          className={`menu-item ${isActive ? 'active' : ''} ${isGroup ? 'menu-group' : ''} ${sidebarCollapsed ? 'collapsed' : ''}`}
          style={{ paddingLeft: `${20 + level * 16}px` }}
          onClick={() => handleMenuClick(menu)}
        >
          <div className="menu-item-content">
            {menu.icon && (
              <span className="menu-icon">
                <i className={menu.icon}></i>
              </span>
            )}
            {!sidebarCollapsed && (
              <>
                <span className="menu-text">{menu.title || menu.name}</span>
                {(hasChildren || isGroup) && (
                  <span className={`expand-icon ${isExpanded ? 'expanded' : ''}`}>
                    <i className="fas fa-chevron-down"></i>
                  </span>
                )}
              </>
            )}
          </div>
        </div>
        
        {(hasChildren || isGroup) && isExpanded && !sidebarCollapsed && (
          <div className="submenu">
            {menu.children && menu.children.map(child => renderMenuItem(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  // 获取要显示的菜单
  const getDisplayMenus = () => {
    const filteredMenus = getFilteredMenus();
    return buildMenuTree(filteredMenus);
  };

  return (
    <aside className={`sidebar ${sidebarCollapsed ? 'collapsed' : ''}`}>
      {/* Logo区域 */}
      <div className="sidebar-header">
        <div className="logo">
          {!sidebarCollapsed && <h2>Go-React Admin</h2>}
          {sidebarCollapsed && <h2>GRA</h2>}
        </div>
        <button className="collapse-btn" onClick={toggleSidebar}>
          <i className={`fas fa-${sidebarCollapsed ? 'angle-right' : 'angle-left'}`}></i>
        </button>
      </div>

      {/* 菜单区域 */}
      <nav className="sidebar-menu">
        <div className="menu-list">
          {getDisplayMenus().map(menu => renderMenuItem(menu))}
        </div>
      </nav>

      {/* 底部信息 */}
      {!sidebarCollapsed && (
        <div className="sidebar-footer">
          <div className="version-info">
            <small>v1.0.0</small>
          </div>
        </div>
      )}
    </aside>
  );
};

export default Sidebar;