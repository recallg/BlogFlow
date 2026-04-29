<script setup lang="ts">
import { ref, onMounted, computed } from "vue";
import { useRouter } from "vue-router";
import { articleApi, authApi } from "../api";
import type { Article, User } from "../types";

const router = useRouter();
const articles = ref<Article[]>([]);
const loading = ref(true);
const searchQuery = ref("");
const currentPage = ref(1);
const totalArticles = ref(0);
const pageSize = 9;
const mouseX = ref(0);
const mouseY = ref(0);
const heroRef = ref<HTMLElement | null>(null);
const currentUser = ref<User | null>(null);
const showUserMenu = ref(false);
const showChangePassword = ref(false);
const changePasswordForm = ref({
    old_password: "",
    new_password: "",
    confirm_password: "",
});
const changePasswordError = ref("");
const changePasswordSuccess = ref("");
const changingPassword = ref(false);

const totalPages = computed(() => Math.ceil(totalArticles.value / pageSize));

const isLoggedIn = computed(() => {
    const token = localStorage.getItem("token");
    return !!token;
});

const userAvatar = computed(() => {
    if (currentUser.value) {
        return currentUser.value.username[0].toUpperCase();
    }
    const stored = localStorage.getItem("user");
    if (stored) {
        try {
            const u = JSON.parse(stored);
            return u.username ? u.username[0].toUpperCase() : "?";
        } catch {}
    }
    return "?";
});

function getIconByCategory(_category?: string): string {
    const icons = ["💻", "☕", "🎨", "🔬", "✨", "🚀", "🌟", "🔥"];
    return icons[Math.floor(Math.random() * icons.length)];
}

async function fetchArticles() {
    loading.value = true;
    try {
        const params: any = {
            page: currentPage.value,
            page_size: pageSize,
        };
        if (searchQuery.value.trim()) {
            params.search = searchQuery.value.trim();
        }
        const res = await articleApi.getList(params);
        articles.value = res.data.articles || [];
        totalArticles.value = res.data.total || 0;
    } catch (err) {
        console.error("Failed to fetch articles:", err);
        articles.value = [];
    } finally {
        loading.value = false;
    }
}

function viewArticle(id: number) {
    router.push(`/article/${id}`);
}

function createArticle() {
    router.push("/create");
}

function goToPage(page: number) {
    if (page < 1 || page > totalPages.value) return;
    currentPage.value = page;
    fetchArticles();
    window.scrollTo({ top: 0, behavior: "smooth" });
}

function getTagArray(tags: string): string[] {
    if (!tags) return [];
    return tags
        .split(",")
        .map((t) => t.trim())
        .filter(Boolean);
}

function formatDate(dateStr: string): string {
    const d = new Date(dateStr);
    const now = new Date();
    const diff = now.getTime() - d.getTime();
    const hours = Math.floor(diff / (1000 * 60 * 60));
    if (hours < 1) return "刚刚";
    if (hours < 24) return `${hours} 小时前`;
    const days = Math.floor(hours / 24);
    if (days === 1) return "昨天";
    if (days < 7) return `${days} 天前`;
    if (days < 30) return `${Math.floor(days / 7)} 周前`;
    return d.toLocaleDateString("zh-CN", {
        year: "numeric",
        month: "short",
        day: "numeric",
    });
}

function handleMouseMove(e: MouseEvent) {
    if (!heroRef.value) return;
    const rect = heroRef.value.getBoundingClientRect();
    mouseX.value = ((e.clientX - rect.left) / rect.width - 0.5) * 2;
    mouseY.value = ((e.clientY - rect.top) / rect.height - 0.5) * 2;
}

function getReadingTime(content: string): string {
    const words = content.length;
    const minutes = Math.max(1, Math.ceil(words / 300));
    return `${minutes} 分钟`;
}

async function handleChangePassword() {
    const form = changePasswordForm.value;
    if (!form.old_password || !form.new_password || !form.confirm_password) {
        changePasswordError.value = "请填写所有密码字段";
        return;
    }
    if (form.new_password.length < 6) {
        changePasswordError.value = "新密码至少需要 6 个字符";
        return;
    }
    if (form.new_password !== form.confirm_password) {
        changePasswordError.value = "两次输入的新密码不一致";
        return;
    }
    changePasswordError.value = "";
    changePasswordSuccess.value = "";
    changingPassword.value = true;
    try {
        await authApi.changePassword({
            old_password: form.old_password,
            new_password: form.new_password,
        });
        changePasswordSuccess.value = "密码修改成功";
        changePasswordForm.value = {
            old_password: "",
            new_password: "",
            confirm_password: "",
        };
        setTimeout(() => {
            showChangePassword.value = false;
            changePasswordSuccess.value = "";
        }, 1500);
    } catch (err: any) {
        changePasswordError.value =
            err.response?.data?.error || "密码修改失败，请重试";
    } finally {
        changingPassword.value = false;
    }
}

function goLogin() {
    router.push("/login");
}

function goRegister() {
    router.push("/register");
}

function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    currentUser.value = null;
    showUserMenu.value = false;
    router.push("/");
}

function loadUser() {
    const stored = localStorage.getItem("user");
    if (stored) {
        try {
            currentUser.value = JSON.parse(stored);
        } catch {}
    }
}

onMounted(() => {
    fetchArticles();
    loadUser();
});
</script>

<template>
    <div class="home-root">
        <!-- ===== GLOBAL AMBIENT BACKGROUND ===== -->
        <div class="ambient-bg">
            <div class="orb orb-1"></div>
            <div class="orb orb-2"></div>
            <div class="orb orb-3"></div>
            <div class="orb orb-4"></div>
            <div class="grid-pattern"></div>
        </div>

        <!-- ===== HEADER ===== -->
        <header class="site-header">
            <div class="header-inner">
                <div class="logo-group" @click="router.push('/')">
                    <div class="logo-icon">
                        <svg
                            viewBox="0 0 32 32"
                            fill="none"
                            xmlns="http://www.w3.org/2000/svg"
                        >
                            <rect
                                x="2"
                                y="2"
                                width="28"
                                height="28"
                                rx="8"
                                stroke="url(#logo-grad)"
                                stroke-width="2.5"
                            />
                            <path
                                d="M10 22V12l6 6 6-6v10"
                                stroke="url(#logo-grad)"
                                stroke-width="2.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            />
                            <defs>
                                <linearGradient
                                    id="logo-grad"
                                    x1="2"
                                    y1="2"
                                    x2="30"
                                    y2="30"
                                >
                                    <stop stop-color="#6366f1" />
                                    <stop offset="1" stop-color="#ec4899" />
                                </linearGradient>
                            </defs>
                        </svg>
                    </div>
                    <div class="logo-text">
                        <span class="logo-title">BlogFlow</span>
                        <span class="logo-badge">v2.0</span>
                    </div>
                </div>

                <div class="header-actions">
                    <div class="header-search-wrapper">
                        <div class="header-search">
                            <svg
                                class="search-icon"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <circle cx="10" cy="10" r="7" />
                                <path d="m17 17 5 5" />
                            </svg>
                            <input
                                v-model="searchQuery"
                                type="text"
                                placeholder="搜索文章..."
                                @keyup.enter="fetchArticles"
                                class="search-input"
                            />
                            <button
                                v-if="searchQuery"
                                @click="
                                    searchQuery = '';
                                    fetchArticles();
                                "
                                class="search-clear"
                            >
                                ✕
                            </button>
                        </div>
                    </div>

                    <template v-if="isLoggedIn">
                        <button @click="createArticle" class="btn-create">
                            <svg
                                class="btn-icon"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M12 5v14M5 12h14" />
                            </svg>
                            <span>写文章</span>
                        </button>
                        <div
                            class="user-menu-wrapper"
                            @click="showUserMenu = !showUserMenu"
                        >
                            <div class="user-avatar-btn">
                                <span>{{ userAvatar }}</span>
                            </div>
                            <transition name="dropdown">
                                <div
                                    v-if="showUserMenu"
                                    class="user-dropdown"
                                    @click.stop
                                >
                                    <div class="dropdown-header">
                                        <div class="dropdown-avatar">
                                            {{ userAvatar }}
                                        </div>
                                        <div class="dropdown-info">
                                            <span class="dropdown-name">{{
                                                currentUser?.username || "用户"
                                            }}</span>
                                            <span class="dropdown-email">{{
                                                currentUser?.email || ""
                                            }}</span>
                                        </div>
                                    </div>
                                    <div class="dropdown-divider"></div>
                                    <button
                                        @click="
                                            showChangePassword = true;
                                            showUserMenu = false;
                                        "
                                        class="dropdown-item"
                                    >
                                        <svg
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            class="dropdown-icon"
                                        >
                                            <rect
                                                x="3"
                                                y="11"
                                                width="18"
                                                height="11"
                                                rx="2"
                                                ry="2"
                                            />
                                            <path
                                                d="M7 11V7a5 5 0 0 1 10 0v4"
                                            />
                                        </svg>
                                        修改密码
                                    </button>
                                    <button
                                        @click="handleLogout"
                                        class="dropdown-item logout-item"
                                    >
                                        <svg
                                            viewBox="0 0 24 24"
                                            fill="none"
                                            stroke="currentColor"
                                            stroke-width="2"
                                            stroke-linecap="round"
                                            stroke-linejoin="round"
                                            class="dropdown-icon"
                                        >
                                            <path
                                                d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"
                                            />
                                            <polyline
                                                points="16 17 21 12 16 7"
                                            />
                                            <line
                                                x1="21"
                                                y1="12"
                                                x2="9"
                                                y2="12"
                                            />
                                        </svg>
                                        退出登录
                                    </button>
                                </div>
                            </transition>
                        </div>
                    </template>
                    <template v-else>
                        <button @click="goLogin" class="btn-login">登录</button>
                        <button @click="goRegister" class="btn-create">
                            <svg
                                class="btn-icon"
                                viewBox="0 0 24 24"
                                fill="none"
                                stroke="currentColor"
                                stroke-width="2.5"
                                stroke-linecap="round"
                                stroke-linejoin="round"
                            >
                                <path d="M12 5v14M5 12h14" />
                            </svg>
                            <span>注册</span>
                        </button>
                    </template>
                </div>
            </div>
        </header>

        <!-- ===== HERO SECTION ===== -->
        <section ref="heroRef" class="hero" @mousemove="handleMouseMove">
            <div class="hero-particles">
                <div class="particle p1"></div>
                <div class="particle p2"></div>
                <div class="particle p3"></div>
                <div class="particle p4"></div>
                <div class="particle p5"></div>
            </div>

            <div class="hero-content">
                <div class="hero-badge">
                    <span class="badge-dot"></span>
                    发现精彩内容
                    <svg
                        class="badge-arrow"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path d="M5 12h14M12 5l7 7-7 7" />
                    </svg>
                </div>

                <h1 class="hero-title">
                    用文字
                    <span class="gradient-word">记录世界</span>
                </h1>

                <p class="hero-desc">
                    每一个故事都值得被讲述，每一个想法都能照亮他人。
                    <br class="hide-mobile" />
                    在这里，开启你的创作之旅。
                </p>

                <div class="hero-stats">
                    <div class="stat-item">
                        <span class="stat-value">{{ totalArticles }}</span>
                        <span class="stat-label">文章</span>
                    </div>
                    <div class="stat-divider"></div>
                    <div class="stat-item">
                        <span class="stat-value">{{
                            Math.max(1, articles.length)
                        }}</span>
                        <span class="stat-label">作者</span>
                    </div>
                    <div class="stat-divider"></div>
                    <div class="stat-item">
                        <span class="stat-value">{{ totalPages }}</span>
                        <span class="stat-label">页数</span>
                    </div>
                </div>
            </div>

            <div
                class="hero-glow"
                :style="{
                    transform: `translate(calc(-50% + ${mouseX * 30}px), calc(-50% + ${mouseY * 30}px))`,
                }"
            ></div>
        </section>

        <!-- ===== CONTENT SECTION ===== -->
        <section class="content-section">
            <div class="content-header">
                <div class="content-header-left">
                    <h2 class="section-title">
                        <span class="title-accent"></span>
                        最新文章
                    </h2>
                    <p class="section-sub">发现最新鲜的思想与见解</p>
                </div>
                <div class="content-header-right">
                    <span class="articles-count">{{ totalArticles }} 篇</span>
                </div>
            </div>

            <!-- Loading -->
            <div v-if="loading" class="skeleton-grid">
                <div v-for="n in 6" :key="n" class="skeleton-card">
                    <div class="skeleton-shimmer"></div>
                </div>
            </div>

            <!-- Empty -->
            <div v-else-if="articles.length === 0" class="empty-state">
                <div class="empty-illustration">
                    <div class="empty-orb"></div>
                    <span class="empty-emoji">✍️</span>
                </div>
                <h3 class="empty-title">还没有文章</h3>
                <p class="empty-desc">你准备好分享第一个故事了吗？</p>
            </div>

            <!-- Article Grid -->
            <div v-else class="article-grid">
                <article
                    v-for="(article, index) in articles"
                    :key="article.id"
                    @click="viewArticle(article.id)"
                    class="article-card"
                    :style="{ '--delay': `${index * 0.05}s` }"
                >
                    <div class="card-glow"></div>
                    <div class="card-bg"></div>

                    <div class="card-content">
                        <!-- Card Top -->
                        <div class="card-top">
                            <div class="card-icon-row">
                                <span class="card-icon">{{
                                    getIconByCategory(article.tags)
                                }}</span>
                                <span class="card-reading">{{
                                    getReadingTime(article.content)
                                }}</span>
                            </div>
                            <h3 class="card-title">{{ article.title }}</h3>
                            <p class="card-summary">
                                {{
                                    article.summary ||
                                    article.content.slice(0, 120) + "..."
                                }}
                            </p>
                        </div>

                        <!-- Card Tags -->
                        <div
                            v-if="getTagArray(article.tags).length > 0"
                            class="card-tags"
                        >
                            <span
                                v-for="tag in getTagArray(article.tags).slice(
                                    0,
                                    3,
                                )"
                                :key="tag"
                                class="card-tag"
                            >
                                #{{ tag }}
                            </span>
                            <span
                                v-if="getTagArray(article.tags).length > 3"
                                class="card-tag-more"
                            >
                                +{{ getTagArray(article.tags).length - 3 }}
                            </span>
                        </div>

                        <!-- Card Footer -->
                        <div class="card-footer">
                            <div class="card-author">
                                <div class="author-avatar">
                                    {{
                                        article.author
                                            ? article.author[0].toUpperCase()
                                            : "A"
                                    }}
                                </div>
                                <div class="author-info">
                                    <span class="author-name">{{
                                        article.author || "Anonymous"
                                    }}</span>
                                    <span class="author-date">{{
                                        formatDate(article.created_at)
                                    }}</span>
                                </div>
                            </div>
                            <div class="card-arrow">
                                <svg
                                    viewBox="0 0 24 24"
                                    fill="none"
                                    stroke="currentColor"
                                    stroke-width="2"
                                    stroke-linecap="round"
                                    stroke-linejoin="round"
                                >
                                    <path d="M5 12h14M12 5l7 7-7 7" />
                                </svg>
                            </div>
                        </div>
                    </div>
                </article>
            </div>

            <!-- Pagination -->
            <div v-if="totalPages > 1" class="pagination">
                <button
                    @click="goToPage(currentPage - 1)"
                    :disabled="currentPage === 1"
                    class="page-btn prev"
                >
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path d="M19 12H5M12 19l-7-7 7-7" />
                    </svg>
                    上一页
                </button>

                <div class="page-numbers">
                    <template v-for="p in totalPages" :key="p">
                        <button
                            v-if="
                                p === 1 ||
                                p === totalPages ||
                                Math.abs(p - currentPage) <= 1
                            "
                            @click="goToPage(p)"
                            :class="['page-num', { active: p === currentPage }]"
                        >
                            {{ p }}
                        </button>
                        <span
                            v-else-if="
                                p === currentPage - 2 || p === currentPage + 2
                            "
                            class="page-dots"
                            >...</span
                        >
                    </template>
                </div>

                <button
                    @click="goToPage(currentPage + 1)"
                    :disabled="currentPage === totalPages"
                    class="page-btn next"
                >
                    下一页
                    <svg
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2.5"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                    >
                        <path d="M5 12h14M12 5l7 7-7 7" />
                    </svg>
                </button>
            </div>
        </section>

        <!-- ===== CHANGE PASSWORD MODAL ===== -->
        <teleport to="body">
            <div
                v-if="showChangePassword"
                class="modal-overlay"
                @click.self="
                    showChangePassword = false;
                    changePasswordError = '';
                    changePasswordSuccess = '';
                "
            >
                <div class="modal-dialog">
                    <div class="modal-header">
                        <h3 class="modal-title">🔒 修改密码</h3>
                        <button
                            class="modal-close"
                            @click="
                                showChangePassword = false;
                                changePasswordError = '';
                                changePasswordSuccess = '';
                            "
                        >
                            ✕
                        </button>
                    </div>
                    <div class="modal-body">
                        <div v-if="changePasswordSuccess" class="modal-success">
                            <span>✅</span> {{ changePasswordSuccess }}
                        </div>
                        <div v-if="changePasswordError" class="modal-error">
                            <span>⚠️</span> {{ changePasswordError }}
                        </div>
                        <div class="modal-form-group">
                            <label class="modal-label">原密码</label>
                            <input
                                v-model="changePasswordForm.old_password"
                                type="password"
                                class="modal-input"
                                placeholder="输入当前密码"
                            />
                        </div>
                        <div class="modal-form-group">
                            <label class="modal-label">新密码</label>
                            <input
                                v-model="changePasswordForm.new_password"
                                type="password"
                                class="modal-input"
                                placeholder="至少 6 个字符"
                            />
                        </div>
                        <div class="modal-form-group">
                            <label class="modal-label">确认新密码</label>
                            <input
                                v-model="changePasswordForm.confirm_password"
                                type="password"
                                class="modal-input"
                                placeholder="再次输入新密码"
                            />
                        </div>
                    </div>
                    <div class="modal-footer">
                        <button
                            class="modal-btn modal-btn-secondary"
                            @click="
                                showChangePassword = false;
                                changePasswordError = '';
                                changePasswordSuccess = '';
                            "
                        >
                            取消
                        </button>
                        <button
                            class="modal-btn modal-btn-primary"
                            :disabled="changingPassword"
                            @click="handleChangePassword"
                        >
                            {{ changingPassword ? "修改中..." : "确认修改" }}
                        </button>
                    </div>
                </div>
            </div>
        </teleport>

        <!-- ===== FOOTER ===== -->
        <footer class="site-footer">
            <div class="footer-glow"></div>
            <div class="footer-inner">
                <div class="footer-brand">
                    <div class="footer-logo">B</div>
                    <span class="footer-name">BlogFlow</span>
                </div>
                <p class="footer-text">
                    Built with <span class="footer-heart">❤️</span> using Vue 3
                    + Go + SQLite
                </p>
                <div class="footer-tech">
                    <span class="tech-dot" style="--dot-color: #42b883"></span>
                    <span>Vue 3</span>
                    <span class="tech-sep">·</span>
                    <span class="tech-dot" style="--dot-color: #00add8"></span>
                    <span>Go</span>
                    <span class="tech-sep">·</span>
                    <span class="tech-dot" style="--dot-color: #6366f1"></span>
                    <span>SQLite</span>
                </div>
            </div>
        </footer>
    </div>
</template>

<style scoped>
/* ============================================
   ROOT LAYOUT
   ============================================ */
.home-root {
    position: relative;
    min-height: 100vh;
    background: #05050a;
    color: #e2e8f0;
    overflow-x: hidden;
}

/* ============================================
   AMBIENT BACKGROUND
   ============================================ */
.ambient-bg {
    position: fixed;
    inset: 0;
    z-index: 0;
    pointer-events: none;
    overflow: hidden;
}

.orb {
    position: absolute;
    border-radius: 50%;
    filter: blur(80px);
    opacity: 0.5;
}

.orb-1 {
    width: 500px;
    height: 500px;
    top: -200px;
    right: -100px;
    background: #6366f1;
    animation: orbFloat1 8s ease-in-out infinite;
}

.orb-2 {
    width: 400px;
    height: 400px;
    bottom: -150px;
    left: -150px;
    background: #ec4899;
    animation: orbFloat2 10s ease-in-out infinite;
}

.orb-3 {
    width: 300px;
    height: 300px;
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    background: #6366f1;
    opacity: 0.15;
    animation: orbFloat3 12s ease-in-out infinite;
}

.orb-4 {
    width: 200px;
    height: 200px;
    top: 60%;
    right: 10%;
    background: #06b6d4;
    animation: orbFloat4 7s ease-in-out infinite;
}

.grid-pattern {
    position: absolute;
    inset: 0;
    background-image:
        linear-gradient(rgba(99, 102, 241, 0.03) 1px, transparent 1px),
        linear-gradient(90deg, rgba(99, 102, 241, 0.03) 1px, transparent 1px);
    background-size: 60px 60px;
}

@keyframes orbFloat1 {
    0%,
    100% {
        transform: translate(0, 0) scale(1);
    }
    33% {
        transform: translate(30px, -30px) scale(1.1);
    }
    66% {
        transform: translate(-20px, 20px) scale(0.9);
    }
}

@keyframes orbFloat2 {
    0%,
    100% {
        transform: translate(0, 0) scale(1);
    }
    33% {
        transform: translate(-30px, 30px) scale(1.05);
    }
    66% {
        transform: translate(20px, -20px) scale(0.95);
    }
}

@keyframes orbFloat3 {
    0%,
    100% {
        opacity: 0.15;
        transform: translate(-50%, -50%) scale(1);
    }
    50% {
        opacity: 0.25;
        transform: translate(-50%, -50%) scale(1.2);
    }
}

@keyframes orbFloat4 {
    0%,
    100% {
        transform: translate(0, 0);
    }
    50% {
        transform: translate(-30px, -30px);
    }
}

/* ============================================
   HEADER
   ============================================ */
.site-header {
    position: sticky;
    top: 0;
    z-index: 100;
    background: rgba(5, 5, 10, 0.8);
    backdrop-filter: blur(24px) saturate(1.4);
    border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.header-inner {
    max-width: 1280px;
    margin: 0 auto;
    padding: 0 32px;
    height: 72px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 24px;
}

/* Logo */
.logo-group {
    display: flex;
    align-items: center;
    gap: 12px;
    cursor: pointer;
    flex-shrink: 0;
}

.logo-icon {
    width: 36px;
    height: 36px;
}

.logo-icon svg {
    width: 100%;
    height: 100%;
}

.logo-text {
    display: flex;
    align-items: baseline;
    gap: 8px;
}

.logo-title {
    font-size: 1.3rem;
    font-weight: 800;
    letter-spacing: -0.5px;
    background: linear-gradient(135deg, #e0e7ff, #818cf8);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
}

.logo-badge {
    font-size: 0.65rem;
    font-weight: 600;
    padding: 2px 6px;
    border-radius: 4px;
    background: rgba(99, 102, 241, 0.2);
    color: #818cf8;
    border: 1px solid rgba(99, 102, 241, 0.3);
}

/* Header Actions */
.header-actions {
    display: flex;
    align-items: center;
    gap: 16px;
    flex: 1;
    justify-content: flex-end;
}

/* Header Search */
.header-search-wrapper {
    max-width: 320px;
    width: 100%;
}

.header-search {
    position: relative;
    display: flex;
    align-items: center;
}

.search-icon {
    position: absolute;
    left: 14px;
    width: 18px;
    height: 18px;
    color: #475569;
    pointer-events: none;
}

.search-input {
    width: 100%;
    padding: 10px 36px 10px 42px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    color: #e2e8f0;
    font-size: 0.9rem;
    outline: none;
    transition: all 0.3s ease;
}

.search-input::placeholder {
    color: #475569;
}

.search-input:focus {
    border-color: rgba(99, 102, 241, 0.4);
    background: rgba(255, 255, 255, 0.06);
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.search-clear {
    position: absolute;
    right: 8px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: transparent;
    color: #64748b;
    border-radius: 6px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.8rem;
}

.search-clear:hover {
    color: #e2e8f0;
    background: rgba(255, 255, 255, 0.05);
}

/* Create Button */
.btn-create {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: white;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    white-space: nowrap;
    font-family: inherit;
}

.btn-create .btn-icon {
    width: 18px;
    height: 18px;
}

.btn-create:hover {
    transform: translateY(-1px);
    box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

.btn-create:active {
    transform: translateY(0);
}

/* Login Button */
.btn-login {
    padding: 10px 20px;
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 12px;
    background: transparent;
    color: #94a3b8;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: inherit;
    white-space: nowrap;
}

.btn-login:hover {
    border-color: rgba(99, 102, 241, 0.3);
    color: #e2e8f0;
    background: rgba(255, 255, 255, 0.03);
}

/* User Menu */
.user-menu-wrapper {
    position: relative;
    cursor: pointer;
}

.user-avatar-btn {
    width: 38px;
    height: 38px;
    border-radius: 12px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.85rem;
    font-weight: 700;
    color: white;
    transition: all 0.3s ease;
    border: 2px solid transparent;
}

.user-avatar-btn:hover {
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.2);
    transform: scale(1.05);
}

/* Dropdown */
.user-dropdown {
    position: absolute;
    top: calc(100% + 10px);
    right: 0;
    width: 240px;
    background: rgba(12, 12, 30, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 16px;
    backdrop-filter: blur(24px);
    overflow: hidden;
    z-index: 200;
    box-shadow: 0 20px 60px rgba(0, 0, 0, 0.5);
    animation: dropdownIn 0.2s ease-out;
}

@keyframes dropdownIn {
    from {
        opacity: 0;
        transform: translateY(-8px) scale(0.96);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

/* Dropdown transition */
.dropdown-enter-active {
    animation: dropdownIn 0.2s ease-out;
}

.dropdown-leave-active {
    animation: dropdownIn 0.15s ease-in reverse;
}

.dropdown-header {
    padding: 16px;
    display: flex;
    align-items: center;
    gap: 12px;
}

.dropdown-avatar {
    width: 40px;
    height: 40px;
    border-radius: 12px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 1rem;
    font-weight: 700;
    color: white;
    flex-shrink: 0;
}

.dropdown-info {
    display: flex;
    flex-direction: column;
    gap: 2px;
    overflow: hidden;
}

.dropdown-name {
    font-size: 0.95rem;
    font-weight: 600;
    color: #f1f5f9;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.dropdown-email {
    font-size: 0.78rem;
    color: #64748b;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.dropdown-divider {
    height: 1px;
    background: rgba(255, 255, 255, 0.05);
    margin: 0 8px;
}

.dropdown-item {
    display: flex;
    align-items: center;
    gap: 10px;
    width: 100%;
    padding: 12px 16px;
    border: none;
    background: transparent;
    color: #94a3b8;
    font-size: 0.88rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.2s ease;
    font-family: inherit;
}

.dropdown-item:hover {
    background: rgba(255, 255, 255, 0.03);
    color: #e2e8f0;
}

.logout-item:hover {
    color: #fca5a5;
}

.dropdown-icon {
    width: 18px;
    height: 18px;
}

/* ============================================
   HERO
   ============================================ */
.hero {
    position: relative;
    z-index: 1;
    max-width: 1280px;
    margin: 0 auto;
    padding: 80px 32px 100px;
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    overflow: hidden;
}

.hero-particles {
    position: absolute;
    inset: 0;
    pointer-events: none;
}

.particle {
    position: absolute;
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: #6366f1;
    opacity: 0.3;
}

.p1 {
    top: 20%;
    left: 10%;
    animation: particle1 6s ease-in-out infinite;
}
.p2 {
    top: 60%;
    left: 85%;
    width: 4px;
    height: 4px;
    background: #ec4899;
    animation: particle2 8s ease-in-out infinite;
}
.p3 {
    top: 30%;
    left: 75%;
    width: 3px;
    height: 3px;
    background: #06b6d4;
    animation: particle3 7s ease-in-out infinite;
}
.p4 {
    top: 70%;
    left: 20%;
    width: 5px;
    height: 5px;
    background: #f59e0b;
    animation: particle1 9s ease-in-out infinite;
}
.p5 {
    top: 15%;
    left: 50%;
    width: 4px;
    height: 4px;
    background: #818cf8;
    animation: particle2 5s ease-in-out infinite;
}

@keyframes particle1 {
    0%,
    100% {
        transform: translate(0, 0) scale(1);
        opacity: 0.3;
    }
    33% {
        transform: translate(20px, -15px) scale(1.5);
        opacity: 0.5;
    }
    66% {
        transform: translate(-10px, 10px) scale(0.8);
        opacity: 0.2;
    }
}

@keyframes particle2 {
    0%,
    100% {
        transform: translate(0, 0) scale(1);
        opacity: 0.3;
    }
    50% {
        transform: translate(-20px, 20px) scale(1.3);
        opacity: 0.6;
    }
}

@keyframes particle3 {
    0%,
    100% {
        transform: translate(0, 0);
        opacity: 0.3;
    }
    50% {
        transform: translate(15px, -25px) scale(1.2);
        opacity: 0.5;
    }
}

.hero-glow {
    position: absolute;
    width: 400px;
    height: 400px;
    border-radius: 50%;
    background: radial-gradient(
        circle,
        rgba(99, 102, 241, 0.08),
        transparent 70%
    );
    pointer-events: none;
    transition: transform 0.3s ease-out;
    top: 50%;
    left: 50%;
}

.hero-content {
    position: relative;
    z-index: 2;
}

.hero-badge {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 6px 16px 6px 10px;
    border-radius: 20px;
    background: linear-gradient(
        135deg,
        rgba(99, 102, 241, 0.1),
        rgba(236, 72, 153, 0.1)
    );
    border: 1px solid rgba(99, 102, 241, 0.15);
    color: #a5b4fc;
    font-size: 0.85rem;
    font-weight: 500;
    margin-bottom: 28px;
    animation: fadeInUp 0.6s ease-out;
}

.badge-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #6366f1;
    animation: pulse-dot 2s ease-in-out infinite;
}

@keyframes pulse-dot {
    0%,
    100% {
        opacity: 1;
        transform: scale(1);
    }
    50% {
        opacity: 0.5;
        transform: scale(0.8);
    }
}

.badge-arrow {
    width: 14px;
    height: 14px;
}

.hero-title {
    font-size: clamp(2.5rem, 6vw, 4.5rem);
    font-weight: 800;
    line-height: 1.15;
    letter-spacing: -2px;
    color: #f1f5f9;
    margin-bottom: 20px;
    animation: fadeInUp 0.6s ease-out 0.1s both;
}

.gradient-word {
    display: inline-block;
    background: linear-gradient(135deg, #6366f1, #a78bfa, #c084fc, #ec4899);
    background-size: 300% 300%;
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
    background-clip: text;
    animation: gradientShift 4s ease-in-out infinite;
}

@keyframes gradientShift {
    0%,
    100% {
        background-position: 0% 50%;
    }
    50% {
        background-position: 100% 50%;
    }
}

.hero-desc {
    font-size: 1.15rem;
    color: #64748b;
    max-width: 600px;
    margin: 0 auto 40px;
    line-height: 1.7;
    animation: fadeInUp 0.6s ease-out 0.2s both;
}

.hide-mobile {
    display: block;
}

/* Hero Stats */
.hero-stats {
    display: inline-flex;
    align-items: center;
    gap: 24px;
    padding: 16px 32px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.06);
    backdrop-filter: blur(12px);
    animation: fadeInUp 0.6s ease-out 0.3s both;
}

.stat-item {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 2px;
}

.stat-value {
    font-size: 1.6rem;
    font-weight: 800;
    color: #f1f5f9;
    line-height: 1;
}

.stat-label {
    font-size: 0.75rem;
    color: #64748b;
    text-transform: uppercase;
    letter-spacing: 1px;
}

.stat-divider {
    width: 1px;
    height: 32px;
    background: rgba(255, 255, 255, 0.08);
}

@keyframes fadeInUp {
    from {
        opacity: 0;
        transform: translateY(20px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}

/* ============================================
   CONTENT SECTION
   ============================================ */
.content-section {
    position: relative;
    z-index: 1;
    max-width: 1280px;
    margin: 0 auto;
    padding: 0 32px 80px;
}

.content-header {
    display: flex;
    align-items: flex-end;
    justify-content: space-between;
    margin-bottom: 32px;
}

.content-header-left {
    display: flex;
    flex-direction: column;
    gap: 8px;
}

.section-title {
    display: flex;
    align-items: center;
    gap: 12px;
    font-size: 1.6rem;
    font-weight: 700;
    color: #f1f5f9;
}

.title-accent {
    width: 4px;
    height: 24px;
    border-radius: 2px;
    background: linear-gradient(180deg, #6366f1, #ec4899);
}

.section-sub {
    font-size: 0.9rem;
    color: #64748b;
    margin-left: 16px;
}

.content-header-right {
    flex-shrink: 0;
}

.articles-count {
    padding: 6px 14px;
    border-radius: 8px;
    background: rgba(99, 102, 241, 0.1);
    color: #818cf8;
    font-size: 0.85rem;
    font-weight: 600;
}

/* ============================================
   SKELETON LOADING
   ============================================ */
.skeleton-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 24px;
}

.skeleton-card {
    height: 340px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.03);
    border: 1px solid rgba(255, 255, 255, 0.05);
    overflow: hidden;
    position: relative;
}

.skeleton-shimmer {
    position: absolute;
    inset: 0;
    background: linear-gradient(
        90deg,
        transparent,
        rgba(255, 255, 255, 0.03),
        transparent
    );
    background-size: 200% 100%;
    animation: shimmer 1.5s ease-in-out infinite;
}

@keyframes shimmer {
    0% {
        background-position: -200% 0;
    }
    100% {
        background-position: 200% 0;
    }
}

/* ============================================
   EMPTY STATE
   ============================================ */
.empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    padding: 80px 20px;
    text-align: center;
}

.empty-illustration {
    position: relative;
    width: 120px;
    height: 120px;
    margin-bottom: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
}

.empty-orb {
    position: absolute;
    inset: 0;
    border-radius: 50%;
    background: radial-gradient(
        circle,
        rgba(99, 102, 241, 0.15),
        transparent 70%
    );
    animation: pulse-orb 3s ease-in-out infinite;
}

@keyframes pulse-orb {
    0%,
    100% {
        transform: scale(1);
        opacity: 0.5;
    }
    50% {
        transform: scale(1.2);
        opacity: 1;
    }
}

.empty-emoji {
    font-size: 3rem;
    position: relative;
    z-index: 1;
}

.empty-title {
    font-size: 1.5rem;
    font-weight: 700;
    color: #f1f5f9;
    margin-bottom: 8px;
}

.empty-desc {
    color: #64748b;
    margin-bottom: 24px;
}

.btn-empty {
    display: inline-flex;
    align-items: center;
    gap: 8px;
    padding: 12px 24px;
    border: none;
    border-radius: 12px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: white;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: inherit;
}

.btn-empty:hover {
    transform: translateY(-2px);
    box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

/* ============================================
   ARTICLE GRID
   ============================================ */
.article-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 24px;
}

.article-card {
    position: relative;
    cursor: pointer;
    border-radius: 16px;
    overflow: hidden;
    animation: cardIn 0.5s ease-out var(--delay, 0s) both;
}

@keyframes cardIn {
    from {
        opacity: 0;
        transform: translateY(30px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.card-glow {
    position: absolute;
    inset: -2px;
    border-radius: 18px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6, #ec4899, #6366f1);
    background-size: 300% 300%;
    opacity: 0;
    transition: opacity 0.5s ease;
    animation: glowShift 3s ease-in-out infinite;
    z-index: 0;
}

.article-card:hover .card-glow {
    opacity: 0.5;
}

@keyframes glowShift {
    0%,
    100% {
        background-position: 0% 50%;
    }
    50% {
        background-position: 100% 50%;
    }
}

.card-bg {
    position: absolute;
    inset: 2px;
    border-radius: 16px;
    background: rgba(255, 255, 255, 0.02);
    border: 1px solid rgba(255, 255, 255, 0.06);
    transition: all 0.4s ease;
    z-index: 0;
}

.article-card:hover .card-bg {
    background: rgba(255, 255, 255, 0.04);
    border-color: rgba(99, 102, 241, 0.2);
}

.card-content {
    position: relative;
    z-index: 1;
    padding: 28px;
    display: flex;
    flex-direction: column;
    height: 100%;
    min-height: 320px;
}

.card-top {
    flex: 1;
}

.card-icon-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 16px;
}

.card-icon {
    font-size: 2rem;
    line-height: 1;
}

.card-reading {
    font-size: 0.75rem;
    color: #64748b;
    padding: 4px 10px;
    border-radius: 6px;
    background: rgba(255, 255, 255, 0.04);
    border: 1px solid rgba(255, 255, 255, 0.05);
}

.card-title {
    font-size: 1.2rem;
    font-weight: 700;
    color: #f1f5f9;
    line-height: 1.4;
    margin-bottom: 12px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
    transition: color 0.3s ease;
}

.article-card:hover .card-title {
    color: #a5b4fc;
}

.card-summary {
    font-size: 0.88rem;
    color: #64748b;
    line-height: 1.7;
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
}

/* Card Tags */
.card-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-top: 16px;
}

.card-tag {
    padding: 4px 12px;
    font-size: 0.78rem;
    border-radius: 8px;
    background: rgba(99, 102, 241, 0.08);
    color: #818cf8;
    border: 1px solid rgba(99, 102, 241, 0.12);
    transition: all 0.3s ease;
}

.article-card:hover .card-tag {
    background: rgba(99, 102, 241, 0.12);
    border-color: rgba(99, 102, 241, 0.2);
}

.card-tag-more {
    padding: 4px 12px;
    font-size: 0.78rem;
    border-radius: 8px;
    background: rgba(255, 255, 255, 0.03);
    color: #475569;
}

/* Card Footer */
.card-footer {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding-top: 20px;
    margin-top: 20px;
    border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.card-author {
    display: flex;
    align-items: center;
    gap: 10px;
}

.author-avatar {
    width: 32px;
    height: 32px;
    border-radius: 10px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.75rem;
    font-weight: 700;
    color: white;
    flex-shrink: 0;
}

.author-info {
    display: flex;
    flex-direction: column;
    gap: 1px;
}

.author-name {
    font-size: 0.85rem;
    font-weight: 600;
    color: #cbd5e1;
}

.author-date {
    font-size: 0.75rem;
    color: #475569;
}

.card-arrow {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 10px;
    background: rgba(255, 255, 255, 0.03);
    color: #475569;
    transition: all 0.3s ease;
}

.card-arrow svg {
    width: 16px;
    height: 16px;
}

.article-card:hover .card-arrow {
    background: rgba(99, 102, 241, 0.12);
    color: #818cf8;
    transform: translateX(4px);
}

/* ============================================
   PAGINATION
   ============================================ */
.pagination {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 8px;
    margin-top: 56px;
    padding-top: 40px;
    border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.page-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 20px;
    border: 1px solid rgba(255, 255, 255, 0.08);
    border-radius: 12px;
    background: rgba(255, 255, 255, 0.02);
    color: #94a3b8;
    font-size: 0.88rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: inherit;
}

.page-btn svg {
    width: 16px;
    height: 16px;
}

.page-btn:hover:not(:disabled) {
    background: rgba(255, 255, 255, 0.05);
    border-color: rgba(99, 102, 241, 0.2);
    color: #e2e8f0;
}

.page-btn:disabled {
    opacity: 0.3;
    cursor: not-allowed;
}

.page-numbers {
    display: flex;
    align-items: center;
    gap: 4px;
}

.page-num {
    width: 38px;
    height: 38px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    border-radius: 10px;
    background: transparent;
    color: #64748b;
    font-size: 0.88rem;
    font-weight: 500;
    cursor: pointer;
    transition: all 0.3s ease;
    font-family: inherit;
}

.page-num:hover {
    background: rgba(255, 255, 255, 0.05);
    color: #e2e8f0;
}

.page-num.active {
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: white;
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.25);
}

.page-dots {
    color: #475569;
    padding: 0 4px;
    font-size: 0.88rem;
    letter-spacing: 2px;
}

/* ============================================
   CHANGE PASSWORD MODAL
   ============================================ */
.modal-overlay {
    position: fixed;
    inset: 0;
    z-index: 9999;
    background: rgba(0, 0, 0, 0.6);
    backdrop-filter: blur(8px);
    display: flex;
    align-items: center;
    justify-content: center;
    animation: fadeIn 0.2s ease;
}

@keyframes fadeIn {
    from {
        opacity: 0;
    }
    to {
        opacity: 1;
    }
}

.modal-dialog {
    width: 420px;
    max-width: 90vw;
    background: rgba(15, 15, 35, 0.98);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 20px;
    backdrop-filter: blur(24px);
    box-shadow: 0 30px 80px rgba(0, 0, 0, 0.6);
    animation: modalIn 0.25s ease-out;
}

@keyframes modalIn {
    from {
        opacity: 0;
        transform: translateY(20px) scale(0.95);
    }
    to {
        opacity: 1;
        transform: translateY(0) scale(1);
    }
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px 0;
}

.modal-title {
    font-size: 1.15rem;
    font-weight: 700;
    color: #f1f5f9;
}

.modal-close {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: none;
    background: rgba(255, 255, 255, 0.05);
    color: #64748b;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.2s;
    font-size: 0.9rem;
}

.modal-close:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #e2e8f0;
}

.modal-body {
    padding: 20px 24px 24px;
}

.modal-success {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    border-radius: 10px;
    background: rgba(34, 197, 94, 0.1);
    border: 1px solid rgba(34, 197, 94, 0.2);
    color: #86efac;
    font-size: 0.9rem;
    margin-bottom: 16px;
}

.modal-error {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 14px;
    border-radius: 10px;
    background: rgba(239, 68, 68, 0.1);
    border: 1px solid rgba(239, 68, 68, 0.2);
    color: #fca5a5;
    font-size: 0.9rem;
    margin-bottom: 16px;
}

.modal-form-group {
    margin-bottom: 18px;
}

.modal-form-group:last-child {
    margin-bottom: 0;
}

.modal-label {
    display: block;
    font-size: 0.85rem;
    font-weight: 600;
    color: #cbd5e1;
    margin-bottom: 8px;
}

.modal-input {
    width: 100%;
    padding: 12px 16px;
    background: rgba(0, 0, 0, 0.3);
    border: 1px solid rgba(255, 255, 255, 0.1);
    border-radius: 10px;
    color: #f1f5f9;
    font-size: 0.95rem;
    outline: none;
    transition: all 0.3s;
    font-family: inherit;
    box-sizing: border-box;
}

.modal-input:focus {
    border-color: #6366f1;
    box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.15);
}

.modal-input::placeholder {
    color: #475569;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 0 24px 20px;
}

.modal-btn {
    padding: 10px 24px;
    border: none;
    border-radius: 10px;
    font-size: 0.9rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s;
    font-family: inherit;
}

.modal-btn-primary {
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    color: white;
    box-shadow: 0 4px 12px rgba(99, 102, 241, 0.3);
}

.modal-btn-primary:hover:not(:disabled) {
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(99, 102, 241, 0.4);
}

.modal-btn-primary:disabled {
    opacity: 0.5;
    cursor: not-allowed;
}

.modal-btn-secondary {
    background: rgba(255, 255, 255, 0.06);
    color: #94a3b8;
    border: 1px solid rgba(255, 255, 255, 0.1);
}

.modal-btn-secondary:hover {
    background: rgba(255, 255, 255, 0.1);
    color: #e2e8f0;
}

/* ============================================
   FOOTER
   ============================================ */
.site-footer {
    position: relative;
    z-index: 1;
    margin-top: 40px;
    padding: 40px 32px;
    border-top: 1px solid rgba(255, 255, 255, 0.04);
    overflow: hidden;
}

.footer-glow {
    position: absolute;
    top: 0;
    left: 50%;
    transform: translateX(-50%);
    width: 600px;
    height: 1px;
    background: linear-gradient(
        90deg,
        transparent,
        rgba(99, 102, 241, 0.3),
        transparent
    );
}

.footer-inner {
    max-width: 1280px;
    margin: 0 auto;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 16px;
    text-align: center;
}

.footer-brand {
    display: flex;
    align-items: center;
    gap: 10px;
}

.footer-logo {
    width: 32px;
    height: 32px;
    border-radius: 10px;
    background: linear-gradient(135deg, #6366f1, #8b5cf6);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 0.9rem;
    font-weight: 800;
    color: white;
}

.footer-name {
    font-size: 1.1rem;
    font-weight: 700;
    color: #cbd5e1;
}

.footer-text {
    color: #475569;
    font-size: 0.85rem;
}

.footer-heart {
    color: #ec4899;
    display: inline-block;
    animation: heartBeat 1.5s ease-in-out infinite;
}

@keyframes heartBeat {
    0%,
    100% {
        transform: scale(1);
    }
    50% {
        transform: scale(1.2);
    }
}

.footer-tech {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 0.8rem;
    color: #475569;
}

.tech-dot {
    width: 6px;
    height: 6px;
    border-radius: 50%;
    background: var(--dot-color);
    opacity: 0.6;
}

.tech-sep {
    color: rgba(255, 255, 255, 0.08);
}

/* ============================================
   RESPONSIVE
   ============================================ */
@media (max-width: 768px) {
    .header-inner {
        padding: 0 16px;
        height: 64px;
    }

    .header-search-wrapper {
        max-width: 180px;
    }

    .hide-mobile {
        display: none;
    }

    .hero {
        padding: 48px 16px 60px;
    }

    .hero-stats {
        gap: 16px;
        padding: 12px 20px;
    }

    .content-section {
        padding: 0 16px 60px;
    }

    .article-grid {
        grid-template-columns: 1fr;
    }

    .content-header {
        flex-direction: column;
        gap: 12px;
        align-items: flex-start;
    }

    .skeleton-grid {
        grid-template-columns: 1fr;
    }

    .pagination {
        flex-wrap: wrap;
        gap: 6px;
    }

    .page-btn span {
        display: none;
    }

    .btn-create span {
        display: none;
    }

    .user-dropdown {
        width: 200px;
        right: -8px;
    }
}

@media (max-width: 480px) {
    .header-actions .btn-create span {
        display: none;
    }

    .btn-create {
        padding: 10px;
    }

    .card-content {
        padding: 20px;
    }
}
</style>
