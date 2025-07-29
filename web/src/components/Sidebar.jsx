import React, { useState, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useMenu } from '../store/MenuContext';
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
    favoriteMenus,
    recentMenus,
    addFavorite,
    removeFavorite,
    menuSearchText,
    setSearchText,
    fetchUserMenus
  } = useMenu();

  const [expandedMenus, setExpandedMenus] = useState(new Set());
  const [activeView, setActiveView] = useState('all'); // 'all', 'favorites', 'recent'

  useEffect(() => {
    // 获取当前用户的菜单权限
    const userId = localStorage.getItem('userId');
    if (userId) {
      fetchUserMenus(userId);
    }
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
    if (menu.children && menu.children.length > 0) {
      toggleMenuExpand(menu.id);
    } else {
      openTab(menu);
      navigate(menu.path);
    }
  };

  // 处理收藏切换
  const handleFavoriteToggle = (e, menuId) => {
    e.stopPropagation();
    if (favoriteMenus.includes(menuId)) {
      removeFavorite(menuId);
    } else {
      addFavorite(menuId);
    }
  };

  // 渲染菜单项
  const renderMenuItem = (menu, level = 0) => {
    const isActive = location.pathname === menu.path;
    const isExpanded = expandedMenus.has(menu.id);
    const isFavorite = favoriteMenus.includes(menu.id);
    const hasChildren = menu.children && menu.children.length > 0;

    return (
      <div key={menu.id} className="menu-item-container">
        <div
          className={`menu-item ${isActive ? 'active' : ''} ${sidebarCollapsed ? 'collapsed' : ''}`}
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
                <span className="menu-text">{menu.name || menu.title}</span>
                <div className="menu-actions">
                  <button
                    className={`favorite-btn ${isFavorite ? 'active' : ''}`}
                    onClick={(e) => handleFavoriteToggle(e, menu.id)}
                    title={isFavorite ? '取消收藏' : '添加收藏'}
                  >
                    <i className={isFavorite ? 'fas fa-star' : 'far fa-star'}></i>
                  </button>
                  {hasChildren && (
                    <span className={`expand-icon ${isExpanded ? 'expanded' : ''}`}>
                      <i className="fas fa-chevron-down"></i>
                    </span>
                  )}
                </div>
              </>
            )}
          </div>
        </div>
        
        {hasChildren && isExpanded && !sidebarCollapsed && (
          <div className="submenu">
            {menu.children.map(child => renderMenuItem(child, level + 1))}
          </div>
        )}
      </div>
    );
  };

  // 获取要显示的菜单
  const getDisplayMenus = () => {
    const filteredMenus = getFilteredMenus();
    const menuTree = buildMenuTree(filteredMenus);

    switch (activeView) {
      case 'favorites':
        return menuTree.filter(menu => 
          favoriteMenus.includes(menu.id) || 
          (menu.children && menu.children.some(child => favoriteMenus.includes(child.id)))
        );
      case 'recent':
        return menuTree.filter(menu => 
          recentMenus.includes(menu.id) || 
          (menu.children && menu.children.some(child => recentMenus.includes(child.id)))
        );
      default:
        return menuTree;
    }
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

      {/* 搜索区域 */}
      {!sidebarCollapsed && (
        <div className="sidebar-search">
          <div className="search-input-container">
            <i className="fas fa-search search-icon"></i>
            <input
              type="text"
              placeholder="搜索菜单..."
              value={menuSearchText}
              onChange={(e) => setSearchText(e.target.value)}
              className="search-input"
            />
            {menuSearchText && (
              <button
                className="clear-search"
                onClick={() => setSearchText('')}
              >
                <i className="fas fa-times"></i>
              </button>
            )}
          </div>
        </div>
      )}

      {/* 视图切换 */}
      {!sidebarCollapsed && (
        <div className="sidebar-tabs">
          <button
            className={`tab-btn ${activeView === 'all' ? 'active' : ''}`}
            onClick={() => setActiveView('all')}
          >
            <i className="fas fa-list"></i>
            <span>全部</span>
          </button>
          <button
            className={`tab-btn ${activeView === 'favorites' ? 'active' : ''}`}
            onClick={() => setActiveView('favorites')}
          >
            <i className="fas fa-star"></i>
            <span>收藏</span>
          </button>
          <button
            className={`tab-btn ${activeView === 'recent' ? 'active' : ''}`}
            onClick={() => setActiveView('recent')}
          >
            <i className="fas fa-clock"></i>
            <span>最近</span>
          </button>
        </div>
      )}

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