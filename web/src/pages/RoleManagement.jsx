import React, { useState, useEffect } from 'react';
import '../assets/styles/management.css'; // 引入样式文件
import { roleApi } from '../api'; // 引入API
import Modal from 'react-modal'; // 引入模态框组件

const RoleManagement = () => {
  const [roles, setRoles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [currentRole, setCurrentRole] = useState(null); // 用于存储当前编辑或创建的角色

  useEffect(() => {
    // 从API获取角色数据
    const fetchRoles = async () => {
      try {
        setLoading(true);
        const response = await roleApi.getRoleList();
        setRoles(response.data.roles);
      } catch (error) {
        console.error('获取角色数据失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRoles();
  }, []);

  // 创建角色
  const handleCreateRole = async (roleData) => {
    try {
      console.log('发送角色数据:', roleData);
      const response = await roleApi.createRole(roleData);
      console.log('创建角色成功:', response);
      // 重新获取角色列表
      const fetchResponse = await roleApi.getRoleList();
      setRoles(fetchResponse.data.roles);
    } catch (error) {
      console.error('创建角色失败:', error);
      if (error.response) {
        console.error('错误响应:', error.response.data);
        console.error('错误状态:', error.response.status);
      }
    }
  };

  // 更新角色
  const handleUpdateRole = async (id, roleData) => {
    try {
      const response = await roleApi.updateRole(id, roleData);
      // 重新获取角色列表
      const fetchResponse = await roleApi.getRoleList();
      setRoles(fetchResponse.data.roles);
    } catch (error) {
      console.error('更新角色失败:', error);
    }
  };

  // 删除角色
  const handleDeleteRole = async (id) => {
    try {
      const response = await roleApi.deleteRole(id);
      // 重新获取角色列表
      const fetchResponse = await roleApi.getRoleList();
      setRoles(fetchResponse.data.roles);
    } catch (error) {
      console.error('删除角色失败:', error);
    }
  };

  const handleCreateClick = () => {
    setCurrentRole(null); // 创建时清空当前角色
    setModalIsOpen(true);
  };

  return (
    <div className="management-container">
      <h2>角色管理</h2>
      <button 
        style={{ marginBottom: '10px', padding: '5px 10px', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
        onClick={handleCreateClick}
      >
        创建角色
      </button>
      {loading ? (
        <p className="loading">加载中...</p>
      ) : (
        <table className="management-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>角色名称</th>
              <th>描述</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {roles.map(role => (
              <tr key={role.id}>
                <td>{role.id}</td>
                <td>{role.name}</td>
                <td>{role.description}</td>
                <td>{role.status === 1 ? '启用' : '禁用'}</td>
                <td>
                  <button 
                    style={{ marginRight: '5px', padding: '5px 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => {
                      setCurrentRole(role);
                      setModalIsOpen(true);
                    }}
                  >
                    编辑
                  </button>
                  <button 
                    style={{ padding: '5px 10px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => handleDeleteRole(role.id)}
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
        onRequestClose={() => setModalIsOpen(false)}
        contentLabel="角色模态框"
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
        <h2>{currentRole ? '编辑角色' : '创建角色'}</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.target);
          const roleData = {
            name: formData.get('name'),
            description: formData.get('description'),
            status: parseInt(formData.get('status')) || 1,
            tenant_id: 1 // 默认租户ID
          };
          
          if (currentRole) {
            // 更新角色
            handleUpdateRole(currentRole.id, roleData);
          } else {
            // 创建角色
            handleCreateRole(roleData);
          }
          setModalIsOpen(false);
        }}>
          <div style={{ marginBottom: '15px' }}>
            <label>角色名称: </label>
            <input 
              type="text" 
              name="name" 
              defaultValue={currentRole?.name || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>描述: </label>
            <input 
              type="text" 
              name="description" 
              defaultValue={currentRole?.description || ''} 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>状态: </label>
            <select 
              name="status" 
              defaultValue={currentRole?.status || 1} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            >
              <option value={1}>启用</option>
              <option value={2}>禁用</option>
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

export default RoleManagement;