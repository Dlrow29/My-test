# Log CLI

并发日志文件处理工具，支持多文件同时扫描和ERROR日志提取。

## 功能特点

- 🚀 多文件并发处理
- 🔍 ERROR 日志自动提取
- ⚡ sync.WaitGroup 并发控制
- 📊 处理统计信息

## 使用方法

```bash
# 编译
go build -o log-cli main.go

# 处理多个日志文件
./log-cli -files=app1.log,app2.log,app3.log

# 处理单个文件
./log-cli -files=error.log
```

## 参数说明

- `-files`: 逗号分隔的日志文件路径列表

## 输出

扫描所有指定文件，输出包含"ERROR"的日志行，并显示处理统计信息。
