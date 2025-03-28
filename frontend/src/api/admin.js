import request from '@/utils/request'

// 用户管理
export function getUsers() {
  return request({
    url: '/api/admin/users',
    method: 'get'
  })
}

export function addUser(data) {
  return request({
    url: '/api/admin/users',
    method: 'post',
    data
  })
}

export function updateUser(id, data) {
  return request({
    url: `/api/admin/users/${id}`,
    method: 'put',
    data
  })
}

export function deleteUser(id) {
  return request({
    url: `/api/admin/users/${id}`,
    method: 'delete'
  })
}

export function resetPassword(id) {
  return request({
    url: `/api/admin/users/${id}/reset-password`,
    method: 'post'
  })
}

// 内容管理
export function getContents(type) {
  return request({
    url: '/api/admin/contents',
    method: 'get',
    params: { type }
  })
}

export function addContent(data) {
  return request({
    url: '/api/admin/contents',
    method: 'post',
    data
  })
}

export function updateContent(id, data) {
  return request({
    url: `/api/admin/contents/${id}`,
    method: 'put',
    data
  })
}

export function deleteContent(id) {
  return request({
    url: `/api/admin/contents/${id}`,
    method: 'delete'
  })
}

// 资源管理
export function getResources(type) {
  return request({
    url: '/api/admin/resources',
    method: 'get',
    params: { type }
  })
}

export function addResource(data) {
  return request({
    url: '/api/admin/resources',
    method: 'post',
    data
  })
} 