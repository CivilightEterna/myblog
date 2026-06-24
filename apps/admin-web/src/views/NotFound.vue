<template>
  <div class="not-found">
    <h1>404</h1>
    <p>页面未找到</p>
    <el-button type="primary" @click="$router.push('/')">返回仪表盘</el-button>
  </div>
</template>

<style scoped>
.not-found {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: #f3f4f6;
}

.not-found h1 {
  font-size: 6rem;
  font-weight: 800;
  color: #8b5cf6;
  line-height: 1;
}

.not-found p {
  font-size: 1.25rem;
  color: #6b7280;
  margin-bottom: 2rem;
}
</style>
