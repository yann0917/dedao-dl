export type ApiEnvelope<T> = {
  code: number
  msg: string
  data: T
}

export type UserInfo = {
  nickname: string
  avatar: string
  today_study_time: number
  study_serial_days: number
  is_teacher: number
  uid_hazy: string
}

export type AuthMeData = {
  loggedIn: boolean
  user?: UserInfo
}

export type AuthAccount = {
  uidHazy: string
  name: string
  avatar: string
  active: boolean
}

export type QRCodeSession = {
  sessionId: string
  qrCode: string
  qrCodeString: string
  expiresAt: number
}

export type QRCodeStatus = {
  status: number
  expiresAt: number
  user?: UserInfo
}

export type CourseCategory = {
  name: string
  count: number
  category: string
}

export type CourseCategoryResponse = {
  data: {
    list: CourseCategory[]
  }
}

export type PurchasedNavbarChild = {
  name: string
  count: number
  filter: string
  show_count: boolean
}

export type PurchasedNavbarItem = {
  id: number
  name: string
  category: string
  channel_type: string
  item_type: string
  children: PurchasedNavbarChild[]
}

export type PurchasedNavbarResponse = {
  list: PurchasedNavbarItem[]
}

export type CourseListItem = {
  id: number
  class_id: number
  enid: string
  type: number
  title: string
  name: string
  intro: string
  subtitle: string
  cover: string
  icon: string
  index_img: string
  author: string
  lecturer_name: string
  price_desc: string
  price: string
  dd_url: string
  duration: number
  progress: number
  publish_num: number
  course_num: number
  is_group: boolean
  group_id: number
  has_play_auth: boolean
  audio_detail?: {
    alias_id: string
    mp3_play_url: string
    icon: string
  } | null
  odob_group_ext_info?: {
    odob_alias_list: string[]
    progress_percent: number
    audio_detail?: {
      alias_id: string
      mp3_play_url: string
      icon: string
    } | null
  } | null
}

export type CourseListResponse = {
  list: CourseListItem[]
  total: number
  is_more: number
}

export type CourseInfoResponse = {
  class_info: {
    enid: string
    name: string
    intro: string
    highlight: string
    lecturer_name: string
    lecturer_title: string
    lecturer_avatar: string
    lecturer_intro: string
    lecturer_name_and_title: string
    learn_user_count: number
    current_article_count: number
    phase_num: number
    is_finished: number
    square_img: string
    index_img: string
    outline_img: string
    share_summary: string
    price_desc: string
    dd_url: string
    share_url: string
    is_subscribe: number
    is_in_vip: boolean
    is_vip: boolean
    collection: {
      is_collected: boolean
      collection_count: number
    }
  }
  user_type: string
  class_reviews_count: number
  class_comment_info?: {
    count: number
    average_score: string
    comment_list: Array<{
      id: number
      title: string
      score: number
      no_style_content: string
      nickname: string
      avatar_s: string
    }>
  }
  items: Array<{
    title: string
    content: string
  }>
}

export type CourseArticleItem = {
  id: number
  enid: string
  class_id: number
  class_enid: string
  title: string
  logo: string
  summary: string
  publish_time: number
  cur_learn_count: number
  is_buy: boolean
  is_read: boolean
  video_status: number
  audio_alias_ids: string[]
  audio?: AudioDetailResponse
}

export type CourseArticleListResponse = {
  article_list: CourseArticleItem[]
  class_id: number
  ptype: number
  pid: number
  reverse: boolean
  chapter_id: string
  max_id: number
}

export type SearchHotResponse = {
  hot_tab_list: Array<{
    name: string
    list: Array<{
      title: string
      searchKey: string
    }>
  }>
}

export type SearchSuggestResponse = {
  list: Array<{
    list: Array<{
      id: number
      title: string
      author: string
      content: string
      extra: {
        enid: string
        image: string
      }
    }>
  }>
}

export type HomeLabel = {
  enid: string
  name: string
}

export type HomeCategory = {
  englishName: string
  enid: string
  icon: string
  id: number
  labelList: HomeLabel[]
  name: string
  navType: number
  relationId: number
  relationName: string
  type: number
}

export type HomeBanner = {
  beginTime: number
  endTime: number
  id: number
  img: string
  isDelete: number
  moduleId: number
  sort: number
  sourceId: string
  title: string
  url: string
  urlType: number
}

export type HomeModule = {
  description: string
  ext1: string
  ext2: string
  ext3: string
  ext4: string
  ext5: string
  id: number
  isShow: number
  name: string
  sort: number
  title: string
  type: string
}

export type HomeData = {
  moduleList: HomeModule[]
  categoryList: HomeCategory[]
  banner: HomeBanner[]
}

export type HomeNavigation = {
  enid: string
  id: number
  relation_id: number
  relation_name: string
  type: number
  nav_type: number
  name: string
  icon: string
  english_name: string
  label_list: HomeLabel[]
}

export type HomeProductSimple = {
  product_type: number
  product_enid: string
  title: string
  intro: string
  introduction: string
  index_image: string
  score: string
  user_score_count: number
  horizontal_image: string
  learn_user_count: number
  author_list: string[]
  trackinfo: string
  log_type: string
}

export type HomeContentResponse = {
  product_list: HomeProductSimple[]
  current_enid: string
  navigation_list: HomeNavigation[]
  page_id: number
  page_size: number
  is_more: number
  request_id: string
}

export type HomeFreeResource = {
  id: number
  enid: string
  name: string
  intro: string
  logo: string
  product_type: number
  product_id: number
  score: number
  class_type: number
  status: number
}

export type HomeFreeResourcesResponse = {
  list: HomeFreeResource[]
}

export type HomeLabelListResponse = {
  list: HomeNavigation[]
}

export type HomePortalResponse = {
  homeData: HomeData
  freeResources?: HomeFreeResourcesResponse
  freeResourcesError?: string
  ebookLabels?: HomeLabelListResponse
  ebookLabelsError?: string
  ebookContent?: HomeContentResponse
  ebookContentError?: string
  courseLabels?: HomeLabelListResponse
  courseLabelsError?: string
  courseContent?: HomeContentResponse
  courseContentError?: string
}

export type AlgoOption = {
  name: string
  value: string
  sub_options?: AlgoOption[]
}

export type AlgoStrategy = {
  title: string
  is_multiple: boolean
  is_hide: boolean
  options: AlgoOption[]
}

export type AlgoFilterRequest = {
  classfc_name: string
  label_id: string
  nav_type: number
  navigation_id: string
  page: number
  page_size: number
  product_types: string
  request_id: string
  sort_strategy: string
}

export type AlgoFilterResponse = {
  filter: {
    title: string
    product_types: AlgoStrategy
    sort_strategy: AlgoStrategy
    progress_strategy: AlgoStrategy
    buy_strategy: AlgoStrategy
    navigations: AlgoStrategy
    tags: AlgoStrategy
  }
  total: number
  request: AlgoFilterRequest
}

export type AlgoProductItem = {
  item_type: number
  id: number
  product_type: number
  product_id: number
  name: string
  intro: string
  phase_num: number
  learn_user_count: number
  price: number
  current_price: number
  index_img: string
  horizontal_image: string
  lecturers_img: string
  price_desc: string
  is_buy: number
  is_vip: boolean
  trackinfo: string
  log_id: string
  log_type: string
  corner_img: string
  lecturer_name_and_title: string
  lecturers_v_status: number
  corner_img_vertical: string
  lecturer_name: string
  lecturer_title: string
  author_list: string[]
  id_out: string
  is_limit_free: boolean
  has_play_auth: boolean
  button_type: number
  audio_id: string
  alias_id: string
  in_bookrack: boolean
  is_show_newbook: boolean
  is_subscribe: boolean
  online_time: string
  need_login: number
  is_on_bookshelf: boolean
  duration: number
  collection_num: number
  dd_url: string
  learning_days_desc: string
  hot_learn_desc: string
  score: string
  introduction: string
  user_score_count: number
  progress: number
  metrics: string
}

export type AlgoProductResponse = {
  product_list: AlgoProductItem[]
  request_id: string
  total: number
  is_more: number
}

export type EbookDetailResponse = {
  detail?: {
    id: number
    title: string
    cover: string
    count: number
    price: string
    current_price: string
    original_price: string
    book_author: string
    author_info: string
    publish_time: string
    book_intro: string
    author_list: string[]
    can_trial_read: boolean
    trial_read_proportion: string
    with_video: boolean
    enid: string
    rank_name: string
    rank_num: number
    is_vip_book: number
    is_on_bookshelf: boolean
    product_score: string
    read_time: number
    douban_score: string
    classify_name: string
    add_studylist_dd_url: string
  }
  notes?: {
    list: Array<{
      note_id: number
      note: string
      content: string
      note_title: string
      note_line: string
      create_time: number
      notes_owner: {
        nickname: string
        avatar: string
      }
    }>
  }
  notesError?: string
}

export type EbookCommentItem = {
  note_id: number
  note_title: string
  note_line: string
  note_line_style: string
  create_time: number
  notes_owner: {
    name: string
    avatar: string
    slogan: string
  }
  notes_count?: {
    like_count: number
    comment_count: number
  }
}

export type EbookCommentResponse = {
  total: number
  list: EbookCommentItem[]
  ebook_score: {
    average_score: string
    score_info: string[]
    total: string
  }
}

export type AudioDetailResponse = {
  topic_encode_id: string
  audio_id: string
  alias_id: string
  icon: string
  index_img: string
  duration: number
  reader_name: string
  summary: string
  title: string
  mp3_play_url: string
  odob_group_enid: string
  slogan: string
  has_play_auth: boolean
  is_vip: boolean
  play_count: number
  share_url: string
  play_dd_url: string
  package_title: string
  source_name: string
  update_tips: string
  trial_listen_tips: string
}

export type AudioGroupResponse = {
  outside?: {
    spu: {
      title: string
      summary: string
      intro: string
      icon: string
      share_title: string
      share_summary: string
      ddurl: string
      extra: {
        enid: string
        intro_text: string
        newest_intro: string
        odob_consumer_num: number
        rn_learn_count_desc: string
        teacher_avatar: string
        teacher_intro: string
        teacher_name: string
      }
    }
    items: Array<{
      title: string
      summary: string
      intro: string
      icon: string
      ddurl: string
      extra: {
        enid: string
        duration: number
        audio_alias_id: string
        teacher_name: string
        odob_audio_detail: AudioDetailResponse
      }
    }>
    count: number
    current_count: number
  }
  group?: {
    odob_audio_detail_list: Array<AudioDetailResponse>
  }
  groupError?: string
}

export type ArticleDetailResponse = {
  info?: {
    class_id: number
    class_enid: string
    ptype: number
    pid: number
    article_id: number
    dd_article_token: string
    article_info: {
      id: number
      id_str: string
      enid: string
      class_enid: string
      title: string
      logo: string
      summary: string
      publish_time: number
      cur_learn_count: number
      audio_alias_ids: string[]
      is_buy: boolean
      is_read: boolean
      audio?: AudioDetailResponse
    }
    class_info: {
      enid: string
      name: string
      lecturer_name: string
      lecturer_title: string
      logo: string
      square_img: string
      share_summary: string
      dd_url: string
    }
    audio: AudioDetailResponse
  }
  detail?: {
    article: {
      Id: number
      AppId: number
      PublishTime: number
      IdStr: string
      AppIdStr: string
    }
    content: string
  }
  markdown: string
}

export type EbookVIPInfo = {
  nickname: string
  slogan: string
  avatar: string
  month_count: number
  total_count: number
  finished_count: number
  save_price: string
  is_vip: boolean
  expire_time: number
  surplus_time: number
  is_expire: boolean
  price_desc: string
  err_tips: string
}

export type OdobVIPInfo = {
  card: Array<{
    id: number
    name: string
    description: string
    price: string
    origin_price: string
    price_desc: string
    discount_tip: string
    subscribe_desc: string
    welfare_info: string
    is_subscribed: number
    selected: number
    rights: Array<{
      name: string
      right: boolean
    }>
  }>
  user: {
    nickname: string
    slogan: string
    avatar: string
    dd_url: string
    is_vip: boolean
    is_expire: boolean
    err_tips: string
    expire_time: number
    surplus_time: number
    week_count: number
    total_count: number
    save_price: string
    price_desc: string
  }
}

export type UserCenterResponse = {
  user?: UserInfo
  ebookVip?: EbookVIPInfo
  ebookVipError?: string
  odobVip?: OdobVIPInfo
  odobVipError?: string
  accounts: AuthAccount[]
}

async function request<T>(input: RequestInfo, init?: RequestInit): Promise<T> {
  const response = await fetch(input, {
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers ?? {}),
    },
    ...init,
  })

  const body = (await response.json()) as ApiEnvelope<T>
  if (!response.ok || body.code !== 0) {
    throw new Error(body.msg || "请求失败")
  }
  return body.data
}

function buildAlgoQuery(params: Partial<AlgoFilterRequest>) {
  const query = new URLSearchParams()

  if (params.classfc_name) query.set("classfcName", params.classfc_name)
  if (params.label_id) query.set("labelId", params.label_id)
  if (typeof params.nav_type === "number") query.set("navType", String(params.nav_type))
  if (params.navigation_id) query.set("navigationId", params.navigation_id)
  if (typeof params.page === "number") query.set("page", String(params.page))
  if (typeof params.page_size === "number") query.set("pageSize", String(params.page_size))
  if (params.product_types) query.set("productTypes", params.product_types)
  if (params.request_id) query.set("requestId", params.request_id)
  if (params.sort_strategy) query.set("sortStrategy", params.sort_strategy)

  return query.toString()
}

export const api = {
  auth: {
    me: () => request<AuthMeData>("/api/auth/me"),
    createQRCode: () => request<QRCodeSession>("/api/auth/qrcode", { method: "POST" }),
    getQRCodeStatus: (sessionId: string) => request<QRCodeStatus>(`/api/auth/qrcode/${sessionId}/status`),
    accounts: () => request<AuthAccount[]>("/api/auth/accounts"),
    switchAccount: (uidHazy: string) =>
      request<AuthMeData>("/api/auth/switch", {
        method: "POST",
        body: JSON.stringify({ uidHazy }),
      }),
    logout: () => request<AuthMeData>("/api/auth/logout", { method: "POST" }),
  },
  user: {
    info: () => request<UserInfo>("/api/user/info"),
    center: () => request<UserCenterResponse>("/api/user/center"),
  },
  course: {
    categories: () => request<CourseCategoryResponse>("/api/course/categories"),
    navbar: () => request<PurchasedNavbarResponse>("/api/course/navbar"),
    list: (params: URLSearchParams) => request<CourseListResponse>(`/api/course/list?${params.toString()}`),
    info: (enid: string) => request<CourseInfoResponse>(`/api/course/info?enid=${encodeURIComponent(enid)}`),
    articles: (enid: string, options?: { count?: number; maxId?: number; reverse?: boolean }) => {
      const query = new URLSearchParams({
        enid,
        count: String(options?.count ?? 30),
        maxId: String(options?.maxId ?? 0),
        reverse: options?.reverse ? "1" : "0",
      })
      return request<CourseArticleListResponse>(`/api/course/articles?${query.toString()}`)
    },
  },
  search: {
    hot: () => request<SearchHotResponse>("/api/search/hot"),
    suggest: (query: string, searchType = 10) =>
      request<SearchSuggestResponse>(`/api/search/suggest?query=${encodeURIComponent(query)}&searchType=${searchType}`),
  },
  home: {
    portal: () => request<HomePortalResponse>("/api/home/portal"),
    labelContent: (type: 2 | 4, enid: string, page = 0, pageSize = type === 2 ? 10 : 4) =>
      request<HomeContentResponse>(
        `/api/home/label-content?type=${type}&enid=${encodeURIComponent(enid)}&page=${page}&pageSize=${pageSize}`,
      ),
  },
  algo: {
    filter: (params: Partial<AlgoFilterRequest>) => request<AlgoFilterResponse>(`/api/algo/filter?${buildAlgoQuery(params)}`),
    products: (params: Partial<AlgoFilterRequest>) =>
      request<AlgoProductResponse>(`/api/algo/products?${buildAlgoQuery(params)}`),
  },
  ebook: {
    detail: (enid: string) => request<EbookDetailResponse>(`/api/ebook/detail?enid=${encodeURIComponent(enid)}`),
    comments: (enid: string, page = 1, limit = 15, sort = "like_count") =>
      request<EbookCommentResponse>(
        `/api/ebook/comments?enid=${encodeURIComponent(enid)}&page=${page}&limit=${limit}&sort=${encodeURIComponent(sort)}`,
      ),
  },
  audio: {
    detail: (enid: string) => request<AudioDetailResponse>(`/api/audio/detail?enid=${encodeURIComponent(enid)}`),
    group: (enid: string) => request<AudioGroupResponse>(`/api/audio/group?enid=${encodeURIComponent(enid)}`),
  },
  article: {
    detail: (aType: 1 | 2, enid: string) =>
      request<ArticleDetailResponse>(`/api/article/detail?aType=${aType}&enid=${encodeURIComponent(enid)}`),
  },
}
