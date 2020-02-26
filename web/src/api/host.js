import request from '@/utils/request'


export function getselecthost() {
    return request({
        url: '/api/v1/host/select',
        method: 'get',
    })
}

export function gethost(params) {
    return request({
        url: '/api/v1/host',
        method: 'get',
        params: params
    })
}

export function stophost(data) {
    return request({
        url: '/api/v1/host/stop',
        method: 'put',
        data: data
    })
}

export function deletehost(data) {
    return request({
        url: '/api/v1/host',
        method: 'delete',
        data: data
    })
}