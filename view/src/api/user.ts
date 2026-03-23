import request from '@/utils/request'

export interface User {
  id: number
  username: string
  email: string
}

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export interface LoginResponse {
  token: string
  refresh_token: string
  user_id: number
  username: string
  email: string
}

export interface RegisterResponse {
  user_id: number
  username: string
  email: string
}

type RawLoginResponse = Partial<{
  token: string
  Token: string
  access_token: string
  refresh_token: string
  Refresh_token: string
  user_id: number
  UserID: number
  username: string
  Username: string
  email: string
  Email: string
}>

type RawRegisterResponse = Partial<{
  user_id: number
  UserID: number
  username: string
  Username: string
  email: string
  Email: string
}>

const normalizeLoginResponse = (raw: RawLoginResponse): LoginResponse => {
  const token = raw.token ?? raw.Token ?? raw.access_token
  const refreshToken = raw.refresh_token ?? raw.Refresh_token
  const userID = raw.user_id ?? raw.UserID
  const username = raw.username ?? raw.Username
  const email = raw.email ?? raw.Email ?? ''

  if (!token || userID === undefined || !username) {
    throw new Error('登录响应字段缺失')
  }

  return {
    token,
    refresh_token: refreshToken ?? '',
    user_id: userID,
    username,
    email,
  }
}

const normalizeRegisterResponse = (raw: RawRegisterResponse): RegisterResponse => {
  const userID = raw.user_id ?? raw.UserID
  const username = raw.username ?? raw.Username
  const email = raw.email ?? raw.Email

  if (userID === undefined || !username || !email) {
    throw new Error('注册响应字段缺失')
  }

  return {
    user_id: userID,
    username,
    email,
  }
}

// 用户登录
// 对应后端: POST /api/v1/users/login
export const login = (data: LoginRequest) => {
  return request.post<unknown, RawLoginResponse>('/api/v1/users/login', data).then(normalizeLoginResponse)
}

// 用户注册
// 对应后端: POST /api/v1/users/register
export const register = (data: RegisterRequest) => {
  return request
    .post<unknown, RawRegisterResponse>('/api/v1/users/register', data)
    .then(normalizeRegisterResponse)
}

// 刷新 Token
// 对应后端: POST /api/v1/users/refresh
export const refreshToken = (refreshToken: string) => {
  return request.post<unknown, { token: string; refresh_token: string }>('/api/v1/users/refresh', {
    refresh_token: refreshToken,
  })
}
