<template>
  <div class="user-management-container">
    <a-card title="用户管理" class="user-card">
      <template #extra>
        <a-button type="primary" @click="showAddUserModal" class="add-button">
          <plus-outlined /> 添加用户
        </a-button>
      </template>
      
      <a-table
        :columns="userColumns"
        :data-source="users"
        :loading="loading"
        row-key="user_id"
        :pagination="{ pageSize: 10 }"
        class="user-table"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'type'">
            <a-tag :color="getUserTypeColor(record.user_type)" class="user-type-tag">
              {{ getUserTypeName(record.user_type) }}
            </a-tag>
          </template>
          <template v-if="column.key === 'action'">
            <a-space>
              <a-button type="link" @click="showEditUserModal(record)" class="action-button edit-button">
                <template #icon><EditOutlined /></template>
                编辑
              </a-button>
              <a-button type="link" @click="showChangePasswordModal(record)" class="action-button password-button">
                <template #icon><LockOutlined /></template>
                修改密码
              </a-button>
              <a-popconfirm
                title="确定要删除该用户吗？"
                @confirm="handleDeleteUser(record)"
                ok-text="确定"
                cancel-text="取消"
                class="delete-confirm"
              >
                <a-button type="link" danger class="action-button delete-button">
                  <template #icon><DeleteOutlined /></template>
                  删除
                </a-button>
              </a-popconfirm>
            </a-space>
          </template>
        </template>
      </a-table>
    </a-card>

    <!-- 用户表单模态框 -->
    <a-modal
      v-model:visible="userModalVisible"
      :title="userModalTitle"
      @ok="handleUserModalOk"
      :confirm-loading="modalLoading"
      class="user-modal"
    >
      <a-form
        ref="userFormRef"
        :model="userForm"
        :rules="userRules"
        :label-col="{ span: 6 }"
        :wrapper-col="{ span: 16 }"
        class="user-form"
      >
        <a-form-item label="用户ID" name="user_id" v-if="isEdit">
          <a-input v-model:value="userForm.user_id" disabled />
        </a-form-item>
        <a-form-item label="用户ID" name="user_id" v-else>
          <a-input v-model:value="userForm.user_id" placeholder="请输入用户ID（字母、数字）" />
          <div class="form-help-text">用户ID将作为登录账号，创建后不可修改</div>
        </a-form-item>
        <a-form-item label="昵称" name="nickname">
          <a-input v-model:value="userForm.nickname" placeholder="请输入用户昵称" />
        </a-form-item>
        <a-form-item
          label="密码"
          name="password"
          :rules="[{ required: !isEdit, message: '请输入密码' }]"
          v-if="!isEdit"
        >
          <a-input-password v-model:value="userForm.password" placeholder="请输入密码" />
          <div class="form-help-text">建议使用包含字母、数字的组合</div>
        </a-form-item>
        <a-form-item label="用户类型" name="user_type">
          <a-select v-model:value="userForm.user_type" placeholder="请选择用户类型">
            <a-select-option value="student">学生</a-select-option>
            <a-select-option value="teacher">教师</a-select-option>
            <a-select-option value="admin">管理员</a-select-option>
          </a-select>
        </a-form-item>
      </a-form>
    </a-modal>

    <!-- 修改密码模态框 -->
    <ChangePassword 
      v-model:visible="changePasswordVisible" 
      @success="handlePasswordChangeSuccess" 
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, onActivated } from 'vue'
import { message } from 'ant-design-vue'
import { PlusOutlined, EditOutlined, DeleteOutlined, LockOutlined } from '@ant-design/icons-vue'
import { getUserList, addUser, updateUser, deleteUser } from '@/api/user'
import { useUserStore } from '@/store/modules/user'
import ChangePassword from '@/components/user/ChangePassword.vue'

// 获取用户store
const userStore = useUserStore()

// 表格列定义
const userColumns = [
  { title: '用户ID', dataIndex: 'user_id', width: 120 },
  { title: '昵称', dataIndex: 'nickname' },
  { title: '用户类型', dataIndex: 'user_type', key: 'type', width: 100 },
  { title: '操作', key: 'action', width: 280 }
]

// 状态管理
const loading = ref(false)
const modalLoading = ref(false)
const users = ref([])

// 模态框状态
const userModalVisible = ref(false)
const userModalTitle = ref('')
const isEdit = ref(false)
const changePasswordVisible = ref(false)

// 表单引用
const userFormRef = ref()

// 表单数据
const userForm = reactive({
  user_id: '',
  nickname: '',
  password: '',
  user_type: 'student'
})

// 表单验证规则
const userRules = {
  user_id: [
    { required: true, message: '请输入用户ID' },
    { pattern: /^[a-zA-Z0-9_]+$/, message: '用户ID只能包含字母、数字和下划线' }
  ],
  nickname: [{ required: true, message: '请输入昵称' }],
  user_type: [{ required: true, message: '请选择用户类型' }]
}

// 工具函数
const getUserTypeColor = (type) => {
  const typeMap = {
    'student': 'green',
    'teacher': 'blue',
    'admin': 'red'
  }
  return typeMap[type] || 'default'
}

const getUserTypeName = (type) => {
  const typeMap = {
    'student': '学生',
    'teacher': '教师',
    'admin': '管理员'
  }
  return typeMap[type] || '未知'
}

// 加载用户列表
const loadUsers = async () => {
  loading.value = true
  try {
    const response = await getUserList(userStore.token)
    if (response.code === 200) {
      users.value = response.data || []
    } else {
      message.error('加载用户列表失败: ' + (response.message || '未知错误'))
    }
  } catch (error) {
    console.error('加载用户列表错误:', error)
    message.error('加载用户列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 显示修改密码模态框
const showChangePasswordModal = (record) => {
  // 设置当前用户ID到userStore，以便ChangePassword组件使用
  userStore.setCurrentUserId(record.user_id)
  changePasswordVisible.value = true
}

// 处理密码修改成功回调
const handlePasswordChangeSuccess = () => {
  message.success('密码修改成功')
  // 清除当前设置的用户ID
  userStore.setCurrentUserId('')
}

const handleDeleteUser = async (record) => {
  try {
    loading.value = true
    const response = await deleteUser(record.user_id, userStore.token)
    if (response.code === 200) {
      message.success('删除用户成功')
      await loadUsers() // 重新加载用户列表
    } else {
      message.error('删除用户失败: ' + (response.message || '未知错误'))
    }
  } catch (error) {
    console.error('删除用户错误:', error)
    message.error('删除用户失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 模态框相关方法
const showAddUserModal = () => {
  isEdit.value = false
  userForm.user_id = ''
  userForm.nickname = ''
  userForm.password = ''
  userForm.user_type = 'student'
  userModalTitle.value = '添加用户'
  userModalVisible.value = true
}

const showEditUserModal = (record) => {
  isEdit.value = true
  userForm.user_id = record.user_id
  userForm.nickname = record.nickname
  userForm.password = ''
  userForm.user_type = record.user_type
  userModalTitle.value = '编辑用户'
  userModalVisible.value = true
}

const handleUserModalOk = async () => {
  try {
    await userFormRef.value.validate()
    modalLoading.value = true
    
    if (isEdit.value) {
      // 编辑用户
      const response = await updateUser({
        user_id: userForm.user_id,
        nickname: userForm.nickname,
        user_type: userForm.user_type
      }, userStore.token)
      
      if (response.code === 200) {
        message.success('更新用户成功')
        userModalVisible.value = false
        await loadUsers() // 重新加载用户列表
      } else {
        message.error('更新用户失败: ' + (response.message || '未知错误'))
      }
    } else {
      // 添加用户
      const response = await addUser({
        user_id: userForm.user_id,
        nickname: userForm.nickname,
        password: userForm.password,
        user_type: userForm.user_type
      }, userStore.token)
      
      if (response.code === 200) {
        message.success('添加用户成功')
        userModalVisible.value = false
        await loadUsers() // 重新加载用户列表
      } else {
        message.error('添加用户失败: ' + (response.message || '未知错误'))
      }
    }
  } catch (error) {
    console.error('表单验证或提交错误:', error)
    message.error('操作失败: ' + (error.message || '未知错误'))
  } finally {
    modalLoading.value = false
  }
}

// 加载数据
onMounted(async () => {
  await loadUsers()
})

// 当组件被激活时重新加载数据
onActivated(async () => {
  await loadUsers()
})
</script>

<style scoped>
.user-management-container {
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

.user-card {
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

.add-button {
  height: 40px;
  border-radius: 8px;
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  transition: all 0.3s ease;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(33, 150, 243, 0.2);
}

.add-button:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.3);
  background: linear-gradient(120deg, #1976D2, #388E3C);
}

:deep(.ant-table) {
  background: transparent;
  border-radius: 12px;
  overflow: hidden;
  margin-top: 16px;
}

:deep(.ant-table-thead > tr > th) {
  background: rgba(33, 150, 243, 0.05);
  color: #2196F3;
  font-weight: 600;
  border-bottom: 1px solid rgba(33, 150, 243, 0.1);
  transition: background 0.3s ease;
}

:deep(.ant-table-tbody > tr > td) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.04);
  transition: all 0.3s ease;
}

:deep(.ant-table-tbody > tr:hover > td) {
  background: rgba(33, 150, 243, 0.05);
}

.user-type-tag {
  border-radius: 12px;
  padding: 2px 12px;
  font-weight: 500;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.action-button {
  border-radius: 6px;
  transition: all 0.3s ease;
  padding: 4px 8px;
}

.action-button:hover {
  background: rgba(33, 150, 243, 0.1);
  transform: translateY(-1px);
}

.delete-button:hover {
  background: rgba(244, 67, 54, 0.1);
}

.user-modal {
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
:deep(.ant-input-password),
:deep(.ant-select-selector) {
  border-radius: 8px;
  padding: 8px 12px;
  border: 1px solid rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

:deep(.ant-input:hover),
:deep(.ant-input-password:hover),
:deep(.ant-select-selector:hover) {
  border-color: #4CAF50;
}

:deep(.ant-input:focus),
:deep(.ant-input-password:focus),
:deep(.ant-select-selector:focus),
:deep(.ant-input-affix-wrapper-focused) {
  border-color: #2196F3;
  box-shadow: 0 0 0 2px rgba(33, 150, 243, 0.2);
}

.form-help-text {
  font-size: 12px;
  color: #8c9bab;
  margin-top: 4px;
  transition: opacity 0.3s ease;
}

:deep(.ant-btn-primary) {
  background: linear-gradient(120deg, #2196F3, #4CAF50);
  border: none;
  height: 40px;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.3s ease;
}

:deep(.ant-btn-primary:hover) {
  background: linear-gradient(120deg, #1976D2, #388E3C);
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(33, 150, 243, 0.2);
}

:deep(.ant-form-item) {
  margin-bottom: 24px;
  transition: all 0.3s ease;
}
</style>