<template>
  <div class="relation-management">
    <!-- 关系类型选择 -->
    <a-row :gutter="16">
      <a-col :span="24">
        <a-space style="margin-bottom: 16px">
          <a-select 
            v-model:value="relationSourceType" 
            style="width: 150px" 
            placeholder="选择源节点类型"
            @change="handleRelationSourceTypeChange"
            allowClear
          >
            <a-select-option value="chapter">章节</a-select-option>
            <a-select-option value="section">小节</a-select-option>
            <a-select-option value="point">知识点</a-select-option>
          </a-select>
          
          <a-select 
            v-model:value="relationTargetType" 
            style="width: 150px" 
            placeholder="选择目标节点类型"
            :disabled="!relationSourceType"
            @change="handleRelationTargetTypeChange"
            allowClear
          >
            <a-select-option v-if="relationSourceType === 'chapter'" value="section">小节</a-select-option>
            <a-select-option v-if="relationSourceType === 'section'" value="point">知识点</a-select-option>
            <a-select-option v-if="relationSourceType === 'section'" value="section">小节</a-select-option>
            <a-select-option v-if="relationSourceType === 'point'" value="point">知识点</a-select-option>
            <a-select-option v-if="relationSourceType === 'chapter'" value="chapter">章节</a-select-option>
          </a-select>
          
          <a-button 
            type="primary" 
            @click="showAddRelationModal"
            :disabled="!relationSourceType || !relationTargetType"
          >
            <plus-outlined /> 添加关系
          </a-button>
        </a-space>
      </a-col>
    </a-row>
    
    <!-- 加载源节点和目标节点 -->
    <a-row :gutter="16" v-if="relationSourceType && relationTargetType">
      <a-col :span="12">
        <a-card title="源节点" :bordered="false">
          <a-table
            :columns="getSourceNodeColumns()"
            :data-source="getSourceNodeDataSource()"
            :loading="knowledgeAdminStore.loading"
            row-key="id"
            :row-selection="{ 
              type: 'radio', 
              selectedRowKeys: selectedSourceNodes,
              onChange: onSourceNodeSelectionChange
            }"
          />
        </a-card>
      </a-col>
      <a-col :span="12">
        <a-card title="目标节点" :bordered="false">
          <a-table
            :columns="getTargetNodeColumns()"
            :data-source="getTargetNodeDataSource()"
            :loading="knowledgeAdminStore.loading"
            row-key="id"
            :row-selection="{ 
              type: 'radio', 
              selectedRowKeys: selectedTargetNodes,
              onChange: onTargetNodeSelectionChange
            }"
          />
        </a-card>
      </a-col>
    </a-row>
    
    <!-- 现有关系列表 -->
    <a-row :gutter="16" style="margin-top: 16px">
      <a-col :span="24">
        <a-card title="现有关系" :bordered="false">
          <a-table
            :columns="relationColumns"
            :data-source="getRelationDataSource()"
            :loading="knowledgeAdminStore.loading"
            row-key="id"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'relationType'">
                <a-tag :color="getRelationTypeColor(record.relationType)">
                  {{ getRelationTypeName(record.relationType) }}
                </a-tag>
              </template>
              <template v-if="column.key === 'action'">
                <a-space>
                  <a-button type="link" @click="handleEditRelation(record)">
                    编辑
                  </a-button>
                  <a-popconfirm
                    :title="`确定要删除 '${record.sourceName}' 和 '${record.targetName}' 之间的 '${record.relationType}' 关系吗？`"
                    @confirm="handleDeleteRelation(record)"
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

    <!-- 关系表单模态框 -->
    <a-modal
      v-model:visible="relationModalVisible"
      :title="relationModalTitle"
      @ok="handleRelationModalOk"
      :confirm-loading="modalLoading"
    >
      <a-form
        ref="relationFormRef"
        :model="relationForm"
        :rules="relationRules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
      >
        <a-form-item label="源节点" name="sourceName">
          <a-select v-model:value="relationForm.sourceName" :disabled="!!selectedSourceNode">
            <a-select-option v-for="node in getSourceNodeDataSource()" :key="node.name" :value="node.name">
              {{ node.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="目标节点" name="targetName">
          <a-select v-model:value="relationForm.targetName" :disabled="!!selectedTargetNode">
            <a-select-option v-for="node in getTargetNodeDataSource()" :key="node.name" :value="node.name">
              {{ node.name }}
            </a-select-option>
          </a-select>
        </a-form-item>
        <a-form-item label="关系类型" name="relationType">
          <a-select 
            v-model:value="relationForm.relationType"
            :disabled="isHierarchicalRelation()"
          >
            <a-select-option 
              v-for="type in getAvailableRelationTypes()" 
              :key="type.value" 
              :value="type.value"
              :disabled="isRelationTypeDisabled(type.value)"
            >
              {{ type.label }}
            </a-select-option>
          </a-select>
          <div v-if="isHierarchicalRelation()" style="color: #999; font-size: 12px;">
            层级关系（章节-小节、小节-知识点）必须使用"包含"关系类型
          </div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined } from '@ant-design/icons-vue'
import { useKnowledgeAdminStore } from '@/store/modules/knowledgeAdmin'

const knowledgeAdminStore = useKnowledgeAdminStore()

// 状态管理
const relationSourceType = ref(null)
const relationTargetType = ref(null)
const selectedSourceNodes = ref([])
const selectedTargetNodes = ref([])
const selectedSourceNode = ref(null)
const selectedTargetNode = ref(null)
const modalLoading = ref(false)

// 关系表单模态框状态
const relationModalVisible = ref(false)
const relationModalTitle = ref('添加关系')
const relationFormRef = ref()
const relationForm = reactive({
  sourceName: '',
  targetName: '',
  sourceType: '',
  targetType: '',
  relationType: 'CONTAINS',
  id: null
})

// 关系表单验证规则
const relationRules = {
  sourceName: [{ required: true, message: '请选择源节点' }],
  targetName: [{ required: true, message: '请选择目标节点' }]
}

// 关系表格列定义
const relationColumns = [
  { title: '源节点', dataIndex: 'sourceName' },
  { title: '目标节点', dataIndex: 'targetName' },
  { title: '关系类型', dataIndex: 'relationType', key: 'relationType' },
  { title: '操作', key: 'action', width: 150 }
]

// 获取源节点列
const getSourceNodeColumns = () => {
  return [
    { title: '名称', dataIndex: 'name' },
    { title: '描述', dataIndex: 'description', ellipsis: true }
  ]
}

// 获取目标节点列
const getTargetNodeColumns = () => {
  return [
    { title: '名称', dataIndex: 'name' },
    { title: '描述', dataIndex: 'description', ellipsis: true }
  ]
}

// 获取源节点数据源
const getSourceNodeDataSource = () => {
  if (relationSourceType.value === 'chapter') {
    return knowledgeAdminStore.chapters
  } else if (relationSourceType.value === 'section') {
    return knowledgeAdminStore.allSections
  } else if (relationSourceType.value === 'point') {
    return knowledgeAdminStore.allPoints
  }
  return []
}

// 获取目标节点数据源
const getTargetNodeDataSource = () => {
  if (relationTargetType.value === 'chapter') {
    return knowledgeAdminStore.chapters
  } else if (relationTargetType.value === 'section') {
    return knowledgeAdminStore.allSections
  } else if (relationTargetType.value === 'point') {
    return knowledgeAdminStore.allPoints
  }
  return []
}

// 关系源节点类型变化
const handleRelationSourceTypeChange = async (value) => {
  relationSourceType.value = value
  relationTargetType.value = null
  selectedSourceNodes.value = []
  selectedTargetNodes.value = []
  selectedSourceNode.value = null
  selectedTargetNode.value = null
  
  if (!value) {
    // 如果是清空操作，不需要加载数据
    return
  }
  
  // 根据源节点类型加载源节点数据
  if (value === 'chapter') {
    await knowledgeAdminStore.loadChapters()
  } else if (value === 'section') {
    await knowledgeAdminStore.loadAllSections()
  } else if (value === 'point') {
    await knowledgeAdminStore.loadAllPoints()
  }
}

// 源节点选择变化
const onSourceNodeSelectionChange = async (selectedRowKeys, selectedRows) => {
  selectedSourceNodes.value = selectedRowKeys
  if (selectedRows && selectedRows.length > 0) {
    selectedSourceNode.value = selectedRows[0]
    relationForm.sourceName = selectedRows[0].name
    relationForm.sourceType = selectedRows[0].type
    
    // 如果已经选择了目标节点类型，加载相关关系
    if (relationTargetType.value) {
      await loadAndUpdateRelations()
    }
    
    // 如果选择了章节，加载对应的小节
    if (selectedRows[0].type === 'chapter') {
      try {
        await knowledgeAdminStore.loadSections(selectedRows[0].id)
      } catch (error) {
        console.error('加载小节失败:', error)
      }
    } 
    // 如果选择了小节，加载对应的知识点
    else if (selectedRows[0].type === 'section') {
      try {
        await knowledgeAdminStore.loadPoints(selectedRows[0].id)
      } catch (error) {
        console.error('加载知识点失败:', error)
      }
    }
  } else {
    selectedSourceNode.value = null
    relationForm.sourceName = ''
    relationForm.sourceType = ''
  }
}

// 目标节点选择变化
const onTargetNodeSelectionChange = async (selectedRowKeys, selectedRows) => {
  selectedTargetNodes.value = selectedRowKeys
  if (selectedRows && selectedRows.length > 0) {
    selectedTargetNode.value = selectedRows[0]
    relationForm.targetName = selectedRows[0].name
    relationForm.targetType = selectedRows[0].type
    
    // 如果已经选择了源节点，加载相关关系
    if (selectedSourceNode.value) {
      await loadAndUpdateRelations()
    }
  } else {
    selectedTargetNode.value = null
    relationForm.targetName = ''
    relationForm.targetType = ''
  }
}

// 关系目标节点类型变化
const handleRelationTargetTypeChange = async (value) => {
  relationTargetType.value = value
  selectedTargetNodes.value = []
  selectedTargetNode.value = null
  
  if (!value) {
    // 如果是清空操作，不需要加载
    return
  }
  
  console.log(`加载${relationSourceType.value}到${value}的关系`)
  
  // 根据目标节点类型加载目标节点数据
  if (value === 'section') {
    await knowledgeAdminStore.loadAllSections()
  } else if (value === 'point') {
    await knowledgeAdminStore.loadAllPoints()
  } else if (value === 'chapter') {
    await knowledgeAdminStore.loadChapters()
  }
  
  // 加载并更新关系数据
  await loadAndUpdateRelations()
}

// 监听关系源类型和目标类型的变化
watch([relationSourceType, relationTargetType], async ([newSourceType, newTargetType]) => {
  if (newSourceType && newTargetType) {
    // 当源类型和目标类型都有值时，加载关系数据
    await loadAndUpdateRelations()
  }
})

// 判断关系是否为层级关系
const isHierarchicalRelation = () => {
  // 章节→小节，小节→知识点是层级关系，必须使用"包含"关系类型
  return (relationSourceType.value === 'chapter' && relationTargetType.value === 'section') ||
         (relationSourceType.value === 'section' && relationTargetType.value === 'point')
}

// 显示添加关系模态框
const showAddRelationModal = () => {
  relationForm.id = null
  relationForm.sourceName = selectedSourceNode.value ? selectedSourceNode.value.name : ''
  relationForm.targetName = selectedTargetNode.value ? selectedTargetNode.value.name : ''
  relationForm.sourceType = relationSourceType.value
  relationForm.targetType = relationTargetType.value
  
  // 根据节点类型预设关系类型
  if (isHierarchicalRelation()) {
    relationForm.relationType = '包含' // 层级关系固定为"包含"
  } else {
    relationForm.relationType = '相关' // 非层级关系默认为"相关"
  }
  
  relationModalTitle.value = '添加关系'
  relationModalVisible.value = true
}

// 处理关系模态框确认
const handleRelationModalOk = async () => {
  try {
    // 首先验证源节点和目标节点不能相同
    if (relationForm.sourceName === relationForm.targetName) {
      message.error('源节点和目标节点不能是同一个节点')
      return
    }
    
    await relationFormRef.value.validate()
    modalLoading.value = true
    console.log('提交关系数据:', relationForm)
    
    // 强制层级关系使用"包含"类型
    const relationData = {...relationForm}
    if (isHierarchicalRelation()) {
      relationData.relationType = '包含'
    }
    
    // 添加或更新关系
    let success = false
    if (relationForm.id) {
      // 编辑关系
      console.log('更新关系:', relationData)
      success = await knowledgeAdminStore.updateRelation({
        sourceName: relationData.sourceName,
        targetName: relationData.targetName,
        sourceType: relationData.sourceType,
        targetType: relationData.targetType,
        newRelationType: relationData.relationType
      })
    } else {
      // 添加关系
      console.log('添加关系:', relationData)
      success = await knowledgeAdminStore.addRelation({
        sourceName: relationData.sourceName,
        targetName: relationData.targetName,
        sourceType: relationData.sourceType,
        targetType: relationData.targetType,
        relationType: relationData.relationType
      })
    }
    
    if (success) {
      // 只有成功才关闭模态框
      relationModalVisible.value = false
      
      // 重新加载关系数据
      console.log('重新加载关系数据')
      await loadAndUpdateRelations()
    } else {
      console.log('关系操作失败，不关闭模态框')
    }
  } catch (error) {
    console.error('关系表单提交错误:', error)
  } finally {
    modalLoading.value = false
  }
}

// 加载并更新关系数据的函数
const loadAndUpdateRelations = async () => {
  if (!relationSourceType.value || !relationTargetType.value) return
  
  try {
    console.log(`重新加载 ${relationSourceType.value} -> ${relationTargetType.value} 之间的关系`)
    
    // 记录刷新前的时间戳，用于判断数据是否是最新的
    const refreshTimestamp = Date.now()
    
    // 先清空关系列表，避免显示过期数据
    knowledgeAdminStore.relations = []
    
    // 添加重试机制，确保关系数据正确加载
    let retryCount = 0
    const maxRetries = 2
    let success = false
    
    while (!success && retryCount <= maxRetries) {
      try {
        await knowledgeAdminStore.loadRelationsByTypes(
          relationSourceType.value,
          relationTargetType.value
        )
        
        // 确认已加载关系数量
        console.log(`已加载 ${knowledgeAdminStore.relations.length} 个关系`)
        success = true
      } catch (error) {
        retryCount++
        console.error(`加载关系数据失败 (尝试 ${retryCount}/${maxRetries + 1}):`, error)
        
        if (retryCount <= maxRetries) {
          console.log(`${retryCount}秒后重试...`)
          await new Promise(resolve => setTimeout(resolve, retryCount * 1000))
        } else {
          throw error // 重试次数用完，抛出错误
        }
      }
    }
    
    // 添加延迟验证机制，确认数据已正确加载
    if (knowledgeAdminStore.relations.length === 0) {
      console.log('关系列表为空，等待1秒后再次尝试加载...')
      await new Promise(resolve => setTimeout(resolve, 1000))
      
      await knowledgeAdminStore.loadRelationsByTypes(
        relationSourceType.value,
        relationTargetType.value
      )
      
      console.log(`延迟验证：已加载 ${knowledgeAdminStore.relations.length} 个关系`)
    }
  } catch (error) {
    console.error('加载关系数据失败:', error)
    message.error('加载关系数据失败: ' + (error.message || '未知错误'))
  }
}

// 获取关系数据源
const getRelationDataSource = () => {
  if (!relationSourceType.value || !relationTargetType.value) {
    return []
  }
  
  // 当没有加载关系数据时，显示一个空列表
  if (knowledgeAdminStore.relations.length === 0) {
    return []
  }
  
  // 过滤关系数据，确保只显示当前选择的源类型和目标类型的关系
  const filteredRelations = knowledgeAdminStore.relations.filter(relation =>
    relation.sourceType === relationSourceType.value && 
    relation.targetType === relationTargetType.value
  )
  
  // 调试日志
  console.log(`过滤后关系数: ${filteredRelations.length}`, filteredRelations)
  
  return filteredRelations
}

// 处理编辑关系
const handleEditRelation = (record) => {
  relationForm.id = record.id
  relationForm.sourceName = record.sourceName
  relationForm.targetName = record.targetName
  relationForm.sourceType = record.sourceType
  relationForm.targetType = record.targetType
  relationForm.relationType = record.relationType
  relationModalTitle.value = '编辑关系'
  relationModalVisible.value = true
}

// 处理删除关系
const handleDeleteRelation = async (record) => {
  try {
    console.log('删除关系:', record)
    
    if (!record.relationType) {
      message.error('缺少关系类型信息，无法删除')
      return
    }
    
    message.loading({ content: '正在删除关系...', key: 'deleteRelation' })
    
    // 确保传递所有必需的参数，包括关系类型
    const success = await knowledgeAdminStore.deleteRelation({
      sourceName: record.sourceName,
      targetName: record.targetName,
      sourceType: record.sourceType,
      targetType: record.targetType,
      relationType: record.relationType
    })
    
    if (success) {
      message.success({ content: '删除关系成功', key: 'deleteRelation' })
      // 删除成功后刷新关系列表
      await loadAndUpdateRelations()
    } else {
      // 即使返回失败，也尝试刷新一次关系列表，因为可能后端已成功删除
      console.log('尝试强制刷新关系列表...')
      await loadAndUpdateRelations()
      
      // 如果刷新后关系列表中找不到该关系，说明后端实际已删除成功
      const relationStillExists = knowledgeAdminStore.relations.some(
        r => r.sourceName === record.sourceName && 
             r.targetName === record.targetName && 
             r.relationType === record.relationType
      )
      
      if (!relationStillExists) {
        message.success({ content: '关系已成功删除', key: 'deleteRelation' })
      }
    }
  } catch (error) {
    console.error('删除关系错误:', error)
    message.error({ content: `删除关系失败: ${error.message || '未知错误'}`, key: 'deleteRelation' })
    
    // 尝试刷新关系列表
    try {
      await loadAndUpdateRelations()
    } catch (refreshError) {
      console.error('刷新关系列表失败:', refreshError)
    }
  }
}

// 获取可用关系类型
const getAvailableRelationTypes = () => {
  return [
    { value: '包含', label: '包含' },
    { value: '前置', label: '前置' },
    { value: '相关', label: '相关' },
    { value: '扩展', label: '扩展' }
  ]
}

// 判断关系类型是否可用
const isRelationTypeDisabled = (value) => {
  // 如果是层级关系，则只允许"包含"类型
  if (isHierarchicalRelation()) {
    return value !== '包含'
  }
  
  return false // 其他情况都允许
}

// 获取关系类型颜色
const getRelationTypeColor = (type) => {
  const relationColors = {
    '包含': 'blue',
    '前置': 'orange',
    '相关': 'green',
    '扩展': 'purple'
  }
  return relationColors[type] || 'default'
}

// 获取关系类型名称
const getRelationTypeName = (type) => {
  return type // 直接返回中文类型名称
}

// 重置关系选择状态
const resetRelationSelections = () => {
  relationSourceType.value = null
  relationTargetType.value = null
  selectedSourceNodes.value = []
  selectedTargetNodes.value = []
  selectedSourceNode.value = null
  selectedTargetNode.value = null
}

// 加载数据
onMounted(async () => {
  await knowledgeAdminStore.loadChapters()
  await knowledgeAdminStore.loadAllSections()
  await knowledgeAdminStore.loadAllPoints()
  
  // 重置选择状态
  resetRelationSelections()
})
</script>

<style scoped>
/* 整体容器样式 */
.relation-management {
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