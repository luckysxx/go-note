<template>
  <div class="about-container">
    <el-card class="about-card" shadow="never">
      <h1>关于 GoNote</h1>
      <p>
        GoNote 是一个基于 Go 和 Vue 3 构建的现代化代码粘贴板工具，
        已接入 user-platform 统一认证体系（SSO）。
      </p>
      <p>
        旨在提供快速、简洁的代码分享体验。
      </p>

      <el-divider />

      <h3>技术栈</h3>
      <ul>
        <li><strong>Frontend:</strong> Vue 3, TypeScript, Vite, Element Plus, Pinia</li>
        <li><strong>Backend:</strong> Go, Gin, Ent ORM, PostgreSQL, Redis</li>
        <li><strong>认证:</strong> gRPC → user-platform（SSO 统一登录）</li>
        <li><strong>配置:</strong> Viper + godotenv</li>
        <li><strong>共享模块:</strong> common（logger, errs）</li>
      </ul>
    </el-card>
  </div>
</template>

<style scoped>
.about-container {
  max-width: 800px;
  margin: 0 auto;
}

.about-card {
  padding: 20px;
  border-radius: var(--radius-lg);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-md);
}

h1 {
  margin-top: 0;
  color: var(--accent);
}

h3 {
  color: var(--text-primary);
}

ul {
  padding-left: 20px;
  line-height: 1.8;
  color: var(--text-secondary);
}
</style>
