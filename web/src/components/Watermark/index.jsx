import React, { useEffect, useRef } from 'react';
import './Watermark.css';

const Watermark = ({ 
  text = '内部资料 禁止外传', 
  opacity = 0.1, 
  fontSize = 16,
  color = '#000000',
  rotate = -30,
  zIndex = 9999,
  gap = 100,
  offsetLeft = 0,
  offsetTop = 0,
  width = 200,
  height = 200,
  getContainer = () => document.body
}) => {
  const watermarkRef = useRef(null);
  const mutationObserverRef = useRef(null);

  useEffect(() => {
    const container = getContainer();
    
    const createWatermark = () => {
      // 创建canvas生成水印图片
      const canvas = document.createElement('canvas');
      const ctx = canvas.getContext('2d');
      
      canvas.width = width;
      canvas.height = height;
      
      ctx.rotate(rotate * Math.PI / 180);
      ctx.font = `${fontSize}px Arial`;
      ctx.fillStyle = color;
      ctx.textAlign = 'center';
      ctx.textBaseline = 'middle';
      ctx.globalAlpha = opacity;
      
      // 绘制水印文字
      ctx.fillText(text, width / 2, height / 2);
      
      return canvas.toDataURL();
    };

    const renderWatermark = () => {
      // 移除旧的水印
      const oldWatermark = container.querySelector('.watermark-container');
      if (oldWatermark) {
        oldWatermark.remove();
      }

      // 创建新的水印容器
      const watermarkContainer = document.createElement('div');
      watermarkContainer.className = 'watermark-container';
      watermarkContainer.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        width: 100%;
        height: 100%;
        pointer-events: none;
        z-index: ${zIndex};
        background-image: url(${createWatermark()});
        background-repeat: repeat;
        background-position: ${offsetLeft}px ${offsetTop}px;
        background-size: ${gap}px ${gap}px;
        user-select: none;
        -webkit-user-select: none;
        -moz-user-select: none;
        -ms-user-select: none;
      `;

      container.appendChild(watermarkContainer);
      watermarkRef.current = watermarkContainer;
    };

    const observeMutations = () => {
      if (mutationObserverRef.current) {
        mutationObserverRef.current.disconnect();
      }

      mutationObserverRef.current = new MutationObserver((mutations) => {
        mutations.forEach((mutation) => {
          if (mutation.type === 'childList') {
            mutation.removedNodes.forEach((node) => {
              if (node === watermarkRef.current) {
                renderWatermark();
              }
            });
          }
        });
      });

      mutationObserverRef.current.observe(container, {
        childList: true,
        subtree: true
      });
    };

    // 初始渲染
    renderWatermark();
    observeMutations();

    // 监听窗口大小变化
    const handleResize = () => {
      renderWatermark();
    };
    window.addEventListener('resize', handleResize);

    return () => {
      if (mutationObserverRef.current) {
        mutationObserverRef.current.disconnect();
      }
      if (watermarkRef.current) {
        watermarkRef.current.remove();
      }
      window.removeEventListener('resize', handleResize);
    };
  }, [text, opacity, fontSize, color, rotate, zIndex, gap, offsetLeft, offsetTop, width, height, getContainer]);

  return null;
};

export default Watermark;