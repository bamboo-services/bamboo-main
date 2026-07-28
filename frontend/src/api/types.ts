// --------------------------------------------------------------------------------
// Copyright (c) 2016-NOW(至今) 筱锋
// Author: 筱锋「xiao_lfeng」(https://www.x-lf.com)
// --------------------------------------------------------------------------------
// 许可证声明：版权所有 (c) 2016-2026 筱锋。保留所有权利。
// 有关MIT许可证的更多信息，请查看项目根目录下的LICENSE文件或访问：
// https://opensource.org/licenses/MIT
// --------------------------------------------------------------------------------

/**
 * 雪花 ID。后端 SnowflakeID 序列化为 JSON 字符串以避免 JS 精度丢失，
 * 前端经 client.ts 的 reviver 统一还原为 bigint，精确无损。
 * 注意：作为 React key / queryKey / URL 参数时需 toString()。
 */
export type SnowflakeID = bigint

/**
 * 后端统一响应包装（见 bamboo-base-go BaseResponse）
 * - code: HTTP 风格状态码，成功为 200
 * - data: 业务数据，失败时省略
 */
export interface BaseResponse<T> {
  context: string
  output: string
  code: number
  message: string
  error_message?: string
  overhead?: number
  data: T
}

/** 分页信息 */
export interface PaginationInfo {
  page: number
  page_size: number
  total: number
  total_pages: number
  has_next: boolean
  has_prev: boolean
}

/** 分页响应（data + pagination） */
export interface PaginationResponse<T> {
  data: Array<T>
  pagination: PaginationInfo
}

// ---------------------------------------------------------------------------
// 站点信息
// ---------------------------------------------------------------------------

/** 站点信息（GET /api/v1/info/site） */
export interface SiteInfoResponse {
  site_name: string
  site_description: string
  introduction: string
  updated_at: string
}

/** 自我介绍（GET /api/v1/info/about） */
export interface AboutResponse {
  content: string
  updated_at: string
}

/** 更新站点信息请求（PUT /api/v1/admin/info/site） */
export interface UpdateSiteRequest {
  site_name?: string
  site_description?: string
  introduction?: string
}

/** 更新自我介绍请求（PUT /api/v1/admin/info/about） */
export interface UpdateAboutRequest {
  content: string
}

// ---------------------------------------------------------------------------
// 友链分组
// ---------------------------------------------------------------------------

/** 友链分组实体（后端 entity.LinkGroup；created_at 不序列化） */
export interface LinkGroup {
  id: SnowflakeID
  name: string
  description: string | null
  sort_order: number
  status: boolean
  updated_at: string
}

/** 添加友链分组请求（POST /api/v1/admin/groups） */
export interface CreateGroupRequest {
  group_name: string
  group_desc?: string
  group_order?: number
}

/** 更新友链分组请求（PUT /api/v1/admin/groups/:id） */
export interface UpdateGroupRequest {
  group_name?: string
  group_desc?: string
  group_order?: number
  group_status?: number // 0=禁用 1=启用
}

/** 友链分组列表查询参数（GET /api/v1/admin/groups/all） */
export interface GroupListParams {
  status?: number
  name?: string
  with_links?: boolean
  only_enabled?: boolean
  order_by?: 'name' | 'sort_order' | 'created_at'
  order?: 'asc' | 'desc'
}

/** 友链分组分页查询参数（GET /api/v1/admin/groups） */
export interface GroupPageParams extends GroupListParams {
  page?: number
  page_size?: number
}

// ---------------------------------------------------------------------------
// 友链颜色
// ---------------------------------------------------------------------------

/** 友链颜色实体（后端 entity.LinkColor；type 0=普通 1=炫彩） */
export interface LinkColor {
  id: SnowflakeID
  name: string
  type: number
  primary_color: string | null
  sub_color: string | null
  hover_color: string | null
  sort_order: number
  status: boolean
  updated_at: string
}

/** 添加友链颜色请求（POST /api/v1/admin/colors） */
export interface CreateColorRequest {
  color_name: string
  color_type?: number
  primary_color?: string
  sub_color?: string
  hover_color?: string
  color_order?: number
}

/** 更新友链颜色请求（PUT /api/v1/admin/colors/:id） */
export interface UpdateColorRequest {
  color_name?: string
  color_type?: number
  primary_color?: string
  sub_color?: string
  hover_color?: string
  color_order?: number
}

/** 友链颜色列表查询参数（GET /api/v1/admin/colors/all） */
export interface ColorListParams {
  status?: number
  type?: number
  name?: string
  only_enabled?: boolean
  order_by?: 'name' | 'sort_order' | 'created_at'
  order?: 'asc' | 'desc'
}

/** 友链颜色分页查询参数（GET /api/v1/admin/colors） */
export interface ColorPageParams extends ColorListParams {
  page?: number
  page_size?: number
}

// ---------------------------------------------------------------------------
// 友情链接
// ---------------------------------------------------------------------------

/**
 * 友情链接实体（后端 entity.LinkFriend）。
 * - status: 0=待审核 1=已通过 2=已拒绝
 * - is_failure: 0=正常 1=失效
 * - group_f_key / color_f_key: 列表与详情接口已 Preload 的嵌套对象
 * - created_at 不序列化，仅有 updated_at
 */
export interface LinkFriend {
  id: SnowflakeID
  name: string
  url: string
  avatar: string | null
  rss: string | null
  description: string | null
  email: string | null
  group_id: SnowflakeID | null
  color_id: SnowflakeID | null
  sort_order: number
  status: number
  is_failure: number
  level: number
  fail_reason: string | null
  apply_remark: string | null
  review_remark: string | null
  updated_at: string
  group_f_key?: LinkGroup | null
  color_f_key?: LinkColor | null
}

/** 公开友链列表响应（GET /api/v1/links） */
export interface FriendPublicResponse {
  links: Array<LinkFriend>
}

/** 管理端友链分页查询参数（GET /api/v1/admin/links） */
export interface LinkListParams {
  page?: number
  page_size?: number
  link_name?: string
  link_status?: number
  link_fail?: number
  link_group_id?: SnowflakeID
  sort_by?: 'created_at' | 'updated_at' | 'link_order' | 'link_name'
  sort_order?: 'asc' | 'desc'
}

/** 添加友链请求（POST /api/v1/admin/links） */
export interface CreateLinkRequest {
  link_name: string
  link_url: string
  link_avatar?: string
  link_rss?: string
  link_desc?: string
  link_email?: string
  link_group_id?: SnowflakeID
  link_color_id?: SnowflakeID
  link_order?: number
  link_level?: number
  link_apply_remark?: string
}

/** 更新友链请求（PUT /api/v1/admin/links/:id） */
export interface UpdateLinkRequest {
  link_name?: string
  link_url?: string
  link_avatar?: string
  link_rss?: string
  link_desc?: string
  link_email?: string
  link_group_id?: SnowflakeID
  link_color_id?: SnowflakeID
  link_order?: number
  link_level?: number
  link_apply_remark?: string
}

/** 更新友链审核状态请求（PUT /api/v1/admin/links/:id/status） */
export interface UpdateLinkStatusRequest {
  link_status: number
  link_review_remark?: string
}

/** 更新友链失效状态请求（PUT /api/v1/admin/links/:id/fail） */
export interface UpdateLinkFailRequest {
  link_fail: number
  link_fail_reason?: string
}

// ---------------------------------------------------------------------------
// 赞助渠道
// ---------------------------------------------------------------------------

/** 赞助渠道公开条目（GET /api/v1/sponsors/channels） */
export interface SponsorChannel {
  id: SnowflakeID
  name: string
  icon: string | null
  sort_order: number
  status: boolean
  sponsor_count: number
}

/** 赞助渠道管理端实体响应（含描述与时间戳） */
export interface SponsorChannelAdmin {
  id: SnowflakeID
  name: string
  icon: string | null
  description: string | null
  sort_order: number
  status: boolean
  sponsor_count: number
  created_at: string
  updated_at: string
}

/** 添加赞助渠道请求（POST /api/v1/admin/sponsors/channels） */
export interface CreateChannelRequest {
  name: string
  icon?: string
  description?: string
  sort_order?: number
}

/** 更新赞助渠道请求（PUT /api/v1/admin/sponsors/channels/:id） */
export interface UpdateChannelRequest {
  name?: string
  icon?: string
  description?: string
  sort_order?: number
}

/** 赞助渠道分页查询参数（GET /api/v1/admin/sponsors/channels） */
export interface ChannelPageParams {
  page?: number
  page_size?: number
  status?: boolean
  name?: string
  order_by?: 'name' | 'sort_order' | 'created_at'
  order?: 'asc' | 'desc'
}

// ---------------------------------------------------------------------------
// 赞助记录
// ---------------------------------------------------------------------------

/** 赞助记录中的渠道简表 */
export interface SponsorChannelSimple {
  id: SnowflakeID
  name: string
  icon: string | null
}

/** 赞助记录公开条目（GET /api/v1/sponsors/records） */
export interface SponsorRecord {
  id: SnowflakeID
  nickname: string
  redirect_url: string | null
  amount: number
  message: string | null
  sponsor_at: string | null
  channel: SponsorChannelSimple | null
}

/** 赞助记录管理端实体响应 */
export interface SponsorRecordAdmin {
  id: SnowflakeID
  nickname: string
  redirect_url: string | null
  amount: number
  channel_id: SnowflakeID | null
  message: string | null
  sponsor_at: string | null
  sort_order: number
  is_anonymous: boolean
  is_hidden: boolean
  created_at: string
  updated_at: string
  channel?: SponsorChannelSimple | null
}

/** 添加赞助记录请求（POST /api/v1/admin/sponsors/records） */
export interface CreateRecordRequest {
  nickname: string
  redirect_url?: string
  amount: number // 单位：分
  channel_id?: SnowflakeID
  message?: string
  sponsor_at?: string
  sort_order?: number
  is_anonymous?: boolean
  is_hidden?: boolean
}

/** 更新赞助记录请求（PUT /api/v1/admin/sponsors/records/:id） */
export interface UpdateRecordRequest {
  nickname?: string
  redirect_url?: string
  amount?: number
  channel_id?: SnowflakeID
  message?: string
  sponsor_at?: string
  sort_order?: number
  is_anonymous?: boolean
  is_hidden?: boolean
}

/** 赞助记录分页查询参数（GET /api/v1/admin/sponsors/records） */
export interface RecordPageParams {
  page?: number
  page_size?: number
  channel_id?: SnowflakeID
  nickname?: string
  is_anonymous?: boolean
  is_hidden?: boolean
  order_by?: 'nickname' | 'amount' | 'sponsor_at' | 'sort_order' | 'created_at'
  order?: 'asc' | 'desc'
}

// ---------------------------------------------------------------------------
// 认证 / 用户
// ---------------------------------------------------------------------------

/** 系统用户（后端 entity.SystemUser，敏感字段 password/oauth_user_id 不输出） */
export interface UserInfo {
  id: SnowflakeID
  username: string
  email: string
  nickname?: string | null
  avatar?: string | null
  role: string
  status: number
  email_verify: boolean
  last_login_at?: string | null
  created_at: string
  updated_at: string
}

/** 登录请求（POST /api/v1/auth/login） */
export interface LoginRequest {
  username: string
  password: string
}

/** 登录响应（密码登录与 SSO 登录共用） */
export interface LoginResponse {
  user: UserInfo
  token: string
  created_at: string
  expired_at: string
}

/** 当前用户信息响应（GET /api/v1/auth/user） */
export interface UserInfoResponse {
  user: UserInfo
}

/** SSO OAuth 回调返回的令牌（GET /api/v1/sso/oauth/callback） */
export interface OAuthToken {
  access_token: string
  token_type: string
  refresh_token: string
  expiry: string
}

// ---------------------------------------------------------------------------
// 仪表盘
// ---------------------------------------------------------------------------

/** 最近友链申请条目 */
export interface RecentApplicationItem {
  id: SnowflakeID
  name: string
  avatar: string
  url: string
  created_at: string
}

/** 仪表盘统计响应（GET /api/v1/admin/dashboard/stats） */
export interface DashboardStats {
  total_links: number
  pending_links: number
  approved_links: number
  recent_applications: Array<RecentApplicationItem>
}

// ---------------------------------------------------------------------------
// 健康检查
// ---------------------------------------------------------------------------

/** 健康检查系统信息 */
export interface HealthSystemInfo {
  version: string
  environment: string
  platform: string
  go_version: string
}

/** 健康检查运行时信息 */
export interface HealthRuntimeInfo {
  uptime: string
  goroutines: number
  memory_usage: string
  cpu_usage: string
}

/** 健康检查响应（GET /api/v1/public/health） */
export interface HealthResponse {
  status: string
  timestamp: string
  system: HealthSystemInfo
  runtime: HealthRuntimeInfo
}
