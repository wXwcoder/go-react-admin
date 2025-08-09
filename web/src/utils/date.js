import dayjs from 'dayjs';

// 格式化日期时间
export const formatDateTime = (date, format = 'YYYY-MM-DD HH:mm:ss') => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 格式化日期
export const formatDate = (date, format = 'YYYY-MM-DD') => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 格式化时间
export const formatTime = (date, format = 'HH:mm:ss') => {
  if (!date) return '';
  return dayjs(date).format(format);
};

// 获取相对时间
export const fromNow = (date) => {
  if (!date) return '';
  return dayjs(date).fromNow();
};

// 获取时间差（毫秒）
export const getTimeDiff = (start, end) => {
  return dayjs(end).valueOf() - dayjs(start).valueOf();
};

// 检查是否为今天
export const isToday = (date) => {
  if (!date) return false;
  return dayjs(date).isToday();
};

// 检查是否为昨天
export const isYesterday = (date) => {
  if (!date) return false;
  return dayjs(date).isYesterday();
};

// 获取本周的开始和结束时间
export const getWeekRange = () => {
  const start = dayjs().startOf('week');
  const end = dayjs().endOf('week');
  return [start, end];
};

// 获取本月的开始和结束时间
export const getMonthRange = () => {
  const start = dayjs().startOf('month');
  const end = dayjs().endOf('month');
  return [start, end];
};

// 获取时间戳
export const getTimestamp = (date) => {
  if (!date) return Date.now();
  return dayjs(date).valueOf();
};