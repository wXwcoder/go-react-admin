import React, { useState, useEffect } from 'react';
import '../assets/styles/management.css'; // 引入样式文件
import { menuApi } from '../api'; // 引入API
import Modal from 'react-modal'; // 引入模态框组件

const MenuManagement = () => {
  const [menus, setMenus] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [currentMenu, setCurrentMenu] = useState(null); // 用于存储当前编辑或创建的菜单

  useEffect(() => {
    // 从API获取菜单数据
    const fetchMenus = async () => {
      try {
        setLoading(true);
        const response = await menuApi.getMenuList();
        setMenus(response.data.menus);
      } catch (error) {
        console.error('获取菜单数据失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchMenus();
  }, []);

  // 创建菜单
  const handleCreateMenu = async (menuData) => {
    try {
      const response = await menuApi.createMenu(menuData);
      // 重新获取菜单列表
      const fetchResponse = await menuApi.getMenuList();
      setMenus(fetchResponse.data.menus);
    } catch (error) {
      console.error('创建菜单失败:', error);
    }
  };

  // 更新菜单
  const handleUpdateMenu = async (id, menuData) => {
    try {
      const response = await menuApi.updateMenu(id, menuData);
      // 重新获取菜单列表
      const fetchResponse = await menuApi.getMenuList();
      setMenus(fetchResponse.data.menus);
    } catch (error) {
      console.error('更新菜单失败:', error);
    }
  };

  // 删除菜单
  const handleDeleteMenu = async (id) => {
    try {
      const response = await menuApi.deleteMenu(id);
      // 重新获取菜单列表
      const fetchResponse = await menuApi.getMenuList();
      setMenus(fetchResponse.data.menus);
    } catch (error) {
      console.error('删除菜单失败:', error);
    }
  };

  const handleCreateClick = () => {
    setCurrentMenu(null); // 创建时清空当前菜单
    setModalIsOpen(true);
  };

  return (
    <div className="management-container"> {/* 使用management-container类名 */}
      <h2>菜单管理</h2>
      <button 
        style={{ marginBottom: '10px', padding: '5px 10px', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
        onClick={handleCreateClick}
      >
        创建菜单
      </button>
      {loading ? (
        <p className="loading">加载中...</p> // 使用loading类名
      ) : (
        <table className="management-table"> {/* 使用management-table类名 */}
          <thead>
            <tr>
              <th>ID</th>
              <th>菜单名称</th>
              <th>路径</th>
              <th>图标</th>
              <th>父级ID</th>
              <th>排序</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {menus.map(menu => (
              <tr key={menu.id}>
              <td>{menu.id}</td>
              <td>{menu.name}</td>
              <td>{menu.path}</td>
              <td>{menu.icon}</td>
              <td>{menu.parent_id}</td>
              <td>{menu.sort}</td>
              <td>{menu.status === 1 ? '启用' : '禁用'}</td>
              <td>
                  <button 
                    style={{ marginRight: '5px', padding: '5px 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => {
                      setCurrentMenu(menu);
                      setModalIsOpen(true);
                    }}
                  >
                    编辑
                  </button>
                  <button 
                    style={{ padding: '5px 10px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => handleDeleteMenu(menu.id)}
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
        contentLabel="菜单模态框"
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
        <h2>{currentMenu ? '编辑菜单' : '创建菜单'}</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.target);
          const menuData = {
            name: formData.get('name'),
            path: formData.get('path'),
            icon: formData.get('icon'),
            parent_id: parseInt(formData.get('parent_id')),
            sort: parseInt(formData.get('sort')),
            status: parseInt(formData.get('status')) || 1
          };
          
          if (currentMenu) {
            // 更新菜单
            handleUpdateMenu(currentMenu.id, menuData);
          } else {
            // 创建菜单
            handleCreateMenu(menuData);
          }
          setModalIsOpen(false);
        }}>
          <div style={{ marginBottom: '15px' }}>
            <label>菜单名称: </label>
            <input 
              type="text" 
              name="name" 
              defaultValue={currentMenu?.name || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>路径: </label>
            <input 
              type="text" 
              name="path" 
              defaultValue={currentMenu?.path || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>图标: </label>
            <select 
              name="icon" 
              defaultValue={currentMenu?.icon || ''} 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            >
              <option value="">请选择图标</option>
              <option value="dashboard">dashboard</option>
              <option value="user">user</option>
              <option value="setting">setting</option>
              <option value="menu">menu</option>
              <option value="role">role</option>
              <option value="permission">permission</option>
              <option value="tenant">tenant</option>
              <option value="api">api</option>
              <option value="log">log</option>
              <option value="home">home</option>
            </select>
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>父级ID: </label>
            <input 
              type="number" 
              name="parent_id" 
              defaultValue={currentMenu?.parent_id || 0} 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>排序: </label>
            <input 
              type="number" 
              name="sort" 
              defaultValue={currentMenu?.sort || 0} 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>状态: </label>
            <select 
              name="status" 
              defaultValue={currentMenu?.status || 1} 
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

export default MenuManagement;