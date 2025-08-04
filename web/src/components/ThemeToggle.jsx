import React from 'react';
import { useTheme } from '../store/ThemeContext';
import '../assets/styles/ThemeToggle.css';

const ThemeToggle = ({ className = '' }) => {
  const { currentTheme, toggleTheme } = useTheme();

  return (
    <div className={`theme-toggle ${className}`} style={{ whiteSpace: 'nowrap', minWidth: 'max-content' }}>
      <button
        className="theme-toggle-btn"
        onClick={toggleTheme}
        title={`切换到${currentTheme === 'light' ? '暗色' : '亮色'}主题`}
        aria-label={`当前主题: ${currentTheme === 'light' ? '亮色' : '暗色'}，点击切换`}
      >
        <div className={`theme-icon ${currentTheme}`}>
          {currentTheme === 'light' ? (
            // 亮色主题图标 - 太阳
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <circle cx="12" cy="12" r="5"></circle>
              <line x1="12" y1="1" x2="12" y2="3"></line>
              <line x1="12" y1="21" x2="12" y2="23"></line>
              <line x1="4.22" y1="4.22" x2="5.64" y2="5.64"></line>
              <line x1="18.36" y1="18.36" x2="19.78" y2="19.78"></line>
              <line x1="1" y1="12" x2="3" y2="12"></line>
              <line x1="21" y1="12" x2="23" y2="12"></line>
              <line x1="4.22" y1="19.78" x2="5.64" y2="18.36"></line>
              <line x1="18.36" y1="5.64" x2="19.78" y2="4.22"></line>
            </svg>
          ) : (
            // 暗色主题图标 - 月亮
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              strokeLinecap="round"
              strokeLinejoin="round"
            >
              <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z"></path>
            </svg>
          )}
        </div>
        <span className="theme-text">
          {currentTheme === 'light' ? '亮色' : '暗色'}
        </span>
      </button>
    </div>
  );
};

export default ThemeToggle;