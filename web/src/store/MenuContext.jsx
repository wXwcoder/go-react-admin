import React, { createContext, useContext, useReducer, useEffect } from 'react';
import { menuApi, permissionApi } from '../api';

// 初始状态
const initialState = {
  menus: [],
  userMenus: [],
  openTabs: [{ id: 'dashboard', name: '仪表板', path: '/dashboard', closable: false }],
  activeTab: 'dashboard',
  sidebarCollapsed: false,
  loading: false,
  error: null
};

// Action类型
const ActionTypes = {
  SET_LOADING: 'SET_LOADING',
  SET_ERROR: 'SET_ERROR',
  SET_MENUS: 'SET_MENUS',
  SET_USER_MENUS: 'SET_USER_MENUS',
  ADD_TAB: 'ADD_TAB',
  REMOVE_TAB: 'REMOVE_TAB',
  SET_ACTIVE_TAB: 'SET_ACTIVE_TAB',
  CLOSE_OTHER_TABS: 'CLOSE_OTHER_TABS',
  CLOSE_ALL_TABS: 'CLOSE_ALL_TABS',
  TOGGLE_SIDEBAR: 'TOGGLE_SIDEBAR',
  REFRESH_TAB: 'REFRESH_TAB'
};

// Reducer
const menuReducer = (state, action) => {
  switch (action.type) {
    case ActionTypes.SET_LOADING:
      return { ...state, loading: action.payload };
    
    case ActionTypes.SET_ERROR:
      return { ...state, error: action.payload, loading: false };
    
    case ActionTypes.SET_MENUS:
      return { ...state, menus: action.payload, loading: false };
    
    case ActionTypes.SET_USER_MENUS:
      return { ...state, userMenus: action.payload, loading: false };
    
    case ActionTypes.ADD_TAB:
      const existingTab = state.openTabs.find(tab => tab.id === action.payload.id);
      if (existingTab) {
        return { ...state, activeTab: action.payload.id };
      }
      return {
        ...state,
        openTabs: [...state.openTabs, action.payload],
        activeTab: action.payload.id
      };
    
    case ActionTypes.REMOVE_TAB:
      const newTabs = state.openTabs.filter(tab => tab.id !== action.payload);
      let newActiveTab = state.activeTab;
      
      if (state.activeTab === action.payload && newTabs.length > 0) {
        const removedIndex = state.openTabs.findIndex(tab => tab.id === action.payload);
        newActiveTab = newTabs[Math.max(0, removedIndex - 1)]?.id || newTabs[0]?.id;
      }
      
      return {
        ...state,
        openTabs: newTabs,
        activeTab: newActiveTab
      };
    
    case ActionTypes.SET_ACTIVE_TAB:
      return { ...state, activeTab: action.payload };
    
    case ActionTypes.CLOSE_OTHER_TABS:
      const keepTabs = state.openTabs.filter(tab => 
        tab.id === action.payload || !tab.closable
      );
      return {
        ...state,
        openTabs: keepTabs,
        activeTab: action.payload
      };
    
    case ActionTypes.CLOSE_ALL_TABS:
      const nonClosableTabs = state.openTabs.filter(tab => !tab.closable);
      return {
        ...state,
        openTabs: nonClosableTabs,
        activeTab: nonClosableTabs[0]?.id || 'dashboard'
      };
    
    case ActionTypes.TOGGLE_SIDEBAR:
      return { ...state, sidebarCollapsed: !state.sidebarCollapsed };
    
    case ActionTypes.REFRESH_TAB:
      return { ...state, refreshKey: Date.now() };
    
    default:
      return state;
  }
};

// Context
const MenuContext = createContext();

// Provider组件
export const MenuProvider = ({ children }) => {
  const [state, dispatch] = useReducer(menuReducer, initialState);

  // 获取所有菜单
  const fetchMenus = async () => {
    try {
      dispatch({ type: ActionTypes.SET_LOADING, payload: true });
      const response = await menuApi.getMenuList();
      dispatch({ type: ActionTypes.SET_MENUS, payload: response.data.menus || [] });
    } catch (error) {
      dispatch({ type: ActionTypes.SET_ERROR, payload: error.message });
    }
  };

  // 获取用户权限菜单
  const fetchUserMenus = async () => {
    try {
      dispatch({ type: ActionTypes.SET_LOADING, payload: true });
      const response = await menuApi.getUserMenus();
      const menus = response.data.data || [];
      dispatch({ type: ActionTypes.SET_USER_MENUS, payload: menus });
    } catch (error) {
      dispatch({ type: ActionTypes.SET_ERROR, payload: error.message });
    }
  };

  // 打开新标签页
  const openTab = (menu) => {
    // 只有非菜单组的项目才能打开标签页
    if (menu.type !== 'group') {
      const tab = {
        id: menu.id,
        name: menu.title || menu.name,
        path: menu.path,
        closable: true
      };
      dispatch({ type: ActionTypes.ADD_TAB, payload: tab });
    }
  };

  // 关闭标签页
  const closeTab = (tabId) => {
    dispatch({ type: ActionTypes.REMOVE_TAB, payload: tabId });
  };

  // 设置活动标签页
  const setActiveTab = (tabId) => {
    dispatch({ type: ActionTypes.SET_ACTIVE_TAB, payload: tabId });
  };

  // 关闭其他标签页
  const closeOtherTabs = (tabId) => {
    dispatch({ type: ActionTypes.CLOSE_OTHER_TABS, payload: tabId });
  };

  // 关闭所有标签页
  const closeAllTabs = () => {
    dispatch({ type: ActionTypes.CLOSE_ALL_TABS });
  };

  // 切换侧边栏
  const toggleSidebar = () => {
    dispatch({ type: ActionTypes.TOGGLE_SIDEBAR });
  };



  // 刷新标签页
  const refreshTab = () => {
    dispatch({ type: ActionTypes.REFRESH_TAB });
  };

  // 过滤菜单（根据权限）
  const getFilteredMenus = () => {
    return state.userMenus.length > 0 ? state.userMenus : state.menus;
  };

  // 构建菜单树
  const buildMenuTree = (menus) => {
    const menuMap = new Map();
    const roots = [];

    // 创建菜单映射
    menus.forEach(menu => {
      menuMap.set(menu.id, { ...menu, children: [] });
    });

    // 构建树结构
    menus.forEach(menu => {
      const menuNode = menuMap.get(menu.id);
      if (menu.parent_id && menuMap.has(menu.parent_id)) {
        menuMap.get(menu.parent_id).children.push(menuNode);
      } else {
        roots.push(menuNode);
      }
    });

    return roots.sort((a, b) => (a.sort || 0) - (b.sort || 0));
  };

  const value = {
    ...state,
    fetchMenus,
    fetchUserMenus,
    openTab,
    closeTab,
    setActiveTab,
    closeOtherTabs,
    closeAllTabs,
    toggleSidebar,
    refreshTab,
    getFilteredMenus,
    buildMenuTree
  };

  return (
    <MenuContext.Provider value={value}>
      {children}
    </MenuContext.Provider>
  );
};

// Hook
export const useMenu = () => {
  const context = useContext(MenuContext);
  if (!context) {
    throw new Error('useMenu must be used within a MenuProvider');
  }
  return context;
};

export default MenuContext;