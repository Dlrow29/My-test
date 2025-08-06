# Downloader

一个简单高效的多线程文件下载器，使用 Go 语言开发。

## 功能特点

- 🚀 支持多线程并发下载，显著提升下载速度
- 📁 支持自定义输出文件名和路径  
- ⚙️ 可配置并发线程数量
- 💻 简洁的命令行界面

## 安装与编译

确保已安装 Go 环境（版本 1.16+），然后执行：

```bash
go build -o downloader main.go downloader.go
```

## 使用方法

### 命令行参数

| 参数 | 说明 | 是否必填 | 默认值 |
|------|------|----------|--------|
| `-u` | 下载链接 URL | 是 | - |
| `-o` | 输出文件名 | 是 | - |
| `-n` | 并发线程数 | 否 | CPU 核心数 |

### 使用示例

```bash
# 基本使用
./downloader -u https://example.com/file.zip -o file.zip

# 指定并发线程数
./downloader -u https://example.com/largefile.zip -o largefile.zip -n 8

# 下载到指定目录
./downloader -u https://example.com/document.pdf -o /path/to/document.pdf -n 4
```

## 许可证

MIT License
