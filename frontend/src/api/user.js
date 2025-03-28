import request from '@/utils/request'

// 用户登录
export async function login(username, password) {
  // 构建表单数据
  const formData = new FormData()
  formData.append('id', username)
  formData.append('password', password)
  
  return request({
    url: '/api/user/login',
    method: 'post',
    data: formData
  })
}

// 用户登出
export async function logout() {
  // 后端可能不需要实际的登出请求，前端清理即可
  // 但为了兼容性，我们保留请求
  return Promise.resolve({ ret: "0", msg: "登出成功" })
}

// 校验用户token
export async function checkToken(token) {
  return request({
    url: '/api/user/checktoken',
    method: 'get',
    params: { token }
  })
}

// 修改密码
export async function changePassword(token, newPassword, userId = null) {
  const params = { 
    token,
    new_password: newPassword
  };
  
  if (userId) {
    params.user_id = userId;
  }
    
  return request({
    url: '/api/user/change-password',
    method: 'post',
    params
  })
}

// 用户注册（仅限学生）
export async function register(username, password, nickname) {
  const formData = new FormData()
  formData.append('id', username)
  formData.append('password', password)
  formData.append('nickname', nickname)
  
  const response = await request({
    url: '/api/user/register',
    method: 'post',
    data: formData
  })
  
  return response
}

/**
 * 获取用户列表
 * @param {string} token 用户token
 * @returns {Promise} 用户列表
 */
export function getUserList(token) {
  return request({
    url: '/api/user/list',
    method: 'get',
    params: { token }
  })
}

/**
 * 添加新用户
 * @param {Object} data 用户数据
 * @param {string} token 用户token
 * @returns {Promise} 添加结果
 */
export function addUser(data, token) {
  return request({
    url: '/api/user/add-user',
    method: 'post',
    data,
    params: { token }
  })
}

/**
 * 更新用户信息
 * @param {Object} data 用户数据
 * @param {string} token 用户token
 * @returns {Promise} 更新结果
 */
export function updateUser(data, token) {
  return request({
    url: '/api/user/update',
    method: 'post',
    data,
    params: { token }
  })
}

/**
 * 删除用户
 * @param {string} userId 用户ID
 * @param {string} token 用户token
 * @returns {Promise} 删除结果
 */
export function deleteUser(userId, token) {
  return request({
    url: '/api/user/delete',
    method: 'post',
    data: { user_id: userId },
    params: { token }
  })
} 