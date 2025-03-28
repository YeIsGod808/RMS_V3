<template>
  <div class="node-management">
    <a-row :gutter="16">
      <a-col :span="24">
        <a-space style="margin-bottom: 16px">
          <a-select 
            v-model:value="nodeType" 
            style="width: 120px" 
            placeholder="选择节点类型"
            @change="handleNodeTypeChange"
            allowClear
          >
            <a-select-option v-for="type in knowledgeAdminStore.nodeTypes" :key="type.value" :value="type.value">
              {{ type.label }}
            </a-select-option>
          </a-select>
          
          <template v-if="nodeType === 'section'">
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
          </template>
          
          <template v-if="nodeType === 'point'">
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
          </template>
          
          <a-button type="primary" @click="showAddNodeModal">
            <plus-outlined /> 添加节点
          </a-button>
        </a-space>
      </a-col>
    </a-row>
    
    <!-- 节点列表 -->
    <a-row :gutter="16">
      <a-col :span="24">
        <a-card title="知识节点列表" :bordered="false">
          <a-table
            :columns="nodeColumns"
            :data-source="getNodeDataSource()"
            :loading="knowledgeAdminStore.loading"
            row-key="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'type'">
                <a-tag :color="getNodeTypeColor(record.type)">
                  {{ getNodeTypeName(record.type) }}
                </a-tag>
              </template>
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" @click="handleEditNode(record)">
                    编辑
                  </a-button>
                  <a-popconfirm
                    title="确定要删除该节点吗？删除将会影响相关联的节点与关系。"
                    @confirm="handleDeleteNode(record)"
                  >
                    <a-button type="link" danger>删除</a-button>
                  </a-popconfirm>
                </a-space>
              </template>
            </template>
          </a-table>
        </a-card>
      </a-col>
    </a-row>

    <!-- 节点表单模态框 -->
    <a-modal
      v-model:visible="nodeModalVisible"
      :title="nodeModalTitle"
      @ok="handleNodeModalOk"
      :confirm-loading="modalLoading"
    >
      <a-form
        ref="nodeFormRef"
        :model="nodeForm"
        :rules="nodeRules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
      >
        <a-form-item label="名称" name="name">
          <a-input v-model:value="nodeForm.name" />
        </a-form-item>
        <a-form-item label="类型" name="type">
          <a-select v-model:value="nodeForm.type" :disabled="nodeForm.id || !!nodeType">
            <a-select-option v-for="type in knowledgeAdminStore.nodeTypes" :key="type.value" :value="type.value">
              {{ type.label }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="描述" name="description">
          <a-textarea v-model:value="nodeForm.description" :rows="4" />
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useKnowledgeAdminStore } from '@/store/modules/knowledgeAdmin'

const knowledgeAdminStore = useKnowledgeAdminStore()

// 状态管理
const nodeType = ref('')
const selectedChapter = ref(null)
const selectedSection = ref(null)
const modalLoading = ref(false)

// 节点表单模态框状态
const nodeModalVisible = ref(false)
const nodeModalTitle = ref('添加节点')
const nodeFormRef = ref()
const nodeForm = reactive({
  id: null,
  name: '',
  type: '',
  description: ''
})

// 节点表单验证规则
const nodeRules = {
  name: [{ required: true, message: '请输入节点名称' }],
  type: [{ required: true, message: '请选择节点类型' }]
}

// 节点表格列定义
const nodeColumns = [
  { title: '名称', dataIndex: 'name' },
  { title: '类型', dataIndex: 'type', key: 'type' },
  { title: '描述', dataIndex: 'description', ellipsis: true },
  { title: '操作', key: 'action', width: 150 }
]

// 获取节点类型颜色
const getNodeTypeColor = (type) => {
  switch (type) {
    case 'chapter': return 'blue'
    case 'section': return 'green'
    case 'point': return 'orange'
    default: return 'default'
  }
}

// 获取节点类型名称
const getNodeTypeName = (type) => {
  switch (type) {
    case 'chapter': return '章节'
    case 'section': return '小节'
    case 'point': return '知识点'
    default: return type
  }
}

// 获取节点数据源
const getNodeDataSource = () => {
  if (nodeType.value === 'chapter') {
    return knowledgeAdminStore.chapters
  } else if (nodeType.value === 'section') {
    if (selectedChapter.value) {
      return knowledgeAdminStore.sections
    } else {
      return knowledgeAdminStore.allSections
    }
  } else if (nodeType.value === 'point') {
    if (selectedSection.value) {
      return knowledgeAdminStore.points
    } else if (selectedChapter.value) {
      return knowledgeAdminStore.points
    } else {
      return knowledgeAdminStore.allPoints
    }
  }
  return []
}

// 知识图谱节点类型变化
const handleNodeTypeChange = async (value) => {
  nodeType.value = value
  selectedChapter.value = null
  selectedSection.value = null
  
  if (!value) {
    // 如果是清空操作，加载所有章节
    await knowledgeAdminStore.loadChapters()
    return
  }
  
  if (value === 'chapter') {
    await knowledgeAdminStore.loadChapters()
  } else if (value === 'section') {
    // 当选择小节节点类型时，加载所有小节
    await knowledgeAdminStore.loadChapters()
    await knowledgeAdminStore.loadAllSections()
  } else if (value === 'point') {
    // 当选择知识点节点类型时，加载所有知识点
    await knowledgeAdminStore.loadChapters()
    await knowledgeAdminStore.loadAllSections()
    await knowledgeAdminStore.loadAllPoints()
  }
}

// 选择章节变化
const handleChapterChange = async (value) => {
  selectedChapter.value = value
  selectedSection.value = null
  
  if (!value) {
    // 如果是清空操作
    if (nodeType.value === 'section') {
      await knowledgeAdminStore.loadAllSections()
    } else if (nodeType.value === 'point') {
      await knowledgeAdminStore.loadAllPoints()
    }
    return
  }
  
  if (nodeType.value === 'section' || nodeType.value === 'point') {
    // 获取章节ID
    const chapter = knowledgeAdminStore.chapters.find(c => c.name === value)
    if (chapter) {
      await knowledgeAdminStore.loadSectionsByChapterId(chapter.id)
    }
  }
}

// 选择小节变化
const handleSectionChange = async (value) => {
  selectedSection.value = value
  
  if (!value) {
    // 如果是清空操作，但保留章节筛选
    if (nodeType.value === 'point' && selectedChapter.value) {
      const chapter = knowledgeAdminStore.chapters.find(c => c.name === selectedChapter.value)
      if (chapter) {
        await knowledgeAdminStore.loadSectionsByChapterId(chapter.id)
      }
    }
    return
  }
  
  if (nodeType.value === 'point') {
    // 获取小节ID
    const section = knowledgeAdminStore.sections.find(s => s.name === value)
    if (section) {
      await knowledgeAdminStore.loadPointsBySectionId(section.id)
    }
  }
}

// 显示添加节点模态框
const showAddNodeModal = () => {
  nodeForm.id = null
  nodeForm.name = ''
  nodeForm.description = ''
  nodeForm.type = nodeType.value || 'chapter'
  nodeModalTitle.value = '添加节点'
  nodeModalVisible.value = true
}

// 处理编辑节点
const handleEditNode = (record) => {
  nodeForm.id = record.id
  nodeForm.name = record.name
  nodeForm.type = record.type
  nodeForm.description = record.description || ''
  nodeModalTitle.value = '编辑节点'
  nodeModalVisible.value = true
}

// 处理删除节点
const handleDeleteNode = async (record) => {
  try {
    await knowledgeAdminStore.deleteNode({
      name: record.name,
      nodeType: record.type
    })
  } catch (error) {
    console.error(error)
  }
}

// 处理节点模态框确认
const handleNodeModalOk = async () => {
  try {
    await nodeFormRef.value.validate()
    modalLoading.value = true
    
    let success = false
    if (nodeForm.id) {
      // 编辑节点
      success = await knowledgeAdminStore.updateNode({
        name: nodeForm.name,
        nodeType: nodeForm.type,
        propertyName: 'description',
        newValue: nodeForm.description
      })
    } else {
      // 添加节点
      success = await knowledgeAdminStore.addNode({
        name: nodeForm.name,
        type: nodeForm.type,
        description: nodeForm.description
      })
    }
    
    if (success) {
      // 只有成功才关闭模态框
      nodeModalVisible.value = false
    }
  } catch (error) {
    console.error(error)
  } finally {
    modalLoading.value = false
  }
}

// 加载数据
onMounted(async () => {
  await knowledgeAdminStore.loadChapters()
})
</script>

<style scoped>
/* 整体容器样式 */
.node-management {
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

/* 卡片样式 */
:deep(.ant-card) {
  border-radius: 12px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(12px);
  transition: all 0.3s ease;
  border: 1px solid rgba(255, 255, 255, 0.2);
  overflow: hidden;
  margin-bottom: 16px;
}

:deep(.ant-card:hover) {
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.12);
  transform: translateY(-2px);
}

:deep(.ant-card-head) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
  padding: 16px 24px;
}

:deep(.ant-card-head-title) {
  font-size: 18px;
  font-weight: 600;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  -webkit-background-clip: text;
  /* 添加标准属性 */
  background-clip: text;
  -webkit-text-fill-color: transparent;
}

/* 表格样式 */
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
  background-color: rgba(33, 150, 243, 0.05);
}

/* 按钮样式 */
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

:deep(.ant-btn-link) {
  transition: all 0.3s ease;
}

:deep(.ant-btn-link:hover) {
  transform: translateY(-1px);
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

/* 选择器样式 */
:deep(.ant-select) {
  border-radius: 8px;
}

:deep(.ant-select-selector) {
  border-radius: 8px !important;
  border: 1px solid rgba(33, 150, 243, 0.2) !important;
  transition: all 0.3s ease !important;
}

:deep(.ant-select:not(.ant-select-disabled):hover .ant-select-selector) {
  border-color: #2196F3 !important;
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.1) !important;
}

:deep(.ant-select-focused:not(.ant-select-disabled) .ant-select-selector) {
  border-color: #2196F3 !important;
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2) !important;
}

/* 标签样式 */
:deep(.ant-tag) {
  border-radius: 4px;
  padding: 2px 8px;
  font-weight: 500;
  border: none;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

/* 模态框样式 */
:deep(.ant-modal-content) {
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1);
}

:deep(.ant-modal-header) {
  padding: 16px 24px;
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

:deep(.ant-modal-footer) {
  border-top: 1px solid rgba(0, 0, 0, 0.06);
  padding: 16px 24px;
}

/* 表单样式 */
:deep(.ant-form-item-label > label) {
  color: #333;
  font-weight: 500;
}

:deep(.ant-form-item-explain-error) {
  margin-top: 4px;
  font-size: 12px;
}

/* 空状态提示 */
:deep(.ant-empty) {
  margin: 32px 0;
}

:deep(.ant-empty-description) {
  color: #999;
}

/* 间距调整 */
.ant-space {
  margin-bottom: 16px;
}
</style>