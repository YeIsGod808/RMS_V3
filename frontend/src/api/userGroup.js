import request from '@/utils/request'

/**
 * 获取用户所在的用户组列表
 * @param {string} token 用户token
 * @returns {Promise} 用户组列表
 */
export function getGroupList(token) {
  return request({
    url: '/api/user-group/get-groups',
    method: 'get',
    params: { token }
  })
}

/**
 * 获取所有用户组列表（仅管理员）
 * @param {string} token 用户token
 * @returns {Promise} 所有用户组列表
 */
export function getAllGroups(token) {
  return request({
    url: '/api/user-group/all-groups',
    method: 'get',
    params: { token }
  })
}

/**
 * 获取用户组详情及成员
 * @param {number} groupId 用户组ID
 * @param {string} token 用户token
 * @returns {Promise} 用户组详情
 */
export function getGroupUsers(groupId, token) {
  return request({
    url: '/api/user-group/get-user',
    method: 'get',
    params: { group_id: groupId, token }
  })
}

/**
 * 创建新用户组
 * @param {string} name 用户组名称
 * @param {string} token 用户token
 * @returns {Promise} 创建结果
 */
export function createGroup(name, token) {
  return request({
    url: '/api/user-group/create',
    method: 'post',
    params: { name, token }
  })
}

/**
 * 编辑用户组名称
 * @param {number} groupId 用户组ID
 * @param {string} name 新的用户组名称
 * @param {string} token 用户token
 * @returns {Promise} 编辑结果
 */
export function editGroupName(groupId, name, token) {
  return request({
    url: '/api/user-group/edit-name',
    method: 'post',
    params: { group_id: groupId, name, token }
  })
}

/**
 * 删除用户组
 * @param {number} groupId 用户组ID
 * @param {string} token 用户token
 * @returns {Promise} 删除结果
 */
export function deleteGroup(groupId, token) {
  return request({
    url: '/api/user-group/delete-group',
    method: 'post',
    params: { group_id: groupId, token }
  })
}

/**
 * 添加用户到用户组
 * @param {number} groupId 用户组ID
 * @param {Array<string>} users 用户ID数组
 * @param {string} token 用户token
 * @returns {Promise} 添加结果
 */
export function addGroupUsers(groupId, users, token) {
  return request({
    url: '/api/user-group/add-user',
    method: 'post',
    params: { group_id: groupId, token },
    data: { Users: users }
  })
}

/**
 * 从用户组中删除用户
 * @param {number} groupId 用户组ID
 * @param {Array<string>} users 用户ID数组
 * @param {string} token 用户token
 * @returns {Promise} 删除结果
 */
export function deleteGroupUsers(groupId, users, token) {
  return request({
    url: '/api/user-group/delete-user',
    method: 'post',
    params: { group_id: groupId, token },
    data: { Users: users }
  })
} 