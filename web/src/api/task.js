import request from '@/utils/request'

export function gettask(params) {
  return request({
    url: '/api/v1/task',
    method: 'get',
    params: params
  })
}

export function createtask(data) {
  return request({
    url: '/api/v1/task',
    method: 'post',
    data: data
  })
}

export function changetask(data) {
  return request({
    url: '/api/v1/task',
    method: 'put',
    data: data
  })
}

export function deletetask(data) {
  return request({
    url: '/api/v1/task',
    method: 'delete',
    data: data
  })
}

export function killtask(data) {
  return request({
    url: '/api/v1/task/kill',
    method: 'put',
    data: data
  })
}

export function runtask(data) {
  return request({
    url: '/api/v1/task/run',
    method: 'put',
    data: data
  })
}

export function getrunningtasks(params) {
  return request({
    url: '/api/v1/task/running',
    method: 'get',
    params: params
  })
}

export function gettaskLog(params) {
  return request({
    url: '/api/v1/task/log',
    method: 'get',
    params: params
  })
}

export function gettaskLogTree(params) {
  return request({
    url: '/api/v1/task/log/tree',
    method: 'get',
    params: params
  })
}

export function parsecron(params) {
  return request({
    url: '/api/v1/task/cron',
    method: 'get',
    params: params
  })
}

export function getselecttask() {
  return request({
    url: '/api/v1/task/select',
    method: 'get'
  })
}

export function clonetask(data) {
  return request({
    url: '/api/v1/task/clone',
    method: 'post',
    data: data
  })
}

export function cleantasklog(data) {
  return request({
    url: '/api/v1/task/log',
    method: 'delete',
    data: data
  })
}