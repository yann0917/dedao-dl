# dedao-dl 排障指南

排查原则：先复述用户目标与关键报错，再按「登录 → 参数/ID → 依赖 → 下载/权限」顺序定位，每步给可执行命令。

## 登录与鉴权

- 报错「请先前往 https://www.dedao.cn 登录得到账户」：先 `dedao-dl login -q` 或 `dedao-dl login -c "<cookie>"`
- 扫码后仍未生效：`dedao-dl who` 验证当前账号；必要时 `dedao-dl users` + `dedao-dl su <uid>` 切换
- cookie 登录失败：确认 cookie 来自 `https://www.dedao.cn` 且未过期
- 多账号结果「不对」：可能当前活跃账号不是预期账号，用 `who` 确认、`su` 切换

## 参数与 ID

- 「参数错误」「文章ID错误」：多半把 enid 当数字 ID（或反之）传入。数字 ID 先拉列表建映射（`course`/`ebook`/`odob`），拿不到映射就用 URL 的 `id` 字符串
- `article --aid` 必须配合课程 ID 或课程 enid；只有文章 enid 时用 `--articleEnID`
- `dlo` 推荐直接用音频 URL 的 `id`（topic_id_str），不必依赖书架命中
- `course/odob/ebook --page --limit` 只传一个不生效，必须成对
- search 结果里拿错字段：正确路径是 `list[].list[].extra.enid`，不是 `id/goods_id`

## 依赖与环境

- PDF 失败：检查 `wkhtmltopdf` 是否已安装并在 PATH（`which wkhtmltopdf`）
- 音频处理失败：检查 `ffmpeg` 是否可执行
- Docker 容器：不含 `wkhtmltopdf`，PDF 下载不可用，只出 mp3/markdown
- Web UI 起不来或占用端口：默认 127.0.0.1:17878，可用 `--port` 换端口；脚本/无终端环境加 `--open=false`，用完结束进程

## 下载与内容

- PDF 批量转换触发频率限制（如 496）：降低并发、稍后重试
- 权限类错误：确认账号已购买该内容或具备 VIP 权限（`vip-ebook`/`vip-odob` 查看）
- 下载产物位置：当前目录 `output/`；空间不足或想重来：`dedao-dl clean output`、`dedao-dl clean cache`

## 输出与显示

- 终端表格过宽：换大窗口，或用 `--json` 自行格式化
- agent/脚本处理：一律 `dedao-dl --json <command>`，字段路径参见命令参考对应说明
