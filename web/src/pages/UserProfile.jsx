import React, { useState, useEffect } from 'react';
import { Form, Input, Button, Avatar, Upload, message, Card, Row, Col } from 'antd';
import { UploadOutlined, UserOutlined } from '@ant-design/icons';
import { userApi } from '../api';

const { TextArea } = Input;

const UserProfile = () => {
  const [form] = Form.useForm();
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(false);
  const [uploadLoading, setUploadLoading] = useState(false);
  const [avatarUrl, setAvatarUrl] = useState('');

  useEffect(() => {
    fetchUserInfo();
  }, []);

  const fetchUserInfo = async () => {
    try {
      setLoading(true);
      const response = await userApi.getCurrentUserInfo();
      const userData = response.data.user;
      setUser(userData);
      setAvatarUrl(userData.avatar || '');
      form.setFieldsValue({
        username: userData.username,
        email: userData.email,
        phone: userData.phone,
        real_name: userData.real_name,
        nickname: userData.nickname,
        bio: userData.bio,
      });
    } catch (error) {
      message.error('获取用户信息失败');
      console.error('获取用户信息失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleUpdateProfile = async (values) => {
    try {
      setLoading(true);
      await userApi.updateUser(user.id, {
        ...values,
        avatar: avatarUrl,
      });
      message.success('个人信息更新成功');
      
      // 更新本地存储的用户信息
      const updatedUser = { ...user, ...values, avatar: avatarUrl };
      setUser(updatedUser);
      localStorage.setItem('userInfo', JSON.stringify(updatedUser));
    } catch (error) {
      message.error('更新个人信息失败');
      console.error('更新个人信息失败:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleAvatarChange = async (info) => {
    if (info.file.status === 'uploading') {
      setUploadLoading(true);
      return;
    }
    if (info.file.status === 'done') {
      try {
        // 这里假设后端有上传头像的接口
        const response = await userApi.uploadAvatar(info.file.originFileObj);
        const avatarUrl = response.data.avatar_url;
        setAvatarUrl(avatarUrl);
        message.success('头像上传成功');
      } catch (error) {
        message.error('头像上传失败');
        console.error('头像上传失败:', error);
      } finally {
        setUploadLoading(false);
      }
    }
  };

  const beforeUpload = (file) => {
    const isJpgOrPng = file.type === 'image/jpeg' || file.type === 'image/png';
    if (!isJpgOrPng) {
      message.error('只能上传JPG/PNG格式的图片!');
    }
    const isLt2M = file.size / 1024 / 1024 < 2;
    if (!isLt2M) {
      message.error('图片大小不能超过2MB!');
    }
    return isJpgOrPng && isLt2M;
  };

  return (
    <div style={{ padding: '24px' }}>
      <Card title="个人信息" style={{ maxWidth: 800, margin: '0 auto' }}>
        <Row gutter={24}>
          <Col span={8} style={{ textAlign: 'center' }}>
            <div style={{ marginBottom: 24 }}>
              <Upload
                name="avatar"
                showUploadList={false}
                customRequest={handleAvatarChange}
                beforeUpload={beforeUpload}
              >
                <Avatar
                  size={120}
                  src={avatarUrl}
                  icon={!avatarUrl && <UserOutlined />}
                  style={{ cursor: 'pointer' }}
                />
              </Upload>
              <div style={{ marginTop: 8 }}>
                <Button
                  icon={<UploadOutlined />}
                  loading={uploadLoading}
                  onClick={() => document.querySelector('input[type="file"]').click()}
                >
                  更换头像
                </Button>
              </div>
            </div>
            <div>
              <h3>{user?.real_name || user?.username}</h3>
              <p style={{ color: '#666' }}>{user?.email}</p>
            </div>
          </Col>
          <Col span={16}>
            <Form
              form={form}
              layout="vertical"
              onFinish={handleUpdateProfile}
              initialValues={{
                username: '',
                email: '',
                phone: '',
                real_name: '',
                nickname: '',
                bio: '',
              }}
            >
              <Form.Item
                label="用户名"
                name="username"
                rules={[{ required: true, message: '请输入用户名' }]}
              >
                <Input placeholder="请输入用户名" />
              </Form.Item>

              <Form.Item
                label="邮箱"
                name="email"
                rules={[
                  { required: true, message: '请输入邮箱' },
                  { type: 'email', message: '请输入有效的邮箱地址' }
                ]}
              >
                <Input placeholder="请输入邮箱" />
              </Form.Item>

              <Form.Item
                label="手机号"
                name="phone"
                rules={[{ pattern: /^1[3-9]\d{9}$/, message: '请输入有效的手机号' }]}
              >
                <Input placeholder="请输入手机号" />
              </Form.Item>

              <Form.Item
                label="真实姓名"
                name="real_name"
              >
                <Input placeholder="请输入真实姓名" />
              </Form.Item>

              <Form.Item
                label="昵称"
                name="nickname"
              >
                <Input placeholder="请输入昵称" />
              </Form.Item>

              <Form.Item
                label="个人简介"
                name="bio"
              >
                <TextArea rows={4} placeholder="请输入个人简介" maxLength={200} />
              </Form.Item>

              <Form.Item>
                <Button type="primary" htmlType="submit" loading={loading}>
                  保存修改
                </Button>
              </Form.Item>
            </Form>
          </Col>
        </Row>
      </Card>
    </div>
  );
};

export default UserProfile;