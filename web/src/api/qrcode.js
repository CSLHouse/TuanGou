import service from '@/utils/request'

export const getQrCodeList = (params) => {
    return service({
        url: '/qrcode/list',
        method: 'get',
        params
    })
}

export const deleteQrCode = (params) => {
    return service({
        url: '/qrcode/delete',
        method: 'delete',
        params,
    })
}

export const updateQrCode = (data) => {
    return service({
        url: '/qrcode/update',
        method: 'put',
        data
    })
}

export const createQrCode = (data) => {
    return service({
        url: '/qrcode/create',
        method: 'post',
        data
    })
}

export const downloadQrCode = (params) => {
    return service({
        url: '/qrcode/download',
        method: 'get',
        params
    })
}
