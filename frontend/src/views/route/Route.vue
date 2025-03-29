<template>
  <div class="route-container">
    <a-card class="route-card">
      <div class="route-header">
        <h2>学习路径规划</h2>
        <p class="route-description">选择起点和终点，系统将为您推荐最佳的学习路径</p>
      </div>

      <div class="route-form">
        <a-form layout="vertical">
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="起点类型">
                <a-select v-model:value="startPointType" placeholder="请选择起点类型">
                  <a-select-option value="chapter">章节</a-select-option>
                  <a-select-option value="section">小节</a-select-option>
                  <a-select-option value="point">知识点</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="起点名称">
                <a-input v-model:value="startPointName" placeholder="请输入起点名称" />
              </a-form-item>
            </a-col>
          </a-row>
          <a-row :gutter="24">
            <a-col :span="12">
              <a-form-item label="终点类型">
                <a-select v-model:value="endPointType" placeholder="请选择终点类型">
                  <a-select-option value="chapter">章节</a-select-option>
                  <a-select-option value="section">小节</a-select-option>
                  <a-select-option value="point">知识点</a-select-option>
                </a-select>
              </a-form-item>
            </a-col>
            <a-col :span="12">
              <a-form-item label="终点名称">
                <a-input v-model:value="endPointName" placeholder="请输入终点名称" />
              </a-form-item>
            </a-col>
          </a-row>
          <div class="form-actions">
            <a-button type="primary" :loading="loading" @click="handleSearch">
              生成学习路径
            </a-button>
          </div>
        </a-form>
      </div>

      <!-- 学习路径展示区域 -->
      <div v-if="pathData" class="path-container">
        <div class="path-visualization">
          <div v-for="(node, index) in pathData.nodes" :key="index" class="node-container">
            <div class="node" :class="{ 'start-node': index === 0, 'end-node': index === pathData.nodes.length - 1 }">
              {{ node }}
            </div>
            <div v-if="index < pathData.nodes.length - 1" class="relationship">
              <div class="arrow"></div>
              <span class="relationship-text">{{ pathData.relationships[index] }}</span>
            </div>
          </div>
        </div>
        <div class="path-description">
          <a-typography-title :level="5">学习路径说明：</a-typography-title>
          <a-typography-paragraph>
            <ul>
              <li v-for="(desc, index) in pathData.description" :key="index">
                {{ desc }}
              </li>
            </ul>
          </a-typography-paragraph>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { message } from 'ant-design-vue'
import { getLearningPath } from '@/api/analysis'

// 表单数据
const startPointType = ref('')
const startPointName = ref('')
const endPointType = ref('')
const endPointName = ref('')
const loading = ref(false)
const pathData = ref(null)

// 获取学习路径
const handleSearch = async () => {
  if (!startPointType.value || !startPointName.value || !endPointType.value || !endPointName.value) {
    message.warning('请填写完整的起点和终点信息')
    return
  }

  loading.value = true
  try {
    const res = await getLearningPath(
      startPointType.value,
      startPointName.value,
      endPointType.value,
      endPointName.value
    )
    if (res.code === 200) {
      pathData.value = res.data
    } else {
      message.error(res.message || '获取学习路径失败')
    }
  } catch (error) {
    console.error('获取学习路径失败:', error)
    message.error('获取学习路径失败，请稍后重试')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.route-container {
  padding: 24px;
  background:rgb(231, 235, 242);
  min-height: calc(100vh - 180px);
  animation: fadeIn 0.5s ease-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.route-card {
  background: rgba(245, 240, 230, 0.95);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
}

.route-header {
  text-align: center;
  margin-bottom: 32px;
  padding: 16px;
  background: linear-gradient(120deg, rgba(33, 150, 243, 0.05), rgba(76, 175, 80, 0.05));
  border-radius: 12px;
}

.route-header h2 {
  font-size: 24px;
  margin-bottom: 8px;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  font-weight: 600;
}

.route-description {
  color: #666;
  font-size: 14px;
  line-height: 1.6;
}

.route-form {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px;
  background: linear-gradient(120deg, rgba(225, 230, 234, 0.05), rgba(162, 240, 164, 0.05));
  border-radius: 16px;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.3);
  transition: all 0.3s ease;
}

.route-form:hover {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08);
}

:deep(.ant-form-item) {
  margin-bottom: 24px;
  transition: all 0.3s ease;
}

:deep(.ant-form-item:hover) {
  transform: translateX(4px);
}

.form-actions {
  text-align: center;
  margin-top: 24px;
}

:deep(.ant-btn-primary) {
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  transition: all 0.3s ease;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
}

:deep(.ant-btn-primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
}

:deep(.ant-form-item-label > label) {
  font-weight: 600;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  background-clip: text;
  -webkit-text-fill-color: transparent;
  font-size: 15px;
  margin-bottom: 8px;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
}

:deep(.ant-select-selector),
:deep(.ant-input) {
  border-radius: 12px;
  border: 2px solid #e5e7eb;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-size: 15px;
  padding: 8px 16px;
  height: auto;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(8px);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
}

:deep(.ant-select-selector:hover),
:deep(.ant-input:hover) {
  border-color: #2196F3;
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.15);
  transform: translateY(-1px);
  background: rgba(255, 255, 255, 0.95);
}

:deep(.ant-select-focused .ant-select-selector),
:deep(.ant-input:focus) {
  border-color: #2196F3;
  box-shadow: 0 0 0 3px rgba(33, 150, 243, 0.15);
  outline: none;
  background: rgba(255, 255, 255, 1);
  transform: translateY(-2px);
}

:deep(.ant-select-selection-placeholder),
:deep(.ant-input::placeholder) {
  color: #94a3b8;
  font-size: 14px;
  font-style: italic;
  opacity: 0.8;
  transition: opacity 0.3s ease;
}

:deep(.ant-select:hover .ant-select-selection-placeholder),
:deep(.ant-input:hover::placeholder) {
  opacity: 0.6;
}

:deep(.ant-select-dropdown) {
  background: rgba(255, 255, 255, 0.98);
  backdrop-filter: blur(8px);
  border-radius: 12px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.3);
}

:deep(.ant-select-item) {
  padding: 10px 16px;
  border-radius: 8px;
  margin: 4px 8px;
  transition: all 0.3s ease;
}

:deep(.ant-select-item-option-active) {
  background: rgba(33, 150, 243, 0.08);
}

:deep(.ant-select-item-option-selected) {
  background: rgba(33, 150, 243, 0.12);
  font-weight: 600;
  color: #2196F3;
}

:deep(.ant-select-selection-item) {
  font-weight: 500;
  color: #1f2937;
}

/* 路径可视化样式 */
.path-container {
  margin-top: 40px;
  padding: 24px;
  background: rgba(250, 250, 250, 0.95);
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.2);
  animation: fadeIn 0.5s ease-out;
}

.path-visualization {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-wrap: wrap;
  gap: 20px;
  margin-bottom: 32px;
  padding: 24px;
  overflow-x: auto;
  background: linear-gradient(120deg, rgba(33, 150, 243, 0.02), rgba(76, 175, 80, 0.02));
  border-radius: 12px;
}

.node-container {
  display: flex;
  align-items: center;
  animation: slideIn 0.3s ease-out;
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(-10px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

.node {
  padding: 14px 28px;
  background: linear-gradient(120deg, #e6f7ff, #f0f7ff);
  border: 2px solid #1890ff;
  border-radius: 10px;
  color: #1890ff;
  font-weight: 500;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(24, 144, 255, 0.1);
}

.node:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(24, 144, 255, 0.15);
}

.start-node {
  background: linear-gradient(120deg, #f6ffed, #f0fce8);
  border-color: #52c41a;
  color: #52c41a;
  box-shadow: 0 2px 8px rgba(82, 196, 26, 0.1);
}

.start-node:hover {
  box-shadow: 0 4px 12px rgba(82, 196, 26, 0.15);
}

.end-node {
  background: linear-gradient(120deg, #fff7e6, #fff4e0);
  border-color: #faad14;
  color: #faad14;
  box-shadow: 0 2px 8px rgba(250, 173, 20, 0.1);
}

.end-node:hover {
  box-shadow: 0 4px 12px rgba(250, 173, 20, 0.15);
}

.relationship {
  display: flex;
  align-items: center;
  margin: 0 8px;
  position: relative;
}

.arrow {
  width: 64px;
  height: 2px;
  background: linear-gradient(90deg, #1890ff, #4CAF50);
  position: relative;
  transition: all 0.3s ease;
}

.arrow::after {
  content: '';
  position: absolute;
  right: -1px;
  top: 50%;
  transform: translateY(-50%);
  border-left: 8px solid #4CAF50;
  border-top: 6px solid transparent;
  border-bottom: 6px solid transparent;
  transition: all 0.3s ease;
}

.relationship-text {
  font-size: 13px;
  color: #666;
  margin-top: -32px;
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  background: transparent;
  padding: 4px 12px;
  border-radius: 6px;
  box-shadow: none;
  z-index: 1;
  white-space: nowrap;
}

.path-description {
  margin-top: 32px;
  padding: 20px;
  background: rgba(255, 255, 255, 0.95);
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.2);
}

:deep(.ant-typography-title) {
  margin-bottom: 16px !important;
  color: #2196F3;
  font-weight: 600;
}

:deep(.ant-typography-paragraph) {
  color: #666;
}

:deep(.ant-typography-paragraph ul) {
  padding-left: 20px;
}

:deep(.ant-typography-paragraph li) {
  margin-bottom: 8px;
  line-height: 1.6;
}

/* 响应式调整 */
@media (max-width: 768px) {
  .route-container {
    padding: 12px;
  }

  .route-form {
    padding: 16px;
  }

  .path-visualization {
    padding: 12px;
    gap: 8px;
  }

  .node {
    padding: 8px 16px;
    font-size: 14px;
  }

  .relationship {
    margin: 0 8px;
  }

  .arrow {
    width: 24px;
  }
}
</style>