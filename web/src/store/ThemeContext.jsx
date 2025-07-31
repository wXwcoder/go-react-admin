import React, { createContext, useContext, useState, useEffect } from 'react';

// 创建主题上下文
const ThemeContext = createContext();

// 主题配置
export const themes = {
  light: {
    name: 'light',
    '--bg-primary': '#ffffff',
    '--bg-secondary': '#f5f5f5',
    '--bg-tertiary': '#fafafa',
    '--text-primary': '#262626',
    '--text-secondary': '#595959',
    '--text-tertiary': '#8c8c8c',
    '--border-primary': '#d9d9d9',
    '--border-secondary': '#e8e8e8',
    '--primary-color': '#1890ff',
    '--primary-hover': '#40a9ff',
    '--success-color': '#52c41a',
    '--warning-color': '#faad14',
    '--error-color': '#ff4d4f',
    '--sidebar-bg': '#001529',
    '--sidebar-text': '#ffffff',
    '--header-bg': '#ffffff',
    '--header-text': '#262626',
    '--card-bg': '#ffffff',
    '--shadow': '0 2px 8px rgba(0, 0, 0, 0.15)',
  },
  dark: {
    name: 'dark',
    '--bg-primary': '#141414',
    '--bg-secondary': '#1f1f1f',
    '--bg-tertiary': '#262626',
    '--text-primary': '#ffffff',
    '--text-secondary': '#d9d9d9',
    '--text-tertiary': '#bfbfbf',
    '--border-primary': '#434343',
    '--border-secondary': '#303030',
    '--primary-color': '#177ddc',
    '--primary-hover': '#3c9ae8',
    '--success-color': '#49aa19',
    '--warning-color': '#d89614',
    '--error-color': '#dc4446',
    '--sidebar-bg': '#0c0c0c',
    '--sidebar-text': '#ffffff',
    '--header-bg': '#1f1f1f',
    '--header-text': '#ffffff',
    '--card-bg': '#1f1f1f',
    '--shadow': '0 2px 8px rgba(0, 0, 0, 0.45)',
  }
};

// 主题提供者组件
export const ThemeProvider = ({ children }) => {
  // 从localStorage获取保存的主题，默认使用light主题
  const [currentTheme, setCurrentTheme] = useState(() => {
    const savedTheme = localStorage.getItem('theme');
    return savedTheme && themes[savedTheme] ? savedTheme : 'light';
  });

  // 应用主题到DOM
  const applyTheme = (themeName) => {
    const theme = themes[themeName];
    if (!theme) return;

    const root = document.documentElement;
    Object.entries(theme).forEach(([property, value]) => {
      if (property !== 'name') {
        root.style.setProperty(property, value);
      }
    });

    // 设置data-theme属性用于CSS选择器
    root.setAttribute('data-theme', themeName);

    // 添加主题类名到body
    document.body.className = document.body.className
      .replace(/theme-\w+/g, '')
      .trim() + ` theme-${themeName}`;
  };

  // 切换主题
  const toggleTheme = () => {
    const newTheme = currentTheme === 'light' ? 'dark' : 'light';
    setCurrentTheme(newTheme);
    localStorage.setItem('theme', newTheme);
  };

  // 设置指定主题
  const setTheme = (themeName) => {
    if (themes[themeName]) {
      setCurrentTheme(themeName);
      localStorage.setItem('theme', themeName);
    }
  };

  // 初始化主题
  useEffect(() => {
    applyTheme(currentTheme);
  }, [currentTheme]);

  // 监听系统主题变化
  useEffect(() => {
    const mediaQuery = window.matchMedia('(prefers-color-scheme: dark)');
    const handleChange = (e) => {
      const systemTheme = e.matches ? 'dark' : 'light';
      const savedTheme = localStorage.getItem('theme');
      
      // 如果用户没有手动设置主题，跟随系统
      if (!savedTheme) {
        setCurrentTheme(systemTheme);
      }
    };

    mediaQuery.addEventListener('change', handleChange);
    return () => mediaQuery.removeEventListener('change', handleChange);
  }, []);

  const value = {
    currentTheme,
    toggleTheme,
    setTheme,
    themes: Object.keys(themes),
    themeConfig: themes[currentTheme]
  };

  return (
    <ThemeContext.Provider value={value}>
      {children}
    </ThemeContext.Provider>
  );
};

// 自定义Hook使用主题
export const useTheme = () => {
  const context = useContext(ThemeContext);
  if (!context) {
    throw new Error('useTheme必须在ThemeProvider内部使用');
  }
  return context;
};