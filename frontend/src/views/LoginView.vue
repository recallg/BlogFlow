<template>
  <div class="login-page">
    <!-- ambient background -->
    <div class="ambient-bg">
      <div class="orb orb-1"></div>
      <div class="orb orb-2"></div>
      <div class="orb orb-3"></div>
    </div>

    <div class="login-container">
      <!-- brand -->
      <div class="brand" @click="$router.push('/')">
        <div class="brand-icon">
          <svg viewBox="0 0 32 32" fill="none">
            <rect x="2" y="2" width="28" height="28" rx="8" stroke="url(#lg)" stroke-width="2.5"/>
            <path d="M10 22V12l6 6 6-6v10" stroke="url(#lg)" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"/>
            <defs><linearGradient id="lg" x1="2" y1="2" x2="30" y2="30"><stop stop-color="#6366f1"/><stop offset="1" stop-color="#ec4899"/></linearGradient></defs>
          </svg>
        </div>
        <span class="brand-name">BlogFlow</span>
      </div>

      <!-- card -->
      <div class="auth-card">
        <div class="card-glow"></div>
        <div class="card-content">
          <h1 class="card-title">欢迎回来</h1>
          <p class="card-sub">登录你的账户继续创作</p>

          <form @submit.prevent="handleLogin" class="auth-form">
            <!-- email -->
            <div class="field">
              <label class="field-label">
                <svg class="field-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="4" width="20" height="16" rx="2"/><path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/></svg>
                邮箱
              </label>
              <input v-model="form.email" type="email" class="field-input" placeholder="your@email.com" required />
            </div>

            <!-- password -->
            <div class="field">
              <label class="field-label">
                <svg class="field-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="11" width="18" height="11" rx="2" ry="2"/><path d="M7 11V7a5 5 0 0 1 10 0v4"/></svg>
                密码
              </label>
              <div class="password-wrapper">
                <input v-model="form.password" :type="showPwd ? 'text' : 'password'" class="field-input" placeholder="••••••••" required />
                <button type="button" class="pwd-toggle" @click="showPwd = !showPwd">
                  {{ showPwd ? '🙈' : '👁️' }}
                </button>
              </div>
            </div>

            <!-- error -->
            <div v-if="error" class="error-msg">
              <span>⚠️</span> {{ error }}
            </div>

            <!-- submit -->
            <button type="submit" class="btn-submit" :disabled="loading">
              <span v-if="loading" class="spinner"></span>
              <span v-else>登录</span>
            </button>
          </form>

          <div class="auth-footer">
            还没有账户？
            <router-link to="/register" class="auth-link">立即注册</router-link>
          </div>

          <router-link to="/" class="back-link">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M19 12H5M12 19l-7-7 7-7"/></svg>
            返回首页
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useRouter } from "vue-router";
import { authApi } from "../api";

const router = useRouter();
const loading = ref(false);
const error = ref("");
const showPwd = ref(false);

const form = ref({
  email: "",
  password: "",
});

async function handleLogin() {
  if (!form.value.email || !form.value.password) return;

  loading.value = true;
  error.value = "";

  try {
    const res = await authApi.login({
      email: form.value.email,
      password: form.value.password,
    });

    const { token, user } = res.data;
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(user));

    router.push("/");
  } catch (e: any) {
    error.value = e.response?.data?.error || "登录失败，请检查邮箱和密码";
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
.login-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #05050a;
  padding: 24px;
  position: relative;
  overflow: hidden;
}

/* ambient */
.ambient-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
}

.orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(100px);
  opacity: 0.4;
}

.orb-1 {
  width: 500px; height: 500px;
  top: -200px; right: -100px;
  background: #6366f1;
  animation: float1 10s ease-in-out infinite;
}

.orb-2 {
  width: 400px; height: 400px;
  bottom: -150px; left: -150px;
  background: #ec4899;
  animation: float2 12s ease-in-out infinite;
}

.orb-3 {
  width: 300px; height: 300px;
  top: 40%; left: 60%;
  background: #06b6d4;
  opacity: 0.2;
  animation: float1 8s ease-in-out infinite reverse;
}

@keyframes float1 {
  0%, 100% { transform: translate(0, 0); }
  33% { transform: translate(30px, -40px); }
  66% { transform: translate(-20px, 20px); }
}
@keyframes float2 {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-40px, 30px); }
}

/* layout */
.login-container {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

/* brand */
.brand {
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
}

.brand-icon {
  width: 40px;
  height: 40px;
}

.brand-icon svg {
  width: 100%;
  height: 100%;
}

.brand-name {
  font-size: 1.5rem;
  font-weight: 800;
  background: linear-gradient(135deg, #e0e7ff, #818cf8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* card */
.auth-card {
  position: relative;
  width: 100%;
  border-radius: 24px;
  overflow: hidden;
}

.card-glow {
  position: absolute;
  inset: -2px;
  border-radius: 26px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6, #ec4899, #6366f1);
  background-size: 300% 300%;
  animation: glowShift 4s ease-in-out infinite;
  opacity: 0.3;
}

@keyframes glowShift {
  0%, 100% { background-position: 0% 50%; }
  50% { background-position: 100% 50%; }
}

.card-content {
  position: relative;
  background: rgba(10, 10, 25, 0.95);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  padding: 40px 32px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.card-title {
  font-size: 1.8rem;
  font-weight: 800;
  color: #f1f5f9;
  text-align: center;
  margin-bottom: 6px;
}

.card-sub {
  font-size: 0.9rem;
  color: #64748b;
  text-align: center;
  margin-bottom: 32px;
}

/* form */
.auth-form {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field-label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #cbd5e1;
}

.field-icon {
  width: 16px;
  height: 16px;
  color: #6366f1;
}

.field-input {
  width: 100%;
  padding: 14px 16px;
  background: rgba(0, 0, 0, 0.4);
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 12px;
  color: #f1f5f9;
  font-size: 0.95rem;
  outline: none;
  transition: all 0.3s ease;
  font-family: inherit;
}

.field-input::placeholder {
  color: #334155;
}

.field-input:focus {
  border-color: rgba(99, 102, 241, 0.5);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.1), 0 0 20px rgba(99, 102, 241, 0.05);
}

.password-wrapper {
  position: relative;
}

.pwd-toggle {
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

/* submit */
.btn-submit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 14px;
  margin-top: 8px;
  border: none;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: white;
  font-size: 1rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.3s ease;
  font-family: inherit;
}

.btn-submit:hover:not(:disabled) {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

.btn-submit:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner {
  width: 20px;
  height: 20px;
  border: 3px solid rgba(255,255,255,0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* footer */
.auth-footer {
  text-align: center;
  margin-top: 24px;
  color: #64748b;
  font-size: 0.9rem;
}

.auth-link {
  color: #818cf8;
  font-weight: 600;
  text-decoration: none;
  transition: color 0.3s;
}

.auth-link:hover {
  color: #a5b4fc;
  text-decoration: underline;
}

.back-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 16px;
  color: #475569;
  font-size: 0.85rem;
  text-decoration: none;
  transition: color 0.3s;
}

.back-link svg {
  width: 14px;
  height: 14px;
}

.back-link:hover {
  color: #94a3b8;
}

/* responsive */
@media (max-width: 480px) {
  .login-page {
    padding: 16px;
  }

  .card-content {
    padding: 28px 20px;
  }

  .card-title {
    font-size: 1.5rem;
  }
}
</style>
