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
      const response = await permissionApi.assignRolePermissions(data);
      
      if (response.data.code === 200) {
        alert('权限分配成功');
        setModalIsOpen(false);
        // 刷新数据
        await fetchInitialData();
      } else {
        alert('权限分配失败: ' + response.data.msg);
      }
    } catch (error) {
      console.error('分配权限失败:', error);
      const errorMsg = error.response?.data?.msg || error.message || '权限分配失败';
      alert('权限分配失败: ' + errorMsg);
    }
  };

  // 保存用户角色
  const handleSaveUserRoles = async () => {
    try {
      const data = {
        user_id: currentItem.id,
        role_ids: selectedPermissions.roles
      };
      const response = await permissionApi.assignUserRoles(data);
      
      if (response.data.code === 200) {
        alert('角色分配成功');
        setModalIsOpen(false);
        // 刷新数据
        await fetchInitialData();
      } else {
        alert('角色分配失败: ' + response.data.msg);
      }
    } catch (error) {
      console.error('分配角色失败:', error);
      const errorMsg = error.response?.data?.msg || error.message || '角色分配失败';
      alert('角色分配失败: ' + errorMsg);
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

  // 处理批量权限选择
  const handleBatchPermissionChange = (type, ids, checked) => {
    setSelectedPermissions(prev => {
      const current = prev[type] || [];
      if (checked) {
        return {
          ...prev,
          [type]: [...new Set([...current, ...ids])]
        };
      } else {
        return {
          ...prev,
          [type]: current.filter(id => !ids.includes(id))
        };
      }
    });
  };

  // 全选/取消全选
  const handleSelectAll = (type, allIds) => {
    const current = selectedPermissions[type] || [];
    const shouldSelectAll = current.length < allIds.length;
    handleBatchPermissionChange(type, allIds, shouldSelectAll);
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

    // 按父节点分组菜单
    const groupMenusByParent = () => {
      const groups = {};
      const rootMenus = menus.filter(menu => !menu.parent_id || menu.parent_id === 0);
      
      rootMenus.forEach(rootMenu => {
        const children = menus.filter(menu => menu.parent_id === rootMenu.id);
        if (children.length > 0) {
          groups[rootMenu.id] = {
            parent: rootMenu,
            children: children
          };
        } else {
          groups[rootMenu.id] = {
            parent: rootMenu,
            children: [rootMenu]
          };
        }
      });

      // 处理没有父节点的独立菜单
      menus.forEach(menu => {
        if (!groups[menu.id] && (!menu.parent_id || menu.parent_id === 0)) {
          groups[menu.id] = {
            parent: menu,
            children: [menu]
          };
        }
      });

      return groups;
    };

    // 按API路径分组
    const groupApisByPath = () => {
      const groups = {};
      
      apis.forEach(api => {
        let groupKey = '其他';
        const path = api.path.toLowerCase();
        
        // 按API路径特征分组
        if (path.includes('/user') || path.includes('/users')) {
          groupKey = '用户管理';
        } else if (path.includes('/role') || path.includes('/roles')) {
          groupKey = '角色管理';
        } else if (path.includes('/menu') || path.includes('/menus')) {
          groupKey = '菜单管理';
        } else if (path.includes('/api') || path.includes('/apis')) {
          groupKey = 'API管理';
        } else if (path.includes('/permission') || path.includes('/permissions')) {
          groupKey = '权限管理';
        } else if (path.includes('/tenant') || path.includes('/tenants')) {
          groupKey = '租户管理';
        } else if (path.includes('/log') || path.includes('/logs')) {
          groupKey = '日志管理';
        } else if (path.includes('/dynamic') || path.includes('/table')) {
          groupKey = '动态数据';
        } else {
          // 提取路径的第一个有意义的部分
          const pathParts = api.path.split('/').filter(part => part && !part.startsWith(':'));
          if (pathParts.length > 0) {
            const firstPart = pathParts[0];
            if (firstPart && firstPart.length > 1) {
              groupKey = firstPart.charAt(0).toUpperCase() + firstPart.slice(1);
            }
          }
        }
        
        if (!groups[groupKey]) {
          groups[groupKey] = [];
        }
        groups[groupKey].push(api);
      });
      
      // 按字母顺序排序分组
      const sortedGroups = {};
      Object.keys(groups).sort().forEach(key => {
        sortedGroups[key] = groups[key];
      });
      
      return sortedGroups;
    };

    if (currentItem.type === 'role') {
      const menuGroups = groupMenusByParent();
      const apiGroups = groupApisByPath();

      return (
        <div>
          <h3>配置角色权限 - {currentItem.data.role?.name}</h3>
          
          <div style={{ marginBottom: '20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
              <h4 style={{ margin: 0 }}>菜单权限</h4>
              <button
                onClick={() => handleSelectAll('menus', menus.map(m => m.id))}
                style={{
                  padding: '5px 10px',
                  backgroundColor: '#17a2b8',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontSize: '12px'
                }}
              >
                {selectedPermissions.menus.length === menus.length ? '取消全选' : '全选'}
              </button>
            </div>
            
            <div style={{ maxHeight: '300px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {Object.entries(menuGroups).map(([groupId, group]) => {
                const groupMenuIds = group.children.map(child => child.id);
                const selectedInGroup = groupMenuIds.filter(id => selectedPermissions.menus.includes(id));
                const isGroupSelected = selectedInGroup.length === groupMenuIds.length;
                
                return (
                  <div key={groupId} style={{ marginBottom: '15px', border: '1px solid #eee', padding: '10px', borderRadius: '4px' }}>
                    <div style={{ marginBottom: '8px', fontWeight: 'bold', backgroundColor: '#f8f9fa', padding: '5px', borderRadius: '3px' }}>
                      <label>
                        <input
                          type="checkbox"
                          checked={isGroupSelected}
                          onChange={(e) => handleBatchPermissionChange('menus', groupMenuIds, e.target.checked)}
                        />
                        <span style={{ marginLeft: '5px' }}>
                          {group.parent.title} {group.children.length > 1 ? `(${group.children.length}项)` : ''}
                        </span>
                      </label>
                    </div>
                    
                    <div style={{ marginLeft: '20px' }}>
                      {group.children.map(menu => (
                        <div key={menu.id} style={{ marginBottom: '3px' }}>
                          <label style={{ fontSize: '14px' }}>
                            <input
                              type="checkbox"
                              checked={selectedPermissions.menus.includes(menu.id)}
                              onChange={(e) => handlePermissionChange('menus', menu.id, e.target.checked)}
                            />
                            <span style={{ marginLeft: '5px' }}>
                              {menu.title} {menu.path && menu.path !== menu.title ? `(${menu.path})` : ''}
                            </span>
                          </label>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
            </div>
          </div>

          <div style={{ marginBottom: '20px' }}>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
              <h4 style={{ margin: 0 }}>API权限</h4>
              <button
                onClick={() => handleSelectAll('apis', apis.map(a => a.id))}
                style={{
                  padding: '5px 10px',
                  backgroundColor: '#17a2b8',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontSize: '12px'
                }}
              >
                {selectedPermissions.apis.length === apis.length ? '取消全选' : '全选'}
              </button>
            </div>
            
            <div style={{ maxHeight: '300px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {Object.entries(apiGroups).map(([groupKey, groupApis]) => {
                const groupApiIds = groupApis.map(api => api.id);
                const selectedInGroup = groupApiIds.filter(id => selectedPermissions.apis.includes(id));
                const isGroupSelected = selectedInGroup.length === groupApiIds.length;
                
                return (
                  <div key={groupKey} style={{ marginBottom: '15px', border: '1px solid #eee', padding: '10px', borderRadius: '4px' }}>
                    <div style={{ marginBottom: '8px', fontWeight: 'bold', backgroundColor: '#f8f9fa', padding: '5px', borderRadius: '3px' }}>
                      <label>
                        <input
                          type="checkbox"
                          checked={isGroupSelected}
                          onChange={(e) => handleBatchPermissionChange('apis', groupApiIds, e.target.checked)}
                        />
                        <span style={{ marginLeft: '5px' }}>
                          /{groupKey} ({groupApis.length}项)
                        </span>
                      </label>
                    </div>
                    
                    <div style={{ marginLeft: '20px' }}>
                      {groupApis.map(api => (
                        <div key={api.id} style={{ marginBottom: '3px' }}>
                          <label style={{ fontSize: '14px' }}>
                            <input
                              type="checkbox"
                              checked={selectedPermissions.apis.includes(api.id)}
                              onChange={(e) => handlePermissionChange('apis', api.id, e.target.checked)}
                            />
                            <span style={{ marginLeft: '5px' }}>
                              {api.method} {api.path} - {api.description}
                            </span>
                          </label>
                        </div>
                      ))}
                    </div>
                  </div>
                );
              })}
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
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginBottom: '10px' }}>
              <h4 style={{ margin: 0 }}>选择角色</h4>
              <button
                onClick={() => handleSelectAll('roles', roles.map(r => r.id))}
                style={{
                  padding: '5px 10px',
                  backgroundColor: '#17a2b8',
                  color: 'white',
                  border: 'none',
                  borderRadius: '4px',
                  cursor: 'pointer',
                  fontSize: '12px'
                }}
              >
                {selectedPermissions.roles?.length === roles.length ? '取消全选' : '全选'}
              </button>
            </div>
            
            <div style={{ maxHeight: '300px', overflowY: 'auto', border: '1px solid #ddd', padding: '10px' }}>
              {roles.map(role => (
                <div key={role.id} style={{ marginBottom: '8px', padding: '5px', border: '1px solid #eee', borderRadius: '4px' }}>
                  <label>
                    <input
                      type="checkbox"
                      checked={selectedPermissions.roles?.includes(role.id) || false}
                      onChange={(e) => handlePermissionChange('roles', role.id, e.target.checked)}
                    />
                    <span style={{ marginLeft: '5px' }}>
                      {role.name} - {role.description}
                    </span>
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