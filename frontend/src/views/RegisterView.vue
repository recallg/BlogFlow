<template>
  <div class="register-page">
    <!-- Ambient Background -->
    <div class="ambient-bg">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <!-- Particle Effects -->
    <div class="particles">
      <div class="particle p1"></div>
      <div class="particle p2"></div>
      <div class="particle p3"></div>
      <div class="particle p4"></div>
    </div>

    <div class="auth-container">
      <!-- Brand Section -->
      <div class="brand-section">
        <div class="brand-content">
          <div class="brand-icon" @click="$router.push('/')">
            <svg viewBox="0 0 32 32" fill="none" xmlns="http://www.w3.org/2000/svg">
              <rect x="2" y="2" width="28" height="28" rx="8" stroke="url(#brand-grad)" stroke-width="2.5"/>
              <path d="M10 22V12l6 6 6-6v10" stroke="url(#brand-grad)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
              <defs>
                <linearGradient id="brand-grad" x1="2" y1="2" x2="30" y2="30">
                  <stop stop-color="#6366f1"/>
                  <stop offset="1" stop-color="#ec4899"/>
                </linearGradient>
              </defs>
            </svg>
          </div>
          <h1 class="brand-title">BlogFlow</h1>
          <p class="brand-desc">加入我们的创作者社区<br/>让每一个想法都闪闪发光</p>
          <div class="brand-features">
            <div class="feature-item">
              <span class="feature-icon">✍️</span>
              <span>自由创作</span>
            </div>
            <div class="feature-item">
              <span class="feature-icon">🌐</span>
              <span>分享观点</span>
            </div>
            <div class="feature-item">
              <span class="feature-icon">🤝</span>
              <span>结识同好</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Form Section -->
      <div class="form-section">
        <div class="form-card">
          <div class="form-header">
            <h2 class="form-title">创建账号</h2>
            <p class="form-subtitle">开始你的创作之旅</p>
          </div>

          <form @submit.prevent="handleRegister" class="auth-form">
            <!-- 用户名 -->
            <div class="input-group">
              <label class="input-label">
                <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
                  <circle cx="12" cy="7" r="4"/>
                </svg>
                用户名
              </label>
              <input
                v-model="form.username"
                type="text"
                class="form-input"
                placeholder="输入用户名"
                required
                autocomplete="username"
              />
            </div>

            <!-- 邮箱 -->
            <div class="input-group">
              <label class="input-label">
                <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="2" y="4" width="20" height="16" rx="2"/>
                  <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
                </svg>
                邮箱
              </label>
              <input
                v-model="form.email"
                type="email"
                class="form-input"
                placeholder="输入邮箱地址"
                required
                autocomplete="email"
              />
            </div>

            <!-- 密码 -->
            <div class="input-group">
              <label class="input-label">
                <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                密码
              </label>
              <div class="password-wrapper">
                <input
                  v-model="form.password"
                  :type="showPassword ? 'text' : 'password'"
                  class="form-input"
                  placeholder="至少 6 位密码"
                  required
                  autocomplete="new-password"
                  minlength="6"
                />
                <button type="button" class="password-toggle" @click="showPassword = !showPassword">
                  {{ showPassword ? '🙈' : '👁️' }}
                </button>
              </div>
            </div>

            <!-- 确认密码 -->
            <div class="input-group">
              <label class="input-label">
                <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                  <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
                确认密码
              </label>
              <input
                v-model="form.confirmPassword"
                :type="showPassword ? 'text' : 'password'"
                class="form-input"
                placeholder="再次输入密码"
                required
                autocomplete="new-password"
              />
            </div>

            <!-- 错误提示 -->
            <div v-if="error" class="error-msg">
              <svg class="error-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                <circle cx="12" cy="12" r="10"/>
                <line x1="12" y1="8" x2="12" y2="12"/>
                <line x1="12" y1="16" x2="12.01" y2="16"/>
              </svg>
              {{ error }}
            </div>

            <!-- 提交按钮 -->
            <button type="submit" class="submit-btn" :disabled="loading">
              <span v-if="loading" class="loading-spinner"></span>
              <span v-else class="btn-text">注册</span>
            </button>
          </form>

          <!-- 底部链接 -->
          <div class="form-footer">
            <span class="footer-text">已有账号？</span>
            <router-link to="/login" class="footer-link">立即登录</router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from "vue";
import { useRouter } from "vue-router";
import { authApi } from "../api";

const router = useRouter();
const loading = ref(false);
const error = ref("");
const showPassword = ref(false);

const form = reactive({
  username: "",
  email: "",
  password: "",
  confirmPassword: "",
});

async function handleRegister() {
  error.value = "";

  // 基础验证
  if (!form.username.trim()) {
    error.value = "请输入用户名";
    return;
  }
  if (!form.email.trim()) {
    error.value = "请输入邮箱";
    return;
  }
  if (form.password.length < 6) {
    error.value = "密码至少 6 位";
    return;
  }
  if (form.password !== form.confirmPassword) {
    error.value = "两次输入的密码不一致";
    return;
  }

  loading.value = true;
  try {
    const res = await authApi.register({
      username: form.username.trim(),
      email: form.email.trim(),
      password: form.password,
    });
    const { token, user } = res.data;
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(user));
    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || "注册失败，请重试";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* ===== Page Layout ===== */
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #05050a;
  position: relative;
  overflow: hidden;
  padding: 24px;
}

/* ===== Ambient Background ===== */
.ambient-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
}

.orb-1 {
  width: 500px;
  height: 500px;
  top: -200px;
  right: -100px;
  background: #6366f1;
  opacity: 0.3;
  animation: orb1 8s ease-in-out infinite;
}

.orb-2 {
  width: 400px;
  height: 400px;
  bottom: -150px;
  left: -100px;
  background: #ec4899;
  opacity: 0.25;
  animation: orb2 10s ease-in-out infinite;
}

.orb-3 {
  width: 300px;
  height: 300px;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: #6366f1;
  opacity: 0.1;
  animation: orb3 12s ease-in-out infinite;
}

@keyframes orb1 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(40px, -30px) scale(1.1); }
}

@keyframes orb2 {
  0%, 100% { transform: translate(0, 0) scale(1); }
  50% { transform: translate(-30px, 30px) scale(1.05); }
}

@keyframes orb3 {
  0%, 100% { opacity: 0.1; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 0.2; transform: translate(-50%, -50%) scale(1.3); }
}

/* ===== Particles ===== */
.particles {
  position: fixed;
  inset: 0;
  pointer-events: none;
}

.particle {
  position: absolute;
  width: 4px;
  height: 4px;
  border-radius: 50%;
  background: #6366f1;
  opacity: 0.3;
}

.p1 { top: 20%; left: 10%; animation: float1 6s ease-in-out infinite; }
.p2 { top: 70%; left: 80%; width: 3px; height: 3px; background: #ec4899; animation: float2 8s ease-in-out infinite; }
.p3 { top: 30%; left: 85%; width: 5px; height: 5px; background: #06b6d4; animation: float3 7s ease-in-out infinite; }
.p4 { top: 80%; left: 15%; background: #f59e0b; animation: float1 9s ease-in-out infinite; }

@keyframes float1 {
  0%, 100% { transform: translate(0, 0); opacity: 0.3; }
  50% { transform: translate(20px, -20px); opacity: 0.6; }
}

@keyframes float2 {
  0%, 100% { transform: translate(0, 0); opacity: 0.3; }
  50% { transform: translate(-15px, 15px); opacity: 0.5; }
}

@keyframes float3 {
  0%, 100% { transform: translate(0, 0); opacity: 0.3; }
  50% { transform: translate(10px, -25px); opacity: 0.5; }
}

/* ===== Auth Container ===== */
.auth-container {
  display: flex;
  width: 100%;
  max-width: 960px;
  min-height: 600px;
  border-radius: 24px;
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.06);
  backdrop-filter: blur(20px);
  overflow: hidden;
  position: relative;
  z-index: 1;
  animation: fadeIn 0.6s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(20px) scale(0.98); }
  to { opacity: 1; transform: translateY(0) scale(1); }
}

/* ===== Brand Section ===== */
.brand-section {
  flex: 1;
  background: linear-gradient(135deg, rgba(99, 102, 241, 0.08), rgba(236, 72, 153, 0.05));
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
  position: relative;
  overflow: hidden;
}

.brand-section::before {
  content: "";
  position: absolute;
  inset: 0;
  background:
    radial-gradient(circle at 20% 50%, rgba(99, 102, 241, 0.05), transparent 50%),
    radial-gradient(circle at 80% 50%, rgba(236, 72, 153, 0.05), transparent 50%);
}

.brand-content {
  text-align: center;
  position: relative;
  z-index: 1;
}

.brand-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 24px;
  cursor: pointer;
  transition: transform 0.3s ease;
}

.brand-icon:hover {
  transform: scale(1.05);
}

.brand-icon svg {
  width: 100%;
  height: 100%;
}

.brand-title {
  font-size: 2.2rem;
  font-weight: 800;
  background: linear-gradient(135deg, #e0e7ff, #818cf8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  margin-bottom: 12px;
}

.brand-desc {
  color: #64748b;
  font-size: 1rem;
  line-height: 1.6;
  margin-bottom: 32px;
}

.brand-features {
  display: flex;
  flex-direction: column;
  gap: 16px;
  align-items: center;
}

.feature-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 20px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.05);
  color: #94a3b8;
  font-size: 0.9rem;
  width: 180px;
  transition: all 0.3s ease;
}

.feature-item:hover {
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(99, 102, 241, 0.15);
  transform: translateX(4px);
}

.feature-icon {
  font-size: 1.2rem;
}

/* ===== Form Section ===== */
.form-section {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
}

.form-card {
  width: 100%;
  max-width: 400px;
}

.form-header {
  text-align: center;
  margin-bottom: 32px;
}

.form-title {
  font-size: 1.8rem;
  font-weight: 700;
  color: #f1f5f9;
  margin-bottom: 8px;
}

.form-subtitle {
  color: #64748b;
  font-size: 0.95rem;
}

/* ===== Form Fields ===== */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.input-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.input-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 500;
  color: #94a3b8;
}

.input-icon {
  width: 16px;
  height: 16px;
}

.form-input {
  width: 100%;
  padding: 14px 16px;
  background: rgba(0, 0, 0, 0.3);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 0.95rem;
  outline: none;
  transition: all 0.3s ease;
  font-family: inherit;
  box-sizing: border-box;
}

.form-input::placeholder {
  color: #475569;
}

.form-input:focus {
  border-color: rgba(99, 102, 241, 0.4);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1);
}

.password-wrapper {
  position: relative;
}

.password-wrapper .form-input {
  padding-right: 48px;
}

.password-toggle {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  font-size: 1.1rem;
  cursor: pointer;
  padding: 4px;
  line-height: 1;
}

/* ===== Error Message ===== */
.error-msg {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 12px 16px;
  border-radius: 10px;
  background: rgba(239, 68, 68, 0.1);
  border: 1px solid rgba(239, 68, 68, 0.2);
  color: #fca5a5;
  font-size: 0.88rem;
}

.error-icon {
  width: 18px;
  height: 18px;
  flex-shrink: 0;
}

/* ===== Submit Button ===== */
.submit-btn {
  width: 100%;
  padding: 14px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 48px;
  margin-top: 4px;
}

.submit-btn:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.loading-spinner {
  width: 22px;
  height: 22px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ===== Footer ===== */
.form-footer {
  text-align: center;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

.footer-text {
  color: #64748b;
  font-size: 0.9rem;
}

.footer-link {
  color: #818cf8;
  font-size: 0.9rem;
  font-weight: 600;
  text-decoration: none;
  margin-left: 6px;
  transition: color 0.3s ease;
}

.footer-link:hover {
  color: #a5b4fc;
  text-decoration: underline;
}

/* ===== Responsive ===== */
@media (max-width: 768px) {
  .register-page {
    padding: 16px;
    align-items: flex-start;
    padding-top: 40px;
  }

  .auth-container {
    flex-direction: column;
    max-width: 100%;
    min-height: auto;
  }

  .brand-section {
    display: none;
  }

  .form-section {
    padding: 32px 24px;
  }

  .form-card {
    max-width: 100%;
  }
}
</style>
