import request from '@/utils/request'

/**
 * 上传资源文件到知识点
 * @param {FormData} data 包含资源信息和文件的FormData对象
 * @returns {Promise} 上传结果
 */
export function uploadResource(data) {
  return request({
    url: '/api/knowledge/uploadResource',
    method: 'post',
    data,
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 60000 // 增加超时时间到60秒，文件上传可能需要更长时间
  })
}

/**
 * 获取知识点关联的视频资源
 * @param {string} pointId 知识点ID
 * @returns {Promise} 视频资源列表
 */
export function getVideosByPointId(pointId) {
  return request({
    url: '/api/knowledge/videosByPointId',
    method: 'get',
    params: {
      pointId: pointId
    }
  })
}

/**
 * 获取知识点关联的课件资源
 * @param {string} pointId 知识点ID
 * @returns {Promise} 课件资源列表
 */
export function getCoursewaresByPointId(pointId) {
  return request({
    url: '/api/knowledge/coursewaresByPointId',
    method: 'get',
    params: {
      pointId: pointId
    }
  })
}

/**
 * 获取知识点关联的练习题资源
 * @param {string} pointId 知识点ID
 * @returns {Promise} 练习题资源列表
 */
export function getExercisesByPointId(pointId) {
  return request({
    url: '/api/knowledge/exercisesByPointId',
    method: 'get',
    params: {
      pointId: pointId
    }
  })
}

/**
 * 删除视频资源
 * @param {number} videoId 视频ID
 * @param {number} pointId 知识点ID
 * @returns {Promise} 删除结果
 */
export function deleteVideo(videoId, pointId) {
  return request({
    url: '/api/knowledge/deleteVideo',
    method: 'post',
    params: {
      videoId,
      pointId
    }
  })
}

/**
 * 删除练习题资源
 * @param {number} exerciseId 练习题ID
 * @param {number} pointId 知识点ID
 * @returns {Promise} 删除结果
 */
export function deleteExercise(exerciseId, pointId) {
  return request({
    url: '/api/knowledge/deleteExercise',
    method: 'post',
    params: {
      exerciseId,
      pointId
    }
  })
}

/**
 * 删除课件资源
 * @param {number} coursewareId 课件ID
 * @param {number} pointId 知识点ID
 * @returns {Promise} 删除结果
 */
export function deleteCourseware(coursewareId, pointId) {
  return request({
    url: '/api/knowledge/deleteCourseware',
    method: 'post',
    params: {
      coursewareId,
      pointId
    }
  })
} 