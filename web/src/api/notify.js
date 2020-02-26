import request from '@/utils/request'


export function getnotify() {
    return request({
        url: '/api/v1/notify',
        method: 'get',
    })
}

export function readnotify(data) {
    return request({
        url: '/api/v1/notify',
        method: 'put',
        data: data
    })
}