import axios from "axios";

const api = axios.create({
  baseURL: "/api",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
});

// Request interceptor - 自动附加 token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem("token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Response interceptor - 401 时清除 token
api.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response?.status === 401) {
      localStorage.removeItem("token");
      localStorage.removeItem("user");
      // 不在 API 层跳转，交由组件处理
    }
    return Promise.reject(error);
  },
);

// ========== 认证 API ==========
export const authApi = {
  register(data: { username: string; email: string; password: string }) {
    return api.post("/auth/register", data);
  },

  login(data: { email: string; password: string }) {
    return api.post("/auth/login", data);
  },

  getMe() {
    return api.get("/auth/me");
  },

  changePassword(data: { old_password: string; new_password: string }) {
    return api.put("/auth/change-password", data);
  },
};

// ========== 文章 API ==========
export const articleApi = {
  getList(params?: { page?: number; page_size?: number; search?: string }) {
    return api.get("/articles", { params });
  },

  getById(id: number) {
    return api.get(`/articles/${id}`);
  },

  create(data: {
    title: string;
    content: string;
    summary?: string;
    author?: string;
    tags?: string;
    cover_url?: string;
  }) {
    return api.post("/articles", data);
  },

  update(
    id: number,
    data: {
      title: string;
      content: string;
      summary?: string;
      author?: string;
      tags?: string;
      cover_url?: string;
    },
  ) {
    return api.put(`/articles/${id}`, data);
  },

  delete(id: number) {
    return api.delete(`/articles/${id}`);
  },
};

export default api;
