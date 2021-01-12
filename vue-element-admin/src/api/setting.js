import request from '@/utils/request'
import { httphost } from '@/utils/global'
 
export function UpdateIndexPics(data) {
  return request({
    url: '',
    method: 'post',
    data,
    baseURL: httphost + '/UpdateIndexPics'
  })
}
 

export function GetIndexPics() {
  return request({
    url: '',
    method: 'get',
    baseURL: httphost + '/GetIndexPics'
  })
}




 export function PostWebsiteInfo(data) {
  return request({
    url: '',
    method: 'post',
     data,
    baseURL: httphost + '/PostWebsiteInfo'
  })
}
 



 export function GetWebsiteInfo() {
  return request({
    url: '',
    method: 'get',
    baseURL: httphost + '/GetWebsiteInfo'
  })
}
 
 