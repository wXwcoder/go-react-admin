# 前端水印功能文档

## 功能概述

本项目为前端页面添加了全面的水印功能，用于保护用户数据安全，防止敏感信息泄露。

## 功能特点

- **动态用户信息水印**：根据当前登录用户显示个性化水印
- **防删除保护**：使用MutationObserver监控水印元素，防止被恶意删除
- **响应式设计**：窗口大小变化时自动重新渲染水印
- **可配置参数**：支持透明度、字体大小、颜色、旋转角度、间距等自定义设置
- **用户友好的设置界面**：通过模态框进行水印配置

## 技术实现

### 核心组件

1. **Watermark组件** (`src/components/Watermark/index.jsx`)
   - 使用Canvas生成水印图片
   - 支持自定义文字、样式、位置等参数
   - 内置防删除机制

2. **useWatermark Hook** (`src/hooks/useWatermark.js`)
   - 提供简洁的水印使用方式
   - 支持动态配置更新
   - 自动处理清理和重新渲染

3. **WatermarkSettings组件** (`src/components/WatermarkSettings.jsx`)
   - 提供图形化配置界面
   - 支持实时预览效果
   - 配置数据持久化到localStorage

### 使用方法

#### 基础使用

在需要添加水印的组件中使用useWatermark hook：

```javascript
import useWatermark from './hooks/useWatermark';

const MyComponent = () => {
  useWatermark({
    text: '机密文档',
    opacity: 0.1,
    fontSize: 16,
    color: '#000000',
    rotate: -30,
    gap: 100
  });

  return <div>内容</div>;
};
```

#### 动态用户信息

自动获取当前登录用户信息作为水印：

```javascript
useWatermark({
  text: `${username} - ${realName}`,
  userInfo: true
});
```

#### 可配置水印

通过设置界面配置水印参数：

1. 点击顶部导航栏的齿轮图标打开水印设置
2. 调整各项参数
3. 点击确定保存配置

## 配置选项

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| enabled | boolean | true | 是否启用水印 |
| text | string | '内部资料 禁止外传' | 水印文字内容 |
| opacity | number | 0.08 | 透明度 (0.01-0.5) |
| fontSize | number | 14 | 字体大小 (10-30px) |
| color | string | '#000000' | 文字颜色 |
| rotate | number | -30 | 旋转角度 (-90°到90°) |
| gap | number | 150 | 水印间距 (50-300px) |
| userInfo | boolean | true | 是否显示用户信息 |

## 安全特性

1. **防删除机制**
   - 使用MutationObserver监控DOM变化
   - 水印被删除时自动重新创建
   - 防止通过开发者工具删除水印

2. **防篡改机制**
   - 水印样式设置为`pointer-events: none`
   - 禁止用户选择和交互
   - 使用固定定位确保覆盖整个页面

3. **数据保护**
   - 水印内容包含用户信息，便于追踪泄露源
   - 敏感信息通过水印标识

## 样式定制

可以通过CSS进一步定制水印样式：

```css
.watermark-container {
  /* 自定义水印容器样式 */
}

/* 防止水印被删除 */
.watermark-container[data-watermark] {
  display: block !important;
  visibility: visible !important;
}
```

## 注意事项

1. 水印功能仅在用户登录后生效
2. 配置信息保存在浏览器localStorage中
3. 刷新页面后配置仍然有效
4. 不同用户的水印配置相互独立

## 故障排除

### 水印不显示
- 检查是否启用了水印功能
- 确认用户已登录
- 查看浏览器控制台是否有错误信息

### 水印被删除
- 系统会自动重新创建被删除的水印
- 如遇到持续问题，请刷新页面

### 配置不生效
- 检查localStorage中是否正确保存了配置
- 确认配置格式正确
- 尝试恢复默认设置