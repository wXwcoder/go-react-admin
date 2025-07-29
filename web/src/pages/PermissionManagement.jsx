import React, { useState, useEffect } from 'react';
import '../assets/styles/management.css';
import { permissionApi, roleApi, userApi, menuApi, apiApi } from '../api';
import Modal from 'react-modal';

const PermissionManagement = () => {
  const [activeTab, setActiveTab] = useState('rolePermissions');
  const [roles, setRoles] = useState([]);
  const [users, setUsers] = useState([]);
  const [menus, setMenus] = useState([]);
  const [apis, setApis] = useState([]);
  const [loading, setLoading] = useState(false);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [currentItem, setCurrentItem] = useState(null);
  const [selectedPermissions, setSelectedPermissions] = useState({
    menus: [],
    apis: []
  });

  useEffect(() => {
    fetchInitialData();
  }, []);

  const fetchInitialData = async () => {
    try {
      setLoading(true);
      const [rolesRes, usersRes, menusRes, apisRes] = await Promise.all([
        roleApi.getRoleList(),
        userApi.getUserList(),
        menuApi.getMenuList(),
        apiApi.getApiList()
      ]);
      
      setRoles(rolesRes.data.roles || []);
      setUsers(usersRes.data.users || []);
      setMenus(menusRes.data.menus || []);
      setApis(apisRes.data.apis || []);
    } catch (error) {
      console.error('获取数据失败:', error);
    } finally {
      setLoading(false);
    }
  };

  // 获取角色权限
  const handleGetRolePermissions = async (roleId) => {
    try {
      const response = await permissionApi.getRolePermissions(roleId);
      const permissions = response.data.data;
      setCurrentItem({ type: 'role', id: roleId, data: permissions });
      setSelectedPermissions({
        menus: permissions.menus?.map(m => m.id) || [],
        apis: permissions.apis?.map(a => a.id) || []
      });
      setModalIsOpen(true);
    } catch (error) {
      console.error('获取角色权限失败:', error);
    }
  };

  // 获取用户角色
  const handleGetUserRoles = async (userId) => {
    try {
      const response = await permissionApi.getUserRoles(userId);
      const userRoles = response.data.data;
      setCurrentItem({ type: 'user', id: userId, data: userRoles });
      setSelectedPermissions({
        roles: userRoles.roles?.map(r => r.id) || []
      });
      setModalIsOpen(true);
    } catch (error) {
      console.error('获取用户角色失败:', error);
    }
  };

  // 保存角色权限
  const handleSaveRolePermissions = async () => {
    try {
      const data = {
        role_id: currentItem.id,
        menu_ids: selectedPermissions.menus,
        api_ids: selectedPermissions.apis
      };
      await permissionApi.assignRolePermissions(data);
      alert('权限分配成功');
      setModalIsOpen(false);
    } catch (error) {
      console.error('分配权限失败:', error);
      alert('权限分配失败');
    }
  };

  // 保存用户角色
  const handleSaveUserRoles = async () => {
    try {
      const data = {
        user_id: currentItem.id,
        role_ids: selectedPermissions.roles
      };
      await permissionApi.assignUserRoles(data);
      alert('角色分配成功');
      setModalIsOpen(false);
    } catch (error) {
      console.error('分配角色失败:', error);
      alert('角色分配失败');
    }
  };

  // 处理权限选择
  const handlePermissionChange = (type, id, checked) => {
    setSelectedPermissions(prev => ({
      ...prev,
      [type]: checked 
        ? [...prev[type], id]
        : prev[type].filter(item => item !== id)
    }));
  };

  const renderRolePermissions = () => (
    <div>
      <h3>角色权限管理</h3>
      {loading ? (
        <p className="loading">加载中...</p>
      ) : (
        <table className="management-table">
          <thead>
            <tr>
              <th>角色ID</th>
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
                    style={{ 
                      padding: '5px 10px', 
                      backgroundColor: '#007bff', 
                      color: 'white', 
                      border: 'none', 
                      borderRadius: '4px', 
                      cursor: 'pointer' 
                    }}
                    onClick={() => handleGetRolePermissions(role.id)}
                  >
                    配置权限
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );

  const renderUserRoles = () => (
    <div>
      <h3>用户角色管理</h3>
      {loading ? (
        <p className="loading">加载中...</p>
      ) : (
        <table className="management-table">
          <thead>
            <tr>
              <th>用户ID</th>
              <th>用户名</th>
              <th>昵称</th>
              <th>邮箱</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {users.map(user => (
              <tr key={user.id}>
                <td>{user.id}</td>
                <td>{user.username}</td>
                <td>{user.nickname}</td>
                <td>{user.email}</td>
                <td>{user.status === 1 ? '启用' : '禁用'}</td>
                <td>
                  <button 
                    style={{ 
                      padding: '5px 10px', 
                      backgroundColor: '#28a745', 
                      color: 'white', 
                      border: 'none', 
                      borderRadius: '4px', 
                      cursor: 'pointer' 
                    }}
                    onClick={() => handleGetUserRoles(user.id)}
                  >
                    分配角色
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  );

  const renderPermissionModal = () => {
    if (!currentItem) return null;

    if (currentItem.type === 'role') {
      return (
        <div>
          <h3>配置角色权限 - {currentItem.data.role?.name}</h3>
          
          <div style={{ marginBottom: '20px' }}>
            <h4>菜单权限</h4>
            <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {menus.map(menu => (
                <div key={menu.id} style={{ marginBottom: '5px' }}>
                  <label>
                    <input
                      type="checkbox"
                      checked={selectedPermissions.menus.includes(menu.id)}
                      onChange={(e) => handlePermissionChange('menus', menu.id, e.target.checked)}
                    />
                    <span style={{ marginLeft: '5px' }}>{menu.title} ({menu.path})</span>
                  </label>
                </div>
              ))}
            </div>
          </div>

          <div style={{ marginBottom: '20px' }}>
            <h4>API权限</h4>
            <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {apis.map(api => (
                <div key={api.id} style={{ marginBottom: '5px' }}>
                  <label>
                    <input
                      type="checkbox"
                      checked={selectedPermissions.apis.includes(api.id)}
                      onChange={(e) => handlePermissionChange('apis', api.id, e.target.checked)}
                    />
                    <span style={{ marginLeft: '5px' }}>{api.method} {api.path} - {api.description}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>

          <div>
            <button 
              onClick={handleSaveRolePermissions}
              style={{ 
                padding: '10px 20px', 
                backgroundColor: '#007bff', 
                color: 'white', 
                border: 'none', 
                borderRadius: '4px', 
                cursor: 'pointer',
                marginRight: '10px'
              }}
            >
              保存
            </button>
            <button 
              onClick={() => setModalIsOpen(false)}
              style={{ 
                padding: '10px 20px', 
                backgroundColor: '#6c757d', 
                color: 'white', 
                border: 'none', 
                borderRadius: '4px', 
                cursor: 'pointer' 
              }}
            >
              取消
            </button>
          </div>
        </div>
      );
    }

    if (currentItem.type === 'user') {
      return (
        <div>
          <h3>分配用户角色 - {currentItem.data.user?.username}</h3>
          
          <div style={{ marginBottom: '20px' }}>
            <h4>选择角色</h4>
            <div style={{ maxHeight: '200px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {roles.map(role => (
                <div key={role.id} style={{ marginBottom: '5px' }}>
                  <label>
                    <input
                      type="checkbox"
                      checked={selectedPermissions.roles?.includes(role.id) || false}
                      onChange={(e) => handlePermissionChange('roles', role.id, e.target.checked)}
                    />
                    <span style={{ marginLeft: '5px' }}>{role.name} - {role.description}</span>
                  </label>
                </div>
              ))}
            </div>
          </div>

          <div>
            <button 
              onClick={handleSaveUserRoles}
              style={{ 
                padding: '10px 20px', 
                backgroundColor: '#007bff', 
                color: 'white', 
                border: 'none', 
                borderRadius: '4px', 
                cursor: 'pointer',
                marginRight: '10px'
              }}
            >
              保存
            </button>
            <button 
              onClick={() => setModalIsOpen(false)}
              style={{ 
                padding: '10px 20px', 
                backgroundColor: '#6c757d', 
                color: 'white', 
                border: 'none', 
                borderRadius: '4px', 
                cursor: 'pointer' 
              }}
            >
              取消
            </button>
          </div>
        </div>
      );
    }

    return null;
  };

  return (
    <div className="management-container">
      <h2>权限管理</h2>
      
      {/* 标签页导航 */}
      <div style={{ marginBottom: '20px', borderBottom: '1px solid #ddd' }}>
        <button
          style={{
            padding: '10px 20px',
            backgroundColor: activeTab === 'rolePermissions' ? '#007bff' : '#f8f9fa',
            color: activeTab === 'rolePermissions' ? 'white' : '#333',
            border: 'none',
            borderRadius: '4px 4px 0 0',
            cursor: 'pointer',
            marginRight: '5px'
          }}
          onClick={() => setActiveTab('rolePermissions')}
        >
          角色权限
        </button>
        <button
          style={{
            padding: '10px 20px',
            backgroundColor: activeTab === 'userRoles' ? '#007bff' : '#f8f9fa',
            color: activeTab === 'userRoles' ? 'white' : '#333',
            border: 'none',
            borderRadius: '4px 4px 0 0',
            cursor: 'pointer'
          }}
          onClick={() => setActiveTab('userRoles')}
        >
          用户角色
        </button>
      </div>

      {/* 标签页内容 */}
      {activeTab === 'rolePermissions' && renderRolePermissions()}
      {activeTab === 'userRoles' && renderUserRoles()}

      {/* 权限配置模态框 */}
      <Modal
        isOpen={modalIsOpen}
        onRequestClose={() => setModalIsOpen(false)}
        contentLabel="权限配置"
        style={{
          content: {
            top: '50%',
            left: '50%',
            right: 'auto',
            bottom: 'auto',
            marginRight: '-50%',
            transform: 'translate(-50%, -50%)',
            width: '600px',
            maxHeight: '80vh',
            padding: '20px',
            overflow: 'auto'
          }
        }}
      >
        {renderPermissionModal()}
      </Modal>
    </div>
  );
};

export default PermissionManagement;