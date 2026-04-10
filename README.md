# aimux API Gateway

<!-- Keep these links. Translations will automatically update with the README. -->
[中文](https://github.com/xanygo/aimux) |
[Deutsch](https://zdoc.app/de/xanygo/aimux) | 
[English](https://zdoc.app/en/xanygo/aimux) | 
[Español](https://zdoc.app/es/xanygo/aimux) | 
[français](https://zdoc.app/fr/xanygo/aimux) | 
[日本語](https://zdoc.app/ja/xanygo/aimux) | 
[한국어](https://zdoc.app/ko/xanygo/aimux) | 
[Português](https://zdoc.app/pt/xanygo/aimux) | 
[Русский](https://zdoc.app/ru/xanygo/aimux) 


管理外部依赖的 API 以及其 APIKey。
数据加密存储，避免在配置文件中或者环境变量中存储原始明文的 APIKey。

使用 `aimux` 后，在 `openclaw.json` 或者环境变量中存储的是 `aimux` 分配的 APIKey，即使泄露也不会有任何影响。

![功能介绍](https://github.com/user-attachments/assets/8b015708-2d24-4f02-aed0-dec31e6c3f17)

## 1. 特性：
 
- **易用**：提供 Web 表单管理页面
- **跨平台**：使用 Go 开发，支持 Windows、Linux、Macos 等系统
- **安全**：数据加密存储
- **存储**：数据支持存储在本地文件系统、Redis
- **实用**：同一个 API，支持多个下游按照权重轮询使用、支持多个模型(model)按照权重轮询使用
- **可观察**：支持 RPC Dump，可以将请求和响应数据存储到日志目录(已脱敏)，以方便分析


## 2. 安装：

### 2.1 使用 docker-compose
下载 [docker-compose.yml](./docker-compose.yml),
修改其中的账号、密码等，在同目录下创建 data 和 log 目录之后，使用
```
docker compose up
```
启动运行。

之后可以通过 `http://127.0.0.1:8201/admin/` 访问管理页面。
API 地址的前缀则为 `http://127.0.0.1:8200/` 。

### 2.2 使用 go install
```
go install github.com/xanygo/aimux@master
```
将 [app.yml](./conf/app.yml) 放到 `/home/work/aimux/conf` 目录中，并修改账号密码等。
```bash
cd /home/work/aimux
aimux
```

### 2.2 下载二进制
在 [releases 页面](https://github.com/xanygo/aimux/releases) 下载编译好的二进制。配置运行同上。

