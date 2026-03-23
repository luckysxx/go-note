<template>
  <div class="paste-page">
    <div class="page-header">
      <div>
        <p class="eyebrow">Snippet Detail</p>
        <h1>{{ snippet?.title || '代码片段详情' }}</h1>
        <p class="lede">查看代码片段内容、语言与更新时间。</p>
      </div>
      <div class="header-actions">
        <el-button v-if="canEdit" size="small" plain @click="goToEdit">
          <el-icon>
            <Edit />
          </el-icon>
          <span>编辑片段</span>
        </el-button>
        <el-button size="small" type="primary" plain @click="copyContent" :disabled="!snippet?.content">
          <el-icon>
            <CopyDocument />
          </el-icon>
          <span>复制内容</span>
        </el-button>
      </div>
    </div>

    <el-card v-if="loading" class="state-card" shadow="never">
      <el-skeleton :rows="12" animated />
    </el-card>

    <el-result v-else-if="error" icon="warning" title="获取失败" :sub-title="error" class="state-card">
      <template #extra>
        <el-button type="primary" @click="retry">重试</el-button>
      </template>
    </el-result>

    <el-card v-else-if="snippet" class="viewer-card" shadow="hover">
      <div class="card-header">
        <div class="title-area">
          <div class="id-chip">#{{ snippet.id }}</div>
          <el-tag size="small" effect="plain">{{ snippet.language || 'text' }}</el-tag>
          <el-tag size="small" effect="light" :type="snippet.visibility === 'public' ? 'success' : 'info'">
            {{ snippet.visibility === 'public' ? '公开' : '私有' }}
          </el-tag>
        </div>
        <div class="meta">更新于 {{ formatDate(snippet.updated_at) }}</div>
      </div>

      <div class="code-shell">
        <div class="code-header">
          <span class="dot red"></span>
          <span class="dot yellow"></span>
          <span class="dot green"></span>
          <span class="filename">{{ snippet.title }}.{{ snippet.language || 'txt' }}</span>
        </div>
        <div class="code-body">
          <div class="gutter">
            <span v-for="(_, idx) in codeLines" :key="idx">{{ idx + 1 }}</span>
          </div>
          <pre><code v-html="highlightedCode"></code></pre>
        </div>
      </div>

      <div class="footer-actions">
        <el-tag size="small" effect="light">创建于 {{ formatDate(snippet.created_at) }}</el-tag>
        <div class="spacer" />
        <el-button size="small" @click="copyContent">
          <el-icon>
            <CopyDocument />
          </el-icon>
          <span>复制代码</span>
        </el-button>
      </div>
    </el-card>

    <el-empty v-else description="等待拉取数据" class="state-card" />
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getSnippet, type Snippet } from '@/api/paste'
import { ElMessage } from 'element-plus'
import { CopyDocument, Edit } from '@element-plus/icons-vue'
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
import typescript from 'highlight.js/lib/languages/typescript'
import python from 'highlight.js/lib/languages/python'
import go from 'highlight.js/lib/languages/go'
import java from 'highlight.js/lib/languages/java'
import cpp from 'highlight.js/lib/languages/cpp'
import rust from 'highlight.js/lib/languages/rust'
import sql from 'highlight.js/lib/languages/sql'
import json from 'highlight.js/lib/languages/json'
import yaml from 'highlight.js/lib/languages/yaml'
import markdown from 'highlight.js/lib/languages/markdown'
import xml from 'highlight.js/lib/languages/xml'
import css from 'highlight.js/lib/languages/css'
import bash from 'highlight.js/lib/languages/bash'
import 'highlight.js/styles/tokyo-night-dark.css'
import { useAuthStore } from '@/stores/auth'

// 注册语言
hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('typescript', typescript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('go', go)
hljs.registerLanguage('java', java)
hljs.registerLanguage('cpp', cpp)
hljs.registerLanguage('rust', rust)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('json', json)
hljs.registerLanguage('yaml', yaml)
hljs.registerLanguage('markdown', markdown)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('css', css)
hljs.registerLanguage('shell', bash)
hljs.registerLanguage('bash', bash)

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const loading = ref(false)
const error = ref('')
const snippet = ref<Snippet | null>(null)
const highlightedCode = ref('')

const snippetId = computed(() => (route.params.id as string) || '')
const codeLines = computed(() => snippet.value?.content?.split('\n') ?? [])
const canEdit = computed(() => {
  if (!snippet.value?.owner_id || !authStore.user?.id) {
    return true
  }
  return snippet.value.owner_id === authStore.user.id
})

const fetchPaste = async (id: string) => {
  if (!id) {
    error.value = '无效的代码片段 ID'
    return
  }
  loading.value = true
  error.value = ''
  try {
    const currentPaste = await getSnippet(id)
    if (!currentPaste.content) {
      throw new Error('后端返回缺少 content 字段')
    }
    snippet.value = currentPaste
    // 应用代码高亮
    await nextTick()
    const lang = currentPaste.language === 'text' ? 'plaintext' : currentPaste.language
    try {
      const result = hljs.highlight(currentPaste.content, { language: lang, ignoreIllegals: true })
      highlightedCode.value = result.value
    } catch {
      // 语言不支持时回退到自动检测
      const result = hljs.highlightAuto(currentPaste.content)
      highlightedCode.value = result.value
    }
  } catch {
    snippet.value = null
    error.value = '代码片段不存在或你无权查看'
  } finally {
    loading.value = false
  }
}

const retry = () => {
  fetchPaste(snippetId.value)
}

onMounted(() => {
  if (snippetId.value) {
    fetchPaste(snippetId.value)
  }
})

watch(
  () => route.params.id,
  (val) => {
    if (typeof val === 'string' && val) {
      fetchPaste(val)
    }
  },
)

const copyContent = async () => {
  if (!snippet.value?.content) return
  try {
    await navigator.clipboard.writeText(snippet.value.content)
    ElMessage.success('复制成功')
  } catch {
    ElMessage.error('复制失败')
  }
}

const goToEdit = () => {
  if (!snippetId.value) return
  router.push(`/snippets/${snippetId.value}/edit`)
}

const formatDate = (value: string) => {
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return date.toLocaleString()
}
</script>

<style scoped lang="scss">
.paste-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
  max-width: 1100px;
  margin: 0 auto;
  padding: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  background: rgba(64, 158, 255, 0.12);
  color: #2468c0;
  border-radius: 999px;
  font-weight: 700;
  width: fit-content;
}

.lede {
  margin: 6px 0 0;
  color: #596273;
  line-height: 1.5;
}

.header-actions {
  display: flex;
  gap: 8px;
  align-items: center;
}

.state-card {
  border-radius: 14px;
  border: 1px solid #e4e7ed;
  padding: 20px;
}

.viewer-card {
  border-radius: 16px;
  border: 1px solid #e4e7ed;

  :deep(.el-card__body) {
    padding: 0;
  }
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 14px 16px;
  border-bottom: 1px solid #f0f2f5;
}

.title-area {
  display: flex;
  align-items: center;
  gap: 10px;
}

.id-chip {
  padding: 6px 10px;
  background: #eef5ff;
  color: #1f2d3d;
  border-radius: 8px;
  font-weight: 700;
  letter-spacing: 0.2px;
}

.meta {
  color: #909399;
  font-size: 13px;
}

.code-shell {
  margin: 16px;
  border-radius: 12px;
  overflow: hidden;
  border: 1px solid #e4e7ed;
  background: #0f172a;
  color: #e2e8f0;
}

.code-header {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px 14px;
  background: rgba(255, 255, 255, 0.06);
  border-bottom: 1px solid rgba(255, 255, 255, 0.08);

  .dot {
    width: 10px;
    height: 10px;
    border-radius: 50%;
    display: inline-block;
  }

  .red {
    background: #ff5f56;
  }

  .yellow {
    background: #f6bd3b;
  }

  .green {
    background: #51c353;
  }

  .filename {
    margin-left: auto;
    font-family: 'Fira Code', 'Consolas', monospace;
    opacity: 0.8;
  }
}

.code-body {
  display: grid;
  grid-template-columns: auto 1fr;
  max-height: 520px;
  overflow: hidden;
}

.gutter {
  background: rgba(255, 255, 255, 0.04);
  color: rgba(226, 232, 240, 0.6);
  padding: 14px 12px;
  text-align: right;
  font-family: 'Fira Code', 'Consolas', monospace;
  font-size: 13px;
  min-width: 48px;
  line-height: 1.6;
  border-right: 1px solid rgba(255, 255, 255, 0.08);

  span {
    display: block;
  }
}

pre {
  margin: 0;
  background: transparent;
  color: inherit;
  font-family: 'Fira Code', 'Consolas', monospace;
  padding: 14px 16px;
  overflow: auto;
  white-space: pre;
  line-height: 1.6;
}

.footer-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px 16px;
  border-top: 1px solid #f0f2f5;
}

.spacer {
  flex: 1;
}

@media (max-width: 960px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .card-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .card-actions {
    align-self: flex-end;
  }
}
</style>
