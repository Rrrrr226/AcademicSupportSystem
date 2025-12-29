import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/api/users/login',
    method: 'post',
    data
  })
}

export function getPersonInfo(data) {
  const token = localStorage.getItem('token');
  return request({
    url: '/user/v1',
    method: 'get',
    params: data,
    headers: { Authorization: `Bearer ${token}` }
  })
}

export function modifyUserInfo(data) {
  const token = localStorage.getItem('token');
  return request({
    url: '/user/v1/direct/modify',
    method: 'post',
    data,
    headers: { Authorization: `Bearer ${token}` }
  })
}