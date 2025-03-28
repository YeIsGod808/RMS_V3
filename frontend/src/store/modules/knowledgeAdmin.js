import { defineStore } from 'pinia'
import { 
  getChapterGraph, 
  getSectionGraph, 
  getPointGraph,
  getAllSections,
  getAllPoints,
  getSectionsByChapterId,
  getPointsBySectionId,
  addNode,
  updateNode,
  deleteNode,
  addRelation,
  updateRelation,
  deleteRelation,
  getRelationsByTypes
} from '@/api/knowledge'
import { message } from 'ant-design-vue'

export const useKnowledgeAdminStore = defineStore('knowledgeAdmin', {
  state: () => ({
    chapters: [],
    sections: [],
    points: [],
    allSections: [], // 存储所有小节
    allPoints: [], // 存储所有知识点
    currentChapter: null,
    currentSection: null,
    loading: false,
    relations: [], // 存储关系数据
    nodeTypes: [
      { value: 'chapter', label: '章节' },
      { value: 'section', label: '小节' },
      { value: 'point', label: '知识点' }
    ],
    relationTypes: [
      { value: 'CONTAINS', label: '包含关系' },
      { value: 'RELATED_TO', label: '相关关系' },
      { value: 'PREREQUISITE', label: '前置关系' }
    ]
  }),

  getters: {
    chapterOptions() {
      return this.chapters.map(chapter => ({
        value: chapter.name,
        label: chapter.name,
        id: chapter.id
      }))
    },
    sectionOptions() {
      return this.sections.map(section => ({
        value: section.name,
        label: section.name,
        id: section.id
      }))
    },
    pointOptions() {
      return this.points.map(point => ({
        value: point.name,
        label: point.name,
        id: point.id
      }))
    },
    allSectionOptions() {
      return this.allSections.map(section => ({
        value: section.name,
        label: section.name,
        id: section.id
      }))
    },
    allPointOptions() {
      return this.allPoints.map(point => ({
        value: point.name,
        label: point.name,
        id: point.id
      }))
    }
  },

  actions: {
    // 加载章节节点
    async loadChapters() {
      try {
        this.loading = true
        const response = await getChapterGraph()
        const data = response.data || response
        
        if (data && data.nodes) {
          this.chapters = data.nodes.filter(node => node.type === 'chapter')
          
          // 提取关系数据
          if (data.links) {
            const chapterRelations = data.links.map(link => {
              const sourceNode = data.nodes.find(node => node.id === link.source)
              const targetNode = data.nodes.find(node => node.id === link.target)
              
              if (sourceNode && targetNode) {
                return {
                  id: `${sourceNode.name}-${targetNode.name}-${link.type}`,
                  sourceName: sourceNode.name,
                  targetName: targetNode.name,
                  sourceType: sourceNode.type,
                  targetType: targetNode.type,
                  relationType: link.type
                }
              }
              return null
            }).filter(Boolean)
            
            // 更新关系数据
            this.updateRelationsInStore(chapterRelations)
          }
        } else {
          this.chapters = []
        }
      } catch (error) {
        message.error('加载章节失败')
        console.error(error)
        this.chapters = []
      } finally {
        this.loading = false
      }
    },

    // 加载小节节点
    async loadSections(chapterId) {
      try {
        this.loading = true
        const response = await getSectionGraph(chapterId)
        const data = response.data || response
        
        if (data && data.nodes) {
          this.sections = data.nodes.filter(node => node.type === 'section')
          this.currentChapter = chapterId
        } else {
          this.sections = []
        }
      } catch (error) {
        message.error('加载小节失败')
        console.error(error)
        this.sections = []
      } finally {
        this.loading = false
      }
    },

    // 加载知识点节点
    async loadPoints(sectionId) {
      try {
        this.loading = true
        const response = await getPointGraph(sectionId)
        const data = response.data || response
        
        if (data && data.nodes) {
          this.points = data.nodes.filter(node => node.type === 'point')
          this.currentSection = sectionId
        } else {
          this.points = []
        }
      } catch (error) {
        message.error('加载知识点失败')
        console.error(error)
        this.points = []
      } finally {
        this.loading = false
      }
    },

    // 加载所有小节（新增）
    async loadAllSections() {
      try {
        this.loading = true
        const response = await getAllSections()
        const data = response.data || response
        
        if (data && data.nodes) {
          this.allSections = data.nodes.filter(node => node.type === 'section')
          
          // 提取关系数据
          if (data.links) {
            const sectionRelations = data.links.map(link => {
              const sourceNode = data.nodes.find(node => node.id === link.source)
              const targetNode = data.nodes.find(node => node.id === link.target)
              
              if (sourceNode && targetNode) {
                return {
                  id: `${sourceNode.name}-${targetNode.name}-${link.type}`,
                  sourceName: sourceNode.name,
                  targetName: targetNode.name,
                  sourceType: sourceNode.type,
                  targetType: targetNode.type,
                  relationType: link.type
                }
              }
              return null
            }).filter(Boolean)
            
            // 更新关系数据
            this.updateRelationsInStore(sectionRelations)
          }
        } else {
          this.allSections = []
        }
      } catch (error) {
        message.error('加载所有小节失败')
        console.error(error)
        this.allSections = []
      } finally {
        this.loading = false
      }
    },

    // 加载所有知识点（新增）
    async loadAllPoints() {
      try {
        this.loading = true
        const response = await getAllPoints()
        const data = response.data || response
        
        if (data && data.nodes) {
          this.allPoints = data.nodes.filter(node => node.type === 'point')
          
          // 提取关系数据
          if (data.links) {
            const pointRelations = data.links.map(link => {
              const sourceNode = data.nodes.find(node => node.id === link.source)
              const targetNode = data.nodes.find(node => node.id === link.target)
              
              if (sourceNode && targetNode) {
                return {
                  id: `${sourceNode.name}-${targetNode.name}-${link.type}`,
                  sourceName: sourceNode.name,
                  targetName: targetNode.name,
                  sourceType: sourceNode.type,
                  targetType: targetNode.type,
                  relationType: link.type
                }
              }
              return null
            }).filter(Boolean)
            
            // 更新关系数据
            this.updateRelationsInStore(pointRelations)
          }
        } else {
          this.allPoints = []
        }
      } catch (error) {
        message.error('加载所有知识点失败')
        console.error(error)
        this.allPoints = []
      } finally {
        this.loading = false
      }
    },

    // 通过章节ID加载小节（新增）
    async loadSectionsByChapterId(chapterId) {
      try {
        this.loading = true
        const response = await getSectionsByChapterId(chapterId)
        const data = response.data || response
        
        if (data && data.nodes) {
          this.sections = data.nodes.filter(node => node.type === 'section')
          this.currentChapter = chapterId
        } else {
          this.sections = []
        }
      } catch (error) {
        message.error('通过章节ID加载小节失败')
        console.error(error)
        this.sections = []
      } finally {
        this.loading = false
      }
    },

    // 通过小节ID加载知识点（新增）
    async loadPointsBySectionId(sectionId) {
      try {
        this.loading = true
        const response = await getPointsBySectionId(sectionId)
        const data = response.data || response
        
        if (data && data.nodes) {
          this.points = data.nodes.filter(node => node.type === 'point')
          this.currentSection = sectionId
        } else {
          this.points = []
        }
      } catch (error) {
        message.error('通过小节ID加载知识点失败')
        console.error(error)
        this.points = []
      } finally {
        this.loading = false
      }
    },

    // 添加节点
    async addNode(nodeData) {
      try {
        this.loading = true
        const response = await addNode(nodeData)
        const data = response.data || response
        
        // 检查是否返回了已存在节点的消息
        if (data && data.message && data.message.includes('已存在')) {
          message.warning(`节点 "${nodeData.name}" 已存在，请使用其他名称`)
          return false
        }
        
        message.success('添加节点成功')
        
        // 根据节点类型刷新相应数据
        if (nodeData.type === 'chapter') {
          await this.loadChapters()
        } else if (nodeData.type === 'section') {
          await this.loadAllSections()
          if (this.currentChapter) {
            await this.loadSections(this.currentChapter)
          }
        } else if (nodeData.type === 'point') {
          await this.loadAllPoints()
          if (this.currentSection) {
            await this.loadPoints(this.currentSection)
          }
        }
        
        return true
      } catch (error) {
        message.error('添加节点失败')
        console.error(error)
        return false
      } finally {
        this.loading = false
      }
    },

    // 更新节点
    async updateNode(nodeData) {
      try {
        this.loading = true
        await updateNode(nodeData)
        message.success('更新节点成功')
        
        // 根据节点类型刷新相应数据
        if (nodeData.nodeType === 'chapter') {
          await this.loadChapters()
        } else if (nodeData.nodeType === 'section') {
          await this.loadAllSections()
          if (this.currentChapter) {
            await this.loadSections(this.currentChapter)
          }
        } else if (nodeData.nodeType === 'point') {
          await this.loadAllPoints()
          if (this.currentSection) {
            await this.loadPoints(this.currentSection)
          }
        }
        
        return true
      } catch (error) {
        message.error('更新节点失败')
        console.error(error)
        return false
      } finally {
        this.loading = false
      }
    },

    // 删除节点
    async deleteNode(nodeData) {
      try {
        this.loading = true
        await deleteNode(nodeData)
        message.success('删除节点成功')
        
        // 根据节点类型刷新相应数据
        if (nodeData.nodeType === 'chapter') {
          await this.loadChapters()
        } else if (nodeData.nodeType === 'section') {
          await this.loadAllSections()
          if (this.currentChapter) {
            await this.loadSections(this.currentChapter)
          }
        } else if (nodeData.nodeType === 'point') {
          await this.loadAllPoints()
          if (this.currentSection) {
            await this.loadPoints(this.currentSection)
          }
        }
        
        return true
      } catch (error) {
        message.error('删除节点失败')
        console.error(error)
        return false
      } finally {
        this.loading = false
      }
    },

    // 添加关系
    async addRelation(relationData) {
      try {
        this.loading = true
        const response = await addRelation(relationData)
        const data = response.data || response
        
        // 检查是否返回了已存在关系的消息
        if (data && data.message && data.message.includes('已存在')) {
          message.warning(`节点 "${relationData.sourceName}" 与 "${relationData.targetName}" 之间的关系已存在`)
          return false
        }
        
        message.success('添加关系成功')
        
        // 重新加载相关类型的关系
        await this.loadRelationsByTypes(relationData.sourceType, relationData.targetType)
        
        return true
      } catch (error) {
        console.error('添加关系失败:', error)
        
        // 处理后端返回的400错误（关系已存在）
        if (error.response && error.response.status === 400) {
          const errorMessage = error.response.data && error.response.data.message
            ? error.response.data.message
            : `节点 "${relationData.sourceName}" 与 "${relationData.targetName}" 之间的关系已存在`;
          
          message.warning(errorMessage)
          return false
        } else {
          message.error('添加关系失败')
        }
        
        return false
      } finally {
        this.loading = false
      }
    },

    // 更新关系
    async updateRelation(relationData) {
      try {
        this.loading = true
        await updateRelation(relationData)
        if (response.code === 200) {
          message.success('更新关系成功')
          // 重新加载相关类型的关系
          await this.loadRelationsByTypes(relationData.sourceType, relationData.targetType)
          return true   
        }
        else {
          message.error(response.message ||'更新关系失败')
          return false
        }
      } catch (error) {
        message.error('更新关系失败')
        console.error('更新关系错误:', error)
        return false
      } finally {
        this.loading = false
      }
    },

    // 删除关系
    async deleteRelation(relationData) {
      try {
        this.loading = true
        
        // 确保关系类型参数存在
        if (!relationData.relationType) {
          console.error('删除关系时缺少关系类型参数', relationData)
          message.error('删除关系失败: 缺少关系类型参数')
          return false
        }
        
        const response = await deleteRelation(relationData)
        
        // 检查响应是否成功
        if (response.code === 200) {
          message.success('删除关系成功');

          // 重新加载相关类型的关系
          await this.loadRelationsByTypes(relationData.sourceType, relationData.targetType);

          return true;
        } else {
          message.error(response.data || '删除关系失败');
          return false;
        }
      } catch (error) {
        message.error('删除关系失败')
        console.error(error)
        return false
      } finally {
        this.loading = false
      }
        
        // 处理错误响应
        if (error.response && error.response.data) {
          const errorMsg = error.response.data.message || '删除关系失败'
          message.error(errorMsg)
        } else {
          message.error(`删除关系失败: ${error.message || '未知错误'}`)
        }
        
        // 尝试刷新关系列表，以防关系已被成功删除但返回了错误
        try {
          await this.loadRelationsByTypes(relationData.sourceType, relationData.targetType)
        } catch (refreshError) {
          console.error('刷新关系列表失败:', refreshError)
        } 
        

    },

    // 更新关系存储
    updateRelationsInStore(newRelations) {
      // 移除重复关系
      const existingIds = this.relations.map(r => r.id)
      const uniqueNewRelations = newRelations.filter(r => !existingIds.includes(r.id))
      
      // 添加新关系
      this.relations = [...this.relations, ...uniqueNewRelations]
    },

    // 加载特定类型节点之间的关系
    async loadRelationsByTypes(sourceType, targetType) {
      try {
        this.loading = true
        // 清空当前关系数据，避免显示过期数据
        this.relations = []
        
        console.log(`加载关系数据: ${sourceType} -> ${targetType}`)
        
        const response = await getRelationsByTypes({
          source_type: sourceType,
          target_type: targetType
        })
        
        // 处理响应数据
        if (!response) {
          console.error('关系API返回空响应')
          this.relations = []
          return
        }
        
        console.log('API返回数据:', response)
        
        // 提取实际数据
        let data
        
        if (response.data !== undefined) {
          // 处理标准响应格式
          if (typeof response.data === 'object' && response.data !== null) {
            // 检查是否是包含data字段的标准格式
            if (response.data.data !== undefined) {
              data = response.data.data
              console.log('使用response.data.data数据')
            } else {
              data = response.data
              console.log('使用response.data数据')
            }
          } else {
            data = response.data
            console.log('使用response.data(非对象)数据')
          }
        } else {
          // 直接使用response作为数据
          data = response
          console.log('使用整个response作为数据')
        }
        
        // 处理不同的数据格式
        if (Array.isArray(data)) {
          // 新的API格式 - 直接返回关系数组
          console.log(`获取到${data.length}个关系(数组格式)`)
          
          // 将API返回的关系数据映射到我们需要的格式
          this.relations = data.map(relation => ({
            id: `${relation.sourceName}-${relation.targetName}-${relation.relationType}`,
            sourceName: relation.sourceName,
            targetName: relation.targetName,
            sourceType: sourceType,
            targetType: targetType,
            relationType: relation.relationType
          }))
        } else if (data && data.relations && Array.isArray(data.relations)) {
          // 带relations字段的格式
          console.log(`获取到${data.relations.length}个关系(relations字段)`)
          
          this.relations = data.relations.map(relation => ({
            id: `${relation.source_name || relation.sourceName}-${relation.target_name || relation.targetName}-${relation.type || relation.relationType}`,
            sourceName: relation.source_name || relation.sourceName,
            targetName: relation.target_name || relation.targetName,
            sourceType: relation.source_type || sourceType,
            targetType: relation.target_type || targetType,
            relationType: relation.type || relation.relationType
          }))
        } else if (data && data.nodes && data.links) {
          // 图结构格式
          console.log(`获取到图结构数据: ${data.nodes.length}个节点, ${data.links.length}个连接`)
          
          const relations = data.links.map(link => {
            // 确保source和target是对象或字符串的情况都能处理
            const sourceId = typeof link.source === 'object' ? link.source.id : link.source
            const targetId = typeof link.target === 'object' ? link.target.id : link.target
            
            const sourceNode = data.nodes.find(node => node.id === sourceId)
            const targetNode = data.nodes.find(node => node.id === targetId)
            
            if (sourceNode && targetNode && 
                sourceNode.type === sourceType && 
                targetNode.type === targetType) {
              return {
                id: `${sourceNode.name}-${targetNode.name}-${link.type}`,
                sourceName: sourceNode.name,
                targetName: targetNode.name,
                sourceType: sourceNode.type,
                targetType: targetNode.type,
                relationType: link.type
              }
            }
            return null
          }).filter(Boolean)
          
          this.relations = relations
        } else {
          console.warn('未识别的关系数据格式:', data)
          this.relations = []
        }
        
        console.log('处理后的关系数据:', this.relations)
      } catch (error) {
        message.error(`加载${sourceType}到${targetType}的关系失败`)
        console.error('加载关系数据错误:', error)
        this.relations = []
      } finally {
        this.loading = false
      }
    }
  }
}) 