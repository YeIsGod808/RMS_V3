<template>
  <div class="user-group-management-container">
    <a-card title="用户组管理" class="group-card">
      <!-- 操作按钮 -->
      <template #extra>
        <a-button type="primary" @click="showCreateGroupModal" class="create-button">
          <plus-outlined /> 创建用户组
        </a-button>
      </template>
      
      <!-- 用户组列表 -->
      <a-table
        :columns="groupColumns"
        :data-source="groups"
        :loading="loading"
        rowKey="group_id"
        :pagination="{ pageSize: 10 }"
        class="group-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'actions'">
            <a-space>
              <a-button type="link" @click="showGroupDetailModal(record)" class="action-button view-button">
                <template #icon><team-outlined /></template>
                查看成员
              </a-button>
              <a-button type="link" @click="showEditGroupModal(record)" class="action-button edit-button">
                <template #icon><edit-outlined /></template>
                编辑
              </a-button>
              <a-popconfirm
                title="确定要删除此用户组吗？"
                @confirm="deleteUserGroup(record)"
                ok-text="确定"
                cancel-text="取消"
                class="delete-confirm"
              >
                <a-button type="link" danger class="action-button delete-button">
                  <template #icon><delete-outlined /></template>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>
    
    <!-- 创建用户组模态框 -->
    <a-modal
      v-model:visible="createGroupModalVisible"
      title="创建用户组"
      @ok="handleCreateGroup"
      :confirm-loading="confirmLoading"
      class="group-modal"
    >
      <a-form :model="groupForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }" class="group-form">
        <a-form-item label="用户组名称" name="name" :rules="[{ required: true, message: '请输入用户组名称' }]">
          <a-input v-model:value="groupForm.name" placeholder="请输入用户组名称" />
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 编辑用户组模态框 -->
    <a-modal
      v-model:visible="editGroupModalVisible"
      title="编辑用户组"
      @ok="handleEditGroup"
      :confirm-loading="confirmLoading"
      class="group-modal"
    >
      <a-form :model="groupForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }" class="group-form">
        <a-form-item label="用户组名称" name="name" :rules="[{ required: true, message: '请输入用户组名称' }]">
          <a-input v-model:value="groupForm.name" placeholder="请输入用户组名称" />
        </a-form-item>
      </a-form>
    </a-modal>
    
    <!-- 用户组详情模态框 -->
    <a-modal
      v-model:visible="groupDetailModalVisible"
      title="用户组详情"
      width="800px"
      :footer="null"
      class="detail-modal"
    >
      <div v-if="selectedGroup" class="group-detail-container">
        <div class="detail-header">
          <h3>{{ selectedGroup.name }}</h3>
          <p>创建者: {{ selectedGroup.owner?.nickname || selectedGroup.owner?.user_id || '未知' }}</p>
        </div>
        
        <a-divider class="custom-divider" />
        
        <div class="group-members">
          <div class="members-header">
            <h4>用户组成员 ({{ groupMembers.length || 0 }}人)</h4>
            <a-button type="primary" size="small" @click="showAddMemberModal" class="add-member-button">
              <plus-outlined /> 添加成员
            </a-button>
          </div>
          
          <a-table
            :columns="memberColumns"
            :data-source="groupMembers"
            :loading="memberLoading"
            rowKey="user_id"
            size="small"
            class="members-table"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.key === 'actions'">
                <a-popconfirm
                  title="确定要移除此成员吗？"
                  @confirm="handleRemoveMember(record)"
                  ok-text="确定"
                  cancel-text="取消"
                >
                  <a-button type="link" danger size="small" class="remove-button">
                    <template #icon><user-delete-outlined /></template>
                    移除
                  </a-button>
                </a-popconfirm>
              </template>
            </template>
          </a-table>
        </div>
      </div>
      <div v-else>
        <a-empty description="未找到用户组信息" />
      </div>
    </a-modal>
    
    <!-- 添加成员模态框 -->
    <a-modal
      v-model:visible="addMemberModalVisible"
      title="添加用户组成员"
      @ok="handleAddMember"
      :confirm-loading="confirmLoading"
      class="member-modal"
    >
      <a-form :model="memberForm" :label-col="{ span: 6 }" :wrapper-col="{ span: 16 }" class="member-form">
        <a-form-item label="用户ID" name="users" :rules="[{ required: true, message: '请输入用户ID' }]">
          <a-textarea 
            v-model:value="memberForm.users" 
            placeholder="请输入用户ID，多个ID用逗号分隔" 
            :rows="4" 
            class="member-input"
          />
          <div class="form-help-text">多个用户ID之间请用逗号分隔，例如: user1,user2,user3</div>
        </a-form-item>
      </a-form>
    </a-modal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { message } from 'ant-design-vue'
import { 
  PlusOutlined, 
  EditOutlined, 
  DeleteOutlined, 
  TeamOutlined,
  UserDeleteOutlined
} from '@ant-design/icons-vue'
import { useUserStore } from '@/store/modules/user'
import {
  getAllGroups,
  getGroupUsers,
  createGroup,
  editGroupName,
  deleteGroup,
  addGroupUsers,
  deleteGroupUsers
} from '@/api/userGroup'

// 获取用户store
const userStore = useUserStore()

// 表格列定义
const groupColumns = [
  { title: '用户组ID', dataIndex: 'group_id', key: 'group_id', width: 100 },
  { title: '用户组名称', dataIndex: 'name', key: 'name' },
  { title: '操作', key: 'actions', width: 250 }
]

const memberColumns = [
  { title: '用户ID', dataIndex: 'user_id', key: 'user_id' },
  { title: '用户昵称', dataIndex: 'nickname', key: 'nickname' },
  { title: '操作', key: 'actions', width: 100 }
]

// 状态
const loading = ref(false)
const memberLoading = ref(false)
const confirmLoading = ref(false)
const groups = ref([])
const groupMembers = ref([])
const selectedGroup = ref(null)

// 模态框状态
const createGroupModalVisible = ref(false)
const editGroupModalVisible = ref(false)
const groupDetailModalVisible = ref(false)
const addMemberModalVisible = ref(false)

// 表单数据
const groupForm = reactive({
  name: '',
  group_id: null
})

const memberForm = reactive({
  users: ''
})

// 加载用户组列表
const loadGroups = async () => {
  loading.value = true
  try {
    const response = await getAllGroups(userStore.token)
    if (response.ret === "0") {
      groups.value = response.groups || []
    } else {
      message.error('加载用户组失败')
    }
  } catch (error) {
    console.error('加载用户组错误:', error)
    message.error('加载用户组失败')
  } finally {
    loading.value = false
  }
}

// 加载用户组成员
const loadGroupMembers = async (groupId) => {
  memberLoading.value = true
  try {
    const response = await getGroupUsers(groupId, userStore.token)
    if (response.ret === "0" && response.group) {
      // 设置整个组信息，包括创建者和成员
      selectedGroup.value = response.group;
      // 设置成员列表
      groupMembers.value = response.group.users || [];
      
      console.log('加载用户组详情成功:', response.group);
      return true;
    } else {
      message.error('加载用户组成员失败');
      console.error('加载用户组详情失败:', response);
      return false;
    }
  } catch (error) {
    console.error('加载用户组成员错误:', error);
    message.error('加载用户组成员失败');
    return false;
  } finally {
    memberLoading.value = false;
  }
}

// 显示创建用户组模态框
const showCreateGroupModal = () => {
  groupForm.name = ''
  groupForm.group_id = null
  createGroupModalVisible.value = true
}

// 显示编辑用户组模态框
const showEditGroupModal = (group) => {
  groupForm.name = group.name
  groupForm.group_id = group.group_id
  editGroupModalVisible.value = true
}

// 显示用户组详情模态框
const showGroupDetailModal = async (group) => {
  // 先把基本信息设置上，以便用户可以看到组的基本信息
  selectedGroup.value = {
    ...group,
    owner: { user_id: '', nickname: '加载中...' },
    users: []
  };
  
  // 打开模态框
  groupDetailModalVisible.value = true;
  
  // 加载完整的组信息（包括创建者和成员）
  const success = await loadGroupMembers(group.group_id);
  
  if (!success) {
    message.error('加载组详情失败，请稍后重试');
  }
}

// 显示添加成员模态框
const showAddMemberModal = () => {
  memberForm.users = ''
  addMemberModalVisible.value = true
}

// 处理创建用户组
const handleCreateGroup = async () => {
  if (!groupForm.name) {
    message.error('请输入用户组名称')
    return
  }
  
  confirmLoading.value = true
  try {
    const response = await createGroup(groupForm.name, userStore.token)
    if (response.ret === "0") {
      message.success('创建用户组成功')
      createGroupModalVisible.value = false
      groupForm.name = ''
      await loadGroups()
    } else {
      message.error('创建用户组失败: ' + response.msg)
    }
  } catch (error) {
    console.error('创建用户组错误:', error)
    message.error('创建用户组失败')
  } finally {
    confirmLoading.value = false
  }
}

// 处理编辑用户组
const handleEditGroup = async () => {
  if (!groupForm.name) {
    message.error('请输入用户组名称')
    return
  }
  
  confirmLoading.value = true
  try {
    const response = await editGroupName(groupForm.group_id, groupForm.name, userStore.token)
    if (response.ret === "0") {
      message.success('编辑用户组成功')
      editGroupModalVisible.value = false
      await loadGroups()
    } else {
      message.error('编辑用户组失败: ' + response.msg)
    }
  } catch (error) {
    console.error('编辑用户组错误:', error)
    message.error('编辑用户组失败')
  } finally {
    confirmLoading.value = false
  }
}

// 删除用户组
const deleteUserGroup = async (group) => {
  loading.value = true
  try {
    const response = await deleteGroup(group.group_id, userStore.token)
    if (response.ret === "0") {
      message.success('删除用户组成功')
      await loadGroups()
    } else {
      message.error('删除用户组失败: ' + response.msg)
    }
  } catch (error) {
    console.error('删除用户组错误:', error)
    message.error('删除用户组失败')
  } finally {
    loading.value = false
  }
}

// 添加用户组成员
const handleAddMember = async () => {
  if (!memberForm.users) {
    message.error('请输入用户ID')
    return
  }
  
  // 解析用户ID
  const userIds = memberForm.users.split(/[,，\s]+/).filter(id => id.trim())
  if (userIds.length === 0) {
    message.error('请输入有效的用户ID')
    return
  }
  
  confirmLoading.value = true
  try {
    const response = await addGroupUsers(selectedGroup.value.group_id, userIds, userStore.token)
    if (response.ret === "0") {
      message.success('添加成员成功')
      addMemberModalVisible.value = false
      memberForm.users = ''
      await loadGroupMembers(selectedGroup.value.group_id)
    } else {
      message.error('添加成员失败: ' + response.msg)
    }
  } catch (error) {
    console.error('添加成员错误:', error)
    message.error('添加成员失败')
  } finally {
    confirmLoading.value = false
  }
}

// 删除用户组成员
const handleRemoveMember = async (member) => {
  memberLoading.value = true
  try {
    const response = await deleteGroupUsers(selectedGroup.value.group_id, [member.user_id], userStore.token)
    if (response.ret === "0") {
      message.success('移除成员成功')
      await loadGroupMembers(selectedGroup.value.group_id)
    } else {
      message.error('移除成员失败: ' + response.msg)
    }
  } catch (error) {
    console.error('移除成员错误:', error)
    message.error('移除成员失败')
  } finally {
    memberLoading.value = false
  }
}

// 初始化加载
onMounted(async () => {
  await loadGroups()
})
</script>

<style scoped>
.user-group-management-container {
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

.group-card {
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

.group-table {
  margin-top: 16px;
}

:deep(.ant-table) {
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.04);
}

:deep(.ant-table-thead > tr > th) {
  background: rgba(33, 150, 243, 0.05);
  font-weight: 600;
}

:deep(.ant-table-tbody > tr:hover > td) {
  background: rgba(33, 150, 243, 0.03);
}

.action-button {
  transition: all 0.2s ease;
  border-radius: 4px;
  padding: 4px 8px;
}

.action-button:hover {
  background: rgba(33, 150, 243, 0.05);
  transform: translateY(-1px);
}

.view-button {
  color: #2196F3;
}

.edit-button {
  color: #4CAF50;
}

.delete-button {
  color: #f5222d;
}

.detail-header {
  margin-bottom: 16px;
  padding: 16px;
  background: linear-gradient(120deg, rgba(33, 150, 243, 0.05), rgba(76, 175, 80, 0.05));
  border-radius: 12px;
}

.detail-header h3 {
  font-size: 20px;
  margin-bottom: 8px;
  color: #2196F3;
}

.custom-divider {
  margin: 24px 0;
  background: linear-gradient(90deg, rgba(33, 150, 243, 0.2), rgba(76, 175, 80, 0.2));
  height: 2px;
}

.members-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  padding: 0 8px;
}

.members-header h4 {
  margin: 0;
  font-size: 16px;
  color: #333;
  font-weight: 600;
}

.add-member-button {
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  box-shadow: 0 2px 4px rgba(33, 150, 243, 0.2);
  transition: all 0.3s ease;
}

.add-member-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 3px 6px rgba(33, 150, 243, 0.3);
}

.members-table {
  border-radius: 8px;
  overflow: hidden;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.03);
}

.remove-button {
  transition: all 0.2s ease;
}

.remove-button:hover {
  background: rgba(245, 34, 45, 0.05);
}

.form-help-text {
  font-size: 12px;
  color: #999;
  margin-top: 4px;
}

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
  padding: 12px 24px;
}

:deep(.ant-btn-primary) {
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  box-shadow: 0 2px 4px rgba(33, 150, 243, 0.2);
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary:hover) {
  transform: translateY(-1px);
  box-shadow: 0 3px 6px rgba(33, 150, 243, 0.3);
}

.group-detail-container {
  animation: fadeIn 0.3s ease-out;
}

.member-input {
  border-radius: 8px;
  transition: all 0.3s ease;
}

.member-input:focus {
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}
</style>