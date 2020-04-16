import request from '@/utils/request'


export function queryinstallstatus() {
    return request({
        url: '/api/v1/install/status',
        method: 'get',
    })
}


export function startinstall(data) {
    return request({
        url: '/api/v1/install',
        method: 'post',
        data: data
    })
}

export function queryversion() {
    return request({
        url: '/api/v1/install/version',
        method: 'get',
    })
}