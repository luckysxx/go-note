<template>
  <div class="home-page">
    <section class="hero">
      <div class="hero-content">
        <p class="eyebrow">
          <el-icon><Connection /></el-icon>
          <span>My Snippets</span>
        </p>
        <h1>管理你的<span class="gradient-text">代码片段</span></h1>
        <p class="lede">像 GitHub Gist 一样浏览、搜索、编辑你自己的代码片段。</p>
      </div>
      <div class="hero-actions">
        <el-input
          v-model="keyword"
          placeholder="按标题搜索..."
          clearable
          class="search-input"
          @input="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
        <el-button type="primary" class="create-btn" @click="router.push('/snippets/new')">
          <el-icon><Plus /></el-icon>
          <span>新建片段</span>
        </el-button>
      </div>
    </section>

    <div class="list-section">
      <template v-if="loading">
        <el-skeleton :rows="8" animated />
      </template>

      <el-result v-else-if="error" icon="warning" title="加载失败" :sub-title="error">
        <template #extra>
          <el-button type="primary" @click="fetchSnippets">重试</el-button>
        </template>
      </el-result>

      <div v-else-if="filteredSnippets.length === 0" class="empty-state">
        <div class="empty-icon">📝</div>
        <h3>暂无代码片段</h3>
        <p>创建你的第一个代码片段，开始分享吧</p>
        <el-button type="primary" @click="router.push('/snippets/new')">创建第一个片段</el-button>
      </div>

      <div v-else class="snippet-grid">
        <div
          v-for="snippet in filteredSnippets"
          :key="snippet.id"
          class="snippet-card"
          @click="openSnippet(snippet.id)"
        >
          <div class="card-header">
            <span class="lang-badge">{{ snippet.language }}</span>
            <el-tag size="small" :type="snippet.visibility === 'public' ? 'success' : 'info'">
              {{ snippet.visibility === 'public' ? '公开' : '私有' }}
            </el-tag>
          </div>
          <h3 class="card-title">{{ snippet.title }}</h3>
          <p class="card-time">{{ formatDate(snippet.updated_at) }}</p>
          <div class="card-actions">
            <el-button text size="small" @click.stop="openSnippet(snippet.id)">
              <el-icon><View /></el-icon>
              查看
            </el-button>
            <el-button text size="small" @click.stop="editSnippet(snippet.id)">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import { Connection, Edit, Plus, Search, View } from '@element-plus/icons-vue'
import { listMySnippets, type Snippet } from '@/api/paste'

const router = useRouter()
const loading = ref(false)
const error = ref('')
const keyword = ref('')
const snippets = ref<Snippet[]>([])

const filteredSnippets = computed(() => {
  const currentKeyword = keyword.value.trim().toLowerCase()
  if (!currentKeyword) {
    return snippets.value
  }
  return snippets.value.filter((item) => item.title.toLowerCase().includes(currentKeyword))
})

const fetchSnippets = async () => {
  loading.value = true
  error.value = ''
  try {
    snippets.value = await listMySnippets()
  } catch (err) {
    snippets.value = []
    error.value = err instanceof Error ? err.message : '无法加载你的代码片段'
  } finally {
    loading.value = false
  }
}

const openSnippet = (id: string | number) => {
  router.push(`/snippets/${id}`)
}

const editSnippet = (id: string | number) => {
  router.push(`/snippets/${id}/edit`)
}

const handleSearch = () => {
  return
}

const formatDate = (value: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}

onMounted(() => {
  fetchSnippets()
})
</script>

<style scoped lang="scss">
.home-page {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.hero {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-xl);
  padding: 32px;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 24px;
  box-shadow: var(--shadow-md);
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: -50%;
    right: -20%;
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, var(--accent-glow), transparent 70%);
    pointer-events: none;
    opacity: 0.6;
  }

  .hero-content {
    position: relative;
    z-index: 1;
  }

  h1 {
    margin: 12px 0 8px;
    font-size: 32px;
    font-weight: 800;
    color: var(--text-primary);
    letter-spacing: -0.5px;
  }

  .gradient-text {
    background: linear-gradient(135deg, var(--accent), #7c3aed);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
  }

  .lede {
    margin: 0;
    color: var(--text-secondary);
    line-height: 1.6;
  }
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  padding: 6px 14px;
  border-radius: 999px;
  background: var(--accent-glow);
  color: var(--accent);
  font-weight: 700;
  font-size: 12px;
  letter-spacing: 0.3px;
  margin: 0;
}

.hero-actions {
  display: flex;
  gap: 12px;
  align-items: center;
  position: relative;
  z-index: 1;
  flex-shrink: 0;

  .search-input {
    width: 260px;
  }

  .create-btn {
    background: linear-gradient(135deg, var(--accent), #7c3aed);
    border: none;
    font-weight: 600;
    border-radius: var(--radius-sm);
    transition: all 0.2s ease;

    &:hover {
      transform: translateY(-1px);
      box-shadow: 0 8px 24px var(--accent-glow);
    }
  }
}

.list-section {
  min-height: 200px;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;

  .empty-icon {
    font-size: 48px;
    margin-bottom: 12px;
  }

  h3 {
    color: var(--text-primary);
    margin: 0 0 8px;
    font-size: 20px;
  }

  p {
    color: var(--text-muted);
    margin: 0 0 20px;
  }
}

.snippet-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.snippet-card {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  padding: 20px;
  cursor: pointer;
  transition: all 0.25s ease;
  box-shadow: var(--shadow-sm);

  &:hover {
    border-color: var(--accent);
    transform: translateY(-2px);
    box-shadow: 0 12px 32px rgba(99, 102, 241, 0.1);
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;
  }

  .lang-badge {
    background: var(--bg-subtle);
    color: var(--accent);
    padding: 4px 10px;
    border-radius: 6px;
    font-size: 12px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
  }

  .card-title {
    margin: 0 0 8px;
    font-size: 16px;
    font-weight: 600;
    color: var(--text-primary);
    line-height: 1.4;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .card-time {
    margin: 0 0 12px;
    font-size: 13px;
    color: var(--text-muted);
  }

  .card-actions {
    display: flex;
    gap: 8px;
    border-top: 1px solid var(--border-light);
    padding-top: 12px;
    margin-top: auto;
  }
}

@media (max-width: 860px) {
  .hero {
    flex-direction: column;
  }

  .hero-actions {
    width: 100%;

    .search-input {
      width: 100%;
    }
  }

  .snippet-grid {
    grid-template-columns: 1fr;
  }
}
</style>
