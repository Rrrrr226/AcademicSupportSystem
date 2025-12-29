import axios from 'axios';

// 配置axios基础URL
axios.defaults.baseURL = 'http://localhost:9001';

// 添加响应拦截器
axios.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export const getSubjectLink = (staffId) => {
  return axios.get(`/subject/get/links/${staffId}`);
};