import { message } from 'antd';

// 配置全局message设置
message.config({
  top: 24,
  duration: 2,
  maxCount: 3,
  rtl: false,
  prefixCls: 'ant-message',
});

// 自定义message工具
const messageUtils = {
  success: (content, duration, onClose) => {
    return message.success({
      content,
      duration,
      onClose,
      icon: <span className="anticon anticon-check-circle" style={{ color: '#52c41a' }}>
        <svg viewBox="64 64 896 896" focusable="false" data-icon="check-circle" width="1em" height="1em" fill="currentColor" aria-hidden="true">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm193.5 301.7l-210.6 292a31.8 31.8 0 01-51.7 0L318.5 484.9c-3.8-5.3 0-12.7 6.5-12.7h46.9c10.2 0 19.9 4.9 25.9 13.3l71.2 98.8 157.2-218c6-8.3 15.6-13.3 25.9-13.3H699c6.5 0 10.3 7.4 6.5 12.7z"></path>
        </svg>
      </span>
    });
  },
  
  error: (content, duration, onClose) => {
    return message.error({
      content,
      duration,
      onClose,
      icon: <span className="anticon anticon-close-circle" style={{ color: '#ff4d4f' }}>
        <svg viewBox="64 64 896 896" focusable="false" data-icon="close-circle" width="1em" height="1em" fill="currentColor" aria-hidden="true">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm165.4 618.2l-66-.3L512 563.4l-99.3 118.4-66.1.3c-4.4 0-8-3.5-8-8 0-1.9.7-3.7 1.9-5.2l130.1-155L340.5 338c-1.2-1.5-1.9-3.3-1.9-5.2 0-4.4 3.6-8 8-8l66.1.3L512 464.6l99.3-118.4 66-.3c4.4 0 8 3.5 8 8 0 1.9-.7 3.7-1.9 5.2L553.5 514l130 155c1.2 1.5 1.9 3.3 1.9 5.2 0 4.4-3.6 8-8 8z"></path>
        </svg>
      </span>
    });
  },
  
  warning: (content, duration, onClose) => {
    return message.warning({
      content,
      duration,
      onClose,
      icon: <span className="anticon anticon-exclamation-circle" style={{ color: '#faad14' }}>
        <svg viewBox="64 64 896 896" focusable="false" data-icon="exclamation-circle" width="1em" height="1em" fill="currentColor" aria-hidden="true">
          <path d="M512 64C264.6 64 64 264.6 64 512s200.6 448 448 448 448-200.6 448-448S759.4 64 512 64zm-32 232c0-4.4 3.6-8 8-8h48c4.4 0 8 3.6 8 8v272c0 4.4-3.6 8-8 8h-48c-4.4 0-8-3.6-8-8V296zm32 440a48.01 48.01 0 010-96 48.01 48.01 0 010 96z"></path>
        </svg>
      </span>
    });
  },
  
  info: (content, duration, onClose) => {
    return message.info({
      content,
      duration,
      onClose
    });
  },
  
  loading: (content, duration, onClose) => {
    return message.loading({
      content,
      duration,
      onClose
    });
  }
};

export default messageUtils;

// 导出原始的message对象以便需要时使用
export { message as originalMessage };