import { defineStore } from 'pinia'
import { 
  getUsers, addUser, updateUser, deleteUser, resetPassword,
  getContents, addContent, updateContent, deleteContent,
  getResources, addResource
} from '@/api/admin'
import { deleteVideo, deleteExercise, deleteCourseware } from '@/api/resource'
import { message } from 'ant-design-vue'

export const useAdminStore = defineStore('admin', {
  state: () => ({
    users: [],
    contents: [],
    resources: [],
    loading: false
  }),

  actions: {
    async loadUsers() {
      try {
        this.loading = true
        const data = await getUsers()
        this.users = data
      } catch (error) {
        message.error('加载用户列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },

    async loadContents(type) {
      try {
        this.loading = true
        const data = await getContents(type)
        this.contents = data
      } catch (error) {
        message.error('加载内容列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },

    async loadResources(type) {
      try {
        this.loading = true
        const data = await getResources(type)
        this.resources = data
      } catch (error) {
        message.error('加载资源列表失败')
        console.error(error)
      } finally {
        this.loading = false
      }
    },

    // 内容管理
    async addContent(contentData) {
      try {
        await addContent(contentData)
        message.success('添加内容成功')
        await this.loadContents(contentData.type)
      } catch (error) {
        message.error('添加内容失败')
        throw error
      }
    },

    async updateContent(id, contentData) {
      try {
        await updateContent(id, contentData)
        message.success('更新内容成功')
        await this.loadContents(contentData.type)
      } catch (error) {
        message.error('更新内容失败')
        throw error
      }
    },

    async deleteContent(id, type) {
      try {
        await deleteContent(id)
        message.success('删除内容成功')
        await this.loadContents(type)
      } catch (error) {
        message.error('删除内容失败')
        throw error
      }
    },

    // 资源管理
    async addResource(resourceData) {
      try {
        await addResource(resourceData)
        message.success('添加资源成功')
        await this.loadResources(resourceData.type)
      } catch (error) {
        message.error('添加资源失败')
        throw error
      }
    },

    async deleteResource(id, type) {
      try {
        let response
        if (type === 'video') {
          response = await deleteVideo(id, this.selectedPointId)
        } else if (type === 'exercise') {
          response = await deleteExercise(id, this.selectedPointId)
        } else if (type === 'courseware') {
          response = await deleteCourseware(id, this.selectedPointId)
        } else {
          throw new Error(`不支持的资源类型: ${type}`)
        }
        
        return response
      } catch (error) {
        console.error('删除资源失败:', error)
        throw error
      }
    }
  }
}) 