# Task CLI

一个简单的任务管理命令行工具，使用 Go 语言开发。

## 功能特点

- ✅ 任务的增删改查
- 📝 支持任务状态管理 (todo/in-progress/done)
- 💾 JSON 文件持久化存储
- ⏰ 自动记录创建和更新时间

## 使用方法

```bash
# 编译
go build -o task-cli main.go

# 添加任务
./task-cli add "完成项目文档"

# 查看所有任务
./task-cli list

# 更新任务状态
./task-cli update 1 done

# 删除任务
./task-cli delete 1
```

## 任务状态

- `todo` - 待办
- `in-progress` - 进行中  
- `done` - 已完成
