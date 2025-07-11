import axios from 'axios';

const BASE_URL = 'http://localhost:9001'; // 根据实际端口调整

export const register = (data) =>
  axios.post(`${BASE_URL}/user/v1/direct/register`, data);

export const login = (data) =>
  axios.post(`${BASE_URL}/user/v1/direct/login`, data);

export const getSubjectLinks = (staffId, token) =>
  axios.get(`${BASE_URL}/subject/get/links`, {
    params: { staff_id: staffId },
    headers: { Authorization: `Bearer ${token}` },
  });