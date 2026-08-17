# dedao-dl

> 🦉 《得到》 APP 课程下载工具，扫码或者使用 cookie 登录后，可在终端查看已购买的课程，听书书架，电子书架，锦囊，推荐话题等

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/yann0917/dedao-dl)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/yann0917/dedao-dl)

欢迎体验桌面客户端 [dedao-gui](https://github.com/yann0917/dedao-gui)

## 特别声明

仅供个人学习使用，请尊重版权，内容版权均为得到所有，请勿传播内容！！！

仅供个人学习使用，请尊重版权，内容版权均为得到所有，请勿传播内容！！！

仅供个人学习使用，请尊重版权，内容版权均为得到所有，请勿传播内容！！！

## 特性

* 可查看**购买**的课程，课程文章内容
* 可查看听书书架，电子书架列表
* 可查看已购买的锦囊
* 可查看知识城邦推荐话题精选内容
* 课程可生成PDF，文稿生成 Markdown 文档，也可生成 mp3 文件
* 每天听本书可下载音频，文稿生成 pdf、Markdown 文档
* 电子书可下载 html, pdf, epub
* 电子书读书笔记可导出为 markdown
* 可切换登录账号

## 安装

### 安装依赖

`dedao-dl` 支持markdown文本下载，pdf下载，以及音频下载，请按照自己的下载需求，安装下列依赖：

#### pdf下载

* wkhtmltopdf
  > 课程和电子书转 PDF 需要借助[wkhtmltopdf](https://wkhtmltopdf.org/downloads.html)

#### 音频下载

* ffmpeg
  > 音频需要借助 [ffmpeg](https://ffmpeg.org/) 合成

#### markdown文本下载

不需要额外安装依赖

### 使用二进制文件安装

进入[下载列表](https://github.com/yann0917/dedao-dl/releases),下载对应的系统版本，下载后即可使用。

### 使用 `go` 安装

安装go，版本需大于1.23，并设置GOPATH环境变量, 并在PATH中添加$GOPATH/bin

使用如下命令安装：

`go install github.com/yann0917/dedao-dl@latest`

### 使用 Docker 运行

> 为了加快 build 速度，`alpine` 镜像源已修改为阿里镜像。(docker 内没有安装 wkhtmltopdf 不能下载PDF)

如果不想在本地安装 `ffmpeg` 则提供了 `docker` 环境，参考以下命令构建并使用容器执行相关命令。

```bash
# build
docker build https://github.com/yann0917/dedao-dl.git#main -t dedao

# 登录
docker run -v `pwd`/config.json:/app/config.json -it --rm dedao login -q

docker run -v `pwd`/config.json:/app/config.json -it --rm dedao cat
# 查看课程
docker run -v `pwd`/config.json:/app/config.json -it --rm dedao course

# 查看电子书
docker run -v `pwd`/config.json:/app/config.json -it --rm dedao ebook

# 下载课程
docker run -v `pwd`/output:/app/output -v `pwd`/config.json:/app/config.json -it --rm dedao dl xxx
# 下载每天听本书
docker run -v `pwd`/output:/app/output -v `pwd`/config.json:/app/config.json -it --rm dedao dlo xxx

```

#### 删除 `docker images` 中为 `none` 的镜像

```bash
docker ps -a | grep "Exited" | awk '{print $1 }'|xargs docker stop
docker ps -a | grep "Exited" | awk '{print $1 }'|xargs docker rm
docker images|grep none|awk '{print $3 }'|xargs docker rmi
```

## 使用方法

* 使用 `dedao-dl login -q` **同时支持「得到App」和「微信」扫码扫码登录**，或者电脑端登录 [得到](https://www.dedao.cn) 生成 cookie 使用 `dedao-dl login -c "xxxxxxxx"` 登录 ☆☆☆☆☆

`dedao-dl -h` 可查看帮助说明，每个命令都有 `-h` 参数可查看该命令的用法

```bash
Usage:
  dedao-dl [command]


Available Commands:
  ace         获取我的锦囊
  article     获取文章详情
  cat         获取课程分类
  clean       清理 output 或 .cache 目录
  channel     学习圈相关操作
  course      获取我购买过课程
  dl          下载已购买课程，并转换成 PDF & 音频
  dle         下载电子书
  dlo         下载每天听本书音频 & 文稿
  ebook       获取我的电子书架
  free        获取免费专区课程
  help        Help about any command
  login       登录得到 pc 端 https://www.dedao.cn
  odob        获取我的听书书架
  recent      查询用户最近学习情况
  su          切换登录过的账号
  topic       获取推荐话题列表
  users       查看登录过的用户列表
  web         启动 Web UI 与 API 服务
  who         查看当前登录的用户
```

### Web UI 使用说明

项目内置了 Web UI（前端页面）和 Web API（后端 gin 服务），可通过 `web` 命令启动：

```bash
# 默认启动地址：http://127.0.0.1:17878 ，并自动打开浏览器
dedao-dl web

# 指定监听地址/端口（适合局域网访问或端口冲突时）
dedao-dl web --host 0.0.0.0 --port 17878

# 仅启动服务，不自动打开浏览器
dedao-dl web --open=false
```

说明：

* Web 服务读取当前目录下的 `config.json`（与 CLI 共用配置与登录信息）
* 无需提前在 CLI 登录：打开 Web 页面后即可使用二维码登录；若已在 CLI 登录过，也可直接进入
* 退出：在终端按 `Ctrl + C`，服务会进行优雅关闭

Web 页面功能（随版本迭代可能略有调整）：

* 扫码登录：打开页面即可扫码登录
* 学习工作台：统一查看已购课程 / 听书 / 电子书 / 锦囊（含分组）与基础信息
* 内容详情：课程详情与文章列表、听书详情与文稿入口、电子书详情与书评
* 下载导出：一键发起下载任务（MP3 / PDF / Markdown / HTML / EPUB）
* 下载进度：在页面底部查看下载队列与进度
* 得到榜单：查看课程 / 听书 / 图书等榜单内容

`dedao-dl cat` 获取课程分类

```text
+---+----------+------+----------+
| # |   名称   |  统计 |  分类标签  |
+---+----------+------+----------+
| 0 | 全部     | 1696 | all      |
| 1 | 课程     |   64 | bauhinia |
| 2 | 听书书架  | 1407 | odob     |
| 3 | 电子书架  |  210 | ebook    |
| 4 | 锦囊      |   15 | compass  |
+---+----------+------+----------+
```

`dedao-dl course` 获取我购买过课程

```text
+----+-----+-------------------------+--------------------------------+---------------------+--------+----------+
| #  | ID  |          课程名称        |                作者             |      购买日期       |  价格  | 学习进度 |
+----+-----+-------------------------+--------------------------------+---------------------+--------+----------+
|  0 |  51 | 张潇雨·个人投资课          | 张潇雨·对冲基金管理人            | 2020-08-17 09:59:56 | 199.00 |      100 |
|  1 | 560 | 吴军·硅谷来信³（年度日更） | 吴军·计算机科学家                 | 2020-11-05 09:28:23 | 299.00 |       14 |
|  2 | 571 | 年度得到·何帆中国经济报告  | 何帆·著名经济学者                 | 2020-12-10 07:22:28 |  99.00 |       36 |
|  3 |  21 | 古典·超级个体              | 古典·资深生涯规划师             | 2019-09-23 17:11:02 | 199.00 |       84 |
|  4 |  79 | 老喻的人生算法课           | 老喻·思考者                     | 2019-09-23 17:13:51 |  99.00 |      100 |
|  5 |  84 | 给忙碌者的女性健康课       | 常青·协和医学院博士               | 2020-10-27 21:30:55 |  19.90 |      100 |
|  6 |  16 | 陈海贤·自我发展心理学      | 陈海贤·心理学博士                 | 2019-09-23 17:13:46 |  99.00 |      100 |
|  7 | 544 | 听书番外篇                | 阿狮·每天听本书负责人             | 2020-08-18 15:43:44 |   0.00 |      100 |
|  8 | 530 | 跟李松蔚学心理咨询         | 李松蔚·心理咨询师                 | 2020-07-01 11:01:44 |  99.00 |      100 |
|  9 | 484 | 王立铭·抑郁症医学课        | 王立铭·浙江大学生命科学研究院教授    | 2020-10-27 13:15:36 |  39.90 |      100 |
| 10 |  82 | 陈海贤·亲密关系30讲        | 陈海贤·心理学博士                 | 2019-11-05 00:02:21 |  99.00 |      100 |
+----+-----+-------------------------+---------------------------------+---------------------+--------+----------+
```

`course/odob/ebook` 支持分页与排序参数：

```bash
dedao-dl course --page 1 --limit 18
dedao-dl course --order buy --page 1 --limit 18
dedao-dl course --group-id 12345 --page 1 --limit 18
dedao-dl odob --page 1 --limit 18
dedao-dl odob --group-id 12345 --page 1 --limit 18
dedao-dl ebook --page 1 --limit 18
dedao-dl ebook --group-id 12345 --page 1 --limit 18
```

参数规则说明：

* `--page` 和 `--limit` 需要同时传；都不传时保持原逻辑（自动拉取全部）
* 分页模式（同时传 `--page` 和 `--limit`）下不展开分组，只展示当前页原始列表
* `course --order` 支持 `study`（默认）和 `buy`（最近购买）
* `odob --order`、`ebook --order` 仅支持 `study`

`dedao-dl recent` 查询最近学习情况（默认使用当前登录用户 uid_hazy）

```bash
dedao-dl recent -h
dedao-dl recent
dedao-dl recent --page-size 20 --max-id 0
dedao-dl recent --product-type 66 --filter-product-type=true
dedao-dl recent --uid-hazy <uid_hazy>
dedao-dl --json recent
```

参数说明：

* `--uid-hazy` 默认自动读取当前登录用户，也可显式指定
* `--page-size` 每页数量，默认 `20`
* `--max-id` 分页游标，默认 `0`
* `--product-type` 产品类型过滤（如 `66`）
* `--filter-product-type` 是否按 `product_type` 过滤，默认 `true`

`dedao-dl clean` 清理工作目录下的输出或缓存目录

```bash
dedao-dl clean output
dedao-dl clean cache
```

说明：

* `clean output` 会清空并重建 `output/` 目录
* `clean cache` 会先关闭 BadgerDB，再清空并重建 `.cache/` 目录

`dedao-dl free` 获取免费课程列表

```
┌────┬────────────────────────────────┬─────────────────────────────┬───────┬──────────────────────────────────────────────────────────┐
│ #  │          ID ( ENID )           │            NAME             │ SCORE │                          INTRO                           │
├────┼────────────────────────────────┼─────────────────────────────┼───────┼──────────────────────────────────────────────────────────┤
│ 0  │ nb9L2q1e3OxKBPNsdoJrgN8P0Rwo6B │ 得到头条                    │ 4.50  │ 拥抱变化，学点知识                                       │
│ 1  │ mlEA1baQN7WKeRBspoV8L9OgkryvYw │ 文明之旅                    │ 4.90  │ 一档持续20年人文历史类长视频节目，每周三更新             │
│ 2  │ ElLD8OrepAxVvGMs4kJ2oybGdmBnvM │ 长谈                        │ 4.90  │ 深度对谈栏目，每周六0点更新。以一灯传诸灯 终至万灯皆明！ │
│ 3  │ b0rNAzaYOj7VyPMs09K8P54m6wlk12 │ 得到精选                    │ 4.10  │ 每天为你精选好课程，创造与知识的偶遇                     │
│ 4  │ LZ1RgB0EW3NK0wjsLbXkP7vj68pDeA │ 脱不花·怎样成为高效学习的人 │ 4.90  │ 让每个人都能从知识中获得力量                             │
│ 5  │ 5L9DznlwYyOVdwasGdKmbWABv0Zk4a │ 罗辑思维·启发俱乐部         │ 4.40  │ 和你一起终身学习                                         │
│ 6  │ QG0xjrRlNamKYdesjdXg5k293pyOqM │ 跟钱炜练习冥想              │ 4.90  │ 为内心赋能                                               │
│ 7  │ QG0xjrRlNamKYdesGdXg5k293pyOqM │ 李彦宏·智能交通7讲          │ 4.80  │ 100分钟带你探索未来出行                                  │
│ 8  │ 7Dl2p3wZn89JA5RhNrJ0LqayRA6EG1 │ 如何玩转智能手机            │ 4.80  │ 得到给咱爸妈的知识礼物                                   │
│ 9  │ 93N5e6Rya4ZJjQYsmDVOmGApwlo08D │ 邵恒头条                    │ 4.90  │ 把你的知识世界往前推进一步                               │
│ 10 │ qBr4kj5gLNYKNdPsm2JdemExy36GaA │ 听书番外篇                  │ 4.80  │ 听书番外篇，知识有意思                                   │
└────┴────────────────────────────────┴─────────────────────────────┴───────┴──────────────────────────────────────────────────────────┘
```

`dedao-dl course -i xxx` 查看课程信息

```text
专栏名称：张潇雨·个人投资课
专栏作者：张潇雨·对冲基金管理人
更新进度：60/60
课程亮点：1. 这是一门适用于普通投资者的课程。即使你没有专业技巧，没有额外的时间精力，也可以通过课程找到适合自己的投资策略与工具，建立自己的投资体系。

2. 这门课程能帮你重建正确的投资常识，从源头上杜绝投资失误。个人投资本来是一条宽敞平坦的康庄大道。只是有太多错误的岔路需要避免。课程的前14讲都在给你指出岔路，带你回到正途。

3. 投资的难点永远在于实操。课程中有16讲实操内容，帮你解决不知道投资产品怎么选、选什么的问题，以及当极端情况出现时你该怎么办。

4. 张潇雨老师在投资银行、私募基金、对冲基金、家族办公室都曾任职，可以说管理过各种各样的钱。所以，他能告诉你为什么很多看似高深的投资方法恰恰不能用于个人投资。

5. 这一次，知识真的就是财富。

+---+------+----------------------------+------+---------------------+--------------+
| # |  ID  |            章节            | 讲数 |      更新时间       | 是否更新完成 |
+---+------+----------------------------+------+---------------------+--------------+
| 0 | 1040 | 导论(3讲)                  |    3 | 2019-05-05 23:52:31 | ✔            |
| 1 | 1041 | 模块一：市场规律(6讲)      |    6 | 2019-06-02 00:26:29 | ✔            |
| 2 | 1042 | 模块二：投资工具(6讲)      |    6 | 2019-06-02 00:26:36 | ✔            |
| 3 | 1043 | 模块三：自我局限(6讲)      |    6 | 2019-06-02 00:26:43 | ✔            |
| 4 | 1044 | 模块四：投资组合构建(9讲)  |    9 | 2019-06-02 00:26:50 | ✔            |
| 5 | 1503 | 模块五：投资实战(16讲)     |   16 | 2020-09-14 14:46:30 | ✔            |
| 6 | 1045 | 结语(1讲)                  |    1 | 2020-09-14 15:07:06 | ✔            |
| 7 | 1046 | 附录(2讲)                  |    2 | 2020-09-14 15:07:00 | ✔            |
| 8 | 1502 | 加餐丨货币的基础逻辑(11讲) |   11 | 2020-09-14 14:48:57 | ✔            |
+---+------+----------------------------+------+---------------------+--------------+
```

`dedao-dl article -i xxx` 查看课程文章信息

```markdown
+----+-------+----------------------------------------------------------+---------------------+----------+
| #  |  ID   |                         课程名称                         |      更新时间       | 是否阅读 |
+----+-------+----------------------------------------------------------+---------------------+----------+
|  0 | 86033 | 发刊词丨这一次，知识就是财富                             | 2019-05-05 23:44:19 | ✔        |
|  1 | 86037 | 01丨普通投资者的优势                                     | 2019-05-05 23:44:33 | ✔        |
|  2 | 86038 | 02丨普通投资者的劣势                                     | 2019-05-05 23:44:57 | ✔        |
|  3 | 86040 | 03丨多元分散：房子还会不会是表现最好的资产？             | 2019-05-05 23:45:24 | ✔        |
|  4 | 86051 | 04丨择时陷阱：“低买高卖”可以实现么？                     | 2019-05-07 00:00:00 | ✔        |
|  5 | 86074 | 05丨宏观迷信：自下而上的投资逻辑                         | 2019-05-08 00:00:01 | ✔        |
|  6 | 86080 | 06丨风险度量：怎样真正地理解风险和亏损？                 | 2019-05-09 00:00:02 | ✔        |
|  7 | 86081 | 07丨海外配置：人人都可以买下全球市场                     | 2019-05-10 00:00:00 | ✔        |
|  8 | 86082 | 08丨第一模块问答：该何时卖出一只股票？                   | 2019-05-11 00:00:01 | ✔        |
|  9 | 86083 | 09丨指数基金：买到伟大公司的最好机会                     | 2019-05-12 00:00:00 | ✔        |
+----+-------+----------------------------------------------------------+---------------------+----------+
```

`dedao-dl dl 123 -t 1 -m -c -o` 下载课程ID 123 的所有课程

* -t 下载格式, 1:mp3, 2:PDF文档, 3:markdown文档 (default 1)
* -m 是否合并课程内容（针对markdown文档），默认不合并
* -c 是否下载热门留言（针对markdown文档），默认不下载
* -o 是否按顺序展示, 如果为true, 则文件名前缀会加上序号, 如 `00x.`

注意：生成 PDF 的时候，操作过于频繁会触发 `496 NoCertificate` , 因此每次生成一次PDF sleep 0~5秒, 尽管如此，还是有极大可能触发操作频繁图形验证。

`dedao-dl dle 123 -t 1` 下载电子书，先通过 `dedao-dl ebook` 获取要下载的电子书 id,  下载格式, 1:html, 2:PDF文档, 3:epub (default 1)

`dedao-dl dlo 123 -t 1` 下载听书ID 123 的音频或文稿, 先通过 `dedao-dl odob` 获取要下载的听书 id, -t 下载格式, 1:mp3, 2:PDF文档, 3:markdown文档 (default 1)

`./dedao-dl ebook notes -i 158162` 查看电子书id = xxx 的读书笔记, 先通过 `dedao-dl ebook` 获取要下载的电子书 id

`./dedao-dl ebook 158162 -t4` 下载电子书id = xxx 的读书笔记, 先通过 `dedao-dl ebook` 获取要下载的电子书 id，-t4 表示下载 markdown 格式的读书笔记

## Skills 使用说明

仓库内置两个面向 agent 的技能说明文件，位于 `skills/` 目录：

* `skills/dedao-dl-commands/SKILL.md`：纯命令速查，适合“命令怎么写、参数怎么传、给我可复制命令”的场景
* `skills/dedao-dl-usage/SKILL.md`：完整用法与排障，适合“从登录到下载流程”和“报错排查”的场景

推荐使用方式：

* 只要命令：优先使用 `dedao-dl-commands`
* 要步骤和排查：使用 `dedao-dl-usage`
* 面向 agent 自动化时，默认使用 JSON 输出：`dedao-dl --json <command> ...`
* 不确定参数时，先执行：`dedao-dl <command> -h`

## References

* [geektime-dl](https://github.com/mmzou/geektime-dl)
* [annie](https://github.com/iawia002/annie)

## Buy me a coffee ☕️

<html>
    <table style="margin-left: auto; margin-right: auto;">
        <tr>
            <td>
               <img src="docs/img/mm_facetoface_collect_qrcode_1678971248686.png"></img>
            </td>
            <td>
                <img src="docs/img/mm_1678972065469.png"></img>
            </td>
        </tr>
    </table>
</html>

![follow me](/docs/img/scan_search_green.png)

## Star History

<a href="https://www.star-history.com/?repos=yann0917%2Fdedao-dl&type=date&legend=top-left">
 <picture>
   <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=yann0917/dedao-dl&type=date&theme=dark&legend=top-left&sealed_token=FcuGWUCiv7tHfnB-5rVV_z0MrtdIHt3_GHVuK7gkPQuZB5cYK36v8Zu8WGfVssCJaX2Pb0Rg8qW2VzBK" />
   <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=yann0917/dedao-dl&type=date&legend=top-left&sealed_token=FcuGWUCiv7tHfnB-5rVV_z0MrtdIHt3_GHVuK7gkPQuZB5cYK36v8Zu8WGfVssCJaX2Pb0Rg8qW2VzBK" />
   <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=yann0917/dedao-dl&type=date&legend=top-left&sealed_token=FcuGWUCiv7tHfnB-5rVV_z0MrtdIHt3_GHVuK7gkPQuZB5cYK36v8Zu8WGfVssCJaX2Pb0Rg8qW2VzBK" />
 </picture>
</a>

## License

[MIT](./LICENSE) © yann0917

## Support

[![jetbrains](https://s1.ax1x.com/2020/03/26/G9uQoR.png)](https://www.jetbrains.com/?from=dedao-dl)

---
