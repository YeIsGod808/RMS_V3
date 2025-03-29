<template>
  <div class="graph-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="title">知识图谱</h1>
        <div class="subtitle">探索数据库知识的连接与关系</div>
      </div>
      <div class="header-actions">
        <a-space>
          <a-select
            v-model:value="viewMode"
            style="width: 120px"
            @change="handleViewModeChange"
          >
            <a-select-option value="force">力导向图</a-select-option>
          </a-select>
          <a-button type="primary" @click="handleRefresh">
            <template #icon><ReloadOutlined /></template>
            刷新
          </a-button>
        </a-space>
      </div>
    </div>

    <div class="graph-container">

      <div class="graph-card">
        <a-spin :spinning="loading">
          <knowledge-graph 
            :mode="viewMode"
            @node-click="handleNodeClick"
          />
        </a-spin>
      </div>
      
      <div class="info-panel" v-if="selectedNode">
        <div class="panel-header">
          <h3>{{ selectedNode.name }}</h3>
          <a-button type="text" @click="selectedNode = null">
            <CloseOutlined />
          </a-button>
        </div>
        <div class="panel-content">
          <p class="description">{{ selectedNode.description }}</p>
          <div class="stats">
            <div class="stat-item">
              <span class="label">关联节点</span>
              <span class="value">{{ selectedNode.connections }}</span>
            </div>
            <div class="stat-item">
              <span class="label">重要程度</span>
              <span class="value">
                <a-rate :value="selectedNode.importance" disabled />
              </span>
            </div>
          </div>
          <a-divider />
          <div class="actions">
            <a-button type="primary" @click="handleExplore">
              深入探索
            </a-button>
            <a-button @click="handleResources">
              相关资源
            </a-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useGraphStore } from '@/store/modules/graph'
import KnowledgeGraph from '@/components/graph/KnowledgeGraph.vue'
import { ReloadOutlined, CloseOutlined } from '@ant-design/icons-vue'

const graphStore = useGraphStore()
const loading = ref(false)
const viewMode = ref('force')
const selectedNode = ref(null)

onMounted(() => {
  loadGraphData()
})

const loadGraphData = async () => {
  loading.value = true
  try {
    await graphStore.loadChapterGraph()
  } finally {
    loading.value = false
  }
}

const handleViewModeChange = (mode) => {
  viewMode.value = mode
}

const handleRefresh = () => {
  loadGraphData()
}

const handleNodeClick = (node) => {
  selectedNode.value = node
}

const handleExplore = () => {
  // 实现深入探索逻辑
}

const handleResources = () => {
  // 实现查看相关资源逻辑
}

</script>

<style scoped>
.graph-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #f5f0e6 0%, #eae3d7 100%);
  padding: 24px;
  transition: background 0.3s ease;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding: 24px;
  background: rgba(245, 240, 230, 0.95);
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
}

.header-content .title {
  font-size: 32px;
  font-weight: 700;
  margin: 0;
  background: linear-gradient(120deg, #2196F3, #00BCD4);
  -webkit-background-clip: text;
  /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
  letter-spacing: 0.5px;
}
.header-content .subtitle {
  color: #5c6b7c;
  margin-top: 8px;
  font-size: 16px;
  letter-spacing: 0.3px;
}

.graph-container {
  display: flex;
  gap: 24px;
  margin-top: 24px;
}

.graph-card {
  flex: 1;
  background: rgba(245, 240, 230, 0.95);
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  backdrop-filter: blur(12px);
  min-height: 600px;
  transition: all 0.3s ease;
  border: 1px solid rgba(245, 240, 230, 0.3);
}

.graph-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
}

.info-panel {
  width: 320px;
  background: rgba(245, 240, 230, 0.95);
  border-radius: 16px;
  backdrop-filter: blur(12px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
  border: 1px solid rgba(245, 240, 230, 0.3);
  animation: slideIn 0.3s ease-out;
}

.panel-header {
  padding: 20px;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.panel-header h3 {
  margin: 0;
  font-size: 20px;
  color: #2196F3;
  font-weight: 600;
}

.panel-content {
  padding: 20px;
}

.description {
  color: #5c6b7c;
  line-height: 1.6;
  font-size: 15px;
  margin-bottom: 20px;
}

.stats {
  margin: 20px 0;
  background: rgba(245, 240, 230, 0.8);
  padding: 16px;
  border-radius: 12px;
}

.stat-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding: 8px 0;
}

.stat-item:last-child {
  margin-bottom: 0;
}

.stat-item .label {
  color: #5c6b7c;
  font-size: 14px;
}

.stat-item .value {
  font-weight: 600;
  color: #2196F3;
  font-size: 15px;
}

.actions {
  display: flex;
  gap: 12px;
  margin-top: 24px;
}

.actions .ant-btn {
  flex: 1;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

.actions .ant-btn-primary {
  background: linear-gradient(120deg, #2196F3, #00BCD4);
  border: none;
}

.actions .ant-btn-primary:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.2);
}

@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateX(20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}
</style>
