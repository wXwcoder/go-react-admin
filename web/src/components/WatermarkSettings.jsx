import React, { useState, useEffect } from 'react';
import { Modal, Form, Input, InputNumber, Switch, Slider, Space, Button } from 'antd';
import '../hooks/useWatermark';

const WatermarkSettings = ({ visible, onClose }) => {
  const [form] = Form.useForm();
  const [settings, setSettings] = useState(() => {
    const saved = localStorage.getItem('watermarkSettings');
    return saved ? JSON.parse(saved) : {
      enabled: true,
      text: '内部资料 禁止外传',
      opacity: 0.08,
      fontSize: 18,
      color: '#000000',
      rotate: -30,
      gap: 150,
      userInfo: true
    };
  });

  const [userInfo, setUserInfo] = useState(null);

  useEffect(() => {
    const userStr = localStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        setUserInfo(user);
      } catch (error) {
        console.error('解析用户信息失败:', error);
      }
    }
  }, []);

  const generateWatermarkText = (values) => {
    if (values.userInfo && userInfo) {
      const username = userInfo.username || '用户';
      const realName = userInfo.realName || '';
      return realName ? `${username} - ${realName}` : username;
    }
    return values.text || '内部资料 禁止外传';
  };

  const handleOk = () => {
    form.validateFields().then((values) => {
      const newSettings = {
        ...values,
        text: generateWatermarkText(values)
      };
      
      setSettings(newSettings);
      localStorage.setItem('watermarkSettings', JSON.stringify(newSettings));
      
      // 更新全局水印配置
      window.dispatchEvent(new CustomEvent('watermarkSettingsChanged', { detail: newSettings }));
      
      onClose();
    });
  };

  const handleReset = () => {
    const defaultSettings = {
      enabled: true,
      text: '内部资料 禁止外传',
      opacity: 0.08,
      fontSize: 14,
      color: '#000000',
      rotate: -30,
      gap: 150,
      userInfo: true
    };
    
    form.setFieldsValue(defaultSettings);
    setSettings(defaultSettings);
    localStorage.setItem('watermarkSettings', JSON.stringify(defaultSettings));
    window.dispatchEvent(new CustomEvent('watermarkSettingsChanged', { detail: defaultSettings }));
  };

  return (
    <Modal
      title="水印设置"
      open={visible}
      onCancel={onClose}
      onOk={handleOk}
      width={600}
      footer={[
        <Button key="reset" onClick={handleReset}>
          恢复默认
        </Button>,
        <Button key="cancel" onClick={onClose}>
          取消
        </Button>,
        <Button key="ok" type="primary" onClick={handleOk}>
          确定
        </Button>
      ]}
    >
      <Form
        form={form}
        layout="vertical"
        initialValues={settings}
      >
        <Form.Item
          label="启用水印"
          name="enabled"
          valuePropName="checked"
        >
          <Switch />
        </Form.Item>

        <Form.Item
          label="显示用户信息"
          name="userInfo"
          valuePropName="checked"
          tooltip="在水印中显示当前登录用户的信息"
        >
          <Switch />
        </Form.Item>

        <Form.Item
          label="自定义文字"
          name="text"
          rules={[{ required: true, message: '请输入水印文字' }]}
        >
          <Input placeholder="请输入水印文字" />
        </Form.Item>

        <Form.Item
          label="透明度"
          name="opacity"
          rules={[{ required: true }]}
        >
          <Slider
            min={0.01}
            max={0.5}
            step={0.01}
            marks={{
              0.01: '0.01',
              0.1: '0.1',
              0.2: '0.2',
              0.3: '0.3',
              0.4: '0.4',
              0.5: '0.5'
            }}
          />
        </Form.Item>

        <Form.Item
          label="字体大小"
          name="fontSize"
          rules={[{ required: true }]}
        >
          <InputNumber min={10} max={30} />
        </Form.Item>

        <Form.Item
          label="颜色"
          name="color"
          rules={[{ required: true }]}
        >
          <Input type="color" />
        </Form.Item>

        <Form.Item
          label="旋转角度"
          name="rotate"
          rules={[{ required: true }]}
        >
          <Slider
            min={-90}
            max={90}
            marks={{
              '-90': '-90°',
              '-45': '-45°',
              0: '0°',
              45: '45°',
              90: '90°'
            }}
          />
        </Form.Item>

        <Form.Item
          label="间距"
          name="gap"
          rules={[{ required: true }]}
        >
          <InputNumber min={50} max={300} />
        </Form.Item>
      </Form>
    </Modal>
  );
};

export default WatermarkSettings;