import axios from 'axios';

const BASE_URL = 'http://localhost:9001'; // 根据实际端口调整

// HDUHelp 三方登录配置（测试环境使用本地地址）
const THIRD_PARTY_API = BASE_URL;
const THIRD_PARTY_PLATFORM = 'HDUHelp';
const CALLBACK_URL = `${window.location.origin}/login/callback`;

export const getSubjectLinks = (staffId, token) =>
  axios.get(`${BASE_URL}/subject/get/links`, {
    params: { staff_id: staffId },
    headers: { Authorization: `Bearer ${token}` },
  });

/**
 * 获取三方登录跳转地址
 * @param {string} from - 来源页面，默认 '/'
 * @returns {Promise} 包含跳转 URL 的响应
 */
export const getThirdPartyJumpUrl = (from = '/') => {
  return axios.get(`${THIRD_PARTY_API}/user/v1/third/jump`, {
    params: {
      platform: THIRD_PARTY_PLATFORM,
      from: from,
      callback: CALLBACK_URL
    },
    headers: {
      'Accept': 'application/json, text/plain, */*',
      'Cache-Control': 'no-cache',
    }
  });
};

/**
 * 三方登录回调 - 用 code 换取 token
 * @param {string} code - 授权码
 * @param {string} state - 状态参数
 * @param {string|null} ticket - ticket 参数，可能为 null
 * @returns {Promise} 包含 token 的响应
 */
export const thirdPartyCallback = ({ code, state, ticket = null }) => {
  return axios.post(`${THIRD_PARTY_API}/user/v1/third/callback`, {
    ticket: ticket,
    state: state,
    code: code,
    callback: CALLBACK_URL
  }, {
    headers: {
      'Accept': 'application/json, text/plain, */*',
      'Content-Type': 'application/json',
      'Cache-Control': 'no-cache',
    }
  });
};

// ============ 管理员相关 API ============

/**
 * 管理员登录
 * @param {string} username - 用户名
 * @param {string} password - 密码
 * @returns {Promise} 包含 token 的响应
 */
export const managerLogin = (username, password) => {
  return axios.post(`${BASE_URL}/managers/login`, {
    username,
    password
  });
};

/**
 * 获取管理员信息
 * @param {string} token - 管理员 token
 * @returns {Promise} 管理员信息
 */
export const getManagerInfo = (token) => {
  return axios.get(`${BASE_URL}/managers/list`, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

/**
 * 获取管理员列表
 * @param {string} token - 管理员 token
 * @returns {Promise} 管理员列表
 */
export const getManagerList = (token) => {
  return axios.get(`${BASE_URL}/managers/list`, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

/**
 * 添加管理员
 * @param {string} staffId - 学号/工号
 * @param {string} token - 管理员 token
 * @returns {Promise} 添加结果
 */
export const addManager = (staffId, token) => {
  return axios.post(`${BASE_URL}/managers/add`, { staffId }, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

/**
 * 删除管理员
 * @param {string} staffId - 管理员学号
 * @param {string} token - 管理员 token
 * @returns {Promise} 删除结果
 */
export const deleteManager = (staffId, token) => {
  return axios.post(`${BASE_URL}/managers/delete`, {
    staffId
  }, {
    headers: { Authorization: `Bearer ${token}` }
  });
};

/**
 * 导入学生科目（Excel 上传）
 * @param {File} file - Excel 文件
 * @param {string} token - 管理员 token
 * @returns {Promise} 导入结果
 */
export const importStudentSubjects = (file, token) => {
  const formData = new FormData();
  formData.append('file', file);
  
  return axios.post(`${BASE_URL}/managers/import/students`, formData, {
    headers: {
      Authorization: `Bearer ${token}`,
      'Content-Type': 'multipart/form-data'
    }
  });
};

/**
 * 下载学生科目导入模板
 * @returns {string} 模板下载地址
 */
export const downloadImportTemplate = () => {
  return `${BASE_URL}/managers/import/template`;
};

/**
 * 获取学科列表
 * @param {string} token - 管理员 token
 * @param {number} page - 页码
 * @param {number} pageSize - 每页数量
 * @returns {Promise} 学科列表
 */
export const getSubjectList = (token, page = 1, pageSize = 10) =>
  axios.get(`${BASE_URL}/subject/v1/list`, {
    params: { page, page_size: pageSize },
    headers: { Authorization: `Bearer ${token}` },
  });

/**
 * 添加学科
 * @param {Object} data - 学科数据 { subject_name, subject_link }
 * @param {string} token - 管理员 token
 * @returns {Promise} 添加结果
 */
export const addSubject = (data, token) =>
  axios.post(`${BASE_URL}/subject/v1/add`, data, {
    headers: { Authorization: `Bearer ${token}` },
  });

/**
 * 删除学科
 * @param {number} id - 学科ID
 * @param {string} token - 管理员 token
 * @returns {Promise} 删除结果
 */
export const deleteSubject = (id, token) =>
  axios.delete(`${BASE_URL}/subject/v1/delete/${id}`, {
    headers: { Authorization: `Bearer ${token}` },
  });

/**
 * 更新学科
 * @param {Object} data - 学科数据 { subject_id, subject_name, subject_link }
 * @param {string} token - 管理员 token
 * @returns {Promise} 更新结果
 */
export const updateSubject = (data, token) =>
  axios.post(`${BASE_URL}/subject/v1/update`, data, {
    headers: { Authorization: `Bearer ${token}` },
  });