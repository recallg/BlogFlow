<template>
    <div class="create-article-page">
        <nav class="top-nav">
            <button class="nav-back" @click="$router.push('/')">
                <svg
                    viewBox="0 0 24 24"
                    fill="none"
                    stroke="currentColor"
                    stroke-width="2.5"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    class="nav-icon"
                >
                    <path d="M19 12H5M12 19l-7-7 7-7" />
                </svg>
                <span>返回首页</span>
            </button>
            <span class="nav-brand">BlogFlow</span>
            <div class="nav-spacer"></div>
        </nav>
        <div class="page-header">
            <h1 class="page-title">
                <span class="gradient-text">✦ 创作新文章</span>
            </h1>
            <p class="page-subtitle">分享你的想法与见解，让世界听见你的声音</p>
        </div>

        <div class="form-container">
            <form @submit.prevent="handleSubmit" class="article-form">
                <div class="form-group">
                    <label class="form-label">
                        <span class="label-icon">📝</span>
                        标题
                    </label>
                    <input
                        v-model="form.title"
                        type="text"
                        class="form-input"
                        placeholder="输入文章标题..."
                        required
                    />
                    <div class="input-glow"></div>
                </div>

                <div class="form-row">
                    <div class="form-group flex-1">
                        <label class="form-label">
                            <span class="label-icon">✍️</span>
                            作者
                        </label>
                        <input
                            v-model="form.author"
                            type="text"
                            class="form-input"
                            placeholder="作者名称"
                        />
                        <div class="input-glow"></div>
                    </div>

                    <div class="form-group flex-1">
                        <label class="form-label">
                            <span class="label-icon">🏷️</span>
                            标签
                        </label>
                        <input
                            v-model="form.tags"
                            type="text"
                            class="form-input"
                            placeholder="用逗号分隔，如: 技术,前端,Vue"
                        />
                        <div class="input-glow"></div>
                    </div>
                </div>

                <div class="form-group">
                    <label class="form-label">
                        <span class="label-icon">📄</span>
                        摘要
                    </label>
                    <textarea
                        v-model="form.summary"
                        class="form-input form-textarea"
                        placeholder="文章摘要（可选，不填将自动生成）"
                        rows="3"
                    ></textarea>
                    <div class="input-glow"></div>
                </div>

                <div class="form-group">
                    <label class="form-label">
                        <span class="label-icon">📖</span>
                        内容
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
                                placeholder="开始写作吧..."
                                required
                            ></textarea>
                        </div>
                        <div class="preview-pane">
                            <div class="pane-header">
                                <span class="pane-title">预览</span>
                            </div>
                            <div class="preview-scroll">
                                <div
                                    class="markdown-body"
                                    v-html="renderedContent"
                                ></div>
                                <div
                                    v-if="!form.content"
                                    class="preview-placeholder"
                                >
                                    <span>输入 Markdown 内容后实时预览</span>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="editor-toolbar">
                        <span class="toolbar-info">{{ contentLength }} 字</span>
                    </div>
                </div>

                <div class="form-actions">
                    <button
                        type="button"
                        class="btn btn-secondary"
                        @click="$router.push('/')"
                    >
                        ✕ 取消
                    </button>
                    <button
                        type="submit"
                        class="btn btn-primary"
                        :disabled="!form.title || !form.content"
                    >
                        <span class="btn-icon">✦</span>
                        发布文章
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { useRouter } from "vue-router";
import { marked } from "marked";
import { articleApi } from "../api";

const router = useRouter();

const form = ref({
    title: "",
    content: "",
    summary: "",
    author: "",
    tags: "",
    cover_url: "",
});

const contentLength = computed(() => form.value.content.length);

const renderedContent = computed(() => {
    if (!form.value.content) return "";
    return marked.parse(form.value.content, { breaks: true }) as string;
});

const handleSubmit = async () => {
    try {
        const { data } = await articleApi.create({
            title: form.value.title,
            content: form.value.content,
            summary: form.value.summary || undefined,
            author: form.value.author || undefined,
            tags: form.value.tags || undefined,
        });
        router.push(`/article/${data.id}`);
    } catch (error) {
        console.error("Failed to create article:", error);
    }
};
</script>

<style scoped>
.create-article-page {
    max-width: 1100px;
    margin: 0 auto;
    padding: 40px 24px 80px;
}

.page-header {
    text-align: center;
    margin-bottom: 48px;
}

.page-title {
    font-size: 2.5rem;
    font-weight: 800;
    margin-bottom: 12px;
}

.gradient-text {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.page-subtitle {
    color: #94a3b8;
    font-size: 1.1rem;
}

.form-container {
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 24px;
    padding: 40px;
    backdrop-filter: blur(20px);
}

.form-group {
    margin-bottom: 28px;
    position: relative;
}

.form-row {
    display: flex;
    gap: 20px;
}

.flex-1 {
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

.label-icon {
    margin-right: 6px;
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
}

.form-input:focus {
    border-color: #667eea;
    box-shadow:
        0 0 0 3px rgba(102, 126, 234, 0.15),
        0 0 20px rgba(102, 126, 234, 0.1);
}

.form-input::placeholder {
    color: #475569;
}

.form-textarea {
    resize: vertical;
    min-height: 80px;
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
    font-family: "SF Mono", "Fira Code", "Consolas", monospace;
    tab-size: 2;
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

/* Markdown preview styling overrides for dark theme */
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
    font-family: "SF Mono", "Fira Code", "Consolas", monospace;
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

.editor-toolbar {
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

.btn-icon {
    margin-right: 8px;
}

/* Top Navigation */
.top-nav {
    display: flex;
    align-items: center;
    margin-bottom: 32px;
    padding-bottom: 16px;
    border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.nav-back {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    background: rgba(255, 255, 255, 0.06);
    border: 1px solid rgba(255, 255, 255, 0.1);
    color: #94a3b8;
    padding: 10px 20px;
    border-radius: 10px;
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s;
    font-family: inherit;
}

.nav-back:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #e2e8f0;
    transform: translateX(-2px);
}

.nav-icon {
    width: 18px;
    height: 18px;
}

.nav-brand {
    font-size: 1.1rem;
    font-weight: 700;
    background: linear-gradient(135deg, #667eea, #764ba2);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    margin-left: 20px;
}

.nav-spacer {
    flex: 1;
}
</style>
