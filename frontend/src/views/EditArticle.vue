<template>
  <div class="edit-article">
      <div class="page-header">
        <h1 class="page-title">✏️ 编辑文章</h1>
        <p class="page-subtitle">修改你的创作内容</p>
      </div>

      <div class="edit-tip">
        <span class="tip-icon">💡</span>
        支持 Markdown 语法，右侧可实时预览
      </div>

    <form @submit.prevent="handleSubmit" class="article-form">
      <div class="form-group">
        <label class="form-label">标题 *</label>
        <input
          v-model="form.title"
          type="text"
          class="form-input"
          placeholder="输入文章标题..."
          required
        />
      </div>

      <div class="form-row">
        <div class="form-group">
          <label class="form-label">作者</label>
          <input
            v-model="form.author"
            type="text"
            class="form-input"
            placeholder="作者名称"
          />
        </div>
        <div class="form-group">
          <label class="form-label">标签 (逗号分隔)</label>
          <input
            v-model="form.tags"
            type="text"
            class="form-input"
            placeholder="vue, golang, web"
          />
        </div>
      </div>

      <div class="form-group">
        <label class="form-label">摘要</label>
        <textarea
          v-model="form.summary"
          class="form-textarea"
          rows="3"
          placeholder="文章摘要..."
        />
      </div>

      <div class="form-group">
        <label class="form-label">
          内容 *
          <span class="label-badge">Markdown</span>
        </label>
        <div class="split-editor">
          <div class="editor-pane">
            <div class="pane-header">
              <span class="pane-title">编辑</span>
            </div>
            <textarea
              v-model="form.content"
              class="editor-textarea"
              placeholder="开始写作..."
              required
            ></textarea>
          </div>
          <div class="preview-pane">
            <div class="pane-header">
              <span class="pane-title">预览</span>
            </div>
            <div class="preview-scroll">
              <div class="markdown-body" v-html="renderedContent"></div>
              <div v-if="!form.content" class="preview-placeholder">
                <span>输入 Markdown 内容后实时预览</span>
              </div>
            </div>
          </div>
        </div>
        <div class="editor-footer">
          <span class="toolbar-info">{{ contentLength }} 字</span>
        </div>
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-secondary" @click="$router.back()">
          ← 返回
        </button>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? '保存中...' : '💾 保存更改' }}
        </button>
      </div>
    </form>

    <div v-if="error" class="toast toast-error">
      {{ error }}
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { marked } from 'marked'
import { articleApi } from '../api'

const route = useRoute()
const router = useRouter()
const loading = ref(false)
const error = ref('')

const form = ref({
  title: '',
  content: '',
  summary: '',
  author: '',
  tags: '',
})

const contentLength = computed(() => form.value.content.length)

const renderedContent = computed(() => {
  if (!form.value.content) return ''
  return marked.parse(form.value.content, { breaks: true }) as string
})

onMounted(async () => {
  try {
    const id = Number(route.params.id)
    const res = await articleApi.getById(id)
    const article = res.data
    form.value = {
      title: article.title,
      content: article.content,
      summary: article.summary,
      author: article.author,
      tags: article.tags,
    }
  } catch {
    error.value = '加载文章失败'
  }
})

async function handleSubmit() {
  if (!form.value.title || !form.value.content) return

  loading.value = true
  error.value = ''
  try {
    const id = Number(route.params.id)
    await articleApi.update(id, {
      title: form.value.title,
      content: form.value.content,
      summary: form.value.summary,
      author: form.value.author || 'Anonymous',
      tags: form.value.tags,
    })
    router.push({ name: 'article-detail', params: { id } })
  } catch {
    error.value = '保存失败，请重试'
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.edit-article {
  max-width: 1100px;
  margin: 0 auto;
  padding: 40px 24px 80px;
}

.page-header {
  text-align: center;
  margin-bottom: 24px;
}

.page-title {
  font-size: 2.2rem;
  font-weight: 800;
  margin-bottom: 8px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.page-subtitle {
  color: #94a3b8;
  font-size: 1.05rem;
}

.edit-tip {
  text-align: center;
  margin-bottom: 32px;
  color: #64748b;
  font-size: 0.9rem;
  background: rgba(102, 126, 234, 0.06);
  border: 1px solid rgba(102, 126, 234, 0.12);
  border-radius: 10px;
  padding: 10px 20px;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  width: auto;
  margin-left: auto;
  margin-right: auto;
}

.tip-icon {
  font-size: 1rem;
}

.article-form {
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 24px;
  padding: 40px;
  backdrop-filter: blur(20px);
}

.form-group {
  margin-bottom: 28px;
}

.form-row {
  display: flex;
  gap: 20px;
}

.form-row .form-group {
  flex: 1;
}

.form-label {
  display: block;
  font-size: 0.9rem;
  font-weight: 600;
  color: #e2e8f0;
  margin-bottom: 10px;
  letter-spacing: 0.5px;
}

.label-badge {
  display: inline-block;
  font-size: 0.7rem;
  font-weight: 700;
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: #fff;
  padding: 2px 10px;
  border-radius: 20px;
  margin-left: 8px;
  vertical-align: middle;
  letter-spacing: 0.5px;
}

.form-input {
  width: 100%;
  padding: 14px 18px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 1rem;
  transition: all 0.3s;
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
}

.form-input:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15), 0 0 20px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
  color: #475569;
}

.form-textarea {
  width: 100%;
  padding: 14px 18px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 1rem;
  transition: all 0.3s;
  outline: none;
  font-family: inherit;
  box-sizing: border-box;
  resize: vertical;
  min-height: 80px;
}

.form-textarea:focus {
  border-color: #667eea;
  box-shadow: 0 0 0 3px rgba(102, 126, 234, 0.15), 0 0 20px rgba(102, 126, 234, 0.1);
}

/* Split Editor */
.split-editor {
  display: flex;
  gap: 0;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  overflow: hidden;
  background: rgba(0, 0, 0, 0.3);
}

.editor-pane,
.preview-pane {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 450px;
}

.editor-pane {
  border-right: 1px solid rgba(255, 255, 255, 0.06);
}

.pane-header {
  padding: 10px 16px;
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
  font-size: 0.8rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1px;
  color: #64748b;
}

.pane-title {
  opacity: 0.7;
}

.editor-textarea {
  flex: 1;
  width: 100%;
  padding: 18px 20px;
  background: transparent;
  border: none;
  color: #e2e8f0;
  font-size: 0.95rem;
  line-height: 1.8;
  resize: vertical;
  outline: none;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
  tab-size: 2;
  box-sizing: border-box;
}

.editor-textarea::placeholder {
  color: #475569;
}

.preview-scroll {
  flex: 1;
  overflow-y: auto;
  padding: 18px 20px;
}

.preview-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  min-height: 200px;
  color: #475569;
  font-size: 0.9rem;
}

/* Markdown preview styling */
.preview-scroll :deep(.markdown-body) {
  background: transparent;
  color: #e2e8f0;
  font-size: 0.95rem;
  line-height: 1.8;
}

.preview-scroll :deep(.markdown-body h1),
.preview-scroll :deep(.markdown-body h2),
.preview-scroll :deep(.markdown-body h3),
.preview-scroll :deep(.markdown-body h4) {
  color: #f1f5f9;
  border-bottom-color: rgba(255, 255, 255, 0.08);
  margin-top: 1.2em;
  margin-bottom: 0.6em;
}

.preview-scroll :deep(.markdown-body h1) {
  font-size: 1.8rem;
}

.preview-scroll :deep(.markdown-body h2) {
  font-size: 1.4rem;
}

.preview-scroll :deep(.markdown-body p) {
  margin-bottom: 0.8em;
}

.preview-scroll :deep(.markdown-body code) {
  background: rgba(102, 126, 234, 0.12);
  color: #a5b4fc;
  padding: 2px 8px;
  border-radius: 4px;
  font-size: 0.85em;
  font-family: 'SF Mono', 'Fira Code', 'Consolas', monospace;
}

.preview-scroll :deep(.markdown-body pre) {
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
  padding: 16px;
  overflow-x: auto;
}

.preview-scroll :deep(.markdown-body pre code) {
  background: transparent;
  padding: 0;
  color: #e2e8f0;
}

.preview-scroll :deep(.markdown-body blockquote) {
  border-left: 3px solid #667eea;
  color: #94a3b8;
  padding: 0.5em 1em;
  background: rgba(102, 126, 234, 0.05);
  border-radius: 0 8px 8px 0;
}

.preview-scroll :deep(.markdown-body ul),
.preview-scroll :deep(.markdown-body ol) {
  padding-left: 1.5em;
  margin-bottom: 0.8em;
}

.preview-scroll :deep(.markdown-body li) {
  margin-bottom: 0.3em;
}

.preview-scroll :deep(.markdown-body a) {
  color: #818cf8;
}

.preview-scroll :deep(.markdown-body hr) {
  border-color: rgba(255, 255, 255, 0.08);
  margin: 1.5em 0;
}

.preview-scroll :deep(.markdown-body img) {
  max-width: 100%;
  border-radius: 8px;
}

.preview-scroll :deep(.markdown-body table) {
  border-collapse: collapse;
  margin-bottom: 1em;
}

.preview-scroll :deep(.markdown-body th),
.preview-scroll :deep(.markdown-body td) {
  border: 1px solid rgba(255, 255, 255, 0.1);
  padding: 8px 14px;
}

.preview-scroll :deep(.markdown-body th) {
  background: rgba(255, 255, 255, 0.04);
  font-weight: 600;
}

.editor-footer {
  display: flex;
  justify-content: flex-end;
  margin-top: 8px;
}

.toolbar-info {
  font-size: 0.8rem;
  color: #64748b;
  background: rgba(0, 0, 0, 0.5);
  padding: 4px 12px;
  border-radius: 20px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 40px;
  padding-top: 28px;
  border-top: 1px solid rgba(255, 255, 255, 0.06);
}

.btn {
  padding: 14px 32px;
  border: none;
  border-radius: 12px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s;
  font-family: inherit;
}

.btn-primary {
  background: linear-gradient(135deg, #667eea, #764ba2);
  color: white;
  box-shadow: 0 4px 15px rgba(102, 126, 234, 0.3);
}

.btn-primary:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(102, 126, 234, 0.4);
}

.btn-primary:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.06);
  color: #94a3b8;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #e2e8f0;
}

.toast {
  position: fixed;
  bottom: 24px;
  right: 24px;
  padding: 14px 24px;
  border-radius: 12px;
  font-weight: 500;
  z-index: 1000;
  animation: slideUp 0.3s ease;
}

.toast-error {
  background: rgba(239, 68, 68, 0.15);
  border: 1px solid rgba(239, 68, 68, 0.3);
  color: #fca5a5;
}

@keyframes slideUp {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
</style>
