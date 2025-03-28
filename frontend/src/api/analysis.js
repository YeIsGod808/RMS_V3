import request from '@/utils/request'

// 知识点关联度分析
export function analyzeConnections(pointId) {
  return request({
    url: '/api/knowledge/analyze/connections',
    method: 'get',
    params: { point_id: pointId }
  })
}

// 获取学习难度分布
export function getDifficultyDistribution(pointId) {
  return request({
    url: '/api/knowledge/learningDifficulty',
    method: 'get',
    params: { point_id: pointId }
  })
}

// 获取推荐学习路径
export function getLearningPath(startPointType, startPointName,endPointType, endPointName) {
  return request({
    url: '/api/knowledge/pathRecommend',
    method: 'post',
    params: {
      start_point_type: startPointType,
      start_point_name: startPointName,
      end_point_type: endPointType,
      end_point_name: endPointName
    }
  })
}

// 上传文件构建知识图谱
// 知识图谱自动构建
// 参数必须是data
export function autoConstruct(data) {
  return request({
    url: '/api/knowledge/autoConstruct',
    method: 'post',
    data, // 必须是data
    headers: {
      'Content-Type': 'multipart/form-data'
    },
    timeout: 60000
  })
}