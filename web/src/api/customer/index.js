import customerMessageAPI from './message';
import customerAPI from './customer';

// 统一导出客户相关API
export {
  customerMessageAPI,
  customerAPI
};

// 为了向后兼容，也提供默认导出
export default {
  message: customerMessageAPI,
  customer: customerAPI
};