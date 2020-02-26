import request from '@/utils/request'


export function gethostgroup(params) {
    return request({
        url: '/api/v1/hostgroup',
        method: 'get',
        params: params
    })
}

export function getselecthostgroup() {
    return request({
        url: '/api/v1/hostgroup/select',
        method: 'get',
    })
}

export function createhostgroup(data) {
    return request({
        url: '/api/v1/hostgroup',
        method: 'post',
        data: data
    })
}

export function deletehostgroup(data) {
    return request({
        url: '/api/v1/hostgroup',
        method: 'delete',
        data: data
    })
}

export function changehostgroup(data) {
    return request({
        url: '/api/v1/hostgroup',
        method: 'put',
        data: data
    })
}


export function gethostsbyhgid(params) {
    return request({
        url: '/api/v1/hostgroup/hosts',
        method: 'get',
        params: params
    })
}