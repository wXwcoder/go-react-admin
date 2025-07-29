import React, { useState, useRef, useEffect } from 'react';
import { useNavigate, useLocation } from 'react-router-dom';
import { useMenu } from '../store/MenuContext';
import '../assets/styles/TabBar.css';

const TabBar = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const {
    openTabs,
    activeTab,
    setActiveTab,
    closeTab,
    closeOtherTabs,
    closeAllTabs,
    refreshTab
  } = useMenu();

  const [contextMenu, setContextMenu] = useState({ visible: false, x: 0, y: 0, tabId: null });
  const [draggedTab, setDraggedTab] = useState(null);
  const tabsRef = useRef(null);

  // 处理标签页点击
  const handleTabClick = (tab) => {
    setActiveTab(tab.id);
    navigate(tab.path);
  };

  // 处理标签页关闭
  const handleTabClose = (e, tabId) => {
    e.stopPropagation();
    closeTab(tabId);
  };

  // 处理右键菜单
  const handleContextMenu = (e, tabId) => {
    e.preventDefault();
    setContextMenu({
      visible: true,
      x: e.clientX,
      y: e.clientY,
      tabId
    });
  };

  // 关闭右键菜单
  const closeContextMenu = () => {
    setContextMenu({ visible: false, x: 0, y: 0, tabId: null });
  };

  // 处理右键菜单操作
  const handleContextMenuAction = (action) => {
    const { tabId } = contextMenu;
    
    switch (action) {
      case 'refresh':
        refreshTab();
        break;
      case 'close':
        closeTab(tabId);
        break;
      case 'closeOthers':
        closeOtherTabs(tabId);
        break;
      case 'closeAll':
        closeAllTabs();
        break;
      case 'closeLeft':
        const leftTabs = openTabs.slice(0, openTabs.findIndex(tab => tab.id === tabId));
        leftTabs.forEach(tab => {
          if (tab.closable) closeTab(tab.id);
        });
        break;
      case 'closeRight':
        const rightTabs = openTabs.slice(openTabs.findIndex(tab => tab.id === tabId) + 1);
        rightTabs.forEach(tab => {
          if (tab.closable) closeTab(tab.id);
        });
        break;
      default:
        break;
    }
    
    closeContextMenu();
  };

  // 拖拽开始
  const handleDragStart = (e, tab) => {
    setDraggedTab(tab);
    e.dataTransfer.effectAllowed = 'move';
  };

  // 拖拽结束
  const handleDragEnd = () => {
    setDraggedTab(null);
  };

  // 拖拽悬停
  const handleDragOver = (e) => {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
  };

  // 拖拽放置
  const handleDrop = (e, targetTab) => {
    e.preventDefault();
    
    if (!draggedTab || draggedTab.id === targetTab.id) return;

    // 重新排序标签页逻辑可以在这里实现
    // 这里简化处理，实际项目中可能需要更复杂的排序逻辑
  };

  // 点击外部关闭右键菜单
  useEffect(() => {
    const handleClickOutside = () => {
      if (contextMenu.visible) {
        closeContextMenu();
      }
    };

    document.addEventListener('click', handleClickOutside);
    return () => document.removeEventListener('click', handleClickOutside);
  }, [contextMenu.visible]);

  // 滚动到活动标签页
  const scrollToActiveTab = () => {
    if (tabsRef.current) {
      const activeTabElement = tabsRef.current.querySelector('.tab-item.active');
      if (activeTabElement) {
        activeTabElement.scrollIntoView({ behavior: 'smooth', inline: 'center' });
      }
    }
  };

  useEffect(() => {
    scrollToActiveTab();
  }, [activeTab]);

  return (
    <div className="tab-bar">
      <div className="tab-bar-content">
        {/* 标签页列表 */}
        <div className="tab-container" ref={tabsRef}>
          <div className="tabs-list" style={{display: 'flex', flexDirection: 'row', height: '100%'}}>
            {openTabs.map((tab) => (
              <div
                key={tab.id}
                className={`tab-item ${activeTab === tab.id ? 'active' : ''} ${!tab.closable ? 'non-closable' : ''}`}
                onClick={() => handleTabClick(tab)}
                onContextMenu={(e) => handleContextMenu(e, tab.id)}
                draggable={tab.closable}
                onDragStart={(e) => handleDragStart(e, tab)}
                onDragEnd={handleDragEnd}
                onDragOver={handleDragOver}
                onDrop={(e) => handleDrop(e, tab)}
              >
                <span className="tab-icon">
                  <i className="fas fa-file-alt"></i>
                </span>
                <span className="tab-title" title={tab.name}>
                  {tab.name}
                </span>
                {tab.closable && (
                  <button
                    className="tab-close"
                    onClick={(e) => handleTabClose(e, tab.id)}
                    title="关闭标签页"
                  >
                    <i className="fas fa-times"></i>
                  </button>
                )}
              </div>
            ))}
          </div>
        </div>

        {/* 标签页操作按钮 */}
        <div className="tab-actions">
          <button
            className="tab-action-btn"
            onClick={refreshTab}
            title="刷新当前页面"
          >
            <i className="fas fa-sync-alt"></i>
          </button>
          <button
            className="tab-action-btn"
            onClick={() => closeOtherTabs(activeTab)}
            title="关闭其他标签页"
          >
            <i className="fas fa-window-close"></i>
          </button>
          <button
            className="tab-action-btn"
            onClick={closeAllTabs}
            title="关闭所有标签页"
          >
            <i className="fas fa-times-circle"></i>
          </button>
        </div>
      </div>

      {/* 右键菜单 */}
      {contextMenu.visible && (
        <div
          className="tab-context-menu"
          style={{
            left: contextMenu.x,
            top: contextMenu.y
          }}
        >
          <div className="context-menu-item" onClick={() => handleContextMenuAction('refresh')}>
            <i className="fas fa-sync-alt"></i>
            <span>刷新</span>
          </div>
          {openTabs.find(tab => tab.id === contextMenu.tabId)?.closable && (
            <div className="context-menu-item" onClick={() => handleContextMenuAction('close')}>
              <i className="fas fa-times"></i>
              <span>关闭</span>
            </div>
          )}
          <div className="context-menu-item" onClick={() => handleContextMenuAction('closeOthers')}>
            <i className="fas fa-window-close"></i>
            <span>关闭其他</span>
          </div>
          <div className="context-menu-item" onClick={() => handleContextMenuAction('closeLeft')}>
            <i className="fas fa-arrow-left"></i>
            <span>关闭左侧</span>
          </div>
          <div className="context-menu-item" onClick={() => handleContextMenuAction('closeRight')}>
            <i className="fas fa-arrow-right"></i>
            <span>关闭右侧</span>
          </div>
          <div className="context-menu-divider"></div>
          <div className="context-menu-item" onClick={() => handleContextMenuAction('closeAll')}>
            <i className="fas fa-times-circle"></i>
            <span>关闭所有</span>
          </div>
        </div>
      )}
    </div>
  );
};

export default TabBar;