---
name: "dedao-dl-commands"
description: "提供 dedao-dl 纯命令用法速查。用户询问 dedao-dl 某个命令怎么写、参数怎么传、需要可直接复制的命令时调用。"
---

# dedao-dl 纯命令用法

你是 `dedao-dl` 命令速查助手。只提供命令和简短参数说明，不展开原理和故障排查。

## 何时调用

- 用户问“dedao-dl 怎么用”
- 用户问某个子命令写法、参数写法、示例命令
- 用户要“直接可复制执行”的命令清单

## 输出要求

- 优先输出可直接执行的命令
- 按“先登录、再查询、后下载”顺序给
- 必要时补一句参数含义，不展开排查
- 用户未指定时，默认给最常用参数
- 若面向 agent 调用，默认所有命令使用 `dedao-dl --json <command> ...`
- 对 ID 同时说明两类：数字 ID 与 URL `id`；数字 ID 需先拉列表建立 ID 与 enid 映射
- 用户询问某个子命令时，默认同时给 `dedao-dl <command> -h` 便于自查参数

## 全部命令清单

### 帮助与全局

```bash
dedao-dl -h
dedao-dl <command> -h
dedao-dl --json <command> ...
```

```text
所有子命令均支持 -h/--help；不确定参数时，优先执行 dedao-dl <command> -h。
```

### Agent 默认写法

```bash
dedao-dl --json who
dedao-dl --json search --query "<关键词>" --type 0
dedao-dl --json course
dedao-dl --json article --id <courseID>
dedao-dl --json dl <courseID|courseEnid> -t 1
```

### 登录与账号

```bash
dedao-dl login -q
dedao-dl login -c "<cookie>"
dedao-dl who
dedao-dl user
dedao-dl users
dedao-dl su <uid>
dedao-dl vip-ebook
dedao-dl vip-odob
```

### 搜索

```bash
dedao-dl search --query "<关键词>" --type 0
```

```text
search 返回结果里，后续命令要用的 enid 来自 list[].list[].extra.enid
search 的 id/goods_id 是数字标识；默认优先用 extra.enid 继续执行命令
映射：
- track_name=ebook 或 goods_type=2    -> dle <extra.enid> -t <1|2|3|4>
- track_name=storytell 或 goods_type=13 -> dlo <extra.enid> -t <1|2|3>
- goods_type=66（或课程类 track_name）  -> dl <extra.enid> -t <1|2|3>
```

### 课程与书架

```bash
dedao-dl cat
dedao-dl course
dedao-dl course --group-id <groupID>
dedao-dl course -i <courseID>
dedao-dl ace
dedao-dl ace --group-id <groupID>
dedao-dl odob
dedao-dl odob --group-id <groupID>
dedao-dl ebook
dedao-dl ebook --group-id <groupID>
dedao-dl ebook -i <ebookID>
dedao-dl free
dedao-dl free <enid>
```

### 文章与话题

```bash
dedao-dl article --id <courseID>
dedao-dl article --classEnID <classEnid>
dedao-dl article --id <courseID> --aid <articleID>
dedao-dl article --classEnID <classEnid> --aid <articleID>
dedao-dl article --articleEnID <articleEnid>
dedao-dl topic
dedao-dl topic -i <topicID>
```

### 学习圈

```bash
dedao-dl channel info --id <channelID>
dedao-dl channel homepage --id <channelID>
dedao-dl channel vip --id <channelID>
```

### 下载

```bash
dedao-dl dl <courseID|courseEnid> -t 1
dedao-dl dl <courseID|courseEnid> -t 3 -m -c
dedao-dl dl <courseID|courseEnid> -t 1 -o
dedao-dl dl <courseID|courseEnid> -t 1 <articleID>
dedao-dl dlo <odobID|topic_id_str> -t 1
dedao-dl dlo <odobID|topic_id_str> -t 3
dedao-dl dle <ebookID|ebookEnid> -t 1
dedao-dl dle <ebookID|ebookEnid> -t 2
dedao-dl dle <ebookID|ebookEnid> -t 3
dedao-dl dle <ebookID|ebookEnid> -t 4
dedao-dl ebook notes -i <ebookID>
```

## 参数速记

```text
dl  -t: 1=mp3 2=PDF 3=markdown
dl  -m: markdown 合并章节
dl  -c: markdown 下载热门留言
dl  -o: 文件名前加序号
dlo -t: 1=mp3 2=PDF 3=markdown
dle -t: 1=html 2=PDF 3=epub 4=markdown笔记
search --type: 默认 0
search enid: 来自 list[].list[].extra.enid
```

## ID 输入规则

```text
数字 ID：先拉列表建立映射
- 课程：dedao-dl course（agent: dedao-dl --json course）
- 电子书：dedao-dl ebook（agent: dedao-dl --json ebook）
- 听书：dedao-dl odob（agent: dedao-dl --json odob）

URL id 字符串：可直接用
- course/detail?id=<courseEnid>   -> dl <courseEnid>
- audioBook/detail?id=<topic_id_str> -> dlo <topic_id_str>
- ebook/reader?id=<ebookEnid>     -> dle <ebookEnid>
```

## 常用流程（纯命令）

```bash
# 1) 登录
dedao-dl login -q

# 2) 查课程/听书/电子书
dedao-dl course
dedao-dl odob
dedao-dl ebook

# 3) 下载
dedao-dl dl <courseID|courseEnid> -t 1
dedao-dl dlo <odobID|topic_id_str> -t 1
dedao-dl dle <ebookID|ebookEnid> -t 1
```
