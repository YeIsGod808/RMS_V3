<template>
  <div class="resource-management">
    <a-card class="resource-card">
      <template #extra>
        <a-space>
          <a-select v-model:value="resourceType" style="width: 120px">
            <a-select-option value="video">视频</a-select-option>
            <a-select-option value="courseware">课件</a-select-option>
            <a-select-option value="exercise">练习题</a-select-option>
          </a-select>
          <a-button type="primary" class="create-button" @click="showAddResourceModal">
            <plus-outlined /> 添加资源
          </a-button>
        </a-space>
      </template>
      
      <!-- 资源筛选 -->
      <a-row :gutter="16" style="margin-bottom: 16px">
        <a-col :span="24">
          <a-space>
            <a-select 
              v-model:value="selectedChapter" 
              style="width: 200px" 
              placeholder="选择章节"
              @change="handleChapterChange"
              allowClear
            >
              <a-select-option v-for="chapter in knowledgeAdminStore.chapterOptions" :key="chapter.value" :value="chapter.value">
                {{ chapter.label }}
              </a-select-option>
            </a-select>
            
            <a-select 
              v-model:value="selectedSection" 
              style="width: 200px" 
              placeholder="选择小节"
              @change="handleSectionChange"
              :disabled="!selectedChapter"
              allowClear
            >
              <a-select-option v-for="section in knowledgeAdminStore.sectionOptions" :key="section.value" :value="section.value">
                {{ section.label }}
              </a-select-option>
            </a-select>
            
            <a-select 
              v-model:value="selectedPoint" 
              style="width: 200px" 
              placeholder="选择知识点"
              @change="handlePointChange"
              :disabled="!selectedSection"
              allowClear
            >
              <a-select-option v-for="point in knowledgeAdminStore.pointOptions" :key="point.value" :value="point.value">
                {{ point.label }}
              </a-select-option>
            </a-select>
          </a-space>
        </a-col>
      </a-row>
      
      <a-table
        :columns="resourceColumns"
        :data-source="resources"
        :loading="loading"
        row-key="id"
        class="resource-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'">
            <a-tag :color="getResourceTypeColor(record.type)" class="resource-type-tag">
              {{ getResourceTypeName(record.type) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'difficulty' && record.difficulty">
            <a-tag :color="getDifficultyColor(record.difficulty)" class="difficulty-tag">
              {{ getDifficultyText(record.difficulty) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'created_at' && record.created_at">
            {{ formatDateTime(record.created_at) }}
          </template>
          <template v-if="column.key === 'updated_at' && record.updated_at">
            {{ formatDateTime(record.updated_at) }}
          </template>
          <template v-if="column.key === 'preview'">
            <a-button 
              type="link" 
              class="preview-button"
              @click="handlePreviewResource(record)" 
              v-if="record.type === 'video' ? record.play_url : 
                    record.type === 'courseware' ? record.courseware_url : 
                    record.type === 'exercise' ? record.exercise_url : 
                    record.url"
            >
              预览
            </a-button>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-popconfirm
                title="确定要删除该资源吗？"
                @confirm="handleDeleteResource(record)"
              >
                <a-button type="link" danger class="delete-button">删除</a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 资源表单模态框 -->
    <a-modal
      v-model:visible="resourceModalVisible"
      :title="resourceModalTitle"
      @ok="handleResourceModalOk"
      :confirm-loading="modalLoading"
      :width="700"
      class="resource-modal"
    >
      <a-form
        ref="resourceFormRef"
        :model="resourceForm"
        :rules="resourceRules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
      >
        <a-form-item label="标题" name="title">
          <a-input v-model:value="resourceForm.title" />
        </a-form-item>
        
        <a-form-item label="资源类型" name="resourceType">
          <a-select v-model:value="resourceForm.resourceType">
            <a-select-option value="video">视频</a-select-option>
            <a-select-option value="courseware">课件</a-select-option>
            <a-select-option value="exercise">练习题</a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="知识点" name="pointName">
          <a-select 
            v-model:value="resourceForm.pointName"
            show-search
            :filter-option="filterOption"
            placeholder="请选择知识点"
          >
            <a-select-option v-for="point in knowledgeAdminStore.allPointOptions" :key="point.value" :value="point.value">
              {{ point.label }}
            </a-select-option>
          </a-select>
        </a-form-item>
        
        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="resourceForm.description" :rows="4" />
        </a-form-item>
        
        <template v-if="resourceForm.resourceType === 'exercise'">
          <a-form-item label="难度" name="difficulty">
            <a-select v-model:value="resourceForm.difficulty">
              <a-select-option value="easy">简单</a-select-option>
              <a-select-option value="medium">中等</a-select-option>
              <a-select-option value="hard">困难</a-select-option>
            </a-select>
          </a-form-item>
        </template>
        
        <a-form-item label="上传方式" name="uploadType">
          <a-radio-group v-model:value="uploadType">
            <a-radio value="file">文件上传</a-radio>
            <a-radio value="link">链接上传</a-radio>
          </a-radio-group>
        </a-form-item>
        
        <a-form-item label="资源链接" name="resourceLink" v-if="uploadType === 'link'">
          <a-input v-model:value="resourceForm.resourceLink" placeholder="请输入资源链接URL" />
        </a-form-item>
        
        <a-form-item label="上传文件" name="file" v-if="uploadType === 'file'">
          <a-upload
            v-model:file-list="fileList"
            :before-upload="beforeUpload"
            :remove="handleRemove"
            :max-count="1"
          >
            <a-button>
              <upload-outlined /> 选择文件
            </a-button>
            <template #itemRender="{ file }">
              <a-space>
                <file-outlined />
                <span>{{ file.name }}</span>
                <a-button type="link" danger @click="() => handleRemove(file)">删除</a-button>
              </a-space>
            </template>
          </a-upload>
          <div style="color: #999; font-size: 12px; margin-top: 8px;">
            <template v-if="resourceForm.resourceType === 'video'">
              支持的视频格式: MP4, AVI, MOV, WMV, MKV
            </template>
            <template v-else-if="resourceForm.resourceType === 'courseware'">
              支持的课件格式: PDF, PPTX, DOC, DOCX
            </template>
            <template v-else-if="resourceForm.resourceType === 'exercise'">
              支持的练习题格式: PDF, DOCX
            </template>
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 资源预览模态框 -->
    <a-modal
      v-model:visible="previewModalVisible"
      :title="previewModalTitle"
      :footer="null"
      :width="800"
      @cancel="handlePreviewModalClose"
      class="preview-modal"
    >
      <template v-if="previewResource && previewResource.url">
        <template v-if="isVideoUrl(previewResource.url) && previewResource.cover_url">
          <video 
            controls 
            style="width: 100%; max-height: 500px;" 
            :src="previewResource.url"
            :poster="previewResource.cover_url"
            ref="videoPlayer"
          ></video>
        </template>
        <template v-else-if="isPdfUrl(previewResource.url)">
          <div style="height: 500px; width: 100%;">
            <iframe 
              :src="previewResource.url" 
              style="width: 100%; height: 100%; border: none;"
            ></iframe>
          </div>
        </template>
        <template v-else>
          <div style="text-align: center; padding: 20px;">
            <p>无法直接预览该资源，请点击下方按钮跳转查看</p>
            <a :href="previewResource.url" target="_blank">
              <a-button type="primary">跳转查看</a-button>
            </a>
          </div>
        </template>
      </template>
      <template v-else>
        <div style="text-align: center; padding: 20px;">
          <a-empty description="资源URL不存在" />
        </div>
      </template>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch, computed } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, UploadOutlined, FileOutlined } from '@ant-design/icons-vue'
import { useAdminStore } from '@/store/modules/admin'
import { useKnowledgeAdminStore } from '@/store/modules/knowledgeAdmin'
import { uploadResource, getVideosByPointId, getExercisesByPointId, getCoursewaresByPointId, 
         deleteVideo, deleteExercise, deleteCourseware } from '@/api/resource'

const adminStore = useAdminStore()
const knowledgeAdminStore = useKnowledgeAdminStore()

// 表格列定义
const resourceColumns = [
  { title: '编号', dataIndex: 'ID', width: 80 },
  { title: '标题', dataIndex: 'title' },
  { title: '类型', dataIndex: 'type', key: 'type', width: 100 },
  { title: '难度', dataIndex: 'difficulty', key: 'difficulty', width: 100 },
  { title: '知识点', dataIndex: 'pointName' },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '创建时间', dataIndex: 'created_at', key: 'created_at', width: 170 },
  { title: '更新时间', dataIndex: 'updated_at', key: 'updated_at', width: 170 },
  { title: '预览', key: 'preview', width: 80 },
  { title: '操作', key: 'action', width: 100 }
]

// 状态管理
const resourceType = ref('video')
const modalLoading = ref(false)
const loading = ref(false)
const resources = ref([])
const selectedChapter = ref(null)
const selectedSection = ref(null)
const selectedPoint = ref(null)
const selectedPointId = ref(null)
const uploadType = ref('file')

// 模态框状态
const resourceModalVisible = ref(false)
const resourceModalTitle = ref('')
const previewModalVisible = ref(false)
const previewModalTitle = ref('')
const previewResource = ref(null)
const videoPlayer = ref(null)

// 表单引用
const resourceFormRef = ref()
const fileList = ref([])

// 表单数据
const resourceForm = reactive({
  title: '',
  resourceType: 'video',
  pointName: '',
  description: '',
  difficulty: 'medium',
  file: null,
  resourceLink: ''
})

// 表单验证规则
const resourceRules = {
  title: [{ required: true, message: '请输入标题' }],
  resourceType: [{ required: true, message: '请选择资源类型' }],
  pointName: [{ required: true, message: '请选择知识点' }],
  difficulty: [{ required: true, message: '请选择难度', trigger: 'change' }],
  resourceLink: [{ 
    required: false,
    message: '请输入有效的URL',
    type: 'url',
    trigger: 'blur' 
  }]
}

// 获取资源类型颜色
const getResourceTypeColor = (type) => {
  switch (type) {
    case 'video': return 'blue'
    case 'courseware': return 'green'
    case 'exercise': return 'orange'
    default: return 'default'
  }
}

// 获取资源类型名称
const getResourceTypeName = (type) => {
  switch (type) {
    case 'video': return '视频'
    case 'courseware': return '课件'
    case 'exercise': return '练习题'
    default: return type
  }
}

// 获取难度颜色
const getDifficultyColor = (difficulty) => {
  switch (difficulty) {
    case 'easy': return 'green'
    case 'medium': return 'blue'
    case 'hard': return 'red'
    default: return 'default'
  }
}

// 获取难度文本
const getDifficultyText = (difficulty) => {
  switch (difficulty) {
    case 'easy': return '简单'
    case 'medium': return '中等'
    case 'hard': return '困难'
    default: return '未知'
  }
}

// 事件处理函数
const handleDeleteResource = async (record) => {
  try {
    if (!selectedPointId.value) {
      message.error('无法获取知识点ID，删除失败')
      return
    }
    
    
    if (!record.ID) {
      message.error('无法获取资源ID，删除失败')
      return
    }
    
    if (record.type === 'video') {
      await deleteVideo(record.ID, selectedPointId.value)
    } else if (record.type === 'exercise') {
      await deleteExercise(record.ID, selectedPointId.value)
    } else if (record.type === 'courseware') {
      await deleteCourseware(record.ID, selectedPointId.value)
    }
    
    message.success('删除资源成功')
    loadResources()
  } catch (error) {
    console.error('删除资源失败:', error)
    message.error('删除资源失败: ' + (error.message || '未知错误'))
  }
}

// 模态框相关方法
const showAddResourceModal = () => {
  resourceForm.title = ''
  resourceForm.resourceType = resourceType.value
  resourceForm.pointName = selectedPoint.value || ''
  resourceForm.description = ''
  resourceForm.difficulty = 'medium'
  resourceForm.resourceLink = ''
  
  fileList.value = []
  uploadType.value = 'file'
  resourceModalTitle.value = '添加资源'
  resourceModalVisible.value = true
}

// 处理资源模态框确认
const handleResourceModalOk = async () => {
  try {
    await resourceFormRef.value.validate()
    modalLoading.value = true
    
    // 添加资源
    const formData = new FormData()
    
    // 添加基础字段
    formData.append('title', resourceForm.title)
    formData.append('resource_type', resourceForm.resourceType)
    formData.append('point_name', resourceForm.pointName)
    
    if (resourceForm.description) {
      formData.append('description', resourceForm.description)
    }
    
    if (resourceForm.resourceType === 'exercise' && resourceForm.difficulty) {
      formData.append('difficulty', resourceForm.difficulty)
    }
    
    // 根据上传方式处理
    if (uploadType.value === 'link') {
      // 链接上传
      if (!resourceForm.resourceLink) {
        message.error('请输入资源链接')
        modalLoading.value = false
        return
      }
      
      // 添加链接
      formData.append('resource_link', resourceForm.resourceLink)
      console.log('准备上传链接:', resourceForm.resourceLink)
    } else {
      // 文件上传
      if (fileList.value.length === 0) {
        message.error('请选择要上传的文件')
        modalLoading.value = false
        return
      }
      
      
      // 获取文件对象
      const file = fileList.value[0]
      const fileObj = file.originFileObj || file
      

      formData.append('data', fileObj)
      
      console.log('准备上传文件:', fileObj.name)
    }
    
    // // 调试输出FormData内容
    // console.log('FormData内容:')
    // for (let [key, value] of formData.entries()) {
    //   console.log(`${key}: ${value instanceof File ? value.name : value}`)
    // }
    
    try {
      const response = await uploadResource(formData)
      console.log('上传响应:', response)
      
      if (response.code === 200) {
        message.success('上传资源成功')
        resourceModalVisible.value = false
        loadResources()
      } else {
        message.error('上传资源失败: ' + (response.message || '未知错误'))
      }
    } catch (error) {
      console.error('上传资源失败:', error)
      if (error.response && error.response.data) {
        console.error('错误详情:', error.response.data)
        message.error('上传资源失败: ' + (error.response.data.message || error.message || '未知错误'))
      } else {
        message.error('上传资源失败: ' + (error.message || '未知错误'))
      }
    }
  } catch (error) {
    console.error('表单验证失败:', error)
    message.error('表单验证失败，请检查输入')
  } finally {
    modalLoading.value = false
  }
}

// 上传文件相关方法
const beforeUpload = (file) => {
  const isValidType = checkFileType(file, resourceForm.resourceType)
  if (!isValidType) {
    message.error('文件类型不支持！')
    return false
  }
  
  const isLt100M = file.size / 1024 / 1024 < 100
  if (!isLt100M) {
    message.error('文件大小不能超过100MB！')
    return false
  }
  
  fileList.value = [file]
  
  resourceForm.file = file
  
  console.log('文件已选择:', file.name, '大小:', file.size, '类型:', file.type)
  return false
}

const checkFileType = (file, resourceType) => {
  const fileName = file.name.toLowerCase()
  
  if (resourceType === 'video') {
    return fileName.endsWith('.mp4') || fileName.endsWith('.avi') || 
           fileName.endsWith('.mov') || fileName.endsWith('.wmv') || fileName.endsWith('.mkv')
  } else if (resourceType === 'courseware') {
    return fileName.endsWith('.pdf') || fileName.endsWith('.pptx') || 
           fileName.endsWith('.doc') || fileName.endsWith('.docx')
  } else if (resourceType === 'exercise') {
    return fileName.endsWith('.pdf') || fileName.endsWith('.docx')
  }
  
  return false
}

const handleRemove = (file) => {
  fileList.value = []
  resourceForm.file = null
  return true
}

const handleChapterChange = async (value) => {
  selectedChapter.value = value
  selectedSection.value = null
  selectedPoint.value = null
  selectedPointId.value = null
  resources.value = [] // 清空资源列表
  
  if (!value) {
    return
  }
  
  const chapter = knowledgeAdminStore.chapters.find(c => c.name === value)
  if (chapter) {
    await knowledgeAdminStore.loadSectionsByChapterId(chapter.id)
  }
}

const handleSectionChange = async (value) => {
  selectedSection.value = value
  selectedPoint.value = null
  selectedPointId.value = null
  resources.value = [] // 清空资源列表
  
  if (!value) {
    return
  }
  
  const section = knowledgeAdminStore.sections.find(s => s.name === value)
  if (section) {
    await knowledgeAdminStore.loadPointsBySectionId(section.id)
  }
}

const handlePointChange = async (value) => {
  selectedPoint.value = value
  
  if (!value) {
    selectedPointId.value = null
    resources.value = []
    return
  }
  
  const point = knowledgeAdminStore.allPoints.find(p => p.name === value)
  if (point) {
    selectedPointId.value = point.id
    await loadResources()
  }
}

const loadResources = async () => {
  if (!selectedPointId.value) {
    resources.value = []
    return
  }
  
  loading.value = true
  try {
    let data = []
    
    if (resourceType.value === 'video') {
      const response = await getVideosByPointId(selectedPointId.value)
      data = response.data || []
    } else if (resourceType.value === 'exercise') {
      const response = await getExercisesByPointId(selectedPointId.value)
      data = response.data || []
    } else if (resourceType.value === 'courseware') {
      const response = await getCoursewaresByPointId(selectedPointId.value)
      data = response.data || []
    }
    
    resources.value = data.map(item => ({
      ...item,
      type: resourceType.value,
      pointName: selectedPoint.value
    }))
  } catch (error) {
    console.error('加载资源失败:', error)
    message.error('加载资源失败')
    resources.value = []
  } finally {
    loading.value = false
  }
}

const handlePreviewResource = (resource) => {
  console.log('预览资源:', resource);
  
  let resourceUrl;
  if (resource.type === 'video') {
    resourceUrl = resource.play_url || resource.url;
  } else if (resource.type === 'courseware') {
    resourceUrl = resource.courseware_url || resource.url;
  } else if (resource.type === 'exercise') {
    resourceUrl = resource.exercise_url || resource.url;
  } else {
    resourceUrl = resource.url;
  }
  
  if (!resourceUrl) {
    message.error('资源URL不存在，无法预览');
    console.error('资源URL不存在:', resource);
    return;
  }
  
  let url = resourceUrl;
  if (url && !url.startsWith('http://') && !url.startsWith('https://')) {
    const baseUrl = window.location.origin;
    url = url.startsWith('/') ? `${baseUrl}${url}` : `${baseUrl}/${url}`;
    console.log('完整URL:', url);
  }
  
  
  previewResource.value = {
    ...resource,
    url: url,
  };
  
  previewModalTitle.value = resource.title || '资源预览';
  previewModalVisible.value = true;
}

const isVideoUrl = (url) => {
  if (!url) return false;
  const lowerUrl = url.toLowerCase();
  return lowerUrl.endsWith('.mp4') || lowerUrl.endsWith('.avi') || 
         lowerUrl.endsWith('.mov') || lowerUrl.endsWith('.wmv') ||
         lowerUrl.includes('video') || lowerUrl.includes('mkv');
}

const isPdfUrl = (url) => {
  if (!url) return false;
  const lowerUrl = url.toLowerCase();
  return lowerUrl.endsWith('.pdf') || lowerUrl.includes('pdf');
}

const filterOption = (input, option) => {
  return option.value.toLowerCase().indexOf(input.toLowerCase()) >= 0;
}

watch(resourceType, async () => {
  if (selectedPointId.value) {
    await loadResources()
  }
})

onMounted(async () => {
  await knowledgeAdminStore.loadChapters()
  await knowledgeAdminStore.loadAllSections()
  await knowledgeAdminStore.loadAllPoints()
})

const handlePreviewModalClose = () => {
  if (videoPlayer.value && isVideoUrl(previewResource.value?.url)) {
    videoPlayer.value.pause()
    videoPlayer.value.currentTime = 0
  }
  
  previewResource.value = null
  previewModalVisible.value = false
}

// 格式化日期时间
const formatDateTime = (dateTimeStr) => {
  if (!dateTimeStr) return '-';
  
  // 尝试转换时间字符串为日期对象
  const date = new Date(dateTimeStr);
  
  // 检查日期是否有效
  if (isNaN(date.getTime())) {
    return dateTimeStr; // 无效日期则返回原始字符串
  }
  
  // 格式化为 YYYY-MM-DD HH:MM:SS
  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, '0');
  const day = String(date.getDate()).padStart(2, '0');
  const hours = String(date.getHours()).padStart(2, '0');
  const minutes = String(date.getMinutes()).padStart(2, '0');
  const seconds = String(date.getSeconds()).padStart(2, '0');
  
  return `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
}
</script>

<style scoped>
.resource-management {
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  0% {
    opacity: 0;
    transform: translateY(10px);
  }
  100% {
    opacity: 1;
    transform: translateY(0);
  }
}

.resource-card {
  background: rgba(255, 255, 255, 0.95);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

:deep(.ant-card-head) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  padding: 16px 24px;
}

:deep(.ant-card-head-title) {
  font-size: 20px;
  font-weight: 600;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
    /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

.create-button {
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  transition: all 0.3s ease;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
}

.create-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
}

.resource-table {
  margin-top: 16px;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

:deep(.ant-table) {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
}

:deep(.ant-table-thead > tr > th) {
  background: linear-gradient(120deg, rgba(33, 150, 243, 0.08), rgba(76, 175, 80, 0.08));
  font-weight: 600;
  color: #1976D2;
  border-bottom: 1px solid rgba(33, 150, 243, 0.1);
  padding: 14px 16px;
  font-size: 14px;
  letter-spacing: 0.3px;
}

:deep(.ant-table-tbody > tr > td) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
  padding: 12px 16px;
  font-size: 14px;
  color: #333;
}

:deep(.ant-table-tbody > tr:hover > td) {
  background: rgba(33, 150, 243, 0.05);
}

:deep(.ant-table-tbody > tr:nth-child(even) > td) {
  background-color: rgba(250, 250, 250, 0.8);
}

:deep(.ant-table-tbody > tr:nth-child(even):hover > td) {
  background-color: rgba(33, 150, 243, 0.05);
}

.resource-type-tag, .difficulty-tag {
  border-radius: 12px;
  padding: 2px 12px;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.preview-button, .delete-button {
  border-radius: 6px;
  transition: all 0.3s ease;
  padding: 4px 8px;
}

.preview-button:hover {
  background: rgba(33, 150, 243, 0.1);
  transform: translateY(-1px);
}

.delete-button:hover {
  background: rgba(244, 67, 54, 0.1);
  transform: translateY(-1px);
}

.resource-modal, .preview-modal {
  border-radius: 16px;
  overflow: hidden;
}

:deep(.ant-modal-content) {
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

:deep(.ant-modal-header) {
  padding: 20px 24px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

:deep(.ant-modal-title) {
  font-size: 18px;
  font-weight: 600;
  color: #2196F3;
}

:deep(.ant-modal-body) {
  padding: 24px;
}

:deep(.ant-form-item-label > label) {
  color: #5c6b7c;
  font-weight: 500;
}

:deep(.ant-input),
:deep(.ant-select-selector),
:deep(.ant-radio-wrapper),
:deep(.ant-textarea) {
  border-radius: 10px;
  transition: all 0.3s ease;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05);
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
}

:deep(.ant-input:focus),
:deep(.ant-select-focused .ant-select-selector),
:deep(.ant-textarea:focus) {
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2), 0 2px 8px rgba(33, 150, 243, 0.1);
  border-color: #2196F3;
  transform: translateY(-1px);
}

:deep(.ant-select-item-option-selected) {
  background-color: rgba(33, 150, 243, 0.1);
  color: #2196F3;
  font-weight: 500;
}

:deep(.ant-select-item) {
  padding: 8px 12px;
  border-radius: 6px;
  margin: 2px 6px;
  transition: all 0.2s ease;
}

:deep(.ant-select-item:hover) {
  background-color: rgba(33, 150, 243, 0.05);
}

:deep(.ant-select) {
  .ant-select-selector {
    border-radius: 10px !important;
    padding: 4px 12px !important;
    height: 40px !important;
    display: flex;
    align-items: center;
    border: 1px solid rgba(0, 0, 0, 0.1) !important;
    background: #fff;
    box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
    transition: all 0.3s ease;

    .ant-select-selection-search-input {
      height: 38px !important;
    }

    .ant-select-selection-placeholder {
      color: #8c9bab;
      line-height: 38px !important;
    }

    .ant-select-selection-item {
      line-height: 38px !important;
      color: #2c3e50;
      font-weight: 500;
    }
  }

  &:hover .ant-select-selector {
    border-color: #4CAF50 !important;
    box-shadow: 0 2px 8px rgba(76, 175, 80, 0.1);
    transform: translateY(-1px);
  }

  &.ant-select-focused .ant-select-selector {
    border-color: #2196F3 !important;
    box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2), 0 4px 12px rgba(33, 150, 243, 0.1) !important;
    transform: translateY(-1px);
  }

  &.ant-select-disabled .ant-select-selector {
    background-color: #f5f7fa !important;
    border-color: #e4e7ed !important;
    box-shadow: none !important;
    cursor: not-allowed;

    .ant-select-selection-item {
      color: #909399;
    }
  }
}

:deep(.ant-select-dropdown) {
  padding: 8px;
  border-radius: 12px;
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.08);
  
  .ant-select-item {
    padding: 8px 12px;
    border-radius: 8px;
    margin: 2px 0;
    transition: all 0.2s ease;
    
    &:hover {
      background-color: rgba(33, 150, 243, 0.05);
    }
    
    &.ant-select-item-option-selected {
      background-color: rgba(33, 150, 243, 0.1);
      color: #2196F3;
      font-weight: 500;
    }
  }
}

:deep(.ant-radio-checked .ant-radio-inner) {
  border-color: #2196F3;
  background-color: #2196F3;
}

:deep(.ant-upload-list-item) {
  border-radius: 8px;
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary) {
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
}
</style>