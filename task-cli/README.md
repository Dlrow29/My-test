# Task-CLI

一个用 Go 编写的极简命令行任务管理器，支持增删改查与状态标记。来源：[https://roadmap.sh/projects/task-tracker](https://roadmap.sh/projects/task-tracker)


## 安装 & 运行

```bash
git clone https://github.com/Dlrow29/go-lab.git
cd go-lab/task-cli
go run . add "Buy milk"
```

## 命令速查

| 命令 | 说明 |
|---|---|
| `go run . add <description>` | 新建任务 |
| `go run . update <id> <new>` | 修改描述 |
| `go run . delete <id>` | 删除任务 |
| `go run . mark-in-progress <id>` | 设为进行中 |
| `go run . mark-done <id>` | 设为已完成 |
| `go run . list`          | 列出全部任务                          |
| `go run . list <status>` | 按状态过滤：todo / in-progress / done |


## 数据文件

任务保存在当前目录的 tasks.json，纯文本，可随时查看。


