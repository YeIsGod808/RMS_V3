<template>
  <a-modal
    class="password-modal"
    :visible="visible"
    title="修改密码"
    @ok="handleSubmit"
    ok-text="确认"
    @cancel="handleCancel"
    cancel-text="取消"
    :confirmLoading="loading"
    :width="460"
  >
    <a-form
      :model="formState"
      :rules="rules"
      ref="formRef"
      layout="vertical"
    >
      <a-form-item name="oldPassword" label="当前密码" v-if="!isAdminChangingOther">
        <a-input-password 
          v-model:value="formState.oldPassword" 
          placeholder="请输入当前密码"
          class="custom-input"
          size="large"
        >
          <template #prefix>
            <LockOutlined class="input-icon" />
          </template>
        </a-input-password>
      </a-form-item>
      
      <a-form-item name="newPassword" label="新密码">
        <a-input-password 
          v-model:value="formState.newPassword" 
          placeholder="请输入新密码"
          class="custom-input"
          size="large"
        >
          <template #prefix>
            <SafetyOutlined class="input-icon" />
          </template>
        </a-input-password>
      </a-form-item>
      
      <a-form-item name="confirmPassword" label="确认新密码">
        <a-input-password 
          v-model:value="formState.confirmPassword" 
          placeholder="请再次输入新密码"
          class="custom-input"
          size="large"
        >
          <template #prefix>
            <CheckCircleOutlined class="input-icon" />
          </template>
        </a-input-password>
      </a-form-item>
    </a-form>
  </a-modal>
</template>

<style scoped lang="scss">
.password-modal {
  :deep(.ant-modal-content) {
    border-radius: 16px;
    overflow: hidden;
    background: rgba(255, 255, 255, 0.98);
    backdrop-filter: blur(10px);
    box-shadow: 0 10px 40px rgba(0, 0, 0, 0.15);
    transform: translateY(0);
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      transform: translateY(-2px);
      box-shadow: 0 15px 50px rgba(0, 0, 0, 0.2);
    }
  }

  :deep(.ant-modal-header) {
    background: transparent;
    border-bottom: none;
    padding: 28px 32px 16px;

    .ant-modal-title {
      font-size: 26px;
      font-weight: 600;
      color: #1a1a2e;
      text-align: center;
      letter-spacing: 0.8px;
      font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
    }
  }

  :deep(.ant-modal-body) {
    padding: 16px 32px 32px;
  }

  :deep(.ant-form-item) {
    margin-bottom: 28px;

    .ant-form-item-label {
      padding-bottom: 10px;
      
      label {
        color: #2c3e50;
        font-weight: 500;
        font-size: 16px;
        transition: color 0.3s ease;
        font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;
      }
    }

    &:hover .ant-form-item-label label {
      color: #1890ff;
    }
  }

  .custom-input {
    border-radius: 14px;
    border: 2px solid #e8e8e8;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);

    &:hover {
      border-color: #40a9ff;
      box-shadow: 0 0 0 3px rgba(24, 144, 255, 0.1);
      transform: translateY(-1px);
    }

    &:focus-within {
      border-color: #1890ff;
      box-shadow: 0 0 0 4px rgba(24, 144, 255, 0.15);
      transform: translateY(-1px);
    }

    :deep(.ant-input) {
      padding: 14px 14px 14px 48px;
      font-size: 15px;
      font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;

      &::placeholder {
        color: #a0aec0;
        font-size: 14px;
      }
    }

    .input-icon {
      color: #a0aec0;
      font-size: 20px;
      margin-left: 4px;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    }

    &:hover .input-icon,
    &:focus-within .input-icon {
      color: #1890ff;
      transform: scale(1.1);
    }
  }

  :deep(.ant-modal-footer) {
    border-top: none;
    padding: 0 32px 28px;
    text-align: center;

    .ant-btn {
      height: 46px;
      padding: 0 32px;
      font-size: 16px;
      border-radius: 14px;
      font-weight: 500;
      transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
      font-family: 'PingFang SC', 'Microsoft YaHei', sans-serif;

      &.ant-btn-primary {
        background: linear-gradient(135deg, #1890ff 0%, #0056b3 100%);
        border: none;
        box-shadow: 0 4px 15px rgba(24, 144, 255, 0.25);

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 6px 20px rgba(24, 144, 255, 0.35);
        }

        &:active {
          transform: translateY(1px);
        }
      }

      &.ant-btn-default {
        border: 2px solid #e8e8e8;
        color: #666;

        &:hover {
          border-color: #1890ff;
          color: #1890ff;
          background: rgba(24, 144, 255, 0.05);
          transform: translateY(-1px);
        }
      }
    }
  }
}
</style>
<script setup>
import { ref, reactive, defineProps, defineEmits, computed } from 'vue'
import { message } from 'ant-design-vue'
import { changePassword } from '@/api/user'
import { useUserStore } from '@/store/modules/user'
import { 
  LockOutlined, 
  SafetyOutlined, 
  CheckCircleOutlined 
} from '@ant-design/icons-vue'

const props = defineProps({
  visible: Boolean
})

const emit = defineEmits(['update:visible', 'success'])
const userStore = useUserStore()
const formRef = ref(null)
const loading = ref(false)

// 判断是否是管理员为他人修改密码
const isAdminChangingOther = computed(() => 
  userStore.currentUserId && userStore.isAdmin
)

const formState = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 动态计算表单验证规则
const rules = computed(() => {
  const baseRules = {
    newPassword: [
      { required: true, message: '请输入新密码' },
      { pattern: /^([a-z]|[A-Z]|[0-9]){8,16}$/, message: '密码只能包含字母和数字，长度8-16位' }
    ],
    confirmPassword: [
      { required: true, message: '请确认新密码' },
      { validator: validateConfirmPassword }
    ]
  }

  // 只有在不是管理员修改他人密码时才需要原密码
  if (!isAdminChangingOther.value) {
    baseRules.oldPassword = [
      { required: true, message: '请输入当前密码' }
    ]
  }

  return baseRules
})

// 验证两次密码是否一致
function validateConfirmPassword(rule, value) {
  if (value !== formState.newPassword) {
    return Promise.reject('两次输入的密码不一致')
  }
  return Promise.resolve()
}

// 提交表单
const handleSubmit = async () => {
  try {
    await formRef.value.validate()
    loading.value = true
    
    // 获取要修改密码的用户ID
    const targetUserId = userStore.currentUserId || userStore.username
    
    // 调用修改密码API，确保传递token
    // 如果是管理员为他人修改密码，传递目标用户ID
    const response = await changePassword(
      userStore.token, 
      formState.newPassword, 
      isAdminChangingOther.value ? targetUserId : null
    )
    
    
    if (response.ret === '0') {
      //message.success('密码修改成功')
      resetForm()
      emit('success')
      emit('update:visible', false)
      
      // 如果是管理员为他人修改密码，清除currentUserId
      if (isAdminChangingOther.value) {
        userStore.setCurrentUserId('')
      }
    } else {
      message.error(response.msg || '密码修改失败')
    }
  } catch (error) {
    if (error.response && error.response.status === 401) {
      message.error('未授权，请重新登录')
    } else if (error.message) {
      message.error(error.message)
    }
  } finally {
    loading.value = false
  }
}

// 取消操作
const handleCancel = () => {
  resetForm()
  emit('update:visible', false)
}

// 重置表单
const resetForm = () => {
  formRef.value?.resetFields()
  Object.assign(formState, {
    oldPassword: '',
    newPassword: '',
    confirmPassword: ''
  })
}
</script>