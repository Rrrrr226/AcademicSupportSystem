<template>
  <div class="subjects-container">
    <div class="header">
      <h2>学科列表</h2>
      <el-button type="text" @click="$router.push('/person-info')">
        <el-icon><UserFilled /></el-icon>
        个人信息
      </el-button>
    </div>
    <!-- Existing code for subject list -->
  </div>
</template>

<script setup>
import { UserFilled } from '@element-plus/icons-vue'
// Existing script code
</script>

<style scoped>
.subjects-container {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

/* Existing styles */
</style>