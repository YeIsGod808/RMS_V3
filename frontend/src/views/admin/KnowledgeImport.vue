<template>
  <div class="knowledge-import-container">
    <a-card title="知识图谱一键导入" class="import-card">
      <p class="description">上传JSON格式的知识图谱数据文件，系统将自动构建知识图谱。</p>
      
      <a-upload
        name="file"
        :multiple="false"
        :before-upload="beforeUpload"
        :custom-request="customUploadRequest"
        :show-upload-list="false"
        class="upload-area"
      >
        <div class="upload-content" :class="{ 'is-dragover': isDragover }" 
             @dragover="isDragover = true" 
             @dragleave="isDragover = false">
          <p class="upload-icon">
            <upload-outlined />
          </p>
          <p class="upload-text">点击或拖拽文件到此区域上传</p>
          <p class="upload-hint">支持上传JSON格式文件</p>
        </div>
      </a-upload>
      
      <div class="upload-status" v-if="uploadStatus !== ''">
        <a-alert
          :message="uploadStatus"
          :type="uploadStatusType"
          show-icon
          class="status-alert"
        />
      </div>
      
      <div class="upload-progress" v-if="uploading">
        <a-progress :percent="uploadPercent" status="active" />
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { UploadOutlined } from '@ant-design/icons-vue'
import axios from 'axios'
import { autoConstruct } from '@/api/analysis'

// 上传状态
const uploading = ref(false)
const uploadPercent = ref(0)
const uploadStatus = ref('')
const uploadStatusType = ref('info')
const isDragover = ref(false)

// 上传前验证文件
const beforeUpload = (file) => {
  const isJSON = file.type === 'application/json' || file.name.endsWith('.json')
  if (!isJSON) {
    message.error('只能上传JSON格式的文件!')
    return false
  }
  
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    message.error('文件大小不能超过10MB!')
    return false
  }
  
  return true
}

// 自定义上传请求
const customUploadRequest = async ({ file }) => {
  uploading.value = true
  uploadStatus.value = '正在上传文件...'
  uploadStatusType.value = 'info'
  uploadPercent.value = 0

  try {
    // 创建FormData对象
    const formData = new FormData()
    formData.append('file', file)
    // 遍历 FormData 内容
    for (let [key, value] of formData.entries()) {
      console.log(key, value);
    }

    // 发送带FormData的请求
    const response = await autoConstruct(formData)

    // 上传成功
    uploadStatus.value = '知识图谱导入成功!'
    uploadStatusType.value = 'success'
    message.success('知识图谱导入成功!')
  } catch (error) {
    // 上传失败
    console.error('上传失败:', error)
    uploadStatus.value = `上传失败: ${error.response?.data?.message || error.message || '未知错误'}`
    uploadStatusType.value = 'error'
    message.error('知识图谱导入失败!')
  } finally {
    uploading.value = false
    uploadPercent.value = 100
    
    // 5秒后清除状态
    setTimeout(() => {
      if (uploadStatusType.value === 'success') {
        uploadStatus.value = ''
      }
    }, 5000)
  }
}
</script>

<style scoped>
.knowledge-import-container {
  width: 100%;
}

.import-card {
  margin-bottom: 24px;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.description {
  margin-bottom: 20px;
  color: #5c6b7c;
  font-size: 14px;
}

.upload-area {
  width: 100%;
  margin-bottom: 20px;
}

.upload-content {
  padding: 40px 0;
  text-align: center;
  background: rgba(240, 247, 255, 0.5);
  border: 2px dashed #d9e6f2;
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.3s ease;
}

.upload-content:hover, .upload-content.is-dragover {
  border-color: #2196F3;
  background: rgba(33, 150, 243, 0.05);
}

.upload-icon {
  font-size: 48px;
  color: #2196F3;
  margin-bottom: 16px;
}

.upload-text {
  font-size: 16px;
  color: #333;
  margin-bottom: 8px;
}

.upload-hint {
  font-size: 14px;
  color: #999;
}

.upload-status {
  margin-top: 16px;
}

.status-alert {
  margin-bottom: 16px;
}

.upload-progress {
  margin-top: 16px;
}
</style>