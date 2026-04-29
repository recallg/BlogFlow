export interface Article {
  id: number;
  title: string;
  content: string;
  summary: string;
  author: string;
  tags: string;
  cover_url: string;
  created_at: string;
  updated_at: string;
}

export interface ArticleListResponse {
  articles: Article[];
  total: number;
  page: number;
  page_size: number;
}

export interface CreateArticleInput {
  title: string;
  content: string;
  summary?: string;
  author?: string;
  tags?: string;
  cover_url?: string;
}

export interface UpdateArticleInput {
  title?: string;
  content?: string;
  summary?: string;
  author?: string;
  tags?: string;
  cover_url?: string;
}

// ===== 用户认证类型 =====

export interface User {
  id: number;
  username: string;
  email: string;
  avatar: string;
  bio: string;
  created_at: string;
}

export interface AuthResponse {
  token: string;
  user: User;
}

export interface LoginInput {
  email: string;
  password: string;
}

export interface RegisterInput {
  username: string;
  email: string;
  password: string;
}
