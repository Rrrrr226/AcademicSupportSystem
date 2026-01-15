import axios from 'axios';
import { BASE_URL } from '../../api/config';

// 配置axios基础URL
axios.defaults.baseURL = BASE_URL;

// 添加响应拦截器
axios.interceptors.response.use(
  response => response,
  error => {
    console.error('API Error:', error.response?.data || error.message);
    return Promise.reject(error);
  }
);

export const getSubjectLink = (staffId) => {
  const token = localStorage.getItem('token');
  return axios.get(`/subject/get/links/${staffId}`, {
    headers: { Authorization: `Bearer ${token}` }
  });
};