import request from '@/utils/request'

export interface Snippet {
  id: string | number
  title: string
  content: string
  language: string
  created_at: string
  updated_at: string
  visibility?: 'private' | 'public'
  owner_id?: number
}

type RawSnippet = Partial<{
  id: string | number
  ID: string | number
  title: string
  Title: string
  short_link: string
  ShortLink: string
  content: string
  Content: string
  language: string
  Language: string
  created_at: string
  CreatedAt: string
  updated_at: string
  UpdatedAt: string
  visibility: 'private' | 'public'
  Visibility: 'private' | 'public'
  owner_id: number
  OwnerID: number
}>

const normalizeSnippet = (raw: RawSnippet): Snippet => {
  const id = raw.id ?? raw.ID
  const title = raw.title ?? raw.Title ?? raw.short_link ?? raw.ShortLink
  const content = raw.content ?? raw.Content
  const language = raw.language ?? raw.Language
  const createdAt = raw.created_at ?? raw.CreatedAt
  const updatedAt = raw.updated_at ?? raw.UpdatedAt ?? createdAt

  if (id === undefined || !title || !content || !language || !createdAt || !updatedAt) {
    throw new Error('Snippet 响应字段缺失')
  }

  return {
    id,
    title,
    content,
    language,
    created_at: createdAt,
    updated_at: updatedAt,
    visibility: raw.visibility ?? raw.Visibility,
    owner_id: raw.owner_id ?? raw.OwnerID,
  }
}

export interface SaveSnippetRequest {
  title: string
  content: string
  language: string
  visibility?: 'private' | 'public'
}

// 获取我的代码片段
// 对应后端: GET /api/v1/me/pastes
export const listMySnippets = () => {
  return request.get<unknown, RawSnippet[]>('/api/v1/me/pastes').then((list) => list.map(normalizeSnippet))
}

// 获取代码片段详情
// 对应后端: GET /api/v1/pastes/:id
export const getSnippet = (id: string) => {
  return request.get<unknown, RawSnippet>(`/api/v1/pastes/${id}`).then(normalizeSnippet)
}

// 创建代码片段
// 对应后端: POST /api/v1/pastes
export const createSnippet = (data: SaveSnippetRequest) => {
  return request.post<unknown, RawSnippet>('/api/v1/pastes', data).then(normalizeSnippet)
}

// 更新代码片段
// 对应后端: PUT /api/v1/pastes/:id
export const updateSnippet = (id: string, data: SaveSnippetRequest) => {
  return request.put<unknown, RawSnippet>(`/api/v1/pastes/${id}`, data).then(normalizeSnippet)
}
