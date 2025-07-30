import React, { useState, useEffect } from 'react';
import '../assets/styles/management.css'; // 引入样式文件
import { userApi, roleApi, permissionApi } from '../api'; // 引入API
import Modal from 'react-modal'; // 引入模态框组件

const UserManagement = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [currentUser, setCurrentUser] = useState(null); // 用于存储当前编辑或创建的用户
  const [selectedRoles, setSelectedRoles] = useState([]); // 存储选择的角色
  const [availableRoles, setAvailableRoles] = useState([]); // 可用角色列表

  useEffect(() => {
    fetchUsers();
    fetchRoles();
  }, []);

  // 从API获取用户数据及角色信息
  const fetchUsers = async () => {
    try {
      setLoading(true);
      const response = await userApi.getUserList();
      const usersData = response.data.users || [];
      
      // 获取每个用户的角色信息
      const usersWithRoles = await Promise.all(
        usersData.map(async (user) => {
          try {
            const roleResponse = await permissionApi.getUserRoles(user.id);
            const roles = roleResponse.data.data?.roles || [];
            return {
              ...user,
              roles: roles // 存储角色对象数组
            };
          } catch (error) {
            console.error(`获取用户 ${user.id} 角色失败:`, error);
            return {
              ...user,
              roles: []
            };
          }
        })
      );
      
      setUsers(usersWithRoles);
    } catch (error) {
      console.error('获取用户数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 从API获取角色列表
  const fetchRoles = async () => {
    try {
      const response = await roleApi.getRoleList();
      const roles = response.data.roles || response.data || [];
      setAvailableRoles(roles);
    } catch (error) {
      console.error('获取角色列表失败:', error);
      // 如果获取失败，使用默认角色
      setAvailableRoles([
        { id: 1, name: '管理员' },
        { id: 2, name: '普通用户' },
        { id: 3, name: '编辑' },
        { id: 4, name: '访客' }
      ]);
    }
  };

  // 创建用户并分配角色
  const handleCreateUser = async (userData, roleIds) => {
    try {
      console.log('Sending user data:', userData);
      const createResponse = await userApi.createUser(userData);
      const newUser = createResponse.data.user;
      
      // 分配角色
      if (roleIds && roleIds.length > 0) {
        await permissionApi.assignUserRoles({
          user_id: newUser.id,
          role_ids: roleIds
        });
      }
      
      // 重新获取用户列表
      await fetchUsers();
    } catch (error) {
      console.error('创建用户失败:', error);
      if (error.response) {
        console.error('Error response:', error.response);
      }
    }
  };

  // 更新用户并分配角色
  const handleUpdateUser = async (id, userData, roleIds) => {
    try {
      // 更新用户基本信息
      await userApi.updateUser(id, userData);
      
      // 更新角色
      await permissionApi.assignUserRoles({
        user_id: id,
        role_ids: roleIds || []
      });
      
      // 重新获取用户列表
      await fetchUsers();
    } catch (error) {
      console.error('更新用户失败:', error);
    }
  };

  // 删除用户
  const handleDeleteUser = async (id) => {
    try {
      await userApi.deleteUser(id);
      // 重新获取用户列表
      const fetchResponse = await userApi.getUserList();
      setUsers(fetchResponse.data.users);
    } catch (error) {
      console.error('删除用户失败:', error);
    }
  };

  const handleCreateClick = () => {
    setCurrentUser(null); // 创建时清空当前用户
    setSelectedRoles([]); // 清空选择的角色
    setModalIsOpen(true);
  };

  return (
    <div className="management-container">
      <h2>用户管理</h2>
      <button 
        style={{ marginBottom: '10px', padding: '5px 10px', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
        onClick={handleCreateClick}
      >
        创建用户
      </button>
      {loading ? (
        <p className="loading">加载中...</p>
      ) : (
        <table className="management-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>用户名</th>
              <th>邮箱</th>
              <th>角色</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {users.map(user => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.username}</td>
                <td>{user.email}</td>
                <td>{user.roles ? user.roles.map(r => r.name).join(', ') : ''}</td>
                <td>{user.status === 1 ? '启用' : '禁用'}</td>
                <td>
                  <button 
                    style={{ marginRight: '5px', padding: '5px 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => {
                      setCurrentUser(user);
                      // 设置选中的角色ID
                      const roleIds = user.roles ? user.roles.map(r => r.id) : [];
                      setSelectedRoles(roleIds);
                      setModalIsOpen(true);
                    }}
                  >
                    编辑
                  </button>
                  <button 
                    style={{ padding: '5px 10px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => handleDeleteUser(user.id)}
                  >
                    删除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
      {/* 模态框 */}
      <Modal
        isOpen={modalIsOpen}
        onRequestClose={() => {
          setModalIsOpen(false);
          setSelectedRoles([]); // 清空选择的角色
        }}
        contentLabel="用户模态框"
        style={{
          content: {
            top: '50%',
            left: '50%',
            right: 'auto',
            bottom: 'auto',
            marginRight: '-50%',
            transform: 'translate(-50%, -50%)',
            width: '500px',
            padding: '20px'
          }
        }}
      >
        <h2>{currentUser ? '编辑用户' : '创建用户'}</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.target);
          const userData = {
            username: formData.get('username'),
            email: formData.get('email'),
            status: parseInt(formData.get('status')) || 1
          };
          
          if (currentUser) {
            // 更新用户
            handleUpdateUser(currentUser.id, userData, selectedRoles);
          } else {
            // 创建用户
            handleCreateUser(userData, selectedRoles);
          }
          setModalIsOpen(false);
          setSelectedRoles([]); // 清空选择的角色
        }}>
          <div style={{ marginBottom: '15px' }}>
            <label>用户名: </label>
            <input 
              type="text" 
              name="username" 
              defaultValue={currentUser?.username || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>邮箱: </label>
            <input 
              type="email" 
              name="email" 
              defaultValue={currentUser?.email || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>角色: </label>
            <div style={{ marginTop: '5px' }}>
              {availableRoles.length > 0 ? (
                availableRoles.map(role => (
                  <label key={role.id} style={{ display: 'block', marginBottom: '5px' }}>
                    <input
                      type="checkbox"
                      value={role.id}
                      checked={selectedRoles.includes(role.id)}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedRoles([...selectedRoles, role.id]);
                        } else {
                          setSelectedRoles(selectedRoles.filter(r => r !== role.id));
                        }
                      }}
                      style={{ marginRight: '5px' }}
                    />
                    {role.name}
                  </label>
                ))
              ) : (
                <p style={{ color: '#666', fontStyle: 'italic' }}>正在加载角色列表...</p>
              )}
            </div>
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>状态: </label>
            <select 
              name="status" 
              defaultValue={currentUser?.status || 1} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            >
              <option value={1}>启用</option>
              <option value={0}>禁用</option>
            </select>
          </div>
          <div>
            <button 
              type="submit" 
              style={{ padding: '10px 20px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer', marginRight: '10px' }}
            >
              保存
            </button>
            <button 
              type="button" 
              onClick={() => setModalIsOpen(false)}
              style={{ padding: '10px 20px', backgroundColor: '#6c757d', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
            >
              取消
            </button>
          </div>
        </form>
      </Modal>
    </div>
  );
};

export default UserManagement;