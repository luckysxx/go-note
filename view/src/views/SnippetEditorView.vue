<template>
  <div class="editor-page">
    <div class="page-header">
      <div>
        <p class="eyebrow">{{ isEditMode ? '编辑片段' : '新建片段' }}</p>
        <h1>{{ isEditMode ? '更新你的代码片段' : '创建一个新的代码片段' }}</h1>
      </div>
      <el-space>
        <el-button @click="goBack">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveSnippet">保存</el-button>
      </el-space>
    </div>

    <el-card class="editor-card" shadow="never" v-loading="loading">
      <el-form label-position="top">
        <el-form-item label="标题">
          <el-input v-model="form.title" maxlength="120" show-word-limit placeholder="例如：用户鉴权中间件" />
        </el-form-item>

        <div class="meta-row">
          <el-form-item label="语言" class="meta-item">
            <el-select v-model="form.language" filterable>
              <el-option label="Plain Text" value="text" />
              <el-option label="Go" value="go" />
              <el-option label="JavaScript" value="javascript" />
              <el-option label="TypeScript" value="typescript" />
              <el-option label="Python" value="python" />
              <el-option label="Java" value="java" />
              <el-option label="C++" value="cpp" />
              <el-option label="Rust" value="rust" />
              <el-option label="SQL" value="sql" />
              <el-option label="JSON" value="json" />
              <el-option label="YAML" value="yaml" />
              <el-option label="Markdown" value="markdown" />
              <el-option label="HTML" value="html" />
              <el-option label="CSS" value="css" />
              <el-option label="Shell" value="shell" />
            </el-select>
          </el-form-item>

          <el-form-item label="可见性" class="meta-item">
            <el-select v-model="form.visibility">
              <el-option label="私有" value="private" />
              <el-option label="公开" value="public" />
            </el-select>
          </el-form-item>
        </div>

        <el-form-item label="代码内容">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="18"
            resize="none"
            spellcheck="false"
            class="code-input"
            placeholder="粘贴或输入你的代码..."
          />
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { createSnippet, getSnippet, updateSnippet } from '@/api/paste'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const saving = ref(false)

const snippetId = computed(() => (route.params.id as string) || '')
const isEditMode = computed(() => !!snippetId.value)

const form = reactive({
  title: '',
  content: '',
  language: 'text',
  visibility: 'private' as 'private' | 'public',
})

const loadSnippet = async () => {
  if (!isEditMode.value) {
    return
  }

  loading.value = true
  try {
    const snippet = await getSnippet(snippetId.value)
    form.title = snippet.title
    form.content = snippet.content
    form.language = snippet.language
    form.visibility = snippet.visibility ?? 'private'
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '加载失败')
    router.replace('/')
  } finally {
    loading.value = false
  }
}

const saveSnippet = async () => {
  if (!form.title.trim()) {
    ElMessage.warning('请输入标题')
    return
  }

  if (!form.content.trim()) {
    ElMessage.warning('请输入代码内容')
    return
  }

  saving.value = true
  try {
    if (isEditMode.value) {
      const updated = await updateSnippet(snippetId.value, {
        title: form.title,
        content: form.content,
        language: form.language,
        visibility: form.visibility,
      })
      ElMessage.success('更新成功')
      router.replace(`/snippets/${updated.id}`)
    } else {
      const created = await createSnippet({
        title: form.title,
        content: form.content,
        language: form.language,
        visibility: form.visibility,
      })
      ElMessage.success('创建成功')
      router.replace(`/snippets/${created.id}`)
    }
  } catch (err) {
    ElMessage.error(err instanceof Error ? err.message : '保存失败')
  } finally {
    saving.value = false
  }
}

const goBack = () => {
  if (isEditMode.value) {
    router.push(`/snippets/${snippetId.value}`)
    return
  }
  router.push('/')
}

onMounted(() => {
  loadSnippet()
})
</script>

<style scoped lang="scss">
.editor-page {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;

  h1 {
    margin: 8px 0 0;
    font-size: 28px;
    color: #1f2d3d;
  }
}

.eyebrow {
  margin: 0;
  display: inline-flex;
  align-items: center;
  padding: 4px 10px;
  border-radius: 999px;
  background: rgba(103, 194, 58, 0.15);
  color: #3a7d1a;
  font-weight: 700;
  font-size: 12px;
}

.editor-card {
  border-radius: 12px;

  :deep(.el-card__body) {
    padding: 20px;
  }
}

.meta-row {
  display: grid;
  grid-template-columns: repeat(2, minmax(180px, 240px));
  gap: 12px;
}

.meta-item {
  margin-bottom: 0;
}

.code-input {
  :deep(.el-textarea__inner) {
    font-family: 'Fira Code', 'Consolas', monospace;
    font-size: 14px;
    line-height: 1.6;
  }
}

@media (max-width: 860px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .meta-row {
    grid-template-columns: 1fr;
  }
}
</style>
