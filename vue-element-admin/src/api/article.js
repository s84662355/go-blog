import request from '@/utils/request'
import { httphost } from '@/utils/global'

export function articleList(query) {
  return request({
    url: '',
    method: 'get',
    params: query,
    baseURL: httphost + '/article/list'
  })
}

 

export function cateInfo(query) {
  return request({
    url: '',
    method: 'get',
    params: query,
    baseURL: httphost + '/cate/info'
  })
}

export function cateList() {
  return request({
    url: '',
    method: 'get',
    baseURL: httphost + '/cate/list'
  })
}


export function cateCreate(query) {
  return request({
    url: '',
    method: 'post',
    params: query,
    baseURL: httphost + '/cate/create'
  })
}

export function cateUpdate (query) {
  return request({
    url: '',
    method: 'post',
    params: query,
    baseURL: httphost + '/cate/update'
  })
}

export function fetchList(query) {
  return request({
    url: '',
    method: 'get',
    params: query,
    baseURL: httphost + '/articles/list'
  })
}

export function fetchArticle(id) {
  return request({
    url: '',
    method: 'get',
    params: { id },
    baseURL: httphost + '/article/detail'
  })
}

export function fetchPv(pv) {
  return request({
    url: '/article/pv',
    method: 'get',
    params: { pv }
  })
}

export function createArticle(data) {
  return request({
    url: '',
    method: 'post',
    data,
    baseURL: httphost + '/article/create'
  })
}

export function updateArticle(data) {
  return request({
    url: '',
    method: 'post',
    data,
    baseURL: httphost + '/article/update'
  })
}

export function delImage(url) {
  return request({
    url: '',
    method: 'get',
    params: { url },
    baseURL: httphost + '/del/image'
  })
}
