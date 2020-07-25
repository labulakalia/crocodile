import request from '@/utils/request'



export function login(data) {
  return request({
    url: '/api/v1/user/login',
    method: 'post',
    auth: data
  })
}

export function logout() {
  return request({
    url: '/api/v1/user/logout',
    method: 'post',
  })
}

export function getselectuser() {
  return request({
    url: '/api/v1/user/select',
    method: 'get',

  })
}


export function getallusers(params) {
  return request({
    url: '/api/v1/user/all',
    method: 'get',
    params: params
  })
}

export function getInfo() {
  return request({
    url: '/api/v1/user/info',
    method: 'get'
  })
}

export function changeselfinfo(data) {
  return request({
    url: '/api/v1/user/info',
    method: 'put',
    data: data
  })
}

// /api/v1/user/admin
export function adminchangeinfo(data) {
  return request({
    url: '/api/v1/user/admin',
    method: 'put',
    data: data
  })
}

export function admindeleteuser(data) {
  return request({
    url: '/api/v1/user/admin',
    method: 'delete',
    data: data
  })
}


export function createuser(data) {
  return request({
    url: '/api/v1/user/registry',
    method: 'post',
    data: data
  })
}

export function getalarmstatus() {
  return request({
    url: '/api/v1/user/alarmstatus',
    method: 'get'
  })
}

export function getoperatelog(params) {
  return request({
    url: '/api/v1/user/operate',
    method: 'get',
    params: params
  })
}

