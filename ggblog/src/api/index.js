import request from '@/utils/request'
import { httphost } from '@/utils/global'


export function setting(key) {
    return request({
        url: httphost +'/setting/'+key,
        method: 'get',
        
    })
}


export function fetchList(params) {
    return request({
        url: httphost +'/article',
        method: 'get',
        params: params
    })
}

export function fetchFocus() {
    return request({
        url: httphost +'/GetFeatureSetting',
        method: 'get',
        params: {}
    })
}

export function fetchCategory() {
    return request({
        url:httphost + '/cate',
        method: 'get',
        params: {}
    })
}

export function fetchFriend() {
    return request({
        url: httphost +'/friend',
        method: 'get',
        params: {}
    })
}

export function fetchSocial() {
    return request({
        url: httphost +'/social',
        method: 'get',
        params: {}
    });
}

export function fetchSiteInfo() {
    return request({
        url:httphost + '/site',
        method: 'get',
        params: {}
    })
}

 

export function getArticle(id) {
    return request({
        url: httphost +'/article/'+id,
        method: 'get',
        params: {}
    })
}

export function fetchComment() {
    return request({
        url: httphost +'/comment',
        method: 'get',
        params: {}
    })
}
