<template>
  <div class="detail-page">
    <!-- Loading -->
    <div v-if="loading" class="loading-container">
      <div class="loader"></div>
    </div>

    <!-- Error -->
    <div v-else-if="error" class="error-container">
      <div class="error-icon">⚠️</div>
      <h2>{{ error }}</h2>
      <router-link to="/" class="back-link">← 返回首页</router-link>
    </div>

    <!-- Article Content -->
    <article v-else-if="article" class="article-container">
      <!-- Back button -->
      <router-link to="/" class="back-btn">
        <span class="back-icon">←</span>
        <span>返回文章列表</span>
      </router-link>

      <!-- Cover Image -->
      <div v-if="article.cover_url" class="cover-image-wrapper">
        <div class="cover-image-container">
          <img :src="article.cover_url" :alt="article.title" class="cover-image" />
          <div class="cover-overlay"></div>
        </div>
      </div>

      <!-- Article Header -->
      <div class="article-header" :class="{ 'no-cover': !article.cover_url }">
        <div class="header-content">
          <h1 class="article-title">{{ article.title }}</h1>
          <div class="article-meta">
            <div class="meta-item">
              <span class="meta-icon">✍️</span>
              <span>{{ article.author || 'Anonymous' }}</span>
            </div>
            <div class="meta-item">
              <span class="meta-icon">📅</span>
              <span>{{ formatDate(article.created_at) }}</span>
            </div>
            <div v-if="article.tags" class="tags">
              <span v-for="tag in parseTags(article.tags)" :key="tag" class="tag">
                #{{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Article Body -->
      <div class="article-body glass-card">
        <div class="article-content" v-html="renderContent(article.content)"></div>
      </div>

      <!-- Action Buttons -->
      <div class="action-bar glass-card">
        <router-link :to="`/edit/${article.id}`" class="action-btn edit-btn">
          ✏️ 编辑文章
        </router-link>
        <button @click="confirmDelete" class="action-btn delete-btn">
          🗑️ 删除文章
        </button>
      </div>
    </article>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { articleApi } from '../api'
import type { Article } from '../types'

const route = useRoute()
const router = useRouter()

const article = ref<Article | null>(null)
const loading = ref(true)
const error = ref('')

const fetchArticle = async () => {
  const id = Number(route.params.id)
  if (!id) {
    error.value = '无效的文章ID'
    loading.value = false
    return
  }

  try {
    const res = await articleApi.getById(id)
    article.value = res.data
  } catch (e: any) {
    error.value = e.response?.data?.error || '加载文章失败'
  } finally {
    loading.value = false
  }
}

const formatDate = (dateStr: string) => {
  const d = new Date(dateStr)
  return d.toLocaleDateString('zh-CN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  })
}

const parseTags = (tags: string) => {
  return tags
    .split(/[,，、\s]+/)
    .map((t) => t.trim())
    .filter(Boolean)
}

const renderContent = (content: string) => {
  return content
    .replace(/\n/g, '<br>')
    .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
    .replace(/\*(.*?)\*/g, '<em>$1</em>')
}

const confirmDelete = async () => {
  if (!confirm('确定要删除这篇文章吗？')) return
  try {
    await articleApi.delete(article.value!.id)
    router.push('/')
  } catch (e: any) {
    alert(e.response?.data?.error || '删除失败')
  }
}

onMounted(fetchArticle)
</script>

<style scoped>
.detail-page {
  max-width: 900px;
  margin: 0 auto;
  padding: 2rem 1.5rem;
}

.loading-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 50vh;
}

.loader {
  width: 48px;
  height: 48px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top-color: #6366f1;
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.error-container {
  text-align: center;
  padding: 4rem 2rem;
}

.error-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
}

.error-container h2 {
  color: #f87171;
  margin-bottom: 1rem;
}

.back-link {
  color: #818cf8;
  text-decoration: none;
  font-size: 1.1rem;
  transition: color 0.3s;
}

.back-link:hover {
  color: #6366f1;
}

/* Back button */
.back-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  color: #a5b4fc;
  text-decoration: none;
  margin-bottom: 2rem;
  padding: 0.5rem 1rem;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  transition: all 0.3s ease;
  font-size: 0.95rem;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  transform: translateX(-4px);
  color: #fff;
}

.back-icon {
  font-size: 1.2rem;
}

/* Cover Image */
.cover-image-wrapper {
  margin: -2rem -1.5rem 0;
}

.cover-image-container {
  position: relative;
  width: 100%;
  height: 400px;
  overflow: hidden;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.cover-overlay {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 40%;
  background: linear-gradient(to top, #0a0a1a, transparent);
}

/* Article Header */
.article-header {
  margin-bottom: 2rem;
  position: relative;
}

.article-header.no-cover {
  padding-top: 1rem;
}

.article-title {
  font-size: 2.5rem;
  font-weight: 800;
  line-height: 1.2;
  background: linear-gradient(135deg, #e0e7ff, #818cf8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 1rem;
}

.article-meta {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 1.5rem;
  color: #94a3b8;
  font-size: 0.95rem;
}

.meta-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.meta-icon {
  font-size: 1rem;
}

.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.2rem 0.8rem;
  border-radius: 20px;
  background: rgba(99, 102, 241, 0.15);
  color: #a5b4fc;
  font-size: 0.85rem;
  border: 1px solid rgba(99, 102, 241, 0.3);
}

/* Article Body */
.article-body {
  padding: 2.5rem;
  margin-bottom: 1.5rem;
}

.article-content {
  line-height: 1.9;
  color: #cbd5e1;
  font-size: 1.1rem;
}

.article-content :deep(p) {
  margin-bottom: 1rem;
}

.article-content :deep(br) {
  display: block;
  margin: 0.5rem 0;
}

.article-content :deep(strong) {
  color: #e2e8f0;
}

.article-content :deep(em) {
  color: #a5b4fc;
}

/* Action Bar */
.action-bar {
  display: flex;
  gap: 1rem;
  padding: 1.5rem;
  justify-content: center;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 2rem;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  font-size: 1rem;
  font-weight: 500;
  text-decoration: none;
  transition: all 0.3s ease;
}

.edit-btn {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
}

.edit-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(99, 102, 241, 0.4);
}

.delete-btn {
  background: linear-gradient(135deg, #ef4444, #dc2626);
  color: white;
}

.delete-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(239, 68, 68, 0.4);
}

/* Glass card */
.glass-card {
  background: rgba(255, 255, 255, 0.03);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

/* Responsive */
@media (max-width: 768px) {
  .detail-page {
    padding: 1rem;
  }

  .cover-image-container {
    height: 250px;
  }

  .cover-image-wrapper {
    margin: -1rem -1rem 0;
  }

  .article-title {
    font-size: 1.8rem;
  }

  .article-body {
    padding: 1.5rem;
  }

  .article-content {
    font-size: 1rem;
  }

  .action-bar {
    flex-direction: column;
  }

  .action-btn {
    justify-content: center;
  }
}
</style>
