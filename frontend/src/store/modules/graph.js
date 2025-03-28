import { defineStore } from 'pinia'
import { 
  getChapterGraph, 
  getSectionGraph, 
  getPointGraph,
  getKnowledgePoint,
  getKnowledgeVideos,
  getKnowledgeExercises,
  getKnowledgeCourseware,
  searchKnowledge
} from '@/api/knowledge'
import { message } from 'ant-design-vue'

export const useGraphStore = defineStore('graph', {
  state: () => ({
    currentGraphData: null,
    currentChapter: null,
    currentSection: null,
    selectedPoint: null,
    loading: false,
    videos: [],
    exercises: []
  }),

  getters: {
    currentLevel: (state) => {
      if (state.currentSection) return 'point'
      if (state.currentChapter) return 'section'
      return 'chapter'
    }
  },

  actions: {
    async loadChapterGraph() {
      try {
        this.loading = true
        const data = await getChapterGraph()
        console.log('章节图谱数据:', data)
        
        // 检查响应数据结构
        if (!data) {
          console.error('API返回数据为空')
          throw new Error('加载章节图谱失败：返回数据为空')
        }
        
        // 处理可能的数据包装（检查是否嵌套在data字段中）
        let graphData = data
        if (data.data && (data.data.nodes || data.data.links)) {
          console.log('数据被包装在data字段中，提取内部数据')
          graphData = data.data
        }
        
        // 验证数据格式
        if (!graphData.nodes || !Array.isArray(graphData.nodes)) {
          console.error('API返回数据格式不正确，缺少nodes数组', graphData)
          throw new Error('加载章节图谱失败：数据格式不正确')
        }
        
        // 处理节点类型，确保所有节点都有正确的type属性
        graphData.nodes = graphData.nodes.map(node => {
          // 如果节点没有type属性，根据ID或其他属性推断类型
          if (!node.type) {
            node.type = 'chapter' // 默认为章节类型
          }
          return node
        })
        
        console.log('设置图谱数据，节点数量:', graphData.nodes.length)
        this.currentGraphData = graphData
        this.currentChapter = null
        this.currentSection = null
        this.selectedPoint = null
      } catch (error) {
        console.error('加载章节图谱失败:', error)
        message.error('加载章节图谱失败')
        this.currentGraphData = null
      } finally {
        this.loading = false
      }
    },

    async loadSectionGraph(chapterId) {
      try {
        this.loading = true
        const data = await getSectionGraph(chapterId)
        console.log('小节图谱数据:', data)
        
        // 检查响应数据结构
        if (!data) {
          console.error('API返回数据为空')
          throw new Error('加载小节图谱失败：返回数据为空')
        }
        
        // 处理可能的数据包装（检查是否嵌套在data字段中）
        let graphData = data
        if (data.data && (data.data.nodes || data.data.links)) {
          console.log('数据被包装在data字段中，提取内部数据')
          graphData = data.data
        }
        
        // 验证数据格式
        if (!graphData.nodes || !Array.isArray(graphData.nodes)) {
          console.error('API返回数据格式不正确，缺少nodes数组', graphData)
          throw new Error('加载小节图谱失败：数据格式不正确')
        }
        
        // 处理节点类型，确保所有节点都有正确的type属性
        graphData.nodes = graphData.nodes.map(node => {
          // 如果节点没有type属性，根据ID或其他属性推断类型
          if (!node.type) {
            node.type = 'section' // 默认为小节类型
          }
          return node
        })
        
        console.log('设置小节图谱数据，节点数量:', graphData.nodes.length)
        
        // 查找对应的章节信息
        const chapter = {
          id: chapterId,
          name: `章节 ${chapterId}`, // 默认名称，可能需要从其他地方获取
          type: 'chapter'
        }
        
        this.currentGraphData = graphData
        this.currentChapter = chapter
        this.currentSection = null
        this.selectedPoint = null
      } catch (error) {
        message.error('该章节未添加小节')
        console.error('加载小节图谱失败:', error)
        // 保持当前数据不变，或设置为null取决于您的需求
      } finally {
        this.loading = false
      }
    },

    async loadPointGraph(sectionId) {
      try {
        this.loading = true
        const data = await getPointGraph(sectionId)
        console.log('知识点图谱数据:', data)
        
        // 检查响应数据结构
        if (!data) {
          console.error('API返回数据为空')
          throw new Error('加载知识点图谱失败：返回数据为空')
        }
        
        // 处理可能的数据包装（检查是否嵌套在data字段中）
        let graphData = data
        if (data.data && (data.data.nodes || data.data.links)) {
          console.log('数据被包装在data字段中，提取内部数据')
          graphData = data.data
        }
        
        // 验证数据格式
        if (!graphData.nodes || !Array.isArray(graphData.nodes)) {
          console.error('API返回数据格式不正确，缺少nodes数组', graphData)
          throw new Error('加载知识点图谱失败：数据格式不正确')
        }
        
        // 处理节点类型，确保所有节点都有正确的type属性
        graphData.nodes = graphData.nodes.map(node => {
          // 如果节点没有type属性，根据ID或其他属性推断类型
          if (!node.type) {
            node.type = 'point' // 默认为知识点类型
          }
          return node
        })
        
        console.log('设置知识点图谱数据，节点数量:', graphData.nodes.length)
        
        // 查找对应的小节信息
        const section = {
          id: sectionId,
          name: `小节 ${sectionId}`, // 默认名称，可能需要从其他地方获取
          type: 'section'
        }
        
        this.currentGraphData = graphData
        this.currentSection = section
        this.selectedPoint = null
      } catch (error) {
        message.error('此小节未添加知识点')
        console.error('加载知识点图谱失败:', error)
        // 保持当前数据不变
      } finally {
        this.loading = false
      }
    },

    async loadPointDetails(pointId) {
      try {
        this.loading = true
        console.log('加载知识点详情:', pointId)
        
        // 获取相关视频
        const videosResponse = await getKnowledgeVideos(pointId)
        console.log('相关视频:', videosResponse)
        
        let videos = []
        if (videosResponse.data) {
          // 为每个视频添加url字段，方便前端统一处理
          videos = videosResponse.data.map(video => ({
            ...video,
            url: video.play_url // 添加url字段映射到play_url
          }))
        }
        
        // 获取相关练习
        const exercisesResponse = await getKnowledgeExercises(pointId)
        console.log('相关练习:', exercisesResponse)
        
        let exercises = []
        if (exercisesResponse.data) {
          // 为每个练习添加url字段，方便前端统一处理
          exercises = exercisesResponse.data.map(exercise => ({
            ...exercise,
            url: exercise.exercise_url // 添加url字段映射到exercise_url
          }))
        }
        
        // 获取相关课件
        const coursewareResponse = await getKnowledgeCourseware(pointId)
        console.log('相关课件:', coursewareResponse)
        
        let courseware = []
        if (coursewareResponse.data) {
          // 为每个课件添加url字段，方便前端统一处理
          courseware = coursewareResponse.data.map(cw => ({
            ...cw,
            url: cw.courseware_url // 添加url字段映射到courseware_url
          }))
        }
        
        // 构建知识点详情
        const currentNode = this.currentGraphData.nodes.find(node => node.id === pointId)
        if (!currentNode) {
          throw new Error('未找到对应的知识点')
        }
        
        this.selectedPoint = {
          ...currentNode,
          videos,
          exercises,
          courseware
        }
        
        console.log('知识点详情加载完成:', this.selectedPoint)
        return this.selectedPoint
      } catch (error) {
        message.error('加载知识点详情失败: ' + (error.message || '未知错误'))
        console.error('加载知识点详情失败:', error)
        return null
      } finally {
        this.loading = false
      }
    },

    async searchNodes(keyword) {
      try {
        const resultsData = await searchKnowledge(keyword)
        console.log('搜索结果:', resultsData)
        
        let results = resultsData
        if (resultsData.data) {
          results = resultsData.data
        }
        
        return results || []
      } catch (error) {
        message.error('搜索失败')
        console.error(error)
        return []
      }
    },

    clearSelection() {
      this.selectedPoint = null
      this.videos = []
      this.exercises = []
    }
  }
}) 