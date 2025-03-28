import request from '@/utils/request'

// 获取章节图谱
export function getChapterGraph() {
  return request({
    url: '/api/knowledge/chapter',
    method: 'get'
  })
}

// 获取小节图谱（通过章节ID）
export function getSectionGraph(chapterId) {
  return request({
    url: '/api/knowledge/sectionByID',
    method: 'get',
    params: { chapter_id: chapterId }
  })
}

// 获取知识点图谱（通过小节ID）
export function getPointGraph(sectionId) {
  return request({
    url: '/api/knowledge/pointByID',
    method: 'get',
    params: { section_id: sectionId }
  })
}

// 直接获取所有小节图谱
export function getAllSections() {
  return request({
    url: '/api/knowledge/section',
    method: 'get'
  })
}

// 直接获取所有知识点图谱
export function getAllPoints() {
  return request({
    url: '/api/knowledge/point',
    method: 'get'
  })
}

// 通过章节ID获取小节
export function getSectionsByChapterId(chapterId) {
  return request({
    url: '/api/knowledge/sectionByID',
    method: 'get',
    params: { chapter_id: chapterId }
  })
}

// 通过小节ID获取知识点
export function getPointsBySectionId(sectionId) {
  return request({
    url: '/api/knowledge/pointByID',
    method: 'get',
    params: { section_id: sectionId }
  })
}

// 获取知识点详情
export function getKnowledgePoint(pointId) {
  return request({
    url: `/api/knowledge/point/${pointId}`,
    method: 'get'
  })
}

// 获取知识点相关视频
export function getKnowledgeVideos(pointId) {
  return request({
    url: `/api/knowledge/videosByPointId`,
    method: 'get',
    params: {
      pointId
    }
  })
}

// 获取知识点相关练习
export function getKnowledgeExercises(pointId) {
  return request({
    url: `/api/knowledge/exercisesByPointId`,
    method: 'get',
    params: {
      pointId
    }
  })
}

// 获取知识点相关课件
export function getKnowledgeCourseware(pointId) {
  return request({
    url: `/api/knowledge/coursewaresByPointId`,
    method: 'get',
    params: {
      pointId
    }
  })
}

// 搜索知识点
export function searchKnowledge(keyword) {
  return request({
    url: '/api/knowledge/searchByKeyword',
    method: 'get',
    params: { keyword }
  })
}

// 添加新的节点
export function addNode(data) {
  return request({
    url: '/api/knowledge/addNode',
    method: 'post',
    params: {
      name: data.name,
      type: data.type,
      description: data.description
    }
  })
}

// 更新节点属性
export function updateNode(data) {
  return request({
    url: '/api/knowledge/updateNode',
    method: 'post',
    params: {
      name: data.name,
      node_type: data.nodeType,
      property_name: data.propertyName,
      new_value: data.newValue
    }
  })
}

// 删除节点
export function deleteNode(data) {
  return request({
    url: '/api/knowledge/deleteNode',
    method: 'post',
    params: {
      name: data.name,
      node_type: data.nodeType
    }
  })
}

// 添加节点间关系
export function addRelation(data) {
  return request({
    url: '/api/knowledge/addLink',
    method: 'post',
    params: {
      source_name: data.sourceName,
      target_name: data.targetName,
      relation_type: data.relationType,
      source_type: data.sourceType,
      target_type: data.targetType
    }
  })
}

// 更新节点间关系
export function updateRelation(data) {
  return request({
    url: '/api/knowledge/updateLink',
    method: 'post',
    params: {
      source_name: data.sourceName,
      target_name: data.targetName,
      new_relation_type: data.newRelationType,
      source_type: data.sourceType,
      target_type: data.targetType,
      old_relation_type: data.oldRelationType
    }
  })
}

// 删除节点间关系
export function deleteRelation(data) {
  return request({
    url: '/api/knowledge/deleteLink',
    method: 'post',
    params: {
      source_name: data.sourceName,
      target_name: data.targetName,
      source_type: data.sourceType,
      target_type: data.targetType,
      relation_type: data.relationType
    }
  })
}

// 获取特定类型节点之间的关系
export function getRelationsByTypes(data) {
  return request({
    url: '/api/knowledge/relation',
    method: 'get',
    params: {
      source_type: data.source_type,
      target_type: data.target_type
    }
  })
} 
