# dedao-dl 全命令参考

由 `dedao-dl <command> -h` 的真实输出整理。全局 flag：`--json`（以 JSON 输出结果）；所有子命令支持 `-h/--help`。位置参数规则以源码为准：`dl` 1-2 个参数，`su` 恰好 1 个，`free` 0-1 个。

## login — 登录得到 PC 端

```text
用法: dedao-dl login [flags]
  -c, --cookie string   cookie from https://www.dedao.cn
  -q, --qrcode          扫码登录 https://www.dedao.cn
```

示例：`dedao-dl login -q`、`dedao-dl login -c "<cookie>"`

## who / user / users / su — 账号

```text
dedao-dl who            当前登录用户
dedao-dl user           用户信息（默认摘要字段；--json 看完整字段 avatar/vip_user 等）
dedao-dl users          登录过的用户列表
dedao-dl su <uid>       切换账号（位置参数恰好 1 个，切换后自动执行 who）
```

## search — 搜索

```text
用法: dedao-dl search [flags]
  -q, --query string   搜索关键词
  -t, --type int       搜索类型，默认 0
```

示例：`dedao-dl search --query "基层中国的运行逻辑" --type 0`

结果字段：后续命令用的 enid 在 `list[].list[].extra.enid`；`id/goods_id` 是数字标识，优先用 `extra.enid`。

## course / odob / ebook / ace — 书架与列表

四个命令的 flags 结构一致：

```text
  -g, --group-id int   分组ID，显示指定分组内的内容
  -i, --id int         详情 ID（ace=锦囊、course=课程、odob=听书、ebook=电子书）
  -l, --limit int      每页数量（与 --page 一起使用）
  -p, --page int       页码（与 --limit 一起使用）
      --order string   course: study(默认)|buy；odob/ebook: 仅 study（默认）
```

- `--page` 与 `--limit` 必须同时传；都不传则自动拉全量
- 传了分页时不展开分组，仅返回当前页原始列表

示例：`dedao-dl course`、`dedao-dl course -i 12345`、`dedao-dl ebook --group-id 12345 --page 1 --limit 20`

## ebook notes — 电子书笔记列表

```text
用法: dedao-dl ebook notes [flags]
  -i, --id int   电子书ID
```

示例：`dedao-dl ebook notes -i 12456`

## article — 文章详情 / 列表

```text
用法: dedao-dl article [flags]
  -a, --aid int              文章id
  -c, --classEnID string     课程enid
  -e, --articleEnID string   文章enid
  -i, --id int               课程id
```

组合：`--id` 或 `--classEnID` 列文章列表；再配 `--aid` 看单篇正文（Markdown 输出）；只有文章 enid 时直接 `--articleEnID`。

示例：`dedao-dl article --i 12345`、`dedao-dl article --i 12345 --a 67890`、`dedao-dl article --e yyyyy`

## dl — 下载已购课程

```text
用法: dedao-dl dl <courseID|courseEnid> [articleID]   （位置参数 1-2 个）
  -t, --downloadType int   1:mp3（默认） 2:PDF文档 3:markdown文档
  -m, --merge              合并课程章节文稿（仅 markdown）
  -c, --comment            下载课程热门留言（仅 markdown）
  -o, --order              文件名前缀加序号（00x.）
```

enid 即课程详情链接的 `id` 参数，如 `https://www.dedao.cn/course/detail?id=ZWyMAOLnR4xJ1vqse8X65QaE8YG29k`。

示例：`dedao-dl dl 123 -t 1 -m`、`dedao-dl dl ZWy... -t 3 -m -c`、`dedao-dl dl ZWy... -t 1 67890`（只下载单篇）

## dlo — 下载每天听本书

```text
用法: dedao-dl dlo <odobID|topic_id_str>
  -t, --downloadType int   1:mp3（默认） 2:PDF文档 3:markdown文档
```

建议优先用音频详情链接的 `id`（topic_id_str）：`https://www.dedao.cn/audioBook/detail?id=ZV1po7jlgBdObqZXy0GNR4wLzxya5v`

示例：`dedao-dl dlo ZV1po7jlgBdObqZXy0GNR4wLzxya5v -t 1`

## dle — 下载电子书

```text
用法: dedao-dl dle <ebookID|ebookEnid>
  -t, --downloadType int   1:html（默认） 2:PDF文档 3:epub 4:markdown笔记
```

enid 即阅读链接的 `id` 参数：`https://www.dedao.cn/ebook/reader?id=N5lDqb9b47...`

示例：`dedao-dl dle 123 -t 1`、`dedao-dl dle N5lDqb9b47... -t 3`

## free — 免费专区

```text
用法: dedao-dl free [enid]   （位置参数 0-1 个）
不带参数 = 免费课程列表；带 enid = 该课程详情
```

## cat — 课程分类

```text
dedao-dl cat    课程分类标签及数量统计
```

## topic — 推荐话题

```text
用法: dedao-dl topic [flags]
  -i, --id string   话题id（不带 = 话题列表；带 = 该话题前 40 条精选内容）
```

## channel — 学习圈

```text
dedao-dl channel info --id <channelID>        学习圈信息
dedao-dl channel homepage --id <channelID>    首页分类
dedao-dl channel vip --id <channelID>         VIP/权限信息
```

## recent — 最近学习

```text
用法: dedao-dl recent [flags]
      --filter-product-type   是否按 product_type 过滤（默认 true）
      --max-id int            分页游标，默认 0
      --page-size int         每页数量，默认 20
      --product-type string   产品类型过滤（如 66=电子书；默认不过滤）
      --uid-hazy string       用户 uid_hazy（默认自动读当前登录用户）
```

示例：`dedao-dl recent --page-size 50 --product-type 66`

## vip-ebook / vip-odob — VIP 信息

```text
dedao-dl vip-ebook    电子书 VIP（配额、是否会员/过期等；--json 看完整字段）
dedao-dl vip-odob     每天听本书 VIP（卡片名称/价格/订阅状态等）
```

## clean — 清理目录

```text
dedao-dl clean output   清理 output 目录
dedao-dl clean cache    清理 .cache 目录
```

## web — Web UI 与 API

```text
用法: dedao-dl web [flags]
      --host string   监听地址（默认 127.0.0.1）
      --open          启动后自动打开浏览器（默认 true；脚本环境用 --open=false）
      --port int      监听端口（默认 17878）
```

## completion / help

cobra 自动生成：`dedao-dl completion [bash|zsh|fish|powershell]`、`dedao-dl help [command]`。
