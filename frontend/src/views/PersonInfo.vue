<template>
  <div class="person-info">
    <el-card class="info-card">
      <template #header>
        <div class="card-header">
          <span>个人信息</span>
        </div>
      </template>
      
      <el-form label-width="100px">
        <el-form-item label="用户名">
          <span>{{ userInfo.username }}</span>
        </el-form-item>
        <el-form-item label="姓名">
          <span>{{ userInfo.name }}</span>
        </el-form-item>
        <el-form-item label="邮箱">
          <span>{{ userInfo.email || '未设置' }}</span>
        </el-form-item>
        <el-form-item label="电话">
          <span>{{ userInfo.phone || '未设置' }}</span>
        </el-form-item>
      </el-form>

      <el-button type="primary" @click="showChangePasswordDialog">修改密码</el-button>
    </el-card>

    <el-dialog v-model="dialogVisible" title="修改密码" width="30%">
      <el-form :model="passwordForm" :rules="rules" ref="passwordFormRef" label-width="100px">
        <el-form-item label="新密码" prop="password">
          <el-input v-model="passwordForm.password" type="password" show-password />
        </el-form-item>
        <el-form-item label="确认密码" prop="confirmPassword">
          <el-input v-model="passwordForm.confirmPassword" type="password" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">取消</el-button>
          <el-button type="primary" @click="handleChangePassword">确认</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getPersonInfo, modifyUserInfo } from '@/api/user'

const userInfo = ref({})
const dialogVisible = ref(false)
const passwordFormRef = ref(null)
const passwordForm = ref({
  password: '',
  confirmPassword: ''
})

const rules = {
  password: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, message: '密码长度至少为6位', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    {
      validator: (rule, value, callback) => {
        if (value !== passwordForm.value.password) {
          callback(new Error('两次输入密码不一致'))
        } else {
          callback()
        }
      },
      trigger: 'blur'
    }
  ]
}

const fetchUserInfo = async () => {
  try {
    const userId = localStorage.getItem('userId')
    const response = await getPersonInfo({ userId })
    userInfo.value = response.data
  } catch (error) {
    ElMessage.error('获取用户信息失败')
  }
}

const showChangePasswordDialog = () => {
  dialogVisible.value = true
  passwordForm.value = {
    password: '',
    confirmPassword: ''
  }
}

const handleChangePassword = async () => {
  await passwordFormRef.value.validate(async (valid) => {
    if (valid) {
      try {
        const userId = localStorage.getItem('userId')
        await modifyUserInfo({
          userId,
          password: passwordForm.value.password
        })
        ElMessage.success('密码修改成功')
        dialogVisible.value = false
      } catch (error) {
        ElMessage.error('密码修改失败')
      }
    }
  })
}

onMounted(() => {
  fetchUserInfo()
})
</script>

<style scoped>
.person-info {
  padding: 20px;
}

.info-card {
  max-width: 600px;
  margin: 0 auto;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>