import React, { useState, useEffect } from 'react';
import '../assets/styles/management.css'; // 引入样式文件
import { apiApi } from '../api'; // 引入API
import Modal from 'react-modal'; // 引入模态框组件

const ApiManagement = () => {
  const [apis, setApis] = useState([]);
  const [loading, setLoading] = useState(true);
  const [modalIsOpen, setModalIsOpen] = useState(false);
  const [currentApi, setCurrentApi] = useState(null); // 用于存储当前编辑或创建的API

  useEffect(() => {
    // 从API获取API数据
    const fetchApis = async () => {
      try {
        setLoading(true);
        const response = await apiApi.getApiList();
        setApis(response.data.apis);
      } catch (error) {
        console.error('获取API数据失败:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchApis();
  }, []);

  // 创建API
  const handleCreateApi = async (apiData) => {
    try {
      const response = await apiApi.createApi(apiData);
      // 重新获取API列表
      const fetchResponse = await apiApi.getApiList();
      setApis(fetchResponse.data.apis);
    } catch (error) {
      console.error('创建API失败:', error);
    }
  };

  // 更新API
  const handleUpdateApi = async (id, apiData) => {
    try {
      const response = await apiApi.updateApi(id, apiData);
      // 重新获取API列表
      const fetchResponse = await apiApi.getApiList();
      setApis(fetchResponse.data.apis);
    } catch (error) {
      console.error('更新API失败:', error);
    }
  };

  // 删除API
  const handleDeleteApi = async (id) => {
    try {
      const response = await apiApi.deleteApi(id);
      // 重新获取API列表
      const fetchResponse = await apiApi.getApiList();
      setApis(fetchResponse.data.apis);
    } catch (error) {
      console.error('删除API失败:', error);
    }
  };

  const handleCreateClick = () => {
    setCurrentApi(null); // 创建时清空当前API
    setModalIsOpen(true);
  };

  // 保存API（创建或更新）
  const handleSaveApi = async (apiData) => {
    try {
      if (currentApi) {
        // 更新API
        await handleUpdateApi(currentApi.id, apiData);
      } else {
        // 创建API
        await handleCreateApi(apiData);
      }
      setModalIsOpen(false);
    } catch (error) {
      console.error('保存API失败:', error);
    }
  };

  return (
    <div className="management-container">
      <h2>API管理</h2>
      <button 
        style={{ marginBottom: '10px', padding: '5px 10px', backgroundColor: '#28a745', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
        onClick={handleCreateClick}
      >
        创建API
      </button>
      {loading ? (
        <p className="loading">加载中...</p>
      ) : (
        <table className="management-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>路径</th>
              <th>方法</th>
              <th>分类</th>
              <th>描述</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody>
            {apis.map((api) => (
              <tr key={api.id}>
                <td>{api.id}</td>
                <td>{api.path}</td>
                <td>{api.method}</td>
                <td>{api.category}</td>
                <td>{api.description}</td>
                <td>{api.status === 1 ? '启用' : '禁用'}</td>
                <td>
                  <button 
                    style={{ marginRight: '5px', padding: '5px 10px', backgroundColor: '#007bff', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => {
                      setCurrentApi(api);
                      setModalIsOpen(true);
                    }}
                  >
                    编辑
                  </button>
                  <button 
                    style={{ padding: '5px 10px', backgroundColor: '#dc3545', color: 'white', border: 'none', borderRadius: '4px', cursor: 'pointer' }}
                    onClick={() => handleDeleteApi(api.id)}
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
        contentLabel="API模态框"
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
        <h2>{currentApi ? '编辑API' : '创建API'}</h2>
        <form onSubmit={(e) => {
          e.preventDefault();
          const formData = new FormData(e.target);
          const apiData = {
            path: formData.get('path'),
            method: formData.get('method'),
            category: formData.get('category'),
            description: formData.get('description'),
            status: parseInt(formData.get('status')) || 1
          };
          handleSaveApi(apiData);
        }}>
          <div style={{ marginBottom: '15px' }}>
            <label>路径: </label>
            <input 
              type="text" 
              name="path" 
              defaultValue={currentApi?.path || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>方法: </label>
            <input 
              type="text" 
              name="method" 
              defaultValue={currentApi?.method || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>分类: </label>
            <input 
              type="text" 
              name="category" 
              defaultValue={currentApi?.category || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>描述: </label>
            <input 
              type="text" 
              name="description" 
              defaultValue={currentApi?.description || ''} 
              required 
              style={{ width: '100%', padding: '8px', marginTop: '5px' }}
            />
          </div>
          <div style={{ marginBottom: '15px' }}>
            <label>状态: </label>
            <select 
              name="status" 
              defaultValue={currentApi?.status || 1} 
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

export default ApiManagement;