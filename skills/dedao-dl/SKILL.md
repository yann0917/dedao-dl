---
name: dedao-dl
description: dedao-dl 命令行工具的全量用法与排障：登录得到账号、搜索/浏览课程·听书·电子书、下载 MP3/PDF/Markdown/EPUB/读书笔记、学习圈、最近学习、多账号切换。凡用户提到 dedao-dl、"得到"下载或导出、贴出 dedao.cn 链接想下载内容、询问某条命令参数写法、或 dedao-dl 命令报错时都应触发——即使用户没有明确说"dedao-dl"这个词。
---

# dedao-dl 使用助手

`dedao-dl` 是得到（dedao.cn）的命令行下载器：下载已购课程（MP3/PDF/Markdown）、每天听本书（MP3/PDF/Markdown）、电子书（HTML/PDF/EPUB/读书笔记）。产物默认写在当前目录 `output/` 下，缓存位于 `.cache/`。

## 前置依赖

- PDF 导出需要 `wkhtmltopdf` 在 PATH 中
- 音频处理需要 `ffmpeg` 可执行
- Docker 容器内不含 `wkhtmltopdf`，PDF 下载不可用

## 核心工作流

```bash
# 1) 登录（二选一）
dedao-dl login -q                 # 扫码登录
dedao-dl login -c "<cookie>"      # cookie 来自 https://www.dedao.cn

# 2) 验证登录
dedao-dl who

# 3) 浏览 / 搜索，拿到 ID
dedao-dl course                   # 已购课程
dedao-dl odob                     # 听书书架
dedao-dl ebook                    # 电子书架
dedao-dl search --query "基层中国的运行逻辑" --type 0
dedao-dl free                     # 免费专区

# 4) 下载
dedao-dl dl  <courseID|courseEnid> -t 1     # 课程
dedao-dl dlo <odobID|topic_id_str> -t 1     # 听书
dedao-dl dle <ebookID|ebookEnid>   -t 1     # 电子书
```

## Agent 使用规则

- 面向 agent / 脚本处理时，所有命令加全局 `--json`：`dedao-dl --json <command> ...`
- 不确定参数时先执行 `dedao-dl <command> -h` 自查
- **ID 有两类，别混用**：
  - 数字 ID（courseID/ebookID…）：必须先拉列表建立映射——课程 `course`、电子书 `ebook`、听书 `odob`
  - 字符串 enid：可直接用，来自 URL 的 `id` 参数
- **search 结果继续操作时，enid 取 `list[].list[].extra.enid`**，不要用 `id/goods_id`；再按类型分发：
  - `track_name=ebook` 或 `goods_type=2` → `dle <extra.enid>`
  - `track_name=storytell` 或 `goods_type=13` → `dlo <extra.enid>`
  - `goods_type=66`（课程类）→ `dl <extra.enid>`
- **URL 自动识别**（用户提供 dedao.cn 链接时提取 `id` 参数）：
  - `/course/detail?id=<课程enid>` → `dl <课程enid>`
  - `/course/article?id=<文章enid>` → **不能直接 `dl <文章enid>`**（实测报 `104000 服务异常`：文章 enid 不是课程 enid）。需先取所属课程 enid 与文章数字 ID，再 `dl <课程enid> <文章数字ID>`：
    ```bash
    # POST /pc/bauhinia/pc/article/info   body: {"detail_id":"<文章enid>"}
    # → c.article_info.class_enid（课程 enid）、c.article_info.id（文章数字 ID）
    dedao-dl dl <class_enid> -t 1 <article_info.id>
    ```
    注意：`article --articleEnID` 只输出正文，**不含 `class_enid`**，无法用它反推所属课程
  - `/audioBook/detail?id=<id>` → `dlo <id>`
  - `/ebook/reader?id=<id>` 或 `/ebook/detail?id=<id>` → `dle <id>`
- 下载格式 `-t`：
  - `dl` / `dlo`：1=mp3（默认） 2=PDF 3=markdown
  - `dle`：1=html（默认） 2=PDF 3=epub 4=markdown 笔记
- `dl` 专属：`-m` 合并章节文稿、`-c` 下载热门留言（仅 markdown）、`-o` 文件名加序号前缀；`dl` 还可追加第 2 个位置参数只下载单篇文章：`dl <课程enid> <articleID>`
  - `articleID` 是**课程文章列表里的数字 id**（如 `115529`），不是 URL 里的 enid，也不是文章详情里的 `dd_article_id` 雪花值（如 `1880700885751824556`）。用 `dedao-dl --json article -c <课程enid>` 列出
  - 传错不会退化成整课程下载：过滤条件会命中 0 篇，结果为空
- 列表分页：`course/odob/ebook` 的 `--page` 与 `--limit` 必须同时传；都不传则自动拉全量。`--order` 仅 course 支持 `study|buy`，odob/ebook 仅 `study`

## 命令地图

| 意图 | 命令 |
|---|---|
| 登录/账号 | `login -q`、`login -c <cookie>`、`who`、`user`、`users`、`su <uid>` |
| 搜索/浏览 | `search -q <关键词>`、`cat`（分类）、`course`、`odob`、`ebook`、`free`、`ace`（锦囊）、`topic` |
| 详情 | `course -i <id>`、`ebook -i <id>`、`article --id <courseID>`、`article --articleEnID <enid>`、`free <enid>` |
| 下载 | `dl`（课程）、`dlo`（听书）、`dle`（电子书） |
| 笔记 | `ebook notes -i <ebookID>`、`dle <id> -t 4` |
| 学习圈 | `channel info`、`homepage`、`vip --id <channelID>` |
| 学习记录/VIP | `recent`、`vip-ebook`、`vip-odob` |
| 维护 | `clean output`、`clean cache`、`web`（Web UI，默认 127.0.0.1:17878） |

## 完整参考

- 每个子命令的全部 flags、位置参数与示例：**读 [references/commands.md](references/commands.md)**
- 报错排查（登录失败、参数/ID 错误、依赖缺失、频率限制）：**读 [references/troubleshooting.md](references/troubleshooting.md)**

## 注意事项

- 排查前先复述用户目标（看列表/详情/下载/笔记/报错），再给最短可执行命令，之后才补可选参数
- PDF 批量转换可能触发频率限制（如 496），降低并发重试
- 部分内容需要购买或 VIP 权限，报权限类错误先确认账号资产
- `web` 命令会常驻进程（默认 127.0.0.1:17878，`--open=false` 不自动开浏览器）；在无人值守/用户无终端的环境里，用完要主动结束该进程
