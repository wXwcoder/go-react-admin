import { useEffect, useRef } from 'react';

const useWatermark = (options = {}) => {
  const {
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
    container = document.body,
    enabled = true
  } = options;

  const watermarkRef = useRef(null);
  const observerRef = useRef(null);

  useEffect(() => {
    if (!enabled) {
      // 如果禁用水印，移除现有的水印
      if (watermarkRef.current) {
        watermarkRef.current.remove();
        watermarkRef.current = null;
      }
      if (observerRef.current) {
        observerRef.current.disconnect();
        observerRef.current = null;
      }
      return;
    }

    const createWatermarkCanvas = () => {
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
      
      // 支持多行水印
      const lines = Array.isArray(text) ? text : [text];
      const lineHeight = fontSize * 1.5;
      const startY = (height - (lines.length - 1) * lineHeight) / 2;
      
      lines.forEach((line, index) => {
        ctx.fillText(line, width / 2, startY + index * lineHeight);
      });
      
      return canvas.toDataURL();
    };

    const renderWatermark = () => {
      // 移除旧水印
      if (watermarkRef.current) {
        watermarkRef.current.remove();
      }

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
        background-image: url(${createWatermarkCanvas()});
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
      if (observerRef.current) {
        observerRef.current.disconnect();
      }

      observerRef.current = new MutationObserver((mutations) => {
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

      observerRef.current.observe(container, {
        childList: true,
        subtree: true
      });
    };

    renderWatermark();
    observeMutations();

    const handleResize = () => {
      renderWatermark();
    };
    window.addEventListener('resize', handleResize);

    return () => {
      if (observerRef.current) {
        observerRef.current.disconnect();
      }
      if (watermarkRef.current) {
        watermarkRef.current.remove();
      }
      window.removeEventListener('resize', handleResize);
    };
  }, [text, opacity, fontSize, color, rotate, zIndex, gap, offsetLeft, offsetTop, width, height, container, enabled]);

  return {
    update: (newOptions) => {
      Object.assign(options, newOptions);
    },
    remove: () => {
      if (watermarkRef.current) {
        watermarkRef.current.remove();
        watermarkRef.current = null;
      }
      if (observerRef.current) {
        observerRef.current.disconnect();
        observerRef.current = null;
      }
    }
  };
};

export default useWatermark;